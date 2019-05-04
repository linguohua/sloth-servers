package prunfast

import (
	"container/list"
	"gconst"
	"pokerface"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"

	"github.com/golang/protobuf/proto"
)

const (
	maxReplayRoomNumber = 50
)

// LoopContext 打牌循环上下文
// 主要保存最近动作的玩家
// 录像会保存到redis，一是供玩家回放，二是出问题（例如玩家指出数值计算不正确等）我们可以回放查找bug
// 极端情况下，recorder marshal后的大小接近3K
// 因此保存redis的时候需要考虑保存速度，以及对redis内存的压力
type LoopContext struct {
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
	var roomType32 = int32(myRoomType)
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
func (lc *LoopContext) addActionWithCards(who *PlayerHolder, action ActionType,
	msgCardHand *pokerface.MsgCardHand, qaIndex int, flags pokerface.SRFlags) {
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
	log.Printf("chair:%d, a:%d, qi:%d,flag:%d\n", a.GetChairID(),
		a.GetAction(), a.GetQaIndex(), a.GetFlags())
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

func (lc *LoopContext) snapshootReplayRecordSummary(room *Room) {
	// 附加本手牌结果概要
	var replayRecordSummary = &pokerface.MsgReplayRecordSummary{}
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
