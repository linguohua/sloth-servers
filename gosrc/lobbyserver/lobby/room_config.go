package lobby

import (
	"encoding/json"
	"gconst"

	log "github.com/sirupsen/logrus"

	"github.com/garyburd/redigo/redis"
)

var (
	// RoomConfigs room configs
	RoomConfigs = make(map[string]string)
)

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

	// 开房类型，比如立即开局、代人开房
	OpenType int `json:"openType"`

	// 判断是否是比赛场，0表示普通场，1表示比赛场
	Race int `json:"is_race"`
	// 比赛房进入托管状态的超时时间
	ProxyTimeout int `json:"proxytimeout"`

	// 判断是否是牌友群
	IsGroup bool `json:"isClub"`
}

func loadAllRoomConfigFromRedis() {
	conn := pool.Get()
	defer conn.Close()

	values, err := redis.Strings(conn.Do("HGETALL", gconst.LobbyRoomConfigTable))
	if err != nil {
		log.Println("loadAllRoomConfig, err:", err)
	}

	for i := 0; i < len(values); i = i + 2 {
		var key = values[i]
		var value = values[i+1]
		RoomConfigs[key] = value
	}
}

// ParseRoomConfigFromString json parse
func ParseRoomConfigFromString(roomConfigString string) *RoomConfigJSON {
	var roomConfigJSON = &RoomConfigJSON{}
	err := json.Unmarshal([]byte(roomConfigString), roomConfigJSON)
	if err != nil {
		log.Panicln("parseRoomConfigString", err)
		return nil
	}

	return roomConfigJSON
}

// GetRoomConfig 导出房间配置
func GetRoomConfig(cfgID string) *RoomConfigJSON {
	cfg, ok := RoomConfigs[cfgID]
	if !ok {
		return nil
	}

	return ParseRoomConfigFromString(cfg)
}
