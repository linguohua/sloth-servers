package zjmahjong

// HStatis 一手牌统计信息
type HStatis struct {
	// locked
	latestDiscardedTileLocked *Tile // 用于锁定自己刚打出的牌不能立即吃进来、碰进来
	latestChowPongTileLocked  *Tile // 用于锁定刚吃，碰进来的牌，不能立即打出同样的牌
	isWinAbleLocked           bool  // 用于锁定过手胡，例如玩家可胡3万，6万，他摸到3万可以胡可是他选择过并打出3万，其他人打6万他也不可以吃铳胡

	isRichi bool

	// 用于记录玩家本手牌已经进行了多少个动作
	// 目前主要用于monkey房间测试时，发送动作
	// 记录提示给客户端
	actionCounter int
}

// reset 重置
func (hs *HStatis) reset() {
	hs.isRichi = false
	hs.actionCounter = 0
	hs.resetLocked()
}

// resetLocked 重置locked
func (hs *HStatis) resetLocked() {
	hs.latestChowPongTileLocked = InvalidTile
	hs.latestDiscardedTileLocked = InvalidTile
	hs.isWinAbleLocked = false
}

func newHStatis() *HStatis {
	h := &HStatis{}
	h.reset()

	return h
}
