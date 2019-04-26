package prunfast

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

	roomScoreConfig *RoomScoreConfig
}

func newRoomConfig() *RoomConfig {
	rc := &RoomConfig{}

	rc.playerNumAcquired = 3
	rc.payNum = 0
	rc.handNum = 2
	rc.payType = 0

	rc.roomScoreConfig = defaultScoreConfig
	return rc
}

func newRoomConfigFromJSON(configJSON *RoomConfigJSON) *RoomConfig {
	rc := newRoomConfig()

	rc.playerNumAcquired = configJSON.PlayerNumAcquired
	rc.payNum = configJSON.PayNum
	rc.payType = configJSON.PayType
	rc.handNum = configJSON.HandNum

	if rc.playerNumAcquired == 0 {
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
