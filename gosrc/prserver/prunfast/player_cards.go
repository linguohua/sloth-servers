package prunfast

import (
	"container/list"
	log "github.com/sirupsen/logrus"
	"pokerface"
)

type cardList = *list.List
type handList = *list.List

// PlayerCards 玩家牌表
type PlayerCards struct {
	// 玩家手上的牌列表
	hand cardList

	//  玩家打出的牌
	discarded handList

	slots []int

	handSlotCached bool

	host *PlayerHolder
}

// newPlayerCards 新建一个PlayerCards
func newPlayerCards(host *PlayerHolder) *PlayerCards {
	pt := PlayerCards{}
	pt.slots = make([]int, CARDMAX)

	// lists
	pt.hand = list.New()
	pt.discarded = list.New()

	pt.host = host
	return &pt
}

// isEmpty 没有任何牌
func (pt *PlayerCards) isEmpty() bool {
	return pt.hand.Len() < 1 &&
		pt.discarded.Len() < 1
}

// hand2Slots 把手牌放到slots上，以便胡牌计算
func (pt *PlayerCards) hand2Slots() {
	if pt.handSlotCached {
		return
	}

	for i := range pt.slots {
		pt.slots[i] = 0
	}

	for e := pt.hand.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Card)
		pt.slots[t.cardID/4]++
	}

	pt.handSlotCached = true
}

// removeCardInHandWithCard 从手牌移除摸个card
func (pt *PlayerCards) removeCardInHandWithCard(t *Card) {
	for e := pt.hand.Back(); e != nil; e = e.Prev() {
		if t == e.Value.(*Card) {
			pt.hand.Remove(e)
			pt.handSlotCached = false
			break
		}
	}
}

// cardCountInHand 手牌的数量
func (pt *PlayerCards) cardCountInHand() int {
	return pt.hand.Len()
}

// cardCountInHand 某种牌的在手牌中的个数
func (pt *PlayerCards) cardCountInHandOf(cardID int) int {
	count := 0
	for e := pt.hand.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Card)
		if t.cardID == cardID {
			count++
		}
	}
	return count
}

// cardInHandOf 查找手牌上某个ID的牌
func (pt *PlayerCards) cardInHandOf(cardID int) *Card {
	for e := pt.hand.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Card)
		if t.cardID == cardID {
			return t
		}
	}
	return nil
}

// clear 清空牌数据
func (pt *PlayerCards) clear() {
	pt.hand.Init()
	pt.discarded.Init()

	pt.handSlotCached = false
}

// addDiscardedCardHand 增加一张打出的牌到出牌列表
func (pt *PlayerCards) addDiscardedCardHand(h *CardHand) {

	pt.discarded.PushBack(h)
}

// addHandCard 增加一张手牌
func (pt *PlayerCards) addHandCard(t *Card) {
	// 大丰关张每个人仅允许16张牌
	if pt.cardCountInHand() >= 16 {
		log.Panic("Total cards must less than 16")
		return
	}

	pt.hand.PushBack(t)
	pt.handSlotCached = false
}

// hasCardInHand 手牌列表上是否有某张牌
func (pt *PlayerCards) hasCardInHand(cardID int) bool {
	for e := pt.hand.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Card)
		if t.cardID == cardID {
			return true
		}
	}
	return false
}

func (pt *PlayerCards) hasMsgCardHandOnHand(cardHand *pokerface.MsgCardHand) bool {
	pt.hand2Slots()

	for _, c := range cardHand.Cards {
		cardID := int(c)
		if pt.slots[cardID/4] < 1 {
			return false
		}
	}

	return true
}

func (pt *PlayerCards) removeMsgCardHandFromHand(cardHand *pokerface.MsgCardHand) *CardHand {
	cards := make([]*Card, 0, len(cardHand.Cards))
	for _, c := range cardHand.Cards {
		cardID := int(c)
		card := pt.removeCardInHand(cardID)
		if card == nil {
			log.Panicln("failed to remove Card inHand:", cardID)
		}

		cards = append(cards, card)
	}

	cardHandNew := newCardHand(CardHandType(cardHand.GetCardHandType()), cards)
	return cardHandNew
}

// removeCardInHand 从手牌列表中删除一张牌
func (pt *PlayerCards) removeCardInHand(cardID int) *Card {
	for e := pt.hand.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Card)
		if t.cardID == cardID {
			pt.hand.Remove(e)
			pt.handSlotCached = false
			return t
		}
	}
	return nil
}

// latestHandCard 手牌列表中的最后一张牌
func (pt *PlayerCards) latestHandCard() *Card {
	if pt.hand.Len() < 1 {
		log.Panic("latestHandCard, hand is empty")
		return nil
	}

	return pt.hand.Back().Value.(*Card)
}

// hand2IDList 手牌列表序列化到ID list，用于消息发送
func (pt *PlayerCards) hand2IDList() []int32 {
	int32List := make([]int32, pt.hand.Len())
	var i = 0
	for e := pt.hand.Front(); e != nil; e = e.Next() {
		int32List[i] = int32(e.Value.(*Card).cardID)
		i++
	}

	return int32List
}

// discardedCardHand2MsgCardHands 面子牌组列表序列化到列表，用于消息发送
func (pt *PlayerCards) discardedCardHand2MsgCardHands() []*pokerface.MsgCardHand {
	msgCardHands := make([]*pokerface.MsgCardHand, pt.discarded.Len())
	var i = 0

	for e := pt.discarded.Front(); e != nil; e = e.Next() {
		cardHand := e.Value.(*CardHand)
		msgCardHand := cardHand.cardHand2MsgCardHand()
		msgCardHands[i] = msgCardHand
		i++
	}

	return msgCardHands
}

// hasCardHandGreatThan 是否有牌组比cardHand大
func (pt *PlayerCards) hasCardHandGreatThan(cardHand *CardHand) bool {
	msgCardHand := cardHand.cardHand2MsgCardHand()
	return pt.hasMsgCardHandGreatThan(msgCardHand)
}

// hasMsgCardHandGreatThan 注意msgCardHand必须是经过服务器转换过其CardHandType，确保其合法性
func (pt *PlayerCards) hasMsgCardHandGreatThan(msgCardHand *pokerface.MsgCardHand) bool {
	var result bool
	isBomb := false
	switch CardHandType(msgCardHand.GetCardHandType()) {
	case CardHandType_Bomb:
		result = pt.hasBombGreatThan(msgCardHand)
		isBomb = true
		break
	case CardHandType_Flush:
		result = pt.hasFlushGreatThan(msgCardHand)
		break
	case CardHandType_Single:
		result = pt.hasSingleGreatThan(msgCardHand)
		break
	case CardHandType_Pair:
		result = pt.hasPairGreatThan(msgCardHand)
		break
	case CardHandType_Pair2X:
		result = pt.hasPair2XGreatThan(msgCardHand)
		break
	case CardHandType_Triplet:
		result = pt.hasTripletGreatThan(msgCardHand)
		break
	case CardHandType_Triplet2X:
		result = pt.hasTriplet2XGreatThan(msgCardHand)
		break
	case CardHandType_Triplet2X2Pair:
		result = pt.hasTriplet2X2PairGreatThan(msgCardHand)
		break
	case CardHandType_TripletPair:
		result = pt.hasTripletPairGreatThan(msgCardHand)
		break
	}

	if !result && !isBomb {
		return pt.hasBombOnHand()
	}

	return result
}

func (pt *PlayerCards) hasBombGreatThan(msgCardHand *pokerface.MsgCardHand) bool {
	if msgCardHand.GetCardHandType() != int32(CardHandType_Bomb) {
		log.Panicln("cardHand type not bomb")
	}

	pt.hand2Slots()
	bombCardSuitID := int(msgCardHand.Cards[0]) / 4
	for newBombRankID := bombCardSuitID + 1; newBombRankID < AH/4; newBombRankID++ {
		if pt.slots[newBombRankID] > 3 {
			return true
		}
	}

	if pt.slots[AH/4] > 2 {
		return true
	}

	return false
}

func (pt *PlayerCards) hasBombOnHand() bool {
	pt.hand2Slots()

	for newBombRankID := 0; newBombRankID < AH/4; newBombRankID++ {
		if pt.slots[newBombRankID] > 3 {
			return true
		}
	}

	if pt.slots[AH/4] > 2 {
		return true
	}

	if pt.slots[R3H/4] > 2 {
		count := 0
		for e := pt.hand.Front(); e != nil; e = e.Next() {
			t := e.Value.(*Card)
			if t.cardID/4 == R3H/4 && t.cardID != R3H {
				count++
			}
		}

		if count == 3 {
			return true
		}
	}

	return false
}

func (pt *PlayerCards) convertMsgCardHand(cards []int32) *pokerface.MsgCardHand {
	player := pt.host
	msgCardHand := agariConvertMsgCardHand(cards)
	if msgCardHand == nil {
		log.Println("failed to convertMsgCardHand, maybe invalid input cards, chair:", player.chairID)
		return nil
	}

	if !pt.hasMsgCardHandOnHand(msgCardHand) {
		log.Println("failed to convertMsgCardHand, hand has no cards:", player.chairID)
		return nil
	}

	if msgCardHand.GetCardHandType() == int32(CardHandType_Triplet) && ((msgCardHand.Cards[0] / 4) == (R3H / 4)) {
		// 3个三:不含有红桃3则认为是炸弹
		found := false
		for _, c := range msgCardHand.Cards {
			if c == R3H {
				found = true
				break
			}
		}
		if !found {
			var cardHandType32 = int32(CardHandType_Bomb)
			msgCardHand.CardHandType = &cardHandType32

			log.Println("three 3, and no R3H, convert triplet to bomb, chair:", player.chairID)
		}
	}

	return msgCardHand
}

func (pt *PlayerCards) hasFlushGreatThan(msgCardHand *pokerface.MsgCardHand) bool {
	if msgCardHand.GetCardHandType() != int32(CardHandType_Flush) {
		log.Panicln("cardHand type not flush")
	}

	// 注意由于最小的顺子是12345，而现在需要寻找一个比他大的，因此肯定不需要考虑ACE绕过来的情况
	pt.hand2Slots()
	flushLen := len(msgCardHand.Cards)
	bombCardSuitID := int(msgCardHand.Cards[0]) / 4 // 最大的顺子牌
	for newBombRankID := AH / 4; newBombRankID > bombCardSuitID; {
		testBombRankID := newBombRankID
		found := true
		for i := 0; i < flushLen; i++ {
			if pt.slots[testBombRankID-i] < 1 {
				newBombRankID = testBombRankID - i - 1
				found = false
				break
			}
		}

		if found {
			return true
		}
	}

	return false
}

func (pt *PlayerCards) hasFlushOnHand() bool {
	// 从3到ACE寻找一个可用顺子
	pt.hand2Slots()
	for start := 1; start < RankEnd/4; start++ {
		flushBegin := start
		count := 0
		for ; flushBegin < RankEnd/4; flushBegin++ {
			if pt.slots[flushBegin] < 1 {
				break
			}
			count++
		}

		if count >= 5 {
			return true
		}
	}

	return false
}

func (pt *PlayerCards) allSingleCardWith(lastDiscarded int) bool {
	pt.hand2Slots()

	// 先把最后打出的那张牌补上
	pt.slots[lastDiscarded/4]++
	if pt.hasFlushOnHand() {
		// 把cache失效
		pt.handSlotCached = false
		return false
	}

	// 把cache失效
	pt.handSlotCached = false
	// 对子，3张
	for _, c := range pt.slots {
		if c > 1 {
			return false
		}
	}

	return true
}

func (pt *PlayerCards) hasSingleGreatThan(msgCardHand *pokerface.MsgCardHand) bool {
	if msgCardHand.GetCardHandType() != int32(CardHandType_Single) {
		log.Panicln("cardHand type not single")
	}

	pt.hand2Slots()
	bombCardSuitID := int(msgCardHand.Cards[0]) / 4
	if bombCardSuitID == 0 {
		// 2 是最大的
		return false
	}

	if pt.slots[0] > 0 {
		// 2 是最大的
		return true
	}

	for newBombRankID := bombCardSuitID + 1; newBombRankID < RankEnd/4; newBombRankID++ {
		if pt.slots[newBombRankID] > 0 {
			return true
		}
	}

	return false
}

func (pt *PlayerCards) hasPair2XGreatThan(msgCardHand *pokerface.MsgCardHand) bool {
	if msgCardHand.GetCardHandType() != int32(CardHandType_Pair2X) {
		log.Panicln("cardHand type not Pair2X")
	}

	pt.hand2Slots()
	flushLen := len(msgCardHand.Cards)
	bombCardSuitID := int(msgCardHand.Cards[0]) / 4 // 最大的顺子牌
	for newBombRankID := AH / 4; newBombRankID > bombCardSuitID; {
		testBombRankID := newBombRankID
		found := true
		for i := 0; i < flushLen/2; i++ {
			if pt.slots[testBombRankID-i] < 2 {
				newBombRankID = testBombRankID - i - 1
				found = false
				break
			}
		}

		if found {
			return true
		}
	}

	return false
}

func (pt *PlayerCards) hasPairGreatThan(msgCardHand *pokerface.MsgCardHand) bool {
	if msgCardHand.GetCardHandType() != int32(CardHandType_Pair) {
		log.Panicln("cardHand type not pair")
	}

	pt.hand2Slots()
	bombCardSuitID := int(msgCardHand.Cards[0]) / 4
	for newBombRankID := bombCardSuitID + 1; newBombRankID < RankEnd/4; newBombRankID++ {
		if pt.slots[newBombRankID] > 1 {
			return true
		}
	}

	return false
}

func (pt *PlayerCards) hasTripletGreatThan(msgCardHand *pokerface.MsgCardHand) bool {
	if msgCardHand.GetCardHandType() != int32(CardHandType_Triplet) {
		log.Panicln("cardHand type not Triplet")
	}

	pt.hand2Slots()
	bombCardSuitID := int(msgCardHand.Cards[0]) / 4
	for newBombRankID := bombCardSuitID + 1; newBombRankID < RankEnd/4; newBombRankID++ {
		if pt.slots[newBombRankID] > 2 {
			return true
		}
	}

	return false
}

func (pt *PlayerCards) hasTriplet2XGreatThan(msgCardHand *pokerface.MsgCardHand) bool {
	if msgCardHand.GetCardHandType() != int32(CardHandType_Triplet2X) {
		log.Panicln("cardHand type not Triplet2X")
	}

	pt.hand2Slots()
	flushLen := len(msgCardHand.Cards)
	bombCardSuitID := int(msgCardHand.Cards[0]) / 4 // 最大的顺子牌
	for newBombRankID := AH / 4; newBombRankID > bombCardSuitID; {
		testBombRankID := newBombRankID
		found := true
		for i := 0; i < flushLen/3; i++ {
			if pt.slots[testBombRankID-i] < 3 {
				newBombRankID = testBombRankID - i - 1
				found = false
				break
			}
		}

		if found {
			return true
		}
	}

	return false
}

func (pt *PlayerCards) hasTriplet2X2PairGreatThan(msgCardHand *pokerface.MsgCardHand) bool {
	if msgCardHand.GetCardHandType() != int32(CardHandType_Triplet2X2Pair) {
		log.Panicln("cardHand type not Triplet2X2Pair")
	}

	pairLength := 4
	if len(msgCardHand.Cards) > 10 {
		pairLength = 6
	}

	pt.hand2Slots()
	flushLen := len(msgCardHand.Cards) - pairLength
	pairCountNeed := flushLen / 3
	bombCardSuitID := int(msgCardHand.Cards[0]) / 4 // 最大的顺子牌
	for newBombRankID := AH / 4; newBombRankID > bombCardSuitID; {
		testBombRankID := newBombRankID
		found := true
		for i := 0; i < flushLen/3; i++ {
			if pt.slots[testBombRankID-i] < 3 {
				newBombRankID = testBombRankID - i - 1
				found = false
				break
			}
		}

		if found {
			// 搜索N个对子
			left := newBombRankID + 1 - flushLen/3
			right := newBombRankID

			pairCount := 0
			for testPair := 0; testPair < left; testPair++ {
				if pt.slots[testPair] > 1 {
					pairCount++
				}
			}

			for testPair := right + 1; testPair < RankEnd/4; testPair++ {
				if pt.slots[testPair] > 1 {
					pairCount++
				}
			}

			if pairCount >= pairCountNeed {
				return true
			}

			newBombRankID--
		}
	}

	return false
}

func (pt *PlayerCards) hasTripletPairGreatThan(msgCardHand *pokerface.MsgCardHand) bool {
	if msgCardHand.GetCardHandType() != int32(CardHandType_TripletPair) {
		log.Panicln("cardHand type not TripletPair")
	}

	pt.hand2Slots()
	flushLen := len(msgCardHand.Cards) - 2
	bombCardSuitID := int(msgCardHand.Cards[0]) / 4 // 最大的顺子牌
	for newBombRankID := AH / 4; newBombRankID > bombCardSuitID; {
		testBombRankID := newBombRankID
		found := true
		for i := 0; i < flushLen/3; i++ {
			if pt.slots[testBombRankID-i] < 3 {
				newBombRankID = testBombRankID - i - 1
				found = false
				break
			}
		}

		if found {
			// 搜索一个对子
			left := newBombRankID + 1 - flushLen/3
			right := newBombRankID

			pairCount := 0
			for testPair := 0; testPair < left; testPair++ {
				if pt.slots[testPair] > 1 {
					pairCount++
				}
			}

			for testPair := right + 1; testPair < RankEnd/4; testPair++ {
				if pt.slots[testPair] > 1 {
					pairCount++
				}
			}

			if pairCount > 0 {
				return true
			}

			newBombRankID--
		}
	}

	return false
}

// isMsgCardHandGreatThan 比较两个牌的大小
func isMsgCardHandGreatThan(prevCardHand *CardHand, msgCardHand *pokerface.MsgCardHand) bool {
	if msgCardHand.GetCardHandType() == int32(CardHandType_Bomb) {
		if prevCardHand.ht != CardHandType_Bomb {
			return true
		}

		return int(msgCardHand.Cards[0]/4) > (prevCardHand.cards[0].cardID / 4)
	}

	// 前一个是炸弹
	if prevCardHand.ht == CardHandType_Bomb {
		return false
	}

	if int32(prevCardHand.ht) != msgCardHand.GetCardHandType() {
		// 必须类型匹配
		return false
	}

	if len(prevCardHand.cards) != len(msgCardHand.Cards) {
		// 张数匹配
		return false
	}

	// 单张时2最大
	if prevCardHand.ht == CardHandType_Single {
		if prevCardHand.cards[0].cardID/4 == 0 {
			return false
		}

		if msgCardHand.Cards[0]/4 == 0 {
			return true
		}
	}

	return int(msgCardHand.Cards[0]/4) > prevCardHand.cards[0].cardID/4
}

func intArrayEquals(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
