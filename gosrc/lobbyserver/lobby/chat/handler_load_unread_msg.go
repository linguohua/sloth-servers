package chat

import (
	"gconst"
	"lobbyserver/lobby"
	log "github.com/sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	"net/http"
	"strconv"
)

func loadUnreadMsg(userID string, cursor int)([]*lobby.MsgChat, int32) {
	chatMsgs := make([]*lobby.MsgChat, 0)

	conn := lobby.Pool().Get()
	defer conn.Close()

	vs, err := redis.Values(conn.Send("HSCAN", gconst.LobbyChatMessagePrefix+userID, cursor), nil)
	if err != nil {
		log.Error("loadUnreadMsg error:", err)

		return chatMsgs, 0
	}

	nexCursor, _ := redis.Int(vs[0], nil)

	vs, _ = redis.Values(vs[1], nil)

	for i := 1; i < len(vs); i = i+ 2 {
		buf, _:= redis.Bytes(vs[i], nil)

		chatMsg := &lobby.MsgChat{}
		err := proto.Unmarshal(buf, chatMsg)
		if err != nil {
			log.Error("error:", err)
			continue
		}

		chatMsgs = append(chatMsgs, chatMsg)

	 }

	return chatMsgs, int32(nexCursor)
}

// onMessageChat 处理聊天消息
func handlerLoadUnreadMsg(w http.ResponseWriter, r *http.Request, userID string) {
	log.Println("handlerChat, userID:", userID)

	cursor := r.URL.Query().Get("cursor")

	cursorInt, _:= strconv.Atoi(cursor)

	msgs, cursorInt32 := loadUnreadMsg(userID, cursorInt)

	reply := &lobby.MsgLoadUnreadChatReply{}
	reply.Msgs = msgs
	reply.Cursor = &cursorInt32

	buf, err := proto.Marshal(reply)
	if err != nil {
		log.Error("handlerLoadUnreadMsg, err:", err)
		return
	}

	w.Write(buf)
}