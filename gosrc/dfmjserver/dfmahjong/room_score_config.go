package dfmahjong

import (
	"log"
)

// ScoreTrimFunc 分数裁剪函数
type ScoreTrimFunc func(float32) float32

// GreatWinPointsTrimFunc 辣子裁剪
type GreatWinPointsTrimFunc func(float32, float32) float32

// RoomScoreConfig 分数配置
type RoomScoreConfig struct {
	// 每个花牌对应的花分
	scorePerFlower float32

	// 小胡底分
	miniWinBaseScore float32
	// 大胡辣子封顶
	maxGreatWinPoints float32

	// 大胡辣子分
	greatWinScore float32

	// 小胡分封顶：普通情形
	miniWinLimitMultipleNormal float32

	// 小胡分封顶：连庄时情形
	miniWinLimitMultipleContinuousBanker float32

	// 大胡连庄单辣子时，乘以倍数
	greatWinContinuousBankerMultiple float32

	// // 小胡连庄时，乘以倍数
	miniWinContinuousBankerMultiple float32

	// 每一个类型对应的大胡配置
	greatWinPointMap map[GreatWinType]float32

	// 坐园子保护线
	loseProtectScore int

	preMiniWinTrimFunc  ScoreTrimFunc
	postMiniWinTrimFunc ScoreTrimFunc

	preGreatWinTrimFunc  ScoreTrimFunc
	postGreatWinTrimFunc ScoreTrimFunc

	greatWinPointTrimFunc GreatWinPointsTrimFunc
}

func trimFuncNone(value float32) float32 {
	log.Println("trimFuncNone")
	return value
}

func trimFuncCeil(value float32) float32 {
	log.Println("trimFuncCeil")
	v := roundFloat32(value)

	return float32(v)
}

func trimFunc510(value float32) float32 {
	log.Println("trimFunc510")
	v := roundFloat32(value)
	x := v / 10
	y := v % 10
	if y != 0 {
		if y <= 5 {
			return float32(x*10 + 5)
		} else if y > 5 {
			return float32((x + 1) * 10)
		}
	}

	return float32(v)
}

func trimFunc10(value float32) float32 {
	log.Println("trimFunc10")
	v := roundFloat32(value)
	x := v / 10
	y := v % 10
	if y != 0 {
		return float32((x + 1) * 10)
	}

	return float32(v)
}

func trimGreatWinPoint2Ciel(value float32, ceil float32) float32 {
	log.Println("trimGreatWinPoint2Ciel")
	if value < ceil && (ceil-value) < 1.0 {
		return ceil
	}

	return value
}

func trimGreatWinPointNone(value float32, ceil float32) float32 {
	log.Println("trimGreatWinPointNone")
	return value
}

func newRoomScoreConfig(greatWinTrimType int) *RoomScoreConfig {
	rsc := &RoomScoreConfig{}

	rsc.greatWinPointMap = make(map[GreatWinType]float32)
	var greatWinPointMap = rsc.greatWinPointMap

	// 根据辣子封顶类型，填充不同的花分，底分，后面改从redis拉取
	switch greatWinTrimType {
	case 0: // 20/40 辣子区间
		rsc.scorePerFlower = 0.2                       // 每一个花的花分
		rsc.miniWinBaseScore = 1.0                     // 小胡底分
		rsc.maxGreatWinPoints = 2.0                    // 大胡辣子封顶
		rsc.greatWinScore = 20                         // 大胡单辣子分
		rsc.miniWinLimitMultipleNormal = 1.0           // 小胡封顶相当于多少辣子
		rsc.miniWinLimitMultipleContinuousBanker = 1.5 // 在连庄配置下小胡封顶相当于多少辣子
		rsc.greatWinContinuousBankerMultiple = 1.5     // 大胡连庄时倍数
		rsc.miniWinContinuousBankerMultiple = 2.0      // 小胡连庄时乘以倍数

		rsc.loseProtectScore = -80 // 坐园子保护线

		rsc.preMiniWinTrimFunc = trimFuncNone
		rsc.postMiniWinTrimFunc = trimFuncCeil

		rsc.preGreatWinTrimFunc = trimFuncNone
		rsc.postGreatWinTrimFunc = trimFuncCeil

		rsc.greatWinPointTrimFunc = trimGreatWinPointNone

		greatWinPointMap[GreatWinType_enumGreatWinType_ChowPongKong] = 1.0       // 独钓的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_FinalDraw] = 1.0          // 海底捞月的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_PongKong] = 1.0           // 碰碰胡的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_PureSame] = 1.0           // 清一色的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_MixedSame] = 1.0          // 混一色的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_ClearFront] = 1.5         // 大门清的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_SevenPair] = 1.5          // 七对的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_GreatSevenPair] = 2.0     // 豪华大七对的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_Heaven] = 2.0             // 天胡的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_AfterConcealedKong] = 1.0 // 暗杠胡的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_AfterExposedKong] = 1.0   // 明杠胡的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_Richi] = 1.0              // 起手听的辣子数

		greatWinPointMap[GreatWinType_enumGreatWinType_MixSameWithMeld] = 1.0          // 有落地的混一色的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_MixSameWithFlowerNoMeld] = 2.0  // 有花无落地的混一色的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_PureSameWithMeld] = 1.5         // 有落地的清一色的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_PureSameWithFlowerNoMeld] = 2.0 // 有花无落地的清一色的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_PongKongWithFlowerNoMeld] = 2.0 // 有花无落地的碰碰胡辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_RobKong] = 1.0                  // 明杠冲（抢杠胡）

		greatWinPointMap[GreatWinType_enumGreatWinType_OpponentsRichi] = 1.0 // 别人起手报听，也得到了一辣子
		break
	case 1: // 30/60  辣子区间
		rsc.scorePerFlower = 0.3                       // 每一个花的花分
		rsc.miniWinBaseScore = 1.0                     // 小胡底分
		rsc.maxGreatWinPoints = 2.0                    // 大胡辣子封顶
		rsc.greatWinScore = 30                         // 大胡单辣子分
		rsc.miniWinLimitMultipleNormal = 1.0           // 小胡封顶相当于多少辣子
		rsc.miniWinLimitMultipleContinuousBanker = 1.5 // 在连庄配置下小胡封顶相当于多少辣子
		rsc.greatWinContinuousBankerMultiple = 1.5     // 大胡连庄时倍数
		rsc.miniWinContinuousBankerMultiple = 2        // 小胡连庄时乘以倍数

		rsc.loseProtectScore = -120 // 坐园子保护线

		rsc.preMiniWinTrimFunc = trimFuncNone
		rsc.postMiniWinTrimFunc = trimFuncCeil

		rsc.preGreatWinTrimFunc = trimFuncNone
		rsc.postGreatWinTrimFunc = trimFuncCeil

		rsc.greatWinPointTrimFunc = trimGreatWinPointNone

		greatWinPointMap[GreatWinType_enumGreatWinType_ChowPongKong] = 1.0       // 独钓的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_FinalDraw] = 1.0          // 海底捞月的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_PongKong] = 1.0           // 碰碰胡的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_PureSame] = 1.0           // 清一色的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_MixedSame] = 1.0          // 混一色的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_ClearFront] = 1.5         // 大门清的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_SevenPair] = 1.5          // 七对的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_GreatSevenPair] = 2.0     // 豪华大七对的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_Heaven] = 2.0             // 天胡的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_AfterConcealedKong] = 1.0 // 暗杠胡的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_AfterExposedKong] = 1.0   // 明杠胡的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_Richi] = 1.0              // 起手听的辣子数

		greatWinPointMap[GreatWinType_enumGreatWinType_MixSameWithMeld] = 1.0          // 有落地的混一色的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_MixSameWithFlowerNoMeld] = 2.0  // 有花无落地的混一色的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_PureSameWithMeld] = 1.5         // 有落地的清一色的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_PureSameWithFlowerNoMeld] = 2.0 // 有花无落地的清一色的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_PongKongWithFlowerNoMeld] = 2.0 // 有花无落地的碰碰胡辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_RobKong] = 1.0                  // 明杠冲（抢杠胡）

		greatWinPointMap[GreatWinType_enumGreatWinType_OpponentsRichi] = 1.0 // 别人起手报听，也得到了一辣子
		break
	case 2: // 50/100/150 辣子区间
		rsc.scorePerFlower = 0.5                       // 每一个花的花分
		rsc.miniWinBaseScore = 2.0                     // 小胡底分
		rsc.maxGreatWinPoints = 3.0                    // 大胡辣子封顶
		rsc.greatWinScore = 50                         // 大胡单辣子分
		rsc.miniWinLimitMultipleNormal = 1.0           // 小胡封顶相当于多少辣子
		rsc.miniWinLimitMultipleContinuousBanker = 1.5 // 在连庄配置下小胡封顶相当于多少辣子
		rsc.greatWinContinuousBankerMultiple = 1.5     // 大胡连庄时倍数
		rsc.miniWinContinuousBankerMultiple = 2.0      // 小胡连庄时乘以倍数

		rsc.loseProtectScore = -200 // 坐园子保护线

		rsc.preMiniWinTrimFunc = trimFunc510
		rsc.postMiniWinTrimFunc = trimFuncNone

		rsc.preGreatWinTrimFunc = trimFuncNone
		rsc.postGreatWinTrimFunc = trimFunc510

		rsc.greatWinPointTrimFunc = trimGreatWinPoint2Ciel

		greatWinPointMap[GreatWinType_enumGreatWinType_ChowPongKong] = 1.0       // 独钓的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_FinalDraw] = 1.0          // 海底捞月的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_PongKong] = 1.0           // 碰碰胡的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_PureSame] = 1.0           // 清一色的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_MixedSame] = 1.0          // 混一色的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_ClearFront] = 2.0         // 大门清的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_SevenPair] = 2.0          // 七对的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_GreatSevenPair] = 3.0     // 豪华大七对的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_Heaven] = 3.0             // 天胡的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_AfterConcealedKong] = 1.0 // 暗杠胡的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_AfterExposedKong] = 1.0   // 明杠胡的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_Richi] = 1.0              // 起手听的辣子数

		greatWinPointMap[GreatWinType_enumGreatWinType_MixSameWithMeld] = 1.0          // 有落地的混一色的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_MixSameWithFlowerNoMeld] = 2.0  // 有花无落地的混一色的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_PureSameWithMeld] = 2.0         // 有落地的清一色的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_PureSameWithFlowerNoMeld] = 3.0 // 有花无落地的清一色的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_PongKongWithFlowerNoMeld] = 2.0 // 有花无落地的碰碰胡辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_RobKong] = 1.0                  // 明杠冲（抢杠胡）

		greatWinPointMap[GreatWinType_enumGreatWinType_OpponentsRichi] = 1.0 // 别人起手报听，也得到了一辣子

		break
	case 3: // 100/200/300
		rsc.scorePerFlower = 1.0                       // 每一个花的花分
		rsc.miniWinBaseScore = 5.0                     // 小胡底分
		rsc.maxGreatWinPoints = 3.0                    // 大胡辣子封顶
		rsc.greatWinScore = 100                        // 大胡单辣子分
		rsc.miniWinLimitMultipleNormal = 1.0           // 小胡封顶相当于多少辣子
		rsc.miniWinLimitMultipleContinuousBanker = 1.5 // 在连庄配置下小胡封顶相当于多少辣子
		rsc.greatWinContinuousBankerMultiple = 1.5     // 大胡连庄时倍数
		rsc.miniWinContinuousBankerMultiple = 2        // 小胡连庄时乘以倍数

		rsc.loseProtectScore = -300 // 坐园子保护线

		rsc.preMiniWinTrimFunc = trimFunc10
		rsc.postMiniWinTrimFunc = trimFuncNone

		rsc.preGreatWinTrimFunc = trimFuncNone
		rsc.postGreatWinTrimFunc = trimFunc10

		rsc.greatWinPointTrimFunc = trimGreatWinPoint2Ciel

		greatWinPointMap[GreatWinType_enumGreatWinType_ChowPongKong] = 1.0       // 独钓的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_FinalDraw] = 1.0          // 海底捞月的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_PongKong] = 1.0           // 碰碰胡的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_PureSame] = 1.0           // 清一色的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_MixedSame] = 1.0          // 混一色的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_ClearFront] = 2.0         // 大门清的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_SevenPair] = 2.0          // 七对的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_GreatSevenPair] = 3.0     // 豪华大七对的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_Heaven] = 3.0             // 天胡的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_AfterConcealedKong] = 1.0 // 暗杠胡的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_AfterExposedKong] = 1.0   // 明杠胡的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_Richi] = 1.0              // 起手听的辣子数

		greatWinPointMap[GreatWinType_enumGreatWinType_MixSameWithMeld] = 1.0          // 有落地的混一色的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_MixSameWithFlowerNoMeld] = 2.0  // 有花无落地的混一色的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_PureSameWithMeld] = 2.0         // 有落地的清一色的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_PureSameWithFlowerNoMeld] = 3.0 // 有花无落地的清一色的辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_PongKongWithFlowerNoMeld] = 2.0 // 有花无落地的碰碰胡辣子数
		greatWinPointMap[GreatWinType_enumGreatWinType_RobKong] = 1.0                  // 明杠冲（抢杠胡）

		greatWinPointMap[GreatWinType_enumGreatWinType_OpponentsRichi] = 1.0 // 别人起手报听，也得到了一辣子
		break
	}

	return rsc
}
