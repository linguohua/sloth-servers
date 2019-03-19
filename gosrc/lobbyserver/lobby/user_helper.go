package lobby

import (
	"gconst"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/golang/protobuf/proto"
)

// ReplyHTTPWithProto reply http
func ReplyHTTPWithProto(w http.ResponseWriter, pb proto.Message, ops int32) {
	accessoryMessage := &AccessoryMessage{}
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

// ParseAccessoryMessage parse message
func ParseAccessoryMessage(r *http.Request) (accMsg *AccessoryMessage, errCode int32) {
	if r.ContentLength < 1 {
		log.Println("parseAccessoryMessage failed, content length is zero")
		errCode = int32(MsgError_ErrRequestInvalidParam)
		return
	}

	message := make([]byte, r.ContentLength)
	n, _ := r.Body.Read(message)
	if n != int(r.ContentLength) {
		log.Println("parseAccessoryMessage failed, can't read request body")
		errCode = int32(MsgError_ErrRequestInvalidParam)
		return
	}

	accessoryMessage := &AccessoryMessage{}
	err := proto.Unmarshal(message, accessoryMessage)
	if err != nil {
		log.Println("parseAccessoryMessage failed, Unmarshal msg error:", err)
		errCode = int32(MsgError_ErrDecode)
		return
	}

	accMsg = accessoryMessage
	errCode = int32(MsgError_ErrSuccess)
	return
}

func converGameServerErrCode2AccServerErrCode(gameServerErrCode int32) int32 {
	var errCode = gameServerErrCode
	if errCode == int32(gconst.SSMsgError_ErrEncode) {
		errCode = int32(MsgError_ErrEncode)
	} else if errCode == int32(gconst.SSMsgError_ErrDecode) {
		errCode = int32(MsgError_ErrDecode)
	} else if errCode == int32(gconst.SSMsgError_ErrRoomExist) {
		errCode = int32(MsgError_ErrGameServerRoomExist)
	} else if errCode == int32(gconst.SSMsgError_ErrNoRoomConfig) {
		errCode = int32(MsgError_ErrGameServerNoRoomConfig)
	} else if errCode == int32(gconst.SSMsgError_ErrUnsupportRoomType) {
		errCode = int32(MsgError_ErrGameServerUnsupportRoomType)
	} else if errCode == int32(gconst.SSMsgError_ErrDecodeRoomConfig) {
		errCode = int32(MsgError_ErrGameServerDecodeRoomConfig)
	} else if errCode == int32(gconst.SSMsgError_ErrRoomNotExist) {
		errCode = int32(MsgError_ErrGameServerRoomNotExist)
	}

	return errCode

}
