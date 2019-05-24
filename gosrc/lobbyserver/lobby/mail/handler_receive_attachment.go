package mail

import (
	"gconst"
	"lobbyserver/lobby"
	log "github.com/sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	"github.com/julienschmidt/httprouter"
	"net/http"
	// "strconv"
)

// onMessageChat 处理聊天消息
func handlerReceiveAttahment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("userID")
	mailID := r.URL.Query().Get("mailID")

	log.Printf("handlerReceiveAttahment, userID:%s, mailID:%s", userID, mailID)

	conn := lobby.Pool().Get()
	defer conn.Close()

	buf, err := redis.Bytes(conn.Do("HGET", gconst.LobbyMailPrefix+userID, mailID))
	if err != nil {
		log.Error("handlerReceiveAttahment, read redis error:", err)
	}

	mail := &lobby.MsgMail{}
	err = proto.Unmarshal(buf, mail)
	if err != nil {
		log.Error("handlerReceiveAttahment, unmarshal error:", err)
		return
	}

	if mail.Attachments.GetIsReceive() {
		log.Error("handlerReceiveAttahment, has been receive attachment")
		return
	}

	attachmentNum := mail.Attachments.GetNum()
	if attachmentNum < 0 {
		log.Error("handlerReceiveAttahment, receive attachment num can not less than 0")
		return
	}

	isReceive := true
	mail.Attachments.IsReceive = &isReceive



	// 更新用户附件
	if mail.Attachments.GetType() == int32(lobby.MailAttachmentType_Diamond) {
		mySQLUtil := lobby.MySQLUtil()
		diamond, errCode := mySQLUtil.UpdateDiamond(userID, int64(attachmentNum))
		if errCode != 0 {
			log.Error("handlerReceiveAttahment, update diamond errCode:", errCode)
			return
		}

		conn.Do("HSET", gconst.LobbyUserTablePrefix+userID, "diamond", diamond)

		lobby.UpdateDiamond(userID, uint64(diamond))
	}

	buf, err = proto.Marshal(mail)
	if err != nil {
		log.Error("handlerReceiveAttahment, marshal error:", err)
		return
	}

	_, err = conn.Do("HSET", gconst.LobbyMailPrefix+userID, mailID, buf)
	if err != nil {
		log.Error("handlerReceiveAttahment set mail error:", err)
	}

	w.Write([]byte("ok"))
}