package zjmahjong

import (
	"gconst"
	"log"

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
}

func newRoomConfig() *RoomConfig {
	rc := &RoomConfig{}

	rc.playerNumAcquired = 4
	rc.payNum = 0
	rc.handNum = 100

	rc.payType = 0

	return rc
}

func newRoomConfigFromJSON(configJSON *RoomConfigJSON) *RoomConfig {
	rc := newRoomConfig()

	rc.playerNumAcquired = configJSON.PlayerNumAcquired
	rc.payNum = configJSON.PayNum
	rc.payType = configJSON.PayType
	rc.handNum = configJSON.HandNum

	log.Printf("newRoomConfigFromJSON,json:%+v\n", configJSON)

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
