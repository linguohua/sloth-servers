package pddz

import (
	"gconst"
	log "github.com/sirupsen/logrus"

	"github.com/garyburd/redigo/redis"
)

// CardNumType 数量
type CardNumType int

// MiniWinAbleType 小胡胡牌条件类型
type MiniWinAbleType int

// WashoutActionType 荒庄后动作，例如爬坡，例如落庄
type WashoutActionType int

const (
	aapay                                       = 1
	cardTypeNum108            CardNumType       = 1
	cardTypeNum136            CardNumType       = 0
	miniWinAble1              MiniWinAbleType   = 0
	miniWinAble2              MiniWinAbleType   = 1
	washoutActionNone         WashoutActionType = 0
	washoutActionAllDouble    WashoutActionType = 1
	washoutActionBankerDouble WashoutActionType = 2

	maxRobLandlordCount = 2
)

var (
	defaultScoreConfig = newRoomScoreConfig()
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

	// 基础分数
	baseScore int

	// 封顶倍数
	multipleLimit int

	noPriority4Farmer bool

	// 庄家选择方法，0表示下一局由先出完牌的当，1表示随机选择，2表示逆时针下一个人做庄；
	// 不管如何，第一局肯定是随机选择庄家
	nextBankerMethod int

	// 是否是叫分叫地主方式
	isCallWithScore bool

	// 是否加倍
	isCallDoubleEnabled bool

	roomScoreConfig *RoomScoreConfig
}

func newRoomConfig() *RoomConfig {
	rc := &RoomConfig{}

	rc.playerNumAcquired = 3
	rc.payNum = 0
	rc.handNum = 100
	rc.payType = 0
	rc.baseScore = 100
	rc.multipleLimit = 32
	rc.noPriority4Farmer = true
	rc.isCallDoubleEnabled = false

	// 专用于兰州斗地主
	rc.isCallWithScore = false

	rc.roomScoreConfig = defaultScoreConfig
	return rc
}

func newRoomConfigFromJSON(configJSON *RoomConfigJSON) *RoomConfig {
	rc := newRoomConfig()

	rc.playerNumAcquired = configJSON.PlayerNumAcquired
	rc.payNum = configJSON.PayNum
	rc.payType = configJSON.PayType
	rc.handNum = configJSON.HandNum
	rc.baseScore = configJSON.BaseScore
	rc.multipleLimit = configJSON.MultipleLimit
	rc.nextBankerMethod = configJSON.NextBankerMethod
	rc.isCallDoubleEnabled = configJSON.IsCallDoubleEnabled

	rc.isCallWithScore = configJSON.IsCallWithScore

	if rc.playerNumAcquired != 3 && rc.playerNumAcquired != 2 {
		log.Println("force room-config playerNumAcuqired to 3, current:", rc.playerNumAcquired)
		rc.playerNumAcquired = 3
	}

	log.Printf("newRoomConfigFromJSON,json:%+v\n", configJSON)

	rc.roomScoreConfig = defaultScoreConfig
	return rc
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

	// 房间类型
	RoomType int `json:"roomType"`

	// 底分
	BaseScore int `json:"baseScore"`

	// 封顶倍数
	MultipleLimit int `json:"multipleLimit"`

	// 庄家选择方法，0表示下一局由先出完牌的当，1表示随机选择，2表示逆时针下一个人做庄；
	// 不管如何，第一局肯定是随机选择庄家
	NextBankerMethod int `json:"nextBankerMethod"`

	// 是否启用加倍流程
	IsCallDoubleEnabled bool `json:"isCallDoubleEnabled"`

	// 是否叫分叫地主方式
	IsCallWithScore bool `json:"isCallWithScore"`

	// 限制同IP
	LimitSameIP bool `json:"limitSameIp"`

	// 限制同位置
	LimitSameLocation bool `json:"limitSameLocation"`
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
