package dfmahjong

import "mahjong"

func isGreatWinTypePongkong(winType int) bool {
	return (winType&int(GreatWinType_enumGreatWinType_PongKongWithFlowerNoMeld)) != 0 ||
		(winType&int(GreatWinType_enumGreatWinType_PongKong)) != 0
}

// calcGreatWinTileType 计算跟行牌无关的牌型大胡
// 1.清一色
// 2.混一色
// 3.碰碰胡
// 4.七对
func calcGreatWinTileType(s *SPlaying, player *PlayerHolder) (int, float32) {
	var tiles = player.tiles
	var points float32
	var winType = 0
	var roomConfig = s.room.config

	// 碰碰胡：由4副刻子或杠，和1对相同的牌组成的胡牌
	if tiles.isAllTripletOrKong() {
		// 无下地，有花
		if tiles.exposedMeldCount() == 0 && tiles.flowerTileCount() > 0 {
			winType |= int(GreatWinType_enumGreatWinType_PongKongWithFlowerNoMeld)
			var gp = roomConfig.greatWinPointMap(GreatWinType_enumGreatWinType_PongKongWithFlowerNoMeld)
			points += gp
			s.cl.Println("GWT:PongKong, with flower no meld:", gp)
		} else {
			winType |= int(GreatWinType_enumGreatWinType_PongKong)
			var gp = roomConfig.greatWinPointMap(GreatWinType_enumGreatWinType_PongKong)
			points += gp
			s.cl.Println("GWT:PongKong, normal:", gp)
		}
	}

	if tiles.suitTypeCount() == 1 {
		if tiles.honorTypeCount() > 0 {
			//--由于箭牌（中发白）已经归类为花牌，因此不需要考虑箭牌数量
			//if tiles.dragonTypeCount() == 0 {
			if tiles.exposedMeldCount() > 0 {
				// 混一色：万条筒其中一种与风牌结合一起，且产生的胡牌,且有下地
				winType |= int(GreatWinType_enumGreatWinType_MixSameWithMeld)
				var gp = roomConfig.greatWinPointMap(GreatWinType_enumGreatWinType_MixSameWithMeld)
				points += gp
				s.cl.Println("GWT:MixedSame, with meld:", gp)
			} else {
				if tiles.flowerTileCount() > 0 {
					// 混一色：万条筒其中一种与风牌结合一起，且产生的胡牌。没有落地，且有花
					winType |= int(GreatWinType_enumGreatWinType_MixSameWithFlowerNoMeld)
					var gp = roomConfig.greatWinPointMap(GreatWinType_enumGreatWinType_MixSameWithFlowerNoMeld)
					points += gp
					s.cl.Println("GWT:MixedSame, with flower no meld:", gp)
				} else {
					// 混一色：万条筒其中一种与风牌结合一起，且产生的胡牌。
					winType |= int(GreatWinType_enumGreatWinType_MixedSame)
					var gp = roomConfig.greatWinPointMap(GreatWinType_enumGreatWinType_MixedSame)
					points += gp
					s.cl.Println("GWT:MixedSame, normal:", gp)
				}
			}
			//}
		} else {
			if tiles.exposedMeldCount() > 0 {
				// 清一色：一色牌组成的胡牌。有落地
				winType |= int(GreatWinType_enumGreatWinType_PureSameWithMeld)
				var gp = roomConfig.greatWinPointMap(GreatWinType_enumGreatWinType_PureSameWithMeld)
				points += gp
				s.cl.Println("GWT:PureSame, with meld:", gp)
			} else {
				if tiles.flowerTileCount() > 0 {
					// 清一色：一色牌组成的胡牌。没有落地，且有花
					winType |= int(GreatWinType_enumGreatWinType_PureSameWithFlowerNoMeld)
					var gp = roomConfig.greatWinPointMap(GreatWinType_enumGreatWinType_PureSameWithFlowerNoMeld)
					points += gp
					s.cl.Println("GWT:PureSame, with flower no meld:", gp)
				} else {
					// 清一色：一色牌组成的胡牌。
					winType |= int(GreatWinType_enumGreatWinType_PureSame)
					var gp = roomConfig.greatWinPointMap(GreatWinType_enumGreatWinType_PureSame)
					points += gp
					s.cl.Println("GWT:PureSame:", gp)
				}
			}
		}
	}

	// 七对，豪华七对
	var st = tiles.calc7Pair()
	switch st {
	case GreatWinType_enumGreatWinType_GreatSevenPair:
		winType |= int(GreatWinType_enumGreatWinType_GreatSevenPair)
		var gp = roomConfig.greatWinPointMap(GreatWinType_enumGreatWinType_GreatSevenPair)
		points += gp
		s.cl.Println("GWT:GreatSevenPair:", gp)
		break
	case GreatWinType_enumGreatWinType_SevenPair:
		winType |= int(GreatWinType_enumGreatWinType_SevenPair)
		var gp = roomConfig.greatWinPointMap(GreatWinType_enumGreatWinType_SevenPair)
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
	var roomConfig = s.room.config

	if !tiles.winAble() {
		s.cl.Panic("calcGreatWinning, not winable")
		return
	}

	// 计算牌型性质的大胡
	winType, points = calcGreatWinTileType(s, player)

	// var selfDrawn = s.lctx.isSelfDraw(player)
	if !selfDrawn && s.lctx.isRobKong() {
		// 如果是最后动作是加杠，则表明是抢杠胡
		winType |= int(GreatWinType_enumGreatWinType_RobKong)
		var gp = roomConfig.greatWinPointMap(GreatWinType_enumGreatWinType_RobKong)
		points += gp
		s.cl.Println("GWT:RobKong:", gp)
	}

	// 独钓：吃碰杠一起12只，剩余1只，胡剩余的1只。
	if tiles.meldCount() == 4 && tiles.tileCountInHand() == 2 {
		winType |= int(GreatWinType_enumGreatWinType_ChowPongKong)
		var gp = roomConfig.greatWinPointMap(GreatWinType_enumGreatWinType_ChowPongKong)
		points += gp
		s.cl.Println("GWT:ChowPongKong:", gp)
	}

	// 海底捞月：摸牌池最后一张牌，并产生胡牌。
	// 需求变更：海底捞只要牌堆空了不管自摸还是吃铳都算
	if s.tileMgr.wallEmpty() {
		winType |= int(GreatWinType_enumGreatWinType_FinalDraw)
		var gp = roomConfig.greatWinPointMap(GreatWinType_enumGreatWinType_FinalDraw)
		points += gp
		s.cl.Println("GWT:FinalDraw:", gp)
	}

	// 大门清：胡牌时，无吃碰杠，且没有抓过花
	// 修正：原本是exposedMeldCount()，也即是暗杠不破坏大门清
	//  但是现在确认，暗杠破坏大门清，因此改为meldCount
	if tiles.meldCount() == 0 && tiles.flowerTileCount() == 0 {
		winType |= int(GreatWinType_enumGreatWinType_ClearFront)
		var gp = roomConfig.greatWinPointMap(GreatWinType_enumGreatWinType_ClearFront)
		points += gp
		s.cl.Println("GWT:ClearFront:", gp)
	}

	// 天胡
	if player == s.room.bankerPlayer() && player.hStatis.actionCounter == 1 && selfDrawn {
		winType |= int(GreatWinType_enumGreatWinType_Heaven)
		var gp = roomConfig.greatWinPointMap(GreatWinType_enumGreatWinType_Heaven)
		points += gp
		s.cl.Println("GWT:Heaven:", gp)
	}

	// 立听胡牌 player.hStatis.isRichi
	if player.hStatis.isRichi {
		winType |= int(GreatWinType_enumGreatWinType_Richi)
		var gp = roomConfig.greatWinPointMap(GreatWinType_enumGreatWinType_Richi)
		points += gp
		s.cl.Println("GWT:Richi:", gp)
	}

	// 放铳者立听
	// if !selfDrawn {
	// 	srAction := s.lctx.getLastNonDrawAction()
	// 	if srAction != nil {
	// 		chucker := s.room.getPlayerByChairID(int(srAction.GetChairID()))
	// 		if chucker != nil && chucker.hStatis.isRichi {
	// 			winType |= int(GreatWinType_enumGreatWinType_OpponentsRichi)
	// 			var gp = roomConfig.greatWinPointMap(GreatWinType_enumGreatWinType_OpponentsRichi)
	// 			points += gp
	// 			s.cl.Println("GWT:Opponents Richi:", gp)
	// 		}
	// 	}
	// }

	// 暗杠胡：手牌里有3只一样的牌，同时胡第4只1样的牌。（必须自摸）
	// 注意不是岭上开花
	if selfDrawn && tiles.tileCountInHandOf(tiles.latestHandTile().tileID) == 4 {
		winType |= int(GreatWinType_enumGreatWinType_AfterConcealedKong)
		var gp = roomConfig.greatWinPointMap(GreatWinType_enumGreatWinType_AfterConcealedKong)
		points += gp
		s.cl.Println("GWT:AfterConcealedKong:", gp)
	}

	// 明杠胡：碰牌后，依然胡碰的那只牌。（必须自摸）
	// 注意不是岭上开花
	if selfDrawn && tiles.hasPongOf(tiles.latestHandTile().tileID) {
		winType |= int(GreatWinType_enumGreatWinType_AfterExposedKong)
		var gp = roomConfig.greatWinPointMap(GreatWinType_enumGreatWinType_AfterExposedKong)
		points += gp
		s.cl.Println("GWT:AfterExposedKong:", gp)
	}

	// 大胡连庄
	// if roomConfig.isDoubleScoreWhenContinuousBanker && player.gStatis.isContinuousBanker {
	// 	var prePoints = points
	// 	points = prePoints * roomConfig.roomScoreConfig.greatWinContinuousBankerMultiple
	// 	sc.isContinuousBanker = true
	// 	s.cl.Printf("GWT:continuous banker, points:%f x %f => %f\n", prePoints, roomConfig.roomScoreConfig.greatWinContinuousBankerMultiple, points)
	// }

	sc.greatWinType = winType
	pointsTrim := roomConfig.roomScoreConfig.greatWinPointTrimFunc(points, roomConfig.roomScoreConfig.maxGreatWinPoints)
	s.cl.Printf("great win points trimFunc: %f=>%f\n", points, pointsTrim)
	sc.fGreatWinPoints = pointsTrim

	// fTrimGreatWinPoints 仅仅用于客户端显示，不参与计算
	if pointsTrim > roomConfig.roomScoreConfig.maxGreatWinPoints {
		sc.fTrimGreatWinPoints = roomConfig.roomScoreConfig.maxGreatWinPoints
	} else {
		sc.fTrimGreatWinPoints = pointsTrim
	}

	s.cl.Printf("great win point:%f, type:%d\n", points, winType)
}

func calcPay2WinnerOfMultipleAndTrim(loser *PlayerHolder, winner *PlayerHolder, room *Room, winScore float32, useGreatWinSpace bool) int {
	roomConfig := room.config
	roomScoreConfig := roomConfig.roomScoreConfig
	var score2PayUnTrim float32
	var limitScore float32
	var scoreAfterTrimFunc float32
	var score2Pay int

	// 使用大胡区间
	if useGreatWinSpace {
		room.cl.Println("winner use greatwin, winScore:", winScore)
		greatWinScore := roomScoreConfig.preGreatWinTrimFunc(winScore)
		room.cl.Printf("score before preTrimFunc:%f, after:%f\n", winScore, greatWinScore)

		if roomConfig.isDoubleScoreWhenContinuousBanker {
			if winner.gStatis.isContinuousBanker {
				winner.sctx.isContinuousBanker = true
				winner.sctx.fContinuousBankerMultiple = roomScoreConfig.greatWinContinuousBankerMultiple
				score2PayUnTrim = greatWinScore * roomScoreConfig.greatWinContinuousBankerMultiple
				room.cl.Println("winner is greatwin and continuous banker, score2PayUnTrim:", score2PayUnTrim)
			} else if loser.gStatis.isContinuousBanker {
				loser.sctx.isContinuousBanker = true
				loser.sctx.fContinuousBankerMultiple = roomScoreConfig.greatWinContinuousBankerMultiple
				score2PayUnTrim = greatWinScore * roomScoreConfig.greatWinContinuousBankerMultiple

				room.cl.Println("winner is greatwin, loser is continuous banker, score2PayUnTrim:", score2PayUnTrim)
			} else {
				score2PayUnTrim = greatWinScore
				room.cl.Println("winner is greatwin, no continuous banker, score2PayUnTrim:", score2PayUnTrim)
			}

			if room.markup > 0 {
				markupMultiple := float32(1.5)
				score2PayUnTrim = score2PayUnTrim * markupMultiple
				winner.sctx.fContinuousBankerMultiple = markupMultiple
				room.cl.Println("room is in markup state, score2PayUnTrim double:", score2PayUnTrim)
			}
		} else {
			score2PayUnTrim = greatWinScore
			room.cl.Println("winner is greatwin, no continuous banker config, score2PayUnTrim:", score2PayUnTrim)
		}

		limitScore = roomScoreConfig.greatWinScore * roomScoreConfig.maxGreatWinPoints
		// if score2PayUnTrim > limitGreatWinScore {
		room.cl.Printf("score2PayUnTrim %f, limit to:%f\n", score2PayUnTrim, limitScore)
		// 	score2PayUnTrim = limitGreatWinScore
		// }
		scoreBeforeTrimFunc := (score2PayUnTrim)
		scoreAfterTrimFunc = roomScoreConfig.postGreatWinTrimFunc(scoreBeforeTrimFunc)
		room.cl.Printf("score before postTrimFunc:%f, after:%f\n", scoreBeforeTrimFunc, scoreAfterTrimFunc)
	} else {
		// 小胡
		room.cl.Println("winner use miniwin, winScore:", winScore)
		fMiniWinUnTrimScore := roomScoreConfig.preMiniWinTrimFunc(winScore)
		room.cl.Printf("score before preTrimFunc:%f, after:%f\n", winScore, fMiniWinUnTrimScore)
		if roomConfig.isDoubleScoreWhenContinuousBanker {
			if winner.gStatis.isContinuousBanker {
				winner.sctx.isContinuousBanker = true
				winner.sctx.fContinuousBankerMultiple = roomScoreConfig.miniWinContinuousBankerMultiple
				score2PayUnTrim = fMiniWinUnTrimScore * roomScoreConfig.miniWinContinuousBankerMultiple

				limitScore = roomScoreConfig.greatWinScore * roomScoreConfig.miniWinLimitMultipleContinuousBanker

				room.cl.Println("winner is miniwin and continuous banker, score2PayUnTrim:", score2PayUnTrim)
			} else if loser.gStatis.isContinuousBanker {
				loser.sctx.isContinuousBanker = true
				loser.sctx.fContinuousBankerMultiple = roomScoreConfig.miniWinContinuousBankerMultiple
				score2PayUnTrim = fMiniWinUnTrimScore * roomScoreConfig.miniWinContinuousBankerMultiple

				limitScore = roomScoreConfig.greatWinScore * roomScoreConfig.miniWinLimitMultipleContinuousBanker

				room.cl.Println("winner is miniwin, loser is continuous banker, score2PayUnTrim:", score2PayUnTrim)
			} else {
				score2PayUnTrim = fMiniWinUnTrimScore
				limitScore = roomScoreConfig.greatWinScore * roomScoreConfig.miniWinLimitMultipleNormal
				room.cl.Println("winner is miniwin, no continuous banker, score2PayUnTrim:", score2PayUnTrim)
			}

			// if score2PayUnTrim > limitMiniWinScore {
			room.cl.Printf("score2PayUnTrim %f, with continuous banker config limit to:%f\n", score2PayUnTrim, limitScore)

			if room.markup > 0 {
				markupMultiple := float32(2.0)
				score2PayUnTrim = score2PayUnTrim * markupMultiple
				winner.sctx.fContinuousBankerMultiple = markupMultiple
				limitScore = roomScoreConfig.greatWinScore * roomScoreConfig.miniWinLimitMultipleContinuousBanker
				room.cl.Printf("room is in markup state, score2PayUnTrim double:%f, limitScore:%f\n", score2PayUnTrim, limitScore)
			}
			// 	score2PayUnTrim = limitMiniWinScore
			// }
		} else {
			score2PayUnTrim = fMiniWinUnTrimScore
			limitScore = roomScoreConfig.greatWinScore * roomScoreConfig.miniWinLimitMultipleNormal
			// if score2PayUnTrim > limitMiniWinScore {
			room.cl.Printf("score2PayUnTrim %f, without continuous banker config limit to:%f\n", score2PayUnTrim, limitScore)
			// 	score2PayUnTrim = limitMiniWinScore
			// }
		}

		scoreBeforeTrimFunc := (score2PayUnTrim)
		scoreAfterTrimFunc = roomScoreConfig.postMiniWinTrimFunc(scoreBeforeTrimFunc)
		room.cl.Printf("score before postTrimFunc:%f, after:%f\n", scoreBeforeTrimFunc, scoreAfterTrimFunc)
	}

	if scoreAfterTrimFunc > limitScore {
		room.cl.Printf("scoreAfterTrimFunc %f > limitScore %f, trim to limitScore\n", scoreAfterTrimFunc, limitScore)
		scoreAfterTrimFunc = limitScore
	}

	score2Pay = roundFloat32(scoreAfterTrimFunc)
	return score2Pay
}

func calcPay2Winner(loser *PlayerHolder, winner *PlayerHolder, room *Room) int {
	room.cl.Printf("calcPay2Winner, greatWinTrimType:%d, loser %d, winner %d\n", room.config.greatWinTrimType, loser.chairID, winner.chairID)

	roomConfig := room.config
	roomScoreConfig := roomConfig.roomScoreConfig

	score2Pay := 0
	useGreatWin := false
	winScore := float32(0)

	if winner.sctx.greatWinType != 0 {
		useGreatWin = true
		winScore = winner.sctx.fGreatWinPoints * roomScoreConfig.greatWinScore
	} else {
		winScore = winner.sctx.fMiniWinUnTrimScore
	}

	// 如果loser是报听状态，则需要赔付一个辣子，而且裁剪按照大胡来裁剪
	if loser.hStatis.isRichi {
		if useGreatWin {
			winScore += roomScoreConfig.greatWinScore
			room.cl.Printf("loser %d is in richi state, pay addition 1 point %f to winner %d\n",
				loser.chairID, roomScoreConfig.greatWinScore, winner.chairID)
		} else {
			useGreatWin = true
			winScore = roomScoreConfig.greatWinScore
			room.cl.Printf("loser %d is in richi state, reset winner %d miniwin to 1 point %f \n",
				loser.chairID, winner.chairID, roomScoreConfig.greatWinScore)
		}

		loser.sctx.isRichiPay1P = true
	}

	score2Pay = calcPay2WinnerOfMultipleAndTrim(loser, winner, room, winScore, useGreatWin)
	return score2Pay
}

func pay2Winner(loser *PlayerHolder, winner *PlayerHolder, room *Room) {

	score2Pay := calcPay2Winner(loser, winner, room)

	loser.sctx.getPayTarget(winner).totalWinScore -= score2Pay
	winner.sctx.getPayTarget(loser).totalWinScore += score2Pay

	room.cl.Printf("loser:%d pay score %d 2 winner %d\n", loser.chairID, score2Pay, winner.chairID)
}

// calcWinChuckGreatWinType 用于计算吃铳胡时，是否可以形成大胡
func calcWinChuckGreatWinType(s *SPlaying, player *PlayerHolder, chuckTile *Tile) bool {
	oldSctx := player.sctx
	var sctx = &ScoreContext{}
	player.sctx = sctx
	tiles := player.tiles
	tiles.temporaryHandAdd(chuckTile)
	calcGreatWinning(s, player, false)
	player.sctx = nil
	tiles.temporaryHandRemove(chuckTile)
	player.sctx = oldSctx
	return sctx.greatWinType != int(GreatWinType_enumGreatWinType_None) &&
		sctx.greatWinType != int(GreatWinType_enumGreatWinType_OpponentsRichi)
}

// readyHandPongKong 听牌能听出碰碰胡
func readyHandPongKong(s *SPlaying, player *PlayerHolder) bool {
	s.cl.Printf("readyHandGreatWin, player %d \n", player.chairID)
	tiles := player.tiles
	tiles.hand2Slots()
	slots := tiles.slots

	tileIDs := make([]int, 0, 34)

	for tryTile := MAN; tryTile < FlowerBegin; tryTile++ {
		slots[tryTile]++
		winAble := isWinable(slots)

		if winAble {
			tileIDs = append(tileIDs, tryTile)
		}
		slots[tryTile]--
	}

	for _, tid := range tileIDs {
		chuckTile := &Tile{tileID: tid}
		tiles.temporaryHandAdd(chuckTile)
		winType, _ := calcGreatWinTileType(s, player)
		tiles.temporaryHandRemove(chuckTile)

		if isGreatWinTypePongkong(winType) {
			return true
		}
	}

	return false
}

// calcMiniWinning 小胡计分的公式为：总分=（1底分+花分）x2连庄x2杠开/杠冲x2小门清x2自摸+墩子分
// 其中，
// 选自摸加双，如果自摸：（花分+底分）X2+墩子分
//    选连庄，如果胡牌：（花分+底分）X2+墩子分
//    选连庄自摸，如果自摸胡牌：（花分+底分）X2x2+墩子分
//    底分为牌局前预设.
func calcMiniWinning(s *SPlaying, player *PlayerHolder, selfDraw bool) {
	//var selfDraw = s.lctx.isSelfDraw(player)
	var tiles = player.tiles
	sc := player.sctx
	roomConfig := s.room.config
	roomScoreConfig := roomConfig.roomScoreConfig

	if !tiles.winAble() {
		s.cl.Panic("calcMiniWinning, not winable")
		return
	}

	var multiple float32 = 1.0
	var miniWinType = 0

	// 花分（包括花牌的花分，以及碰杠，墩子花分）
	var flowerCount = tiles.allFlowerScoreCount(selfDraw)

	// 如果没有任何花，奖励10分
	if tiles.flowerTileCount() < 1 && tiles.meldCount() > 0 {
		flowerCount += 10
		miniWinType |= int(MiniWinType_enumMiniWinType_NoFlowers)
		s.cl.Println("MWT:NoFlowers")
	}

	var flowerScore = float32(flowerCount) * roomScoreConfig.scorePerFlower

	var baseScore = roomScoreConfig.miniWinBaseScore
	var doubleAble = flowerScore + baseScore

	s.cl.Printf("mini-Calc, flowerScore:%f*%d=>%f, baseScore:%f, doubleAble:%f\n", roomScoreConfig.scorePerFlower,
		flowerCount, flowerScore, baseScore, doubleAble)

	sc.fMiniWinBasicScore = baseScore
	sc.fMiniWinFlowerScore = flowerScore
	// 选自摸加双，如果自摸：（花分+底分）X2
	// 自摸X2，为可选玩法的自摸加双，当玩家没勾选时，则自摸X2不计入公式
	if s.room.config.isDoubleScoreWhenSelfDrawn && selfDraw {
		multiple = multiple * 2.0
		miniWinType |= int(MiniWinType_enumMiniWinType_SelfDraw)
		s.cl.Println("MWT:SelfDraw")
	}

	// 连庄统一在pay2Winner计算

	// 杠冲x2，小胡胡牌方式为，某玩家杠牌后出一只，我刚好胡牌。或者对方补花后出的牌，我刚好胡牌
	// 杠牌者放铳
	if ok := s.lctx.isXKong2Discarded(player); ok {
		multiple = multiple * 2.0
		miniWinType |= int(MiniWinType_enumMiniWinType_Kong2Discard)
		s.cl.Println("MWT:Kong2Discard")
	}

	// 杠开X2，小胡胡牌的方式为，玩家杠牌后，摸牌后胡牌。或者补花后，摸牌胡牌。注意：因为这个杠开带自摸属性，如果选了自摸还要再X2。
	if s.lctx.isXKong2SelfDraw(player) {
		multiple = multiple * 2.0
		miniWinType |= int(MiniWinType_enumMiniWinType_Kong2SelfDraw)
		s.cl.Println("MWT:Kong2SelfDraw")
	}

	// 小门清X2，小胡胡牌方式为，胡牌者门前有花，但无吃碰杠。
	// 修正：之前是调用meldCount()，也即是暗杠一样取消小门清
	//  经确认，暗杠不影响小门清，因此，改为调用exposedMeldCount()
	// 2019年1月修改： && tiles.flowerTileCount() > 0，去掉这个条件约束，也就是只要没有除暗杠之外的落地牌组
	//  即算是小门清
	if tiles.exposedMeldCount() == 0 {
		multiple = multiple * 2.0
		miniWinType |= int(MiniWinType_enumMiniWinType_SecondFrontClear)
		s.cl.Println("MWT:SecondFrontClear")
	}

	s.cl.Printf("mini win multiple:%f, win type:%d\n", multiple, miniWinType)

	var sum = doubleAble * multiple

	sc.fMiniMultiple = multiple

	sc.fMiniWinUnTrimScore = (sum)
	sc.miniWinType = miniWinType

	s.cl.Printf("mini-calc, mini win un-trim score:%f\n", sum)
}

// calcSpecialFlowerScore 计算墩子分
func calcSpecialFlowerScore(s *SPlaying, player *PlayerHolder) {
	s.cl.Printf("calcSpecialFlowerScore, player chairID:%d, userID:%s\n", player.chairID, player.userID())
	var tiles = player.tiles
	sc := player.sctx

	var ckongSpecialCount = tiles.concealedKongCount()
	var ckongSpecialScore = s.room.config.dunzi4ConcealedKong * ckongSpecialCount
	s.cl.Printf("SpecialX, ckong %d *%d=>%d\n", ckongSpecialCount, s.room.config.dunzi4ConcealedKong, ckongSpecialScore)

	var specialScore = ckongSpecialScore

	var ekongSpecialCount = tiles.exposedKongCount()
	var ekongSpecialScore = s.room.config.dunzi4ExposedKong * ekongSpecialCount
	s.cl.Printf("SpecialX, ekong %d *%d=>%d\n", ekongSpecialCount, s.room.config.dunzi4ExposedKong, ekongSpecialScore)
	specialScore += ekongSpecialScore

	var f4SpecialCount = tiles.quadFlowerCount()
	var f4SpecialScore = s.room.config.dunzi4QuadFlower * f4SpecialCount
	s.cl.Printf("SpecialX, f4 %d *%d=>%d\n", f4SpecialCount, s.room.config.dunzi4QuadFlower, f4SpecialScore)
	specialScore += f4SpecialScore

	sc.specialScore = (specialScore)
	s.cl.Printf("SpecialX: %d\n", specialScore)
}

// calcFakers 计算包牌关系
// 函数名字Faker，骗子，包牌事实上类似于骗子，也即是有意让某人胡牌
// 需求变更：吃椪杠3以上，而且是大胡或者能听出大胡，则是包牌关系
func calcFakers(s *SPlaying, player *PlayerHolder) []*PlayerHolder {
	fakers := make([]*PlayerHolder, 0, len(s.players))

	for _, p := range s.players {
		if p == player {
			continue
		}

		if fakeLike(s, p, player) || fakeLike(s, player, p) {
			fakers = append(fakers, p)
		}
	}

	return fakers
}

// fakeLike 包牌关系。
// 需求变更：
// 1. 只要吃椪杠4次，必包
// 2. 只要吃椪杠3次，且吃椪杠别人者，是清一色或者混一色牌子（注意不一定成型），包牌
// 3. 只要椪杠3次，且椪杠别人者，是碰碰胡或者听成碰碰胡，包牌
func fakeLike(s *SPlaying, p1 *PlayerHolder, p2 *PlayerHolder) bool {
	s.cl.Printf("fakeLike-calc p1:%d, p2:%d\n", p1.chairID, p2.chairID)
	var tiles2 = p2.tiles
	//var tiles1 = p1.tiles

	var chow = tiles2.chowCountFrom(p1)
	var pong = tiles2.pongCountFrom(p1)
	var kong = tiles2.kongCountFrom(p1)

	xcount := pong + chow + kong

	if xcount < 3 {
		return false
	}

	//  只要吃椪杠4次，必包
	if xcount == 4 {
		s.cl.Printf("fakeLike-result p1:%d, p2:%d\n", p1.chairID, p2.chairID)
		return true
	}

	// 只要吃椪杠3次，且吃椪杠别人者，是清一色或者混一色牌子（注意不一定成型），包牌
	if tiles2.suitTypeCount() == 1 {
		s.cl.Printf("fakeLike-result p1:%d, p2:%d\n", p1.chairID, p2.chairID)
		return true
	}

	// s.cl.Println("fakeLike-calc, p2 greatwinType:", p2.sctx.greatWinType)
	// 只要椪杠3次，且椪杠别人者，是碰碰胡或者听成碰碰胡，包牌
	fake := false
	if isGreatWinTypePongkong(p2.sctx.greatWinType) {
		fake = true
	} else if p2.tiles.agariTileCount() < 14 {
		fake = readyHandPongKong(s, p2)
	}

	if fake {
		s.cl.Printf("fakeLike-result p1:%d, p2:%d\n", p1.chairID, p2.chairID)
		return true
	}

	return false
}

func calcFakersScores(s *SPlaying, fakers []*PlayerHolder, winner *PlayerHolder) {
	// 包牌者需要多付出其他人损失的分数
	for _, faker := range fakers {
		// victimLoseScore := 0
		// ABC3个人，C放炮B胡牌，AB包牌，那么：
		// 如果A连庄，A付出包牌分得考虑连庄
		// 如果B连庄，A付出包牌分得考虑连庄
		// 如果C连庄，A付出包牌分得不，不，不考虑连庄
		// 检查赢牌者的收益列表，对于每一个输者，包牌者都要跟着输多一次
		for _, pc := range winner.sctx.orderPlayerSctxs {
			// 分数小于0，表示victim需要支付给fwiner，因此faker也需要支付给fwiner
			if pc.totalWinScore > 0 {
				floser := pc.target
				if floser == faker {
					continue
				}

				victimLoseScore := calcPay2Winner(faker, winner, s.room)

				faker.sctx.getPayTarget(winner).fakeWinScore -= (victimLoseScore)
				winner.sctx.getPayTarget(faker).fakeWinScore += (victimLoseScore)
				s.cl.Printf("faker:%d pay %d to fwiner:%d\n", faker.chairID, victimLoseScore, winner.chairID)
			}
		}
		//faker.gStatis.isContinuousBanker = oldIsContinuousBanker
	}
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

	var fakers = calcFakers(s, winner)
	if len(fakers) > 0 {
		calcFakersScores(s, fakers, winner)
	}

	// 墩子分
	paySpecialScore(s)

	// 最终计分，此时才考虑进园子保护
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
	shouldPayTrim := s.room.loseProtectTrimPay(loser, shouldPay)
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

func paySpecialScore(s *SPlaying) {
	// 考虑墩子分
	for _, p := range s.players {

		calcSpecialFlowerScore(s, p)

		if p.sctx.specialScore == 0 {
			continue
		}

		for _, p1 := range s.players {
			if p1 != p {
				specialScore := p.sctx.specialScore
				//specialScoreTrim := s.room.loseProtectTrimPay(p1, specialScore)
				//s.cl.Printf("pay special score lose protected trim:%d=>%d\n", specialScore, specialScoreTrim)
				p1.sctx.getPayTarget(p).totalWinScore -= specialScore
				p.sctx.getPayTarget(p1).totalWinScore += specialScore
			}
		}
	}
}

// basicScoreCalc 计算单个玩家基础的得分
func basicScoreCalc(s *SPlaying, player *PlayerHolder, selfDraw bool) {
	sc := player.sctx
	calcGreatWinning(s, player, selfDraw)

	if sc.greatWinType == 0 {
		calcMiniWinning(s, player, selfDraw)
	}

	// if sc.greatWinType != 0 {
	// 	var roomScoreConfig = s.room.config.roomScoreConfig
	// 	sc.scoreWithoutTrim = roundFloat32(sc.fTrimGreatWinPoints * roomScoreConfig.greatWinScore)
	// 	s.cl.Printf("basicScoreCalc great-win, fTrimGreatWinPoints:%f * %f up-trim to baseWinScore:%d\n", sc.fTrimGreatWinPoints,
	// 		roomScoreConfig.greatWinScore, sc.baseWinScore)
	// } else {
	// 	sc.scoreWithoutTrim = roundFloat32(sc.fMiniWinTrimScore)
	// 	s.cl.Printf("basicScoreCalc mini-win, fMiniWinTrimScore:%f, up-trim to baseWinScore:%d\n", sc.fMiniWinTrimScore, sc.baseWinScore)
	// }
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

		// s.cl.Printf("calcFinalResultWithChucker, winChuck chairID:%d, base win score:%d\n", winner.chairID, winner.sctx.baseWinScore)
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

	// 包牌
	// 对每一个赢家，计算和他形成包牌关系的人
	for _, p := range s.players {
		if !p.tiles.winAble() || p == chucker {
			continue
		}

		winner := p
		var fakers = calcFakers(s, winner)
		if len(fakers) < 1 {
			continue
		}

		calcFakersScores(s, fakers, winner)
	}

	// 墩子分
	paySpecialScore(s)

	// 最终计分，此时才考虑进园子保护
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
