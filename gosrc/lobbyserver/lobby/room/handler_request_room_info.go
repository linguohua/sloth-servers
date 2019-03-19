package room

import (
	"fmt"
	"gconst"
	"lobbyserver/lobby"
	"lobbyserver/lobby/pay"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
)

func replyJoinClubRoomError(w http.ResponseWriter, errorCode int32, clubID string) {
	conn := lobby.Pool().Get()
	defer conn.Close()
	clubNumber, _ := redis.String(conn.Do("HGET", gconst.ClubTablePrefix+clubID, "clubNumber"))

	msgRequestRoomInfoRsp := &lobby.MsgRequestRoomInfoRsp{}
	msgRequestRoomInfoRsp.Result = proto.Int32(errorCode)
	var errString = lobby.ErrorString[errorCode]
	errString = fmt.Sprintf(errString, clubNumber)
	msgRequestRoomInfoRsp.RetMsg = &errString

	reply(w, msgRequestRoomInfoRsp, int32(lobby.MessageCode_OPRequestRoomInfo))
}

func replyRequestRoomInfo(w http.ResponseWriter, errorCode int32, roomInfo *lobby.RoomInfo) {

	msgRequestRoomInfoRsp := &lobby.MsgRequestRoomInfoRsp{}
	msgRequestRoomInfoRsp.Result = proto.Int32(errorCode)
	var errString = lobby.ErrorString[errorCode]
	msgRequestRoomInfoRsp.RetMsg = &errString
	msgRequestRoomInfoRsp.RoomInfo = roomInfo

	reply(w, msgRequestRoomInfoRsp, int32(lobby.MessageCode_OPRequestRoomInfo))
}

func isFullRoom(roomID string, userID string, conn redis.Conn, roomConfigString string) bool {
	// 判断房间是否已经满
	buf, err := redis.Bytes(conn.Do("HGET", gconst.GsRoomTablePrefix+roomID, "players"))
	if err != nil {
		log.Println("Get room players failed:", err)
		return false
	}

	userIDList := &gconst.SSMsgUserIDList{}
	err = proto.Unmarshal(buf, userIDList)
	if err != nil {
		log.Println("Unmarshal failed:", err)
		return false
	}

	var isUserInRoom = false
	var userIDs = userIDList.GetUserIDs()
	for _, uid := range userIDs {
		if uid == userID {
			isUserInRoom = true
			break
		}
	}

	roomConfigJSON := lobby.ParseRoomConfigFromString(roomConfigString)
	// 如果PlayerNumAcquired为0，则无限人数
	if roomConfigJSON.PlayerNumAcquired == 0 {
		return false
	}

	if !isUserInRoom && roomConfigJSON.PlayerNumAcquired == len(userIDs) {
		return true
	}

	return false

}

func isClubMember(clubID string, userID string) bool {
	conn := lobby.Pool().Get()
	defer conn.Close()

	result, err := redis.Int(conn.Do("SISMEMBER", gconst.ClubMemberSetPrefix+clubID, userID))
	if err != nil {
		log.Println("isCanJoinClubRoom, error:", err)
		return false
	}

	if result != 1 {
		return false
	}

	return true
}

func isGroupMember(groupID string, userID string) bool {
	conn := lobby.Pool().Get()
	defer conn.Close()

	vs, err := redis.Strings(conn.Do("HGETALL", gconst.UserClubTablePrefix+userID))
	if err != nil {
		return false
	}

	for i := 0; i < len(vs); i = i + 2 {
		id := vs[i]
		if id == groupID {
			return true
		}
	}

	return false
}

func handlerRequestRoomInfo(w http.ResponseWriter, r *http.Request, userID string) {
	log.Println("handlerRequestRoomInfo call, userID:", userID)
	// 1. 从请求中获取房间6位数字ID
	// 2. 检查房间有效性，比如是否存在，是否已经满了
	// 3. 用房间6位数字ID去请求房间id
	// 4. 获取房间所在服务器的ID
	// 5. 用服务器ID去获取服务器的URL
	// if isUserInBlacklist(userID) {
	// 	log.Println("onMessageRequestRoomInfo,user in black list")
	// 	replyRequestRoomInfo(w, int32(lobby.MsgError_ErrUserInBlacklist), nil)
	// 	return
	// }

	accessoryMessage, errCode := parseAccessoryMessage(r)
	if errCode != int32(lobby.MsgError_ErrSuccess) {
		replyRequestRoomInfo(w, errCode, nil)
	}

	bytes := accessoryMessage.GetData()

	var msgRequestRoomInfo = &lobby.MsgRequestRoomInfo{}
	err := proto.Unmarshal(bytes, msgRequestRoomInfo)
	if err != nil {
		log.Println("onMessageRequestRoomInfo,1 Unmarshal err:", err)
		replyRequestRoomInfo(w, int32(lobby.MsgError_ErrDecode), nil)
		return
	}

	var roomNumber = msgRequestRoomInfo.GetRoomNumber()
	if roomNumber == "" {
		log.Println("onMessageRequestRoomInfo roomNumber is empty")
		replyRequestRoomInfo(w, int32(lobby.MsgError_ErrRoomNumberIsEmpty), nil)
		return
	}

	var lastRoomInfo = loadLastRoomInfo(userID)
	if lastRoomInfo != nil {
		if lastRoomInfo.GetRoomNumber() == roomNumber {
			replyRequestRoomInfo(w, int32(lobby.MsgError_ErrSuccess), lastRoomInfo)
			return
		}

		log.Printf("handlerRequestRoomInfo, User %s in other room, roomNumber: %s, roomId:%s", userID, lastRoomInfo.GetRoomNumber(), lastRoomInfo.GetRoomID())
		replyRequestRoomInfo(w, int32(lobby.MsgError_ErrUserInOtherRoom), lastRoomInfo)
		return
	}

	conn := lobby.Pool().Get()
	defer conn.Close()

	roomID, err := redis.String(conn.Do("HGET", gconst.RoomNumberTable+roomNumber, "roomID"))
	if err != nil && err != redis.ErrNil {
		log.Println("onMessageRequestRoomInfo get roomID err: ", err)
		replyRequestRoomInfo(w, int32(lobby.MsgError_ErrDatabase), nil)
		return
	}

	if roomID == "" {
		log.Println("roomNumber not exist")
		replyRequestRoomInfo(w, int32(lobby.MsgError_ErrRoomNumberNotExist), nil)
		return
	}

	values, err := redis.Strings(conn.Do("HMGET", gconst.RoomTablePrefix+roomID, "roomConfigID", "gameServerID", "groupID", "roomType", "arenaID"))
	if err != nil {
		log.Println("onMessageRequestRoomInfo get roomConfigID, gameServerID err: ", err)
		replyRequestRoomInfo(w, int32(lobby.MsgError_ErrDatabase), nil)
		return
	}

	var roomConfigID = values[0]
	var gameServerID = values[1]
	var groupID = values[2]
	var roomType = values[3]
	var arenaID = values[4]
	// var clubID  = values[2]
	// if clubID != "" &&  !isClubMember(clubID, userID){
	// 	log.Println("onMessageRequestRoomInfo get roomConfigID, gameServerID err: ", err)
	// 	replyJoinClubRoomError(w, int32(MsgError_ErrUserCanNotJoinCLubRoom), clubID)
	// 	return
	// }

	roomConfig, err := redis.String(conn.Do("HGET", gconst.RoomConfigTable, roomConfigID))
	if err != nil {
		log.Println("onMessageRequestRoomInfo get roomConfig err: ", err)
		replyRequestRoomInfo(w, int32(lobby.MsgError_ErrDatabase), nil)
		return
	}

	if isFullRoom(roomID, userID, conn, roomConfig) {
		log.Printf("onMessageRequestRoomInfo room: %s, roomNum:%s, is full", roomID, roomNumber)
		replyRequestRoomInfo(w, int32(lobby.MsgError_ErrRoomIsFull), nil)
		return
	}

	roomConfigJSON := lobby.ParseRoomConfigFromString(roomConfig)
	if groupID != "" && roomConfigJSON.PayType == pay.GroupPay && !isGroupMember(groupID, userID) {
		replyRequestRoomInfo(w, int32(lobby.MsgError_ErrUserCanNotJoinCLubRoom), nil)
		return
	}

	var gameServerURL = getGameServerURL(gameServerID)
	roomTypeInt, _ := strconv.Atoi(roomType)

	var roomInfo = &lobby.RoomInfo{}
	roomInfo.RoomID = &roomID
	roomInfo.RoomNumber = &roomNumber
	roomInfo.GameServerURL = &gameServerURL
	roomInfo.Config = &roomConfig
	roomInfo.ArenaID = &arenaID

	var propCfg = getPropCfg(roomTypeInt)
	roomInfo.PropCfg = &propCfg

	log.Printf("handlerRequestRoomInfo, userID: %s, roomNumber:%s, roomID:%s, gameServerURL:%s", userID, roomNumber, roomID, gameServerURL)

	replyRequestRoomInfo(w, int32(lobby.MsgError_ErrSuccess), roomInfo)
}
