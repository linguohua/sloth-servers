package mail

import (
	"gconst"
	"lobbyserver/lobby"
	log "github.com/sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func loadMail(userID string, cursor int, count int)([]*lobby.MsgMail, int32) {
	mails := make([]*lobby.MsgMail, 0)

	conn := lobby.Pool().Get()
	defer conn.Close()

	mailIDs, err := redis.Strings(conn.Do("ZREVRANGE", gconst.LobbyMailSortSetPrefix+userID, cursor, cursor + count - 1))
	if err != nil {
		log.Error("loadUnreadMsg, get mail ids error:", err)

		return mails, 0
	}

	conn.Send("MULTI")
	for _, id := range mailIDs {
		conn.Send("HGET", gconst.LobbyMailPrefix+userID, id)
	}

	vs, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		log.Error("loadUnreadMsg, get mail error:", err)

		return mails, 0
	}

	log.Println("vs length:", len(vs))
	for _, v := range vs {
		buf, err := redis.Bytes(v, nil)
		if err != nil {
			log.Error("loadUnreadMsg, get mail buf error:", err)
			continue
		}

		msgMail := &lobby.MsgMail{}
		err = proto.Unmarshal(buf, msgMail)
		if err != nil {
			log.Error("loadUnreadMsg, unmarshal mail buf error:", err)
			continue
		}

		mails = append(mails, msgMail)

	}
	// _, err = conn.Do("EXEC")
	// if err != nil {
	// 	log.Println("saveChatMsg err: ", err)
	// }

	nexCursor := 0
	if len(mailIDs) == count {
		nexCursor = cursor + count
	}

	return mails, int32(nexCursor)
}

// 拉取邮件
func handlerLoadMail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("userID")

	log.Println("handlerLoadMail, userID:", userID)

	cursor := r.URL.Query().Get("cursor")
	count := r.URL.Query().Get("count")

	cursorInt, _:= strconv.Atoi(cursor)
	countInt, _:= strconv.Atoi(count)
	if countInt == 0 {
		 countInt = 10
	}

	mails, cursorInt32 := loadMail(userID, cursorInt, countInt)
	reply := &lobby.MsgLoadMail{}
	reply.Mails = mails
	reply.Cursor = &cursorInt32

	buf, err := proto.Marshal(reply)
	if err != nil {
		log.Error("handlerLoadUnreadMsg, err:", err)
		return
	}

	log.Println("handlerLoadMail:, buf:", buf)
	w.Write(buf)
}