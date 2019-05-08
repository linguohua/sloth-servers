package room

import (
	"gconst"
	"lobbyserver/lobby"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"

	"io/ioutil"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
)

func replyRequestRoomInfo(w http.ResponseWriter, errorCode int32, roomInfo *lobby.RoomInfo) {

	msgRequestRoomInfoRsp := &lobby.MsgRequestRoomInfoRsp{}
	msgRequestRoomInfoRsp.Result = proto.Int32(errorCode)
	var errString = lobby.ErrorString[errorCode]
	msgRequestRoomInfoRsp.RetMsg = &errString
	msgRequestRoomInfoRsp.RoomInfo = roomInfo

	bytes, err := proto.Marshal(msgRequestRoomInfoRsp)
	if err != nil {
		log.Panic("reply msg, marshal msg failed")
		return
	}

	w.Write(bytes)
}

func getModuleCfg(r *http.Request, roomInfo *lobby.RoomInfo) string {
	roomConfigJSON := lobby.ParseRoomConfigFromString(roomInfo.GetConfig())

	query := r.URL.Query()
	isForceUpgrade := query.Get("forceUpgrade")
	if isForceUpgrade != "true" {
		return ""
	}

	// 游戏模块名称
	query.Set("qMod", roomConfigJSON.ModuleName)
	//不知道版本是多少，填最低版本，会获取到最高版本的配置
	r.URL.RawQuery = query.Encode()

	log.Println(" r.URL.Query:", r.URL.Query())
	updatUtil := lobby.UpdateUtil()
	moduleCfg := updatUtil.GetModuleCfg(r)

	log.Println("moduleCfg:", moduleCfg)
	return moduleCfg
}

func isFullRoom(roomID string, userID string, conn redis.Conn, roomConfigString string) bool {
	// 判断房间是否已经满
	buf, err := redis.Bytes(conn.Do("HGET", gconst.GameServerRoomTablePrefix+roomID, "players"))
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

func handlerRequestRoomInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("userID")
	log.Println("handlerRequestRoomInfo call, userID:", userID)

	// 1. 从请求中获取房间6位数字ID
	// 2. 检查房间有效性，比如是否存在，是否已经满了
	// 3. 用房间6位数字ID去请求房间id
	// 4. 获取房间所在服务器的ID
	// 5. 用服务器ID去获取服务器的URL
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("handlerCreateRoom error:", err)
		return
	}

	var msgRequestRoomInfo = &lobby.MsgRequestRoomInfo{}
	err = proto.Unmarshal(body, msgRequestRoomInfo)
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
		// TODO: 获取模块版本更新信息
		moduleCfg := getModuleCfg(r, lastRoomInfo)
		lastRoomInfo.ModuleCfg = &moduleCfg

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

	roomID, err := redis.String(conn.Do("HGET", gconst.LobbyRoomNumberTablePrefix+roomNumber, "roomID"))
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

	values, err := redis.Strings(conn.Do("HMGET", gconst.LobbyRoomTablePrefix+roomID, "roomConfigID", "gameServerID", "roomType"))
	if err != nil {
		log.Println("onMessageRequestRoomInfo get roomConfigID, gameServerID err: ", err)
		replyRequestRoomInfo(w, int32(lobby.MsgError_ErrDatabase), nil)
		return
	}

	var roomConfigID = values[0]
	var gameServerID = values[1]
	var roomType = values[2]

	roomConfig, err := redis.String(conn.Do("HGET", gconst.LobbyRoomConfigTable, roomConfigID))
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

	// roomConfigJSON := lobby.ParseRoomConfigFromString(roomConfig)
	// if groupID != "" && roomConfigJSON.PayType == pay.GroupPay && !isGroupMember(groupID, userID) {
	// 	replyRequestRoomInfo(w, int32(lobby.MsgError_ErrUserCanNotJoinCLubRoom), nil)
	// 	return
	// }

	roomTypeInt, _ := strconv.Atoi(roomType)

	var roomInfo = &lobby.RoomInfo{}
	roomInfo.RoomID = &roomID
	roomInfo.RoomNumber = &roomNumber
	roomInfo.GameServerID = &gameServerID
	roomInfo.Config = &roomConfig

	var propCfg = getPropCfg(roomTypeInt)
	roomInfo.PropCfg = &propCfg

	moduleCfg := getModuleCfg(r, lastRoomInfo)
	lastRoomInfo.ModuleCfg = &moduleCfg

	log.Printf("handlerRequestRoomInfo, userID: %s, roomNumber:%s, roomID:%s, GameServerID:%s", userID, roomNumber, roomID, gameServerID)

	replyRequestRoomInfo(w, int32(lobby.MsgError_ErrSuccess), roomInfo)
}
