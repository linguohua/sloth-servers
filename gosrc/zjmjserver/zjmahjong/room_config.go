package zjmahjong

import (
	"gconst"

	log "github.com/sirupsen/logrus"

	"github.com/garyburd/redigo/redis"
)

const (
	aapay = 1
)

// RoomConfig 游戏房间配置
type RoomConfig struct {
	/// 牌局要求的参与玩家数量
	playerNumAcquired int
	/// 付费数目，钻石个数
	payNum int

	/// 牌局费用支付类型
	/// 例如创建房间者支付，或者分摊
	// 0表示房主支付，1表示分摊
	payType int

	// 局数
	handNum int

	// 没有风牌
	noWind bool

	// 放杠杠爆后全包
	afterKongChuckerPayForAll bool

	// 杠爆 双倍
	afterKongX2 bool

	// 海底捞双倍
	finalDrawX2 bool

	// 小七对双倍
	sevenPairX2 bool

	// 大七对4倍
	greatSevenPairX4 bool

	// 全风牌双倍
	allWindX2 bool

	// 清一色双倍
	pureSameX2 bool

	// 碰碰胡双倍
	pongpongX2 bool

	// 天胡5倍
	heavenX5 bool

	// 十三幺10倍
	thirteenOrphanX10 bool

	// 抽马个数
	horseCount int

	// 底分
	baseScore int

	// 封顶倍数
	trimMultiple int
}

func newRoomConfig() *RoomConfig {
	rc := &RoomConfig{}

	rc.playerNumAcquired = 4
	rc.payNum = 0
	rc.handNum = 100

	rc.payType = 0

	rc.afterKongChuckerPayForAll = true
	rc.afterKongX2 = true
	rc.finalDrawX2 = true
	rc.sevenPairX2 = true
	rc.greatSevenPairX4 = true
	rc.allWindX2 = true
	rc.pureSameX2 = true
	rc.pongpongX2 = true
	rc.heavenX5 = true
	rc.thirteenOrphanX10 = true

	rc.horseCount = 6
	rc.baseScore = 1

	rc.trimMultiple = 0

	return rc
}

func newRoomConfigFromJSON(configJSON *RoomConfigJSON) *RoomConfig {
	rc := newRoomConfig()

	rc.playerNumAcquired = configJSON.PlayerNumAcquired
	rc.payNum = configJSON.PayNum
	rc.payType = configJSON.PayType
	rc.handNum = configJSON.HandNum

	log.Printf("newRoomConfigFromJSON,json:%+v\n", configJSON)

	rc.noWind = configJSON.NoWind
	rc.afterKongChuckerPayForAll = configJSON.AfterKongChuckerPayForAll
	rc.afterKongX2 = configJSON.AfterKongX2
	rc.finalDrawX2 = configJSON.FinalDrawX2
	rc.sevenPairX2 = configJSON.SevenPairX2
	rc.greatSevenPairX4 = configJSON.GreatSevenPairX4
	rc.allWindX2 = configJSON.AllWindX2
	rc.pureSameX2 = configJSON.PureSameX2
	rc.pongpongX2 = configJSON.PongpongX2
	rc.heavenX5 = configJSON.HeavenX5
	rc.thirteenOrphanX10 = configJSON.ThirteenOrphanX10

	switch configJSON.HorseNumberType {
	case 0:
		rc.horseCount = 4
	case 1:
		rc.horseCount = 6
	case 2:
		rc.horseCount = 8
	case 3:
		rc.horseCount = 12
	default:
		rc.horseCount = 6
	}

	switch configJSON.HorseNumberType {
	case 0:
		rc.baseScore = 1
	case 1:
		rc.baseScore = 2
	case 2:
		rc.baseScore = 5
	case 3:
		rc.baseScore = 10
	default:
		rc.baseScore = 2
	}

	switch configJSON.TrimType {
	case 1:
		rc.trimMultiple = 8
	case 2:
		rc.trimMultiple = 16
	}

	return rc
}

// 天胡倍数
func (rc *RoomConfig) heavenPoint() int {
	if rc.heavenX5 {
		return 5
	}

	return 1
}

// 海底捞倍数
func (rc *RoomConfig) finalDrawPoint() int {
	if rc.finalDrawX2 {
		return 2
	}

	return 1
}

// 杠爆倍数
func (rc *RoomConfig) afterKongPint() int {
	if rc.afterKongX2 {
		return 2
	}

	return 1
}

func (rc *RoomConfig) pureSamePoint() int {
	if rc.pureSameX2 {
		return 2
	}

	return 1
}

func (rc *RoomConfig) sevenPairPoint() int {
	if rc.sevenPairX2 {
		return 2
	}

	return 1
}

func (rc *RoomConfig) greatSevenPairPoint() int {
	if rc.greatSevenPairX4 {
		return 4
	}

	return 2
}

func (rc *RoomConfig) pongPongPoint() int {
	if rc.pongpongX2 {
		return 2
	}

	return 1
}

func (rc *RoomConfig) thirteenOrphansPoint() int {
	if rc.thirteenOrphanX10 {
		return 10
	}

	return 5
}

func (rc *RoomConfig) allWinPoint() int {
	if rc.allWindX2 {
		return 2
	}

	return 1
}

// RoomConfigJSON 游戏房间配置
type RoomConfigJSON struct {
	/// 牌局要求的参与玩家数量
	PlayerNumAcquired int `json:"playerNumAcquired"`
	/// 付费数目，钻石个数
	PayNum int `json:"payNum"`

	/// 牌局费用支付类型
	/// 例如创建房间者支付，或者分摊
	// 0表示房主支付，1表示分摊
	PayType int `json:"payType"`

	// 局数
	HandNum int `json:"handNum"`

	// 没有风牌
	NoWind bool `json:"noWind"`

	// 放杠杠爆后全包
	AfterKongChuckerPayForAll bool `json:"afterKongChuckerPayForAll"`

	// 杠爆 双倍
	AfterKongX2 bool `json:"afterKongX2"`

	// 海底捞双倍
	FinalDrawX2 bool `json:"finalDrawX2"`

	// 小七对双倍
	SevenPairX2 bool `json:"sevenPairX2"`

	// 大七对4倍
	GreatSevenPairX4 bool `json:"greatSevenPairX4"`

	// 全风牌双倍
	AllWindX2 bool `json:"allWindX2"`

	// 清一色双倍
	PureSameX2 bool `json:"pureSameX2"`

	// 碰碰胡双倍
	PongpongX2 bool `json:"pongpongX2"`

	// 天胡5倍
	HeavenX5 bool `json:"heavenX5"`

	// 十三幺10倍
	ThirteenOrphanX10 bool `json:"thirteenOrphanX10"`

	// 抽马类型
	HorseNumberType int `json:"HorseNumberType"`

	// 底分类型
	BaseScoreType int `json:"baseScoreType"`

	// 封顶类型
	TrimType int `json:"trimType"`
}

func loadRoomConfigFromRedis(roomConfigID string) []byte {
	// 获取redis链接，并退出函数时释放
	conn := pool.Get()
	defer conn.Close()

	buf, err := redis.Bytes(conn.Do("hget", gconst.LobbyRoomConfigTable, roomConfigID))
	if err != nil {
		log.Println("load room config from redis :", err)
		return nil
	}

	return buf
}
