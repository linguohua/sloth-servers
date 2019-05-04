package dfmahjong

var (
	dict     = make(map[string]int)
	dictName = make(map[int]string)
)

// MonkeyUserTilesCfg 用来测试
type MonkeyUserTilesCfg struct {
	handTiles   []int
	flowerTiles []int
	actionTips  []string
	userID      string

	isBanker bool
	index    int
}

// newMonkeyUserTilesCfg 新建玩家发牌配置
func newMonkeyUserTilesCfg(isBanker bool, idx int) *MonkeyUserTilesCfg {
	tuc := MonkeyUserTilesCfg{}
	tuc.isBanker = isBanker
	tuc.index = idx
	return &tuc
}

// isValid 判断手牌是否有效
func (tuc *MonkeyUserTilesCfg) isValid() bool {
	if tuc.isBanker {
		if tuc.handTiles != nil && len(tuc.handTiles) == 14 {
			return true
		}
		//log.Println("if tuc.handTiles != nil && len(tuc.handTiles) == 14 not true, tuc.handTiles len: ", len(tuc.handTiles))
		//log.Println("tuc: ", tuc)
		return false
	}

	if tuc.handTiles != nil && len(tuc.handTiles) == 13 {
		return true
	}

	//log.Println("tuc.handTiles != nil && len(tuc.handTiles) == 13 not true")
	//log.Println("tuc: ", tuc)
	return false
}

// setHandTiles 设置发牌时的手牌序列
func (tuc *MonkeyUserTilesCfg) setHandTiles(hands []string) {

	var index = 0
	hl := make([]int, len(hands))
	for _, hand := range hands {
		if hand != "" {
			hl[index] = dict[hand]
			index++
		}
	}

	tuc.handTiles = hl[0:index]
}

// setFlowerTiles 设置发牌时的花牌序列
func (tuc *MonkeyUserTilesCfg) setFlowerTiles(flowers []string) {

	var index = 0
	fl := make([]int, len(flowers))
	for _, flower := range flowers {
		if flower != "" {
			fl[index] = dict[flower]
			index++
		}
	}

	tuc.flowerTiles = fl[0:index]
}

// setActionTips 设置动作提示，用于提示客户端下一步动作是什么
func (tuc *MonkeyUserTilesCfg) setActionTips(actionTips []string) {

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
	dict["1万"] = MAN1
	dict["2万"] = MAN2
	dict["3万"] = MAN3
	dict["4万"] = MAN4
	dict["5万"] = MAN5
	dict["6万"] = MAN6
	dict["7万"] = MAN7
	dict["8万"] = MAN8
	dict["9万"] = MAN9

	dict["1筒"] = PIN1
	dict["2筒"] = PIN2
	dict["3筒"] = PIN3
	dict["4筒"] = PIN4
	dict["5筒"] = PIN5
	dict["6筒"] = PIN6
	dict["7筒"] = PIN7
	dict["8筒"] = PIN8
	dict["9筒"] = PIN9

	dict["1条"] = SOU1
	dict["2条"] = SOU2
	dict["3条"] = SOU3
	dict["4条"] = SOU4
	dict["5条"] = SOU5
	dict["6条"] = SOU6
	dict["7条"] = SOU7
	dict["8条"] = SOU8
	dict["9条"] = SOU9

	dict["东"] = TON
	dict["南"] = NAN
	dict["西"] = SHA
	dict["北"] = PEI
	dict["白"] = HAK
	dict["发"] = HAT
	dict["中"] = CHU
	dict["梅"] = PLUM
	dict["兰"] = ORCHID
	dict["竹"] = BAMBOO
	dict["菊"] = CHRYSANTHEMUM
	dict["春"] = SPRING
	dict["夏"] = SUMMER
	dict["秋"] = AUTUMN
	dict["冬"] = WINTER

	dictName[MAN1] = "1万"
	dictName[MAN2] = "2万"
	dictName[MAN3] = "3万"
	dictName[MAN4] = "4万"
	dictName[MAN5] = "5万"
	dictName[MAN6] = "6万"
	dictName[MAN7] = "7万"
	dictName[MAN8] = "8万"
	dictName[MAN9] = "9万"

	dictName[PIN1] = "1筒"
	dictName[PIN2] = "2筒"
	dictName[PIN3] = "3筒"
	dictName[PIN4] = "4筒"
	dictName[PIN5] = "5筒"
	dictName[PIN6] = "6筒"
	dictName[PIN7] = "7筒"
	dictName[PIN8] = "8筒"
	dictName[PIN9] = "9筒"

	dictName[SOU1] = "1条"
	dictName[SOU2] = "2条"
	dictName[SOU3] = "3条"
	dictName[SOU4] = "4条"
	dictName[SOU5] = "5条"
	dictName[SOU6] = "6条"
	dictName[SOU7] = "7条"
	dictName[SOU8] = "8条"
	dictName[SOU9] = "9条"

	dictName[TON] = "东"
	dictName[NAN] = "南"
	dictName[SHA] = "西"
	dictName[PEI] = "北"
	dictName[HAK] = "白"
	dictName[HAT] = "发"
	dictName[CHU] = "中"
	dictName[PLUM] = "梅"
	dictName[ORCHID] = "兰"
	dictName[BAMBOO] = "竹"
	dictName[CHRYSANTHEMUM] = "菊"
	dictName[SPRING] = "春"
	dictName[SUMMER] = "夏"
	dictName[AUTUMN] = "秋"
	dictName[WINTER] = "冬"
}

func init() {
	initDict()
}
