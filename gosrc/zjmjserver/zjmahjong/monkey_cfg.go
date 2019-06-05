package zjmahjong

import log "github.com/sirupsen/logrus"

// MonkeyCfg 测试用
type MonkeyCfg struct {
	monkeyUserTilesCfgList []*MonkeyUserTilesCfg
	draws                  []int
	name                   string
	isForceConsistent      bool
	isContinuousBanker     bool
}

// newMonkeyCfg 新建一个monkey config
func newMonkeyCfg(name string) *MonkeyCfg {
	handTileConfig := &MonkeyCfg{}
	handTileConfig.name = name
	handTileConfig.monkeyUserTilesCfgList = make([]*MonkeyUserTilesCfg, 0, 4)
	return handTileConfig
}

// playerCount 需要的玩家个数
func (mtc *MonkeyCfg) playerCount() int {
	return len(mtc.monkeyUserTilesCfgList)
}

// getMonkeyUserTilesCfg 获得玩家的发牌配置
func (mtc *MonkeyCfg) getMonkeyUserTilesCfg(userID string) *MonkeyUserTilesCfg {
	for _, cfg := range mtc.monkeyUserTilesCfgList {
		if cfg.userID == userID {
			return cfg
		}
	}

	log.Panicln("no monkey tiles config for:", userID)
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

// isValid 判断配置是否有效
func (mtc *MonkeyCfg) isValid() bool {
	for _, tilesUserCfg := range mtc.monkeyUserTilesCfgList {
		if !tilesUserCfg.isValid() {
			log.Println("!tilesUserCfg.isValid()")
			return false
		}
	}

	if len(mtc.monkeyUserTilesCfgList) < 2 {
		log.Println("len(mtc.tilesPairList) < 2")
		return false
	}

	// 检查牌的张数是否在要求范围内
	slots := make([]int, TILEMAX)
	for _, tuc := range mtc.monkeyUserTilesCfgList {
		for _, handTile := range tuc.handTiles {
			slots[handTile]++
		}

	}

	for _, t := range mtc.draws {
		slots[t]++
	}

	// 非花牌，最多有4张
	for tid, slot := range slots {
		if slot > 4 {
			log.Println("slot > 4 :", tid)
			return false
		}
	}

	// 春夏秋冬，梅兰竹菊只有1张
	for i := PLUM; i < TILEMAX; i++ {
		if slots[i] > 1 {
			log.Println("flow slots[i] > 1:", dictName[i])
			return false
		}
	}

	return true
}
