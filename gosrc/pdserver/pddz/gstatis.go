package pddz

import (
	"math"
)

// GStatis 一局牌的统计信息
type GStatis struct {
	roundScore          int // 四舍五入后的分数
	miniWinCounter      int // 新疆麻将修改为小胡次数
	winSelfDrawnCounter int
	greatWinCounter     int // 新疆麻将修改为大胡次数
	kongerCounter       int // 包庄次数
	winRobKongCounter   int
	isContinuousBanker  bool
}

func (gs *GStatis) reset() {
	gs.roundScore = 0
}

func newGStatis() *GStatis {
	g := &GStatis{}
	return g
}

// roundFloat32 不采用四舍五入，直接向上取整
func roundFloat32(f float32) int {
	// v := int(math.Floor(float64(f)))
	// remain := f - float32(v)
	// if remain >= 0.5 {
	// 	return v + 1
	// }

	v := int(math.Ceil(float64(f)))
	return v
}
