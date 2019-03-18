package lobby

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/golang/protobuf/proto"
)

func handleUniCast(w http.ResponseWriter, r *http.Request) {
	log.Println("handleUniCast")
	var userID = r.URL.Query().Get("userID")
	if userID == "" {
		w.WriteHeader(404)
		w.Write([]byte("Need user id !"))
		return
	}

	if r.ContentLength < 1 {
		w.WriteHeader(404)
		w.Write([]byte("content is emtpy"))
		return
	}

	message := make([]byte, r.ContentLength)
	n, _ := r.Body.Read(message)
	if n != int(r.ContentLength) {
		w.WriteHeader(404)
		w.Write([]byte("Read message not match origin lenght"))
		return
	}

	user := userMgr.getUserByID(userID)
	if user == nil {
		log.Println("user offline")
		w.Write([]byte("User offline !"))
		return
	}

	var msgCode = int32(MessageCode_OPUnicast)
	accessoryMessage := &AccessoryMessage{}
	accessoryMessage.Ops = &msgCode
	accessoryMessage.Data = message

	bytes, err := proto.Marshal(accessoryMessage)
	if err != nil {
		log.Panic("sendMsg, marshal msg failed")
		return
	}

	user.send(bytes)

	var msg = fmt.Sprintf("Send message to user %s success!", userID)
	w.Write([]byte(msg))
	return

}
