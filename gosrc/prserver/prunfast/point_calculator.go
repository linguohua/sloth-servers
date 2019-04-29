package prunfast

// payCardScore2Winner 支付牌分
func payCardScore2Winner(chucker *PlayerHolder, loser *PlayerHolder, winner *PlayerHolder, s *SPlaying) {
	s.cl.Printf("pay2Winner, loser %d, winner %d\n", loser.chairID, winner.chairID)

	// curBanker := s.room.bankerPlayer()
	var score2Pay = 0
	var remainCardCount = loser.cards.cardCountInHand()

	var markupMultiple int
	if remainCardCount < 10 {
		markupMultiple = 1
	} else if remainCardCount < 16 {
		markupMultiple = 2
	} else {
		markupMultiple = 3
	}

	s.cl.Printf("loser %d has %d cards remain, multiple:%d\n", loser.chairID, remainCardCount, markupMultiple)

	score2Pay = remainCardCount * markupMultiple

	var payer *PlayerHolder
	if chucker == nil {
		payer = loser
	} else {
		payer = chucker
		s.cl.Printf("pay2Winner, chucker %d, pay to winner %d for loser %d\n", chucker.chairID, winner.chairID, loser.chairID)
	}

	payer.sctx.getPayTarget(winner).cardWinScore -= score2Pay
	winner.sctx.getPayTarget(payer).cardWinScore += score2Pay

	s.cl.Printf("payer:%d pay card-win-score %d 2 winner %d\n", loser.chairID, score2Pay, winner.chairID)
}

// collectMyEarn 直接地收取某个玩家所赢的钱，如果输家不够，输家就进入保护状态，而不为输家去收取其他人的钱
func collectMyEarn(s *SPlaying, player *PlayerHolder) {
	for _, pc := range player.sctx.orderPlayerSctxs {
		if !pc.hasFinalPay && pc.finalWinScore > 0 {
			loser := pc.target
			winner := player

			finalPay(s, winner, loser, pc)
		}
	}
}

// finalPay 带保护的支付
func finalPay(s *SPlaying, winner *PlayerHolder, loser *PlayerHolder, pc *PlayerScoreContext) {
	shouldPay := pc.finalWinScore
	shouldPayTrim := shouldPay
	s.cl.Printf("player:%d final pay 2 payer:%d, trim:%d=>%d\n", loser.chairID, winner.chairID, shouldPay, shouldPayTrim)

	pc.hasFinalPay = true
	pc.finalWinScore = shouldPayTrim
	winner.gStatis.roundScore += shouldPayTrim

	loserPC := loser.sctx.getPayTarget(winner)
	loserPC.finalWinScore = -shouldPayTrim
	loserPC.hasFinalPay = true
	loser.gStatis.roundScore -= shouldPayTrim
}

// doFinalPay 最终计分，orderPlayers赢家按照逆时针排在前端
func doFinalPay(s *SPlaying, orderPlayers []*PlayerHolder) {
	s.cl.Println("doFinalPay")
	// roomConfig := s.room.config

	// 汇总各种得分
	for _, p := range orderPlayers {
		for _, pc := range p.sctx.orderPlayerSctxs {
			if pc.hasCalc {
				continue
			}

			if pc.cardWinScore != 0 {
				pc.finalWinScore += pc.cardWinScore
			}

			pc.hasCalc = true
			pc2 := pc.target.sctx.getPayTarget(p)
			pc2.finalWinScore = -pc.finalWinScore
			pc2.hasCalc = true
		}
	}

	// 检查玩家得分者
	for _, p := range orderPlayers {
		winner := p
		for _, pc := range p.sctx.orderPlayerSctxs {
			// 能够从某人身上赢钱
			if !pc.hasFinalPay && pc.finalWinScore > 0 {
				loser := pc.target

				// 输家先把所有他该得到的钱收回来，以便付给赢家
				// 注意只有链条上的直接输家才能收取其所赢的钱，链条上的下一个输家是没有机会收取其所赢的钱的
				// 例如，A收取B的钱，B可以收取其他人的钱，但假如B要收取C的钱，此时C就不能像B一样收取其他人的钱
				collectMyEarn(s, loser)

				finalPay(s, winner, loser, pc)
			}
		}
	}
}

// calcFinalResultSelfDraw 计算自摸胡牌时的得分结果
func calcFinalResultSelfDraw(s *SPlaying, winner *PlayerHolder) {
	s.cl.Printf("calcFinalResultSelfDraw, winner chairID:%d, userID:%s\n", winner.chairID, winner.userID())
	for _, p := range s.players {
		p.sctx = &ScoreContext{}
		p.sctx.initPlayerScoreContext(s.cardMgr.getOrderPlayers(p), s.room)
	}

	// 自摸胡牌只有一个赢牌者
	winner.sctx.winType = int(HandOverType_enumHandOverType_Win_SelfDrawn)

	var chucker *PlayerHolder
	chucker = calcFake(s, winner)

	if chucker != nil {
		chucker.sctx.winType = int(HandOverType_enumHandOverType_Chucker)
	}

	var losers = s.cardMgr.getOrderPlayers(winner)
	// 每一个输牌者付分
	for _, loser := range losers {
		// 此处考虑爬坡、落庄之类
		payCardScore2Winner(chucker, loser, winner, s)
	}

	// 最终计分
	orderPlayers := make([]*PlayerHolder, 0, len(s.players))
	orderPlayers = append(orderPlayers, winner)
	orderPlayers = append(orderPlayers, s.cardMgr.getOrderPlayers(winner)...)
	doFinalPay(s, orderPlayers)
}

func calcFake(s *SPlaying, winner *PlayerHolder) *PlayerHolder {
	s.cl.Println("calcFake, winner:", winner.chairID)
	leftOpponent := s.cardMgr.leftOpponent(winner)
	e := winner.cards.discarded.Back()
	var msgCardHand = e.Value.(*CardHand)

	if len(msgCardHand.cards) >= 4 {
		s.cl.Printf("winner:%d not in alarm state, lastHandCardLength:%d, not fake\n", winner.chairID, len(msgCardHand.cards))
		return nil
	}

	prevAction := s.lctx.prevprev()
	if prevAction == nil {
		s.cl.Println("calcFake, prevAction == nil, not fake")
		return nil
	}

	if prevAction.GetChairID() != int32(leftOpponent.chairID) {
		s.cl.Printf("calcFake, chair %d not match %d, not fake\n", prevAction.GetChairID(), leftOpponent.chairID)
		return nil
	}

	// 表示PayerReAction
	if leftOpponent.hStatis.lastExpectedType == 1 {
		s.cl.Printf("calcFake, leftOpponent %d lastExpectedType 1, player-reaction\n", leftOpponent.chairID)
		// 炸弹那就没有办法了，不能算包牌，这里考虑到报警了，那么只有一种炸弹：就是3个ace
		if msgCardHand.ht == (CardHandType_Bomb) {
			s.cl.Printf("calcFake, winner %d last hand is bomb(3ace), not fake\n", winner.chairID)
			return nil
		}

		// 如果上家只能选择过，则不算包牌
		if leftOpponent.hStatis.lastExpectedActions == int(ActionType_enumActionType_SKIP) {
			s.cl.Printf("calcFake, leftOpponent %d only skipable, not fake\n", leftOpponent.chairID)
			return nil
		}

		// 上家可以打牌，却选择过，包牌
		if prevAction.GetAction() == int32(ActionType_enumActionType_SKIP) {
			s.cl.Printf("calcFake, leftOpponent %d can discard but choose skip, fake\n", leftOpponent.chairID)
			return leftOpponent
		}

		// 上家选择打牌，检查他的牌是不是可以的最大的牌型
		e2 := leftOpponent.cards.discarded.Back()
		var msgCardHand2 = e2.Value.(*CardHand)
		if leftOpponent.cards.hasCardHandGreatThan(msgCardHand2) {
			s.cl.Printf("calcFake, leftOpponent %d discarded not the greatest one, fake\n", leftOpponent.chairID)
			return leftOpponent
		}
	} else {
		s.cl.Printf("calcFake, leftOpponent %d lastExpectedType 0, player-action\n", leftOpponent.chairID)
		// 表示PlayerAction
		// 炸弹那就没有办法了，不能算包牌，这里考虑到报警了，那么只有一种炸弹：就是3个ace
		if msgCardHand.ht == (CardHandType_Bomb) {
			s.cl.Printf("calcFake, winner %d last hand is bomb(3ace), not fake\n", winner.chairID)
			return nil
		}

		if msgCardHand.ht == CardHandType_Single {
			// 检查是否必然打单张，如果不是则认为是包牌
			e2 := leftOpponent.cards.discarded.Back()
			var msgCardHand2 = e2.Value.(*CardHand)
			leftOpponentDiscarded := msgCardHand2.cards[0].cardID
			// 检查是否只能打出一个单张
			if !leftOpponent.cards.allSingleCardWith(leftOpponentDiscarded) {
				s.cl.Printf("calcFake, leftOpponent %d FREE discard single hand  %d, but it can discard pair/triplet/flush or others actually, fake\n",
					leftOpponent.chairID, leftOpponentDiscarded)
				return leftOpponent
			}

			msgCardHand22 := msgCardHand2.cardHand2MsgCardHand()
			if leftOpponent.cards.hasSingleGreatThan(msgCardHand22) {
				// 如果最后打出的牌，不是最大的单张，则包牌
				s.cl.Printf("calcFake, leftOpponent %d FREE discard non-greatest single hand %d, fake\n",
					leftOpponent.chairID, leftOpponentDiscarded)
				return leftOpponent
			}
		} else {
			// 对子，三张的情形下必是包牌
			return leftOpponent
		}
	}

	return nil
}

// calcFinalResultWashout 流局计算
func calcFinalResultWashout(s *SPlaying) {
	s.cl.Println("calcFinalResultWashout")

	for _, p := range s.players {
		p.sctx = &ScoreContext{}
		p.sctx.initPlayerScoreContext(s.cardMgr.getOrderPlayers(p), s.room)
	}
}
