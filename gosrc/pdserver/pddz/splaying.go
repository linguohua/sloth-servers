package pddz

import (
	"pokerface"
	"runtime/debug"

	"github.com/sirupsen/logrus"

	"github.com/golang/protobuf/proto"
)

// SPlaying 正在游戏状态
type SPlaying struct {
	taskDiscardReAction *TaskPlayerReAction
	taskPlayerAction    *TaskPlayerAction
	taskCallDoulbe      *TaskCallDouble

	//taskFirstReadyHand  *TaskFirstReadyHand
	cardMgr *CardMgr

	players []*PlayerHolder
	room    *Room
	lctx    *LoopContext

	cl *logrus.Entry

	lastAwardCards []int32
}

// newSPlaying 新建playing 状态机
func newSPlaying(room *Room) *SPlaying {
	s := &SPlaying{}
	s.room = room
	s.players = room.players
	s.cardMgr = newCardMgr(room, room.players)
	s.cl = room.cl

	return s
}

func (s *SPlaying) getStateName() string {
	return "SPlaying"
}

// onPlayerEnter 处理玩家进入消息，playing状态下不应该允许玩家进入
func (s *SPlaying) onPlayerEnter(player *PlayerHolder) {
	s.cl.Panic("Playing state no allow player enter")
}

// onPlayerLeave 处理玩家离开事件，playing状态下，并不从
// room中删除player对象
func (s *SPlaying) onPlayerLeave(player *PlayerHolder) {
	// 通知其他玩家，player离线，但是不从room删除player
	// 一直等待player上线，或者，等待其他玩家解散牌局
	player.state = pokerface.PlayerState_PSOffline
	s.room.updateRoomInfo2All()
	//s.room.writePlayerLeave2Redis(player, false)
}

// onPlayerReEnter 处理玩家重入事件，主要是掉线恢复
func (s *SPlaying) onPlayerReEnter(player *PlayerHolder) {
	player.state = pokerface.PlayerState_PSPlaying
	s.room.updateRoomInfo2All()

	// 掉线恢复
	// 先发送牌数据
	msgRestore := serializeMsgRestore(s, player)
	player.sendMsg(msgRestore, int32(pokerface.MessageCode_OPRestore))

	// 根据当前的gameLoop的等待状态，给玩家重新发送最近一个消息
	if s.taskPlayerAction != nil {
		s.taskPlayerAction.onPlayerRestore(player)
	} else if s.taskDiscardReAction != nil {
		s.taskDiscardReAction.onPlayerRestore(player)
	} else if s.taskCallDoulbe != nil {
		s.taskCallDoulbe.onPlayerRestore(player)
	}
}

// getStateConst state const
func (s *SPlaying) getStateConst() pokerface.RoomState {
	return pokerface.RoomState_SRoomPlaying
}

func (s *SPlaying) onStateEnter() {
	s.cl.Println("SPlaying enter")

	for _, p := range s.players {
		p.state = pokerface.PlayerState_PSPlaying
		p.resetForNextHand()
	}

	// 房间做一些准备
	s.room.handBegin()

	// 用户状态已经改变，发送更新到所有客户端
	s.room.updateRoomInfo2All()

	// 进入游戏循环
	go s.gameLoop()
}

func (s *SPlaying) onStateLeave() {
	// 需要把所有正在等待的Task cancel掉
	// 否则gameLoop的go routine不能退出
	s.cl.Println("SPlaying leave")

	if s.taskPlayerAction != nil {
		s.taskPlayerAction.cancel()
	}

	if s.taskDiscardReAction != nil {
		s.taskDiscardReAction.cancel()
	}

	if s.taskCallDoulbe != nil {
		s.taskCallDoulbe.cancel()
	}
}

func (s *SPlaying) onMessage(iu IUser, gmsg *pokerface.GameMessage) {
	var msgCode = pokerface.MessageCode(gmsg.GetOps())

	switch msgCode {
	case pokerface.MessageCode_OPAction:
		actionMsg := &pokerface.MsgPlayerAction{}
		err := proto.Unmarshal(gmsg.GetData(), actionMsg)
		if err == nil {
			s.onActionMessage(iu, actionMsg)
		} else {
			s.cl.Println("onMessage unmarshal error:", err)
		}
		break
	default:
		s.cl.Println("onMessage unsupported msgCode:", msgCode)
		break
	}
}

func (s *SPlaying) onActionMessage(iu IUser, msg *pokerface.MsgPlayerAction) {
	var player = s.getPlayerByID(iu.userID())

	if player == nil {
		s.cl.Panic("onActionMessage error, player is nil")
		return
	}

	if !s.room.qaIndexExpected(int(msg.GetQaIndex())) {
		s.cl.Printf("OnMessageAction error, qaIndex %d not expected, userId:%s\n", msg.GetQaIndex(), iu.userID())
		return
	}

	if player.expectedAction&int(msg.GetAction()) == 0 {
		s.cl.Printf("OnMessageAction allow actions %d not match %d, userId:%s\n", player.expectedAction, msg.GetAction(), iu.userID())
		s.cl.Panic("action not expected")
		return
	}

	// reset player expected actions
	player.expectedAction = 0

	var action = ActionType(msg.GetAction())
	switch action {
	case ActionType_enumActionType_SKIP:
		onMessageSkip(s, player, msg)
		break

	case ActionType_enumActionType_DISCARD:
		onMessageDiscard(s, player, msg)
		break

	case ActionType_enumActionType_Call:
		onMessageCall(s, player, msg)
		break

	case ActionType_enumActionType_Rob:
		onMessageRob(s, player, msg)
		break

	case ActionType_enumActionType_CallDouble:
		onMessageCallDouble(s, player, msg)
		break

	case ActionType_enumActionType_CallWithScore:
		onMessageCallWithScore(s, player, msg)
		break

	default:
		s.cl.Panic("OnMessageAction unsupported action code")
		break
	}
}

func (s *SPlaying) getPlayerByID(usreID string) *PlayerHolder {
	for _, v := range s.players {
		if v.userID() == usreID {
			return v
		}
	}
	return nil
}

// getPlayerByChairID 根据chairID获得player对象
func (s *SPlaying) getPlayerByChairID(chairID int) *PlayerHolder {
	for _, v := range s.players {
		if v.chairID == chairID {
			return v
		}
	}
	return nil
}

// firstDiscardedPlayer 第一个出牌者：拥有红桃3者先出牌
func (s *SPlaying) firstDiscardedPlayer() *PlayerHolder {
	return s.room.landlordPlayer()
}

// gameLoop 消息主循环
//  1. 开局处理，起手听之类
//  2. 进入主循环
//      2.0 为某人抽牌
// 		2.1 通知并等待某人出牌等
//      2.2 通知其他人，某人出牌以及等待其他人动作
func (s *SPlaying) gameLoop() {
	s.cl.Println("game loop start")

	// 如果gameLoop中的goroutine出错，则该房间挂死，但是不影响其他房间
	defer func() {
		if r := recover(); r != nil {
			roomExceptionCount++
			debug.PrintStack()
			s.cl.Printf("-----This ROOM will die:%s, Recovered in gameLoop:%v\n", s.room.ID, r)
		}
	}()

	loopKeep := true

	for {
		// 新建上下文
		s.lctx = newLoopContext(s)
		// 发牌
		s.cardMgr.drawForAll()

		// 给所有客户端发牌
		for _, player := range s.players {
			var msgDeal = serializeMsgDeal(s, player)
			player.sendDealMsg(msgDeal)
		}

		// 保存发牌数据
		s.lctx.snapshootDealActions()

		var called bool
		if s.room.config.isCallWithScore {
			// 等待玩家们叫地主，如果没有人叫地主，则重新发牌
			loopKeep, called = s.waitPlayersActionCallWithScore()
			if !loopKeep {
				return
			}
		} else {
			// 等待玩家们叫地主，如果没有人叫地主，则重新发牌
			loopKeep, called = s.waitPlayersActionCall()
			if !loopKeep {
				return
			}
		}

		if !called {
			s.cl.Println("no one call, re-draw for all")
			// 重置发牌管理器
			s.cardMgr = newCardMgr(s.room, s.room.players)
			for _, p := range s.players {
				// 清空手牌，以便重新发牌
				p.resetForNextHand()
			}
		} else {
			break
		}
	}

	robRoundCount := 1
	if s.room.config.isCallWithScore {
		// 这种玩法没有抢地主流程
		robRoundCount = 0
	} else {
		if s.room.is2PlayerRoom() {
			robRoundCount = maxRobLandlordCount / 2
		}
	}

	for ; robRoundCount > 0; robRoundCount-- {
		// 抢地主
		loopKeep, hasRob := s.waitPlayersActionRob()
		if !loopKeep {
			return
		}

		if !hasRob {
			break
		}
	}

	// 给地主3张底牌
	s.awardLandlordLast3()

	if s.room.config.isCallDoubleEnabled {
		// 请求所有人叫加倍
		loopKeep = s.waitPlayersActionCallDouble()
		if !loopKeep {
			return
		}
	}

	// firstDiscardedPlayer 会重设banker玩家
	currentDiscardPlayer := s.firstDiscardedPlayer()
	s.cl.Printf("first discarded player is landlord player:%s:[%d]", currentDiscardPlayer.userID(),
		currentDiscardPlayer.chairID)

	for {
		// 等待玩家出牌或者过
		loopKeep = s.waitPlayerAction(currentDiscardPlayer)
		if !loopKeep {
			break
		}

		var cardHand = s.taskPlayerAction.applyCardHand
		var action = s.taskPlayerAction.replyAction
		s.taskPlayerAction = nil

		// 玩家选择出牌，因此需要考虑其他玩家是否能够针对本次出牌动作，例如吃椪杠胡等
		if action == (ActionType_enumActionType_DISCARD) {
			remain := currentDiscardPlayer.cards.cardCountInHand()
			winAbleRemain := s.room.winAbleRemainCount(currentDiscardPlayer)

			if remain <= winAbleRemain {
				// 玩家打完了手牌，则结束
				s.cl.Printf("player %s remain:%d, <= winAbleRemain:%d, end round", currentDiscardPlayer.userID(),
					remain, winAbleRemain)
				s.onHandWinnerBornSelfDraw(currentDiscardPlayer)
				break
			}

			loopKeep, nextPlayer := s.waitOpponentsAction(currentDiscardPlayer, cardHand)
			if !loopKeep {
				break
			}

			currentDiscardPlayer = nextPlayer

			continue
		}

		s.cl.Panic("game loop should not be here\n")
	}

	if !s.room.isForMonkey && s.lctx.recorder.GetIsHandOver() {
		s.lctx.dump2Redis(s)
	}

	s.lctx = nil
}

// waitPlayersActionCallWithScore 等待玩家叫地主（叫分叫地主）
func (s *SPlaying) waitPlayersActionCallWithScore() (bool, bool) {
	s.cl.Println("waitPlayersActionCallWithScore")
	// 从庄家开始逐个询问
	banker := s.room.bankerPlayer()
	orderPlayers := s.cardMgr.getOrderPlayersWithFirst(banker)

	preScore := 0
	for _, curPlayer := range orderPlayers {
		// var newDraw = s.lctx.isSelfDraw(curPlayer)
		var qaIndex = s.room.nextQAIndex()

		allowFlags := 1

		if curPlayer.cards.has2jokerOr42() {
			allowFlags = (1 << uint(3))
		} else {
			nextScore := 0
			if preScore == 0 {
				nextScore = 0
			} else {
				nextScore = preScore + 1
			}

			for ix := nextScore; ix <= 3; ix++ {
				allowFlags |= (1 << uint(ix))
			}
		}

		// 填写客户端可以操作的动作
		actions := int(ActionType_enumActionType_CallWithScore)

		curPlayer.expectedAction = actions
		msgAllowPlayerAction := serializeMsgAllowedForDiscard(s, curPlayer, actions, qaIndex)
		timeout32 := int32(allowFlags)
		msgAllowPlayerAction.TimeoutInSeconds = &timeout32

		curPlayer.sendActoinAllowedMsg(msgAllowPlayerAction)
		if s.room.isForceConsistent() {
			s.sendMonkeyTips(curPlayer)
		}

		// 发给其他人，其他人也看到此时等待该玩家叫地主
		for _, p := range s.players {
			if p != curPlayer {
				p.sendActoinAllowedMsg(msgAllowPlayerAction)
			}
		}

		s.taskPlayerAction = newTaskPlayerAction(curPlayer, actions)
		s.taskPlayerAction.s = s

		s.taskPlayerAction.allowFlags = allowFlags

		s.taskPlayerAction.actionMsgForRestore = msgAllowPlayerAction

		result := s.taskPlayerAction.wait()

		if !result {
			// 牌局要结束了（例如被解散了）
			return false, false
		}

		taskPlayerAction := s.taskPlayerAction
		s.taskPlayerAction = nil

		action := taskPlayerAction.replyAction
		if action == (ActionType_enumActionType_CallWithScore) {

			selectScore := taskPlayerAction.flags

			s.cl.Printf("player %s:[%d] call-with-score reply, selectScore:%d", curPlayer.userID(), curPlayer.chairID, selectScore)

			// 记录动作
			s.lctx.addActionWithCards(curPlayer, action, nil, qaIndex, pokerface.SRFlags(selectScore))

			// 发送通知给所有客户端
			// 动作者自身需要收到听牌列表更新
			msgActionResultNotify := serializeMsgActionResultNotifyForNoCard((ActionType_enumActionType_CallWithScore), curPlayer)
			cardsInWall32 := int32(selectScore)
			msgActionResultNotify.CardsInWall = &cardsInWall32 // CardsInWall做特殊用途，0表示不叫地主，其他值表示叫了地主

			for _, p := range s.players {
				p.sendActionResultNotify(msgActionResultNotify)
			}

			if selectScore != 0 {
				// 叫了地主
				s.room.markup = selectScore
				s.cl.Printf("player %s:[%d] call-with-score %d, room markup changed", curPlayer.userID(),
					curPlayer.chairID, selectScore)

				// 记录下最高叫分
				preScore = selectScore
				s.room.landlordUserID = curPlayer.userID()
			}

			if selectScore == 3 {
				// 玩家叫了地主，叫到最高分，则不再继续流程
				return true, true
			}

			// 玩家不叫地主
			curPlayer.hStatis.isSkipCall = true
		}
	}

	if preScore == 0 {
		// 没有人叫地主, 强制庄家为地址
		s.room.landlordUserID = banker.userID()

		s.cl.Printf("no body call-with-score, force to banker user:%s", s.room.landlordUserID)

		// 发送通知给所有客户端
		// 动作者自身需要收到听牌列表更新
		msgActionResultNotify := serializeMsgActionResultNotifyForNoCard((ActionType_enumActionType_CallWithScore), banker)
		cardsInWall32 := int32(1)
		msgActionResultNotify.CardsInWall = &cardsInWall32 // CardsInWall做特殊用途，0表示不叫地主，其他值表示叫了地主

		for _, p := range s.players {
			p.sendActionResultNotify(msgActionResultNotify)
		}
	}

	return true, true
}

// waitPlayersActionCall 等待玩家叫地主
func (s *SPlaying) waitPlayersActionCall() (bool, bool) {
	s.cl.Println("waitPlayersActionCall")
	// 从庄家开始逐个询问
	orderPlayers := s.cardMgr.getOrderPlayersWithFirst(s.room.bankerPlayer())

	for _, curPlayer := range orderPlayers {
		// var newDraw = s.lctx.isSelfDraw(curPlayer)
		var qaIndex = s.room.nextQAIndex()

		// 填写客户端可以操作的动作
		actions := int(ActionType_enumActionType_Call)

		curPlayer.expectedAction = actions
		msgAllowPlayerAction := serializeMsgAllowedForDiscard(s, curPlayer, actions, qaIndex)
		curPlayer.sendActoinAllowedMsg(msgAllowPlayerAction)
		if s.room.isForceConsistent() {
			s.sendMonkeyTips(curPlayer)
		}

		// 发给其他人，其他人也看到此时等待该玩家叫地主
		for _, p := range s.players {
			if p != curPlayer {
				p.sendActoinAllowedMsg(msgAllowPlayerAction)
			}
		}

		s.taskPlayerAction = newTaskPlayerAction(curPlayer, actions)
		s.taskPlayerAction.s = s

		result := s.taskPlayerAction.wait()

		if !result {
			// 牌局要结束了（例如被解散了）
			return false, false
		}

		taskPlayerAction := s.taskPlayerAction
		s.taskPlayerAction = nil

		action := taskPlayerAction.replyAction
		if action == (ActionType_enumActionType_Call) {

			flags := pokerface.SRFlags_SRNone
			if taskPlayerAction.flags != 0 {
				flags = pokerface.SRFlags_SRRichi
			}

			s.cl.Printf("player %s:[%d] Call reply, flag:%d", curPlayer.userID(), curPlayer.chairID, flags)

			// 记录动作
			s.lctx.addActionWithCards(curPlayer, action, nil, qaIndex, flags)

			// 发送通知给所有客户端
			// 动作者自身需要收到听牌列表更新
			msgActionResultNotify := serializeMsgActionResultNotifyForNoCard((ActionType_enumActionType_Call), curPlayer)
			cardsInWall32 := int32(3)
			if flags != pokerface.SRFlags_SRNone {
				cardsInWall32 = int32(1003)
			}
			msgActionResultNotify.CardsInWall = &cardsInWall32 // CardsInWall做特殊用途，如果客户端收到1003表明玩家叫了地主

			for _, p := range s.players {
				p.sendActionResultNotify(msgActionResultNotify)
			}

			if flags != pokerface.SRFlags_SRNone {
				// 玩家叫了地主
				s.room.landlordUserID = curPlayer.userID()
				return true, true
			}

			// 玩家不叫地主
			curPlayer.hStatis.isSkipCall = true
		}
	}

	// 没有人叫地主
	return true, false
}

// waitPlayersActionCallDouble 等待玩家加注
func (s *SPlaying) waitPlayersActionCallDouble() bool {
	var qaIndex = s.room.nextQAIndex()

	s.taskCallDoulbe = newTaskCallDouble(s)
	taskCallDoulbe := s.taskCallDoulbe

	action := int(ActionType_enumActionType_CallDouble)
	for _, p := range s.players {
		msgAllowedAction := serializeMsgAllowedForDiscard(s, p, action, qaIndex)
		p.expectedAction = action
		p.sendActoinAllowedMsg(msgAllowedAction)

		if s.room.isForceConsistent() {
			s.sendMonkeyTips(p)
		}
	}

	result := taskCallDoulbe.wait()
	if !result {
		return false
	}

	// 立即重置
	s.taskCallDoulbe = nil

	// 记录玩家加注行为
	for e := taskCallDoulbe.waitQueue.Front(); e != nil; e = e.Next() {
		wi := e.Value.(*TaskCallDoubleQueueItem)

		flags := pokerface.SRFlags_SRNone
		if wi.replyFlags != 0 {
			flags = pokerface.SRFlags_SRRichi
		}

		s.lctx.addActionWithCards(wi.player,
			ActionType_enumActionType_CallDouble, nil, qaIndex,
			flags)

		s.cl.Printf("player %s chairID %d call double, flags:%d",
			wi.player.userID(), wi.player.chairID, flags)
	}

	return true
}

// waitPlayersActionRob 等待非地主玩家抢地主
func (s *SPlaying) waitPlayersActionRob() (bool, bool) {
	// 从非地主玩家开始逐个询问
	landlord := s.room.landlordPlayer()
	orderPlayers := s.cardMgr.getOrderPlayers(landlord)
	// 地主也添加到列表中，以便别人抢了地主后，其还可以抢一次
	orderPlayers = append(orderPlayers, landlord)
	hasRob := false
	for _, curPlayer := range orderPlayers {

		// 自己已经是地主，就不要抢了
		if curPlayer.userID() == s.room.landlordUserID {
			continue
		}

		// 不叫地主的玩家没有机会抢地主
		if curPlayer.hStatis.isSkipCall {
			continue
		}

		// var newDraw = s.lctx.isSelfDraw(curPlayer)
		var qaIndex = s.room.nextQAIndex()

		// 填写客户端可以操作的动作
		actions := int(ActionType_enumActionType_Rob)

		curPlayer.expectedAction = actions
		msgAllowPlayerAction := serializeMsgAllowedForDiscard(s, curPlayer, actions, qaIndex)
		curPlayer.sendActoinAllowedMsg(msgAllowPlayerAction)
		if s.room.isForceConsistent() {
			s.sendMonkeyTips(curPlayer)
		}

		// 发给其他人，其他人也看到此时等待该玩家叫地主
		for _, p := range s.players {
			if p != curPlayer {
				p.sendActoinAllowedMsg(msgAllowPlayerAction)
			}
		}

		s.taskPlayerAction = newTaskPlayerAction(curPlayer, actions)
		s.taskPlayerAction.s = s

		result := s.taskPlayerAction.wait()

		if !result {
			// 牌局要结束了（例如被解散了）
			return false, false
		}

		taskPlayerAction := s.taskPlayerAction
		s.taskPlayerAction = nil

		action := taskPlayerAction.replyAction
		if action == (ActionType_enumActionType_Rob) {

			flags := pokerface.SRFlags_SRNone
			if taskPlayerAction.flags != 0 {
				flags = pokerface.SRFlags_SRRichi
			}

			s.cl.Printf("player %s:[%d] Rob reply, flag:%d", curPlayer.userID(), curPlayer.chairID, flags)

			// 记录动作
			s.lctx.addActionWithCards(curPlayer, action, nil, qaIndex, flags)

			// 发送通知给所有客户端
			// 动作者自身需要收到听牌列表更新
			msgActionResultNotify := serializeMsgActionResultNotifyForNoCard((ActionType_enumActionType_Rob), curPlayer)
			cardsInWall32 := int32(3)
			if flags != pokerface.SRFlags_SRNone {
				// 玩家抢了地主
				s.room.landlordUserID = curPlayer.userID()
				s.room.markup = s.room.markup * 2 // 倍数翻倍

				cardsInWall32 = int32(1003 + s.room.markup*10) // cardsInWall32 携带了当前最新的倍数给所有客户端

				hasRob = true
			} else {
				hasRob = false
			}

			msgActionResultNotify.CardsInWall = &cardsInWall32 // CardsInWall做特殊用途

			for _, p := range s.players {
				p.sendActionResultNotify(msgActionResultNotify)
			}
		}
	}

	return true, hasRob
}

// awardLandlordLast3 给予地主玩家3张底牌
func (s *SPlaying) awardLandlordLast3() {
	// 抽剩余3张牌
	landlordPlayer := s.room.landlordPlayer()
	cards := make([]*Card, 0, 3)
	for i := 0; i < 3; i++ {
		ok, c := s.cardMgr.drawForPlayer(landlordPlayer, false)
		if !ok {
			s.cl.Panicln("awardLandlordLast3, falied to draw last 3")
		}

		cards = append(cards, c)
	}

	// 保存一下最后三张，以便掉线恢复
	s.lastAwardCards = make([]int32, len(cards))
	for i, c := range cards {
		s.lastAwardCards[i] = int32(c.cardID)
	}

	s.lctx.addActionWithRawCards(landlordPlayer, ActionType_enumActionType_DRAW, cards, s.room.qaIndex)

	msgActionResultNotify := serializeMsgActionResultNotifyForDraw(landlordPlayer, cards)
	for _, p := range s.players {
		p.sendActionResultNotify(msgActionResultNotify)
	}
}

// onHandWashout 流局处理
func (s *SPlaying) onHandWashout() {
	calcFinalResultWashout(s)

	var msgHandOver = serializeMsgHandOverWashout(s)
	s.room.appendHandScoreRecord(int(HandOverType_enumHandOverType_None))

	s.lctx.finishHandWashout(msgHandOver.Scores)

	// 庄家、风圈计算；切换到Starting状态
	s.onHandOver(true)
	s.room.onHandOver(msgHandOver)
}

// onHandWinnerBornSelfDraw 自摸胡牌处理，发送自摸胡牌结果以及计分等
func (s *SPlaying) onHandWinnerBornSelfDraw(winner *PlayerHolder) {
	s.cl.Println("onHandWinnerBornSelfDraw")

	s.lctx.addActionWithCards(winner, ActionType_enumActionType_Win_SelfDrawn, nil, s.room.qaIndex, pokerface.SRFlags_SRNone)

	calcFinalResultSelfDraw(s, winner)
	var msgHandOver = serializeMsgHandOver(s, int(HandOverType_enumHandOverType_Win_SelfDrawn))

	s.room.appendHandScoreRecord(int(HandOverType_enumHandOverType_Win_SelfDrawn))

	// 记录用户的全局统计数据
	s.hStatis2GStatis()

	// 记录得分
	s.lctx.finishWinnerBorn(msgHandOver.Scores)

	// 庄家、风圈计算；切换到Starting状态
	s.onHandOver(false)
	s.room.onHandOver(msgHandOver)
}

// hStatis2GStatis 记录用户的全局统计数据
func (s *SPlaying) hStatis2GStatis() {
	// 并非流局，则记录所有人的本手牌数据
	for _, p := range s.players {
		// p.gStatis.roundScore += (p.sctx.totalWinScore)

		// if p.sctx.winType == int(HandOverType_enumHandOverType_Win_SelfDrawn) {
		// 	p.gStatis.winSelfDrawnCounter++
		// }

		if p.sctx.calcTotalWinScore() > 0 {
			p.gStatis.winSelfDrawnCounter++
		}

		// // 新疆麻将修改：大胡次数
		// if p.sctx.greatWinCount > 0 {
		// 	p.gStatis.greatWinCounter += p.sctx.greatWinCount
		// }

		// // 新疆麻将修改：小胡次数
		// if p.sctx.miniWinCount > 0 {
		// 	p.gStatis.miniWinCounter += p.sctx.miniWinCount
		// }

		// // 新疆麻将修改：包庄次数
		// if p.sctx.isPayForAll {
		// 	p.gStatis.kongerCounter++
		// }
	}
}

// onHandOver 一手牌结束，不管是流局还是胡牌的，都会调用此函数
func (s *SPlaying) onHandOver(washout bool) {
	s.cl.Println("onHandOver")

	s.chooseBanker(washout)

}

// chooseBanker 选择下一手牌的庄家
// 大丰关张扑克牌没有切换庄家，第一个发牌者以拥有红桃3者为准
func (s *SPlaying) chooseBanker(washout bool) {
	var newBanker *PlayerHolder
	switch s.room.config.nextBankerMethod {
	case 0:
		for _, p := range s.players {
			if p.sctx.isWin() {
				newBanker = p
			}
		}
	case 1:
		idx := s.room.rand.Intn(len(s.players))
		newBanker = s.players[idx]
	case 2:
		curBanker := s.room.bankerPlayer()
		newBanker = s.cardMgr.rightOpponent(curBanker)
	}

	s.room.bankerChange2(newBanker)
}

// waitPlayerAction 等待当前玩家出牌，或者加杠，暗杠，自摸胡牌
func (s *SPlaying) waitPlayerAction(curPlayer *PlayerHolder) bool {
	// var newDraw = s.lctx.isSelfDraw(curPlayer)
	var qaIndex = s.room.nextQAIndex()

	// 填写客户端可以操作的动作
	actions := int(ActionType_enumActionType_DISCARD)

	curPlayer.expectedAction = actions
	msgAllowPlayerAction := serializeMsgAllowedForDiscard(s, curPlayer, actions, qaIndex)
	curPlayer.sendActoinAllowedMsg(msgAllowPlayerAction)
	if s.room.isForceConsistent() {
		s.sendMonkeyTips(curPlayer)
	}

	// 其他玩家只看到当前玩家在出牌，不能看到胡牌杠牌等动作
	msgAllowPlayerAction2 := serializeMsgAllowedForDiscard2Opponent(curPlayer, qaIndex)
	for _, p := range s.players {
		if p != curPlayer {
			p.sendActoinAllowedMsg(msgAllowPlayerAction2)
		}
	}

	s.taskPlayerAction = newTaskPlayerAction(curPlayer, actions)
	s.taskPlayerAction.s = s

	result := s.taskPlayerAction.wait()

	if !result {
		return false
	}

	action := s.taskPlayerAction.replyAction
	if action == (ActionType_enumActionType_DISCARD) {
		// 出牌，从手牌队列移到出牌队列
		s.taskPlayerAction.applyCardHand = s.cardMgr.playerDiscard(curPlayer, s.taskPlayerAction.replyCardHand)
		// 记录出牌数据以及重置一些用于限制玩家动作的开关变量（每次玩家出牌后，这些开关变量都重置）
		curPlayer.hStatis.resetLocked()

		flags := pokerface.SRFlags_SRNone

		// 记录动作
		s.lctx.addActionWithCards(curPlayer, action, s.taskPlayerAction.replyCardHand, qaIndex, flags)

		// 发送出牌结果通知给所有客户端
		// 动作者自身需要收到听牌列表更新
		msgActionResultNotify := serializeMsgActionResultNotifyForDiscardedCard((ActionType_enumActionType_DISCARD), curPlayer,
			s.taskPlayerAction.replyCardHand)

		for _, p := range s.players {
			p.sendActionResultNotify(msgActionResultNotify)
		}
	}
	return true
}

// waitOpponentsAction 等待其他玩家对本次出牌的响应
// 函数会逐个请求其他玩家是否跟，如果没有人跟，则本次出牌玩家继续出牌，如果有人跟，则最后一个跟牌者出牌
func (s *SPlaying) waitOpponentsAction(curDiscardPlayer *PlayerHolder, discardedCardHand *CardHand) (bool, *PlayerHolder) {

	prevDiscardedPlayer := curDiscardPlayer
	prevDiscaredCardHand := discardedCardHand

	// 外层循环表示打到其他人都打不起为止
	// 最终没有人跟的那个出牌者，会回到gameLoop中继续出自由牌
	keepGO := true
	for keepGO {
		// 先置为false，如果本轮循环没有人出牌，则终止
		keepGO = false
		orderPlayers := s.cardMgr.getOrderPlayers(prevDiscardedPlayer)

		// 依次询问其他人是否要跟
		for _, opponent := range orderPlayers {
			taskDiscardReAction := analyseTaskPlayerReAction(s, opponent, prevDiscardedPlayer, prevDiscaredCardHand)
			if taskDiscardReAction == nil {
				// 发送SKIP给所有人
				msgActionResultNotify := serializeMsgActionResultNotifyForNoCard((ActionType_enumActionType_SKIP), opponent)
				for _, p := range s.players {
					p.sendActionResultNotify(msgActionResultNotify)
				}
				continue
			}

			s.taskDiscardReAction = taskDiscardReAction

			qaIndex := s.room.nextQAIndex()
			msgAllowedAction := serializeMsgAllowedForDiscardResponse(opponent, qaIndex, prevDiscaredCardHand, prevDiscardedPlayer)
			opponent.sendReActoinAllowedMsg(msgAllowedAction)
			if s.room.isForceConsistent() {
				s.sendMonkeyTips(opponent)
			}

			// 通知其他人，正在等待该玩家打牌
			msgAllowPlayerAction2 := serializeMsgAllowedForDiscard2Opponent(opponent, qaIndex)
			for _, p := range s.players {
				if p != opponent {
					p.sendActoinAllowedMsg(msgAllowPlayerAction2)
				}
			}

			s.taskDiscardReAction.s = s
			result := s.taskDiscardReAction.wait()
			if !result {
				s.taskDiscardReAction = nil
				return false, nil
			}
			s.taskDiscardReAction = nil

			replyAction := taskDiscardReAction.replyAction
			if replyAction == ActionType_enumActionType_SKIP {
				s.lctx.addActionWithCards(opponent, ActionType_enumActionType_SKIP, nil, qaIndex, pokerface.SRFlags_SRNone)

				// 发送SKIP给所有人
				msgActionResultNotify := serializeMsgActionResultNotifyForNoCard((ActionType_enumActionType_SKIP), opponent)
				for _, p := range s.players {
					p.sendActionResultNotify(msgActionResultNotify)
				}
				continue
			} else if replyAction == (ActionType_enumActionType_DISCARD) {

				s.applyActionForOpponents(taskDiscardReAction, qaIndex)

				remain := opponent.cards.cardCountInHand()
				winAbleRemain := s.room.winAbleRemainCount(opponent)

				if remain <= winAbleRemain {
					s.cl.Printf("player %s remain:%d, <= winAbleRemain:%d, end round", opponent.userID(),
						remain, winAbleRemain)
					s.onHandWinnerBornSelfDraw(opponent)
					return false, nil
				}

				prevDiscardedPlayer = opponent
				prevDiscaredCardHand = taskDiscardReAction.applyCardHand

				// 有人出牌了，退出内部循环，继续外部循环
				keepGO = true
				break
			}
		}
	}

	return true, prevDiscardedPlayer
}

// applyActionForOpponents 对于玩家的响应动作，例如其选择碰，需要落实碰牌操作：修改牌列表等
// 对于玩家的响应动作不能放在handler代码中处理的原因是， 吃椪杠胡有优先级，假如一个玩家选择了碰，而
// 另一个玩家选择了胡，则胡者优先，因此我们需要cache玩家的选择，等待所有人都回复了（或者优先级最高者回复了）
// 再由本函数来落实玩家的选择。优先级处理见TaskDiscardReAction类代码
func (s *SPlaying) applyActionForOpponents(taskDiscardReAction *TaskPlayerReAction, qaIndex int) {
	curPlayer := taskDiscardReAction.waitPlayer

	// 出牌，从手牌队列移到出牌队列
	taskDiscardReAction.applyCardHand = s.cardMgr.playerDiscard(curPlayer, taskDiscardReAction.replyMsgCardHand)
	// 记录出牌数据以及重置一些用于限制玩家动作的开关变量（每次玩家出牌后，这些开关变量都重置）
	curPlayer.hStatis.resetLocked()

	flags := pokerface.SRFlags_SRNone

	// 记录动作
	s.lctx.addActionWithCards(curPlayer, taskDiscardReAction.replyAction, taskDiscardReAction.replyMsgCardHand, qaIndex, flags)

	// 发送出牌结果通知给所有客户端
	// 动作者自身需要收到听牌列表更新
	msgActionResultNotify := serializeMsgActionResultNotifyForDiscardedCard((ActionType_enumActionType_DISCARD), curPlayer,
		taskDiscardReAction.replyMsgCardHand)

	for _, p := range s.players {
		p.sendActionResultNotify(msgActionResultNotify)
	}
}

// 发送monkey提示信息给客户端，以便客户端能够进行正确的操作选择
func (s *SPlaying) sendMonkeyTips(player *PlayerHolder) {
	monkeyUserCardsCfg := s.room.monkeyCfg.getMonkeyUserCardsCfg(player.userID())

	if monkeyUserCardsCfg == nil {
		return
	}

	actionTips := monkeyUserCardsCfg.actionTips
	if len(actionTips) <= player.hStatis.actionCounter {
		return
	}

	tips := actionTips[player.hStatis.actionCounter]

	player.sendTipsString(tips)
}

func sendTips2Player(player *PlayerHolder, tipCode pokerface.TipCode) {
	msgTips := &pokerface.MsgRoomShowTips{}
	var tipCode32 = int32(tipCode)
	msgTips.TipCode = &tipCode32

	player.sendMsg(msgTips, int32(pokerface.MessageCode_OPRoomShowTips))
}
