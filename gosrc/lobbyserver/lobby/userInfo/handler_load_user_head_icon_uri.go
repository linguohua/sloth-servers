package userinfo

import (
	"gconst"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"lobbyserver/lobby"

	"github.com/garyburd/redigo/redis"
	proto "github.com/golang/protobuf/proto"
)

func replyLoadUserHeadIconSuccess(w http.ResponseWriter, msg *lobby.MsgLoadUserHeadIconURIReply) {
	bytes, err := proto.Marshal(msg)
	if err != nil {
		log.Panic("reply msg, marshal msg failed, err:", err)
		return
	}

	w.Write(bytes)
}

func replyLoadUserHeadIconError(w http.ResponseWriter, errorCode int32) {
	var loadUserHeadIconURIReply = &lobby.MsgLoadUserHeadIconURIReply{}
	var code = errorCode
	loadUserHeadIconURIReply.Result = &code
	var errString = lobby.ErrorString[errorCode]
	loadUserHeadIconURIReply.RetMsg = &errString

	bytes, err := proto.Marshal(loadUserHeadIconURIReply)
	if err != nil {
		log.Panic("reply msg, marshal msg failed, err:", err)
		return
	}

	w.Write(bytes)
}

func handleLoadUserHeadIconURI(w http.ResponseWriter, r *http.Request, userID string) {
	log.Println("handleLoadUserHeadIconURI, useID:", userID)

	if r.ContentLength < 1 {
		log.Println("parseAccessoryMessage failed, content length is zero")
		replyLoadUserHeadIconError(w, int32(lobby.MsgError_ErrRequestInvalidParam))
		return
	}

	message := make([]byte, r.ContentLength)
	n, _ := r.Body.Read(message)
	if n != int(r.ContentLength) {
		log.Println("parseAccessoryMessage failed, can't read request body")
		replyLoadUserHeadIconError(w, int32(lobby.MsgError_ErrRequestInvalidParam))
		return
	}

	msg := &lobby.MsgLoadUserHeadIconURI{}
	err := proto.Unmarshal(message, msg)
	if err != nil {
		log.Println("onMessageCreateRoom, Unmarshal err:", err)
		replyLoadUserHeadIconError(w, int32(lobby.MsgError_ErrDecode))
		return
	}

	userIDs := msg.GetUserIDs()
	if len(userIDs) == 0 {
		var replyMsg = &lobby.MsgLoadUserHeadIconURIReply{}
		var errCode = int32(lobby.MsgError_ErrSuccess)
		replyMsg.Result = &errCode
		replyLoadUserHeadIconSuccess(w, replyMsg)
		return
	}

	conn := lobby.Pool().Get()
	defer conn.Close()

	conn.Send("MULTI")

	for _, userID := range userIDs {
		conn.Send("HMGET", gconst.LobbyUserTablePrefix+userID, "userSex", "userLogo")
	}

	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		log.Println("handleLoadUserHeadIconURI， load user head icon from redis err:", err)
		replyLoadUserHeadIconError(w, int32(lobby.MsgError_ErrDatabase))
		return
	}

	var headIconInfos = make([]*lobby.MsgHeadIconInfo, 0, len(values))
	for index, v := range values {
		fileds, _ := redis.Strings(v, nil)
		var headIconInfo = &lobby.MsgHeadIconInfo{}
		var userID = userIDs[index]
		sexUint64, _ := strconv.ParseUint(fileds[0], 10, 32)
		var sex = uint32(sexUint64)
		headIconInfo.Sex = &sex
		headIconInfo.UserID = &userID
		var headIconURI = fileds[1]
		headIconInfo.HeadIconURI = &headIconURI
		headIconInfos = append(headIconInfos, headIconInfo)
	}
	var replyMsg = &lobby.MsgLoadUserHeadIconURIReply{}
	var errCode = int32(lobby.MsgError_ErrSuccess)
	replyMsg.Result = &errCode
	replyMsg.HeadIconInfos = headIconInfos
	replyLoadUserHeadIconSuccess(w, replyMsg)
}

func reply(w http.ResponseWriter, pb proto.Message, ops int32) {
	accessoryMessage := &lobby.AccessoryMessage{}
	accessoryMessage.Ops = &ops

	if pb != nil {
		bytes, err := proto.Marshal(pb)

		if err != nil {
			log.Panic("reply msg, marshal msg failed")
			return
		}
		accessoryMessage.Data = bytes
	}

	bytes, err := proto.Marshal(accessoryMessage)
	if err != nil {
		log.Panic("reply msg, marshal msg failed")
		return
	}

	w.Write(bytes)
}
