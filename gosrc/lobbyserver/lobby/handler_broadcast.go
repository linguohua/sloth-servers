package lobby

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"

	"github.com/golang/protobuf/proto"
)

func handleBroadCast(w http.ResponseWriter, r *http.Request) {
	log.Println("handleBroadCast")
	var cmdString = r.URL.Query().Get("cmd")
	cmd, err := strconv.Atoi(cmdString)
	if err != nil {
		var msg = fmt.Sprintf("Can't parser param %s", cmdString)
		w.WriteHeader(404)
		w.Write([]byte(msg))
		return
	}

	if cmd != int(MessageCode_OPBroadcast) {
		var msg = fmt.Sprintf("cmd %d not broadcast message !", cmd)
		w.WriteHeader(404)
		w.Write([]byte(msg))
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

	var msgCode = int32(MessageCode_OPBroadcast)
	accessoryMessage := &AccessoryMessage{}
	accessoryMessage.Ops = &msgCode
	accessoryMessage.Data = message

	bytes, err := proto.Marshal(accessoryMessage)
	if err != nil {
		log.Panic("sendMsg, marshal msg failed")
		return
	}

	users := userMgr.users
	for _, user := range users {
		user.send(bytes)
	}

	w.Write([]byte("Broadcast message success!"))
	return

}

// 用户推送
func handlePlayerPush(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != "POST" {
		msg := fmt.Sprintf("req metod not post!")
		log.Println(msg)
		w.WriteHeader(404)
		w.Write([]byte(msg))
		return
	}

	vals := r.URL.Query()
	playerID := vals.Get("playerID")
	cmdString := vals.Get("cmd")
	cmd, err := strconv.Atoi(cmdString)
	if err != nil {
		msg := fmt.Sprintf("Can't parser param %s", cmdString)
		w.WriteHeader(404)
		w.Write([]byte(msg))
		return
	}

	var data []byte
	if r.ContentLength > 0 {
		message := make([]byte, r.ContentLength)
		n, _ := r.Body.Read(message)
		if n != int(r.ContentLength) {
			w.WriteHeader(404)
			w.Write([]byte("Read message not match origin lenght"))
			return
		}
		data = message
	}

	msgCode := int32(cmd)
	accessoryMessage := &AccessoryMessage{}
	accessoryMessage.Ops = &msgCode
	accessoryMessage.Data = data

	bytes, err := proto.Marshal(accessoryMessage)
	if err != nil {
		log.Panic("sendMsg, marshal msg failed")
		return
	}

	user := userMgr.getUserByID(playerID)
	if user != nil {
		user.send(bytes)
		msg := fmt.Sprintf("player %s handlePlayerPush", playerID)
		log.Println(msg)
	}

	w.Write([]byte("Broadcast message success!"))
	return
}

// 给所有用户广播
func handleAllPlayerPush(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != "POST" {
		msg := fmt.Sprintf("req method not post!")
		log.Println(msg)
		w.WriteHeader(404)
		w.Write([]byte(msg))
		return
	}

	cmdString := r.URL.Query().Get("cmd")
	cmd, err := strconv.Atoi(cmdString)
	if err != nil {
		msg := fmt.Sprintf("Can't parser param %s", cmdString)
		w.WriteHeader(404)
		w.Write([]byte(msg))
		return
	}

	var data []byte
	if r.ContentLength > 0 {
		message := make([]byte, r.ContentLength)
		n, _ := r.Body.Read(message)
		if n != int(r.ContentLength) {
			w.WriteHeader(404)
			w.Write([]byte("Read message not match origin lenght"))
			return
		}

		data = message
	}

	msgCode := int32(cmd)
	accessoryMessage := &AccessoryMessage{}
	accessoryMessage.Ops = &msgCode
	accessoryMessage.Data = data

	bytes, err := proto.Marshal(accessoryMessage)
	if err != nil {
		log.Panic("sendMsg, marshal msg failed")
		return
	}

	users := userMgr.users
	for _, user := range users {
		user.send(bytes)
	}

	w.Write([]byte("Broadcast message success!"))

	return
}
