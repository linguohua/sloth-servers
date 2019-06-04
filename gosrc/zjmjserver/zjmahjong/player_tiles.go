package zjmahjong

import (
	"container/list"
	"mahjong"

	log "github.com/sirupsen/logrus"
)

type tileList = *list.List
type meldList = *list.List

// PlayerTiles 玩家牌表
type PlayerTiles struct {
	// 玩家手上的牌列表
	hand tileList
	// 花牌:梅兰竹菊、春夏秋冬一共8种且8张（每种一张），
	// 再加上被当花牌的风牌(同种4张),加上中发白12张，一共12种且24张牌
	flowers tileList
	//  玩家打出的牌
	discarded tileList
	// 玩家吃、碰、杠的面子列表
	// 落地牌，胡牌计算时，不允许重新组合
	fixedMelds meldList

	slots []int

	handSlotCached bool

	flowerSlotCached bool

	host *PlayerHolder
}

// newPlayerTiles 新建一个PlayerTiles
func newPlayerTiles(host *PlayerHolder) *PlayerTiles {
	pt := PlayerTiles{}
	pt.slots = make([]int, TILEMAX)

	// lists
	pt.hand = list.New()
	pt.flowers = list.New()
	pt.discarded = list.New()
	pt.fixedMelds = list.New()
	pt.host = host
	return &pt
}

// meldTileCount 面子牌组的所有牌张数
// 注意这里是逻辑张数，例如杠牌的逻辑张数是3，物理张数是4
func (pt *PlayerTiles) meldTileCount() int {
	count := 0
	for e := pt.fixedMelds.Front(); e != nil; e = e.Next() {
		meld := (e.Value).(*Meld)
		count += meld.logicTileCount()
	}

	return count
}

// nonFlowerTileCount 非花牌张数
func (pt *PlayerTiles) nonFlowerTileCount() int {
	return pt.hand.Len() + pt.meldTileCount()
}

// flowerTileCount 花牌张数
func (pt *PlayerTiles) flowerTileCount() int {
	return pt.flowers.Len()
}

// isEmpty 没有任何牌
func (pt *PlayerTiles) isEmpty() bool {
	return pt.hand.Len() < 1 &&
		pt.flowers.Len() < 1 &&
		pt.discarded.Len() < 1 &&
		pt.fixedMelds.Len() < 1
}

// agariTileCount 可以参与胡牌运算的牌张数
func (pt *PlayerTiles) agariTileCount() int {
	return pt.nonFlowerTileCount()
}

// winAble 判断手牌是否可以胡牌
func (pt *PlayerTiles) winAble() bool {
	if pt.agariTileCount() != 14 {
		return false
	}

	pt.hand2Slots()
	return isWinable(pt.slots)
}

// meldCount 面子牌组个数
func (pt *PlayerTiles) meldCount() int {
	return pt.fixedMelds.Len()
}

// hand2Slots 把手牌放到slots上，以便胡牌计算
func (pt *PlayerTiles) hand2Slots() {
	if pt.handSlotCached {
		return
	}

	for i := range pt.slots {
		pt.slots[i] = 0
	}

	for e := pt.hand.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Tile)
		pt.slots[t.tileID]++
	}

	pt.handSlotCached = true
	pt.flowerSlotCached = false
}

// flower2Slots 把花牌放到slots上，以便计算
func (pt *PlayerTiles) flower2Slots() {
	if pt.flowerSlotCached {
		return
	}

	for i := range pt.slots {
		pt.slots[i] = 0
	}

	for e := pt.flowers.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Tile)
		pt.slots[t.tileID]++
	}

	pt.handSlotCached = false
	pt.flowerSlotCached = true
}

// winAbleWith 判断新增一张牌情况下是否可以胡牌
func (pt *PlayerTiles) winAbleWith(latestDiscardTile *Tile) bool {
	pt.hand2Slots()

	var found = false
	var i = latestDiscardTile.tileID
	pt.slots[i]++
	if isWinable(pt.slots) {
		// 可以听
		//tiles.AddReadyHandTile(i);
		found = true
	}
	pt.slots[i]--

	return found
}

// removeTileInDiscarded 从打出的牌列表中删除，倒序搜索
func (pt *PlayerTiles) removeTileInDiscarded(t *Tile) {
	for e := pt.discarded.Back(); e != nil; e = e.Prev() {
		if t == e.Value.(*Tile) {
			pt.discarded.Remove(e)
			break
		}
	}
}

// addMeld 增加一个面子牌组
func (pt *PlayerTiles) addMeld(m *Meld) {
	if pt.nonFlowerTileCount() >= 12 {
		log.Panic("AddMeld tiles count should less than 12")
		return
	}

	pt.fixedMelds.PushBack(m)
}

// tileCountInHand 手牌的数量
func (pt *PlayerTiles) tileCountInHand() int {
	return pt.hand.Len()
}

// tileCountInHand 某种牌的在手牌中的个数
func (pt *PlayerTiles) tileCountInHandOf(tileID int) int {
	count := 0
	for e := pt.hand.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Tile)
		if t.tileID == tileID {
			count++
		}
	}
	return count
}

// concealedKongAble 是否可以暗杠
func (pt *PlayerTiles) concealedKongAble() bool {
	pt.hand2Slots()
	for _, v := range pt.slots {
		if v > 3 {
			return true
		}
	}
	return false
}

// isAllTriplet 胡牌是否由4个刻子和一对构成
// 注意杠牌也是刻子
func (pt *PlayerTiles) isAllTriplet() bool {
	pt.hand2Slots()
	pairs := 0
	for _, v := range pt.slots {
		if v == 2 {
			pairs++
		}
	}

	// 必须只有一个对子
	if pairs != 1 {
		return false
	}

	for _, v := range pt.slots {
		if v != 2 && v != 3 && v != 0 {
			// 如果出现1个牌，或者4个牌，表明是一定有顺子
			return false
		}
	}

	// 落地牌组必须全部是碰牌组
	for e := pt.fixedMelds.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Meld)
		if m.mt != mahjong.MeldType_enumMeldTypeTriplet {
			return false
		}
	}

	return true
}

// isThirteenOrphans 十三幺
func (pt *PlayerTiles) isThirteenOrphans() bool {
	pt.hand2Slots()
	return isAgariThirteenOrphans(pt.slots)
}

// exposedMeldCount 落地牌组中明牌示人的牌组数量
// 也即是除了暗杠之外的meld数量
func (pt *PlayerTiles) exposedMeldCount() int {
	count := 0
	for e := pt.fixedMelds.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Meld)
		if m.mt != mahjong.MeldType_enumMeldTypeConcealedKong {
			count++
		}
	}
	return count
}

// concealedKongCount 暗杠牌组数量
func (pt *PlayerTiles) concealedKongCount() int {
	count := 0
	for e := pt.fixedMelds.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Meld)
		if m.mt == mahjong.MeldType_enumMeldTypeConcealedKong {
			count++
		}
	}
	return count
}

// concealedKongIDList 暗杠牌组数量
func (pt *PlayerTiles) concealedKongIDList() []int32 {
	idlist := make([]int32, 0, 4)

	for e := pt.fixedMelds.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Meld)
		if m.mt == mahjong.MeldType_enumMeldTypeConcealedKong {
			idlist = append(idlist, int32(m.t1.tileID))
		}
	}
	return idlist
}

// exposedKongCount 明杠或者加杠的数量
func (pt *PlayerTiles) exposedKongCount() int {
	count := 0
	for e := pt.fixedMelds.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Meld)
		if m.isExposedKong() {
			count++
		}
	}
	return count
}

// tripletMeldCount 刻子牌组数量
func (pt *PlayerTiles) tripletMeldCount() int {
	count := 0
	for e := pt.fixedMelds.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Meld)
		if m.isTriplet() {
			count++
		}
	}
	return count
}

// sequenceMeldCount 顺子牌组数量
func (pt *PlayerTiles) sequenceMeldCount() int {
	count := 0
	for e := pt.fixedMelds.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Meld)
		if m.isSequence() {
			count++
		}
	}
	return count
}

// kongMeldCount 杠牌组（明，暗，续）数量
func (pt *PlayerTiles) kongMeldCount() int {
	count := 0
	for e := pt.fixedMelds.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Meld)
		if m.isKong() {
			count++
		}
	}
	return count
}

// suitTypeCount 所有索子牌种类
func (pt *PlayerTiles) suitTypeCount() int {
	set := make(map[int]bool)
	// 先检查落地牌组中的索子
	for e := pt.fixedMelds.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Meld)
		if m.isSequence() {
			set[m.t1.suitType()] = true
			set[m.t2.suitType()] = true
			set[m.t3.suitType()] = true
		} else if m.t1.isSuit() {
			set[m.t1.suitType()] = true
		}
	}
	// 再检查手牌上的索子
	for e := pt.hand.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Tile)
		if t.isSuit() {
			set[t.suitType()] = true
		}
	}

	return len(set)
}

// honorTypeCount 字牌种类，字牌是包括风牌和箭牌，但是不包括花牌
func (pt *PlayerTiles) honorTypeCount() int {
	set := make(map[int]bool)
	// 先检查落地牌组中的字牌
	for e := pt.fixedMelds.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Meld)
		if m.t1.isHonor() {
			set[m.t1.tileID] = true
		}
	}
	// 再检查手牌上的字牌
	for e := pt.hand.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Tile)
		if t.isHonor() {
			set[t.tileID] = true
		}
	}

	return len(set)
}

// dragonTypeCount 箭牌种类
func (pt *PlayerTiles) dragonTypeCount() int {
	set := make(map[int]bool)
	// 先检查落地牌组中的箭牌
	for e := pt.fixedMelds.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Meld)
		if m.t1.isDragon() {
			set[m.t1.tileID] = true
		}
	}
	// 再检查手牌上的箭牌
	for e := pt.hand.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Tile)
		if t.isDragon() {
			set[t.tileID] = true
		}
	}

	return len(set)
}

// windTypeCount 风牌种类
func (pt *PlayerTiles) windTypeCount() int {
	set := make(map[int]bool)
	// 先检查落地牌组中的风牌
	for e := pt.fixedMelds.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Meld)
		if m.t1.isWind() {
			set[m.t1.tileID] = true
		}
	}
	// 再检查手牌上的风牌
	for e := pt.hand.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Tile)
		if t.isWind() {
			set[t.tileID] = true
		}
	}

	return len(set)
}

// isSecondClearFrontAble 是否算作小门清
// 需求修正：如果没有吃椪杠他人（注意是他人，因此暗杠不算）行为，且有花分，叫小门清
func (pt *PlayerTiles) isSecondClearFrontAble() bool {
	meldCount := pt.exposedMeldCount()
	if meldCount > 0 {
		// 有吃椪杠他人行为，因此不是小门清
		return false
	}

	// 如果有暗杠，或者有花牌，都可以产生花分，因此算作小门清
	if pt.flowerTileCount() > 0 || pt.concealedKongCount() > 0 {
		return true
	}

	return false
}

// allFlowerScoreCount 计算小胡花分
// 一个碰算1个花，3个一样的牌在手里算2个花（必须胡出刻子），花墩子算8花（还要算花墩子），
// 明杠非花算3花（还要额外算墩子），暗杠非花算4花（还要额外算墩子）
// 为牌局前预设的分数，针对花牌，1个花X分
// a.胡牌为小胡，手里有3只或4只一样的牌，额外要算2个花。
// b.胡牌为小胡，且有明杠，额外要算3个花。
// c.胡牌为小胡，且有暗杠，额外要算4个花。
// d.胡牌为小胡，且有门前有花墩子（4只一样的花牌），算8个花。
// e.门前有除开花墩子外的花，有几个花牌算几个花
func (pt *PlayerTiles) allFlowerScoreCount(selfDraw bool) int {
	count := 0
	var pongCount = pt.tripletMeldCount()
	count += pongCount // 一个碰算1个花
	log.Printf("FlowerX, pong %d *1=> %d\n", pongCount, pongCount*1)

	var ekongCount = pt.exposedKongCount()
	count += ekongCount * 3 // 明杠（包括加杠），额外要算3个花

	log.Printf("FlowerX, ekong %d *3=> %d\n", ekongCount, ekongCount*3)

	var ckongCount = pt.concealedKongCount()
	count += ckongCount * 4 // 暗杠，额外要算4个花
	log.Printf("FlowerX, ckong %d *4=> %d\n", ckongCount, ckongCount*4)

	pt.hand2Slots()

	// --需求修正：必须胡牌牌型中，这几个一样牌真的一幅刻子才可以
	// --例如，123万33万345万678条这种牌胡了，就不能将333万为2个花
	// -- 必须是123条  333万  345万 678条 99饼  才能将333万算2个花
	// 需求修正：如果最后一张牌是形成刻子，且改牌是吃炮而来的，则该刻子只能获得1分而不是两分
	var sameInHandCount, lastTileKotCount = winAbleKotCount(pt.slots, pt.latestHandTile().tileID)
	var lastTileCountInHand = pt.tileCountInHandOf(pt.latestHandTile().tileID)

	log.Printf("FlowerX, same inHand %d *2=> %d\n", sameInHandCount, sameInHandCount*2)
	count += (sameInHandCount * 2)
	if lastTileKotCount > 0 && !selfDraw && lastTileCountInHand < 4 {
		log.Println("last tile is in kot, but chuck from other, can earn only 1 score")
		count--
	}

	pt.flower2Slots()
	var f4Count = pt.quadFlowerCount()
	var fCount = pt.flowerTileCount() - 4*f4Count
	count += 8 * f4Count
	count += fCount

	log.Printf("FlowerX, dunzi %d *8=> %d\n", f4Count, f4Count*8)
	log.Printf("FlowerX, non-dunzi flower %d *1=> %d\n", fCount, fCount)

	return count
}

// quadFlowerCount 4张一样的花牌组，有几个
// 注意这样的花牌组如果有，最多有6个，只可能是风牌当花牌的情形以及中发白
// 因为其他8种花牌，都是每种只得1张牌
// --需求修改：如果春夏秋冬都有，算一个花墩子；如果梅兰竹菊都有，算一个花墩子
func (pt *PlayerTiles) quadFlowerCount() int {
	count := 0
	pt.flower2Slots()
	for _, v := range pt.slots {
		if v == 4 {
			// 花墩子（4只一样的花牌），算8个花。
			count++
		}
	}

	// PLUM          = 34 // 梅
	// ORCHID        = 35 // 兰
	// BAMBOO        = 36 // 竹
	// CHRYSANTHEMUM = 37 // 菊
	if pt.slots[PLUM] > 0 && pt.slots[ORCHID] > 0 &&
		pt.slots[BAMBOO] > 0 && pt.slots[CHRYSANTHEMUM] > 0 {
		count++
	}

	// SPRING        = 38 // 春
	// SUMMER        = 39 // 夏
	// AUTUMN        = 40 // 秋
	// WINTER        = 41 // 冬
	if pt.slots[SPRING] > 0 && pt.slots[SUMMER] > 0 &&
		pt.slots[AUTUMN] > 0 && pt.slots[WINTER] > 0 {
		count++
	}

	if count > 6 {
		log.Panic("quad flower should not great than 6")
	}
	return count
}

// readyHandAble 是否可以听牌
func (pt *PlayerTiles) readyHandAble() bool {
	if pt.agariTileCount() != 13 {
		return false
	}

	pt.hand2Slots()
	found := false

	for i := MAN; i < FlowerBegin; i++ {
		// 跳过当做花牌的风牌
		if i == pt.host.room.pseudoFlowerTileID {
			continue
		}

		pt.slots[i]++
		if isWinable(pt.slots) {
			found = true
			pt.slots[i]--
			break
		}
		pt.slots[i]--
	}

	return found
}

// clear 清空牌数据
func (pt *PlayerTiles) clear() {
	pt.hand.Init()
	pt.flowers.Init()
	pt.fixedMelds.Init()
	pt.discarded.Init()

	pt.handSlotCached = false
	pt.flowerSlotCached = false
}

// addFlowerTile 增加一个花牌
func (pt *PlayerTiles) addFlowerTile(t *Tile) {
	// 非花牌必须小于14个，否则不可能得到抽牌机会
	if pt.nonFlowerTileCount() >= 14 {
		log.Panic("Total tiles must less than 14")
		return
	}

	if !t.isFlower() {
		log.Panic("add flower tile must be flower")
		return
	}

	pt.flowers.PushBack(t)
	pt.flowerSlotCached = false
}

// addDiscardedTile 增加一张打出的牌到出牌列表
func (pt *PlayerTiles) addDiscardedTile(t *Tile) {
	if pt.nonFlowerTileCount() >= 14 {
		log.Panic("Total tiles must less than 14")
		return
	}

	pt.discarded.PushBack(t)
}

// addHandTile 增加一张手牌
func (pt *PlayerTiles) addHandTile(t *Tile) {
	// 非花牌必须小于14个，否则不可能得到抽牌机会
	if pt.nonFlowerTileCount() >= 14 {
		log.Panic("Total tiles must less than 14")
		return
	}

	if t.isFlower() {
		log.Panic("add hand tile must not be flower")
		return
	}

	pt.hand.PushBack(t)
	pt.handSlotCached = false
}

// temporaryHandAdd 临时增加一张手牌
func (pt *PlayerTiles) temporaryHandAdd(t *Tile) {
	// 非花牌必须小于14个，否则不可能得到抽牌机会
	if pt.nonFlowerTileCount() >= 14 {
		log.Panic("Total tiles must less than 14")
		return
	}

	if t.isFlower() {
		log.Panic("add hand tile must not be flower")
		return
	}

	pt.hand.PushBack(t)
	pt.handSlotCached = false
}

// removeTileInHand 从手牌列表中删除一张牌
func (pt *PlayerTiles) temporaryHandRemove(tile *Tile) {
	e := pt.hand.Back()
	if e == nil {
		log.Panicln("temporaryHandRemove, hand is empty")
		return
	}

	lastTile := e.Value.(*Tile)
	if lastTile != tile {
		log.Panicln("temporaryHandRemove, last tile not temporay add one")
		return
	}

	pt.hand.Remove(e)
	pt.handSlotCached = false
}

// addPesudoFlowerTile 增加一张伪花牌（即是风牌当花牌）
func (pt *PlayerTiles) addPesudoFlowerTile(t *Tile) {
	// 非花牌必须小于14个，否则不可能得到抽牌机会
	if pt.nonFlowerTileCount() >= 14 {
		log.Panic("Total tiles must less than 14")
		return
	}
	pt.flowers.PushBack(t)
	pt.flowerSlotCached = false
}

// hasTileInHand 手牌列表上是否有某张牌
func (pt *PlayerTiles) hasTileInHand(tileID int) bool {
	for e := pt.hand.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Tile)
		if t.tileID == tileID {
			return true
		}
	}
	return false
}

// removeTileInHand 从手牌列表中删除一张牌
func (pt *PlayerTiles) removeTileInHand(tileID int) *Tile {
	for e := pt.hand.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Tile)
		if t.tileID == tileID {
			pt.hand.Remove(e)
			pt.handSlotCached = false
			return t
		}
	}
	return nil
}

// chowAbleWith 是否可以吃牌
func (pt *PlayerTiles) chowAbleWith(t *Tile) bool {
	// 只有索子，筒子，万子才可以构成顺子
	if !t.isSuit() {
		return false
	}

	// 如果最后打出去的牌和当前想吃的牌相同，则不能吃
	latestDiscardTileOfSelf := pt.host.hStatis.latestDiscardedTileLocked
	// 需求修正，只剩下4张牌的时候，不需要考虑锁定
	if pt.hand.Len() > 4 {
		if latestDiscardTileOfSelf.tileID == t.tileID {
			return false
		}
	}

	chowIDList := pt.chowAble2IdList(t.tileID)
	return len(chowIDList) > 0
}

// chowAble2IdList 获得可以吃牌序列
// 例如上家打出3万，当前玩家可以（1万，2万，3万）、（2万，3万，4万）
// （3万，4万，5万）几种情况来吃上家打出的3万
// 返回结果数组，每一个元素表示顺子第一个牌，有多少个元素就有多少个顺子，也即是多少种吃法
func (pt *PlayerTiles) chowAble2IdList(tileID int) []int {
	pt.hand2Slots()
	pt.slots[tileID]++

	discardedTileTileID := tileID
	var i, j, k, myLatestDiscardedTileID int
	myLatestDiscardTile := pt.host.hStatis.latestDiscardedTileLocked
	myLatestDiscardedTileID = myLatestDiscardTile.tileID

	var tt = int(tileID / 9)
	var lbound = tt * 9
	var ubound = lbound + 9

	isMyLatestDicardedSuit := myLatestDiscardedTileID >= lbound && myLatestDiscardedTileID < ubound

	var idlist = make([]int, 14)
	var cnt = 0
	for i = tileID - 2; i <= tileID; i++ {
		if i < lbound {
			continue
		}

		j = i + 1
		k = i + 2

		if k >= ubound {
			continue
		}

		if pt.slots[i] > 0 && pt.slots[j] > 0 && pt.slots[k] > 0 {
			// log.Printf("match, %d,%d,%d\n", i, j, k)
			// 如果吃这个牌和手上的某两个牌构成顺子，而之前打出的牌和这两个牌也能构成顺子（可是玩家却选择了不要顺子而打出去）
			// 那么现在也不能吃，例如，玩家本来手上有4,5,6，现在他打出6万，然后上家打出3万，他是不能吃的。因为他本来的4,5万和他刚打出去的6万就是一个顺子
			// 如果牌数少于等于4张，则随便chow
			if pt.hand.Len() > 4 && isMyLatestDicardedSuit {
				if ((i == discardedTileTileID) && (k+1 == myLatestDiscardedTileID)) ||
					((k == discardedTileTileID) && (i-1 == myLatestDiscardedTileID)) {
					// 相邻情形不能胡
					log.Printf("chowAbleWith, can't chow, neighbor: lastDiscardedSelf:%d, chow:%d\n", myLatestDiscardedTileID, tileID)
				} else {
					idlist[cnt] = i
					cnt++
				}
			} else {
				idlist[cnt] = i
				cnt++
			}
		}
	}

	pt.slots[tileID]--
	return idlist[:cnt]
}

// pongAbleWith 是否可以碰牌
func (pt *PlayerTiles) pongAbleWith(t *Tile) bool {
	// 最后一次打出去的牌不能跟本次要碰的牌相同
	latestDiscardTileOfSelf := pt.host.hStatis.latestDiscardedTileLocked
	// 需求修正：如果手上只剩下4张牌，那么不考虑这个锁定
	if pt.hand.Len() > 4 && latestDiscardTileOfSelf.tileID == t.tileID {
		return false
	}

	pt.hand2Slots()
	return pt.slots[t.tileID] >= 2
}

// exposedKongAbleWith 是否可以明杠
func (pt *PlayerTiles) exposedKongAbleWith(t *Tile) bool {
	pt.hand2Slots()
	return pt.slots[t.tileID] == 3
}

// exposedKongAbleWith 是否可以明杠
func (pt *PlayerTiles) concealedKongAbleWith(tileID int) bool {
	pt.hand2Slots()
	return pt.slots[tileID] == 4
}

// calc7Pair 计算7对和豪华七对
// 七对：7对不一样的牌组成的胡牌
// 豪华大七对：首先满足七对，且：有4个同种牌，且胡的那只刚好是4只相同中的1只
// 修正：必须是14张手牌，也即是不能有任何牌组
func (pt *PlayerTiles) calc7Pair() GreatWinType {

	// 不能有任何落地牌组
	if pt.meldCount() > 0 {
		return GreatWinType_None
	}

	pt.hand2Slots()

	// 先检查手牌，必须是对子，或者4只
	for _, v := range pt.slots {
		if v != 2 && v != 0 && v != 4 {
			// 出现1或者3，不是7对
			return GreatWinType_None
		}
	}

	// 不管是自摸，还是胡别人的牌，一定是放在最后一张
	var lastTile = pt.hand.Back().Value.(*Tile)
	// 检查手牌中的有没有4只同样牌，如果有而且等于最后一张牌
	// 那么就是豪华大七对
	for i, v := range pt.slots {
		if v != 4 {
			continue
		}
		if i == lastTile.tileID {
			return GreatWinType_GreatSevenPair
		}
	}

	return GreatWinType_SevenPair
}

// pongCountFrom 某玩家贡献了多少个刻子牌组
func (pt *PlayerTiles) pongCountFrom(p *PlayerHolder) int {
	if pt.host == p {
		return 0
	}

	count := 0
	userID := p.userID()
	for e := pt.fixedMelds.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Meld)
		if !t.isTriplet() {
			continue
		}

		if t.t1.drawBy == userID ||
			t.t2.drawBy == userID ||
			t.t3.drawBy == userID {
			count++
		}
	}
	return count
}

// chowPongExposedKongCount 吃椪牌组个数
func (pt *PlayerTiles) chowPongExposedKongCount() int {
	count := 0
	for e := pt.fixedMelds.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Meld)
		if t.isTriplet() || t.isSequence() || t.isExposedKong() {
			count++
		}
	}
	return count
}

func (pt *PlayerTiles) triplet2KongRollback(kongTile *Tile) {
	for e := pt.fixedMelds.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Meld)
		if t.mt == mahjong.MeldType_enumMeldTypeTriplet2Kong {
			if t.t4 == kongTile {
				t.triplet2KongRollback()
				break
			}
		}
	}
}

// kongCountFrom 某玩家贡献了多少个加杠、明杠牌组
func (pt *PlayerTiles) kongCountFrom(p *PlayerHolder) int {
	if pt.host == p {
		return 0
	}

	count := 0
	userID := p.userID()
	for e := pt.fixedMelds.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Meld)
		if !t.isExposedKong() {
			continue
		}

		if t.t1.drawBy == userID ||
			t.t2.drawBy == userID ||
			t.t3.drawBy == userID ||
			t.t4.drawBy == userID {
			count++
		}
	}
	return count
}

// chowCountFrom 上家贡献了多少个顺子牌组
func (pt *PlayerTiles) chowCountFrom(p *PlayerHolder) int {
	if pt.host == p {
		return 0
	}

	count := 0
	userID := p.userID()
	for e := pt.fixedMelds.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Meld)
		if !t.isSequence() {
			continue
		}

		if t.t1.drawBy == userID ||
			t.t2.drawBy == userID ||
			t.t3.drawBy == userID {
			count++
		}
	}
	return count
}

// tripletMeldWith 获得某牌的碰牌牌组
func (pt *PlayerTiles) tripletMeldWith(tileID int) *Meld {
	for e := pt.fixedMelds.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Meld)
		if m.isTriplet() && m.t1.tileID == tileID {
			return m
		}
	}
	return nil
}

// triplet2KongAbleWith 是否因某牌而可以加杠
func (pt *PlayerTiles) triplet2KongAbleWith(tileID int) bool {
	return pt.hasPongOf(tileID)
}

// hasPongOf 是否有某牌的碰牌牌组
func (pt *PlayerTiles) hasPongOf(tileID int) bool {
	return pt.tripletMeldWith(tileID) != nil
}

// chowPongLockedIDList 和吃牌组构造顺子的牌的ID
func (pt *PlayerTiles) chowPongLockedIDList(chowPongTile *Tile) []int {
	var ids = make([]int, 0, 3)
	if pt.fixedMelds.Len() < 1 {
		return ids
	}

	lastMeld := pt.fixedMelds.Back().Value.(*Meld)

	if lastMeld.isTriplet() {
		ids = append(ids, chowPongTile.tileID)
		return ids
	}

	if lastMeld.isSequence() {
		tileID := chowPongTile.tileID
		log.Printf("chowPongLockedIDList, tileID:%d, t1:%d,t2:%d,t3:%d\n", tileID, lastMeld.t1.tileID, lastMeld.t2.tileID, lastMeld.t3.tileID)
		var tt = int(tileID / 9)
		var lbound = tt * 9
		var ubound = lbound + 9

		if lastMeld.t1 == chowPongTile {
			ids = append(ids, tileID)
			bound := lastMeld.t3.tileID + 1
			if bound < ubound {
				ids = append(ids, bound)
			}
		} else if lastMeld.t2 == chowPongTile {
			ids = append(ids, tileID)
		} else if lastMeld.t3 == chowPongTile {
			ids = append(ids, tileID)
			bound := lastMeld.t1.tileID - 1
			if bound >= lbound {
				ids = append(ids, bound)
			}
		}
	}

	return ids
}

// readyHandTilesWhenThrow 获得可以听的牌列表
// 返回结果是int数组，内容布局：[可以听得牌的id][可以听得牌的剩余数量]...
func (pt *PlayerTiles) readyHandTilesWhenThrow(tileID int, tm *TileMgr) []int32 {
	if tileID != TILEMAX {
		if pt.agariTileCount() != 14 {
			log.Panic("ReadyHandTilesWhenThrow, AgariTileCount != 14")
			return nil
		}
	}

	// slots
	pt.hand2Slots()

	if tileID != TILEMAX {
		if pt.slots[tileID] <= 0 {
			log.Panicf("ReadyHandTilesWhenThrow, no tid %d to throw\n", tileID)
			return nil
		}
	}

	if tileID != TILEMAX {
		// 修改slot
		pt.slots[tileID]--
	}

	remainSlots := tm.tileRemainInHandOrWall()
	// 假如可以听34个牌
	var idlist = make([]int32, 34*2)
	cnt := 0
	for i := MAN; i < FlowerBegin; i++ {
		if i == pt.host.room.pseudoFlowerTileID {
			continue
		}

		var remain = remainSlots[i] - pt.slots[i]
		if i == tileID {
			remain--
		}

		if remain < 1 {
			continue
		}

		pt.slots[i]++
		if isWinable(pt.slots) {
			// 可以听
			idlist[cnt] = int32(i)
			cnt++
			idlist[cnt] = int32(remain)
			cnt++
		}

		pt.slots[i]--
	}

	if tileID != TILEMAX {
		// 还原slot
		pt.slots[tileID]++
	}

	return idlist[:cnt]
}

// latestHandTile 手牌列表中的最后一张牌
func (pt *PlayerTiles) latestHandTile() *Tile {
	if pt.hand.Len() < 1 {
		log.Panic("latestHandTile, hand is empty")
		return nil
	}

	return pt.hand.Back().Value.(*Tile)
}

// discard2IDList 打出的牌列表序列化到ID list，用于消息发送
func (pt *PlayerTiles) discard2IDList() []int32 {
	int32List := make([]int32, pt.discarded.Len())
	var i = 0
	for e := pt.discarded.Front(); e != nil; e = e.Next() {
		int32List[i] = int32(e.Value.(*Tile).tileID)
		i++
	}

	return int32List
}

// flower2IDList 花牌列表序列化到ID list，用于消息发送
func (pt *PlayerTiles) flower2IDList() []int32 {
	int32List := make([]int32, pt.flowers.Len())
	var i = 0
	for e := pt.flowers.Front(); e != nil; e = e.Next() {
		int32List[i] = int32(e.Value.(*Tile).tileID)
		i++
	}

	return int32List
}

// hand2IDList 手牌列表序列化到ID list，用于消息发送
func (pt *PlayerTiles) hand2IDList() []int32 {
	int32List := make([]int32, pt.hand.Len())
	var i = 0
	for e := pt.hand.Front(); e != nil; e = e.Next() {
		int32List[i] = int32(e.Value.(*Tile).tileID)
		i++
	}

	return int32List
}

// melds2MsgMeldTileList 面子牌组列表序列化到列表，用于消息发送
func (pt *PlayerTiles) melds2MsgMeldTileList(mark bool) []*mahjong.MsgMeldTile {
	msgMelds := make([]*mahjong.MsgMeldTile, pt.fixedMelds.Len())
	var i = 0

	for e := pt.fixedMelds.Front(); e != nil; e = e.Next() {
		meld := e.Value.(*Meld)
		msgMeld := pt.meld2MsgMeldTile(meld, mark)
		msgMelds[i] = msgMeld
		i++
	}

	return msgMelds
}

// meld2MsgMeldTile server内部的meld转换为MsgMeldTile，以便发送给客户端
func (pt *PlayerTiles) meld2MsgMeldTile(meld *Meld, mark bool) *mahjong.MsgMeldTile {
	msgMeld := &mahjong.MsgMeldTile{}
	var meldType = int32(meld.mt)
	msgMeld.MeldType = &meldType
	var tile1 = int32(meld.t1.tileID)
	msgMeld.Tile1 = &tile1

	contributor, tile := (pt.getContributorChairID(meld))

	var contributor32 = int32(contributor)
	msgMeld.Contributor = &contributor32

	if meld.mt == mahjong.MeldType_enumMeldTypeConcealedKong && mark {
		tile1 = TILEMAX
	}

	if meld.mt == mahjong.MeldType_enumMeldTypeSequence {
		var chowTile32 = int32(tile.tileID)
		msgMeld.ChowTile = &chowTile32
	}
	return msgMeld
}

// getContributorChairID 获得contributor的座位ID
func (pt *PlayerTiles) getContributorChairID(meld *Meld) (int, *Tile) {
	myslef := pt.host
	var tile *Tile
	var othersUserID string
	var tiles = []*Tile{meld.t1, meld.t2, meld.t3, meld.t4}

	for _, t := range tiles {
		if t.drawBy != "" && t.drawBy != myslef.userID() {
			othersUserID = t.drawBy
			tile = t
			break
		}
	}

	if othersUserID != "" {
		other := myslef.room.getPlayerByUserID(othersUserID)
		return other.chairID, tile
	}

	return myslef.chairID, tile
}

// concealedKongAble2IDList 暗杠到列表，用于发送消息给客户端
func (pt *PlayerTiles) concealedKongAble2IDList() []int {
	pt.hand2Slots()
	idlist := make([]int, 0, 3)
	for tid, i := range pt.slots {
		if i > 3 {
			idlist = append(idlist, tid)
		}
	}

	return idlist
}

// pongMelds 获得所有碰牌的牌组
func (pt *PlayerTiles) pongMelds() []*Meld {
	melds := make([]*Meld, 0, pt.meldCount())
	for e := pt.fixedMelds.Front(); e != nil; e = e.Next() {
		m := e.Value.(*Meld)
		if m.isTriplet() {
			melds = append(melds, m)
		}
	}
	return melds
}

// triplet2KongAble2IDList 可以加杠的列表
func (pt *PlayerTiles) triplet2KongAble2IDList() []int {
	idList := make([]int, 0, pt.meldCount())
	pongMelds := pt.pongMelds()
	for _, pm := range pongMelds {
		if pt.tileCountInHandOf(pm.t1.tileID) > 0 {
			idList = append(idList, pm.t1.tileID)
		}
	}

	return idList
}
