package pddz

import (
	log "github.com/sirupsen/logrus"
	"pokerface"
	"sort"
)

// 定义扑克牌ID
const (
	R2H = 0 // 2
	R2D = 1
	R2C = 2
	R2S = 3

	R3H = 4 // 3
	R3D = 5
	R3C = 6
	R3S = 7

	R4H = 8 // 4
	R4D = 9
	R4C = 10
	R4S = 11

	R5H = 12 // 5
	R5D = 13
	R5C = 14
	R5S = 15

	R6H = 16 // 6
	R6D = 17
	R6C = 18
	R6S = 19

	R7H = 20 // 7
	R7D = 21
	R7C = 22
	R7S = 23

	R8H = 24 // 8
	R8D = 25
	R8C = 26
	R8S = 27

	R9H = 28 // 9
	R9D = 29
	R9C = 30
	R9S = 31

	R10H = 32 // 10
	R10D = 33
	R10C = 34
	R10S = 35

	JH = 36 // Jack
	JD = 37
	JC = 38
	JS = 39

	QH = 40 // Queen
	QD = 41
	QC = 42
	QS = 43

	KH = 44 // King
	KD = 45
	KC = 46
	KS = 47

	AH = 48 // ACE
	AD = 49
	AC = 50
	AS = 51

	RankEnd = 52 // rank 牌分界

	JOB = 52 // joker black，黑小丑
	JOR = 53 // joker red，红小丑

	CARDMAX = 54
)

var (
	agariTable = make(map[int64]int)

	// 0   1   2   3   4   5   6   7   8   9   10  11  12 13
	// 2,  3,  4,  5,  6,  7,  8,  9,  10, J,  Q,  K,  A, JO
	// 13, 1, 2,  3, 4,  5, 6,  7,  8, 9,  10,  11,  12,  14
	rank2Priority = []int{12, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 13}
	priority2Rank = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 0, 13}
)

func calcKey(hai []int32) (int64, []int) {
	// 最多14种牌
	slots := make([]int, 14)

	for _, h := range hai {
		slots[h/4]++
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

	for i := 0; i < 14; i++ {
		slots[i] = 0
	}

	for _, h := range hai {
		slots[h/4]++
	}

	return tag, slots
}

func agariPureTriplet2X(slots []int) bool {
	for i := 0; i < len(slots); i++ {
		s := slots[i]
		if s == 0 {
			continue
		}

		if s != 3 {
			return false
		}

		end := i
		// 开始遍历顺子
		for j := i; j < len(slots); j++ {
			if slots[j] != 3 {
				end = j
				break
			}
		}

		for j := end; j < len(slots); j++ {
			if slots[j] != 0 {
				return false
			}
		}

		break
	}

	return true
}

func agariSearchContinousFlush(slots []int, itemCount int, flushLength int, flushSlots []bool) bool {
	log.Printf("agariSearchContinousFlush, itemCount:%d, flushLength:%d\n", itemCount, flushLength)
	// 跳过2和小丑
	for i := 12; i >= (flushLength); {
		found := true
		// 寻找一个长度为flushLength的顺子
		for j := i; j > (i - flushLength); j-- {
			if slots[j] != itemCount {
				found = false
				i = j - 1
				break
			}
		}

		if found {
			// 填充flushSlots
			for j := i; j > (i - flushLength); j-- {
				flushSlots[j] = true
			}
			return true
		}
	}

	return false
}

func agariVerifyFlushX(slots []int, ct CardHandType, flushLength int, flushSlots []bool) bool {

	switch ct {
	case CardHandType_Flush:
		return agariSearchContinousFlush(slots, 1, flushLength, flushSlots)
	case CardHandType_Pair3X:
		return agariSearchContinousFlush(slots, 2, flushLength, flushSlots)
	case CardHandType_Triplet2X, CardHandType_Triplet2X2Pair, CardHandType_Triplet2X2Single:
		return agariSearchContinousFlush(slots, 3, flushLength, flushSlots)
	default:
	}

	return false
}

// convertMsgCardHand 转换为MsgCardHand
func agariConvertMsgCardHand(hai []int32) *pokerface.MsgCardHand {
	key, slots := calcKey(hai)

	agari, ok := agariTable[key]
	if !ok {
		log.Println("invalid hai")
		return nil
	}

	ct := CardHandType(agari & 0x00ff)
	flushX := (agari >> 16) & 0x00ff

	flushSlots := make([]bool, len(slots))
	// 如果是顺子（单顺；双顺；三顺；飞机带翅膀），则需要检查
	// 顺子连续性；以及顺子不能包含2和王；
	if flushX > 0 {
		if !agariVerifyFlushX(slots, ct, flushX, flushSlots) {
			log.Println("hai not valid flushX")
			return nil
		}
	}

	// 如果是对子，则检查是否大小王，是的话需要转换为火箭
	if ct == CardHandType_Pair && slots[13] == 2 {
		ct = CardHandType_Roket
	}

	// 携带对子的牌型，限制大小王不能组成对子
	if ct != CardHandType_Roket {
		// 携带对子的牌型
		if slots[13] == 2 {
			log.Println("JOKERs normal pair not allowed")
			return nil
		}
	}

	// 如果是飞机带单张翅膀，如果翅膀是3张构成，而且3张也和飞机连续，则需要转换为
	// 飞机，由于单张飞机满足4*N，而3顺是3*M, 而且牌数小于等于20，因此只有12张牌时
	// 才会出现即是3顺也是飞机带翅膀的情况
	if ct == CardHandType_Triplet2X2Single {
		if len(hai) == 12 && agariPureTriplet2X(slots) {
			ct = CardHandType_Triplet2X
		}
	}

	log.Printf("convertMsgCardHand, agarix:%x, ct:%d\n", agari, ct)

	cardHand := &pokerface.MsgCardHand{}
	var cardHandType32 = int32(ct)
	cardHand.CardHandType = &cardHandType32

	sort.Slice(hai, func(i, j int) bool {
		rankI := hai[i] / 4
		rankJ := hai[j] / 4
		if rankI == rankJ {
			return hai[i] > hai[j]
		}

		return rank2Priority[rankI] > rank2Priority[rankJ]
	})

	haiNew := make([]int32, 0, len(hai)+1)

	// 构造序列，如果是带对子则对子位于后面
	switch ct {
	case CardHandType_TripletPair, CardHandType_Triplet2X2Pair,
		CardHandType_TripletSingle, CardHandType_Triplet2X2Single:
		// 顺子3张放前面
		for _, v := range hai {
			if slots[v/4] == 3 && flushSlots[v/4] {
				haiNew = append(haiNew, int32(v))
			}
		}

		// 非顺子3张放后面
		for _, v := range hai {
			if slots[v/4] == 3 && !flushSlots[v/4] {
				haiNew = append(haiNew, int32(v))
			}
		}

		// 2张，1张放最后
		for _, v := range hai {
			if slots[v/4] != 3 {
				haiNew = append(haiNew, int32(v))
			}
		}
	case CardHandType_FourX2Single, CardHandType_FourX2Pair:
		// 对子放后面
		for _, v := range hai {
			if slots[v/4] == 4 {
				haiNew = append(haiNew, int32(v))
			}
		}
		for _, v := range hai {
			if slots[v/4] != 4 {
				haiNew = append(haiNew, int32(v))
			}
		}
		break
	default:
		for _, v := range hai {
			haiNew = append(haiNew, int32(v))
		}
	}

	cardHand.Cards = haiNew
	return cardHand
}

func init() {
	agariTable[0xa98ac7] = 0x80801
	agariTable[0x14d] = 0x30909
	agariTable[0x32dcd5] = 0x71509
	agariTable[0x515a6] = 0x30f0a
	agariTable[0x32dcd4] = 0x5140b
	agariTable[0x1b207] = 0x60601
	agariTable[0x10f447] = 0x70701
	agariTable[0x56ce] = 0x50a05
	agariTable[0x14c] = 0x2080b
	agariTable[0x51613] = 0x4100b
	agariTable[0x8235] = 0x50f09
	agariTable[0x51615] = 0x61209
	agariTable[0x8229] = 0x30c0b
	agariTable[0x1] = 0x103
	agariTable[0x3] = 0x306
	agariTable[0x21e88e] = 0x70e05
	agariTable[0x21] = 0x20609
	agariTable[0xd05] = 0x30c0b
	agariTable[0x5160a] = 0x4100b
	agariTable[0x1f] = 0x407
	agariTable[0x153158e] = 0x81005
	agariTable[0xcef] = 0x2080b
	agariTable[0x51537] = 0x30c0b
	agariTable[0xcfa] = 0x20a0a
	agariTable[0x1fc9bfe] = 0x4140a
	agariTable[0x1fc97a7] = 0x4100b
	agariTable[0x19debd01c7] = 0xc0c01
	agariTable[0xde] = 0x30505
	agariTable[0x8ae] = 0x40805
	agariTable[0x3640e] = 0x60c05
	agariTable[0xd3ed78e] = 0x91205
	agariTable[0xc6ae4a87] = 0x5140b
	agariTable[0xc6ae75ee] = 0x5190a
	agariTable[0x1fc9fe5] = 0x5140b
	agariTable[0x4] = 0x402
	agariTable[0x20] = 0x508
	agariTable[0x2b67] = 0x50501
	agariTable[0x69f6bc7] = 0x90901
	agariTable[0x2964619c7] = 0xb0b01
	agariTable[0x2] = 0x204
	agariTable[0x84746b8e] = 0xa1405
	agariTable[0x1fca03f] = 0x5140b
	agariTable[0x1a6] = 0x80d
	agariTable[0x423a35c7] = 0xa0a01
	agariTable[0x32dc5b] = 0x4100b
	agariTable[0x13de3e8f] = 0x5140b
	agariTable[0x19b] = 0x60c
	agariTable[0x2a] = 0x60c
}
