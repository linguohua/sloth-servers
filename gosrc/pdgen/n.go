package main

import (
	"log"
	"math"
	"pddz"
	"sort"
)

var (
	keyTags = make(map[int64]int)
)

func makeTagValue(haiCount int, ct pddz.CardHandType, flushCount int) int {
	// 第一个字节是类型
	// 第二个字节是牌张数
	// 第三个字节判断是否需要顺子检查
	v := int(ct)
	v |= (haiCount << 8)

	if flushCount > 0 {
		v |= (flushCount << 16)
	}

	return v
}

// 炸弹：四张同数值牌（如四个 7 ）
func genBombTag() {
	keyTags[4] = makeTagValue(4, pddz.CardHandType_Bomb, 0)
}

// 单牌：单个牌（如红桃 5 ）
func genSingleTag() {
	keyTags[1] = makeTagValue(1, pddz.CardHandType_Single, 0)
}

// 对牌：数值相同的两张牌（如梅花 4+ 方块 4 ）
func genPairTag() {
	keyTags[2] = makeTagValue(2, pddz.CardHandType_Pair, 0)
}

// 三张牌：数值相同的三张牌（如三个 J ）
func genTripletTag() {
	keyTags[3] = makeTagValue(3, pddz.CardHandType_Triplet, 0)
}

// 三带一：数值相同的三张牌 + 一张单牌或一对牌。例如： 333+6 或 444+99
func genTripletSingleTag() {
	keyTags[31] = makeTagValue(4, pddz.CardHandType_TripletSingle, 0)
}

// 三带一：数值相同的三张牌 + 一张单牌或一对牌。例如： 333+6 或 444+99
func genTripletPairTag() {
	keyTags[32] = makeTagValue(5, pddz.CardHandType_TripletPair, 0)
}

// 单顺：五张或更多的连续单牌（如： 45678 或 78910JQK ）不包括 2 点和双王。
func genFlushXTag() {
	// 5张
	tag := int64(11111)
	keyTags[tag] = makeTagValue(5, pddz.CardHandType_Flush, 5)

	// 6张以及以上；最多12张
	for i := 0; i < 7; i++ {
		tag = tag*10 + 1
		keyTags[tag] = makeTagValue(i+6, pddz.CardHandType_Flush, 6+i)
	}
}

// 双顺：三对或更多的连续对牌（如： 334455 、 7788991010JJ ）不包括 2 点和双王
func genPair3xTag() {
	// 3对
	tag := int64(222)
	keyTags[tag] = makeTagValue(5, pddz.CardHandType_Pair3X, 3)

	// 4对以及以上；最多10对
	for i := 0; i < 7; i++ {
		tag = tag*10 + 2
		keyTags[tag] = makeTagValue((i+4)*2, pddz.CardHandType_Pair3X, 4+i)
	}
}

// 三顺：二个或更多的连续三张牌（如： 333444 ， 555666777888 ）不包括 2 点和双王
func genTriplet2xTag() {
	// 2个三张
	tag := int64(33)
	keyTags[tag] = makeTagValue(6, pddz.CardHandType_Triplet2X, 2)

	// 3个三张及以上；最多7个三张
	for i := 0; i < 5; i++ {
		tag = tag*10 + 3
		keyTags[tag] = makeTagValue((i+3)*3, pddz.CardHandType_Triplet2X, 3+i)
	}
}

// 飞机带翅膀
func genTriplet2XXTag() {
	genTriplet2XTagPair()
	genTriplet2XTagSingle()
}

// 飞机带翅膀，加对子
func genTriplet2XTagPair() {
	// 最多5个飞机
	for i := 2; i < 6; i++ {
		tag1 := 3
		for j := 1; j < i; j++ {
			tag1 = tag1*10 + 3
		}

		// 带对子
		tag2 := 2
		for j := 1; j < i; j++ {
			tag2 = tag2*10 + 2
		}

		tag := int64(tag1*int(math.Pow(10, float64(i))) + tag2)
		keyTags[tag] = makeTagValue((i)*5, pddz.CardHandType_Triplet2X2Pair, i)
	}
}

// 飞机带翅膀，加单张
func genTriplet2XTagSingle() {
	// 最多5个飞机
	for i := 2; i < 6; i++ {
		genTriplet2XTagSingleDivide(i, i, 3, 0, 0)
	}
}

func genTriplet2XTagSingleDivide(rawFlyLength int, remainFlyLength int, divUnit int, divTimes int, tag2 int) {
	if remainFlyLength == 0 {
		// 最小可以分割单元了
		tag1 := 3
		for j := 1; j < rawFlyLength; j++ {
			tag1 = tag1*10 + 3
		}

		tag := int64(tag1*int(math.Pow(10, float64(divTimes))) + tag2)
		keyTags[tag] = makeTagValue((rawFlyLength)*4, pddz.CardHandType_Triplet2X2Single, rawFlyLength)
		return
	}

	myDivUnit := divUnit
	if myDivUnit > remainFlyLength {
		myDivUnit = remainFlyLength
	}

	for du := myDivUnit; du >= 1; du-- {
		tag2N := tag2*10 + du
		flyLengthN := remainFlyLength - du
		divTimesN := divTimes + 1
		genTriplet2XTagSingleDivide(rawFlyLength, flyLengthN, du, divTimesN, tag2N)
	}
}

// 四带二：四张牌+两手牌。（注意：四带二不是炸弹）如： 5555 + 3 + 8 或 4444 + 55 + 77
func genFourX2PairTag() {
	tag := int64(422)
	keyTags[tag] = makeTagValue(8, pddz.CardHandType_FourX2Pair, 0)
}

// 四带二：四张牌+两手牌。（注意：四带二不是炸弹）如： 5555 + 3 + 8 或 4444 + 55 + 77
func genFourX2SingleTag() {
	tag := int64(411)
	keyTags[tag] = makeTagValue(6, pddz.CardHandType_FourX2Single, 0)

	tag = int64(42)
	keyTags[tag] = makeTagValue(6, pddz.CardHandType_FourX2Single, 0)
}

func genAllTag() {
	genBombTag()
	genSingleTag()
	genPairTag()
	genTripletTag()
	genTripletSingleTag()
	genTripletPairTag()
	genFlushXTag()
	genPair3xTag()

	genTriplet2xTag()
	genTriplet2XXTag()

	genFourX2PairTag()
	genFourX2SingleTag()
}

func calcKeyTag(haiRank []int) int64 {
	// 最多14种牌
	slots := make([]int, 14)

	for _, r := range haiRank {
		slots[r]++
	}

	for _, s := range slots {
		if s > 4 {
			log.Panicln("slots elem great than 3:", s)
		}
	}

	sort.Ints(slots)
	tag := int64(0)
	for i := len(slots) - 1; i >= 0; i-- {
		if slots[i] == 0 {
			break
		}

		tag = tag*10 + int64(slots[i])
	}

	return tag
}

func assertCardHandType(hai []int, cht pddz.CardHandType) {
	tag := calcKeyTag(hai)
	v, ok := keyTags[tag]
	if !ok {
		log.Panicln("tag not valid:", tag)
	}

	// log.Printf("tag:%d, v:%d\n", tag, v)
	cht2 := (v & 0x0f)
	if cht2 != int(cht) {
		log.Panicln("CardHandType not exist:", cht)
	}

	// log.Println("equal:", cht)
}
