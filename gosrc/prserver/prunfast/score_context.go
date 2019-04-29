package prunfast

import (
	log "github.com/sirupsen/logrus"
)

// ScoreContext 得分计算
type ScoreContext struct {
	room    *Room
	winType int

	// 大胡
	greatWinType   int
	greatWinPoints int

	// 包庄
	isPayForAll bool
	// 包庄、爬坡倍数
	markupMultiple int

	// 大小胡计数
	greatWinCount int
	miniWinCount  int

	orderPlayerSctxs []*PlayerScoreContext
}

// PlayerScoreContext 与每一个玩家的分数关系
type PlayerScoreContext struct {
	target        *PlayerHolder
	finalWinScore int // 总分整数即可，不需要浮点数
	//fakeWinScore  int
	cardWinScore int // 胡牌得分
	// kongMultiple int // 杠分倍数

	hasFinalPay bool // 是否两清了
	hasCalc     bool
}

func (sc *ScoreContext) isWin() bool {
	return sc.winType == int(HandOverType_enumHandOverType_Win_Chuck) ||
		sc.winType == int(HandOverType_enumHandOverType_Win_SelfDrawn) || sc.winType == int(HandOverType_enumHandOverType_Win_RobKong)
}

func (sc *ScoreContext) initPlayerScoreContext(orderPlayers []*PlayerHolder, room *Room) {
	sc.room = room
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
			sum += pc.finalWinScore
		}
	}

	return sum
}
