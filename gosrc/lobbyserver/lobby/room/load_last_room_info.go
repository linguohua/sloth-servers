package room

import (
	log "github.com/sirupsen/logrus"
	"gconst"
	"strconv"
	"time"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

func getPropCfg(roomType int)string{
	// conn := pool.Get()
	// defer conn.Close()

	// propCfgString, err := redis.String(conn.Do("HGET", gconst.GamePropsCfgTable, roomType))
	// if err != nil {
	// 	log.Println("loadPropCfg error:", err)
	// 	return ""
	// }
	clientPropCfgMap := clientPropCfgsMap[roomType]

	buf, err := json.Marshal(clientPropCfgMap)
	if err != nil {
		return ""
	}

	return string(buf)
}

func loadUserLastEnterRoomID(userID string) string {
	log.Println("loadUserLastEnterRoomID, userID:", userID)
	var timeAsLeave = 6 * 60 * 60

	conn := pool.Get()
	defer conn.Close()

	fields, err := redis.Strings(conn.Do("HMGET", gconst.PlayerTablePrefix+userID, "enterRoom", "enterTime", "leaveRoom", "leaveTime"))
	if err != nil {
		log.Println("loadLastRoomNumber err:", err)
		return ""
	}

	var enterRoomID = fields[0]
	var enterTime = fields[1]
	var leaveRoomID = fields[2]
	var leaveTime = fields[3]

	if enterRoomID == "" {
		return ""
	}

	enterTimeInt64, _ := strconv.ParseUint(enterTime, 10, 32)
	leaveTimeInt64, _ := strconv.ParseUint(leaveTime, 10, 32)

	var isLeaveRoom = false

	if enterRoomID == leaveRoomID {
		if enterTimeInt64 < leaveTimeInt64 {
			isLeaveRoom = true
		}
	}

	var nowTime = time.Now().Unix()
	var diff = uint64(nowTime) - enterTimeInt64
	if diff > uint64(timeAsLeave) {
		isLeaveRoom = true
	}

	if isLeaveRoom {
		return ""
	}

	return enterRoomID
}

func loadLastRoomInfo(userID string) *RoomInfo {
	log.Println("loadLastRoomInfo, userID:", userID)
	enterRoomID := loadUserLastEnterRoomID(userID)
	if enterRoomID == "" {
		return nil
	}

	conn := pool.Get()
	defer conn.Close()

	values, err := redis.Strings(conn.Do("HMGET", gconst.RoomTablePrefix+enterRoomID, "roomNumber", "roomConfigID", "gameServerID", "roomType", "arenaID", "raceTemplateID"))
	if err != nil {
		log.Println("load room info err:", err)
		return nil
	}

	//log.Println("oadLastRoom, enterRoomID:", enterRoomID)
	var roomNumber = values[0]
	var roomConfigID = values[1]
	var gameServerID = values[2]
	var roomType = values[3]
	var arenaID = values[4]
	var raceTemplateID = values[5]
	//log.Printf("loadLastRoom, roomNumber:%s, roomConfigID:%s, gameServerID:%s\n", roomNumber, roomConfigID, gameServerID)
	var gameServerURL = getGameServerURL(gameServerID)

	if gameServerURL == "" {
		log.Printf("loadLastRoom, roomID:%s, invalid last room record, gameServerURL is nil\n", enterRoomID)
		return nil
	}

	roomTypeInt, _:=strconv.Atoi(roomType)

	var roomInfo = &RoomInfo{}
	roomInfo.RoomID = &enterRoomID
	roomInfo.RoomNumber = &roomNumber
	roomInfo.GameServerURL = &gameServerURL
	roomInfo.ArenaID =&arenaID
	roomInfo.RaceTemplateID = &raceTemplateID
	var propCfg = getPropCfg(roomTypeInt)
	roomInfo.PropCfg = &propCfg

	//log.Println("loadLastRoom, gserverURL:", gameServerURL)
	roomConfig, ok := roomConfigs[roomConfigID]
	if ok {
		roomInfo.Config = &roomConfig
	}

	return roomInfo
}
