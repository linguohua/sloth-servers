package dfmahjong

import "mahjong"

// Meld 面子牌组
// 3个一组：顺子，刻子
// 4个一组：杠
// 对于杠牌，运算时依然当做是一组3张面子牌，例如玩家杠牌时（此时牌张数已经由13张加上新杠的牌变成14张），
// 玩家必须从牌墙上补牌，然后再出一张（也就是此时玩家出完牌后还有14张牌），而不是像吃、碰那样直接出牌，
// 但麻将规定只有胡牌时才允许手上一共14张牌（花牌除外），
// 因此，明杠或者暗杠的面子牌其实是当做3张牌。
//
// 加杠：
// 加杠是指已经碰了牌后来摸到一张同样的牌，由于这张牌往往不可以顺利和其他牌组成面子牌组，因此需要处理掉这张牌，
// 通过“加杠”把这张牌加到之前碰的面子牌组上即可消除这张牌，然后玩家继续摸牌和出牌
type Meld struct {
	mt mahjong.MeldType
	t1 *Tile
	t2 *Tile
	t3 *Tile
	t4 *Tile
}

// logicTileCount 所有面子牌组都当做3个牌
// 顺子、刻子固然是3个
// 但明杠、暗杠都当是明刻子和暗刻子，因此当做3个牌
func (m *Meld) logicTileCount() int {
	return 3
}

// physicTileCount 真正的牌数量
func (m *Meld) physicTileCount() int {
	switch m.mt {
	case mahjong.MeldType_enumMeldTypeConcealedKong:
		return 4
	case mahjong.MeldType_enumMeldTypeExposedKong:
		return 4
	case mahjong.MeldType_enumMeldTypeTriplet2Kong:
		return 4
	case mahjong.MeldType_enumMeldTypeTriplet:
		return 3
	case mahjong.MeldType_enumMeldTypeSequence:
		return 3
	default:
		return 3
	}
}

// isExposedKong 是否明杠或者加杠
func (m *Meld) isExposedKong() bool {
	return m.mt == mahjong.MeldType_enumMeldTypeExposedKong || m.mt == mahjong.MeldType_enumMeldTypeTriplet2Kong
}

// isConcealedKong 是否暗杠
func (m *Meld) isConcealedKong() bool {
	return m.mt == mahjong.MeldType_enumMeldTypeConcealedKong
}

// isTriplet 是否刻子
func (m *Meld) isTriplet() bool {
	return m.mt == mahjong.MeldType_enumMeldTypeTriplet
}

// isSequence 是否顺子
func (m *Meld) isSequence() bool {
	return m.mt == mahjong.MeldType_enumMeldTypeSequence
}

// isKong 是否杠牌（包括：明杠，暗杠，加杠）
func (m *Meld) isKong() bool {
	return m.physicTileCount() == 4
}

// triplet2Kong 把牌组转换为加杠牌组
func (m *Meld) triplet2Kong(tile *Tile) {
	m.mt = mahjong.MeldType_enumMeldTypeTriplet2Kong
	m.t4 = tile
}

// triplet2KongRollback 加杠被抢后回滚到碰牌
func (m *Meld) triplet2KongRollback() {
	m.mt = mahjong.MeldType_enumMeldTypeTriplet
	m.t4 = nil
}
