package zjmahjong

import (
	"mahjong"

	log "github.com/sirupsen/logrus"
)

// serializeMsgRestore 序列化掉线恢复消息给客户端
func serializeMsgRestore(s *SPlaying, player *PlayerHolder) *mahjong.MsgRestore {
	msgRestore := &mahjong.MsgRestore{}

	msgDeal := serializeMsgDeal(s, player)
	msgRestore.MsgDeal = msgDeal

	richiChairs := make([]int32, 0, len(s.players))

	for _, p := range s.players {
		if p.hStatis.isRichi {
			richiChairs = append(richiChairs, int32(p.chairID))
		}
	}

	msgRestore.ReadyHandChairs = richiChairs

	var lastDiscaredChairID32 = int32(5)
	var srAction = s.lctx.getLastNonDrawAction()
	if srAction != nil {
		if srAction.GetAction() == int32(mahjong.ActionType_enumActionType_DISCARD) {
			lastDiscaredChairID32 = (srAction.GetChairID())
		}
	}

	msgRestore.LastDiscaredChairID = &lastDiscaredChairID32

	var isMeNewDraw = false
	srAction = s.lctx.current()
	if srAction != nil && srAction.GetAction() == int32(mahjong.ActionType_enumActionType_DRAW) &&
		srAction.GetChairID() == int32(player.chairID) {
		isMeNewDraw = true
	}

	msgRestore.IsMeNewDraw = &isMeNewDraw

	var waitDiscardReAction = false
	if srAction != nil && srAction.GetAction() == int32(mahjong.ActionType_enumActionType_DISCARD) {
		waitDiscardReAction = true
	}

	log.Println("msgRestore, WaitDiscardReAction:", waitDiscardReAction)
	msgRestore.WaitDiscardReAction = &waitDiscardReAction
	return msgRestore
}

// serializeMsgDeal 序列化发牌消息给客户端
func serializeMsgDeal(s *SPlaying, forwho *PlayerHolder) *mahjong.MsgDeal {
	var msg = &mahjong.MsgDeal{}
	var bankerChairID = int32(s.room.bankerPlayer().chairID)
	msg.BankerChairID = &bankerChairID
	var windFlowerID = int32(0)
	msg.WindFlowerID = &windFlowerID
	var tileInWall = int32(s.tileMgr.tileCountInWall())
	msg.TilesInWall = &tileInWall
	var isContinuousBanker = s.room.bankerPlayer().gStatis.isContinuousBanker
	msg.IsContinuousBanker = &isContinuousBanker

	// 家家庄
	var markup32 = int32(0)
	msg.Markup = &markup32

	// 骰子
	msg.Dice1 = &s.room.dice1
	msg.Dice2 = &s.room.dice2

	playerTileLists := make([]*mahjong.MsgPlayerTileList, len(s.players))

	for i, p := range s.players {
		var tileList *mahjong.MsgPlayerTileList
		if p == forwho {
			tileList = serializeTileListForSelf(p)
		} else {
			tileList = serializeTileListForOpponent(p)
		}
		playerTileLists[i] = tileList
	}

	msg.PlayerTileLists = playerTileLists
	return msg
}

func dice12(room *Room) (int32, int32) {
	dice1 := room.rand.Intn(6) + 1
	dice2 := room.rand.Intn(6) + 1

	return int32(dice1), int32(dice2)
}

// serializeTileListForSelf 序列化牌列表给自己
func serializeTileListForSelf(player *PlayerHolder) *mahjong.MsgPlayerTileList {
	playerTileList := &mahjong.MsgPlayerTileList{}
	tiles := player.tiles

	var chairID = int32(player.chairID)
	playerTileList.ChairID = &chairID

	// 已经打出去的牌
	playerTileList.TilesDiscard = tiles.discard2IDList()

	// 已经落地的面子牌
	playerTileList.Melds = tiles.melds2MsgMeldTileList(false)

	// 花牌
	playerTileList.TilesFlower = tiles.flower2IDList()

	// 手牌
	var tileCountInHand = int32(tiles.tileCountInHand())
	playerTileList.TileCountInHand = &tileCountInHand
	playerTileList.TilesHand = tiles.hand2IDList()

	return playerTileList
}

// serializeTileListForOpponent 序列化牌列表给其他玩家
func serializeTileListForOpponent(player *PlayerHolder) *mahjong.MsgPlayerTileList {
	playerTileList := &mahjong.MsgPlayerTileList{}
	tiles := player.tiles

	var chairID = int32(player.chairID)
	playerTileList.ChairID = &chairID

	// 已经打出去的牌
	playerTileList.TilesDiscard = tiles.discard2IDList()

	// 已经落地的面子牌，暗杠只发标记
	mark := true
	playerTileList.Melds = tiles.melds2MsgMeldTileList(mark)

	// 花牌
	playerTileList.TilesFlower = tiles.flower2IDList()

	// 手牌，只发一个数量
	var tileCountInHand = int32(tiles.tileCountInHand())
	playerTileList.TileCountInHand = &tileCountInHand

	return playerTileList
}

// serializeMsgAllowedForRichi 序列化允许起手听消息，给那些可以起手听的玩家
func serializeMsgAllowedForRichi(s *SPlaying, player *PlayerHolder, qaIndex int) *mahjong.MsgAllowPlayerAction {
	var msg = &mahjong.MsgAllowPlayerAction{}
	var qaIndex32 = int32(qaIndex)
	msg.QaIndex = &qaIndex32
	var allowedActions32 = int32(mahjong.ActionType_enumActionType_FirstReadyHand | mahjong.ActionType_enumActionType_SKIP)
	msg.AllowedActions = &allowedActions32
	var chairID32 = int32(player.chairID)
	msg.ActionChairID = &chairID32
	var timeout32 = int32(15)
	msg.TimeoutInSeconds = &timeout32

	// 听牌提示，也即是听什么牌，还剩下多少张之类的
	var msgReadyHandTip = &mahjong.MsgReadyHandTips{}
	var targetTile32 = int32(0)
	msgReadyHandTip.TargetTile = &targetTile32
	msgReadyHandTip.ReadyHandList = player.tiles.readyHandTilesWhenThrow(TILEMAX, s.tileMgr)

	msg.TipsForAction = []*mahjong.MsgReadyHandTips{msgReadyHandTip}
	return msg
}

// serializeMsgActionResultNotifyForDraw 序列化抽牌结果给其他玩家
func serializeMsgActionResultNotifyForDraw(player *PlayerHolder, tileID int, flowers []*Tile, tileCountInWall int) *mahjong.MsgActionResultNotify {
	var msg = &mahjong.MsgActionResultNotify{}
	var action32 = int32(mahjong.ActionType_enumActionType_DRAW)
	msg.Action = &action32
	var tileCountInWall32 = int32(tileCountInWall)
	msg.TilesInWall = &tileCountInWall32
	var chairID32 = int32(player.chairID)
	msg.TargetChairID = &chairID32
	var tileID32 = int32(tileID)
	msg.ActionTile = &tileID32

	if len(flowers) > 0 {
		flowerIDs := make([]int32, len(flowers))
		for i, f := range flowers {
			flowerIDs[i] = int32(f.tileID)
		}

		msg.NewFlowers = flowerIDs
	}

	return msg
}

// serializeMsgActionResultNotifyForTile 序列化某个玩家的动作结果给其他玩家
func serializeMsgActionResultNotifyForTile(action int, player *PlayerHolder, tileID int) *mahjong.MsgActionResultNotify {
	var msg = &mahjong.MsgActionResultNotify{}
	var action32 = int32(action)
	msg.Action = &action32
	var chairID32 = int32(player.chairID)
	msg.TargetChairID = &chairID32
	var tileID32 = int32(tileID)
	msg.ActionTile = &tileID32

	return msg
}

func serializeMsgActionResultNotifyForDiscardedTile(action int, player *PlayerHolder, tileID int, needWaitReAction bool) *mahjong.MsgActionResultNotify {
	var msg = &mahjong.MsgActionResultNotify{}
	var action32 = int32(action)
	msg.Action = &action32
	var chairID32 = int32(player.chairID)
	msg.TargetChairID = &chairID32
	var tileID32 = int32(tileID)
	msg.ActionTile = &tileID32
	msg.WaitDiscardReAction = &needWaitReAction
	return msg
}

// serializeMsgActionResultNotifyForNoTile 序列化某个玩家的动作结果给其他玩家
func serializeMsgActionResultNotifyForNoTile(actoin int, player *PlayerHolder) *mahjong.MsgActionResultNotify {
	var msg = &mahjong.MsgActionResultNotify{}
	var action32 = int32(actoin)
	msg.Action = &action32
	var chairID32 = int32(player.chairID)
	msg.TargetChairID = &chairID32

	return msg
}

// serializeMsgAllowedForDiscard2Opponent 序列化正在等待某个玩家出牌的消息给其他玩家
func serializeMsgAllowedForDiscard2Opponent(player *PlayerHolder, qaIndex int, actions int) *mahjong.MsgAllowPlayerAction {
	var msg = &mahjong.MsgAllowPlayerAction{}
	var qaIndex32 = int32(qaIndex)
	msg.QaIndex = &qaIndex32
	var allowedActions32 = int32(actions)
	msg.AllowedActions = &allowedActions32
	var chairID32 = int32(player.chairID)
	msg.ActionChairID = &chairID32
	var timeout32 = int32(15)
	msg.TimeoutInSeconds = &timeout32

	return msg
}

// serializeMsgAllowedForDiscard 序列化某个玩家出牌是允许的动作，例如不仅允许他出牌，还允许他暗杠，加杠，自摸胡牌等等
func serializeMsgAllowedForDiscard(s *SPlaying, player *PlayerHolder, actions int, qaIndex int) *mahjong.MsgAllowPlayerAction {
	var msg = &mahjong.MsgAllowPlayerAction{}
	var qaIndex32 = int32(qaIndex)
	msg.QaIndex = &qaIndex32
	var allowedActions32 = int32(actions)
	msg.AllowedActions = &allowedActions32
	var chairID32 = int32(player.chairID)
	msg.ActionChairID = &chairID32
	var timeout32 = int32(15)
	msg.TimeoutInSeconds = &timeout32

	// 仅可以胡牌时，直接返回胡牌，没有其他选项
	if actions == int(mahjong.ActionType_enumActionType_WIN_SelfDrawn) {
		return msg
	}

	discardAble := (actions & int(mahjong.ActionType_enumActionType_DISCARD)) != 0
	// 修正一下expetedAction：当可以出牌时，不允许客户端发送“skip”到服务器
	if discardAble {
		xActions := player.expectedAction
		player.expectedAction = (xActions & (^int(mahjong.ActionType_enumActionType_SKIP)))
	}

	msgMelds := make([]*mahjong.MsgMeldTile, 0, 4)
	// 如果可以暗杠，则加上可以暗杠的牌型
	if actions&int(mahjong.ActionType_enumActionType_KONG_Concealed) != 0 {
		concealedIDList := player.tiles.concealedKongAble2IDList()

		for _, id := range concealedIDList {
			msgMeld := &mahjong.MsgMeldTile{}
			var meldType32 = int32(mahjong.MeldType_enumMeldTypeConcealedKong)
			msgMeld.MeldType = &meldType32
			var tile132 = int32(id)
			msgMeld.Tile1 = &tile132
			msgMeld.Contributor = &chairID32

			msgMelds = append(msgMelds, msgMeld)
		}
	}

	// 可以加杠，加上加杠牌型
	if actions&int(mahjong.ActionType_enumActionType_KONG_Triplet2) != 0 {
		// msgMeld := &MsgMeldTile{}
		// var meldType32 = int32(MeldType_enumMeldTypeTriplet2Kong)
		// msgMeld.MeldType = &meldType32
		// var tile132 = int32(player.tiles.latestHandTile().tileID)
		// msgMeld.Tile1 = &tile132
		// msgMeld.Contributor = &chairID32

		// msgMelds = append(msgMelds, msgMeld)

		triplet2KongIDList := player.tiles.triplet2KongAble2IDList()
		for _, id := range triplet2KongIDList {
			msgMeld := &mahjong.MsgMeldTile{}
			var meldType32 = int32(mahjong.MeldType_enumMeldTypeTriplet2Kong)
			msgMeld.MeldType = &meldType32
			var tile132 = int32(id)
			msgMeld.Tile1 = &tile132
			msgMeld.Contributor = &chairID32

			msgMelds = append(msgMelds, msgMeld)
		}
	}

	if len(msgMelds) > 0 {
		msg.MeldsForAction = msgMelds
	}

	if discardAble {
		// 如果处于听牌状态
		// 则只允许打出刚摸到的牌
		selfDraw := s.lctx.isSelfDraw(player)
		if player.hStatis.isRichi || (selfDraw && s.tileMgr.tileCountInWall() == 0) {
			tid := player.tiles.latestHandTile().tileID
			msgReadyHandTip := &mahjong.MsgReadyHandTips{}
			tid32 := int32(tid)
			msgReadyHandTip.TargetTile = &tid32
			msgReadyHandTip.ReadyHandList = player.tiles.readyHandTilesWhenThrow(tid, s.tileMgr)
			msg.TipsForAction = []*mahjong.MsgReadyHandTips{msgReadyHandTip}

			if allowedActions32 == int32(mahjong.ActionType_enumActionType_SKIP|mahjong.ActionType_enumActionType_DISCARD) {
				allowedActions32 &= ^int32(mahjong.ActionType_enumActionType_SKIP)
			}
			return msg
		}

		// 构建可以打出的牌列表，每一个可以打出的牌，如果有牌可听也一并发给客户端
		tidsHand := int32Distinct(player.tiles.hand2IDList())
		if player.hStatis.latestChowPongTileLocked != InvalidTile && player.tiles.tileCountInHand() > 2 {
			// 移除不能出的牌列表
			deleteIDs := player.tiles.chowPongLockedIDList(player.hStatis.latestChowPongTileLocked)
			for _, dID := range deleteIDs {
				tidsHand = int32ListRemove(tidsHand, dID)
			}
		}

		tipsArray := make([]*mahjong.MsgReadyHandTips, len(tidsHand))
		n := 0
		for i, v := range tidsHand {
			msgReadyHandTips := &mahjong.MsgReadyHandTips{}
			v32 := int32(v)
			msgReadyHandTips.TargetTile = &v32

			readyHandList := player.tiles.readyHandTilesWhenThrow(int(v), s.tileMgr)
			if len(readyHandList) > 0 {
				n += len(readyHandList)
				msgReadyHandTips.ReadyHandList = readyHandList
			}
			tipsArray[i] = msgReadyHandTips
		}

		if n == 0 && (actions&int(mahjong.ActionType_enumActionType_FirstReadyHand) != 0) {
			// 没牌可听，去掉听牌允许
			actions &= ^int(mahjong.ActionType_enumActionType_FirstReadyHand)
			player.expectedAction = actions
			allowedActions32 = int32(actions)
		}

		msg.TipsForAction = tipsArray
	}

	if allowedActions32 == int32(mahjong.ActionType_enumActionType_SKIP|mahjong.ActionType_enumActionType_DISCARD) {
		allowedActions32 &= ^int32(mahjong.ActionType_enumActionType_SKIP)
	}

	return msg
}

// int32Distinct 辅助函数，用户从int32数组中找到不同的元素
func int32Distinct(intList []int32) []int32 {
	dist := make(map[int32]bool)
	distInt32 := make([]int32, 0, len(intList))

	for _, v := range intList {
		if _, ok := dist[v]; !ok {
			dist[v] = true
			distInt32 = append(distInt32, v)
		}
	}

	return distInt32
}

// int32ListRemove 辅助函数，用于从int32数组中移除某个整数
func int32ListRemove(intList []int32, element int) []int32 {
	for i, v := range intList {
		if v == int32(element) {
			result := append(intList[:i], intList[i+1:]...)
			return result
		}
	}

	return intList
}

// serializeMsgAllowedForDiscardResponse 序列化允许的反应动作给玩家，以便其他玩家，对某个玩家的出牌做出反应
func serializeMsgAllowedForDiscardResponse(player *PlayerHolder, qaIndex int,
	discardedTile *Tile, discardPlayer *PlayerHolder) *mahjong.MsgAllowPlayerReAction {
	msg := &mahjong.MsgAllowPlayerReAction{}
	var qaIndex32 = int32(qaIndex)
	msg.QaIndex = &qaIndex32
	var allowedActions32 = int32(player.expectedAction)
	msg.AllowedActions = &allowedActions32
	var chairID32 = int32(player.chairID)
	msg.ActionChairID = &chairID32
	var timeout32 = int32(15)
	msg.TimeoutInSeconds = &timeout32
	var victimChairID32 = int32(discardPlayer.chairID)
	msg.VictimChairID = &victimChairID32
	var victimTileID32 = int32(discardedTile.tileID)
	msg.VictimTileID = &victimTileID32

	actions := player.expectedAction
	msgMelds := make([]*mahjong.MsgMeldTile, 0, 10)

	// 如果可以明杠，加上杠牌的牌型
	if actions&int(mahjong.ActionType_enumActionType_KONG_Exposed) != 0 {
		msgMeld := &mahjong.MsgMeldTile{}
		var meldType32 = int32(mahjong.MeldType_enumMeldTypeExposedKong)
		msgMeld.MeldType = &meldType32
		var tile132 = int32(discardedTile.tileID)
		msgMeld.Tile1 = &tile132
		var contributor32 = int32(discardPlayer.chairID)
		msgMeld.Contributor = &contributor32

		msgMelds = append(msgMelds, msgMeld)
	}

	//  如果可以碰，加上碰牌的牌型
	if actions&int(mahjong.ActionType_enumActionType_PONG) != 0 {
		msgMeld := &mahjong.MsgMeldTile{}
		var meldType32 = int32(mahjong.MeldType_enumMeldTypeTriplet)
		msgMeld.MeldType = &meldType32
		var tile132 = int32(discardedTile.tileID)
		msgMeld.Tile1 = &tile132
		var contributor32 = int32(discardPlayer.chairID)
		msgMeld.Contributor = &contributor32

		msgMelds = append(msgMelds, msgMeld)
	}

	// 如果可以吃，加上吃牌的牌型
	if actions&int(mahjong.ActionType_enumActionType_CHOW) != 0 {
		idList := player.tiles.chowAble2IdList(discardedTile.tileID)
		for _, tid := range idList {
			msgMeld := &mahjong.MsgMeldTile{}
			var meldType32 = int32(mahjong.MeldType_enumMeldTypeSequence)
			msgMeld.MeldType = &meldType32
			var tile132 = int32(tid)
			msgMeld.Tile1 = &tile132
			var contributor32 = int32(discardPlayer.chairID)
			msgMeld.Contributor = &contributor32

			msgMelds = append(msgMelds, msgMeld)
		}
	}

	if len(msgMelds) > 0 {
		msg.MeldsForAction = msgMelds
	}

	return msg
}

// serializeMsgActionResultNotifyForResponse 序列化玩家的动作结果给所有的玩家
func serializeMsgActionResultNotifyForResponse(action int, player *PlayerHolder, newMeld *Meld, tileID int) *mahjong.MsgActionResultNotify {
	msg := &mahjong.MsgActionResultNotify{}
	var action32 = int32(action)
	msg.Action = &action32
	var chairID32 = int32(player.chairID)
	msg.TargetChairID = &chairID32
	var actionTile32 = int32(tileID)
	msg.ActionTile = &actionTile32

	if newMeld != nil {
		msg.ActionMeld = player.tiles.meld2MsgMeldTile(newMeld, false)

		// 有了吃椪杠，暗杠需要明牌显示
		concealedMeldIDs := player.tiles.concealedKongIDList()
		if len(concealedMeldIDs) > 0 {
			msg.NewFlowers = concealedMeldIDs
		}
	}

	return msg
}

// serializeMsgHandOver 序列化单手牌结束的消息（非流局）给客户端
func serializeMsgHandOver(s *SPlaying, handOverType int) *mahjong.MsgHandOver {
	// 计算分数
	var msgHandScore = &mahjong.MsgHandScore{}

	playerScores := make([]*mahjong.MsgPlayerScore, 0, len(s.players))
	banker := s.room.bankerPlayer()

	for _, player := range s.players {
		var sc = player.sctx
		var msgPlayerScore = &mahjong.MsgPlayerScore{}
		var targetChairID32 = int32(player.chairID)
		msgPlayerScore.TargetChairID = &targetChairID32
		var winType32 = int32(sc.winType)
		msgPlayerScore.WinType = &winType32
		var score32 = int32(sc.calcTotalWinScore())
		msgPlayerScore.Score = &score32
		var specialScore32 = int32(sc.horseCount) // 湛江麻用于表示中马个数
		msgPlayerScore.SpecialScore = &specialScore32
		var fakeWinScore32 = int32(0)
		msgPlayerScore.FakeWinScore = &fakeWinScore32

		msgPlayerScore.IsContinuousBanker = &sc.isContinuousBanker
		var continuousBankerMultiple32 = int32(0) // 湛江麻将不用
		msgPlayerScore.ContinuousBankerMultiple = &continuousBankerMultiple32

		if sc.greatWinType != 0 {
			greatWin := &mahjong.MsgPlayerScoreGreatWin{}
			var baseWinScore32 = int32(0)
			greatWin.BaseWinScore = &baseWinScore32
			var greatWinPoints32 = int32(sc.greatWinPoints)
			greatWin.GreatWinPoints = &greatWinPoints32
			var greatWinType32 = int32(sc.greatWinType)
			greatWin.GreatWinType = &greatWinType32
			var trimGreatWinPoints32 = int32(0) // 湛江麻将不用
			greatWin.TrimGreatWinPoints = &trimGreatWinPoints32

			msgPlayerScore.GreatWin = greatWin
		}

		if player == banker {
			if s.horseTiles != nil {
				horseList := make([]int32, len(s.horseTiles))
				for i, t := range s.horseTiles {
					horseList[i] = int32(t.tileID)
				}

				// 湛江麻将把马牌放到庄家的fake list上
				msgPlayerScore.FakeList = horseList
			}
		}
		playerScores = append(playerScores, msgPlayerScore)
	}

	msgHandScore.PlayerScores = playerScores

	// 构造MsgHandOver
	msgHandOver := &mahjong.MsgHandOver{}
	msgHandOver.Scores = msgHandScore

	// 赢牌类型
	var endType32 = int32(handOverType)
	msgHandOver.EndType = &endType32
	var isContinueAble = s.room.isContinuAble()
	msgHandOver.ContinueAble = &isContinueAble

	msgHandOver.PlayerTileLists = serializeTileListsForHandOver(s)

	return msgHandOver
}

// serializeTileListsForHandOver 序列化单手牌结束消息给客户端
func serializeTileListsForHandOver(s *SPlaying) []*mahjong.MsgPlayerTileList {
	playerTileLists := make([]*mahjong.MsgPlayerTileList, len(s.players))

	for i, p := range s.players {
		var tileList *mahjong.MsgPlayerTileList
		tileList = serializeTileListForHandOver(p)
		playerTileLists[i] = tileList
	}

	return playerTileLists
}

// serializeTileListForHandOver 序列化单手牌结束时，玩家的手牌内容给客户端
func serializeTileListForHandOver(player *PlayerHolder) *mahjong.MsgPlayerTileList {
	// 所有用户的手牌列表，如果有暗杠，还需要把暗杠带上，以便客户端显示
	playerTileList := &mahjong.MsgPlayerTileList{}
	tiles := player.tiles

	var chairID = int32(player.chairID)
	playerTileList.ChairID = &chairID

	// 已经落地的面子牌，false表示，所有面子牌都公开
	playerTileList.Melds = tiles.melds2MsgMeldTileList(false)

	// 手牌
	var tileCountInHand = int32(tiles.tileCountInHand())
	playerTileList.TileCountInHand = &tileCountInHand
	playerTileList.TilesHand = tiles.hand2IDList()

	return playerTileList
}

// serializeMsgHandOverWashout 序列化单手牌结束（流局）消息给客户端
func serializeMsgHandOverWashout(s *SPlaying) *mahjong.MsgHandOver {
	msgHandOver := &mahjong.MsgHandOver{}

	var endType32 = int32(mahjong.HandOverType_enumHandOverType_None)
	msgHandOver.EndType = &endType32
	var isContinueAble = s.room.handRoundStarted < s.room.config.handNum
	msgHandOver.ContinueAble = &isContinueAble

	msgHandOver.PlayerTileLists = serializeTileListsForHandOver(s)

	return msgHandOver
}

// serializeMsgGameOver 序列化游戏结束消息给客户端
func serializeMsgGameOver(r *Room) *mahjong.MsgGameOver {
	var playerStats = make([]*mahjong.MsgGameOverPlayerStat, len(r.players))

	for i, p := range r.players {
		var stat = &mahjong.MsgGameOverPlayerStat{}
		var score32 = int32(p.gStatis.roundScore)
		stat.Score = &score32
		var winChuckCounter32 = int32(p.gStatis.winChuckCounter)
		stat.WinChuckCounter = &winChuckCounter32
		var winSelfDrawn32 = int32(p.gStatis.winSelfDrawnCounter)
		stat.WinSelfDrawnCounter = &winSelfDrawn32
		var chucker32 = int32(p.gStatis.chuckerCounter)
		stat.ChuckerCounter = &chucker32
		var chairID32 = int32(p.chairID)
		stat.ChairID = &chairID32

		playerStats[i] = stat
	}

	var msgGameOver = &mahjong.MsgGameOver{}
	msgGameOver.PlayerStats = playerStats

	return msgGameOver
}

// serializeMsgRoomInfo 序列化房间信息给客户端
func serializeMsgRoomInfo(r *Room) *mahjong.MsgRoomInfo {
	msg := &mahjong.MsgRoomInfo{}

	var state32 = int32(r.state.getStateConst())
	msg.State = &state32
	var ownerID = r.ownerID
	msg.OwnerID = &ownerID
	var roomNumber = r.roomNumber
	msg.RoomNumber = &roomNumber
	var handStartted32 = int32(r.handRoundStarted)
	msg.HandStartted = &handStartted32
	var handFinished32 = int32(r.handRoundFinished)
	msg.HandFinished = &handFinished32

	playerInfos := make([]*mahjong.MsgPlayerInfo, len(r.players))
	for i, p := range r.players {
		var msgPlayerInfo = &mahjong.MsgPlayerInfo{}
		var chairID32 = int32(p.chairID)
		msgPlayerInfo.ChairID = &chairID32
		var userID = p.userID()
		msgPlayerInfo.UserID = &userID
		var pstate32 = int32(p.state)
		msgPlayerInfo.State = &pstate32

		var userInfo = p.user.getInfo()
		msgPlayerInfo.Nick = &userInfo.nick
		msgPlayerInfo.Gender = &userInfo.gender
		msgPlayerInfo.HeadIconURI = &userInfo.headIconURI
		msgPlayerInfo.Ip = &userInfo.ip
		msgPlayerInfo.Location = &userInfo.location
		var dfHands = int32(userInfo.dfHands)
		msgPlayerInfo.DfHands = &dfHands

		var diamond32 = int32(userInfo.diamond)
		msgPlayerInfo.Diamond = &diamond32
		var charm32 = int32(userInfo.charm)
		msgPlayerInfo.Charm = &charm32
		var avatarID = int32(userInfo.avatarID)
		msgPlayerInfo.AvatarID = &avatarID
		var clubIDs = userInfo.clubIDs
		msgPlayerInfo.ClubIDs = clubIDs
		var dan = int32(userInfo.dan)
		msgPlayerInfo.Dan = &dan

		playerInfos[i] = msgPlayerInfo
	}

	msg.Players = playerInfos
	msg.ScoreRecords = r.scoreRecords

	return msg
}
