package chat

import (
	"encoding/json"
	gconst "gconst"
	"lobbyserver/lobby"
	log "github.com/sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	"net/http"
	"io/ioutil"
	"fmt"
)

func saveChatMsg(chatMsg *lobby.MsgChat, userIds []string) {
	// 获取redis链接，并退出函数时释放
	conn := lobby.Pool().Get()
	defer conn.Close()

	uid, _ := uuid.NewV4()
	msgID := fmt.Sprintf("%v",uid)

	buf, err := proto.Marshal(chatMsg)
	if err != nil {
		return
	}

	conn.Send("MULTI")
	for _, userID := range userIds {
		conn.Send("HSET", gconst.LobbyChatMessagePrefix+userID, msgID, buf)
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		log.Println("saveChatMsg err: ", err)
	}
}

func filterSensitiveWord(chatMsg *lobby.MsgChat) {
	// 从消息体中取出文本聊天消息
	type MsgContent struct {
		Msg      string `json:"msg"`
		Index    int    `json:"index"`
		URL      string `json:"url"`
		NickName string `json:"nickname"`
		Sex      string `json:"sex"`
	}

	var msgContent = MsgContent{}
	json.Unmarshal(chatMsg.GetData(), &msgContent)

	isReplace, newWord := replaceSensitiveWord(msgContent.Msg, "*")
	if isReplace {
		if newWord == "" {
			return
		}

		msgContent.Msg = newWord
		buf, err := json.Marshal(msgContent)
		if err != nil {
			log.Println("Marshal chat msg content err:", msgContent)
		} else {
			chatMsg.Data = buf
		}
	}
}

// onMessageChat 处理聊天消息
func handlerChat(w http.ResponseWriter, r *http.Request, userID string) {
	log.Println("handlerChat, userID:", userID)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("handlerChat error:", err)
		return
	}


	chatMsg := &lobby.MsgChat{}
	err = proto.Unmarshal(body, chatMsg)
	if err != nil {
		log.Println("onMessageChat decode failed:", err)
		return
	}

	// 过滤文本
	if chatMsg.GetDataType() == int32(lobby.ChatDataType_Text) {
		filterSensitiveWord(chatMsg)

	}

	sessionMgr := lobby.SessionMgr()

	var scope = lobby.ChatScopeType(chatMsg.GetScope())
	switch scope {
	case lobby.ChatScopeType_UniCast:
		to := chatMsg.GetTo()
		if to == "" {
			log.Println("unicast chat must supply target")
			return
		}

		chatMsg.From = &userID
		// 给对方发一份
		ok := sessionMgr.SendProtoMsgTo(to, chatMsg, int32(lobby.MessageCode_OPChat))
		if !ok {
			log.Printf("handlerChat, send msg to %s failed, target user not exists or is offline", to)
		}

		// 给自己也发一份
		ok = sessionMgr.SendProtoMsgTo(userID, chatMsg, int32(lobby.MessageCode_OPChat))
		if !ok {
			log.Printf("handlerChat, send msg to %s failed, target user not exists or is offline", userID)
		}

		break
	case lobby.ChatScopeType_InRoom:
		// 先从redis中读取用户当前所在的房间中的所有用户ID
		userIDList := readUserIDListInRoom(userID)
		chatMsg.From = &userID

		var isIncludeSelf = false
		for _, uID := range userIDList {
			if uID == userID {
				isIncludeSelf = true
				break
			}
		}

		if !isIncludeSelf {
			userIDList = append(userIDList, userID)
		}

		for _, uID := range userIDList {
			ok := sessionMgr.SendProtoMsgTo(uID, chatMsg, int32(lobby.MessageCode_OPChat))
			if !ok {
				log.Printf("handlerChat, send msg to %s failed, target user not exists or is offline", uID)
			}
		}

		saveChatMsg(chatMsg, userIDList)
		break
	default:
		log.Println("not support chat scope:", scope)
		break
	}

	w.Write([]byte("ok"))
}

func readUserIDListInRoom(who string) []string {
	// 获取redis链接，并退出函数时释放
	conn := lobby.Pool().Get()
	defer conn.Close()

	// 首先读取who所在的房间ID
	roomID := lobby.RoomUtil().LoadUserLastEnterRoomID(who)
	if roomID == "" {
		log.Println("readUserIDListInRoom, get user last room failed:")
		return []string{}
	}

	// 接着读取房间内的用户ID列表
	buf, err := redis.Bytes(conn.Do("HGET", gconst.GameServerRoomTablePrefix+roomID, "players"))
	if err != nil {
		log.Println("readUserIDListInRoom, get room players failed:", err)
		return []string{}
	}

	userIDList := &gconst.SSMsgUserIDList{}
	err = proto.Unmarshal(buf, userIDList)
	if err != nil {
		log.Println("readUserIDListInRoom, unmarshal failed:", err)
		return []string{}
	}

	return userIDList.GetUserIDs()
}
