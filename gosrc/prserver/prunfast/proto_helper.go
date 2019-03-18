package prunfast

import (
	"pokerface"
)

// serializeMsgRestore 序列化掉线恢复消息给客户端
func serializeMsgRestore(s *SPlaying, player *PlayerHolder) *pokerface.MsgRestore {
	msgRestore := &pokerface.MsgRestore{}

	msgDeal := serializeMsgDeal(s, player)
	msgRestore.MsgDeal = msgDeal

	return msgRestore
}

// serializeMsgDeal 序列化发牌消息给客户端
func serializeMsgDeal(s *SPlaying, forwho *PlayerHolder) *pokerface.MsgDeal {
	var msg = &pokerface.MsgDeal{}
	var bankerChairID = int32(s.room.bankerPlayer().chairID)
	msg.BankerChairID = &bankerChairID
	var windFlowerID = int32(0)
	msg.WindFlowerID = &windFlowerID
	var cardInWall = int32(s.cardMgr.cardCountInWall())
	msg.CardsInWall = &cardInWall

	var markup = int32(s.room.markup)
	msg.Markup = &markup

	playerCardLists := make([]*pokerface.MsgPlayerCardList, len(s.players))

	for i, p := range s.players {
		var cardList *pokerface.MsgPlayerCardList
		if p == forwho {
			cardList = serializeCardListForSelf(p)
		} else {
			cardList = serializeCardListForOpponent(p)
		}
		playerCardLists[i] = cardList
	}

	msg.PlayerCardLists = playerCardLists
	return msg
}

// serializeCardListForSelf 序列化牌列表给自己
func serializeCardListForSelf(player *PlayerHolder) *pokerface.MsgPlayerCardList {
	playerCardList := &pokerface.MsgPlayerCardList{}
	cards := player.cards

	var chairID = int32(player.chairID)
	playerCardList.ChairID = &chairID

	// 已经打出去的牌
	playerCardList.DiscardedHands = cards.discardedCardHand2MsgCardHands()

	// 手牌
	var cardCountInHand = int32(cards.cardCountInHand())
	playerCardList.CardCountOnHand = &cardCountInHand
	playerCardList.CardsOnHand = cards.hand2IDList()

	return playerCardList
}

// serializeCardListForOpponent 序列化牌列表给其他玩家
func serializeCardListForOpponent(player *PlayerHolder) *pokerface.MsgPlayerCardList {
	playerCardList := &pokerface.MsgPlayerCardList{}
	cards := player.cards

	var chairID = int32(player.chairID)
	playerCardList.ChairID = &chairID

	// 已经打出去的牌
	playerCardList.DiscardedHands = cards.discardedCardHand2MsgCardHands()

	// 手牌
	var cardCountInHand = int32(cards.cardCountInHand())
	playerCardList.CardCountOnHand = &cardCountInHand

	return playerCardList
}

// // serializeMsgAllowedForRichi 序列化允许起手听消息，给那些可以起手听的玩家
// func serializeMsgAllowedForRichi(s *SPlaying, player *PlayerHolder, qaIndex int) *pokerface.MsgAllowPlayerAction {
// 	var msg = &pokerface.MsgAllowPlayerAction{}
// 	var qaIndex32 = int32(qaIndex)
// 	msg.QaIndex = &qaIndex32
// 	var allowedActions32 = int32(ActionType_enumActionType_FirstReadyHand | ActionType_enumActionType_SKIP)
// 	msg.AllowedActions = &allowedActions32
// 	var chairID32 = int32(player.chairID)
// 	msg.ActionChairID = &chairID32
// 	var timeout32 = int32(15)
// 	msg.TimeoutInSeconds = &timeout32

// 	// 听牌提示，也即是听什么牌，还剩下多少张之类的
// 	var msgReadyHandTip = &pokerface.MsgReadyHandTips{}
// 	var targetCard32 = int32(0)
// 	msgReadyHandTip.TargetCard = &targetCard32
// 	msgReadyHandTip.ReadyHandList = player.cards.readyHandCardsWhenThrow(TILEMAX, s.cardMgr)

// 	msg.TipsForAction = []*pokerface.MsgReadyHandTips{msgReadyHandTip}
// 	return msg
// }

// // serializeMsgActionResultNotifyForDraw 序列化抽牌结果给其他玩家
// func serializeMsgActionResultNotifyForDraw(player *PlayerHolder, cardID int, cardCountInWall int, newKongAfterCard *Card) *pokerface.MsgActionResultNotify {
// 	var msg = &pokerface.MsgActionResultNotify{}
// 	var action32 = int32(ActionType_enumActionType_DRAW)
// 	msg.Action = &action32
// 	var cardCountInWall32 = int32(cardCountInWall)
// 	msg.CardsInWall = &cardCountInWall32
// 	var chairID32 = int32(player.chairID)
// 	msg.TargetChairID = &chairID32
// 	var cardID32 = int32(cardID)
// 	msg.ActionCard = &cardID32

// 	if newKongAfterCard != nil {
// 		flowerIDs := []int32{int32(newKongAfterCard.cardID)}
// 		msg.NewFlowers = flowerIDs
// 	}

// 	return msg
// }

// // serializeMsgActionResultNotifyForCard 序列化某个玩家的动作结果给其他玩家
// func serializeMsgActionResultNotifyForCard(action int, player *PlayerHolder, cardID int) *pokerface.MsgActionResultNotify {
// 	var msg = &pokerface.MsgActionResultNotify{}
// 	var action32 = int32(action)
// 	msg.Action = &action32
// 	var chairID32 = int32(player.chairID)
// 	msg.TargetChairID = &chairID32
// 	var cardID32 = int32(cardID)
// 	msg.ActionCard = &cardID32

// 	return msg
// }

// func serializeMsgActionResultNotifyForAccWin(action int, player *PlayerHolder, cardID int, kongFollowLocked bool) *pokerface.MsgActionResultNotify {
// 	var msg = &pokerface.MsgActionResultNotify{}
// 	var action32 = int32(action)
// 	msg.Action = &action32
// 	var chairID32 = int32(player.chairID)
// 	msg.TargetChairID = &chairID32
// 	var cardID32 = int32(cardID)
// 	msg.ActionCard = &cardID32

// 	var kongFollowLockedInt32 int32
// 	if kongFollowLocked {
// 		kongFollowLockedInt32 = 1
// 	}
// 	msg.CardsInWall = &kongFollowLockedInt32

// 	return msg
// }

func serializeMsgActionResultNotifyForDiscardedCard(action ActionType, player *PlayerHolder, msgCardHand *pokerface.MsgCardHand) *pokerface.MsgActionResultNotify {
	var msg = &pokerface.MsgActionResultNotify{}
	var action32 = int32(action)
	msg.Action = &action32
	var chairID32 = int32(player.chairID)
	msg.TargetChairID = &chairID32

	msg.ActionHand = msgCardHand

	return msg
}

// serializeMsgActionResultNotifyForNoCard 序列化某个玩家的动作结果给其他玩家
func serializeMsgActionResultNotifyForNoCard(actoin ActionType, player *PlayerHolder) *pokerface.MsgActionResultNotify {
	var msg = &pokerface.MsgActionResultNotify{}
	var action32 = int32(actoin)
	msg.Action = &action32
	var chairID32 = int32(player.chairID)
	msg.TargetChairID = &chairID32

	return msg
}

// serializeMsgAllowedForDiscard2Opponent 序列化正在等待某个玩家出牌的消息给其他玩家
func serializeMsgAllowedForDiscard2Opponent(player *PlayerHolder, qaIndex int) *pokerface.MsgAllowPlayerAction {
	var msg = &pokerface.MsgAllowPlayerAction{}
	var qaIndex32 = int32(qaIndex)
	msg.QaIndex = &qaIndex32
	var allowedActions32 = int32(ActionType_enumActionType_DISCARD)
	msg.AllowedActions = &allowedActions32
	var chairID32 = int32(player.chairID)
	msg.ActionChairID = &chairID32
	var timeout32 = int32(15)
	msg.TimeoutInSeconds = &timeout32

	return msg
}

// serializeMsgAllowedForDiscard 序列化某个玩家出牌是允许的动作，例如不仅允许他出牌，还允许他暗杠，加杠，自摸胡牌等等
func serializeMsgAllowedForDiscard(s *SPlaying, player *PlayerHolder, actions int, qaIndex int) *pokerface.MsgAllowPlayerAction {
	var msg = &pokerface.MsgAllowPlayerAction{}
	var qaIndex32 = int32(qaIndex)
	msg.QaIndex = &qaIndex32
	var allowedActions32 = int32(actions)
	msg.AllowedActions = &allowedActions32
	var chairID32 = int32(player.chairID)
	msg.ActionChairID = &chairID32
	var timeout32 = int32(15)
	msg.TimeoutInSeconds = &timeout32

	// 如果第一出牌者，必须出带有红桃3的牌组(单张红桃3也可以)
	// 需求变更：只有三个人的时候，才强制出红桃3
	if s.room.config.playerNumAcquired == 3 {
		if player.hStatis.isFirstDiscarded && player.hStatis.actionCounter < 1 {
			timeout32 = 0x010f
		}
	}

	return msg
}

// // int32Distinct 辅助函数，用户从int32数组中找到不同的元素
// func int32Distinct(intList []int32) []int32 {
// 	dist := make(map[int32]bool)
// 	distInt32 := make([]int32, 0, len(intList))

// 	for _, v := range intList {
// 		if _, ok := dist[v]; !ok {
// 			dist[v] = true
// 			distInt32 = append(distInt32, v)
// 		}
// 	}

// 	return distInt32
// }

// // int32ListRemove 辅助函数，用于从int32数组中移除某个整数
// func int32ListRemove(intList []int32, element int) []int32 {
// 	for i, v := range intList {
// 		if v == int32(element) {
// 			result := append(intList[:i], intList[i+1:]...)
// 			return result
// 		}
// 	}

// 	return intList
// }

// serializeMsgAllowedForDiscardResponseRestore 用于离线恢复
func serializeMsgAllowedForDiscardResponseRestore(somebody *PlayerHolder, qaIndex int, discardedCardHand *CardHand, discardPlayer *PlayerHolder) *pokerface.MsgAllowPlayerReAction {
	return serializeMsgAllowedForDiscardResponse(somebody, qaIndex, discardedCardHand, discardPlayer)
}

// serializeMsgAllowedForDiscardResponse 序列化允许的反应动作给玩家，以便其他玩家，对某个玩家的出牌做出反应
func serializeMsgAllowedForDiscardResponse(player *PlayerHolder, qaIndex int, discardedCardHand *CardHand, discardPlayer *PlayerHolder) *pokerface.MsgAllowPlayerReAction {
	msg := &pokerface.MsgAllowPlayerReAction{}
	var qaIndex32 = int32(qaIndex)
	msg.QaIndex = &qaIndex32
	var allowedActions32 = int32(player.expectedAction)
	msg.AllowedActions = &allowedActions32
	var chairID32 = int32(player.chairID)
	msg.ActionChairID = &chairID32
	var timeout32 = int32(15)
	msg.TimeoutInSeconds = &timeout32
	var victimChairID32 = int32(discardPlayer.chairID)
	// msg.VictimChairID = &victimChairID32
	// var victimCardID32 = int32(discardedCard.cardID)
	// msg.VictimCardID = &victimCardID32
	discardAble := allowedActions32&int32(ActionType_enumActionType_DISCARD) != 0
	if discardAble {
		// 如果上一手是ACE，本手有2必须打2
		if discardedCardHand.ht == CardHandType_Single && discardedCardHand.cards[0].cardID/4 == (AH/4) {
			if player.cards.hasCardInHand(R2H) {
				// 不能过，必须打2
				allowedActions32 &= ^int32(ActionType_enumActionType_SKIP)
				timeout32 = (0x020f)
			}
		}
	}

	msg.PrevActionChairID = &victimChairID32
	msg.PrevActionHand = discardedCardHand.cardHand2MsgCardHand()

	return msg
}

// // serializeMsgActionResultNotifyForResponse 序列化玩家的动作结果给所有的玩家
// func serializeMsgActionResultNotifyForResponse(action int, player *PlayerHolder, newMeld *Meld, cardID int) *pokerface.MsgActionResultNotify {
// 	msg := &pokerface.MsgActionResultNotify{}
// 	var action32 = int32(action)
// 	msg.Action = &action32
// 	var chairID32 = int32(player.chairID)
// 	msg.TargetChairID = &chairID32
// 	var actionCard32 = int32(cardID)
// 	msg.ActionCard = &actionCard32

// 	if newMeld != nil {
// 		msg.ActionMeld = player.cards.meld2MsgMeldCard(newMeld, false)

// 		// 有了吃椪杠，暗杠需要明牌显示
// 		concealedMeldIDs := player.cards.concealedKongIDList()
// 		if len(concealedMeldIDs) > 0 {
// 			msg.NewFlowers = concealedMeldIDs
// 		}
// 	}

// 	return msg
// }

// serializeMsgHandOver 序列化单手牌结束的消息（非流局）给客户端
func serializeMsgHandOver(s *SPlaying, handOverType int) *pokerface.MsgHandOver {
	// 计算分数
	var msgHandScore = &pokerface.MsgHandScore{}

	playerScores := make([]*pokerface.MsgPlayerScore, 0, len(s.players))

	for _, player := range s.players {
		var sc = player.sctx
		var msgPlayerScore = &pokerface.MsgPlayerScore{}
		var targetChairID32 = int32(player.chairID)
		msgPlayerScore.TargetChairID = &targetChairID32
		var winType32 = int32(sc.winType)
		msgPlayerScore.WinType = &winType32
		var score32 = int32(sc.calcTotalWinScore()) // 总得分
		msgPlayerScore.Score = &score32
		var specialScore32 = int32(0)
		msgPlayerScore.SpecialScore = &specialScore32 // 新疆麻将本标志用于表示压大胡次数

		var fakeWinScoreFlag32 = int32(0)
		msgPlayerScore.FakeWinScore = &fakeWinScoreFlag32 // 新疆麻将本标志用于表示压小胡次数

		var continuousBankerMultiple = int32(sc.markupMultiple)
		msgPlayerScore.ContinuousBankerMultiple = &continuousBankerMultiple // 新疆麻将本标志用于爬坡落庄倍数

		var isContinuosBanker = sc.isPayForAll
		msgPlayerScore.IsContinuousBanker = &isContinuosBanker // 新疆麻将本标志用于表示是否包庄

		if handOverType != int(HandOverType_enumHandOverType_None) {
			greatWin := &pokerface.MsgPlayerScoreGreatWin{}
			var baseWinScore32 = int32(1)
			greatWin.BaseWinScore = &baseWinScore32
			var greatWinPoints32 = int32(sc.greatWinPoints)
			greatWin.GreatWinPoints = &greatWinPoints32
			var greatWinType32 = int32(sc.greatWinType)
			greatWin.GreatWinType = &greatWinType32
			var trimGreatWinPoints32 = int32(sc.greatWinPoints)
			greatWin.TrimGreatWinPoints = &trimGreatWinPoints32

			var endXScore = int32(0)
			greatWin.ContinuousBankerExtra = &endXScore // 马牌分数
			msgPlayerScore.GreatWin = greatWin
		}
		//}

		playerScores = append(playerScores, msgPlayerScore)
	}

	msgHandScore.PlayerScores = playerScores

	// 构造MsgHandOver
	msgHandOver := &pokerface.MsgHandOver{}
	msgHandOver.Scores = msgHandScore

	// 赢牌类型
	var endType32 = int32(handOverType)
	msgHandOver.EndType = &endType32
	var isContinueAble = s.room.handRoundStarted < s.room.config.handNum
	msgHandOver.ContinueAble = &isContinueAble

	msgHandOver.PlayerCardLists = serializeCardListsForHandOver(s)

	return msgHandOver
}

// serializeCardListsForHandOver 序列化单手牌结束消息给客户端
func serializeCardListsForHandOver(s *SPlaying) []*pokerface.MsgPlayerCardList {
	playerCardLists := make([]*pokerface.MsgPlayerCardList, len(s.players))

	for i, p := range s.players {
		var cardList *pokerface.MsgPlayerCardList
		cardList = serializeCardListForHandOver(p)
		playerCardLists[i] = cardList
	}

	return playerCardLists
}

// serializeCardListForHandOver 序列化单手牌结束时，玩家的手牌内容给客户端
func serializeCardListForHandOver(player *PlayerHolder) *pokerface.MsgPlayerCardList {
	// 所有用户的手牌列表，如果有暗杠，还需要把暗杠带上，以便客户端显示
	return serializeCardListForSelf(player)
}

// serializeMsgHandOverWashout 序列化单手牌结束（流局）消息给客户端
func serializeMsgHandOverWashout(s *SPlaying) *pokerface.MsgHandOver {
	return serializeMsgHandOver(s, int(HandOverType_enumHandOverType_None))
}

// serializeMsgGameOver 序列化游戏结束消息给客户端
func serializeMsgGameOver(r *Room) *pokerface.MsgGameOver {
	var playerStats = make([]*pokerface.MsgGameOverPlayerStat, len(r.players))

	for i, p := range r.players {
		var stat = &pokerface.MsgGameOverPlayerStat{}
		var score32 = int32(p.gStatis.roundScore)
		stat.Score = &score32

		var winChuckCounter32 = int32(p.gStatis.miniWinCounter)
		stat.WinChuckCounter = &winChuckCounter32
		var winSelfDrawn32 = int32(p.gStatis.winSelfDrawnCounter)
		stat.WinSelfDrawnCounter = &winSelfDrawn32

		var chucker32 = int32(p.gStatis.greatWinCounter)
		stat.ChuckerCounter = &chucker32
		var chairID32 = int32(p.chairID)
		stat.ChairID = &chairID32

		var robKongCounter32 = int32(p.gStatis.winRobKongCounter)
		var kongerCounter32 = int32(p.gStatis.kongerCounter)
		stat.RobKongCounter = &robKongCounter32
		stat.KongerCounter = &kongerCounter32

		playerStats[i] = stat
	}

	var msgGameOver = &pokerface.MsgGameOver{}
	msgGameOver.PlayerStats = playerStats

	return msgGameOver
}

// serializeMsgRoomInfo 序列化房间信息给客户端
func serializeMsgRoomInfo(r *Room) *pokerface.MsgRoomInfo {
	msg := &pokerface.MsgRoomInfo{}

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

	playerInfos := make([]*pokerface.MsgPlayerInfo, len(r.players))
	for i, p := range r.players {
		var msgPlayerInfo = &pokerface.MsgPlayerInfo{}
		var chairID32 = int32(p.chairID)
		msgPlayerInfo.ChairID = &chairID32
		var userID = p.userID()
		msgPlayerInfo.UserID = &userID
		var pstate32 = int32(p.state)
		msgPlayerInfo.State = &pstate32

		var userInfo = p.user.getInfo()
		msgPlayerInfo.Nick = &userInfo.nick
		msgPlayerInfo.Sex = &userInfo.sex
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
