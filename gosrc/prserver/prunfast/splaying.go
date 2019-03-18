package prunfast

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
	//taskFirstReadyHand  *TaskFirstReadyHand
	cardMgr *CardMgr

	players []*PlayerHolder
	room    *Room
	lctx    *LoopContext

	cl *logrus.Entry
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
	if s.room.config.playerNumAcquired < 3 {
		// 两人房间，就选第一个人好作为发牌者，依次轮流
		return s.room.bankerPlayer()
	}

	for _, p := range s.players {
		if p.cards.hasCardInHand(R3H) {
			// 记录一下第一个出牌者标志
			p.hStatis.isFirstDiscarded = true
			// 第一个出牌者算做庄家
			s.room.bankerChange2(p)
			return p
		}
	}

	s.cl.Panicln("firstDiscardedPlayer, nobody has R3H")
	return nil
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

	// 发牌
	s.cardMgr.drawForAll()

	// 给所有客户端发牌
	for _, player := range s.players {
		var msgDeal = serializeMsgDeal(s, player)
		player.sendDealMsg(msgDeal)
	}

	s.lctx = newLoopContext(s)

	loopKeep := true
	// firstDiscardedPlayer 会重设banker玩家
	currentDiscardPlayer := s.firstDiscardedPlayer()
	// 保存发牌数据
	s.lctx.snapshootDealActions()

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
			if remain == 0 {
				// 玩家打完了手牌，则结束
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

		if p.sctx.winType == int(HandOverType_enumHandOverType_Win_SelfDrawn) {
			p.gStatis.winSelfDrawnCounter++
		}

		// 新疆麻将修改：大胡次数
		if p.sctx.greatWinCount > 0 {
			p.gStatis.greatWinCounter += p.sctx.greatWinCount
		}

		// 新疆麻将修改：小胡次数
		if p.sctx.miniWinCount > 0 {
			p.gStatis.miniWinCounter += p.sctx.miniWinCount
		}

		// 新疆麻将修改：包庄次数
		if p.sctx.isPayForAll {
			p.gStatis.kongerCounter++
		}
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
	// 两个人的时候，轮流打牌
	if s.room.config.playerNumAcquired < 3 {
		curBanker := s.room.bankerPlayer()
		newBanker := s.cardMgr.rightOpponent(curBanker)

		s.room.bankerChange2(newBanker)
	}
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
				if remain == 0 {
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
