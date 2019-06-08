package zjmahjong

import (
	"mahjong"

	log "github.com/sirupsen/logrus"
)

// calcGreatWinTileType 计算跟行牌无关的牌型大胡
// 1.清一色
// 2.碰碰胡
// 3.七对，大七对
// 4.十三幺
// 5.全风子
func calcGreatWinTileType(s *SPlaying, player *PlayerHolder) (int, int) {
	var tiles = player.tiles
	var points int
	var winType = 0
	roomConfig := s.room.config

	if tiles.suitTypeCount() == 1 {
		if tiles.honorTypeCount() == 0 {
			// 清一色：一色牌组成的胡牌
			winType |= int(GreatWinType_PureSame)
			var gp = roomConfig.pureSamePoint()
			points += gp
			s.cl.Println("GWT:PureSame:", gp)
		}
	}

	// 七对，豪华七对
	var st = tiles.calc7Pair()
	switch st {
	case GreatWinType_GreatSevenPair:
		winType |= int(GreatWinType_GreatSevenPair)
		var gp = roomConfig.greatSevenPairPoint()
		points += gp
		s.cl.Println("GWT:GreatSevenPair:", gp)
		break
	case GreatWinType_SevenPair:
		winType |= int(GreatWinType_SevenPair)
		var gp = roomConfig.sevenPairPoint()
		points += gp
		s.cl.Println("GWT:SevenPair:", gp)
		break
	}

	// 碰碰胡：全部是碰牌牌组的胡牌
	if tiles.isAllTriplet() {
		winType |= int(GreatWinType_PongPong)
		var gp = roomConfig.pongPongPoint()
		points += gp
		log.Println("GWT:PongKong, normal:", gp)
	}

	// 十三幺
	if tiles.isThirteenOrphans() {
		winType |= int(GreatWinType_Thirteen)
		var gp = roomConfig.thirteenOrphansPoint()
		points += gp
		log.Println("GWT:, ThirteenOrphans:", gp)
	}

	// 全风子 : 全部由风牌组成的胡牌
	if tiles.suitTypeCount() == 0 {
		winType |= int(GreatWinType_AllWind)
		var gp = roomConfig.allWindPoint()
		points += gp
		log.Println("GWT:, All Wind:", gp)
	}

	return winType, points
}

// calcGreatWinning 判断玩家的胡牌是否大胡（大丰：辣子胡）
// 先计算牌型倍数
// 再计算行牌倍数
func calcGreatWinning(s *SPlaying, player *PlayerHolder, selfDrawn bool) {
	var points int
	var winType = 0
	var tiles = player.tiles
	sc := player.sctx
	roomConfig := s.room.config

	if !tiles.winAble() {
		s.cl.Panic("calcGreatWinning, not winable")
		return
	}

	// 计算牌型性质的大胡
	winType, points = calcGreatWinTileType(s, player)

	// 抢杠胡，抢续杠后胡牌
	if !selfDrawn && s.lctx.isRobKong() {
		// 如果是最后动作是加杠，则表明是抢杠胡
		winType |= int(GreatWinType_RobKong)
		s.cl.Println("GWT:RobKong")
	}

	// 天胡，庄家起手胡牌
	if player == s.room.bankerPlayer() && player.hStatis.actionCounter == 1 && selfDrawn {
		winType |= int(GreatWinType_Heaven)
		var gp = roomConfig.heavenPoint()
		points += gp
		s.cl.Println("GWT:Heaven:", gp)
	}

	// 自杠胡：暗杠/续杠后，自摸胡牌
	if selfDrawn && s.lctx.isSelfKong(player) {
		winType |= int(GreatWinType_AfterConcealedKong | GreatWinType_AfterKong)
		var gp = roomConfig.afterKongPint()
		points += gp
		s.cl.Println("GWT:AfterConcealedKong:", gp)
	}

	// 放杠胡：明杠，对手出牌放杠，自摸胡牌，对方全包
	if selfDrawn && s.lctx.kongerOf(player, s.room) != nil {
		winType |= int(GreatWinType_AfterExposedKong | GreatWinType_AfterKong)
		var gp = roomConfig.afterKongPint()
		points += gp
		s.cl.Println("GWT:AfterExposedKong:", gp)
	}

	// 海底捞：自摸牌墙最后一张牌而胡牌
	if s.tileMgr.wallEmpty() && selfDrawn {
		winType |= int(GreatWinType_FinalDraw)
		var gp = roomConfig.finalDrawPoint()
		points += gp
		log.Println("GWT:FinalDraw:", gp)
	}

	sc.greatWinType = winType
	sc.greatWinPoints = points

	s.cl.Printf("great win point:%d, type:%d\n", points, winType)
}

// 总分=N X单个输赢
// 单个输赢 = 底分X (基础倍数 ）X 中马倍数
// 其中：
// N表示需要付分玩家数量（有人全包时，该人付N份）；
// 如果当前是2人，N最多等于1；当前3人，N最多等于2，当前4人，N最多等于3；
// 基础倍数=牌型倍数之和；
// 中马倍数=中马个数+1；
func pay2Winner(loser *PlayerHolder, winner *PlayerHolder, room *Room, mutiple int) {
	horseMultiple := (winner.sctx.horseCount + 1)
	baseMutiple := winner.sctx.greatWinPoints

	if baseMutiple == 0 {
		// 如果没有牌型倍数，则置为1，否则乘法运算结果恒为0
		baseMutiple = 1
	}

	trimMultiple := baseMutiple * horseMultiple
	if room.config.trimMultiple > 0 {
		if trimMultiple > room.config.trimMultiple {
			before := trimMultiple
			trimMultiple = room.config.trimMultiple

			room.cl.Printf("pay2Winner, trim %d to %d", before, trimMultiple)
		}
	}

	score2Pay := room.config.baseScore * trimMultiple * mutiple
	room.cl.Printf("%s pay2Winner %s, score2Pay:%d = baseScore:%d X baseMutiple:%d X horseMultiple:%d X mutiple: %d\n",
		loser.userID(), winner.userID(), score2Pay, room.config.baseScore, baseMutiple, horseMultiple, mutiple)

	loser.sctx.getPayTarget(winner).totalWinScore -= score2Pay
	winner.sctx.getPayTarget(loser).totalWinScore += score2Pay

	room.cl.Printf("loser:%d pay score %d 2 winner %d\n", loser.chairID, score2Pay, winner.chairID)
}

// calcKongMultiple 计算杠分
func calcKongMultiple(s *SPlaying, player *PlayerHolder) {
	log.Println("calcKongMultiple for player:", player.chairID)
	tiles := player.tiles
	kongMelds := tiles.kongMelds()

	if len(kongMelds) < 1 {
		return
	}

	for _, m := range kongMelds {
		switch m.mt {
		case mahjong.MeldType_enumMeldTypeConcealedKong:
			payConcealedKong(player, s)
			break
		case mahjong.MeldType_enumMeldTypeExposedKong:
			payExposedKong(player, m, s)
			break
		case mahjong.MeldType_enumMeldTypeTriplet2Kong:
			payTriplet2Kong(player, s)
			break
		}
	}
}

// payConcealedKong 暗杠计分, 每人出2分，共收6分
func payConcealedKong(konger *PlayerHolder, s *SPlaying) {
	// 暗杠每一个人都要给予分数

	for _, p := range s.players {
		if p == konger {
			continue
		}

		multiple := 2

		p.sctx.getPayTarget(konger).kongMultiple -= multiple
		konger.sctx.getPayTarget(p).kongMultiple += multiple
		log.Printf("player :%d pay concealed kong multiple %d to %d\n", p.chairID, multiple, konger.chairID)
	}
}

// payExposedKong 加杠计分, 每人出1分，共3分
func payTriplet2Kong(konger *PlayerHolder, s *SPlaying) {
	// 明杠每一个人都要给予分数
	for _, p := range s.players {
		if p == konger {
			continue
		}

		multiple := 1

		p.sctx.getPayTarget(konger).kongMultiple -= multiple
		konger.sctx.getPayTarget(p).kongMultiple += multiple

		log.Printf("player :%d pay triplet2Kong multiple %d to %d\n", p.chairID, multiple, konger.chairID)
	}
}

// payExposedKong 明杠计分，放杠者出3分，共收3分
func payExposedKong(konger *PlayerHolder, m *Meld, s *SPlaying) {
	loser := s.tileMgr.getContributor(konger, m)
	if loser == nil || loser == konger {
		log.Panicf("payExposedKong failed, can't find contributor player\n")
		return
	}

	multiple := 3

	loser.sctx.getPayTarget(konger).kongMultiple -= multiple
	konger.sctx.getPayTarget(loser).kongMultiple += multiple

	log.Printf("player :%d pay exposed kong multiple %d to %d\n", loser.chairID, multiple, konger.chairID)
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
	roomConfig := s.room.config

	// 汇总各种得分
	for _, p := range orderPlayers {
		for _, pc := range p.sctx.orderPlayerSctxs {
			if pc.hasCalc {
				continue
			}

			targetPlayer := pc.target

			if pc.kongMultiple != 0 {
				before := pc.totalWinScore
				add := pc.kongMultiple * roomConfig.baseScore
				pc.totalWinScore += add

				log.Printf("doFinalPay: player:%d take %d kongMultiple from palyer:%d, totalWinScore %d=>%d\n",
					p.chairID, add, targetPlayer.chairID, before, pc.totalWinScore)
			}

			pc.hasCalc = true
			pc2 := pc.target.sctx.getPayTarget(p)
			pc2.totalWinScore = -pc.totalWinScore
			pc2.hasCalc = true
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

		// 计算马牌
		calcHorse(winner, s)

		winner.sctx.winType = int(mahjong.HandOverType_enumHandOverType_Win_Chuck)
	}

	// 计算分数
	// 为每一个可以胡牌者，放铳者都需要付出对方的分数
	for _, p := range s.players {
		if !p.tiles.winAble() || p == chucker {
			continue
		}

		var winner = p
		pay2Winner(chucker, winner, s.room, 1)
	}

	// 计算每个人的杠牌得分
	for _, p := range s.players {
		calcKongMultiple(s, p)
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

// 计算中马
func calcHorse(winner *PlayerHolder, s *SPlaying) {
	lastTile := winner.tiles.latestHandTile()
	horseType := lastTile.horseType()
	horseCount := s.room.config.horseCount

	horseTileMatchCount := 0
	if s.horseTiles == nil {
		horseTiles := s.tileMgr.drawHorseTiles(horseCount)
		// 记录马牌
		s.horseTiles = horseTiles

		horseTileIDs := make([]int, len(horseTiles))
		for i, t := range horseTiles {
			horseTileIDs[i] = t.tileID
		}

		s.cl.Printf("calcHorse, winner:%s, horseTiles:%+v", winner.userID(), horseTileIDs)

		// 记录马牌列表
		s.lctx.addActionWithTiles(nil, horseTileIDs, mahjong.ActionType_enumActionType_CustomA, 0, mahjong.SRFlags_SRUserReplyOnly, 0)
	}

	for _, ht := range s.horseTiles {
		if ht.horseType() == horseType {
			horseTileMatchCount++
		}
	}

	// 保存中马个数
	winner.sctx.horseCount = horseTileMatchCount
	s.cl.Printf("winner:%s, horseTileMatchCount:%d", winner.userID(), horseTileMatchCount)
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
	sctx := winner.sctx
	roomConfig := s.room.config

	// 计算大小胡
	basicScoreCalc(s, winner, true)

	isAfterExposedKong := (sctx.greatWinType&(int(GreatWinType_AfterExposedKong)) != 0)
	isKongerPayForAll := isAfterExposedKong && roomConfig.afterKongChuckerPayForAll

	// 计算中马
	calcHorse(winner, s)

	if isKongerPayForAll {
		// 一人支付
		konger := s.lctx.kongerOf(winner, s.room)
		pay2Winner(konger, winner, s.room, len(s.players)-1)
	} else {
		// 其他人各自付分
		for _, p := range s.players {
			if p == winner {
				continue
			}

			pay2Winner(p, winner, s.room, 1)
		}
	}

	// 计算每个人的杠牌得分
	for _, p := range s.players {
		calcKongMultiple(s, p)
	}

	orderPlayers := make([]*PlayerHolder, 0, len(s.players))
	orderPlayers = append(orderPlayers, winner)
	orderPlayers = append(orderPlayers, s.tileMgr.getOrderPlayers(winner)...)
	doFinalPay(s, orderPlayers)
}
