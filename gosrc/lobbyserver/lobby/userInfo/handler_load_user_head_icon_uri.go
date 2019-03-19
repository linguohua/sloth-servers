package userInfo

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"gconst"
	"strconv"

	"github.com/garyburd/redigo/redis"
	proto "github.com/golang/protobuf/proto"
	"lobbyserver/lobby/pb"
	"lobbyserver/lobby/errorString"
	"lobbyserver/lobby"
)

func replyLoadUserHeadIconSuccess(w http.ResponseWriter, msg *pb.MsgLoadUserHeadIconURIReply) {
	bytes, err := proto.Marshal(msg)
	if err != nil {
		log.Panic("reply msg, marshal msg failed, err:", err)
		return
	}

	w.Write(bytes)
}

func replyLoadUserHeadIconError(w http.ResponseWriter, errorCode int32) {
	var loadUserHeadIconURIReply = &pb.MsgLoadUserHeadIconURIReply{}
	var code = errorCode
	loadUserHeadIconURIReply.Result = &code
	var errString = errorString.ErrorString[errorCode]
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
		replyLoadUserHeadIconError(w, int32(pb.MsgError_ErrRequestInvalidParam))
		return
	}

	message := make([]byte, r.ContentLength)
	n, _ := r.Body.Read(message)
	if n != int(r.ContentLength) {
		log.Println("parseAccessoryMessage failed, can't read request body")
		replyLoadUserHeadIconError(w, int32(pb.MsgError_ErrRequestInvalidParam))
		return
	}

	msg := &pb.MsgLoadUserHeadIconURI{}
	err := proto.Unmarshal(message, msg)
	if err != nil {
		log.Println("onMessageCreateRoom, Unmarshal err:", err)
		replyLoadUserHeadIconError(w, int32(pb.MsgError_ErrDecode))
		return
	}

	userIDs := msg.GetUserIDs()
	if len(userIDs) == 0 {
		var replyMsg = &pb.MsgLoadUserHeadIconURIReply{}
		var errCode = int32(pb.MsgError_ErrSuccess)
		replyMsg.Result = &errCode
		replyLoadUserHeadIconSuccess(w, replyMsg)
		return
	}

	conn := pool.Get()
	defer conn.Close()

	conn.Send("MULTI")

	for _, userID := range userIDs {
		conn.Send("HMGET", gconst.AsUserTablePrefix+userID, "userSex", "userLogo")
	}

	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		log.Println("handleLoadUserHeadIconURI， load user head icon from redis err:", err)
		replyLoadUserHeadIconError(w, int32(MsgError_ErrDatabase))
		return
	}

	var headIconInfos = make([]*MsgHeadIconInfo, 0, len(values))
	for index, v := range values {
		fileds, _ := redis.Strings(v, nil)
		var headIconInfo = &MsgHeadIconInfo{}
		var userID = userIDs[index]
		sexUint64, _ := strconv.ParseUint(fileds[0], 10, 32)
		var sex = uint32(sexUint64)
		headIconInfo.Sex = &sex
		headIconInfo.UserID = &userID
		var headIconURI = fileds[1]
		headIconInfo.HeadIconURI = &headIconURI
		headIconInfos = append(headIconInfos, headIconInfo)
	}
	var replyMsg = &MsgLoadUserHeadIconURIReply{}
	var errCode = int32(MsgError_ErrSuccess)
	replyMsg.Result = &errCode
	replyMsg.HeadIconInfos = headIconInfos
	replyLoadUserHeadIconSuccess(w, replyMsg)
}
