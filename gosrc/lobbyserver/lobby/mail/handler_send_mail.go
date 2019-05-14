package mail

import (
	"encoding/json"
	"gconst"
	"io/ioutil"
	"lobbyserver/lobby"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"time"
	"fmt"

)

// SendMail 发送邮件
type SendMail struct {
	Mail *lobby.MsgMail `json:"mail"`
	IsAll bool `json:"isAll"`
	Users []string `json:"users"`
}

func saveMail(sendMail *SendMail) {
	// 获取redis链接，并退出函数时释放
	conn := lobby.Pool().Get()
	defer conn.Close()

	mailID := sendMail.Mail.GetId()

	buf, err := proto.Marshal(sendMail.Mail)
	if err != nil {
		return
	}

	conn.Send("MULTI")

	if sendMail.IsAll {
		// TODO: 拉取所有用户，然后把消息保存给所有用户

	} else {
		for _, userID := range sendMail.Users {
			conn.Send("HSET", gconst.LobbyMailPrefix+userID, mailID, buf)
			// sort set
			conn.Send("ZADD", gconst.LobbyMailSortSetPrefix+userID, time.Now().Unix(), mailID)
		}
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		log.Println("saveChatMsg err: ", err)
	}
}

// onMessageChat 处理聊天消息
func handlerSendMail(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// userID := ps.ByName("userID")
	log.Println("handlerSendMail")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("handlerSendMail error:", err)
		return
	}

	sendMail := &SendMail{}
	err = json.Unmarshal(body, sendMail)
	if err != nil {
		log.Println("handlerSendMail decode failed:", err)
		return
	}

	log.Println("sendMail:", sendMail)

	uid, _ := uuid.NewV4()
	mailID := fmt.Sprintf("%v", uid)
	timeStamp := time.Now().Unix()

	sendMail.Mail.Id = &mailID
	sendMail.Mail.TimeStamp = &timeStamp

	if sendMail.Mail.Attachments != nil {
		isReceive := false
		sendMail.Mail.Attachments.IsReceive = &isReceive
	}

	sessionMgr := lobby.SessionMgr()

	if sendMail.IsAll {
		// 发给所有人
		bytes, err := proto.Marshal(sendMail.Mail)
		if err != nil {
			log.Panic("sendMail, marshal msg failed")
			return
		}
		ops := int32(lobby.MessageCode_OPMail)
		lobbyMessage := &lobby.LobbyMessage{}
		lobbyMessage.Ops = &ops
		lobbyMessage.Data = bytes

		bytes, err = proto.Marshal(lobbyMessage)
		if err != nil {
			log.Panic("sendMail, marshal msg failed")
			return
		}

		sessionMgr.Broacast(bytes)
	} else {
		// 发给指定用户
		for _, userID := range sendMail.Users {
			ok := sessionMgr.SendProtoMsgTo(userID, sendMail.Mail, int32(lobby.MessageCode_OPMail))
			if !ok {
				log.Printf("handlerSendMail, send msg to %s failed, target user not exists or is offline", userID)
			}

		}
	}

	saveMail(sendMail)

	w.Write([]byte("ok"))
}