package prunfast

import (
	"bytes"
	"container/list"
	fmt "fmt"
	"gconst"
	"pokerface"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"

	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
)

const (
	// MaxShareAbleID 最小分享ID
	MaxShareAbleID = 99999999
	// MinShareAbleID 最大分享ID
	MinShareAbleID = 10000000

	maxReplayRoomNumber = 50
)

// LoopContext 打牌循环上下文
// 主要保存最近动作的玩家
// 录像会保存到redis，一是供玩家回放，二是出问题（例如玩家指出数值计算不正确等）我们可以回放查找bug
// 极端情况下，recorder marshal后的大小接近3K
// 因此保存redis的时候需要考虑保存速度，以及对redis内存的压力
type LoopContext struct {
	//drawCount int

	//actionCount int
	msgReplayRoom       *pokerface.MsgReplayRoom
	replayRecordSummary *pokerface.MsgReplayRecordSummary

	recorder *pokerface.SRMsgHandRecorder

	actionList *list.List

	s  *SPlaying
	cl *logrus.Entry
}

func newLoopContext(s *SPlaying) *LoopContext {
	ctx := &LoopContext{}
	ctx.s = s
	ctx.recorder = &pokerface.SRMsgHandRecorder{}
	ctx.actionList = list.New()

	ctx.cl = s.cl

	return ctx
}

func (lc *LoopContext) fetchNonUserReplyOnly(stepback int) *pokerface.SRAction {
	var step = 0
	for e := lc.actionList.Back(); e != nil; e = e.Prev() {
		sraction := e.Value.(*pokerface.SRAction)
		if sraction.GetFlags()&int32(pokerface.SRFlags_SRUserReplyOnly) != 0 {
			continue
		}

		step++
		if step > stepback {
			return sraction
		}
	}

	return nil
}

func (lc *LoopContext) current() *pokerface.SRAction {
	return lc.fetchNonUserReplyOnly(0)
}

func (lc *LoopContext) prev() *pokerface.SRAction {
	return lc.fetchNonUserReplyOnly(1)
}

func (lc *LoopContext) prevprev() *pokerface.SRAction {
	return lc.fetchNonUserReplyOnly(2)
}

// unixTimeInMinutes 获取系统时间，并转换为分钟
func unixTimeInMinutes() uint32 {
	return uint32(time.Now().Unix() / 60)
}

// snapshootDealActions 在开始出牌前保存发牌信息
func (lc *LoopContext) snapshootDealActions() {
	room := lc.s.room
	msgRecorder := lc.recorder
	// 记录庄家和风花牌
	var bankerChairID32 = int32(room.bankerPlayer().chairID)
	msgRecorder.BankerChairID = &bankerChairID32
	var windFlowerID32 = int32(0)
	msgRecorder.WindFlowerID = &windFlowerID32
	var isHandOver = false
	msgRecorder.IsHandOver = &isHandOver
	msgRecorder.RoomConfigID = &room.configID
	var timeRecord = unixTimeInMinutes()
	msgRecorder.StartTime = &timeRecord
	var handNum32 = int32(room.handRoundStarted)
	msgRecorder.HandNum = &handNum32
	var isContinuoursBanker = room.bankerPlayer().gStatis.isContinuousBanker
	msgRecorder.IsContinuousBanker = &isContinuoursBanker
	msgRecorder.RoomNumber = &lc.s.room.roomNumber
	var roomType32 = int32(gconst.RoomType_DafengGZ)
	msgRecorder.RoomType = &roomType32

	extra := &pokerface.SRMsgHandRecorderExtra{}
	var markup32 = int32(room.markup)
	extra.Markup = &markup32
	extra.OwnerUserID = &room.ownerID
	msgRecorder.Extra = extra

	// 记录参与玩牌的玩家列表
	playerlist := make([]*pokerface.SRMsgPlayerInfo, len(room.players))
	for i, p := range room.players {
		sp := &pokerface.SRMsgPlayerInfo{}
		var chairID32 = int32(p.chairID)
		sp.ChairID = &chairID32
		var userID = p.userID()
		sp.UserID = &userID
		var nick = p.user.getInfo().nick
		sp.Nick = &nick

		var sex = p.user.getInfo().sex
		sp.Sex = &sex
		var headIconURL = p.user.getInfo().headIconURI
		sp.HeadIconURI = &headIconURL

		var avatarID = int32(p.user.getInfo().avatarID)
		sp.AvatarID = &avatarID

		playerlist[i] = sp
	}

	msgRecorder.Players = playerlist

	// 记录发牌数据
	deals := make([]*pokerface.SRDealDetail, len(room.players))
	for i, p := range room.players {
		deal := &pokerface.SRDealDetail{}
		var chairID32 = int32(p.chairID)
		deal.ChairID = &chairID32

		cards := p.cards
		deal.CardsHand = cards.hand2IDList()
		// deal.CardsFlower = cards.flower2IDList()

		deals[i] = deal
	}

	msgRecorder.Deals = deals
}

// addDrawAction 记录抽牌动作，配牌工具需要根据flags不同对抽牌动作做必要的选择
func (lc *LoopContext) addActionWithCards(who *PlayerHolder, action ActionType, msgCardHand *pokerface.MsgCardHand, qaIndex int, flags pokerface.SRFlags) {
	var msgSRAction = &pokerface.SRAction{}
	var action32 = int32(action)
	msgSRAction.Action = &action32
	var chairID32 = int32(who.chairID)
	msgSRAction.ChairID = &chairID32
	var qaIndex32 = int32(qaIndex)
	msgSRAction.QaIndex = &qaIndex32
	var flags32 = int32(flags)
	msgSRAction.Flags = &flags32

	if msgCardHand != nil {
		cards := make([]int32, 0, len(msgCardHand.Cards)+1)
		cards = append(cards, msgCardHand.GetCardHandType())
		cards = append(cards, msgCardHand.Cards...)
		msgSRAction.Cards = cards
	}

	lc.actionList.PushBack(msgSRAction)
}

// addHandWashout 流局
func (lc *LoopContext) finishHandWashout(handScore *pokerface.MsgHandScore) {
	lc.finishWinnerBorn(handScore)
}

// addWinnerBorn 胡牌
func (lc *LoopContext) finishWinnerBorn(handScore *pokerface.MsgHandScore) {
	msgRecorder := lc.recorder
	var isHandOver = true
	msgRecorder.IsHandOver = &isHandOver
	bytes, err := proto.Marshal(handScore)
	if err == nil {
		msgRecorder.HandScore = bytes
	}

	var timeRecord = unixTimeInMinutes()
	msgRecorder.EndTime = &timeRecord
	lc.actionList2Actions()

	if lc.s.room.isUlimitRound {
		return
	}

	lc.snapshootReplayRecordSummary(lc.s.room)
}

// actionList2Actions 转换到proto的action list
func (lc *LoopContext) actionList2Actions() {
	msgRecorder := lc.recorder
	actions := make([]*pokerface.SRAction, lc.actionList.Len())
	var i = 0
	for e := lc.actionList.Front(); e != nil; e = e.Next() {
		a := e.Value.(*pokerface.SRAction)
		actions[i] = a
		i++
	}

	msgRecorder.Actions = actions
}

// dump 打印
func (lc *LoopContext) dump() {

	for e := lc.actionList.Front(); e != nil; e = e.Next() {
		a := e.Value.(*pokerface.SRAction)
		dumpSRAction(a)
	}

	buf, err := proto.Marshal(lc.recorder)
	if err != nil {
		log.Panicln(err)
	}

	log.Printf("recorder size:%d\n", len(buf))
}

// dumpSRAction 打印action
func dumpSRAction(a *pokerface.SRAction) {
	log.Printf("chair:%d, a:%d, qi:%d,flag:%d\n", a.GetChairID(), a.GetAction(), a.GetQaIndex(), a.GetFlags())
}

// toByteArray 转换为byte数组
func (lc *LoopContext) toByteArray() []byte {
	lc.actionList2Actions()
	buf, err := proto.Marshal(lc.recorder)
	if err != nil {
		log.Panicln(err)
	}

	return buf
}

// dump2Redis 输出到redis，注意此时不能再引用room中的players
func (lc *LoopContext) dump2Redis(s *SPlaying) {
	if lc.s.room.isUlimitRound {
		return
	}

	if lc.recorder.Actions == nil {
		lc.actionList2Actions()
	}

	newUUID, err := uuid.NewV4()
	if err != nil {
		log.Panicln("dump2Redis failed, new uuid error:", err)
	}

	recordID := fmt.Sprintf("%s", newUUID)

	var shareAbleID = lc.randomRecordShareAbleID(recordID)
	log.Printf("dump2Redis, new shareID:%s for Record:%s\n", shareAbleID, recordID)

	lc.recorder.ShareAbleID = &shareAbleID

	buf, err := proto.Marshal(lc.recorder)
	if err != nil {
		log.Panicln(err)
	}

	// 获取redis链接，并退出函数时释放
	conn := pool.Get()
	defer conn.Close()

	_, err = conn.Do("HMSET", gconst.GameServerMJRecorderTablePrefix+recordID, "d", buf, "r", s.room.ID, "cid", s.room.configID)
	if err != nil {
		log.Println(err)
		return
	}

	lc.updateReplayRoom(conn, s.room.ID, recordID, shareAbleID)

	// 为每一个玩家添加回播房间记录
	for _, p := range lc.msgReplayRoom.Players {
		var replayRooms []string
		replayRoomStr, err := redis.String(conn.Do("HGET", gconst.LobbyPlayerTablePrefix+p.GetUserID(), "rr"))
		if err != nil {
			// 由于玩家没有玩过本游戏因此可能mjrc不存在，redis返回nil错误，因此是正常情况不需要输出日志
			//log.Println(err)
			replayRoomStr = ""
		} else {
			replayRooms = strings.Split(replayRoomStr, ",")
		}

		found := false
		if len(replayRooms) > 0 {
			for _, rr := range replayRooms {
				if rr == s.room.ID {
					found = true
					break
				}
			}

		}

		// 房间已经存在
		if found {
			continue
		}

		// 限制每个用户只能保存MJMaxReplayRoomNumber个最近记录
		// 如果裁剪了用户的最近记录，则需要检查记录是否已经无人引用，如果是则彻底删除记录
		if len(replayRooms) >= maxReplayRoomNumber {
			lc.unbindMJReplayRoomIfUseless(conn, replayRooms[0], p.GetUserID())
			replayRooms = replayRooms[1:]
		}

		replayRooms = append(replayRooms, s.room.ID)
		_, err = conn.Do("HSET", gconst.LobbyPlayerTablePrefix+p.GetUserID(), "rr", strArray2Comma(replayRooms))
		if err != nil {
			log.Println(err)
		}
	}

	// 如果房间是俱乐部房间，则需要把回播记录写到俱乐部的回播列表中
	// if lc.s.room.clubID != "" {
	// 	lc.appendCurrentReplayRoom2Club(conn)
	// }
}

func (lc *LoopContext) snapshootReplayRecordSummary(room *Room) {
	// debug.PrintStack()
	// 附加本手牌结果概要
	var replayRecordSummary = &pokerface.MsgReplayRecordSummary{}
	// replayRecordSummary.RecordUUID = &recordID
	// replayRecordSummary.ShareAbleID = &shareAbleID
	replayRecordSummary.StartTime = lc.recorder.StartTime
	var endTime32 = uint32(unixTimeInMinutes())
	replayRecordSummary.EndTime = &endTime32

	lastHand := room.scoreRecords[len(room.scoreRecords)-1]
	// 如果是不流局，则把所有玩家的得分概要保存
	if lastHand.GetEndType() != int32(HandOverType_enumHandOverType_None) {
		replayPlayerScores := make([]*pokerface.MsgReplayPlayerScoreSummary, len(lastHand.PlayerRecords))
		for i, rp := range lastHand.PlayerRecords {
			playerScore := &pokerface.MsgReplayPlayerScoreSummary{}
			playerScore.WinType = rp.WinType
			var chairID32 = int32(room.getPlayerByUserID(rp.GetUserID()).chairID)
			playerScore.ChairID = &chairID32
			playerScore.Score = rp.Score

			replayPlayerScores[i] = playerScore
		}
		replayRecordSummary.PlayerScores = replayPlayerScores
	}

	lc.replayRecordSummary = replayRecordSummary

	var msgReplayRoom = &pokerface.MsgReplayRoom{}
	var recordRoomType32 = int32(gconst.RoomType_DafengGZ)
	msgReplayRoom.RecordRoomType = &recordRoomType32
	msgReplayRoom.RoomNumber = &room.roomNumber
	var startTime = unixTimeInMinutes()
	msgReplayRoom.StartTime = &startTime
	msgReplayRoom.OwnerUserID = &room.ownerID

	// 玩家列表
	var replayPlayers = make([]*pokerface.MsgReplayPlayerInfo, len(room.players))
	for i, p := range room.players {
		rp := &pokerface.MsgReplayPlayerInfo{}
		var chairID32 = int32(p.chairID)
		rp.ChairID = &chairID32
		var userID = p.userID()
		rp.UserID = &userID
		var nick = p.user.getInfo().nick
		rp.Nick = &nick
		var sex = p.user.getInfo().sex
		rp.Sex = &sex
		var headIconURL = p.user.getInfo().headIconURI
		rp.HeadIconURI = &headIconURL

		var totalScore32 = int32(p.gStatis.roundScore)
		rp.TotalScore = &totalScore32

		var avatarID = int32(p.user.getInfo().avatarID)
		rp.AvatarID = &avatarID

		replayPlayers[i] = rp
	}

	log.Println("snapshootReplayRecordSummary, players count:", len(replayPlayers))
	msgReplayRoom.Players = replayPlayers

	lc.msgReplayRoom = msgReplayRoom
}

// updateReplayRoom 更新回播房间记录
func (lc *LoopContext) updateReplayRoom(conn redis.Conn, roomID string, recordID string, shareAbleID string) {

	// 加载房间回播记录列表
	recordsStr := ""
	var bb []byte
	var records []string
	values, err := redis.Values(conn.Do("HMGET", gconst.GameServerMJReplayRoomTablePrefix+roomID, "hr", "d"))
	if err == nil {
		recordsStr, err = redis.String(values[0], err)
		if err == nil {
			records = strings.Split(recordsStr, ",")
		}
		bb, _ = redis.Bytes(values[1], err)
	}

	records = append(records, recordID)

	msgReplayRoom := &pokerface.MsgReplayRoom{}
	if bb != nil && len(bb) > 0 {
		err = proto.Unmarshal(bb, msgReplayRoom)
		if err != nil {
			log.Println("updateReplayRoom proto err:", err)
		}

		// 保存每一个人的总得分
		msgReplayRoom.Players = lc.msgReplayRoom.Players
		log.Println("updateReplayRoom, use old, players count:", len(msgReplayRoom.Players))

	} else {
		msgReplayRoom = lc.msgReplayRoom
		log.Println("updateReplayRoom, use new, players count:", len(msgReplayRoom.Players))
	}

	userIDs := make([]string, len(msgReplayRoom.Players))
	for i, p := range msgReplayRoom.Players {
		userIDs[i] = p.GetUserID()
	}

	var replayRecordSummary = lc.replayRecordSummary
	replayRecordSummary.RecordUUID = &recordID
	replayRecordSummary.ShareAbleID = &shareAbleID

	msgReplayRoom.Records = append(msgReplayRoom.Records, replayRecordSummary)

	var endTime = unixTimeInMinutes()
	msgReplayRoom.EndTime = &endTime

	buf, err := proto.Marshal(msgReplayRoom)
	if err != nil {
		log.Println("updateReplayRoom, marshal error:", err)
		return
	}

	unixTimeInMinutes32 := unixTimeInMinutes()
	// 记录加到房间的回放列表
	conn.Send("MULTI")
	conn.Send("HMSET", gconst.GameServerMJReplayRoomTablePrefix+roomID, "hr", strArray2Comma(records),
		"d", buf, "rrt", int(gconst.RoomType_DafengGZ), "date", unixTimeInMinutes32)

	for _, u := range userIDs {
		conn.Send("SADD", gconst.GameServerReplayRoomsReferenceSetPrefix+roomID, u)
	}
	conn.Do("EXEC")
}

/*
func (lc *LoopContext) appendCurrentReplayRoom2Club(conn redis.Conn) {
	clubID := lc.s.room.clubID
	roomID := lc.s.room.ID

	conn.Send("MULTI")
	conn.Send("SISMEMBER", gconst.ClubReplayRoomsSetPrefix+clubID, roomID)
	conn.Send("LLEN", gconst.ClubReplayRoomsListPrefix+clubID)
	valus, err := redis.Values(conn.Do("EXEC"))

	if err != nil {
		log.Println("appendCurrentReplayRoom2Club, redis err:", err)
		return
	}

	isMember, err := redis.Int(valus[0], nil)
	if err != nil && err != redis.ErrNil {
		log.Println("appendCurrentReplayRoom2Club, isMember redis err:", err)
		return
	}

	if isMember == 0 {
		// 使用LUA脚本来修改俱乐部回播记录，以防止俱乐部被删除后，还会添加记录
		luaScriptClubReplayRoom.Do(conn, clubID, roomID)
	}

	llen, err := redis.Int(valus[1], nil)
	if err != nil && err != redis.ErrNil {
		log.Println("appendCurrentReplayRoom2Club, isMember redis err:", err)
		return
	}

	// 尽量保持俱乐部的回播房间列表长度不要过长
	if llen > gconst.MaxClubReplayRoomsNum {
		// 过长则裁剪列表和解除引用
		removedRoomID, err := redis.String(conn.Do("RPOP", gconst.ClubReplayRoomsListPrefix+clubID))
		if err == nil {
			conn.Send("MULTI")
			conn.Send("SREM", gconst.ClubReplayRoomsSetPrefix+clubID, removedRoomID)
			conn.Send("SREM", gconst.GameServerReplayRoomsReferenceSetPrefix+removedRoomID, clubID)
			conn.Do("EXEC")

			lc.unbindMJReplayRoomIfUseless(conn, removedRoomID, clubID)
		} else {
			log.Println("appendCurrentReplayRoom2Club, RPOP redis err:", err)
		}
	}
}
*/

// strArray2Comma 字符串数据转为逗号分隔字符串
func strArray2Comma(ss []string) string {
	result := ""
	if len(ss) < 1 {
		return result
	}

	for i := 0; i < len(ss)-1; i++ {
		result = result + ss[i] + ","
	}

	result = result + ss[len(ss)-1]

	return result
}

// unbindMJReplayRoomIfUseless 解除回播房间的引用
func (lc *LoopContext) unbindMJReplayRoomIfUseless(conn redis.Conn, roomID string, referenceBy string) {
	conn.Send("MULTI")
	conn.Send("SREM", gconst.GameServerReplayRoomsReferenceSetPrefix+roomID, referenceBy)
	conn.Send("SCARD", gconst.GameServerReplayRoomsReferenceSetPrefix+roomID) // 新的回播房间引用关系是保存于这个set中
	conn.Send("HGET", gconst.GameServerMJReplayRoomTablePrefix+roomID, "u")   // 旧的回播房间引用则是保存在这个field中，以逗号分隔

	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil && err != redis.ErrNil {
		log.Println("unbindMJReplayRoomIfUseless, redis err:", err)
		return
	}

	refSetSize, _ := redis.Int(values[1], nil)
	refUserStr, _ := redis.String(values[2], nil) // 这部分相关代码，是由于需要迁移老的回播房间引用存储方式而保留

	if refSetSize > 0 {
		// 还有其他引用
		return
	}

	needUnBind := true
	// 这部分相关代码，是由于需要迁移老的回播房间引用存储方式而保留
	if refUserStr != "" {

		found := false
		users := strings.Split(refUserStr, ",")
		for i, u := range users {
			if u == referenceBy {
				found = true
				users = append(users[:i], users[i+1:]...)
				break
			}
		}

		if found {
			// 更新引用的用户列表
			conn.Send("MULTI")
			for _, u := range users {
				conn.Send("SADD", gconst.GameServerReplayRoomsReferenceSetPrefix+roomID, u)
			}
			conn.Send("HDEL", gconst.GameServerMJReplayRoomTablePrefix+roomID, "u") // 把旧的清理掉
			conn.Do("EXEC")
		}

		if len(users) > 0 {
			// 还有老数据引用着这个回播房间，暂时不能删除
			needUnBind = false
		}
	}

	if needUnBind {
		// 没有人引用这个记录了，可以删除，删除房间每一局打牌记录
		lc.unbindMJReplayRecords(conn, roomID)

		// 清理该房间的回播记录
		conn.Send("MULTI")
		conn.Send("DEL", gconst.GameServerReplayRoomsReferenceSetPrefix+roomID)
		conn.Send("DEL", gconst.GameServerMJReplayRoomTablePrefix+roomID)
		_, err := conn.Do("EXEC")
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("del replay room:", roomID)
	}
}

// unbindMJReplayRecords 解除回放记录引用
func (lc *LoopContext) unbindMJReplayRecords(conn redis.Conn, roomID string) {
	var records []string
	recordsStr, err := redis.String(conn.Do("HGET", gconst.GameServerMJReplayRoomTablePrefix+roomID, "hr"))
	if err != nil {
	} else {
		records = strings.Split(recordsStr, ",")
	}

	if len(records) < 1 {
		return
	}

	// 先登记到已删除set中
	conn.Send("MULTI")
	for _, r := range records {
		// 添加到已经被删除set中，以便定时清理任务可以通知sqlserver等持久化数据库做清理
		conn.Send("SADD", gconst.GameServerMJRecorderDeletedSet, r)
	}
	conn.Do("EXEC")

	conn.Send("MULTI")
	for _, r := range records {
		conn.Send("HGET", gconst.GameServerMJRecorderTablePrefix+r, "sid")
	}
	sids, _ := redis.Strings(conn.Do("EXEC"))

	// 删除shared id，但是如果record不存在于redis中，则无法获得shared id
	conn.Send("MULTI")
	for _, sid := range sids {
		if sid != "" {
			conn.Send("DEL", gconst.GameServerMJRecorderTablePrefix+sid)
			// 新的sid存在于MJRecorderShareIDTable哈希表中
			conn.Send("HDEL", gconst.GameServerMJRecorderShareIDTable, sid)
			lc.cl.Println("delete hand record sid:", sid)
		}
	}
	conn.Do("EXEC")

	// 先从redis中删除，回播记录可能已经不存在于redis，已经被腾挪到持久化数据库
	conn.Send("MULTI")
	for _, r := range records {
		conn.Send("DEL", gconst.GameServerMJRecorderTablePrefix+r)
		lc.cl.Println("delete hand record:", r)
	}
	conn.Do("EXEC")
}

func loadMJLastRecordForRoom(conn redis.Conn, roomID string) []byte {
	// 加载回播房间记录
	recordsStr, err := redis.String(conn.Do("HGET", gconst.GameServerMJReplayRoomTablePrefix+roomID, "hr"))
	if err != nil {
		recordsStr = ""
	}

	log.Println("recordsStr:", recordsStr)
	records := strings.Split(recordsStr, ",")
	if len(records) < 1 {
		log.Println("user has replay room, but no hand record")
		return nil
	}

	recordID := records[len(records)-1]
	return loadMJRecord(conn, recordID)
}

func loadMJRecord(conn redis.Conn, recordID string) []byte {
	log.Printf("load %s from %s\n", recordID, gconst.GameServerMJRecorderTablePrefix+recordID)
	buf, err := redis.Bytes(conn.Do("HGET", gconst.GameServerMJRecorderTablePrefix+recordID, "d"))
	if err != nil && err != redis.ErrNil {
		log.Println(err)
		return nil
	}

	if err == redis.ErrNil {
		return loadMJRecordFromSQLServer(recordID)
	}

	return buf
}

func loadMJRecordFromSQLServer(recordID string) []byte {
	// conn, err := mssql.StartMssql(gscfg.DbIP, gscfg.DbPort, gscfg.DbUser, gscfg.DbPassword, gscfg.DbName)
	// if err != nil {
	// 	log.Println("handleLoadReplayRecord, StartMssql err:", err)
	// 	return nil
	// }

	// defer conn.Close()

	// var grcRecord = mssql.LoadGRCRcordFromSQLServer(recordID, conn)
	// if grcRecord != nil {
	// 	return grcRecord.RecordData
	// }

	return nil
}

// loadMJLastRecordForUser 从redis加载最后一手牌记录
func loadMJLastRecordForUser(userID string) []byte {
	// 获取redis链接，并退出函数时释放
	conn := pool.Get()
	defer conn.Close()

	replayRoomsStr, err := redis.String(conn.Do("HGET", gconst.LobbyPlayerTablePrefix+userID, "rr"))
	if err != nil {
		log.Println(err)
		return nil
	}

	replayRooms := strings.Split(replayRoomsStr, ",")
	if len(replayRooms) < 1 {
		// 没有数据
		log.Println("user has no mj replay room, userID:", userID)
		return nil
	}

	replayRoom := replayRooms[len(replayRooms)-1]

	return loadMJLastRecordForRoom(conn, replayRoom)
}

// loadMJRoomRecardShareIDs 加载房间所有的分享码ID,MJTestTool会使用一个shareID
// 来加载其所在的房间的所有shareID，然后再逐个回播记录加载，并保存到文件夹中
func loadMJRoomRecardShareIDs(conn redis.Conn, recordID string) ([]byte, error) {

	roomID, err := redis.String(conn.Do("HGET", gconst.GameServerMJRecorderTablePrefix+recordID, "r"))
	if err != nil {
		log.Println("loadMJRoomRecardShareIDs, roomID failed:", err)
		return nil, err
	}

	ridsWithComma, err := redis.String(conn.Do("HGET", gconst.GameServerMJReplayRoomTablePrefix+roomID, "hr"))
	if err != nil {
		log.Println("loadMJRoomRecardShareIDs, hr failed:", err)
		return nil, err
	}

	records := strings.Split(ridsWithComma, ",")
	conn.Send("MULTI")
	for _, r := range records {
		conn.Send("HGET", gconst.GameServerMJRecorderTablePrefix+r, "sid")
	}

	sids, err := redis.Strings(conn.Do("EXEC"))
	if err != nil {
		log.Println("loadMJRoomRecardShareIDs, sids failed:", err)
		return nil, err
	}

	strBytes := bytes.NewBufferString("")
	for _, sid := range sids {
		strBytes.WriteString(sid)
		strBytes.WriteString("\n")
	}

	return strBytes.Bytes(), nil
}

func (lc *LoopContext) randomRecordShareAbleID(recordID string) string {
	maxRetry := 3
	for i := 0; i < maxRetry; i++ {
		shareAbleID := lc.randomRecordShareAbleIDImpl(recordID)
		if shareAbleID != "" {
			return shareAbleID
		}
	}

	log.Println("ERROR, randomRecordShareAbleID failed to alloc shareAble ID")

	return "0000"
}

func (lc *LoopContext) randomRecordShareAbleIDImpl(recordID string) string {
	const maxTry = 20
	rand := lc.s.room.rand
	shareAbleIDArray := make([]string, maxTry)
	for i := 0; i < maxTry; i++ {
		shareAbleID := rand.Intn(MaxShareAbleID-MinShareAbleID) + MinShareAbleID
		shareAbleIDStr := fmt.Sprintf("%d", shareAbleID)
		shareAbleIDArray[i] = shareAbleIDStr
	}

	shareAbleIDStrs := strArray2Comma(shareAbleIDArray)
	return validRedisRandNumber(shareAbleIDStrs, recordID)
}

// 1.检查数据库是否已经存在随机数
// 2.若不存在，则保存到数据库，然后返回这个随机数
func validRedisRandNumber(shareAbleIDStrs string, recordID string) string {
	conn := pool.Get()
	defer conn.Close()

	// luaScript 在startRedis中创建
	randNumber, err := redis.String(luaScript.Do(conn, gconst.GameServerMJRecorderTablePrefix, recordID,
		shareAbleIDStrs, gconst.GameServerMJRecorderShareIDTable))
	if err != nil {
		logrus.Printf("randromNumber error, roomNumbers %s, roomID %s, error:%v \n", shareAbleIDStrs, recordID, err)
	}

	// TODO: 由于是随机测试可用分享号码，因此如果所有测试失败，则会返回空的分享号码，后面可以通过加大测试样本数量，或者
	// 重放尝试N次，或者索性用逐次递增检查的方式来获得可用的分享号码
	return randNumber
}
