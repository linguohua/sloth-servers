package dfmahjong

import (
	"log"
	"mahjong"
)

// ScoreContext 得分计算
type ScoreContext struct {
	winType int

	greatWinType              int
	fGreatWinPoints           float32
	fTrimGreatWinPoints       float32
	fContinuousBankerMultiple float32

	miniWinType         int
	specialScore        int     // 墩子分整数即可，不需要浮点数
	fMiniWinBasicScore  float32 // 小胡基础底分，例如0.3之类
	fMiniWinFlowerScore float32 // 小胡的花牌的花分
	fMiniMultiple       float32 // 小胡倍率
	//fMiniWinTrimScore   float32 // 裁剪后的小胡分数
	fMiniWinUnTrimScore float32 // 未裁剪的小胡分数

	//baseWinScore  int // 基本赢输分数
	// totalWinScore int // 总分整数即可，不需要浮点数

	//continuousBankerExtra int
	isContinuousBanker bool
	//fakers                []*PlayerHolder
	// fakeWinScore int
	isRichiPay1P bool

	orderPlayerSctxs []*PlayerScoreContext
}

// PlayerScoreContext 与每一个玩家的分数关系
type PlayerScoreContext struct {
	target        *PlayerHolder
	totalWinScore int // 总分整数即可，不需要浮点数
	fakeWinScore  int
	hasClear      bool // 是否两清了
}

func (sc *ScoreContext) isWin() bool {
	return sc.winType == int(mahjong.HandOverType_enumHandOverType_Win_Chuck) || sc.winType == int(mahjong.HandOverType_enumHandOverType_Win_SelfDrawn)
}

func (sc *ScoreContext) initPlayerScoreContext(orderPlayers []*PlayerHolder) {
	sc.orderPlayerSctxs = make([]*PlayerScoreContext, len(orderPlayers))
	for i, p := range orderPlayers {
		pc := &PlayerScoreContext{}
		pc.target = p

		sc.orderPlayerSctxs[i] = pc
	}
}

func (sc *ScoreContext) getPayTarget(player *PlayerHolder) *PlayerScoreContext {
	for _, pc := range sc.orderPlayerSctxs {
		if pc.target == player {
			return pc
		}
	}

	log.Panicln("ScoreContext can't find target for:", player.userID())
	return nil
}

func (sc *ScoreContext) calcTotalWinScore() int {
	sum := 0
	for _, pc := range sc.orderPlayerSctxs {
		if pc != nil {
			sum += pc.totalWinScore
		}
	}

	return sum
}

func (sc *ScoreContext) calcTotalFakeScore() int {
	sum := 0
	for _, pc := range sc.orderPlayerSctxs {
		if pc != nil {
			sum += pc.fakeWinScore
		}
	}

	return sum
}
