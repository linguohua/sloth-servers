package pddz

// HStatis 一手牌统计信息
type HStatis struct {
	// locked
	//latestDiscardedCardLocked *Card
	//latestChowPongCardLocked  *Card
	// kongLockedMap    map[int]bool
	// accumulativeWins []*AccumulativeWinContext

	//pongAbleCardLocked int
	// isWinAbleLocked bool

	// isRichi bool
	// 是否加注
	isCallDouble bool

	// 是否不叫地主
	isSkipCall bool

	lastExpectedActions int

	lastExpectedType int

	// 用于记录玩家本手牌已经进行了多少个动作
	// 目前主要用于monkey房间测试时，发送动作
	// 记录提示给客户端
	actionCounter int
}

// reset 重置
func (hs *HStatis) reset() {
	hs.isCallDouble = false
	hs.lastExpectedActions = 0
	// hs.isRichi = false
	hs.actionCounter = 0
	hs.lastExpectedType = 0
	hs.isSkipCall = false

	// hs.accumulativeWins = make([]*AccumulativeWinContext, 0, 4)
	// hs.kongLockedMap = make(map[int]bool)
	hs.resetLocked()
}

// resetLocked 重置locked
func (hs *HStatis) resetLocked() {
	// hs.latestChowPongCardLocked = InvalidCard
	// hs.latestDiscardedCardLocked = InvalidCard
	// hs.pongAbleCardLocked = TILEMAX
	// hs.isWinAbleLocked = false
}

func newHStatis() *HStatis {
	h := &HStatis{}
	h.reset()

	return h
}
