package pddz

// payCardScore2Winner 支付牌分
func payCardScore2Winner(loser *PlayerHolder, winner *PlayerHolder, scoreMultiple int, s *SPlaying) {
	s.cl.Printf("pay2Winner, loser %d, winner %d, score multiple:%d, room markup:%d, baseScore:%d\n",
		loser.chairID, winner.chairID, s.room.markup,
		scoreMultiple, s.room.config.baseScore)

	limitMultiple := scoreMultiple
	if !s.room.config.isCallWithScore {
		// 不是兰州麻将
		limitMultiple = scoreMultiple * s.room.markup
	}

	multipleLimit := s.room.config.multipleLimit
	if multipleLimit != 0 {
		s.cl.Println("room with multipleLimit:", multipleLimit)
		if limitMultiple > multipleLimit {
			s.cl.Printf("scoreMultiple trim to multipleLimit, from %d => %d", limitMultiple, multipleLimit)
			limitMultiple = multipleLimit
		}
	}

	if s.room.config.isCallWithScore {
		// 兰州麻将这里需要乘以markup
		limitMultiple = limitMultiple * s.room.markup
	}

	score2Pay := limitMultiple * s.room.config.baseScore

	loser.sctx.getPayTarget(winner).cardWinScore -= score2Pay
	winner.sctx.getPayTarget(loser).cardWinScore += score2Pay

	s.cl.Printf("payer:%s[%d] pay card-win-score %d[limitMultiple:%d * baseScore:%d] to winner%s[%d]\n",
		loser.userID(), loser.chairID, score2Pay, limitMultiple, s.room.config.baseScore, winner.userID(), winner.chairID)
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

	// 分为两个派别计算输赢和支付
	landlordPlayer := s.room.landlordPlayer()
	farmers := s.cardMgr.getOrderPlayers(landlordPlayer)

	totalDiscardedBombCount := 0
	for _, player := range s.players {
		totalDiscardedBombCount = totalDiscardedBombCount + player.cards.bombOrRocketCountOfDiscarded()
	}

	// 基本分是房间倍数乘以底分
	s.cl.Printf("markup %d X baseScore %d\n", s.room.markup, s.room.config.baseScore)

	scoreMultiple := 1
	// 一个炸弹翻倍
	for i := 0; i < totalDiscardedBombCount; i++ {
		scoreMultiple = scoreMultiple * 2
	}

	s.cl.Printf("consider with bomb count:%d, score has multiple to:%d\n", totalDiscardedBombCount, scoreMultiple)

	if winner == landlordPlayer {
		// 地主赢牌
		// 春天：地主所有牌出完，其他两家一张都未出，地主分数翻倍
		discardedCount := 0
		for _, p := range farmers {
			discardedCount = discardedCount + s.lctx.discardedActionCount(p)
		}

		if discardedCount == 0 {
			landlordPlayer.sctx.spring = true
			scoreMultiple = scoreMultiple * 2
			s.cl.Printf("landlordPlayer's spring, score has multiple X2 to:%d\n", scoreMultiple)
		}

		if landlordPlayer.hStatis.isCallDouble {
			scoreMultiple = scoreMultiple * 2
			s.cl.Printf("landlordPlayer[winner] choose call-double, win score has multiple X2 to:%d\n", scoreMultiple)
		}

		for _, p := range farmers {
			payScoreMultiple := scoreMultiple
			if p.hStatis.isCallDouble {
				payScoreMultiple = payScoreMultiple * 2
				s.cl.Printf("farmer[loser] %d choose call-double, lose score has multiple X2 to:%d\n",
					p.chairID, payScoreMultiple)
			}

			payCardScore2Winner(p, landlordPlayer, payScoreMultiple, s)
		}
	} else {
		// 农民赢牌
		// 反春天：农民中有一家先出完牌，地主只出过一手牌，农民分数翻倍
		discardedCount := s.lctx.discardedActionCount(landlordPlayer)
		if discardedCount == 1 {
			for _, p := range farmers {
				p.sctx.spring = true
			}

			scoreMultiple = scoreMultiple * 2
			s.cl.Printf("farmer's spring, score has multiple X2 to:%d\n", scoreMultiple)
		}

		if landlordPlayer.hStatis.isCallDouble {
			scoreMultiple = scoreMultiple * 2
			s.cl.Printf("landlordPlayer[loser] %d choose call-double, lose score has multiple X2 to:%d\n",
				landlordPlayer.chairID, scoreMultiple)
		}

		for _, p := range farmers {
			payScoreMultiple := scoreMultiple
			if p.hStatis.isCallDouble {
				payScoreMultiple = payScoreMultiple * 2
				s.cl.Printf("farmer[winner] %d choose call-double, win score has multiple X2 to:%d\n", p.chairID,
					payScoreMultiple)
			}

			payCardScore2Winner(landlordPlayer, p, payScoreMultiple, s)
		}
	}

	// 最终计分
	orderPlayers := make([]*PlayerHolder, 0, len(s.players))
	orderPlayers = append(orderPlayers, winner)
	orderPlayers = append(orderPlayers, s.cardMgr.getOrderPlayers(winner)...)
	doFinalPay(s, orderPlayers)
}

// calcFinalResultWashout 流局计算
func calcFinalResultWashout(s *SPlaying) {
	s.cl.Println("calcFinalResultWashout")

	for _, p := range s.players {
		p.sctx = &ScoreContext{}
		p.sctx.initPlayerScoreContext(s.cardMgr.getOrderPlayers(p), s.room)
	}
}
