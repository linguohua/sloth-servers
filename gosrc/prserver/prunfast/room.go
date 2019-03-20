package prunfast

import (
	"fmt"
	"gconst"
	"gpubsub"
	"gscfg"
	"math/rand"
	"pokerface"
	"sort"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

var (
	chairs2P = []int{0, 2}
	chairs3P = []int{0, 1, 2}
	chairs4P = []int{0, 1, 2, 3}

	// SystemFloat2ProtoScale 分数放大倍数，主要是proto不方便存储浮点数
	// SystemFloat2ProtoScale float32 = 10.0
)

// Room 一个房间（等价于一个牌桌）
// 房间最关键的数据是状态机，房间任意时刻都有一个状态机来表达其状态
// 以及处理用户的消息和事件
// 房间的状态机会根据触发条件切换：例如刚开始的时候房间的IDLE状态
// 只要有一个用户进来，就切换到Waiting状态
// 而在Waiting状态下，如果所有玩家都已经离线并且player删除，则转为IDLE状态
// 注意waiting和playing随着一手牌的开始和结束是反复切换的，跟逻辑上的牌局开始（即是开始了第一局就认为牌局开始）
// 是不一样的，状态机主要是用于和客户端消息交互的控制，逻辑上的牌局开始则是通过handRoundFinished
// 以及handRoundStarted判断
type Room struct {
	isUlimitRound     bool   // 无限局数，主要用于客户端长时间测试
	isForMonkey       bool   // 是否用于测试
	bankerUserID      string // 庄家的ID
	bankerSwitchCount int    // 庄家切换次数，用于切换风圈（风圈切换了后pseudoFlowerCardID也会跟着变化）
	qaIndex           int    // 流水号，注意每一手牌开始流水号都会重置
	markup            int    // 连续荒庄计数

	// isKongFollowLocked bool // 新疆麻将特有，房间是否处于锁杠状态

	config    *RoomConfig // 房间的配置对象
	configID  string
	monkeyCfg *MonkeyCfg // 房间发牌配置

	players        []*PlayerHolder // 房间的所有玩家
	playingUserIDs []string        // 游戏进行时，玩家ID列表，如果服务器奔溃恢复，该列表从redis加载，见RoomMgr的restoreRooms

	state IState // 状态机

	ownerID    string // 房间拥有者，就是谁开的房，扣钱的时候用，解散的时候只有房主才可以解散
	chairs     []int  // 房间的座位配置，例如2人房间，座位配置为0,2；也即是相对而坐
	ID         string // 房间的唯一ID，由房间管理服务器根据其规则而定，只需确保唯一即可
	roomNumber string // 房间号
	clubID     string // 俱乐部ID

	handRoundFinished int // 已经完成的手牌轮数
	handRoundStarted  int // 已经开始的手牌轮数，可能和handRoundFinished相等，也可能大1

	disband          *RoomDisband // 解散控制数据结构，当玩家申请解散时生成，解散状态完结后删除
	lastReceivedTime time.Time    // 最后一次消息接收时间，用于判断房间空闲了多长时间

	scoreRecords []*pokerface.MsgRoomHandScoreRecord

	deleteReason pokerface.RoomDeleteReason

	rand *rand.Rand

	rbl                *RoomBlockList
	networkRequestLock *sync.Mutex // 网络请求处理lock，确保任意时刻，房间只处理一个玩家请求

	cl *logrus.Entry
}

// newBaseRoom 新建room对象并做一些基本初始化
func newBaseRoom(ownerID string, clubID string, ID string, roomNumber string) *Room {
	r := &Room{}
	r.ownerID = ownerID
	r.players = make([]*PlayerHolder, 0, 4)
	r.ID = ID
	r.roomNumber = roomNumber
	r.clubID = clubID
	r.lastReceivedTime = time.Now()
	r.playingUserIDs = make([]string, 0, 4)
	r.scoreRecords = make([]*pokerface.MsgRoomHandScoreRecord, 0, 16)
	r.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	r.networkRequestLock = &sync.Mutex{}

	r.cl = logrus.WithFields(logrus.Fields{
		"roomID": r.ID,
		"number": r.roomNumber,
	})

	// 进入空闲状态
	r.state2(nil, pokerface.RoomState_SRoomIdle)

	return r
}

// initChair 初始化座位列表，座位从0开始编号
func (r *Room) initChair() {
	r.chairs = make([]int, 0, 4)
	playerNum := r.config.playerNumAcquired
	switch playerNum {
	case 2:
		r.chairs = append(r.chairs, chairs2P...)
	case 3:
		r.chairs = append(r.chairs, chairs3P...)
	case 4:
		r.chairs = append(r.chairs, chairs4P...)
	}
}

// allocChair 申请一个座位, fixChairID如果大于-1表示指定了座位
func (r *Room) allocChair(fixChairID int) int {
	if len(r.chairs) == 0 {
		r.cl.Panic("no chair id to alloc")
		return -1
	}

	var result = -1
	if fixChairID >= 0 {
		for i, c := range r.chairs {
			if c == fixChairID {
				result = fixChairID
				r.chairs = append(r.chairs[0:i], r.chairs[i+1:]...)
			}
		}
	}

	if result < 0 {
		result = r.chairs[0]
		r.chairs = r.chairs[1:]
	}

	return result
}

// releaseChair 归还一个座位
func (r *Room) releaseChair(chairID int) {
	if len(r.chairs) == r.config.playerNumAcquired {
		r.cl.Panic("chair array is fulled")
		return
	}

	r.chairs = append(r.chairs, chairID)
	// 排序座位
	sort.Ints(r.chairs)
}

// newRoomFromMgr 房间管理服务器请求创建房间
func newRoomFromMgr(ownerID string, clubID string, ID string, configID string,
	configJSON *RoomConfigJSON, roomNuber string) *Room {
	r := newBaseRoom(ownerID, clubID, ID, roomNuber)
	r.isForMonkey = false

	r.config = newRoomConfigFromJSON(configJSON)
	r.initChair()
	r.configID = configID

	return r
}

// newRoomForMonkey 测试管理器请求创建一个房间
func newRoomForMonkey(ownerID string, ID string, roomConfig *RoomConfig) *Room {
	r := newBaseRoom(ownerID, "", ID, ID)
	r.isForMonkey = true
	// 测试用
	r.isUlimitRound = true
	r.config = roomConfig
	r.initChair()
	r.configID = ""
	r.roomNumber = ID
	return r
}

// qaIndexExpected 是否期待的QAIndex
func (r *Room) qaIndexExpected(qaIndex int) bool {
	return r.qaIndex == qaIndex
}

// bankerPlayer 当前的庄家
func (r *Room) bankerPlayer() *PlayerHolder {
	for _, p := range r.players {
		if p.userID() == r.bankerUserID {
			return p
		}
	}

	return nil
}

// getPlayerByUserID 根据userID获取player对象
func (r *Room) getPlayerByUserID(userID string) *PlayerHolder {
	for _, p := range r.players {
		if p.userID() == userID {
			return p
		}
	}

	return nil
}

// getPlayerByChairID 根据chairID获取player对象
func (r *Room) getPlayerByChairID(chairID int) *PlayerHolder {
	for _, p := range r.players {
		if p.chairID == chairID {
			return p
		}
	}

	return nil
}

// nextQAIndex 下一个qaIndex
func (r *Room) nextQAIndex() int {
	r.qaIndex++
	return r.qaIndex
}

// bankerChange2 更改庄家到新玩家
func (r *Room) bankerChange2(newBanker *PlayerHolder) {
	if newBanker.userID() == r.bankerUserID {
		return
	}

	r.cl.Printf("banker change, old:%s, new:%s\n", r.bankerUserID, newBanker.userID())
	r.bankerUserID = newBanker.userID()

	// 计算风圈
	// 从庄家开始逆时针，轮庄，首次定庄后，依次从东风→南风→西风→北风开始，东风为花牌。
	// 如果首次庄家连庄则下局依然为东风为花牌，庄家下庄下家为庄，也是东风为花牌，
	// 轮庄满4人后，则庄的花牌从南风开始，依次类推。
	r.bankerSwitchCount++
}

// onHandOver 一手牌结局
func (r *Room) onHandOver(msgHandOver *pokerface.MsgHandOver) {

	r.handRoundFinished++

	if r.isUlimitRound {
		r.handRoundFinished = 0
	}

	// 写一些玩家统计信息到redis
	r.writePlayersStatis()

	if r.handRoundFinished == r.config.handNum {
		// 已经达到最大局数，需要续费
		r.onGameOver(msgHandOver)

		roomMgr.forceDisbandRoom(r, pokerface.RoomDeleteReason_DisbandMaxHand)
		return
	}

	// 重置qaIndex
	r.qaIndex = 0

	// 下一手牌，所以直接进入等待状态而不是空闲状态
	r.state2(r.state, pokerface.RoomState_SRoomWaiting)

	for _, p := range r.players {
		p.resetForNextHand()
		p.state = pokerface.PlayerState_PSNone

		// 确保状态已经切换到SWaiting后，才发送手牌结果给客户端
		p.sendHandOver(msgHandOver)
	}

	// 所有用户状态已经被改为PlayerState_PSNone
	// 因此通知所有客户端更新用户状态
	r.updateRoomInfo2All()

	// 写分数记录到redis，以便可以奔溃时恢复
	r.writeHandEnd2Redis()
}

// onGameOver 处理牌局结束，通知所有玩家牌局结束
func (r *Room) onGameOver(msgHandOver *pokerface.MsgHandOver) {
	var msg = serializeMsgGameOver(r)

	for _, p := range r.players {
		if msgHandOver != nil {
			p.sendHandOver(msgHandOver)
		}

		p.sendMsg(msg, int32(pokerface.MessageCode_OPGameOver))
	}

	// 如果是俱乐部房间，则需要写大赢家统计
	if r.clubID != "" && len(r.scoreRecords) > 0 {
		r.calcSaveGreatWinnersForClubRoom()
	}
}

// isForceConsistent 房间是否是测试房间并强制一致
func (r *Room) isForceConsistent() bool {
	return r.monkeyCfg != nil && r.monkeyCfg.isForceConsistent
}

// handBegin 房间做一些准备开始游戏
func (r *Room) handBegin() {
	if r.isForceConsistent() {
		// 重设置一下庄家ID
		r.bankerUserID = r.monkeyCfg.monkeyUserCardsCfgList[0].userID
		// r.markup = r.monkeyCfg.markup

		if r.monkeyCfg.isContinuousBanker {
			r.bankerPlayer().gStatis.isContinuousBanker = true
		}
	} else {
		// if nil == r.monkeyCfg && r.handRoundStarted == 0 {
		// 	// 第一局随机庄家
		// 	n := len(r.players)
		// 	x := r.rand.Intn(n)
		// 	log.Println("random banker:", x)
		// 	r.bankerUserID = r.players[x].userID()
		// }
	}

	r.refreshBankerID()

	// 记录计数器
	r.handRoundStarted++
	if r.isUlimitRound {
		r.handRoundStarted = 1
	}

	r.playingUserIDs = r.playingUserIDs[:0]
	for _, p := range r.players {
		r.playingUserIDs = append(r.playingUserIDs, p.userID())
	}

	// 每一局开始都重置锁杠状态
	// r.isKongFollowLocked = false

	r.writeHandBegin2Redis()
}

// refreshBankerID 修正庄家ID，由于庄家可以退出了游戏等
// 需要选择新的庄家
func (r *Room) refreshBankerID() {
	if r.bankerUserID == "" {
		if len(r.players) > 0 {
			r.bankerUserID = r.players[0].userID()
		}
	} else if r.bankerPlayer() == nil {
		if len(r.playingUserIDs) < 1 {
			if len(r.players) > 0 {
				r.bankerUserID = r.players[0].userID()
			}
		}
	}
}

// onUserMessage 处理websocket消息
func (r *Room) onUserMessage(user IUser, msg []byte) {
	r.networkRequestLock.Lock()
	defer r.networkRequestLock.Unlock()

	gmsg := &pokerface.GameMessage{}
	err := proto.Unmarshal(msg, gmsg)
	if err != nil {
		r.cl.Println("onUserMessage, unmarshal error:", err)
		return
	}

	// 记录一下最后一个消息的接收时间
	r.lastReceivedTime = time.Now()

	// 对于玩家请求离开，直接处理，不需要交给状态机
	var msgCode = gmsg.GetOps()
	var handled = false
	switch msgCode {
	case int32(pokerface.MessageCode_OPPlayerLeaveRoom):
		player := r.getPlayerByUserID(user.userID())
		player.allowedLeave = r.allow2Leave(player)

		if !player.allowedLeave {
			msg2 := &pokerface.MsgEnterRoomResult{}
			var status32 = int32(1)
			msg2.Status = &status32
			buf := formatGameMsg(msg2, int32(pokerface.MessageCode_OPPlayerLeaveRoom))

			user.send(buf)
			return
		}

		// 发送结果给客户端
		buf := formatGameMsg(nil, int32(pokerface.MessageCode_OPPlayerLeaveRoom))
		user.send(buf)
		// 断开websocket连接
		user.closeWebsocket()
		handled = true
		break
	case int32(pokerface.MessageCode_OPDisbandRequest):
		r.onDisbandRequest(user, gmsg)
		handled = true
		break
	case int32(pokerface.MessageCode_OPDisbandAnswer):
		r.onDisbandAnswer(user, gmsg)
		handled = true
		break
	case int32(pokerface.MessageCode_OPKickout):
		if r.handRoundStarted > 0 {
			player := r.getPlayerByUserID(user.userID())
			sendKickoutError(player, pokerface.KickoutResult_KickoutResult_FailedGameHasStartted)
			handled = true
		}
		break
	case int32(pokerface.MessageCode_OPDonate):
		r.onDonateRequest(user, gmsg)
		handled = true
		break
	case int32(pokerface.MessageCode_OP2Lobby):
		r.onUser2Lobby(user, gmsg)
		handled = true
		break
	default:
		break
	}

	if handled {
		return
	}

	// 不是房间可以处理的消息，交给状态机
	r.state.onMessage(user, gmsg)
}

// onUser2Lobby 处理玩家返回大厅请求，相当于掉线
func (r *Room) onUser2Lobby(user IUser, gmsg *pokerface.GameMessage) {
	r.cl.Printf("user %s require 2 lobby", user.userID())

	if r.handRoundStarted > 0 {
		// 游戏已经开始，不可以切换
		msg2 := &pokerface.MsgEnterRoomResult{}
		var status32 = int32(1)
		msg2.Status = &status32
		buf := formatGameMsg(msg2, int32(pokerface.MessageCode_OP2Lobby))

		user.send(buf)
		return
	}

	// 发送结果给客户端
	buf := formatGameMsg(nil, int32(pokerface.MessageCode_OP2Lobby))
	user.send(buf)
	// 断开websocket连接
	user.closeWebsocket()
}

// 踢开所有正在房间内的玩家
func (r *Room) kickAll() {
	// 断开玩家的链接
	for _, p := range r.players {
		// 需要调用closeWebsocket而不是detach
		// 否则座位没有被归还，房间无法再次进入
		p.user.closeWebsocket()
	}
}

// 踢开所有正在房间内的玩家
func (r *Room) reset() {
	// r.handRoundFinished = 0
	r.handRoundStarted = r.handRoundFinished
	r.qaIndex = 0

	if r.state.getStateConst() == pokerface.RoomState_SRoomIdle {
		return
	}

	// 首先转到空闲状态以便在玩家断线时删除玩家
	if r.state.getStateConst() == pokerface.RoomState_SRoomPlaying {
		r.state2(r.state, pokerface.RoomState_SRoomWaiting)
	}

	// 特殊处理已经离线的玩家
	offlinePlayers := make([]*PlayerHolder, 0, 4)
	for _, p := range r.players {
		if p.state == pokerface.PlayerState_PSOffline {
			offlinePlayers = append(offlinePlayers, p)
		}
	}
	// 已经离线的玩家需要单独移除，因为其已经没有websocket
	// 没法通过websocket的关闭事件来移除
	for _, p := range offlinePlayers {
		r.state.onPlayerLeave(p)
	}

	// 断开玩家的链接
	for _, p := range r.players {
		// 需要调用closeWebsocket而不是detach
		// 否则座位没有被归还，房间无法再次进入
		p.user.closeWebsocket()
	}
}

// forceDisband 房间解散
// 流程：首先请求房间服务器，解散房间
// 房间服务器答复可以解散后，调用destroy()
// 真正销毁房间
func (r *Room) forceDisband() {
	r.cl.Println("forceDisband room")
	// 先请求房间服务器删除这个房间
	if r.requireRoomServer2Delete() == false {
		r.cl.Println("requireRoomServer2Delete failed, but force to delete room")
	}

	// 销毁房间
	r.destroy()
}

// requireRoomServer2Delete 发送删除房间的请求给房间管理服务器
// 并等待回复
func (r *Room) requireRoomServer2Delete() bool {

	msgDeleteRoom := &gconst.SSMsgGameServer2RoomMgrServerDisbandRoom{}
	msgDeleteRoom.RoomID = &r.ID
	var handFinished = int32(r.handRoundFinished)
	var handStart = int32(r.handRoundStarted)
	msgDeleteRoom.HandFinished = &handFinished
	msgDeleteRoom.HandStart = &handStart

	msgDeleteRoom.PlayerUserIDs = r.playingUserIDs

	msgDeleteRoomBuf, err := proto.Marshal(msgDeleteRoom)
	if err != nil {
		r.cl.Println("requireRoomServer2Delete parse roomConfig err： ", err)
		return false
	}

	msgType := int32(gconst.SSMsgType_Request)
	requestCode := int32(gconst.SSMsgReqCode_DeleteRoom)
	status := int32(gconst.SSMsgError_ErrSuccess)

	msgBag := &gconst.SSMsgBag{}
	msgBag.MsgType = &msgType
	var sn = (pubSubSequnce)
	pubSubSequnce++

	msgBag.SeqNO = &sn
	msgBag.RequestCode = &requestCode
	msgBag.Status = &status
	var url = gscfg.ServerID
	msgBag.SourceURL = &url
	msgBag.Params = msgDeleteRoomBuf

	//等待游戏服务器的回应
	succeed, _ := gpubsub.SendAndWait(gscfg.RoomServerID, msgBag, 10*time.Second)

	if succeed {
		r.cl.Println("room server reply ok for delete room require")
	} else {
		r.cl.Println("wait room server reply timeout, for delete request")
	}

	return succeed
}

// destroy 房间删除处理
func (r *Room) destroy() {
	reasonStr, ok := pokerface.RoomDeleteReason_name[int32(r.deleteReason)]
	if ok {
		r.cl.Printf("room now destroy, room number, reason:%s\n", reasonStr)
	} else {
		r.cl.Printf("room now destroy, room number, reason:%d\n", r.deleteReason)
	}

	// 通知所有人房间已经被删除
	msgDelete := &pokerface.MsgRoomDelete{}
	var reason32 = int32(r.deleteReason)
	msgDelete.Reason = &reason32

	for _, p := range r.players {
		p.sendMsg(msgDelete, int32(pokerface.MessageCode_OPRoomDeleted))
	}

	// 强制停止游戏，状态转换到deleted状态
	r.state2(r.state, pokerface.RoomState_SRoomDeleted)

	// 更新redis数据库
	// 获取redis链接，并退出函数时释放
	conn := pool.Get()
	defer conn.Close()

	// 关闭所有玩家连接
	for _, p := range r.players {
		// 重设为offline
		p.user.detach()
	}

	conn.Send("MULTI")
	for _, p := range r.players {
		conn.Do("HMSET", gconst.PlayerTablePrefix+p.userID(), "leaveRoom", r.ID, "leaveTime", time.Now().Unix())
	}
	conn.Do("EXEC")

	r.players = nil

	// 删除Game server管理的room记录
	conn.Do("DEL", gconst.GsRoomTablePrefix+r.ID)
}

// onUserOffline 处理用户离线，不同的状态下，玩家离线表现不同
// 例如，如果是等待状态，且游戏并没有开始，那么玩家离线后，其player对象会被清除
// 但是如果是游戏正在进行，那么玩家离线，其player对象不会被清除，而一直等待其上线
// 或者直到其他玩家决定解散本局游戏
func (r *Room) onUserOffline(user IUser, lock bool) {
	if lock {
		r.networkRequestLock.Lock()
		defer r.networkRequestLock.Unlock()
	}

	player := r.getPlayerByUserID(user.userID())
	if player == nil {
		r.cl.Println("user off line, but no player, userID:", user.userID())
		return
	}

	// 让状态机处理用户离线
	// 不同状态下对用户离线的处理是不同的，比如Waiting状态，用户离线会把Player删除
	// 也即是Waiting状态下用户随意进出。但在Playing状态下，用户离线Player对象一直保留
	// 除非其他玩家选择解散本局游戏
	r.state.onPlayerLeave(player)

	// AA需要返还钻石
	if player.allowedLeave && r.config.payType == aapay {
		r.notifyReturnDiamond(player)
	}

	// 如果庄家退出且player被销毁则更新banker id
	r.refreshBankerID()

	// if r.disband != nil && r.disband.applicant != player {
	// 	// 如果正在解散，而用户离线，则认为他已经回复同意
	// 	r.disband.onDisbandAnswer(player, nil)
	// }

	r.writePlayerLeave2Redis(player, player.allowedLeave)
}

func (r *Room) allow2Leave(player *PlayerHolder) bool {
	if r.ownerID == player.userID() {
		return false
	}

	if r.handRoundStarted > 0 {
		return false
	}

	if r.config.payType == aapay && player.userID() == r.ownerID {
		return false
	}

	if r.disband != nil {
		return false
	}

	if r.disband != nil {
		return false
	}

	return true
}

// userTryEnter 处理玩家尝试进入房间
func (r *Room) userTryEnter(ws *websocket.Conn, userID string) IUser {
	r.networkRequestLock.Lock()
	defer r.networkRequestLock.Unlock()

	r.cl.Printf("userTryEnter room, userID:%s\n", userID)
	player := r.getPlayerByUserID(userID)

	// 如果房间是monkey房间，且其配置要求强制一致，则userID必须位于配置中
	if r.isForceConsistent() && player == nil {
		monkeyUserCardCfg := r.monkeyCfg.getMonkeyUserCardsCfg(userID)
		if monkeyUserCardCfg == nil {
			sendEnterRoomError(ws, userID, pokerface.EnterRoomStatus_MonkeyRoomUserIDNotMatch)
			return nil
		}

		// 而且玩家进入的顺序必须严格按照配置指定
		loginSeq := len(r.players)
		if loginSeq != monkeyUserCardCfg.index {
			sendEnterRoomError(ws, userID, pokerface.EnterRoomStatus_MonkeyRoomUserLoginSeqNotMatch)
			return nil
		}
	}

	if player != nil {
		// 用户存在，则可能是如下原因：
		// 		客户端代码判断自己已经离线，然后重连服务器
		//		服务器也知道客户端已经离线并在等客户端上线
		return r.onPlayerReconnect(player, ws)
	}

	if len(r.players) == r.config.playerNumAcquired {
		// 已经满员
		sendEnterRoomError(ws, userID, pokerface.EnterRoomStatus_RoomIsFulled)
		return nil
	}

	// 只有空闲状态，或者等待状态才允许玩家进入
	status := r.state.getStateConst()
	if status != pokerface.RoomState_SRoomIdle && status != pokerface.RoomState_SRoomWaiting {
		sendEnterRoomError(ws, userID, pokerface.EnterRoomStatus_RoomPlaying)
		return nil
	}

	// 检查是否在黑名单里面
	if r.rbl != nil && r.rbl.has(userID) {
		sendEnterRoomError(ws, userID, pokerface.EnterRoomStatus_InRoomBlackList)
		return nil
	}

	// 房间正在解散状态，不允许进入
	if r.disband != nil {
		sendEnterRoomError(ws, userID, pokerface.EnterRoomStatus_RoomInApplicateDisband)
		return nil
	}

	var fixChairID = -1
	// 如果游戏已经开始了至少一局，则不再允许其他人进入，仅允许开局时刻的人进入
	if r.handRoundFinished > 0 {
		found := false
		for i, u := range r.playingUserIDs {
			if u == userID {
				found = true
				playerNum := r.config.playerNumAcquired
				switch playerNum {
				case 2:
					fixChairID = chairs2P[i]
				case 3:
					fixChairID = chairs3P[i]
				case 4:
					fixChairID = chairs4P[i]
				}
				break
			}
		}

		if !found {
			sendEnterRoomError(ws, userID, pokerface.EnterRoomStatus_RoomPlaying)
			return nil
		}
	}

	// 如果是AA，则请求房间管理服务器扣钱
	if r.config.payType == aapay && r.ownerID != userID && r.handRoundFinished == 0 {
		if ok, rt := r.takeOffDiamond(userID); !ok {
			sendEnterRoomError(ws, userID, rt)
			return nil
		}
	}

	// 如果房间是俱乐部房间，则要求进入者必须是俱乐部成员
	if r.clubID != "" {
		if !r.isUserClubMember(userID) {
			sendEnterRoomError(ws, userID, pokerface.EnterRoomStatus_NotClubMember)
			return nil
		}
	}

	// 可以进入房间，新建player对象
	guser := newGUser(userID, ws, r)
	chairID := r.allocChair(fixChairID)
	player = newPlayerHolder(r, chairID, guser)

	// 增加到玩家列表
	r.players = append(r.players, player)

	if len(r.scoreRecords) > 0 {
		r.restorePlayerGStatis(player)
	}

	// 根据座位ID排序
	r.sortPlayers()

	// 发送成功进入房间通知给客户端
	sendEnterRoomError(ws, userID, pokerface.EnterRoomStatus_Success)

	if r.monkeyCfg != nil && player == r.bankerPlayer() {
		if r.monkeyCfg.isContinuousBanker {
			player.gStatis.isContinuousBanker = true
		}
	}

	// 更新banker id
	r.refreshBankerID()

	// 通知状态机，状态机根据当前状态是否需要写redis
	r.state.onPlayerEnter(player)

	// 写redis数据库，以便其他服务器能够知道玩家进入该房间
	r.writePlayerEnter2Redis(player)

	return guser
}

// writePlayerEnter2Redis 把用户进入事件写入redis，包括room当前的玩家列表以及玩家最后处于的room
func (r *Room) writePlayerEnter2Redis(player *PlayerHolder) {
	if r.isForMonkey {
		return
	}

	// 获取redis链接，并退出函数时释放
	conn := pool.Get()
	defer conn.Close()

	r.writeOnlinePlayerList2Redis(conn)

	// 写用户的最后所在房间
	conn.Do("HMSET", gconst.PlayerTablePrefix+player.userID(), "enterRoom", r.ID, "enterTime", time.Now().Unix())
}

// writePlayerLeave2Redis 把用户离开事件写入redis，包括room当前的玩家列表以及玩家最后处于的room
func (r *Room) writePlayerLeave2Redis(player *PlayerHolder, clearLastRoom bool) {
	if r.isForMonkey {
		return
	}

	// 获取redis链接，并退出函数时释放
	conn := pool.Get()
	defer conn.Close()

	r.writeOnlinePlayerList2Redis(conn)

	// 写用户的最后所在房间
	if clearLastRoom {
		conn.Do("HMSET", gconst.PlayerTablePrefix+player.userID(), "leaveRoom", r.ID, "leaveTime", time.Now().Unix())
	}
}

// writePlayerList2Redis 复用一个redis conn，把玩家列表写到redis
func (r *Room) writeOnlinePlayerList2Redis(conn redis.Conn) {

	userIDs := make([]string, 0, len(r.players))
	for _, p := range r.players {
		if p.state != pokerface.PlayerState_PSOffline {
			userIDs = append(userIDs, p.userID())
		}
	}

	// 写房间的用户列表
	var msgIDList = &gconst.SSMsgUserIDList{}
	msgIDList.UserIDs = userIDs
	buf, err := proto.Marshal(msgIDList)
	if err != nil {
		r.cl.Println("writeOnlinePlayerList2Redis error:", err)
	} else {
		conn.Do("HSET", gconst.GsRoomTablePrefix+r.ID, "players", buf)
	}
}

func (r *Room) pushState2RoomMgrServer() {
	userIDs := make([]string, len(r.players))
	for i, p := range r.players {
		userIDs[i] = p.userID()
	}
	// 推送状态更新给房间服务器
	var roomNotify = &gconst.SSMsgRoomStateNotify{}
	var state32 = int32(r.state.getStateConst())
	roomNotify.State = &state32
	var roomID = r.ID
	roomNotify.RoomID = &roomID
	var handStartted32 = int32(r.handRoundStarted)
	roomNotify.HandStartted = &handStartted32
	var time32 = uint32(r.lastReceivedTime.Unix() / 60)
	roomNotify.LastActiveTime = &time32
	roomNotify.UserIDs = userIDs
	pushNotify2RoomServer(gconst.SSMsgReqCode_RoomStateNotify, roomNotify)
}

// onPlayerReconnect 处理玩家重入事件
func (r *Room) onPlayerReconnect(player *PlayerHolder, ws *websocket.Conn) IUser {
	r.cl.Printf("onPlayerReconnect, userID:%s, roomNumber:%s\n", player.userID(), r.roomNumber)
	user := player.user

	// 更新用户的个人信息
	user.updateInfo()

	// 更换websocket连接
	user.rebind(ws)

	// 发送成功进入房间通知给客户端
	sendEnterRoomError(ws, player.userID(), pokerface.EnterRoomStatus_Success)

	// 通知状态机
	r.state.onPlayerReEnter(player)

	// 如果房间处于解散状态，发送解散状态通知
	if r.disband != nil {
		r.disband.sendNotify2All()
	}

	// 写redis数据库，以便其他服务器能够知道玩家进入该房间
	r.writePlayerEnter2Redis(player)

	return user
}

// state2 房间切换状态
func (r *Room) state2(oldState IState, newStateCode pokerface.RoomState) {
	r.cl.Println("room state2----")
	if oldState != nil {
		r.cl.Println("oldState:", oldState.getStateName())
	} else {
		r.cl.Println("oldState nil")
	}

	if oldState != r.state {
		r.cl.Printf("old state:%s, not room's current state:%s, state2 failed\n", oldState.getStateName(), r.state.getStateName())
		return
	}

	var newState IState
	switch newStateCode {
	case pokerface.RoomState_SRoomIdle:
		newState = &SIdle{room: r}
		break
	case pokerface.RoomState_SRoomDeleted:
		newState = &SDeleted{room: r}
		break
	case pokerface.RoomState_SRoomWaiting:
		newState = &SWaiting{room: r}
		break
	case pokerface.RoomState_SRoomPlaying:
		newState = newSPlaying(r)
		break
	}

	if newState == nil {
		r.cl.Println("new state nil, failed")
		return
	}

	r.cl.Println("newState:", newState.getStateName())

	r.state = newState

	if oldState != nil {
		oldState.onStateLeave()
	}

	newState.onStateEnter()
}

// stateRemovePlayer 用于状态机从玩家列表中移除玩家
// 玩家移除后，需要归还座位
func (r *Room) stateRemovePlayer(player *PlayerHolder) {
	for i, p := range r.players {
		if p == player {
			r.players = append(r.players[0:i], r.players[i+1:]...)
			// 归还座位
			r.releaseChair(p.chairID)
			break
		}
	}
}

// updateRoomInfo2All 把房间当前状态和玩家数据发给所有用户
func (r *Room) updateRoomInfo2All() {
	if len(r.players) > 0 {
		var msgRoomInfo = serializeMsgRoomInfo(r)
		for _, p := range r.players {
			p.sendMsg(msgRoomInfo, int32(pokerface.MessageCode_OPRoomUpdate))
		}
	}

	r.pushState2RoomMgrServer()
}

// byChairID 根据座位ID排序
type byChairID []*PlayerHolder

func (s byChairID) Len() int {
	return len(s)
}
func (s byChairID) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byChairID) Less(i, j int) bool {
	return s[i].chairID < s[j].chairID
}

// sortPlayers 根据座位ID排序
func (r *Room) sortPlayers() {
	sort.Sort(byChairID(r.players))
}

// sendEnterRoomError 发送进入房间错误给客户端
func sendEnterRoomError(ws *websocket.Conn, userID string, status pokerface.EnterRoomStatus) {
	if status != pokerface.EnterRoomStatus_Success {
		statusStr, ok := pokerface.EnterRoomStatus_name[int32(status)]
		if ok {
			logrus.Printf("sendEnterRoomError, userID:%s, error:%s\n", userID, statusStr)
		} else {
			logrus.Printf("sendEnterRoomError, userID:%s, error:%d\n", userID, status)
		}
	}

	msg := &pokerface.MsgEnterRoomResult{}
	var status32 = int32(status)
	msg.Status = &status32

	buf := formatGameMsg(msg, int32(pokerface.MessageCode_OPPlayerEnterRoom))

	if ws != nil {
		ws.WriteMessage(websocket.BinaryMessage, buf)
	} else {
		logrus.Println("sendEnterRoomError, ws == nil")
	}
}

// formatGameMsg 构建一个game message类型的消息
func formatGameMsg(pb proto.Message, ops int32) []byte {
	gmsg := &pokerface.GameMessage{}
	gmsg.Ops = &ops

	if pb != nil {
		buf, err := proto.Marshal(pb)
		if err != nil {
			logrus.Println("formatGameMsg err:", err)
			return nil
		}
		gmsg.Data = buf
	}

	bytes, err := proto.Marshal(gmsg)
	if err != nil {
		logrus.Println("marshal game msg failed:", err)
		return nil
	}

	return bytes
}

// onDonateRequest 处理玩家道具请求
func (r *Room) onDonateRequest(user IUser, gmsg *pokerface.GameMessage) {
	var msgDonate = &pokerface.MsgDonate{}
	err := proto.Unmarshal(gmsg.Data, msgDonate)
	if err != nil {
		r.cl.Println("onDonateRequest unmarshal error:", err)
		return
	}

	player := r.getPlayerByUserID(user.userID())
	if player == nil {
		return
	}

	toWho := r.getPlayerByChairID(int(msgDonate.GetToChairID()))
	if toWho == nil {
		return
	}

	// 请求房间管理服务器扣除钻石
	// 如果扣除钻石失败，则发送tip给客户端
	// --player.sendTipsCode(TipCode_TCDonateFailedNoEnoughDiamond)

	// 如果扣除钻石成功，则更新本人的钻石数量和对方的魅力值
	// 通过room update消息通知所有人
	// --r.updateRoomInfo2All()
	ssMsgDonate := &gconst.SSMsgDonate{}
	var propsType = msgDonate.GetItemID()
	ssMsgDonate.PropsType = &propsType
	var from = user.userID()
	ssMsgDonate.From = &from
	var to = toWho.userID()
	ssMsgDonate.To = &to

	ssMsgDonateBuf, err := proto.Marshal(ssMsgDonate)
	if err != nil {
		r.cl.Panicln("Marshal ssMsgDonateBuf err： ", err)
	}

	msgType := int32(gconst.SSMsgType_Request)
	requestCode := int32(gconst.SSMsgReqCode_Donate)
	status := int32(gconst.SSMsgError_ErrSuccess)

	msgBag := &gconst.SSMsgBag{}
	msgBag.MsgType = &msgType
	var sn = (pubSubSequnce)
	pubSubSequnce++

	msgBag.SeqNO = &sn
	msgBag.RequestCode = &requestCode
	msgBag.Status = &status
	var url = gscfg.ServerID
	msgBag.SourceURL = &url
	msgBag.Params = ssMsgDonateBuf

	//等待房间服务器的回应
	success, reply := gpubsub.SendAndWait(gscfg.RoomServerID, msgBag, 5*time.Second)
	if !success {
		r.cl.Println("Donate Waiting room server time out")
		return
	}

	status = reply.GetStatus()

	// 如果扣除钻石失败，则发送tip给客户端
	if status == int32(gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough) {
		player.sendTipsCode(pokerface.TipCode_TCDonateFailedNoEnoughDiamond)
		return
	}

	if status != int32(gconst.SSMsgError_ErrSuccess) {
		r.cl.Println("donate failed, errCode:", reply.GetStatus())
		return
	}

	var params = reply.GetParams()
	var donateRsp = &gconst.SSMsgDonateRsp{}
	err = proto.Unmarshal(params, donateRsp)
	if err != nil {
		r.cl.Println("Unmarshal SSMsgDonateRsp err： ", err)
		return
	}

	// 则更新本人的钻石数量和对方的魅力值
	var fromWhoInfo = user.getInfo()
	var diamond = donateRsp.GetDiamond()
	fromWhoInfo.diamond = int(diamond)

	var toWhoInfo = toWho.user.getInfo()
	var charm = donateRsp.GetCharm()
	toWhoInfo.charm = int(charm)

	r.updateRoomInfo2All()

	// 通知所有人打赏成功
	var fromChairID32 = int32(player.chairID)
	msgDonate.FromChairID = &fromChairID32

	if len(r.players) > 0 {
		for _, p := range r.players {
			p.sendMsg(msgDonate, int32(pokerface.MessageCode_OPDonate))
		}
	}
}

// onDisbandRequest 处理玩家解散房间请求
func (r *Room) onDisbandRequest(user IUser, gmsg *pokerface.GameMessage) {
	r.cl.Printf("room %s disband applicate by user %s\n", r.roomNumber, user.userID())

	var stateConst = r.state.getStateConst()
	if stateConst != pokerface.RoomState_SRoomWaiting &&
		stateConst != pokerface.RoomState_SRoomPlaying {
		// 只有这两个状态下才允许解散房间
		r.cl.Println("room state not allowed to disband, userID:", user.userID())
		return
	}

	if r.handRoundStarted < 1 && r.ownerID != user.userID() {
		// 牌局尚未开始，仅允许房主解散房间
		r.cl.Println("game has not start, only owner can disband")
		r.sendDisbandError(user, pokerface.DisbandState_ErrorNeedOwnerWhenGameNotStart)
		return
	}

	if r.disband != nil {
		// 之前已经有人申请了解散，此时正处于等待解散
		r.cl.Println("room in disband wait reply state")
		r.sendDisbandError(user, pokerface.DisbandState_ErrorDuplicateAcquire)
		return
	}

	var disband = &RoomDisband{}
	disband.applicant = r.getPlayerByUserID(user.userID())

	if len(r.players) > 1 {
		disbandWaitItems := make([]*RoomDisbandWait, 0, len(r.players)-1)

		for _, p := range r.players {
			// if p.state == pokerface.PlayerState_PSOffline {
			// 	// 离线玩家不参与解散投票
			// 	continue
			// }

			if p == disband.applicant {
				continue
			}

			var disbandWait = &RoomDisbandWait{}
			disbandWait.player = p

			disbandWaitItems = append(disbandWaitItems, disbandWait)
		}

		if len(disbandWaitItems) > 0 {
			disband.waitItems = disbandWaitItems
		}
	}

	disband.room = r
	r.disband = disband

	// 启动一个routine执行解散操作
	go disband.startDisband()
}

// sendDisbandError 发送解散请求错误给客户端
func (r *Room) sendDisbandError(user IUser, err pokerface.DisbandState) {
	var player = r.getPlayerByUserID(user.userID())
	msg := &pokerface.MsgDisbandNotify{}
	var state32 = int32(err)
	msg.DisbandState = &state32
	var applicant32 = int32(player.chairID)
	msg.Applicant = &applicant32

	buf := formatGameMsg(msg, int32(pokerface.MessageCode_OPDisbandNotify))
	user.send(buf)
}

// onDisbandAnswer 处理其他玩家对解散房间请求的响应
func (r *Room) onDisbandAnswer(user IUser, gmsg *pokerface.GameMessage) {
	if r.disband == nil {
		r.cl.Println("room no in disband waiting state, discard disband answer")
		return
	}

	var msgDisbandAnswer = &pokerface.MsgDisbandAnswer{}
	err := proto.Unmarshal(gmsg.Data, msgDisbandAnswer)
	if err != nil {
		r.cl.Println(err)
		return
	}

	player := r.getPlayerByUserID(user.userID())

	r.disband.onDisbandAnswer(player, msgDisbandAnswer)
}

// appendHandScoreRecord 增加一手牌得分记录
func (r *Room) appendHandScoreRecord(winType int) {
	if r.isUlimitRound {
		return
	}

	record := &pokerface.MsgRoomHandScoreRecord{}
	var endType32 = int32(winType)
	record.EndType = &endType32
	var handIndex32 = int32(r.handRoundStarted)
	record.HandIndex = &handIndex32

	r.scoreRecords = append(r.scoreRecords, record)

	// if winType == int(HandOverType_enumHandOverType_None) {
	// 	// 如果是流局，则没有玩家得分列表
	// 	return
	// }

	playerRecords := make([]*pokerface.PlayerHandScoreRecord, len(r.players))
	for i, p := range r.players {
		playerRecord := &pokerface.PlayerHandScoreRecord{}
		var userID = p.userID()
		playerRecord.UserID = &userID
		var winType32 = int32(p.sctx.winType)
		playerRecord.WinType = &winType32
		var score32 = int32(p.sctx.calcTotalWinScore())
		playerRecord.Score = &score32

		playerRecords[i] = playerRecord
	}

	record.PlayerRecords = playerRecords
}

// writeHandBegin2Redis 写手牌开始信息到redis
func (r *Room) writeHandBegin2Redis() {

	userIDs := make([]string, len(r.players))
	for i, p := range r.players {
		userIDs[i] = p.userID()
	}

	// 写房间的用户列表
	var msgIDList = &gconst.SSMsgUserIDList{}
	msgIDList.UserIDs = userIDs
	buf, err := proto.Marshal(msgIDList)
	if err != nil {
		r.cl.Println("writeHandScoreRecords2Redis proto error:", err)
		return
	}

	// 更新房间状态到redis
	// 获取redis链接，并退出函数时释放
	conn := pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("HMSET", gconst.GsRoomTablePrefix+r.ID, "hrStartted", r.handRoundStarted, "hrfinished", r.handRoundFinished,
		"bankerID", r.bankerUserID,
		"hp", buf)

	conn.Send("HSET", gconst.GameRoomStatistics+r.ID, "hrStartted", r.handRoundStarted)
	// key 24小时后过期
	conn.Send("EXPIRE", gconst.GameRoomStatistics+r.ID, 24*60*60)

	_, err = conn.Do("EXEC")
	if err != nil {
		r.cl.Println("writeHandBegin2Redis error:", err)
	}

	r.cl.Printf("writeHandBegin2Redis completed, bankerID:%s\n", r.bankerUserID)
}

// writePlayersStatis 写玩家统计信息到redis
func (r *Room) writePlayersStatis() {
	if r.isUlimitRound {
		return
	}

	// 更新房间状态到redis
	// 获取redis链接，并退出函数时释放
	conn := pool.Get()
	defer conn.Close()

	conn.Send("MULTI")

	dd := time.Now().Format("20060102")
	// 过期该key，第二天凌晨5点过期
	tOfTomorrow5 := time.Now().AddDate(0, 0, 1)
	tOfTomorrow5 = time.Date(tOfTomorrow5.Year(), tOfTomorrow5.Month(), tOfTomorrow5.Day(),
		0, 0, 0, 0, time.Local)
	tOfTomorrow5 = tOfTomorrow5.Add(5 * time.Hour)
	expireat := tOfTomorrow5.Unix()

	// 房间管理服务器获取统计信息，其中dfHands为玩家一共进行了多少局游戏，dfHMW为最多的单局得分，dfHML为最大得单局失分
	for _, p := range r.players {
		conn.Send("HINCRBY", gconst.PlayerTablePrefix+p.userID(), "dfHands", 1)

		userID := p.userID()
		// g:yyyymmdd:userID:dsu(roomType)
		dailyStatisTable := fmt.Sprintf(gconst.GameServerDailyStatisTablePrefix, dd, userID, gconst.RoomType_DafengGZ)
		if p.sctx != nil {
			var totalWinScore = p.sctx.calcTotalWinScore()
			if totalWinScore != 0 {
				luaScriptForHandScore.Send(conn, gconst.PlayerTablePrefix+p.userID(), totalWinScore)

				if totalWinScore > 0 {
					// 赢牌次数
					conn.Send("HINCRBY", dailyStatisTable, "wh", 1)
				}
			}
		}

		// 完成局数
		conn.Send("HINCRBY", dailyStatisTable, "fh", 1)
		conn.Send("EXPIREAT", dailyStatisTable, expireat)

		// 创建房间并完成次数
		if r.ownerID == userID && r.handRoundFinished == r.config.handNum {
			conn.Send("HINCRBY", dailyStatisTable, "cf", 1)
		}
	}
	conn.Do("EXEC")
}

// writeHandEnd2Redis 写手牌结束信息到redis
func (r *Room) writeHandEnd2Redis() {
	if r.isUlimitRound {
		return
	}
	roomScoreRecords := &pokerface.RoomScoreRecords{}
	roomScoreRecords.ScoreRecords = r.scoreRecords

	bytes, err := proto.Marshal(roomScoreRecords)
	if err != nil {
		r.cl.Println("writeHandEnd2Redis proto error:", err)
		return
	}

	userIDs := make([]string, len(r.players))
	for i, p := range r.players {
		userIDs[i] = p.userID()
	}

	// 写房间的用户列表
	var msgIDList = &gconst.SSMsgUserIDList{}
	msgIDList.UserIDs = userIDs
	buf, err := proto.Marshal(msgIDList)
	if err != nil {
		r.cl.Println("writeHandEnd2Redis proto error:", err)
		return
	}

	// 更新房间状态到redis
	// 获取redis链接，并退出函数时释放
	conn := pool.Get()
	defer conn.Close()

	conn.Do("HMSET", gconst.GsRoomTablePrefix+r.ID, "hrStartted", r.handRoundStarted, "hrfinished", r.handRoundFinished,
		"bankerID", r.bankerUserID,
		"sr", bytes,
		"hp", buf)

	r.cl.Printf("writeHandEnd2Redis completed, bankerID:%s,scoreRecords:%d\n", r.bankerUserID, len(r.scoreRecords))
}

// readHandInfoFromRedis4Restore 恢复房间时，从redis读取手牌信息
func (r *Room) readHandInfoFromRedis4Restore(conn redis.Conn) {
	// 读取一下已经完成的手牌局数，以及玩家列表等一手牌相关的信息
	values, err := redis.Values(conn.Do("hmget", gconst.GsRoomTablePrefix+r.ID, "hrfinished", "bankerID",
		"sr", "hp"))
	if err == nil {
		handRoundFinished, err := redis.Int(values[0], nil)
		if err == nil {
			r.handRoundFinished = handRoundFinished
			r.handRoundStarted = r.handRoundFinished
		}

		bankerUserID, err := redis.String(values[1], nil)
		if err == nil {
			r.bankerUserID = bankerUserID
			r.cl.Println("resotre bankerID:", r.bankerUserID)
		} else {
			//log.Println(err)
		}

		bytes, err := redis.Bytes(values[2], nil)
		if err == nil {
			roomScoreRecords := &pokerface.RoomScoreRecords{}
			err = proto.Unmarshal(bytes, roomScoreRecords)
			if err == nil {
				r.scoreRecords = roomScoreRecords.ScoreRecords
				r.cl.Println("restore room, scoreRecords:", len(r.scoreRecords))
			}
		} else {
			//log.Println(err)
		}

		bytes, err = redis.Bytes(values[3], nil)
		if err == nil {
			var msgIDList = &gconst.SSMsgUserIDList{}
			err = proto.Unmarshal(bytes, msgIDList)
			if err == nil {
				r.playingUserIDs = append(r.playingUserIDs, msgIDList.UserIDs...)
				r.cl.Println("restore room, hplayers count:", len(r.playingUserIDs))
			}
		} else {
			//log.Println(err)
		}
	}
}

// restorePlayersWhen 如果服务器意外终止服务，重启后直接恢复player对象
func (r *Room) restorePlayersWhen() {

	if len(r.playingUserIDs) < 1 {
		return
	}

	for _, userID := range r.playingUserIDs {
		// 可以进入房间，新建player对象
		guser := newGUser(userID, nil, r)
		chairID := r.allocChair(-1)
		player := newPlayerHolder(r, chairID, guser)

		// 增加到玩家列表
		r.players = append(r.players, player)

		if len(r.scoreRecords) > 0 {
			r.restorePlayerGStatis(player)
		}

		player.state = pokerface.PlayerState_PSOffline
	}

	// 根据座位ID排序
	r.sortPlayers()
	// 更新banker id
	r.refreshBankerID()

	// 转换到等待状态
	r.state2(r.state, pokerface.RoomState_SRoomWaiting)
}

// restorePlayerGStatis 从历史记录中恢复player的全局统计数据
func (r *Room) restorePlayerGStatis(player *PlayerHolder) {
	for _, sr := range r.scoreRecords {
		for _, pr := range sr.PlayerRecords {
			if pr.GetUserID() == player.userID() {
				player.gStatis.roundScore += int(pr.GetScore())

				switch pr.GetWinType() {
				case int32(HandOverType_enumHandOverType_Chucker):
					player.gStatis.greatWinCounter++
					break
				case int32(HandOverType_enumHandOverType_Win_SelfDrawn):
					player.gStatis.winSelfDrawnCounter++
					break
				case int32(HandOverType_enumHandOverType_Win_Chuck):
					player.gStatis.miniWinCounter++
					break
				case int32(HandOverType_enumHandOverType_Win_RobKong):
					player.gStatis.winRobKongCounter++
					break
				case int32(HandOverType_enumHandOverType_Konger):
					player.gStatis.kongerCounter++
					break
				}
				break
			}
		}
	}

	r.cl.Printf("restorePlayerGStatis, chucker:%d, winChuck:%d, winSelf:%d, score:%d\n", player.gStatis.greatWinCounter, player.gStatis.miniWinCounter,
		player.gStatis.winSelfDrawnCounter, player.gStatis.roundScore)
}

// notifyReturnDiamond 通知管理服务器返还钻石给玩家
func (r *Room) notifyReturnDiamond(player *PlayerHolder) {
	r.cl.Println("notifyReturnDiamond, userID:", player.userID())

	if r.isForMonkey {
		return
	}

	msgUpdateBalance := &gconst.SSMsgUpdateBalance{}
	var roomID = r.ID
	msgUpdateBalance.RoomID = &roomID
	var userID = player.userID()
	msgUpdateBalance.UserID = &userID

	pushNotify2RoomServer(gconst.SSMsgReqCode_AAExitRoomNotify, msgUpdateBalance)
}

// takeOffDiamond 扣除钻石
func (r *Room) takeOffDiamond(userID string) (bool, pokerface.EnterRoomStatus) {
	r.cl.Println("takeOffDiamond, userID:", userID)

	if r.isForMonkey {
		return true, pokerface.EnterRoomStatus_Success
	}

	takeoffStatus := pokerface.EnterRoomStatus_TakeoffDiamondFailedNotEnough
	// takeoffStatus = EnterRoomStatus_TakeoffDiamondFailedIO

	msgUpdateBalance := &gconst.SSMsgUpdateBalance{}
	var roomID = r.ID
	msgUpdateBalance.RoomID = &roomID
	msgUpdateBalance.UserID = &userID

	msgUpdateBalanceBuf, err := proto.Marshal(msgUpdateBalance)
	if err != nil {
		r.cl.Panicln("Marshal msgUpdateBalance err： ", err)
	}

	msgType := int32(gconst.SSMsgType_Request)
	requestCode := int32(gconst.SSMsgReqCode_AAEnterRoom)
	status := int32(gconst.SSMsgError_ErrSuccess)

	msgBag := &gconst.SSMsgBag{}
	msgBag.MsgType = &msgType
	var sn = (pubSubSequnce)
	pubSubSequnce++

	msgBag.SeqNO = &sn
	msgBag.RequestCode = &requestCode
	msgBag.Status = &status
	var url = gscfg.ServerID
	msgBag.SourceURL = &url
	msgBag.Params = msgUpdateBalanceBuf

	//等待房间服务器的回应
	succeed, reply := gpubsub.SendAndWait(gscfg.RoomServerID, msgBag, 15*time.Second)
	if !succeed {
		r.cl.Println("wait room server time out, TakeoffDiamondFailedIO")
		return false, pokerface.EnterRoomStatus_TakeoffDiamondFailedIO
	}

	status = reply.GetStatus()
	r.cl.Println("takeOffDiamond status:", status)

	if status == int32(gconst.SSMsgError_ErrTakeoffDiamondFailedIO) {
		r.cl.Println("room server TakeoffDiamondFailedIO")
		return false, pokerface.EnterRoomStatus_TakeoffDiamondFailedIO
	}

	if status == int32(gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough) {
		r.cl.Println("room server TakeoffDiamondFailedNotEnough")
		return false, pokerface.EnterRoomStatus_TakeoffDiamondFailedNotEnough
	}

	if status != int32(gconst.SSMsgError_ErrSuccess) {
		r.cl.Println("takeOffDiamond Unkonw error code:", status)
		return false, takeoffStatus
	}

	return true, takeoffStatus
}

// isContinuAble 房间是否还可以继续打牌
func (r *Room) isContinuAble() bool {
	var isContinueAble = r.handRoundStarted < r.config.handNum
	if !isContinueAble {
		return false
	}

	return true
}

func (r *Room) updateUserLocation(userID string, location string) {
	r.cl.Println("updateUserLocation: ", userID)
	msg := &pokerface.MsgUpdateLocation{}
	msg.UserID = &userID
	msg.Location = &location
	buf := formatGameMsg(msg, int32(pokerface.MessageCode_OPUpdateLocation))

	for _, player := range r.players {
		user := player.user
		user.send(buf)
	}
}

func (r *Room) isUserClubMember(userID string) bool {
	// 获取redis链接，并退出函数时释放
	conn := pool.Get()
	defer conn.Close()

	i, err := redis.Int(conn.Do("SISMEMBER", gconst.ClubMemberSetPrefix+r.clubID, userID))
	if err != nil {
		r.cl.Println("isUserClubMember, redis err:", err)
		return false
	}

	if i == 0 {
		return false
	}

	return true
}

// calcSaveGreatWinnersForClubRoom 计算并保存大赢家，注意函数会增加一次对局计数和大赢家计数，因此必须是新对局完成后调用本函数
func (r *Room) calcSaveGreatWinnersForClubRoom() {
	if r.clubID == "" {
		return
	}

	if len(r.scoreRecords) < 1 {
		return
	}

	// 获取redis链接，并退出函数时释放
	conn := pool.Get()
	defer conn.Close()

	records := r.scoreRecords
	scoresMap := make(map[string]int)
	for _, rc := range records {
		playerScores := rc.PlayerRecords
		for _, ps := range playerScores {
			s := scoresMap[ps.GetUserID()]
			s = s + int(ps.GetScore())

			scoresMap[ps.GetUserID()] = s
		}
	}

	// 找到最大的计分
	maxScore := 0
	for _, v := range scoresMap {
		if v > maxScore {
			maxScore = v
		}
	}

	// 找到获得最大计分的玩家userID
	greatWinnerUserIDs := make([]string, 0, len(scoresMap))
	for k, v := range scoresMap {
		if v == maxScore {
			userID := k
			greatWinnerUserIDs = append(greatWinnerUserIDs, userID)
		}
	}

	// 读取俱乐部统计表
	clubID := r.clubID
	timeNow := time.Now()
	targetDate := timeNow.Format("2006-Jan-02")
	scoreSetKey := "club:" + targetDate + ":" + clubID
	clubScoreMap := make(map[string]int)
	for userID := range scoresMap {
		clubScoreMap[userID] = 0
	}

	index := 0
	userIDs := make([]string, len(scoresMap))
	for uID := range scoresMap {
		userIDs[index] = uID
		index++
	}

	conn.Send("MULTI")
	for _, userID := range userIDs {
		conn.Send("ZSCORE", scoreSetKey, userID)
	}

	values, _ := redis.Values(conn.Do("EXEC"))
	index = 0
	// 每个玩家增加一次对局计数
	for _, userID := range userIDs {
		score, _ := redis.Int(values[index], nil)
		clubScoreMap[userID] = incClubScorePlayCount(score)
		index++
	}

	// 大赢家增加一次计数
	for _, userID := range greatWinnerUserIDs {
		score, _ := clubScoreMap[userID]
		clubScoreMap[userID] = incClubScoreWinCount(score)
	}

	greatWinnersStr := strArray2Comma(greatWinnerUserIDs)
	conn.Send("MULTI")
	for userID, score := range clubScoreMap {
		conn.Send("ZADD", scoreSetKey, score, userID)
	}
	conn.Send("HSET", gconst.MJReplayRoomTablePrefix+r.ID, "gw", greatWinnersStr)
	conn.Do("EXEC")
}

func incClubScorePlayCount(score int) int {
	winCount := int((score & 0xffff0000) >> 16)
	playCount := int((score & 0x0000ffff))

	if playCount == 0 {
		playCount = 0xffff
	}

	playCount = playCount - 1

	score = (winCount << 16) | (playCount)

	return score
}

func incClubScoreWinCount(score int) int {
	winCount := int((score & 0xffff0000) >> 16)
	playCount := int((score & 0x0000ffff))

	winCount = winCount + 1

	score = (winCount << 16) | (playCount)

	return score
}
