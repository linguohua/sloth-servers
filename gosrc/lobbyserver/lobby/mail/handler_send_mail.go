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
	"github.com/garyburd/redigo/redis"
	"time"
	"fmt"

)

// SendMail 发送邮件
type SendMail struct {
	Mail *lobby.MsgMail `json:"mail"`
	IsAll bool `json:"isAll"`
	Users []string `json:"users"`
}

func loadAllUserID()[]string {
	conn := lobby.Pool().Get()
	defer conn.Close()

	userIDs, err := redis.Strings(conn.Do("SMEMBERS", gconst.LobbyUserSet))
	if err != nil {
		log.Error("loadAllUserID error:", err)
	}

	return userIDs
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

	userIDs := sendMail.Users
	if sendMail.IsAll {
		userIDs = loadAllUserID()
	}

	conn.Send("MULTI")
	for _, userID := range userIDs {
		conn.Send("HSET", gconst.LobbyMailPrefix+userID, mailID, buf)
		// sort set
		conn.Send("ZADD", gconst.LobbyMailSortSetPrefix+userID, time.Now().Unix(), mailID)
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		log.Println("saveChatMsg err: ", err)
	}

	// 确保每个用户只保留50条
	for _, userID := range userIDs {
		mailIDs, err := redis.Strings(conn.Do("ZREVRANGE", gconst.LobbyMailSortSetPrefix+userID, 50, -1))
		if err != nil {
			log.Error("saveMail, get mail ids error:", err)
			continue
		}

		conn.Send("MULTI")

		for _, mailID := range mailIDs {
			conn.Send("HDEL", gconst.LobbyMailPrefix+userID, mailID)
			conn.Send("ZREM", gconst.LobbyMailSortSetPrefix+userID, mailID)
		}
		_, err = conn.Do("EXEC")
		if err != nil {
			log.Error("saveMail, remov mail error:", err)
		}

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

	saveMail(sendMail)

	sessionMgr := lobby.SessionMgr()

	if sendMail.IsAll {
		// 发给所有人
		ops := int32(lobby.MessageCode_OPMail)
		lobbyMessage := &lobby.LobbyMessage{}
		lobbyMessage.Ops = &ops

		bytes, err := proto.Marshal(lobbyMessage)
		if err != nil {
			log.Panic("sendMail, marshal msg failed")
			return
		}

		sessionMgr.Broacast(bytes)
	} else {
		// 发给指定用户
		for _, userID := range sendMail.Users {
			ok := sessionMgr.SendProtoMsgTo(userID, nil, int32(lobby.MessageCode_OPMail))
			if !ok {
				log.Printf("handlerSendMail, send msg to %s failed, target user not exists or is offline", userID)
			}
		}
	}

	w.Write([]byte("ok"))
}

func sendMail(userID string, content string, title string) {
	uid, _ := uuid.NewV4()
	mailID := fmt.Sprintf("%v", uid)
	timeStamp := time.Now().Unix()

	mail :=  &lobby.MsgMail{}
	mail.Id = &mailID
	mail.TimeStamp = &timeStamp
	mail.Title = &title
	mail.Content = &content
	isRead := false
	mail.IsRead = &isRead

	sendMail := &SendMail{}
	sendMail.Mail = mail
	sendMail.IsAll = false
	sendMail.Users = []string {userID}

	saveMail(sendMail)

	sessionMgr := lobby.SessionMgr()
	ok := sessionMgr.SendProtoMsgTo(userID, nil, int32(lobby.MessageCode_OPMail))
	if !ok {
		log.Printf("handlerSendMail, send msg to %s failed, target user not exists or is offline", userID)
	}
}