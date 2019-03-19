package sessions

import (
	"encoding/json"
	gconst "gconst"
	"lobbyserver/lobby"
	"lobbyserver/lobby/room"

	log "github.com/sirupsen/logrus"

	"github.com/garyburd/redigo/redis"

	"github.com/golang/protobuf/proto"
)

// onMessageChat 处理聊天消息
func onMessageChat(user *User, accessoryMessage *lobby.AccessoryMessage) {
	log.Println("onMessageChat, userID:", user.uID)
	chatMsg := &lobby.MsgChat{}
	err := proto.Unmarshal(accessoryMessage.Data, chatMsg)
	if err != nil {
		log.Println("onMessageChat decode failed:", err)
		return
	}

	// 过滤文本
	if chatMsg.GetDataType() == int32(lobby.ChatDataType_Text) {
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
		// log.Println("msgContent:", msgContent.Msg)
		// if isContainSensitiveWord(msgContent.Msg) {
		// 	log.Println("Chat message content contain sensitive word")
		// 	return
		// }
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

	var scope = lobby.ChatScopeType(chatMsg.GetScope())
	switch scope {
	case lobby.ChatScopeType_UniCast:
		to := chatMsg.GetTo()
		if to == "" {
			log.Println("unicast chat must supply target")
			return
		}

		var toUser = userMgr.getUserByID(to)
		if toUser == nil {
			log.Println("chat error, target user not exists or is offline")
			return
		}

		chatMsg.From = &user.uID
		toUser.sendMsg(chatMsg, int32(lobby.MessageCode_OPChat))

		// 给自己也发一份，以便客户端更新聊天界面
		user.sendMsg(chatMsg, int32(lobby.MessageCode_OPChat))
		break
	case lobby.ChatScopeType_InRoom:
		// 先从redis中读取用户当前所在的房间中的所有用户ID
		userIDList := readUserIDListInRoom(user.userID())
		chatMsg.From = &user.uID

		var selfSent = false
		for _, uID := range userIDList {
			var toUser = userMgr.getUserByID(uID)
			if toUser != nil {
				toUser.sendMsg(chatMsg, int32(lobby.MessageCode_OPChat))

				if uID == user.uID {
					selfSent = true
				}
			}
		}

		// 如果碰巧列表中没有自己，那就给自己发一份
		if !selfSent {
			user.sendMsg(chatMsg, int32(lobby.MessageCode_OPChat))
		}
		break
	default:
		log.Println("not support chat scope:", scope)
		break
	}
}

func readUserIDListInRoom(who string) []string {
	// 获取redis链接，并退出函数时释放
	conn := lobby.Pool().Get()
	defer conn.Close()

	// 首先读取who所在的房间ID
	roomID := room.LoadUserLastEnterRoomID(who)
	if roomID == "" {
		log.Println("readUserIDListInRoom, get user last room failed:")
		return []string{}
	}

	// 接着读取房间内的用户ID列表
	buf, err := redis.Bytes(conn.Do("HGET", gconst.GsRoomTablePrefix+roomID, "players"))
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
