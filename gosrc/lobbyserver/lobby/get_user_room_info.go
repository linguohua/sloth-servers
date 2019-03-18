package lobby

import (
	"fmt"
	"gconst"
	"strconv"
	log "github.com/sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
)

// UserRoomInfo 用户房间信息
type UserRoomInfo struct {
	RoomType          int32
	RoomNumber        string
	PlayerNumAcquired int32
	GameServerPort    int32
}

// GetUserRoomInfo 获取用户房间信息
func GetUserRoomInfo(userID string) (*UserRoomInfo, error) {
	log.Println("GetUserRoomInfo, userID:", userID)
	enterRoomID := loadUserLastEnterRoomID(userID)
	if enterRoomID == "" {
		return nil, fmt.Errorf("Can't get Enter room ID")
	}

	conn := pool.Get()
	defer conn.Close()

	values, err := redis.Strings(conn.Do("HMGET", gconst.RoomTablePrefix+enterRoomID, "roomNumber", "roomConfigID", "gameServerID", "roomType"))
	if err != nil {
		return nil, fmt.Errorf("Get Room Info error %v", err)
	}

	var roomNumber = values[0]
	var roomConfigID = values[1]
	var gameServerID = values[2]
	var roomType = values[3]

	roomTypeInt, err := strconv.Atoi(roomType)
	if err != nil {
		return nil, fmt.Errorf("Parser roomType to int error %v", err)
	}

	port, err := redis.Int(conn.Do("HGET", gconst.GameServerKeyPrefix+gameServerID, "p"))
	if err != nil {
		return nil, fmt.Errorf("Load game server port error:%v", err)
	}

	roomCfg, ok := roomConfigs[roomConfigID]
	if !ok {
		return nil, fmt.Errorf("Can't get room config")
	}

	var roomCfgJSON = parseRoomConfigFromString(roomCfg)

	var userRoomInfo = &UserRoomInfo{}
	userRoomInfo.RoomType = int32(roomTypeInt)
	userRoomInfo.PlayerNumAcquired = int32(roomCfgJSON.PlayerNumAcquired)
	userRoomInfo.GameServerPort = int32(port)
	userRoomInfo.RoomNumber = roomNumber

	return userRoomInfo, nil
}
