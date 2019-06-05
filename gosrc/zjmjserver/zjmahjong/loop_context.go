package zjmahjong

import (
	"container/list"
	"gconst"
	"mahjong"

	log "github.com/sirupsen/logrus"

	"time"

	"github.com/sirupsen/logrus"

	"github.com/garyburd/redigo/redis"

	"github.com/golang/protobuf/proto"
)

// LoopContext 打牌循环上下文
// 主要保存最近动作的玩家
// 录像会保存到redis，一是供玩家回放，二是出问题（例如玩家指出数值计算不正确等）我们可以回放查找bug
// 极端情况下，recorder marshal后的大小接近3K
// 因此保存redis的时候需要考虑保存速度，以及对redis内存的压力
type LoopContext struct {
	replayRecordSummary *mahjong.MsgReplayRecordSummary

	recorder *mahjong.SRMsgHandRecorder

	actionList *list.List

	s *SPlaying

	cl *logrus.Entry
}

func newLoopContext(s *SPlaying) *LoopContext {
	ctx := &LoopContext{}
	ctx.s = s
	ctx.recorder = &mahjong.SRMsgHandRecorder{}
	ctx.actionList = list.New()
	ctx.cl = s.cl

	return ctx
}

func (lc *LoopContext) isSelfDraw(player *PlayerHolder) bool {
	cur := lc.current()
	if cur == nil {
		return false
	}

	chairID := int(cur.GetChairID())
	action := int(cur.GetAction())

	return chairID == player.chairID && action == int(mahjong.ActionType_enumActionType_DRAW)
}

func (lc *LoopContext) isRobKong() bool {
	// 首先从动作队列尾部开始，往回跳过胡牌的动作
	var a *mahjong.SRAction
	winChuckFound := 0
	for e := lc.actionList.Back(); e != nil; e = e.Prev() {
		a = e.Value.(*mahjong.SRAction)

		if a.GetFlags()&int32(mahjong.SRFlags_SRUserReplyOnly) != 0 {
			continue
		}

		if int(a.GetAction()) != int(mahjong.ActionType_enumActionType_WIN_Chuck) {
			break
		} else {
			winChuckFound++
		}
	}

	if a != nil && a.GetAction() == int32(mahjong.ActionType_enumActionType_KONG_Triplet2) && winChuckFound > 0 {
		return true
	}

	return false
}

func (lc *LoopContext) kongerOf(me *PlayerHolder, room *Room) *PlayerHolder {
	// 操作系列： 打牌---杠----摸牌----自摸胡牌
	var cur = lc.current()
	pre := lc.prev()
	prepre := lc.prevprev()
	preprepre := lc.fetchNonUserReplyOnly(3)

	if cur == nil || pre == nil || prepre == nil || preprepre == nil {
		return nil
	}

	if int(cur.GetChairID()) != me.chairID {
		// 自己自摸
		return nil
	}

	if me.chairID != int(pre.GetChairID()) {
		// 上一个，上上一个操作都应该是自己
		return nil
	}

	if me.chairID != int(prepre.GetChairID()) {
		// 上一个，上上一个操作都应该是自己
		return nil
	}

	if me.chairID == int(preprepre.GetChairID()) {
		// 打牌者不能是自己
		return nil
	}

	if int(pre.GetAction()) != int(mahjong.ActionType_enumActionType_DRAW) {
		// 摸牌
		return nil
	}

	if int(prepre.GetAction()) != int(mahjong.ActionType_enumActionType_KONG_Exposed) {
		// 明杠
		return nil
	}

	if int(preprepre.GetAction()) != int(mahjong.ActionType_enumActionType_DISCARD) {
		// 打牌
		return nil
	}

	chairID := int(preprepre.GetChairID())
	return room.getPlayerByChairID(chairID)
}

func (lc *LoopContext) isSelfKong(me *PlayerHolder) bool {
	// 操作系列： 杠（续杠，暗杠）----摸牌----自摸胡牌
	var cur = lc.current()
	pre := lc.prev()
	prepre := lc.prevprev()

	if cur == nil || pre == nil || prepre == nil {
		return false
	}

	if int(cur.GetChairID()) != me.chairID {
		// 自己自摸
		return false
	}

	if me.chairID != int(pre.GetChairID()) {
		// 上一个，上上一个操作都应该是自己
		return false
	}

	if me.chairID != int(prepre.GetChairID()) {
		// 上一个，上上一个操作都应该是自己
		return false
	}

	if int(pre.GetAction()) != int(mahjong.ActionType_enumActionType_DRAW) {
		// 摸牌
		return false
	}

	action := prepre.GetAction()
	if int(action) != int(mahjong.ActionType_enumActionType_KONG_Concealed) &&
		int(action) != int(mahjong.ActionType_enumActionType_KONG_Triplet2) {
		// 续杠或者暗杠
		return false
	}

	return true
}

func (lc *LoopContext) fetchNonUserReplyOnly(stepback int) *mahjong.SRAction {
	var step = 0
	for e := lc.actionList.Back(); e != nil; e = e.Prev() {
		sraction := e.Value.(*mahjong.SRAction)
		if sraction.GetFlags()&int32(mahjong.SRFlags_SRUserReplyOnly) != 0 {
			continue
		}

		step++
		if step > stepback {
			return sraction
		}
	}

	return nil
}

func (lc *LoopContext) getLastNonDrawAction() *mahjong.SRAction {
	var srAction = lc.fetchNonUserReplyOnly(0)
	if srAction.GetAction() != int32(mahjong.ActionType_enumActionType_DRAW) {
		return srAction
	}

	srAction = lc.fetchNonUserReplyOnly(1)
	return srAction
}

func (lc *LoopContext) current() *mahjong.SRAction {
	return lc.fetchNonUserReplyOnly(0)
}

func (lc *LoopContext) prev() *mahjong.SRAction {
	return lc.fetchNonUserReplyOnly(1)
}

func (lc *LoopContext) prevprev() *mahjong.SRAction {
	return lc.fetchNonUserReplyOnly(2)
}

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
	var roomType32 = int32(myRoomType)
	msgRecorder.RoomType = &roomType32

	extra := &mahjong.SRMsgHandRecorderExtra{}
	var markup32 = int32(0)
	extra.Markup = &markup32
	extra.OwnerUserID = &room.ownerID
	msgRecorder.Extra = extra

	// 记录参与玩牌的玩家列表
	playerlist := make([]*mahjong.SRMsgPlayerInfo, len(room.players))
	for i, p := range room.players {
		sp := &mahjong.SRMsgPlayerInfo{}
		var chairID32 = int32(p.chairID)
		sp.ChairID = &chairID32
		var userID = p.userID()
		sp.UserID = &userID
		var nick = p.user.getInfo().nick
		sp.Nick = &nick

		var gender = p.user.getInfo().gender
		sp.Gender = &gender
		var headIconURL = p.user.getInfo().headIconURI
		sp.HeadIconURI = &headIconURL

		var avatarID = int32(p.user.getInfo().avatarID)
		sp.AvatarID = &avatarID

		playerlist[i] = sp
	}

	msgRecorder.Players = playerlist

	// 记录发牌数据
	deals := make([]*mahjong.SRDealDetail, len(room.players))
	for i, p := range room.players {
		deal := &mahjong.SRDealDetail{}
		var chairID32 = int32(p.chairID)
		deal.ChairID = &chairID32

		tiles := p.tiles
		deal.TilesHand = tiles.hand2IDList()
		deal.TilesFlower = tiles.flower2IDList()

		deals[i] = deal
	}

	msgRecorder.Deals = deals
}

// addDrawAction 记录抽牌动作
func (lc *LoopContext) addDrawAction(who *PlayerHolder, tileIDHand int, tileIDsFlower []*Tile, qaIndex int) {
	var msgSRAction = &mahjong.SRAction{}
	var action32 = int32(mahjong.ActionType_enumActionType_DRAW)
	msgSRAction.Action = &action32
	var chairID32 = int32(who.chairID)
	msgSRAction.ChairID = &chairID32
	var qaIndex32 = int32(qaIndex)
	msgSRAction.QaIndex = &qaIndex32
	var flags32 = int32(mahjong.SRFlags_SRNone)
	msgSRAction.Flags = &flags32

	tiles := make([]int32, 1+len(tileIDsFlower))
	i := 0

	for _, t := range tileIDsFlower {
		tiles[i] = int32(t.tileID)
		i++
	}

	tiles[i] = int32(tileIDHand)

	msgSRAction.Tiles = tiles

	lc.actionList.PushBack(msgSRAction)

	//lc.drawCount++
	//lc.actionCount++
}

// addActionWithTile 记录关于牌的动作，例如吃椪杠，注意虽然操作的是面子牌组，但只需要保存牌组第一张牌即可
func (lc *LoopContext) addActionWithTile(who *PlayerHolder, tileID int,
	chowTile int, action mahjong.ActionType, qaIndex int, flags mahjong.SRFlags,
	allowActions int) {
	var msgSRAction = &mahjong.SRAction{}
	var action32 = int32(action)
	msgSRAction.Action = &action32
	var chairID32 = int32(who.chairID)
	msgSRAction.ChairID = &chairID32
	var qaIndex32 = int32(qaIndex)
	msgSRAction.QaIndex = &qaIndex32
	var flags32 = int32(flags)
	msgSRAction.Flags = &flags32

	allowActions32 := int32(allowActions)
	msgSRAction.AllowActions = &allowActions32

	if tileID != InvalidTile.tileID {
		tiles := []int32{int32(tileID)}
		if action == mahjong.ActionType_enumActionType_CHOW {
			tiles = []int32{int32(tileID), int32(chowTile)}
		}

		msgSRAction.Tiles = tiles
	}

	lc.actionList.PushBack(msgSRAction)

	//if flags == 0 {
	//	lc.actionCount++
	//}
}

func (lc *LoopContext) addActionWithTiles(who *PlayerHolder, tileIDs []int,
	action mahjong.ActionType, qaIndex int, flags mahjong.SRFlags,
	allowActions int) {

	var msgSRAction = &mahjong.SRAction{}
	var action32 = int32(action)
	msgSRAction.Action = &action32
	var chairID32 int32
	if who != nil {
		chairID32 = int32(who.chairID)
	}
	msgSRAction.ChairID = &chairID32
	var qaIndex32 = int32(qaIndex)
	msgSRAction.QaIndex = &qaIndex32
	var flags32 = int32(flags)
	msgSRAction.Flags = &flags32

	allowActions32 := int32(allowActions)
	msgSRAction.AllowActions = &allowActions32

	tilesID32 := make([]int32, len(tileIDs))
	for i, tid := range tileIDs {
		tilesID32[i] = int32(tid)
	}

	msgSRAction.Tiles = tilesID32

	lc.actionList.PushBack(msgSRAction)

	//if flags == 0 {
	//	lc.actionCount++
	//}
}

// addHandWashout 流局
func (lc *LoopContext) finishHandWashout() {
	msgRecorder := lc.recorder
	var isHandOver = true
	msgRecorder.IsHandOver = &isHandOver
	var timeRecord = unixTimeInMinutes()
	msgRecorder.EndTime = &timeRecord
	lc.actionList2Actions()

	if lc.s.room.isUlimitRound {
		return
	}

	lc.snapshootReplayRecordSummary(lc.s.room)
}

// addWinnerBorn 胡牌
func (lc *LoopContext) finishWinnerBorn(handScore *mahjong.MsgHandScore) {
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
	actions := make([]*mahjong.SRAction, lc.actionList.Len())
	var i = 0
	for e := lc.actionList.Front(); e != nil; e = e.Next() {
		a := e.Value.(*mahjong.SRAction)
		actions[i] = a
		i++
	}

	msgRecorder.Actions = actions
}

// dump 打印
func (lc *LoopContext) dump() {

	for e := lc.actionList.Front(); e != nil; e = e.Next() {
		a := e.Value.(*mahjong.SRAction)
		lc.dumpSRAction(a)
	}

	buf, err := proto.Marshal(lc.recorder)
	if err != nil {
		lc.cl.Panicln(err)
	}

	lc.cl.Printf("recorder size:%d\n", len(buf))
}

// dumpSRAction 打印action
func (lc *LoopContext) dumpSRAction(a *mahjong.SRAction) {
	lc.cl.Printf("chair:%d, a:%d, qi:%d,flag:%d\n", a.GetChairID(), a.GetAction(), a.GetQaIndex(), a.GetFlags())
}

// toByteArray 转换为byte数组
func (lc *LoopContext) toByteArray() []byte {
	lc.actionList2Actions()
	buf, err := proto.Marshal(lc.recorder)
	if err != nil {
		lc.cl.Panicln(err)
	}

	return buf
}

func (lc *LoopContext) snapshootReplayRecordSummary(room *Room) {
	// debug.PrintStack()
	// 附加本手牌结果概要
	var replayRecordSummary = &mahjong.MsgReplayRecordSummary{}
	replayRecordSummary.StartTime = lc.recorder.StartTime
	var endTime32 = uint32(unixTimeInMinutes())
	replayRecordSummary.EndTime = &endTime32

	lastHand := room.scoreRecords[len(room.scoreRecords)-1]
	// 如果是不流局，则把所有玩家的得分概要保存
	if lastHand.GetEndType() != int32(mahjong.HandOverType_enumHandOverType_None) {
		replayPlayerScores := make([]*mahjong.MsgReplayPlayerScoreSummary, len(lastHand.PlayerRecords))
		for i, rp := range lastHand.PlayerRecords {
			playerScore := &mahjong.MsgReplayPlayerScoreSummary{}
			playerScore.WinType = rp.WinType
			var chairID32 = int32(room.getPlayerByUserID(rp.GetUserID()).chairID)
			playerScore.ChairID = &chairID32
			playerScore.Score = rp.Score

			replayPlayerScores[i] = playerScore
		}
		replayRecordSummary.PlayerScores = replayPlayerScores
	}

	lc.replayRecordSummary = replayRecordSummary
}

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

func loadMJLastRecordForRoom(conn redis.Conn, roomID string) []byte {
	// 加载回播房间记录
	bytes, err := redis.Bytes(conn.Do("HGET", gconst.GameServerMJReplayRoomTablePrefix+roomID, "d"))
	if err != nil {
		log.Println("loadMJLastRecordForRoom, err:", err)
		return nil
	}

	if bytes == nil {
		log.Println("loadMJLastRecordForRoom, bytes is nil, roomID:", roomID)
		return nil
	}

	var replayRoom = mahjong.MsgReplayRoom{}
	err = proto.Unmarshal(bytes, &replayRoom)
	if err != nil {
		log.Println("loadMJLastRecordForRoom, err:", err)
		return nil
	}

	records := replayRoom.GetRecords()
	if len(records) > 0 {
		recordID := records[len(records)-1].GetRecordUUID()

		return loadMJRecord(conn, recordID)
	}

	log.Println("loadMJLastRecordForRoom, no record found for room number:", replayRoom.RoomNumber)
	return nil
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

	replayRoomID, err := redis.String(conn.Do("LINDEX", gconst.GameServerMJReplayRoomListPrefix+userID, -1))

	if err != nil {
		log.Println("loadMJLastRecordForUser failed:", err)
		return nil
	}

	return loadMJLastRecordForRoom(conn, replayRoomID)
}
