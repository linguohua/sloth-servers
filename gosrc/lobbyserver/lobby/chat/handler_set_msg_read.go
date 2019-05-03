package chat

import (
	// "gconst"
	// "lobbyserver/lobby"
	log "github.com/sirupsen/logrus"
	// "github.com/garyburd/redigo/redis"
	// "github.com/golang/protobuf/proto"
	"net/http"
	// "strconv"
)

// onMessageChat 处理聊天消息
func handlerSetMsgRead(w http.ResponseWriter, r *http.Request, userID string) {
	log.Println("handlerChat, userID:", userID)

	// body, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	log.Println("handlerChat error:", err)
	// 	return
	// }

	// msgs, cursorInt32 := loadUnreadMsg(userID, cursorInt)

	// reply := &lobby.MsgLoadUnreadChatReply{}
	// reply.Msgs = msgs
	// reply.Cursor = &cursorInt32

	// buf, err := proto.Marshal(reply)
	// if err != nil {
	// 	log.Error("handlerLoadUnreadMsg, err:", err)
	// 	return
	// }

	// w.Write(buf)
}