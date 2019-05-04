package dfmahjong

import (
	"gconst"
	"log"

	"github.com/garyburd/redigo/redis"
)

const (
	aapay = 1
)

var (
	roomScoreConfigs []*RoomScoreConfig
)

// RoomConfig 游戏房间配置
type RoomConfig struct {
	/// 牌局要求的参与玩家数量
	playerNumAcquired int
	/// 付费数目，钻石个数
	payNum int
	/// 墩子分类型
	dunziPointType int
	/// 牌局费用支付类型
	/// 例如创建房间者支付，或者分摊
	// 0表示房主支付，1表示分摊
	payType int

	// 局数
	handNum int

	// 暗杠的墩子分
	dunzi4ConcealedKong int
	// 明杠的墩子分
	dunzi4ExposedKong int
	// 4个花牌的墩子分
	dunzi4QuadFlower int

	// 自摸X2
	isDoubleScoreWhenSelfDrawn bool
	// 连庄X2
	isDoubleScoreWhenContinuousBanker bool

	// 封顶类型
	greatWinTrimType int

	// 是否启用输分保护，也即是“坐园子”
	// 如果是的话，则玩家输牌到一定的分数后，就不再输分了
	loseProtectEnabled bool

	roomScoreConfig *RoomScoreConfig
}

func newRoomConfig() *RoomConfig {
	rc := &RoomConfig{}

	rc.playerNumAcquired = 4
	rc.payNum = 0
	rc.handNum = 100
	rc.dunziPointType = 0
	rc.payType = 0
	rc.dunzi4ConcealedKong = 2
	rc.dunzi4ExposedKong = 1
	rc.dunzi4QuadFlower = 3
	rc.isDoubleScoreWhenSelfDrawn = true
	rc.isDoubleScoreWhenContinuousBanker = true
	rc.greatWinTrimType = 1
	rc.roomScoreConfig = roomScoreConfigs[rc.greatWinTrimType]
	return rc
}

func newRoomConfigFromJSON(configJSON *RoomConfigJSON) *RoomConfig {
	rc := newRoomConfig()

	rc.playerNumAcquired = configJSON.PlayerNumAcquired
	rc.payNum = configJSON.PayNum
	rc.payType = configJSON.PayType
	rc.handNum = configJSON.HandNum

	rc.dunziPointType = configJSON.DunziPointType
	log.Printf("newRoomConfigFromJSON,json:%+v\n", configJSON)

	switch rc.dunziPointType {
	case 0:
		rc.dunzi4ExposedKong = 1
		rc.dunzi4ConcealedKong = 2
		rc.dunzi4QuadFlower = 3
		break
	case 1:
		rc.dunzi4ExposedKong = 2
		rc.dunzi4ConcealedKong = 4
		rc.dunzi4QuadFlower = 6
		break
	case 2:
		rc.dunzi4ExposedKong = 5
		rc.dunzi4ConcealedKong = 10
		rc.dunzi4QuadFlower = 15
		break
	case 3:
		rc.dunzi4ExposedKong = 10
		rc.dunzi4ConcealedKong = 20
		rc.dunzi4QuadFlower = 30
		break
	}

	rc.isDoubleScoreWhenSelfDrawn = configJSON.DoubleScoreWhenSelfDrawn
	rc.isDoubleScoreWhenContinuousBanker = configJSON.DoubleScoreWhenContinuousBanker

	var fenDingType = configJSON.FengDingType
	rc.greatWinTrimType = fenDingType
	rc.roomScoreConfig = roomScoreConfigs[rc.greatWinTrimType]

	if configJSON.ZuoYuanZiEnabled {
		rc.loseProtectEnabled = true
	}

	return rc
}

func initDefaulScoreMap() {
	roomScoreConfigs = make([]*RoomScoreConfig, 4)
	for i := 0; i < 4; i++ {
		roomScoreConfigs[i] = newRoomScoreConfig(i)
	}
}

func (rc *RoomConfig) greatWinPointMap(greatWinType GreatWinType) float32 {
	var rsc = roomScoreConfigs[rc.greatWinTrimType]
	v, ok := rsc.greatWinPointMap[greatWinType]
	if !ok {
		// 如果没有配置，直接返回封顶
		log.Printf("greatWinPointMap no map for:%d, return max greatWinPoints:%f\n", greatWinType, rsc.maxGreatWinPoints)
		return rsc.maxGreatWinPoints
	}

	return v
}

// RoomConfigJSON 游戏房间配置
type RoomConfigJSON struct {
	/// 牌局要求的参与玩家数量
	PlayerNumAcquired int `json:"playerNumAcquired"`
	/// 付费数目，钻石个数
	PayNum int `json:"payNum"`
	/// 墩子分类型,0表示1分/2分， 1表示2分/4分
	DunziPointType int `json:"dunziPointType"`
	///封顶类型，0表示20/40, 1表示30/60, 2表示50/100/150, 3表示100/200/300
	FengDingType int `json:"fengDingType"`
	/// 牌局费用支付类型
	/// 例如创建房间者支付，或者分摊
	// 0表示房主支付，1表示分摊
	PayType int `json:"payType"`

	// 局数
	HandNum int `json:"handNum"`

	// 自摸X2
	DoubleScoreWhenSelfDrawn bool `json:"doubleScoreWhenSelfDrawn"`
	// 连庄X2
	DoubleScoreWhenContinuousBanker bool `json:"doubleScoreWhenContinuousBanker"`

	ZuoYuanZiEnabled bool `json:"doubleScoreWhenZuoYuanZi"`
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

func init() {
	initDefaulScoreMap()
}
