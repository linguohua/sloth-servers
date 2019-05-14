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
func handlerSetMsgRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("userID")
	mailID := r.URL.Query().Get("mailID")

	log.Printf("handlerSetMsgRead, userID:%s, mailID:%s", userID, mailID)

	conn := lobby.Pool().Get()
	defer conn.Close()

	buf, err := redis.Bytes(conn.Do("HGET", gconst.LobbyMailPrefix+userID, mailID))
	if err != nil {
		log.Error("handlerSetMsgRead, read redis error:", err)
	}

	mail := &lobby.MsgMail{}
	err = proto.Unmarshal(buf, mail)
	if err != nil {
		log.Error("handlerSetMsgRead, unmarshal error:", err)
		return
	}

	isRead := true
	mail.IsRead = &isRead


	buf, err = proto.Marshal(mail)
	if err != nil {
		log.Error("handlerSetMsgRead, marshal error:", err)
		return
	}

	_, err = conn.Do("HSET", gconst.LobbyMailPrefix+userID, mailID, buf)
	if err != nil {
		log.Error("handlerSetMsgRead set mail error:", err)
	}

	w.Write([]byte("ok"))
}