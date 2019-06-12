package pddz

import log "github.com/sirupsen/logrus"

// MonkeyCfg 测试用
type MonkeyCfg struct {
	monkeyUserCardsCfgList []*MonkeyUserCardsCfg
	draws                  []int
	// kongDraws              []int // 杠后牌
	// markup                 int   // 上楼计数
	name               string
	isForceConsistent  bool
	isContinuousBanker bool
}

// newMonkeyCfg 新建一个monkey config
func newMonkeyCfg(name string) *MonkeyCfg {
	handCardConfig := &MonkeyCfg{}
	handCardConfig.name = name
	handCardConfig.monkeyUserCardsCfgList = make([]*MonkeyUserCardsCfg, 0, 4)
	return handCardConfig
}

// playerCount 需要的玩家个数
func (mtc *MonkeyCfg) playerCount() int {
	return len(mtc.monkeyUserCardsCfgList)
}

// getMonkeyUserCardsCfg 获得玩家的发牌配置
func (mtc *MonkeyCfg) getMonkeyUserCardsCfg(userID string) *MonkeyUserCardsCfg {
	for _, cfg := range mtc.monkeyUserCardsCfgList {
		if cfg.userID == userID {
			return cfg
		}
	}

	log.Panicln("no monkey cards config for:", userID)
	return nil
}

// drawSeq 抽牌序列
func (mtc *MonkeyCfg) drawSeq(draws []string) {

	index := 0
	d := make([]int, len(draws))
	for _, draw := range draws {
		if draw != "" {
			d[index] = dict[draw]
			index++
		}
	}

	mtc.draws = d[0:index]
}

// kongDrawSeq 抽牌序列
// func (mtc *MonkeyCfg) kongDrawSeq(draws []string) {

// 	index := 0
// 	d := make([]int, len(draws))
// 	for _, draw := range draws {
// 		if draw != "" {
// 			d[index] = dict[draw]
// 			index++
// 		}
// 	}

// 	mtc.kongDraws = d[0:index]
// }

// isValid 判断配置是否有效
func (mtc *MonkeyCfg) isValid() bool {
	for _, cardsUserCfg := range mtc.monkeyUserCardsCfgList {
		if !cardsUserCfg.isValid() {
			log.Println("!cardsUserCfg.isValid()")
			return false
		}
	}

	if len(mtc.monkeyUserCardsCfgList) < 2 {
		log.Println("len(mtc.cardsPairList) < 2")
		return false
	}

	// 检查牌的张数是否在要求范围内
	slots := make([]int, CARDMAX)
	for _, tuc := range mtc.monkeyUserCardsCfgList {
		for _, handCard := range tuc.handCards {
			slots[handCard]++
		}
	}

	// for _, t := range mtc.draws {
	// 	slots[t]++
	// }

	// for _, t := range mtc.kongDraws {
	// 	slots[t]++
	// }

	// 非花牌，最多有4张
	for tid, slot := range slots {
		if slot > 1 {
			log.Println("slot > 4 :", tid)
			return false
		}
	}

	// 春夏秋冬，梅兰竹菊只有1张
	// for i := RankEnd; i < CARDMAX; i++ {
	// 	if slots[i] > 1 {
	// 		log.Println("flow slots[i] > 1:", dictName[i])
	// 		return false
	// 	}
	// }

	return true
}
