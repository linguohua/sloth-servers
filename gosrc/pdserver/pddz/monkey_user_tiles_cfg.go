package pddz

var (
	dict     = make(map[string]int)
	dictName = make(map[int]string)
)

// MonkeyUserCardsCfg 用来测试
type MonkeyUserCardsCfg struct {
	handCards []int
	//flowerCards []int
	actionTips []string
	userID     string

	isBanker bool
	index    int
}

// newMonkeyUserCardsCfg 新建玩家发牌配置
func newMonkeyUserCardsCfg(isBanker bool, idx int) *MonkeyUserCardsCfg {
	tuc := MonkeyUserCardsCfg{}
	tuc.isBanker = isBanker
	tuc.index = idx
	return &tuc
}

// isValid 判断手牌是否有效
func (tuc *MonkeyUserCardsCfg) isValid() bool {
	if tuc.isBanker {
		if tuc.handCards != nil && len(tuc.handCards) == 17 {
			return true
		}
		//log.Println("if tuc.handCards != nil && len(tuc.handCards) == 14 not true, tuc.handCards len: ", len(tuc.handCards))
		//log.Println("tuc: ", tuc)
		return false
	}

	if tuc.handCards != nil && len(tuc.handCards) == 17 {
		return true
	}

	//log.Println("tuc.handCards != nil && len(tuc.handCards) == 13 not true")
	//log.Println("tuc: ", tuc)
	return false
}

// setHandCards 设置发牌时的手牌序列
func (tuc *MonkeyUserCardsCfg) setHandCards(hands []string) {

	var index = 0
	hl := make([]int, len(hands))
	for _, hand := range hands {
		if hand != "" {
			hl[index] = dict[hand]
			index++
		}
	}

	tuc.handCards = hl[0:index]
}

// setFlowerCards 设置发牌时的花牌序列
// func (tuc *MonkeyUserCardsCfg) setFlowerCards(flowers []string) {

// 	var index = 0
// 	fl := make([]int, len(flowers))
// 	for _, flower := range flowers {
// 		if flower != "" {
// 			fl[index] = dict[flower]
// 			index++
// 		}
// 	}

// 	tuc.flowerCards = fl[0:index]
// }

// setActionTips 设置动作提示，用于提示客户端下一步动作是什么
func (tuc *MonkeyUserCardsCfg) setActionTips(actionTips []string) {

	var index = 0
	tips := make([]string, len(actionTips))
	for _, tip := range actionTips {
		if tip != "" {
			tips[index] = tip
			index++
		}
	}

	tuc.actionTips = tips[0:index]
}

// initDict 初始化两个map，用于转换牌的字符到序号
func initDict() {
	dict["红桃2"] = R2H
	dict["方块2"] = R2D
	dict["梅花2"] = R2C
	dict["黑桃2"] = R2S

	dict["红桃3"] = R3H
	dict["方块3"] = R3D
	dict["梅花3"] = R3C
	dict["黑桃3"] = R3S

	dict["红桃4"] = R4H
	dict["方块4"] = R4D
	dict["梅花4"] = R4C
	dict["黑桃4"] = R4S

	dict["红桃5"] = R5H
	dict["方块5"] = R5D
	dict["梅花5"] = R5C
	dict["黑桃5"] = R5S

	dict["红桃6"] = R6H
	dict["方块6"] = R6D
	dict["梅花6"] = R6C
	dict["黑桃6"] = R6S

	dict["红桃7"] = R7H
	dict["方块7"] = R7D
	dict["梅花7"] = R7C
	dict["黑桃7"] = R7S

	dict["红桃8"] = R8H
	dict["方块8"] = R8D
	dict["梅花8"] = R8C
	dict["黑桃8"] = R8S

	dict["红桃9"] = R9H
	dict["方块9"] = R9D
	dict["梅花9"] = R9C
	dict["黑桃9"] = R9S

	dict["红桃10"] = R10H
	dict["方块10"] = R10D
	dict["梅花10"] = R10C
	dict["黑桃10"] = R10S

	dict["红桃J"] = JH
	dict["方块J"] = JD
	dict["梅花J"] = JC
	dict["黑桃J"] = JS

	dict["红桃Q"] = QH
	dict["方块Q"] = QD
	dict["梅花Q"] = QC
	dict["黑桃Q"] = QS

	dict["红桃K"] = KH
	dict["方块K"] = KD
	dict["梅花K"] = KC
	dict["黑桃K"] = KS

	dict["红桃A"] = AH
	dict["方块A"] = AD
	dict["梅花A"] = AC
	dict["黑桃A"] = AS

	dict["黑小丑"] = JOB
	dict["红小丑"] = JOR

	dictName[R2H] = "红桃2"
	dictName[R2D] = "方块2"
	dictName[R2C] = "梅花2"
	dictName[R2S] = "黑桃2"

	dictName[R3H] = "红桃3"
	dictName[R3D] = "方块3"
	dictName[R3C] = "梅花3"
	dictName[R3S] = "黑桃3"

	dictName[R4H] = "红桃4"
	dictName[R4D] = "方块4"
	dictName[R4C] = "梅花4"
	dictName[R4S] = "黑桃4"

	dictName[R5H] = "红桃5"
	dictName[R5D] = "方块5"
	dictName[R5C] = "梅花5"
	dictName[R5S] = "黑桃5"

	dictName[R6H] = "红桃6"
	dictName[R6D] = "方块6"
	dictName[R6C] = "梅花6"
	dictName[R6S] = "黑桃6"

	dictName[R7H] = "红桃7"
	dictName[R7D] = "方块7"
	dictName[R7C] = "梅花7"
	dictName[R7S] = "黑桃7"

	dictName[R8H] = "红桃8"
	dictName[R8D] = "方块8"
	dictName[R8C] = "梅花8"
	dictName[R8S] = "黑桃8"

	dictName[R9H] = "红桃9"
	dictName[R9D] = "方块9"
	dictName[R9C] = "梅花9"
	dictName[R9S] = "黑桃9"

	dictName[R10H] = "红桃10"
	dictName[R10D] = "方块10"
	dictName[R10C] = "梅花10"
	dictName[R10S] = "黑桃10"

	dictName[JH] = "红桃J"
	dictName[JD] = "方块J"
	dictName[JC] = "梅花J"
	dictName[JS] = "黑桃J"

	dictName[QH] = "红桃Q"
	dictName[QD] = "方块Q"
	dictName[QC] = "梅花Q"
	dictName[QS] = "黑桃Q"

	dictName[KH] = "红桃K"
	dictName[KD] = "方块K"
	dictName[KC] = "梅花K"
	dictName[KS] = "黑桃K"

	dictName[AH] = "红桃A"
	dictName[AD] = "方块A"
	dictName[AC] = "梅花A"
	dictName[AS] = "黑桃A"

	dictName[JOB] = "黑小丑"
	dictName[JOR] = "红小丑"
}

func init() {
	initDict()
}
