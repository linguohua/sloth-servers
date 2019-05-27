package room

import (
	"gconst"
	"gpubsub"
	"lobbyserver/config"
	"lobbyserver/lobby"
	"net/http"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/julienschmidt/httprouter"

	log "github.com/sirupsen/logrus"

	"github.com/golang/protobuf/proto"
)

func deleteRoomReply(w http.ResponseWriter, errCode int32) {
	deleteRoomReply := &lobby.MsgDeleteRoomReply{}
	deleteRoomReply.Result = &errCode
	bytes, err := proto.Marshal(deleteRoomReply)
	if err != nil {
		log.Panic("deleteRoomReply, marshal msg failed")
		return
	}

	w.Write(bytes)
}

func deleteClubRoomInfoFromRedis(roomID string, clubID string) {
	log.Printf("deleteClubRoomInfoFromRedis, roomID:%s, clubID:%s", roomID, clubID)
	// 1. 先删除房间信息
	// 2. 从用户房间列表中删除房间，AA制的房间可能不放在用户的房间列表中
	conn := lobby.Pool().Get()
	defer conn.Close()

	roomNumberString, err := redis.String(conn.Do("HGET", gconst.LobbyRoomTablePrefix+roomID, "roomNumber"))
	if err != nil {
		log.Error("deleteClubRoomInfoFromRedis, error:", err)
		return
	}

	conn.Send("MULTI")
	conn.Send("DEL", gconst.LobbyRoomTablePrefix+roomID)
	conn.Send("DEL", gconst.LobbyRoomNumberTablePrefix+roomNumberString)
	conn.Send("SREM", gconst.LobbyClubRoomSetPrefix+clubID, roomID)
	conn.Send("SREM", gconst.LobbyRoomTableSet, roomID)

	_, err = conn.Do("EXEC")
	if err != nil {
		log.Println("deleteRoomInfoFromRedis err:", err)
		return
	}
}

func forceDeleteRoom(roomID string) (errCode int32) {
	conn := lobby.Pool().Get()
	defer conn.Close()

	fields, err := redis.Strings(conn.Do("HMGET", gconst.LobbyRoomTablePrefix+roomID, "ownerID", "roomType"))
	if err != nil {
		log.Printf("handlerDeleteRoomForClub, get room %s ownerID, roomType from redis err:%v", roomID, err)
		return int32(lobby.MsgError_ErrDatabase)
	}

	var ownerID = fields[0]
	var roomTypeStr = fields[1]
	roomType, err := strconv.Atoi(roomTypeStr)
	if err != nil {
		log.Error("handlerDeleteRoomForClub, Convert roomTypeStr error:", err)
	}

	if ownerID == "" && roomTypeStr == "" {
		return int32(lobby.MsgError_ErrRoomNotExist)
	}

	//请求游戏服务器删除房间
	var msgDeleteRoom = &gconst.SSMsgDeleteRoom{}
	msgDeleteRoom.RoomID = &roomID

	msgDeleteRoomBuf, err := proto.Marshal(msgDeleteRoom)
	if err != nil {
		log.Println("parse roomConfig err： ", err)
		return int32(lobby.MsgError_ErrEncode)
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
			return errCode
		}

		// TODO: 需要退钱给用户
		deleteClubRoomInfoFromRedis(roomID, ownerID)

		return int32(lobby.MsgError_ErrSuccess)
	}

	return int32(lobby.MsgError_ErrRequestGameServerTimeOut)
}

func handlerDeleteClubRoom(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("userID")
	clubID := r.URL.Query().Get("clubID")
	roomID := r.URL.Query().Get("roomID")

	log.Printf("handlerDeleteRoomForClub call, userID:%s, clubID:%s, roomID:%s", userID, clubID, roomID)

	if clubID == "" {
		log.Println("handlerDeleteRoomForClub, need clubID")
		deleteRoomReply(w, int32(lobby.MsgError_ErrRequestInvalidParam))
		return
	}

	if roomID == "" {
		log.Println("handlerDeleteRoomForClub, need roomID")
		deleteRoomReply(w, int32(lobby.MsgError_ErrRequestInvalidParam))
		return
	}

	// 1. 判断牌友圈是否存在
	// 2. 判断用户是否是管理员或者群主
	clubMgr := lobby.ClubMgr()
	club := clubMgr.GetClub(clubID)
	if club == nil {
		log.Printf("handlerDeleteRoomForClub, no club found for %s", clubID)
		deleteRoomReply(w, int32(lobby.MsgError_ErrRequestInvalidParam))
		return
	}

	if !clubMgr.IsUserPermisionDeleteRoom(userID, clubID) {
		log.Printf("handlerDeleteRoomForClub, user %s not allow delete room in club %s", userID, clubID)
		deleteRoomReply(w, int32(lobby.MsgError_ErrOnlyClubCreatorOrManagerAllowDeleteRoom))
		return
	}

	errCode := forceDeleteRoom(roomID)
	deleteRoomReply(w, errCode)
}
