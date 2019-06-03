package zjmahjong

import (
	"mahjong"
	"runtime/debug"

	"github.com/sirupsen/logrus"

	"github.com/golang/protobuf/proto"
)

// SPlaying 正在游戏状态
type SPlaying struct {
	taskDiscardReAction *TaskPlayerReAction
	taskPlayerAction    *TaskPlayerAction
	taskFirstReadyHand  *TaskFirstReadyHand
	tileMgr             *TileMgr

	players []*PlayerHolder
	room    *Room
	lctx    *LoopContext
	cl      *logrus.Entry
}

// newSPlaying 新建playing 状态机
func newSPlaying(room *Room) *SPlaying {
	s := &SPlaying{}
	s.room = room
	s.players = room.players
	s.tileMgr = newTileMgr(room, room.players)
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
	player.state = mahjong.PlayerState_PSOffline
	s.room.updateRoomInfo2All()
	//s.room.writePlayerLeave2Redis(player, false)
}

// onPlayerReEnter 处理玩家重入事件，主要是掉线恢复
func (s *SPlaying) onPlayerReEnter(player *PlayerHolder) {
	player.state = mahjong.PlayerState_PSPlaying
	s.room.updateRoomInfo2All()

	// 掉线恢复
	// 先发送牌数据
	msgRestore := serializeMsgRestore(s, player)
	player.sendMsg(msgRestore, int32(mahjong.MessageCode_OPRestore))

	// 根据当前的gameLoop的等待状态，给玩家重新发送最近一个消息
	if s.taskFirstReadyHand != nil {
		s.taskFirstReadyHand.onPlayerRestore(player)
	} else if s.taskPlayerAction != nil {
		s.taskPlayerAction.onPlayerRestore(player)
	} else if s.taskDiscardReAction != nil {
		s.taskDiscardReAction.onPlayerRestore(player)
	}
}

// getStateConst state const
func (s *SPlaying) getStateConst() mahjong.RoomState {
	return mahjong.RoomState_SRoomPlaying
}

func (s *SPlaying) onStateEnter() {
	s.cl.Println("SPlaying enter")

	for _, p := range s.players {
		p.state = mahjong.PlayerState_PSPlaying
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
	// for _, p := range s.players {
	// 	p.resetForNextHand()
	// }

	// 需要把所有正在等待的Task cancel掉
	// 否则gameLoop的go routine不能退出
	s.cl.Println("SPlaying leave")
	if s.taskFirstReadyHand != nil {
		s.taskFirstReadyHand.cancel()
	}

	if s.taskPlayerAction != nil {
		s.taskPlayerAction.cancel()
	}

	if s.taskDiscardReAction != nil {
		s.taskDiscardReAction.cancel()
	}
}

func (s *SPlaying) onMessage(iu IUser, gmsg *mahjong.GameMessage) {
	var msgCode = mahjong.MessageCode(gmsg.GetOps())

	switch msgCode {
	case mahjong.MessageCode_OPAction:
		actionMsg := &mahjong.MsgPlayerAction{}
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

func (s *SPlaying) onActionMessage(iu IUser, msg *mahjong.MsgPlayerAction) {
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

	var action = mahjong.ActionType(msg.GetAction())
	switch action {
	case mahjong.ActionType_enumActionType_CHOW:
		onMessageChow(s, player, msg)
		break
	case mahjong.ActionType_enumActionType_SKIP:
		onMessageSkip(s, player, msg)
		break
	case mahjong.ActionType_enumActionType_PONG:
		onMessagePong(s, player, msg)
		break
	case mahjong.ActionType_enumActionType_KONG_Exposed:
		onMessageKongExposed(s, player, msg)
		break
	case mahjong.ActionType_enumActionType_KONG_Concealed:
		onMessageKongConcealed(s, player, msg)
		break
	case mahjong.ActionType_enumActionType_WIN_Chuck:
		onMessageWinChuck(s, player, msg)
		break
	case mahjong.ActionType_enumActionType_DISCARD:
		onMessageDiscard(s, player, msg)
		break
	case mahjong.ActionType_enumActionType_WIN_SelfDrawn:
		onMessageWinSelfDraw(s, player, msg)
		break
	case mahjong.ActionType_enumActionType_DRAW:
		s.cl.Println("OnActionMessage unsupported DRAW from client")
		break
	case mahjong.ActionType_enumActionType_FirstReadyHand:
		onMessageFirstReadyHand(s, player, msg)
		break
	case mahjong.ActionType_enumActionType_KONG_Triplet2:
		onMessageTriplet2Kong(s, player, msg)
		break
	case mahjong.ActionType_enumActionType_CustomB:
		onMessageShouldFinalDraw(s, player, msg)
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

	// 发牌，注意庄家此时是获得14张牌
	s.tileMgr.drawForAll()

	// 给所有客户端发牌
	for _, player := range s.players {
		var msgDeal = serializeMsgDeal(s, player)
		player.sendDealMsg(msgDeal)
	}

	s.lctx = newLoopContext(s)
	// 保存发牌数据
	s.lctx.snapshootDealActions()

	// 起手听牌处理
	loopKeep := s.firstReadyHand()
	if !loopKeep {
		return
	}

	currentDiscardPlayer := s.room.bankerPlayer()
	bankerFirstAction := true
	for {

		var needDraw = false
		if currentDiscardPlayer.tiles.agariTileCount() < 14 {
			needDraw = true
		}

		// 玩家抽牌
		if needDraw {
			// drawForPlayer以参数true调用表示保留最后一张
			drawOK := s.drawForPlayer(currentDiscardPlayer, true)
			if !drawOK && s.tileMgr.tileCountInWall() == 1 {
				loopKeep := s.waitFinalDrawQuery(currentDiscardPlayer)
				if !loopKeep {
					break
				}

				// drawForPlayer以参数true调用表示不再保留最后一张
				drawOK = s.drawForPlayer(currentDiscardPlayer, false)
			}

			if !drawOK {
				// 流局
				s.onHandWashout()
				//s.lctx.finishHandWashout()
				break
			}
		}

		// 庄家需要特殊处理，因为庄家发牌时有14张牌而不需要摸牌
		// 但是逻辑上却视为摸牌
		if bankerFirstAction {
			bankerFirstAction = false
			needDraw = true
		}

		// 等待玩家胡牌或者出牌
		loopKeep, tile, action := s.waitPlayerAction(currentDiscardPlayer, needDraw)
		if !loopKeep {
			break
		}

		// 如果玩家胡牌（自摸胡），结束循坏
		if action == int(mahjong.ActionType_enumActionType_WIN_SelfDrawn) {
			s.onHandWinnerBornSelfDraw(currentDiscardPlayer)
			break
		}

		// 玩家选择出牌，因此需要考虑其他玩家是否能够针对本次出牌动作，例如吃椪杠胡等
		if action == int(mahjong.ActionType_enumActionType_DISCARD) {
			loopKeep, nextPlayer := s.waitOpponentsAction(currentDiscardPlayer, tile)
			if !loopKeep {
				break
			}

			currentDiscardPlayer = nextPlayer
			continue
		}

		// 加杠，需要考虑抢杠胡
		if action == int(mahjong.ActionType_enumActionType_KONG_Triplet2) {
			if s.taskDiscardReAction != nil {
				loopKeep := s.waitOpponentsRobKong(currentDiscardPlayer, tile)
				if !loopKeep {
					break
				}
			}
			continue
		}

		// 玩家选择暗杠，需要补牌
		if action == int(mahjong.ActionType_enumActionType_KONG_Concealed) {
			continue
		}

		s.cl.Panic("game loop should not be here\n")
	}

	if !s.room.isForMonkey && s.lctx.recorder.GetIsHandOver() {
		s.lctx.dump2Redis(s)
	}

	s.lctx = nil
}

// firstReadyHand 处理起手听牌
// 这个函数相当于把庄家过滤掉，因为readyHandAble要求牌数一定是13
// 由于庄家发牌的那一刻是发了14张牌，因此，庄家不可能readyHandAble
func (s *SPlaying) firstReadyHand() bool {
	count := 0
	for _, p := range s.players {
		tiles := p.tiles
		if tiles.readyHandAble() {
			count++
		}
	}

	if count < 1 {
		// 没有起手听牌
		return true
	}

	readyHandPlayers := make([]*PlayerHolder, 0, count)
	qaIndex := s.room.nextQAIndex()

	for _, p := range s.players {
		if !p.tiles.readyHandAble() {
			continue
		}
		msgAllowedAction := serializeMsgAllowedForRichi(s, p, qaIndex)
		p.expectedAction = int(msgAllowedAction.GetAllowedActions())
		p.sendActoinAllowedMsg(msgAllowedAction)

		if s.room.isForceConsistent() {
			s.sendMonkeyTips(p)
		}

		readyHandPlayers = append(readyHandPlayers, p)
	}

	s.taskFirstReadyHand = newTaskFirstReadyHand(readyHandPlayers)
	s.taskFirstReadyHand.s = s

	// 发送一个tips给庄家
	// sendTips2Player(s.room.bankerPlayer(), TipCode_TCWaitOpponentsAction)

	result := s.taskFirstReadyHand.wait()

	actions := mahjong.ActionType_enumActionType_FirstReadyHand | mahjong.ActionType_enumActionType_SKIP
	// 记录动作
	for _, wi := range s.taskFirstReadyHand.waitQueue {
		flags := mahjong.SRFlags_SRNone
		if wi.isReply && wi.isRichi {
			flags = mahjong.SRFlags_SRRichi
		}

		s.lctx.addActionWithTile(wi.player, InvalidTile.tileID, 0, mahjong.ActionType_enumActionType_FirstReadyHand, qaIndex, flags,
			int(actions))
	}

	s.taskFirstReadyHand = nil
	return result
}

// drawForPlayer 为玩家抽取一张非花牌的手牌，如果没牌可抽了就返回false
func (s *SPlaying) drawForPlayer(player *PlayerHolder, reserveLast bool) bool {
	ok, hand, flowers := s.tileMgr.drawForPlayer(player, true, reserveLast)

	// 如果有花牌产生，先发花牌
	// 注意到，如果最后没有非花牌可抽，那么客户端得到一张EmptyTile，客户端需要做判断和相应处理
	// 例如客户端不能把这个手牌加入手牌列表
	var newTile = EmptyTile
	if ok {
		newTile = hand
	}

	// 记录抽牌动作
	s.lctx.addDrawAction(player, newTile.tileID, flowers, s.room.qaIndex)

	msgActionResultNotify := serializeMsgActionResultNotifyForDraw(player, newTile.tileID, flowers, s.tileMgr.tileCountInWall())
	player.sendActionResultNotify(msgActionResultNotify)

	if ok {
		// 清空手牌标志
		var tileMark = int32(TILEMAX)
		msgActionResultNotify.ActionTile = &tileMark
	}

	for _, p := range s.players {
		if p != player {
			p.sendActionResultNotify(msgActionResultNotify)
		}
	}
	return ok
}

// onHandWashout 流局处理
func (s *SPlaying) onHandWashout() {
	var msgHandOver = serializeMsgHandOverWashout(s)

	s.room.appendHandScoreRecord(int(mahjong.HandOverType_enumHandOverType_None))

	if !s.room.isUlimitRound {
		s.lctx.finishHandWashout()
	}

	// 庄家、风圈计算；切换到Starting状态
	s.onHandOver(true)
	s.room.onHandOver(msgHandOver)
}

// onHandWinnerBornSelfDraw 自摸胡牌处理，发送自摸胡牌结果以及计分等
func (s *SPlaying) onHandWinnerBornSelfDraw(winner *PlayerHolder) {
	s.cl.Println("onHandWinnerBornSelfDraw")
	calcFinalResultSelfDraw(s, winner)
	var msgHandOver = serializeMsgHandOver(s, int(mahjong.HandOverType_enumHandOverType_Win_SelfDrawn))

	s.room.appendHandScoreRecord(int(mahjong.HandOverType_enumHandOverType_Win_SelfDrawn))

	// 记录用户的全局统计数据
	s.hStatis2GStatis()

	// 记录得分
	if !s.room.isUlimitRound {
		s.lctx.finishWinnerBorn(msgHandOver.Scores)
	}

	// 庄家、风圈计算；切换到Starting状态
	s.onHandOver(false)
	s.room.onHandOver(msgHandOver)
}

// onHandWinnerBornWinChuck 吃铳胡牌处理，发送吃铳胡牌结果以及计分等
func (s *SPlaying) onHandWinnerBornWinChuck(chucker *PlayerHolder) {
	s.cl.Println("onHandWinnerBornWinChuck")
	calcFinalResultWithChucker(s, chucker)

	var msgHandOver = serializeMsgHandOver(s, int(mahjong.HandOverType_enumHandOverType_Win_Chuck))

	s.room.appendHandScoreRecord(int(mahjong.HandOverType_enumHandOverType_Win_Chuck))

	// 记录用户的全局统计数据
	s.hStatis2GStatis()

	// 记录得分
	if !s.room.isUlimitRound {
		s.lctx.finishWinnerBorn(msgHandOver.Scores)
	}

	// 庄家、风圈计算；切换到Waiting状态
	s.onHandOver(false)
	s.room.onHandOver(msgHandOver)
}

// hStatis2GStatis 记录用户的全局统计数据
func (s *SPlaying) hStatis2GStatis() {
	// 并非流局，则记录所有人的本手牌数据
	for _, p := range s.players {
		// p.gStatis.roundScore += (p.sctx.totalWinScore)

		if p.sctx.winType == int(mahjong.HandOverType_enumHandOverType_Win_SelfDrawn) {
			p.gStatis.winSelfDrawnCounter++
		} else if p.sctx.winType == int(mahjong.HandOverType_enumHandOverType_Win_Chuck) {
			p.gStatis.winChuckCounter++
		} else if p.sctx.winType == int(mahjong.HandOverType_enumHandOverType_Chucker) {
			p.gStatis.chuckerCounter++
		}
	}
}

// onHandOver 一手牌结束，不管是流局还是胡牌的，都会调用此函数
func (s *SPlaying) onHandOver(washout bool) {
	s.cl.Println("onHandOver")

	// for _, p := range s.players {
	// 	p.hStatis.reset()
	// }

	s.chooseBanker(washout)

	// if !s.room.isForMonkey {
	// 	s.lctx.dump2Redis(s)
	// }
}

// chooseBanker 选择下一手牌的庄家
func (s *SPlaying) chooseBanker(washout bool) {
	curBanker := s.room.bankerPlayer()

	for _, p := range s.players {
		p.gStatis.isContinuousBanker = false
	}

	var newBanker *PlayerHolder
	// 流局，庄家继续做庄
	if washout {
		newBanker = curBanker
	} else {
		newBanker = s.tileMgr.rightOpponent(curBanker)
	}

	// curBanker.hStatis.isContinuousBanker = false
	s.room.bankerChange2(newBanker)
}

// waitFinalDrawQuery 等待玩家选择是否开最后一张牌
func (s *SPlaying) waitFinalDrawQuery(curPlayer *PlayerHolder) bool {
	var qaIndex = s.room.nextQAIndex()

	actions := int(mahjong.ActionType_enumActionType_CustomB | mahjong.ActionType_enumActionType_SKIP)

	curPlayer.expectedAction = actions
	s.taskPlayerAction = newTaskPlayerAction(curPlayer, actions)
	s.taskPlayerAction.s = s
	s.taskPlayerAction.isForFinalDraw = true

	msgAllowPlayerAction2 := serializeMsgAllowedForDiscard2Opponent(curPlayer, qaIndex, curPlayer.expectedAction)
	for _, p := range s.players {
		p.sendActoinAllowedMsg(msgAllowPlayerAction2)
	}

	result := s.taskPlayerAction.wait()
	if !result {
		return false
	}

	action := s.taskPlayerAction.action
	s.taskPlayerAction = nil

	flags := mahjong.SRFlags_SRNone
	s.lctx.addActionWithTile(curPlayer, InvalidTile.tileID, 0, mahjong.ActionType(action), qaIndex, flags, actions)

	if action == int(mahjong.ActionType_enumActionType_SKIP) {
		s.onHandWashout()
		return false
	}

	msgActionResultNotify := serializeMsgActionResultNotifyForNoTile(int(mahjong.ActionType_enumActionType_CustomB),
		curPlayer)
	for _, p := range s.players {
		p.sendActionResultNotify(msgActionResultNotify)
	}

	return true
}

// waitPlayerAction 等待当前玩家出牌，或者加杠，暗杠，自摸胡牌
func (s *SPlaying) waitPlayerAction(curPlayer *PlayerHolder, newDraw bool) (loopKeep bool, discardedTile *Tile, action int) {
	// var newDraw = s.lctx.isSelfDraw(curPlayer)
	var qaIndex = s.room.nextQAIndex()
	loopKeep = true
	discardedTile = nil
	action = 0

	if newDraw {
		curPlayer.hStatis.resetLocked()
	}

	// 大丰麻将特殊规则：当自己打出一只牌后，别家再打出相同的牌，自己不能立即吃 / 碰，需要再次轮到自己出牌后，才可对此牌进行吃碰；
	// 自己吃 / 碰后不可以打出与用于吃 / 碰的两张牌组合成坎子 / 顺子的牌。
	actions := s.tileMgr.actionForDiscardPlayer(curPlayer, newDraw)

	msgAllowPlayerAction := serializeMsgAllowedForDiscard(s, curPlayer, actions, qaIndex)
	// 修正一下actions
	actions = int(msgAllowPlayerAction.GetAllowedActions())
	curPlayer.expectedAction = actions

	curPlayer.sendActoinAllowedMsg(msgAllowPlayerAction)
	if s.room.isForceConsistent() {
		s.sendMonkeyTips(curPlayer)
	}

	// 其他玩家只看到当前玩家在出牌，不能看到胡牌杠牌等动作
	msgAllowPlayerAction2 := serializeMsgAllowedForDiscard2Opponent(curPlayer, qaIndex,
		int(mahjong.ActionType_enumActionType_DISCARD))
	for _, p := range s.players {
		if p != curPlayer {
			p.sendActoinAllowedMsg(msgAllowPlayerAction2)
		}
	}

	s.taskPlayerAction = newTaskPlayerAction(curPlayer, actions)
	s.taskPlayerAction.s = s

	result := s.taskPlayerAction.wait()

	// var tile = s.taskPlayerAction.tile
	var tileID = s.taskPlayerAction.tileID
	var xflags = s.taskPlayerAction.flags
	action = s.taskPlayerAction.action
	s.taskPlayerAction = nil

	if !result {
		loopKeep = false
		return
	}

	if action == int(mahjong.ActionType_enumActionType_DISCARD) {

		// 出牌，从手牌队列移到出牌队列
		discardedTile = s.tileMgr.playerDiscard(curPlayer, tileID)
		// 记录出牌数据以及重置一些用于限制玩家动作的开关变量（每次玩家出牌后，这些开关变量都重置）
		// curPlayer.hStatis.resetLocked()

		// 记录最后打出的牌
		curPlayer.hStatis.latestDiscardedTileLocked = discardedTile

		flags := mahjong.SRFlags_SRNone

		// 处理庄家起手听
		if curPlayer == s.room.bankerPlayer() && curPlayer.hStatis.actionCounter == 1 {
			if xflags == 1 && curPlayer.tiles.readyHandAble() {
				curPlayer.hStatis.isRichi = true

				flags = mahjong.SRFlags_SRRichi
				var msgActionResultNotify = serializeMsgActionResultNotifyForNoTile(int(mahjong.ActionType_enumActionType_FirstReadyHand), curPlayer)
				for _, p := range s.players {
					p.sendActionResultNotify(msgActionResultNotify)
				}
			}
		}

		// 记录动作
		s.lctx.addActionWithTile(curPlayer, discardedTile.tileID, 0, mahjong.ActionType_enumActionType_DISCARD, qaIndex, flags, actions)

		orderPlayers := s.tileMgr.getOrderPlayers(curPlayer)
		taskDiscardReAction := analyseTaskDiscardReaction(s, orderPlayers, discardedTile, curPlayer)
		var needWaitReAction = false
		if taskDiscardReAction != nil {
			s.taskDiscardReAction = taskDiscardReAction
			needWaitReAction = true
		}

		// 发送出牌结果通知给所有客户端
		// 动作者自身需要收到听牌列表更新
		msgActionResultNotify := serializeMsgActionResultNotifyForDiscardedTile(int(mahjong.ActionType_enumActionType_DISCARD),
			curPlayer, discardedTile.tileID, needWaitReAction)
		for _, p := range s.players {
			p.sendActionResultNotify(msgActionResultNotify)
		}

	} else if action == int(mahjong.ActionType_enumActionType_KONG_Concealed) {
		// 以及实施杠牌：从_hands移除，并加入到面子牌表中
		s.tileMgr.kongConcealed(curPlayer, tileID)

		// 记录动作
		s.lctx.addActionWithTile(curPlayer, tileID, 0, mahjong.ActionType_enumActionType_KONG_Concealed,
			qaIndex, mahjong.SRFlags_SRNone, actions)

		// 发送操作结果通知给所有玩家
		msgActionResultNotify := serializeMsgActionResultNotifyForTile(int(mahjong.ActionType_enumActionType_KONG_Concealed),
			curPlayer, tileID)
		curPlayer.sendActionResultNotify(msgActionResultNotify)

		// 消除tile，仅发送标记到对手客户端
		// 需求修正：如果有吃椪，需要明牌暗杠
		var tilemax32 = msgActionResultNotify.GetActionTile()
		if curPlayer.tiles.chowPongExposedKongCount() == 0 {
			tilemax32 = int32(TILEMAX)
		}
		msgActionResultNotify.ActionTile = &tilemax32
		for _, p := range s.players {
			if p != curPlayer {
				p.sendActionResultNotify(msgActionResultNotify)
			}
		}

	} else if action == int(mahjong.ActionType_enumActionType_KONG_Triplet2) {
		// 实施加杠：从_hands移除，并加入到之前碰牌的面子牌组中
		discardedTile = s.tileMgr.triplet2Kong(curPlayer, tileID)

		// 记录动作
		s.lctx.addActionWithTile(curPlayer, tileID, 0, mahjong.ActionType_enumActionType_KONG_Triplet2,
			qaIndex, mahjong.SRFlags_SRNone, actions)

		orderPlayers := s.tileMgr.getOrderPlayers(curPlayer)
		taskDiscardReAction := analyseTaskTriplet2KongReaction(s, orderPlayers, discardedTile, curPlayer)
		if taskDiscardReAction != nil {
			s.taskDiscardReAction = taskDiscardReAction
		}

		// 发送加杠结果给所有人
		msgActionResultNotify := serializeMsgActionResultNotifyForTile(int(mahjong.ActionType_enumActionType_KONG_Triplet2),
			curPlayer, discardedTile.tileID)
		for _, p := range s.players {
			p.sendActionResultNotify(msgActionResultNotify)
		}
	} else if action == int(mahjong.ActionType_enumActionType_WIN_SelfDrawn) {
		// 自摸时这里不需要处理，gameLoop会处理
		// 记录动作
		s.lctx.addActionWithTile(curPlayer, InvalidTile.tileID, 0, mahjong.ActionType_enumActionType_WIN_SelfDrawn,
			qaIndex, mahjong.SRFlags_SRNone, actions)
	}
	return
}

// waitOpponentsRobKong 等待玩家抢杠胡
func (s *SPlaying) waitOpponentsRobKong(curPlayer *PlayerHolder, curTile *Tile) bool {
	result, _ := s.waitOpponentsAction(curPlayer, curTile)
	return result
}

// waitOpponentsAction 等待其他玩家对本次出牌的响应，例如有的玩家可以吃椪杠，则等其选择是否吃椪杠
func (s *SPlaying) waitOpponentsAction(curDiscardPlayer *PlayerHolder, discardedTile *Tile) (bool, *PlayerHolder) {
	orderPlayers := s.tileMgr.getOrderPlayers(curDiscardPlayer)

	if s.taskDiscardReAction == nil {
		return true, s.tileMgr.rightOpponent(curDiscardPlayer)
	}

	qaIndex := s.room.nextQAIndex()
	for _, p := range orderPlayers {
		if p.expectedAction == 0 {
			continue
		}

		msgAllowedAction := serializeMsgAllowedForDiscardResponse(p, qaIndex, discardedTile, curDiscardPlayer)
		p.sendReActoinAllowedMsg(msgAllowedAction)
		if s.room.isForceConsistent() {
			s.sendMonkeyTips(p)
		}
	}

	// 如果是用于复现问题的房间，则强制等待所有玩家都回复，否则由于客户端
	// 需要接收下一个动作提示，但是由于客户端发送上来的请求可以已经失效（更高优先级的玩家已经回复）
	// 导致该客户端的actionCounter不能增加1
	if s.room.isForceConsistent() {
		s.taskDiscardReAction.forceWaitAllReply = true
	}

	// 发送一个提示给出牌者
	// sendTips2Player(curDiscardPlayer, TipCode_TCWaitOpponentsAction)

	s.taskDiscardReAction.s = s
	result := s.taskDiscardReAction.wait()
	if !result {
		s.taskDiscardReAction = nil
		return false, nil
	}

	var taskDiscardReAction = s.taskDiscardReAction
	// 立即重置变量
	s.taskDiscardReAction = nil

	// 记录所有玩家尝试的action
	for e := taskDiscardReAction.waitQueue.Front(); e != nil; e = e.Next() {
		wi := e.Value.(*ReActionQueueItem)
		tileID := taskDiscardReAction.actionTile.tileID
		if wi.replyAction == int(mahjong.ActionType_enumActionType_CHOW) {
			tileID = int(wi.msgMeldTile.GetTile1())
		}

		s.lctx.addActionWithTile(wi.player, tileID, discardedTile.tileID, mahjong.ActionType(wi.replyAction), qaIndex,
			mahjong.SRFlags_SRUserReplyOnly, wi.actions)
	}

	replayAction := taskDiscardReAction.replayAction()
	// 全部选择过
	if replayAction == int(mahjong.ActionType_enumActionType_SKIP) {
		return true, s.tileMgr.rightOpponent(curDiscardPlayer)
	}

	// 为响应的玩家实施其选择
	s.applyActionForOpponents(taskDiscardReAction)

	// 如果有人胡牌，结束循坏
	if replayAction == int(mahjong.ActionType_enumActionType_WIN_Chuck) {
		s.onHandWinnerBornWinChuck(curDiscardPlayer)
		return false, nil
	}

	// 其他选择：吃、碰、杠，都需要该玩家继续出牌
	return true, taskDiscardReAction.who()
}

// applyActionForOpponents 对于玩家的响应动作，例如其选择碰，需要落实碰牌操作：修改牌列表等
// 对于玩家的响应动作不能放在handler代码中处理的原因是， 吃椪杠胡有优先级，假如一个玩家选择了碰，而
// 另一个玩家选择了胡，则胡者优先，因此我们需要cache玩家的选择，等待所有人都回复了（或者优先级最高者回复了）
// 再由本函数来落实玩家的选择。优先级处理见TaskDiscardReAction类代码
func (s *SPlaying) applyActionForOpponents(taskDiscardReAction *TaskPlayerReAction) {
	wi := taskDiscardReAction.whoWI()
	player := wi.player
	replayAction := taskDiscardReAction.replayAction()
	needNotifyClient := true
	var newMeld *Meld
	latestDiscardedTileID := taskDiscardReAction.actionTile.tileID

	switch replayAction {
	case int(mahjong.ActionType_enumActionType_CHOW):
		// 如果玩家处于过手胡锁定状态，但是发生了吃椪杠，则解除过手胡
		// 备注：2018-3-15 代理杨总和大冯哥都确认这次修改
		player.hStatis.isWinAbleLocked = false

		// 以及实施吃牌：从原拥有者移除，并加入到当前玩家的面子牌表中
		// 设置当前出牌者为player
		newMeld = s.tileMgr.chow(player, taskDiscardReAction)
		player.hStatis.latestChowPongTileLocked = taskDiscardReAction.actionTile

		// 记录动作
		s.lctx.addActionWithTile(player, newMeld.t1.tileID, taskDiscardReAction.actionTile.tileID,
			mahjong.ActionType_enumActionType_CHOW, s.room.qaIndex, mahjong.SRFlags_SRNone, wi.actions)
		break

	case int(mahjong.ActionType_enumActionType_PONG):
		// 如果玩家处于过手胡锁定状态，但是发生了吃椪杠，则解除过手胡
		// 备注：2018-3-15 代理杨总和大冯哥都确认这次修改
		player.hStatis.isWinAbleLocked = false

		player.hStatis.latestChowPongTileLocked = taskDiscardReAction.actionTile
		// 以及实施碰牌：从原拥有者移除，并加入到当前玩家的面子牌表中
		// 设置当前出牌者为player
		newMeld = s.tileMgr.pong(player, taskDiscardReAction)

		// 记录动作
		s.lctx.addActionWithTile(player, newMeld.t1.tileID, 0, mahjong.ActionType_enumActionType_PONG, s.room.qaIndex, mahjong.SRFlags_SRNone, wi.actions)
		break

	case int(mahjong.ActionType_enumActionType_KONG_Exposed):
		// 如果玩家处于过手胡锁定状态，但是发生了吃椪杠，则解除过手胡
		// 备注：2018-3-15 代理杨总和大冯哥都确认这次修改
		player.hStatis.isWinAbleLocked = false

		// 以及实施杠牌：从原拥有者移除，并加入到当前玩家的面子牌表中
		// 设置当前出牌者为player
		newMeld = s.tileMgr.kongExposed(player, taskDiscardReAction)
		// 记录动作
		s.lctx.addActionWithTile(player, newMeld.t1.tileID, 0, mahjong.ActionType_enumActionType_KONG_Exposed, s.room.qaIndex, mahjong.SRFlags_SRNone, wi.actions)
		break

	case int(mahjong.ActionType_enumActionType_WIN_Chuck):
		needNotifyClient = false
		// 可能有多个玩家胡牌，逐个检查
		// 由于最终会下发所有玩家的手牌列表，因此并不需要在此发送被吃的牌到所有的胡牌者客户端
		for e := taskDiscardReAction.waitQueue.Front(); e != nil; e = e.Next() {
			wi := e.Value.(*ReActionQueueItem)
			if wi.replyAction != int(mahjong.ActionType_enumActionType_WIN_Chuck) {
				continue
			}

			p := wi.player
			if s.tileMgr.winChuckAble(p, taskDiscardReAction) {
				// 实施吃铳胡牌
				s.tileMgr.winChuck(p, taskDiscardReAction)

				// 记录动作
				s.lctx.addActionWithTile(p, latestDiscardedTileID, 0,
					mahjong.ActionType_enumActionType_WIN_Chuck, s.room.qaIndex, mahjong.SRFlags_SRNone, wi.actions)
			}
		}
		break

	default:
		needNotifyClient = false
		break
	}

	if needNotifyClient {
		msgActionResultNotify := serializeMsgActionResultNotifyForResponse(replayAction,
			player, newMeld, latestDiscardedTileID)

		for _, p := range s.players {
			p.sendActionResultNotify(msgActionResultNotify)
		}
	}
}

// 发送monkey提示信息给客户端，以便客户端能够进行正确的操作选择
func (s *SPlaying) sendMonkeyTips(player *PlayerHolder) {
	monkeyUserTilesCfg := s.room.monkeyCfg.getMonkeyUserTilesCfg(player.userID())

	if monkeyUserTilesCfg == nil {
		return
	}

	actionTips := monkeyUserTilesCfg.actionTips
	if len(actionTips) <= player.hStatis.actionCounter {
		return
	}

	tips := actionTips[player.hStatis.actionCounter]

	player.sendTipsString(tips)
}

func sendTips2Player(player *PlayerHolder, tipCode mahjong.TipCode) {
	msgTips := &mahjong.MsgRoomShowTips{}
	var tipCode32 = int32(tipCode)
	msgTips.TipCode = &tipCode32

	player.sendMsg(msgTips, int32(mahjong.MessageCode_OPRoomShowTips))
}
