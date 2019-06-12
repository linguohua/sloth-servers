package main

import (
	"fmt"
	"log"
	"pddz"
)

var (
	cards = []int{
		4, 4, 4, 4, // 2,3,4,5
		4, 4, 4, 4, // 6,7,8,9
		4, 4, 4, 4, 4, 2, // 10, j,q,k,ace,joke
	}

	// 0   1   2   3   4   5   6   7   8   9   10  11  12 13
	// 2,  3,  4,  5,  6,  7,  8,  9,  10, J,  Q,  K,  A, JO
	// 13, 1, 2,  3, 4,  5, 6,  7,  8, 9,  10,  11,  12,  14
	rank2Priority = []int{12, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 13}
	priority2Rank = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 0, 13}

	keysMap   = make(map[int64]int)
	totalCalc = 0
)

func resetCards() {
	for i := range cards {
		cards[i] = 4
	}

	cards[13] = 2
}

func calcKey(hai []int) int64 {
	// 最多14种牌
	slots := make([]int, 14)
	// 每一种牌，用3bit表达其张数（最多4张，因此3bit可以）
	// 因此一共14*3 = 42bits，没有超过lua5.1的48位int限制
	var key int64

	for _, h := range hai {
		slots[h]++
	}

	for i, s := range slots {
		if s > 4 {
			log.Panicln("slots elem great than 4:", s)
		}
		if s < 3 {
			slots[i] = s
		}
	}

	pulse := false
	for _, s := range slots {
		var sv = int64(s)
		if sv == 0 {
			if pulse {
				pulse = false
			}
		} else {
			if !pulse && key != 0 {
				key = key * 10
			}

			key = key*10 + sv
			pulse = true
		}
	}

	return key
}

func testTriplet2XX() {
	resetCards()

	hai := make([]int, 52)

	for flushLength := 2; flushLength <= 12; flushLength++ {
		// 三顺从3到ACE，不包含2和大小王
		for begin := 1; begin < (14 - flushLength); begin++ {
			haiCount := 0

			for step := 0; step < flushLength; step++ {
				index := step + begin
				hai[haiCount] = index
				haiCount++

				hai[haiCount] = index
				haiCount++

				hai[haiCount] = index
				haiCount++
				cards[index] = 0
			}

			testTriplet2XXChooseN(hai, haiCount, flushLength)
			testTriplet2XXChooseN2(hai, haiCount, flushLength)
			// key := calcKey(hai[0:haiCount])
			// log.Println("key:", key)
			// totalCalc++
			// keysMap[key] = 1

			// 恢复牌表
			for step := 0; step < flushLength; step++ {
				index := step + begin

				cards[index] = 4
			}
		}
	}

}

func testTriplet2XXChooseN(hai []int, haiCount int, flushLength int) {
	if haiCount > 20 {
		return
	}

	if flushLength == 0 {
		key := calcKey(hai[0:haiCount])
		// log.Println("key:", key)
		totalCalc++
		keysMap[key] = 1

		assertCardHandType(hai[0:haiCount], pddz.CardHandType_Triplet2X2Single)
		return
	}

	// 先选择单张
	for i := 0; i < 14; i++ {
		if cards[i] > 1 {
			cards[i] = cards[i] - 1

			hai[haiCount] = i
			haiCount++

			flushLength--
			testTriplet2XXChooseN(hai, haiCount, flushLength)
			flushLength++

			// 恢复
			cards[i] = cards[i] + 1
			haiCount = haiCount - 1
		}
	}
}

func testTriplet2XXChooseN2(hai []int, haiCount int, flushLength int) {
	if haiCount > 20 {
		return
	}

	if flushLength == 0 {
		key := calcKey(hai[0:haiCount])
		//log.Println("key:", key)
		totalCalc++
		keysMap[key] = 1
		assertCardHandType(hai[0:haiCount], pddz.CardHandType_Triplet2X2Pair)
		return
	}

	// 选择对子，对子不包含鬼
	for i := 0; i < 13; i++ {
		if cards[i] == 4 {
			cards[i] = cards[i] - 2

			hai[haiCount] = i
			haiCount++
			hai[haiCount] = i
			haiCount++

			flushLength--
			testTriplet2XXChooseN2(hai, haiCount, flushLength)
			flushLength++

			// 恢复
			cards[i] = cards[i] + 2
			haiCount = haiCount - 2
		}
	}
}

func testSingle() {
	resetCards()

	hai := make([]int, 52)
	haiCount := 0

	for i := 0; i < 14; i++ {
		haiCount = 0

		hai[haiCount] = i
		haiCount++

		totalCalc++
		assertCardHandType(hai[:haiCount], pddz.CardHandType_Single)
	}
}

func testPair() {
	resetCards()

	hai := make([]int, 52)
	haiCount := 0

	for i := 0; i < 14; i++ {
		haiCount = 0

		hai[haiCount] = i
		haiCount++
		hai[haiCount] = i
		haiCount++

		totalCalc++
		assertCardHandType(hai[:haiCount], pddz.CardHandType_Pair)
	}
}

func testTriplet() {
	resetCards()

	hai := make([]int, 52)
	haiCount := 0

	// 不包含鬼牌
	for i := 0; i < 13; i++ {
		haiCount = 0

		hai[haiCount] = i
		haiCount++
		hai[haiCount] = i
		haiCount++
		hai[haiCount] = i
		haiCount++

		totalCalc++
		assertCardHandType(hai[:haiCount], pddz.CardHandType_Triplet)
	}
}

func testBomb() {
	resetCards()

	hai := make([]int, 52)
	haiCount := 0

	// 不包含鬼牌
	for i := 0; i < 13; i++ {
		haiCount = 0

		hai[haiCount] = i
		haiCount++
		hai[haiCount] = i
		haiCount++
		hai[haiCount] = i
		haiCount++
		hai[haiCount] = i
		haiCount++

		totalCalc++
		assertCardHandType(hai[:haiCount], pddz.CardHandType_Bomb)
	}
}

func testPair3X() {
	resetCards()

	hai := make([]int, 52)
	haiCount := 0

	for xlen := 3; xlen < 11; xlen++ {
		for begin := 0; begin < (14 - xlen); begin++ {
			haiCount = 0
			for i := 0; i < xlen; i++ {
				hai[haiCount] = begin + i
				haiCount++

				hai[haiCount] = begin + i
				haiCount++
			}

			totalCalc++
			assertCardHandType(hai[:haiCount], pddz.CardHandType_Pair3X)
		}
	}
}

func testFlush() {
	resetCards()

	hai := make([]int, 52)
	haiCount := 0

	for xlen := 5; xlen < 13; xlen++ {
		for begin := 0; begin < (14 - xlen); begin++ {
			haiCount = 0
			for i := 0; i < xlen; i++ {
				hai[haiCount] = begin + i
				haiCount++
			}

			totalCalc++
			assertCardHandType(hai[:haiCount], pddz.CardHandType_Flush)
		}
	}
}

func testFourX2() {
	resetCards()

	hai := make([]int, 52)
	haiCount := 0

	// 不包含鬼牌
	for i := 0; i < 13; i++ {
		haiCount = 0
		for j := 0; j < 4; j++ {
			hai[haiCount] = i
			haiCount++
		}

		cards[i] = 0
		testFourX2Pairs(hai, haiCount)
		testFourX2Single(hai, haiCount)

		cards[i] = 4
	}
}

func testFourX2Pairs(hai []int, haiCount int) {
	if haiCount == 8 {
		totalCalc++
		assertCardHandType(hai[:haiCount], pddz.CardHandType_FourX2Pair)

		return
	}

	for i := 0; i < 13; i++ {
		if cards[i] == 4 {
			cards[i] = 2

			hai[haiCount] = i
			haiCount++
			hai[haiCount] = i
			haiCount++

			testFourX2Pairs(hai, haiCount)

			cards[i] = 4
			haiCount = haiCount - 2
		}
	}
}

func testFourX2Single(hai []int, haiCount int) {
	if haiCount == 6 {
		totalCalc++
		assertCardHandType(hai[:haiCount], pddz.CardHandType_FourX2Single)

		return
	}

	for i := 0; i < 14; i++ {
		if cards[i] > 1 {
			cards[i]--

			hai[haiCount] = i
			haiCount++

			testFourX2Single(hai, haiCount)

			cards[i]++
			haiCount = haiCount - 1
		}
	}
}

func nchoosek(m int, k int) int {
	a := 1
	for i := 1; i <= m; i++ {
		a = a * i
	}

	b1 := 1
	b2 := 1
	for i := 1; i <= k; i++ {
		b1 = b1 * i
	}

	for i := 1; i <= (m - k); i++ {
		b2 = b2 * i
	}

	return a / (b1 * b2)
}

func dumpKeysTag() {
	for k, v := range keyTags {
		fmt.Printf("agariTable[0x%x]=0x%x\n", k, v)
	}
}

func main() {
	// genTriplet2XX()

	// for k := range keysMap {
	// 	log.Printf("key:%d\n", k)
	// }

	// log.Printf("keys count:%d, total:%d\n", len(keysMap), totalCalc)
	genAllTag()

	// genTriplet2XTag()
	// for t := range keyTags {
	// 	log.Printf("tag:%d\n", t)
	// }
	dumpKeysTag()

	log.Printf("tag count:%d\n", len(keyTags))

	testTriplet2XX()
	testSingle()
	testPair()
	testTriplet()
	testBomb()

	testPair3X()
	testFlush()
	testFourX2()

	// for k := range keysMap {
	// 	log.Printf("key:%d\n", k)
	// }

	log.Printf("keys count:%d, total:%d\n", len(keysMap), totalCalc)
}
