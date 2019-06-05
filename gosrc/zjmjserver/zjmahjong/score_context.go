package zjmahjong

import (
	"mahjong"

	log "github.com/sirupsen/logrus"
)

// ScoreContext 得分计算
type ScoreContext struct {
	winType int

	greatWinType   int
	greatWinPoints int
	horseCount     int

	isContinuousBanker bool
	orderPlayerSctxs   []*PlayerScoreContext
}

// PlayerScoreContext 与每一个玩家的分数关系
type PlayerScoreContext struct {
	target        *PlayerHolder
	totalWinScore int  // 总分整数即可，不需要浮点数
	hasClear      bool // 是否两清了
}

func (sc *ScoreContext) isWin() bool {
	return sc.winType == int(mahjong.HandOverType_enumHandOverType_Win_Chuck) ||
		sc.winType == int(mahjong.HandOverType_enumHandOverType_Win_SelfDrawn)
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
