package room

import (
	"gconst"
	"gpubsub"
	"lobbyserver/config"
	"lobbyserver/lobby"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"encoding/json"

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

	roomNumberString, err := redis.String(conn.Do("HGET", gconst.LobbyRoomTablePrefix+roomID, "roomNumber"))
	if err != nil {
		log.Println("deleteRoomInfoFromRedis, error:", err)
		return
	}

	conn.Send("MULTI")
	conn.Send("DEL", gconst.LobbyRoomTablePrefix+roomID)
	conn.Send("DEL", gconst.LobbyRoomNumberTablePrefix+roomNumberString)
	conn.Send("HDEL", gconst.LobbyUserTablePrefix+userIDString, "roomID")
	conn.Send("SREM", gconst.LobbyRoomTableSet, roomID)

	_, err = conn.Do("EXEC")
	if err != nil {
		log.Println("deleteRoomInfoFromRedis err:", err)
		return
	}
}

func handlerDeleteRoom(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(lobby.ContextKey("userID")).(string)
	roomID := r.URL.Query().Get("roomID")
	if roomID == "" {
		var errCode = int32(lobby.MsgError_ErrRoomIDIsEmpty)
		replyDeleteRoom(w, errCode, lobby.ErrorString[errCode])
		return
	}

	log.Printf("handlerDeleteRoom, userID:%s, roomID:%s", userID, roomID)

	conn := lobby.Pool().Get()
	defer conn.Close()

	exist, err := redis.Int(conn.Do("HEXISTS", gconst.LobbyRoomTablePrefix+roomID))
	if exist == 0 {
		var errCode = int32(lobby.MsgError_ErrRoomNotExist)
		replyDeleteRoom(w, errCode, lobby.ErrorString[errCode])
		return
	}

	fields, err := redis.Strings(conn.Do("HMGET", gconst.LobbyRoomTablePrefix+roomID, "ownerID", "roomType"))
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
	var gameServerID = loadLatestGameServer(int(roomType))

	succeed, msgBagReply := gpubsub.SendAndWait(gameServerID, msgBag, time.Second)

	if succeed {
		errCode := msgBagReply.GetStatus()
		if errCode != 0 {

			errCode = converGameServerErrCode2AccServerErrCode(errCode)
			replyDeleteRoom(w, errCode, lobby.ErrorString[errCode])
			return
		}

		// TODO: 需要退钱给用户
		deleteRoomInfoFromRedis(roomID, ownerID)

		replyDeleteRoom(w, int32(lobby.MsgError_ErrSuccess), "ok")
	} else {
		var errCode = int32(lobby.MsgError_ErrRequestGameServerTimeOut)
		replyDeleteRoom(w, errCode, lobby.ErrorString[errCode])
	}
}
