package zjmahjong

import "mahjong"

// calcGreatWinTileType 计算跟行牌无关的牌型大胡
// 1.清一色
// 2.混一色
// 3.碰碰胡
// 4.七对
func calcGreatWinTileType(s *SPlaying, player *PlayerHolder) (int, float32) {
	var tiles = player.tiles
	var points float32
	var winType = 0
	if tiles.suitTypeCount() == 1 {
		if tiles.honorTypeCount() > 0 {

			//}
		} else {
			if tiles.exposedMeldCount() == 0 {
				if tiles.flowerTileCount() == 0 {
					// 清一色：一色牌组成的胡牌。
					winType |= int(GreatWinType_PureSame)
					var gp = float32(1.0)
					points += gp
					s.cl.Println("GWT:PureSame:", gp)
				}
			}
		}
	}

	// 七对，豪华七对
	var st = tiles.calc7Pair()
	switch st {
	case GreatWinType_GreatSevenPair:
		winType |= int(GreatWinType_GreatSevenPair)
		var gp = float32(1.0)
		points += gp
		s.cl.Println("GWT:GreatSevenPair:", gp)
		break
	case GreatWinType_SevenPair:
		winType |= int(GreatWinType_SevenPair)
		var gp = float32(1.0)
		points += gp
		s.cl.Println("GWT:SevenPair:", gp)
		break
	}

	return winType, points
}

// calcGreatWinning 判断玩家的胡牌是否大胡（大丰：辣子胡）
// 规则：
// 独钓：吃碰杠一起12只，剩余1只，胡剩余的1只。
// 海底捞月：摸牌池最后一张牌，并产生胡牌。
// 碰碰胡：由4副刻子或杠，和1对相同的牌组成的胡牌。
// 混一色：万条筒其中一种与风牌结合一起，且产生的胡牌。
// 清一色：一色牌组成的胡牌。
// 大门清：胡牌时，无吃碰杠，且没有抓过花。
// 七对：7对不一样的牌组成的胡牌。
// 豪华大七对：有4个同种牌，且胡的那只刚好是4只相同中的1只.
// 天胡：庄家起手摸第14只牌，产生胡牌.
// 暗杠胡：手牌里有3只一样的牌，同时胡第4只1样的牌。
// 明杠胡：直杠名牌后，摸牌产生胡牌。
// 起手报听胡牌：起手报听，报听后胡牌。
func calcGreatWinning(s *SPlaying, player *PlayerHolder, selfDrawn bool) {
	var points float32
	var winType = 0
	var tiles = player.tiles
	sc := player.sctx

	if !tiles.winAble() {
		s.cl.Panic("calcGreatWinning, not winable")
		return
	}

	// 计算牌型性质的大胡
	winType, points = calcGreatWinTileType(s, player)

	// var selfDrawn = s.lctx.isSelfDraw(player)
	if !selfDrawn && s.lctx.isRobKong() {
		// 如果是最后动作是加杠，则表明是抢杠胡
		winType |= int(GreatWinType_RobKong)
		var gp = float32(1.0)
		points += gp
		s.cl.Println("GWT:RobKong:", gp)
	}

	// 天胡
	if player == s.room.bankerPlayer() && player.hStatis.actionCounter == 1 && selfDrawn {
		winType |= int(GreatWinType_Heaven)
		var gp = float32(1.0)
		points += gp
		s.cl.Println("GWT:Heaven:", gp)
	}

	// 暗杠胡：手牌里有3只一样的牌，同时胡第4只1样的牌。（必须自摸）
	// 注意不是岭上开花
	if selfDrawn && tiles.tileCountInHandOf(tiles.latestHandTile().tileID) == 4 {
		winType |= int(GreatWinType_AfterConcealedKong)
		var gp = float32(1.0)
		points += gp
		s.cl.Println("GWT:AfterConcealedKong:", gp)
	}

	// 明杠胡：碰牌后，依然胡碰的那只牌。（必须自摸）
	// 注意不是岭上开花
	if selfDrawn && tiles.hasPongOf(tiles.latestHandTile().tileID) {
		winType |= int(GreatWinType_AfterExposedKong)
		var gp = float32(1.0)
		points += gp
		s.cl.Println("GWT:AfterExposedKong:", gp)
	}

	sc.greatWinType = winType
	sc.fGreatWinPoints = points

	s.cl.Printf("great win point:%f, type:%d\n", points, winType)
}

func calcPay2Winner(loser *PlayerHolder, winner *PlayerHolder, room *Room) int {
	room.cl.Printf("calcPay2Winner, loser %d, winner %d\n", loser.chairID, winner.chairID)

	score2Pay := 0

	return score2Pay
}

func pay2Winner(loser *PlayerHolder, winner *PlayerHolder, room *Room) {

	score2Pay := calcPay2Winner(loser, winner, room)

	loser.sctx.getPayTarget(winner).totalWinScore -= score2Pay
	winner.sctx.getPayTarget(loser).totalWinScore += score2Pay

	room.cl.Printf("loser:%d pay score %d 2 winner %d\n", loser.chairID, score2Pay, winner.chairID)
}

// calcFinalResultSelfDraw 计算自摸胡牌时的得分结果
func calcFinalResultSelfDraw(s *SPlaying, winner *PlayerHolder) {
	s.cl.Printf("calcFinalResultSelfDraw, winner chairID:%d, userID:%s\n", winner.chairID, winner.userID())
	for _, p := range s.players {
		p.sctx = &ScoreContext{}
		p.sctx.initPlayerScoreContext(s.tileMgr.getOrderPlayers(p))
	}

	// 自摸胡牌只有一个赢牌者
	winner.sctx.winType = int(mahjong.HandOverType_enumHandOverType_Win_SelfDrawn)

	// 计算大小胡
	basicScoreCalc(s, winner, true)

	// 其他人各自付分
	for _, p := range s.players {
		if p == winner {
			continue
		}

		pay2Winner(p, winner, s.room)
	}

	orderPlayers := make([]*PlayerHolder, 0, len(s.players))
	orderPlayers = append(orderPlayers, winner)
	orderPlayers = append(orderPlayers, s.tileMgr.getOrderPlayers(winner)...)
	doFinalPay(s, orderPlayers)
}

// collectMyEarn 直接地收取某个玩家所赢的钱，如果输家不够，输家就进入保护状态，而不为输家去收取其他人的钱
func collectMyEarn(s *SPlaying, player *PlayerHolder) {
	for _, pc := range player.sctx.orderPlayerSctxs {
		if !pc.hasClear && pc.totalWinScore > 0 {
			loser := pc.target
			winner := player

			loseProtectPay(s, winner, loser, pc)
		}
	}
}

// loseProtectPay 带保护的支付
func loseProtectPay(s *SPlaying, winner *PlayerHolder, loser *PlayerHolder, pc *PlayerScoreContext) {
	shouldPay := pc.totalWinScore
	shouldPayTrim := shouldPay
	s.cl.Printf("player:%d final pay 2 payer:%d, trim:%d=>%d\n", loser.chairID, winner.chairID, shouldPay, shouldPayTrim)

	pc.hasClear = true
	pc.totalWinScore = shouldPayTrim
	winner.gStatis.roundScore += shouldPayTrim

	loserPC := loser.sctx.getPayTarget(winner)
	loserPC.totalWinScore = -shouldPayTrim
	loserPC.hasClear = true
	loser.gStatis.roundScore -= shouldPayTrim
}

// doFinalPay 最终计分，orderPlayers赢家按照逆时针排在前端
func doFinalPay(s *SPlaying, orderPlayers []*PlayerHolder) {
	s.cl.Println("doFinalPay")
	// 汇总包牌得失分
	for _, p := range orderPlayers {
		for _, pc := range p.sctx.orderPlayerSctxs {
			if pc.fakeWinScore != 0 {
				pc.totalWinScore += pc.fakeWinScore
			}
		}
	}

	// 检查玩家得分者
	for _, p := range orderPlayers {
		winner := p
		for _, pc := range p.sctx.orderPlayerSctxs {
			// 能够从某人身上赢钱
			if !pc.hasClear && pc.totalWinScore > 0 {
				loser := pc.target

				// 输家先把所有他该得到的钱收回来，以便付给赢家
				// 注意只有链条上的直接输家才能收取其所赢的钱，链条上的下一个输家是没有机会收取其所赢的钱的
				// 例如，A收取B的钱，B可以收取其他人的钱，但假如B要收取C的钱，此时C就不能像B一样收取其他人的钱
				collectMyEarn(s, loser)

				loseProtectPay(s, winner, loser, pc)
			}
		}
	}
}

// basicScoreCalc 计算单个玩家基础的得分
func basicScoreCalc(s *SPlaying, player *PlayerHolder, selfDraw bool) {
	calcGreatWinning(s, player, selfDraw)
}

// calcFinalResultWithChucker 计算吃铳胡牌时各个玩家的得分
func calcFinalResultWithChucker(s *SPlaying, chucker *PlayerHolder) {
	s.cl.Printf("calcFinalResultWithChucker, chucker chairID:%d, userID:%s\n", chucker.chairID, chucker.userID())
	for _, p := range s.players {
		p.sctx = &ScoreContext{}
		p.sctx.initPlayerScoreContext(s.tileMgr.getOrderPlayers(p))
	}

	// 放铳者
	chucker.sctx.winType = int(mahjong.HandOverType_enumHandOverType_Chucker)

	// 先计算所有人大小胡牌情况
	for _, p := range s.players {
		if !p.tiles.winAble() {
			continue
		}

		var winner = p
		basicScoreCalc(s, winner, false)
		winner.sctx.winType = int(mahjong.HandOverType_enumHandOverType_Win_Chuck)
	}

	// 计算分数
	// 为每一个可以胡牌者，放铳者都需要付出对方的分数
	for _, p := range s.players {
		if !p.tiles.winAble() || p == chucker {
			continue
		}

		var winner = p
		pay2Winner(chucker, winner, s.room)
	}

	orderPlayers := make([]*PlayerHolder, 0, len(s.players))
	xorderPlayers := s.tileMgr.getOrderPlayers(chucker)
	for _, xp := range xorderPlayers {
		if xp.sctx.winType == int(mahjong.HandOverType_enumHandOverType_Win_Chuck) {
			orderPlayers = append(orderPlayers, xp)
		}
	}

	for _, xp := range xorderPlayers {
		if xp.sctx.winType != int(mahjong.HandOverType_enumHandOverType_Win_Chuck) {
			orderPlayers = append(orderPlayers, xp)
		}
	}

	orderPlayers = append(orderPlayers, chucker)

	doFinalPay(s, orderPlayers)
}
