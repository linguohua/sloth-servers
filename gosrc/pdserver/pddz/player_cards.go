package pddz

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

// has2jokerOr42 是否有4个2或者二大王
func (pt *PlayerCards) has2jokerOr42() bool {
	jcount := 0
	tcount := 0
	for e := pt.hand.Front(); e != nil; e = e.Next() {
		t := e.Value.(*Card)
		rank := t.cardID / 4
		if rank == 13 {
			jcount++
		} else if rank == 0 {
			tcount++
		}
	}

	return jcount == 2 || tcount == 4
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

// bombCountOfDiscarded 打出的牌列表中炸弹的个数
func (pt *PlayerCards) bombOrRocketCountOfDiscarded() int {
	count := 0
	for e := pt.discarded.Front(); e != nil; e = e.Next() {
		cardHand := e.Value.(*CardHand)

		if cardHand.ht == CardHandType_Bomb ||
			cardHand.ht == CardHandType_Roket {
			count++
		}
	}

	return count
}

// addHandCard 增加一张手牌
func (pt *PlayerCards) addHandCard(t *Card) {
	// 斗地主最多20张牌
	if pt.cardCountInHand() >= 21 {
		log.Panic("Total cards must less than 21")
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

	cht := CardHandType(msgCardHand.GetCardHandType())
	// 火箭最大
	if cht == CardHandType_Roket {
		return false
	}

	switch cht {
	case CardHandType_Bomb:
		result = pt.hasBombGreatThan(msgCardHand)
		isBomb = true
		break
	case CardHandType_Single:
		result = pt.hasSingleGreatThan(msgCardHand)
		break
	case CardHandType_Pair:
		result = pt.hasPairGreatThan(msgCardHand)
		break
	case CardHandType_Triplet:
		result = pt.hasTripletGreatThan(msgCardHand)
		break
	case CardHandType_TripletPair:
		result = pt.hasTripletPairGreatThan(msgCardHand)
		break
	case CardHandType_TripletSingle:
		result = pt.hasTripletSingleGreatThan(msgCardHand)
		break
	case CardHandType_Flush:
		result = pt.hasFlushGreatThan(msgCardHand)
		break
	case CardHandType_Pair3X:
		result = pt.hasPair3XGreatThan(msgCardHand)
		break
	case CardHandType_Triplet2X:
		result = pt.hasTriplet2XGreatThan(msgCardHand)
		break
	case CardHandType_Triplet2X2Pair:
		result = pt.hasTriplet2X2PairGreatThan(msgCardHand)
		break
	case CardHandType_Triplet2X2Single:
		result = pt.hasTriplet2X2SingleGreatThan(msgCardHand)
		break
	case CardHandType_FourX2Pair:
		result = pt.hasFourX2PairGreatThan(msgCardHand)
		break
	case CardHandType_FourX2Single:
		result = pt.hasFourX2SingleGreatThan(msgCardHand)
		break
	default:
		log.Panicln("unkonwn CardHandType:", cht)
		return false
	}

	// 已经找到存在更大的牌组
	if result {
		return result
	}

	// 如果不是炸弹，则检查是否有炸弹
	if !isBomb {
		result = pt.hasBombOnHand()
	}

	// 自己手上有炸弹
	if result {
		return result
	}

	// 手上是否有火箭
	return pt.hasRoketOnHand()
}

// hasBombGreatThan 是否有更大的炸弹 四张同数值牌（如四个 7 ）
func (pt *PlayerCards) hasBombGreatThan(msgCardHand *pokerface.MsgCardHand) bool {
	if msgCardHand.GetCardHandType() != int32(CardHandType_Bomb) {
		log.Panicln("cardHand type not bomb")
	}

	pt.hand2Slots()
	rank := int(msgCardHand.Cards[0]) / 4
	priority := rank2Priority[rank]

	// 不考虑双王
	for priorityNew := priority + 1; priorityNew < 13; priorityNew++ {
		rankNew := priority2Rank[priorityNew]
		if pt.slots[rankNew] > 3 {
			return true
		}
	}

	return false
}

// hasSingleGreatThan 是否有更大的单张
func (pt *PlayerCards) hasSingleGreatThan(msgCardHand *pokerface.MsgCardHand) bool {
	if msgCardHand.GetCardHandType() != int32(CardHandType_Single) {
		log.Panicln("cardHand type not single")
	}

	pt.hand2Slots()
	cardID := int(msgCardHand.Cards[0])
	rank := cardID / 4
	if cardID == JOR {
		// 红小丑最大
		return false
	}

	if cardID == JOB {
		// 红小丑最大
		if pt.slots[rank] > 0 {
			return true
		}
	}

	// 比较priority
	priority := rank2Priority[rank]

	for priorityNew := priority + 1; priorityNew < 14; priorityNew++ {
		rankNew := priority2Rank[priorityNew]
		if pt.slots[rankNew] > 0 {
			return true
		}
	}

	return false
}

// hasPairGreatThan 是否有更大的对子
func (pt *PlayerCards) hasPairGreatThan(msgCardHand *pokerface.MsgCardHand) bool {
	if msgCardHand.GetCardHandType() != int32(CardHandType_Pair) {
		log.Panicln("cardHand type not pair")
	}

	pt.hand2Slots()
	rank := int(msgCardHand.Cards[0]) / 4
	priority := rank2Priority[rank]

	// 不考虑双王
	for priorityNew := priority + 1; priorityNew < 13; priorityNew++ {
		rankNew := priority2Rank[priorityNew]
		if pt.slots[rankNew] > 1 {
			return true
		}
	}

	return false
}

// hasTripletGreatThan 是否有更大的三张
func (pt *PlayerCards) hasTripletGreatThan(msgCardHand *pokerface.MsgCardHand) bool {
	if msgCardHand.GetCardHandType() != int32(CardHandType_Triplet) {
		log.Panicln("cardHand type not Triplet")
	}

	pt.hand2Slots()
	rank := int(msgCardHand.Cards[0]) / 4
	priority := rank2Priority[rank]

	// 不考虑双王
	for priorityNew := priority + 1; priorityNew < 13; priorityNew++ {
		rankNew := priority2Rank[priorityNew]
		if pt.slots[rankNew] > 2 {
			return true
		}
	}

	return false
}

// hasTripletPairGreatThan 是否具有更大的三带一对
func (pt *PlayerCards) hasTripletPairGreatThan(msgCardHand *pokerface.MsgCardHand) bool {
	if msgCardHand.GetCardHandType() != int32(CardHandType_TripletPair) {
		log.Panicln("cardHand type not TripletPair")
	}

	pt.hand2Slots()
	rank := int(msgCardHand.Cards[0]) / 4
	priority := rank2Priority[rank]

	// 不考虑双王
	for priorityNew := priority + 1; priorityNew < 13; priorityNew++ {
		testRank := priority2Rank[priorityNew]
		found := true
		if pt.slots[testRank] < 3 {
			found = false
		}

		if found {
			// 搜索一个对子
			leftPriority := priorityNew
			rightPriority := priorityNew

			pairCount := 0
			for testPriority := 0; testPriority < leftPriority && pairCount < 1; testPriority++ {
				testPair := priority2Rank[testPriority]
				if pt.slots[testPair] > 1 {
					pairCount++
				}
			}

			// 对子不考虑大小王
			for testPriority := rightPriority + 1; testPriority < 13 && pairCount < 1; testPriority++ {
				testPair := priority2Rank[testPriority]
				if pt.slots[testPair] > 1 {
					pairCount++
				}
			}

			if pairCount > 0 {
				return true
			}
		}
	}

	return false
}

// hasTripletSingleGreatThan 是否具有更大的三带一张
func (pt *PlayerCards) hasTripletSingleGreatThan(msgCardHand *pokerface.MsgCardHand) bool {
	if msgCardHand.GetCardHandType() != int32(CardHandType_TripletSingle) {
		log.Panicln("cardHand type not TripletSingle")
	}

	pt.hand2Slots()
	rank := int(msgCardHand.Cards[0]) / 4
	priority := rank2Priority[rank]

	// 不考虑双王
	for priorityNew := priority + 1; priorityNew < 13; priorityNew++ {
		testRank := priority2Rank[priorityNew]
		found := true
		if pt.slots[testRank] < 3 {
			found = false
		}

		if found {
			// 搜索一个
			leftPriority := priorityNew
			rightPriority := priorityNew

			pairCount := 0
			for testPriority := 0; testPriority < leftPriority && pairCount < 1; testPriority++ {
				testPair := priority2Rank[testPriority]
				if pt.slots[testPair] > 0 {
					pairCount++
				}
			}

			// 单张考虑大小王
			for testPriority := rightPriority + 1; testPriority < 14 && pairCount < 1; testPriority++ {
				testPair := priority2Rank[testPriority]
				if pt.slots[testPair] > 0 {
					pairCount++
				}
			}

			if pairCount > 0 {
				return true
			}
		}
	}

	return false
}

// hasFourX2SingleGreatThan 是否有更大的四带二
func (pt *PlayerCards) hasFourX2SingleGreatThan(msgCardHand *pokerface.MsgCardHand) bool {
	if msgCardHand.GetCardHandType() != int32(CardHandType_FourX2Single) {
		log.Panicln("cardHand type not FourX2Single")
	}

	pt.hand2Slots()
	rank := int(msgCardHand.Cards[0]) / 4
	priority := rank2Priority[rank]

	// 不考虑双王
	for priorityNew := priority + 1; priorityNew < 13; priorityNew++ {
		testRank := priority2Rank[priorityNew]
		found := true
		if pt.slots[testRank] < 4 {
			found = false
		}

		if found {
			// 搜索一个
			leftPriority := priorityNew
			rightPriority := priorityNew

			pairCount := 0
			for testPriority := 0; testPriority < leftPriority && pairCount < 2; testPriority++ {
				testPair := priority2Rank[testPriority]
				if pt.slots[testPair] > 0 {
					pairCount++
				}
			}

			// 单张考虑双王
			for testPriority := rightPriority + 1; testPriority < 14 && pairCount < 2; testPriority++ {
				testPair := priority2Rank[testPriority]
				if pt.slots[testPair] > 0 {
					pairCount++
				}
			}

			if pairCount > 1 {
				return true
			}
		}
	}

	return false
}

// hasFourX2PairGreatThan 是否有更大的四带二
func (pt *PlayerCards) hasFourX2PairGreatThan(msgCardHand *pokerface.MsgCardHand) bool {
	if msgCardHand.GetCardHandType() != int32(CardHandType_FourX2Pair) {
		log.Panicln("cardHand type not FourX2Pair")
	}

	pt.hand2Slots()
	rank := int(msgCardHand.Cards[0]) / 4
	priority := rank2Priority[rank]

	// 不考虑双王
	for priorityNew := priority + 1; priorityNew < 13; priorityNew++ {
		testRank := priority2Rank[priorityNew]
		found := true
		if pt.slots[testRank] < 4 {
			found = false
		}

		if found {
			// 搜索一个
			leftPriority := priorityNew
			rightPriority := priorityNew

			pairCount := 0
			for testPriority := 0; testPriority < leftPriority && pairCount < 2; testPriority++ {
				testPair := priority2Rank[testPriority]
				if pt.slots[testPair] > 1 {
					pairCount++
				}
			}

			// 对子不考虑双王
			for testPriority := rightPriority + 1; testPriority < 13 && pairCount < 2; testPriority++ {
				testPair := priority2Rank[testPriority]
				if pt.slots[testPair] > 1 {
					pairCount++
				}
			}

			if pairCount > 1 {
				return true
			}
		}
	}

	return false
}

// hasFlushGreatThan 是否有更大的顺子
func (pt *PlayerCards) hasFlushGreatThan(msgCardHand *pokerface.MsgCardHand) bool {
	if msgCardHand.GetCardHandType() != int32(CardHandType_Flush) {
		log.Panicln("cardHand type not flush")
	}

	// 顺子不包括2，以及大小王
	// 注意由于最小的顺子是34567，而现在需要寻找一个比他大的
	pt.hand2Slots()
	flushLen := len(msgCardHand.Cards)
	rank := int(msgCardHand.Cards[0]) / 4 // 最大的顺子牌

	for rankNew := 12; rankNew > rank; {
		testRank := rankNew
		found := true
		for i := 0; i < flushLen; i++ {
			if pt.slots[testRank-i] < 1 {
				rankNew = testRank - i - 1
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

// hasPair3XGreatThan 是否有更大的双顺
func (pt *PlayerCards) hasPair3XGreatThan(msgCardHand *pokerface.MsgCardHand) bool {
	if msgCardHand.GetCardHandType() != int32(CardHandType_Pair3X) {
		log.Panicln("cardHand type not Pair3X")
	}

	// 顺子不包括2，以及大小王
	// 注意由于最小的顺子是34567，而现在需要寻找一个比他大的
	pt.hand2Slots()
	flushLen := len(msgCardHand.Cards) / 2
	rank := int(msgCardHand.Cards[0]) / 4 // 最大的顺子牌

	for rankNew := 12; rankNew > rank; {
		testRank := rankNew
		found := true
		for i := 0; i < flushLen; i++ {
			if pt.slots[testRank-i] < 2 {
				rankNew = testRank - i - 1
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

// hasTriplet2XGreatThan 是否有更大的三顺
func (pt *PlayerCards) hasTriplet2XGreatThan(msgCardHand *pokerface.MsgCardHand) bool {
	if msgCardHand.GetCardHandType() != int32(CardHandType_Triplet2X) {
		log.Panicln("cardHand type not Triplet2X")
	}

	// 顺子不包括2，以及大小王
	// 注意由于最小的顺子是34567，而现在需要寻找一个比他大的
	pt.hand2Slots()
	flushLen := len(msgCardHand.Cards) / 3
	rank := int(msgCardHand.Cards[0]) / 4 // 最大的顺子牌

	for rankNew := 12; rankNew > rank; {
		testRank := rankNew
		found := true
		for i := 0; i < flushLen; i++ {
			if pt.slots[testRank-i] < 3 {
				rankNew = testRank - i - 1
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

// hasTriplet2X2PairGreatThan 是否有更大的三顺带对子
func (pt *PlayerCards) hasTriplet2X2PairGreatThan(msgCardHand *pokerface.MsgCardHand) bool {
	if msgCardHand.GetCardHandType() != int32(CardHandType_Triplet2X2Pair) {
		log.Panicln("cardHand type not Triplet2X2Pair")
	}

	// 顺子不包括2，以及大小王
	// 注意由于最小的顺子是34567，而现在需要寻找一个比他大的
	pt.hand2Slots()
	flushLen := len(msgCardHand.Cards) / 5 // 设长度为三顺长度为x，那么必须有x个对子，因此总牌数是5x
	rank := int(msgCardHand.Cards[0]) / 4  // 最大的顺子牌

	for rankNew := 12; rankNew > rank; {
		testRank := rankNew
		found := true
		for i := 0; i < flushLen; i++ {
			if pt.slots[testRank-i] < 3 {
				rankNew = testRank - i - 1
				found = false
				break
			}
		}

		if found {
			// 寻找N个对子
			// 搜索一个对子
			leftPriority := rank2Priority[testRank-flushLen+1]
			rightPriority := rank2Priority[testRank]

			pairCount := 0
			for testPriority := 0; testPriority < leftPriority && pairCount < flushLen; testPriority++ {
				testPair := priority2Rank[testPriority]
				if pt.slots[testPair] > 1 {
					pairCount++
				}
			}

			// 对子不考虑双王
			for testPriority := rightPriority + 1; testPriority < 13 && pairCount < flushLen; testPriority++ {
				testPair := priority2Rank[testPriority]
				if pt.slots[testPair] > 1 {
					pairCount++
				}
			}

			if pairCount >= flushLen {
				return true
			}

			// 需要移动rankNew
			rankNew--
		}
	}

	return false
}

// hasTriplet2X2SingleGreatThan 是否有更大的三顺带单张
func (pt *PlayerCards) hasTriplet2X2SingleGreatThan(msgCardHand *pokerface.MsgCardHand) bool {
	if msgCardHand.GetCardHandType() != int32(CardHandType_Triplet2X2Single) {
		log.Panicln("cardHand type not Triplet2X2Single")
	}

	// 顺子不包括2，以及大小王
	// 注意由于最小的顺子是34567，而现在需要寻找一个比他大的
	pt.hand2Slots()
	flushLen := len(msgCardHand.Cards) / 4 // 设长度为三顺长度为x，那么必须有x个单张，因此总牌数是4x
	rank := int(msgCardHand.Cards[0]) / 4  // 最大的顺子牌

	for rankNew := 12; rankNew > rank; {
		testRank := rankNew
		found := true
		for i := 0; i < flushLen; i++ {
			if pt.slots[testRank-i] < 3 {
				rankNew = testRank - i - 1
				found = false
				break
			}
		}

		if found {
			// 寻找N个对子
			// 搜索一个对子
			leftPriority := rank2Priority[testRank-flushLen+1]
			rightPriority := rank2Priority[testRank]

			pairCount := 0
			for testPriority := 0; testPriority < leftPriority && pairCount < flushLen; testPriority++ {
				testPair := priority2Rank[testPriority]
				if pt.slots[testPair] > 0 {
					pairCount++
				}
			}

			// 单个需要考虑大小王
			for testPriority := rightPriority + 1; testPriority < 14 && pairCount < flushLen; testPriority++ {
				testPair := priority2Rank[testPriority]
				if pt.slots[testPair] > 0 {
					pairCount++
				}
			}

			if pairCount >= flushLen {
				return true
			}

			// 需要移动rankNew
			rankNew--
		}
	}

	return false
}

// hasRoketOnHand 是否有火箭 即双王（大王和小王），最大的牌
func (pt *PlayerCards) hasRoketOnHand() bool {
	pt.hand2Slots()

	return pt.slots[JOB/4] > 1
}

// hasBombOnHand 是否有炸弹 四张同数值牌（如四个 7 ）
func (pt *PlayerCards) hasBombOnHand() bool {
	pt.hand2Slots()

	for rank := 0; rank < 13; rank++ {
		if pt.slots[rank] > 3 {
			return true
		}
	}

	return false
}

// convertMsgCardHand 把牌列表转换为MsgCardHand对象
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

	return msgCardHand
}

// isMsgCardHandGreatThan 比较两个牌的大小
func isMsgCardHandGreatThan(prevCardHand *CardHand, msgCardHand *pokerface.MsgCardHand) bool {

	cht := CardHandType(msgCardHand.GetCardHandType())
	pcht := prevCardHand.ht

	// 火箭最大
	if cht == (CardHandType_Roket) {
		return true
	}

	if pcht == (CardHandType_Roket) {
		return false
	}

	if cht == (CardHandType_Bomb) {
		// 炸弹仅次于火箭
		if prevCardHand.ht != CardHandType_Bomb {
			return true
		}

		// 都是炸弹，则比较牌大小
		rank1 := int(msgCardHand.Cards[0] / 4)
		rank2 := (prevCardHand.cards[0].cardID / 4)
		return rank2Priority[rank1] > rank2Priority[rank2]
	}

	// 前一个是炸弹而当前不是炸弹
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

	// 同类型，比较牌大小
	rank1 := int(msgCardHand.Cards[0] / 4)
	rank2 := (prevCardHand.cards[0].cardID / 4)

	// 大小王直接比较大小
	if rank1 == rank2 && rank1 == 13 {
		return msgCardHand.Cards[0] > int32(prevCardHand.cards[0].cardID)
	}

	return rank2Priority[rank1] > rank2Priority[rank2]
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

// hasFlushOnHand 手上是否有顺子
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
