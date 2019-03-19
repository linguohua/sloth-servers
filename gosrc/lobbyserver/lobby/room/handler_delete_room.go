package room

import (
	"gconst"
	"lobbyserver/config"
	"lobbyserver/lobby"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"encoding/json"
	"fmt"

	"github.com/garyburd/redigo/redis"
	proto "github.com/golang/protobuf/proto"
)

const (
	//SRoomIdle 房间空闲
	SRoomIdle = 0
	// SRoomWaiting 房间正在等待玩家进入
	SRoomWaiting = 1
	// SRoomPlaying 游戏正在进行中
	SRoomPlaying = 2
	//SRoomDeleted 房间已经删除
	SRoomDeleted = 3
)

func replyDeleteRoom(w http.ResponseWriter, errCode int32, errMsg string) {
	type DeleteRoomReply struct {
		ErrorCode int    `json:"errorCode"`
		ErrorMsg  string `json:"errorMsg"`
	}

	reply := &DeleteRoomReply{}
	reply.ErrorCode = int(errCode)
	reply.ErrorMsg = errMsg

	b, err := json.Marshal(reply)
	if err != nil {
		log.Panicln("genericReply, json marshal error:", err)
		return
	}

	w.Write(b)
}

func deleteRoomInfoFromRedis(roomID string, userIDString string) {
	log.Printf("deleteRoomInfoFromRedis, roomID:%s, userID:%s", roomID, userIDString)
	// 1. 先删除房间信息
	// 2. 从用户房间列表中删除房间，AA制的房间可能不放在用户的房间列表中
	conn := lobby.Pool().Get()
	defer conn.Close()

	vs, err := redis.Strings(conn.Do("HMGET", gconst.RoomTablePrefix+roomID, "roomNumber", "groupID", "ownerID"))
	if err != nil {
		log.Println("deleteRoomInfoFromRedis, error:", err)
		return
	}

	var roomNumberString = vs[0]
	var groupID = vs[1]
	var ownerID = vs[2]

	conn.Send("MULTI")
	conn.Send("DEL", gconst.RoomTablePrefix+roomID)
	conn.Send("DEL", gconst.RoomNumberTable+roomNumberString)
	conn.Send("HDEL", gconst.AsUserTablePrefix+userIDString, "roomID")
	conn.Send("SREM", gconst.RoomTableACCSet, roomID)

	if groupID != "" {
		conn.Send("SREM", gconst.GroupRoomsSetPrefix+groupID, roomID)

		groupMemberRoomsSetKey := fmt.Sprintf(gconst.GroupMemberRoomsSet, groupID, ownerID)
		conn.Send("SREM", groupMemberRoomsSetKey, roomID)
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		log.Println("deleteRoomInfoFromRedis err:", err)
		return
	}

	if userIDString == "" {
		return
	}

	bytes, err := redis.Bytes(conn.Do("HGET", gconst.AsUserTablePrefix+userIDString, "rooms"))
	if err != nil {
		log.Println("deleteRoomInfoFromRedis, error:", err)
		return
	}

	var roomIDList = &lobby.RoomIDList{}
	err = proto.Unmarshal(bytes, roomIDList)
	if err != nil {
		log.Println("deleteRoomInfoFromRedis, error:", err)
		return
	}

	var roomIDs = roomIDList.GetRoomIDs()
	for i, v := range roomIDs {
		if v == roomID {
			roomIDs = append(roomIDs[:i], roomIDs[i+1:]...)
			break
		}
	}

	roomIDList.RoomIDs = roomIDs
	buf, err := proto.Marshal(roomIDList)
	if err != nil {
		log.Println("deleteRoomInfoFromRedis, error:", err)
		return
	}

	conn.Do("HSET", gconst.AsUserTablePrefix+userIDString, "rooms", buf)

}

func handlerDeleteRoom(w http.ResponseWriter, r *http.Request, userID string) {
	roomID := r.URL.Query().Get("roomID")
	if roomID == "" {
		var errCode = int32(lobby.MsgError_ErrRoomIDIsEmpty)
		replyDeleteRoom(w, errCode, lobby.ErrorString[errCode])
		return
	}

	log.Printf("handlerDeleteRoom, userID:%s, roomID:%s", userID, roomID)

	conn := lobby.Pool().Get()
	defer conn.Close()

	exist, err := redis.Int(conn.Do("HEXISTS", gconst.RoomTablePrefix+roomID))
	if exist == 0 {
		var errCode = int32(lobby.MsgError_ErrRoomNotExist)
		replyDeleteRoom(w, errCode, lobby.ErrorString[errCode])
		return
	}

	fields, err := redis.Strings(conn.Do("HMGET", gconst.RoomTablePrefix+roomID, "ownerID", "roomType"))
	if err == redis.ErrNil {
		var errCode = int32(lobby.MsgError_ErrRoomNotExist)
		replyDeleteRoom(w, errCode, lobby.ErrorString[errCode])
		return
	}

	if err != nil {
		var errCode = int32(lobby.MsgError_ErrDatabase)
		replyDeleteRoom(w, errCode, lobby.ErrorString[errCode])
		return
	}

	//检查房间的拥有者
	var ownerID = fields[0]
	if ownerID != userID {
		log.Printf("onMessageDeleteRoom, %s not room creator,cant delete room, owner is %s", userID, ownerID)
		var errCode = int32(lobby.MsgError_ErrNotRoomCreater)
		replyDeleteRoom(w, errCode, lobby.ErrorString[errCode])
		return
	}

	var roomTypeStr = fields[1]
	roomType, err := strconv.Atoi(roomTypeStr)
	if err != nil {
		var errCode = int32(lobby.MsgError_ErrDecode)
		replyDeleteRoom(w, errCode, lobby.ErrorString[errCode])
		return
	}

	//检查房间的运行状态
	// var stateString = fields[1]
	// state, _ := strconv.ParseInt(stateString, 10, 32)
	// if state == SRoomPlaying {
	// 	log.Printf("onMessageDeleteRoom, game is playing")
	// 	var errCode = int32(MsgError_ErrGameIsPlaying)
	// 	replyDeleteRoom(w, errCode, ErrorString[errCode])
	// 	return
	// }

	//请求游戏服务器删除房间
	var msgDeleteRoom = &gconst.SSMsgDeleteRoom{}
	msgDeleteRoom.RoomID = &roomID

	msgDeleteRoomBuf, err := proto.Marshal(msgDeleteRoom)
	if err != nil {
		log.Println("parse roomConfig err： ", err)
		var errCode = int32(lobby.MsgError_ErrEncode)
		replyDeleteRoom(w, errCode, lobby.ErrorString[errCode])
		return
	}

	msgType := int32(gconst.SSMsgType_Request)
	requestCode := int32(gconst.SSMsgReqCode_DeleteRoom)
	status := int32(gconst.SSMsgError_ErrSuccess)

	msgBag := &gconst.SSMsgBag{}
	msgBag.MsgType = &msgType
	var sn = lobby.GenerateSn()
	msgBag.SeqNO = &sn
	msgBag.RequestCode = &requestCode
	msgBag.Status = &status
	var url = config.ServerID
	msgBag.SourceURL = &url
	msgBag.Params = msgDeleteRoomBuf

	// log.Println("roomType:", roomType)
	var gameServerID = getGameServerID(int(roomType))

	succeed, msgBagReply := lobby.SendAndWait(gameServerID, msgBag, time.Second)

	if succeed {
		errCode := msgBagReply.GetStatus()
		if errCode != 0 {

			errCode = converGameServerErrCode2AccServerErrCode(errCode)
			replyDeleteRoom(w, errCode, lobby.ErrorString[errCode])
			return
		}

		deleteRoomInfoFromRedis(roomID, ownerID)

		replyDeleteRoom(w, int32(lobby.MsgError_ErrSuccess), "ok")
		// order := refund2UserAndSave2Redis(roomID, , 0)
		// if order != nil && order.Refund != nil && (order.Refund.Result == int32(gconst.SSMsgError_ErrSuccess)) {
		// 	user.updateMoney(uint32(order.Refund.RemainDiamond))
		// }

		// var releaseRoomRsp = &MsgReleaseRoomRsp{}
		// errCode = int32(MsgError_ErrSuccess)
		// releaseRoomRsp.Result = &errCode
		// var errString = ErrorString[errCode]
		// releaseRoomRsp.RetMsg = &errString
		// releaseRoomRsp.RoomID = &roomID
		// user.sendMsg(releaseRoomRsp, int32(MessageCode_OPDeleteRoom))
	} else {
		var errCode = int32(lobby.MsgError_ErrRequestGameServerTimeOut)
		replyDeleteRoom(w, errCode, lobby.ErrorString[errCode])
	}
}
