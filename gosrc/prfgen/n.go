package main

import (
	"log"
	"math"
	"prserver/prunfast"
	"sort"
)

var (
	keyTags = make(map[int64]int)
)

func makeTagValue(haiCount int, ct prunfast.CardHandType, flushCount int) int {
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
	keyTags[4] = makeTagValue(4, prunfast.CardHandType_Bomb, 0)
}

// 单牌：单个牌（如红桃 5 ）
func genSingleTag() {
	keyTags[1] = makeTagValue(1, prunfast.CardHandType_Single, 0)
}

// 对牌：数值相同的两张牌（如梅花 4+ 方块 4 ）
func genPairTag() {
	keyTags[2] = makeTagValue(2, prunfast.CardHandType_Pair, 0)
}

// 三张牌：数值相同的三张牌（如三个 J ）
// 如果三张是ACE，则认为是炸弹
func genTripletTag() {
	keyTags[3] = makeTagValue(3, prunfast.CardHandType_Triplet, 0)
}

// 三带二：数值相同的三张牌 + 一张单牌或一对牌。例如： 333+6 或 444+99
func genTripletPairTag() {
	keyTags[32] = makeTagValue(5, prunfast.CardHandType_TripletPair, 0)
}

// 单顺：五张或更多的连续单牌（如： 45678 或 78910JQK ）不包括 2 点和双王。
func genFlushXTag() {
	// 5张
	tag := int64(11111)
	keyTags[tag] = makeTagValue(5, prunfast.CardHandType_Flush, 5)

	// 6张以及以上；最多12张
	// 顺子不能包含2，因此只要3到ACE，一共12种
	// [关张每个人发16张牌，一共48张牌（没有小丑因此52张牌，去掉一张ace，去掉三张2, 52-4=48）]
	for i := 0; i < 7; i++ {
		tag = tag*10 + 1
		keyTags[tag] = makeTagValue(i+6, prunfast.CardHandType_Flush, 6+i)
	}
}

// 连对：二个或更多的连续对牌（如： 3344 ， 55667788 ）不包括双王
func genPair2xTag() {
	// 2个对子
	tag := int64(22)
	keyTags[tag] = makeTagValue(6, prunfast.CardHandType_Pair2X, 2)

	// 2个对子及以上；最多8个对子，因为每人最多16张牌
	// [关张每个人发16张牌，一共48张牌（没有小丑因此52张牌，去掉一张ace，去掉三张2, 52-4=48）]
	for i := 0; i < 6; i++ {
		tag = tag*10 + 2
		keyTags[tag] = makeTagValue((i+2)*2, prunfast.CardHandType_Pair2X, 2+i)
	}
}

// 三顺：二个或更多的连续三张牌（如： 333444 ， 555666777888 ）不包括 2 点和双王
func genTriplet2xTag() {
	// 2个三张
	tag := int64(33)
	keyTags[tag] = makeTagValue(6, prunfast.CardHandType_Triplet2X, 2)

	// 3个三张及以上；最多5个三张，因为每人最多16张牌
	// [关张每个人发16张牌，一共48张牌（没有小丑因此52张牌，去掉一张ace，去掉三张2, 52-4=48）]
	for i := 0; i < 3; i++ {
		tag = tag*10 + 3
		keyTags[tag] = makeTagValue((i+3)*3, prunfast.CardHandType_Triplet2X, 3+i)
	}
}

// 飞机带翅膀
func genTriplet2XXTag() {
	genTriplet2XTagPair()
}

// 飞机带翅膀，加对子
func genTriplet2XTagPair() {
	// 3*n + 2*n <= 16
	// 因此，n最多是3，最多3个三张+3个对子
	for i := 2; i < 4; i++ {
		tag1 := 3
		for j := 1; j < i; j++ {
			tag1 = tag1*10 + 3 // 构建3张的tag
		}

		// 带对子
		tag2 := 2
		for j := 1; j < i; j++ {
			tag2 = tag2*10 + 2
		}

		tag := int64(tag1*int(math.Pow(10, float64(i))) + tag2)
		keyTags[tag] = makeTagValue((i)*5, prunfast.CardHandType_Triplet2X2Pair, i)
	}
}

func genAllTag() {
	genBombTag()
	genSingleTag()
	genPairTag()
	genTripletTag()
	genTripletPairTag()
	genFlushXTag()
	genPair2xTag()
	genTriplet2xTag()
	genTriplet2XXTag()
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

func assertCardHandType(hai []int, cht prunfast.CardHandType) {
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
