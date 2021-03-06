package room

import (
	"gconst"
	"lobbyserver/lobby"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
)

func getPropCfg(roomType int) string {
	donateUtil := lobby.DonateUtil()
	return donateUtil.GetRoomPropsCfg(roomType)
}

func loadUserLastEnterRoomID(userID string) string {
	log.Println("loadUserLastEnterRoomID, userID:", userID)
	var timeAsLeave = 6 * 60 * 60

	conn := lobby.Pool().Get()
	defer conn.Close()

	fields, err := redis.Strings(conn.Do("HMGET", gconst.LobbyPlayerTablePrefix+userID, "enterRoom", "enterTime", "leaveRoom", "leaveTime"))
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

func loadLastRoomInfo(userID string) *lobby.RoomInfo {
	log.Println("loadLastRoomInfo, userID:", userID)
	enterRoomID := loadUserLastEnterRoomID(userID)
	if enterRoomID == "" {
		return nil
	}

	conn := lobby.Pool().Get()
	defer conn.Close()

	values, err := redis.Strings(conn.Do("HMGET", gconst.LobbyRoomTablePrefix+enterRoomID, "roomNumber", "roomConfigID", "gameServerID", "roomType", "arenaID", "raceTemplateID"))
	if err != nil {
		log.Println("load room info err:", err)
		return nil
	}

	//log.Println("oadLastRoom, enterRoomID:", enterRoomID)
	var roomNumber = values[0]
	var roomConfigID = values[1]
	var gameServerID = values[2]
	var roomType = values[3]
	//log.Printf("loadLastRoom, roomNumber:%s, roomConfigID:%s, gameServerID:%s\n", roomNumber, roomConfigID, gameServerID)
	roomTypeInt, _ := strconv.Atoi(roomType)

	var roomInfo = &lobby.RoomInfo{}
	roomInfo.RoomID = &enterRoomID
	roomInfo.RoomNumber = &roomNumber
	roomInfo.GameServerID = &gameServerID
	var propCfg = getPropCfg(roomTypeInt)
	roomInfo.PropCfg = &propCfg

	//log.Println("loadLastRoom, gserverURL:", gameServerURL)
	roomConfig, ok := lobby.RoomConfigs[roomConfigID]
	if ok {
		roomInfo.Config = &roomConfig
	}

	return roomInfo
}
