package club

import (
	"net/http"
	"gconst"
	"strconv"
	"strings"
	"lobbyserver/lobby"

	"github.com/garyburd/redigo/redis"
	"github.com/julienschmidt/httprouter"
	proto "github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

// onLoadEvents 加载事件
func onLoadMyApplyEvent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("userID")

	log.Println("onLoadMyApplyEvent, userID:", userID)

	var query = r.URL.Query()
	// 俱乐部ID

	cursorStr := query.Get("cursor")
	cursor := 0
	if cursorStr != "" {
		cursor, _ = strconv.Atoi(cursorStr)
	}

	// 获得redis连接
	conn := lobby.Pool().Get()
	defer conn.Close()

	const maxLoad int = 10

	// 先加载ID
	idStrings, err := redis.Strings(conn.Do("LRANGE", gconst.LobbyClubUserApplicantEventPrefix+userID, cursor,
		cursor+maxLoad-1))
	if err != nil && err != redis.ErrNil {
		log.Println("onLoadMyApplyEvent, redis error:", err)
		sendGenericError(w, ClubOperError_CERR_Database_IO)
		return
	}

	events := make([]*MsgClubEvent, 0, maxLoad)
	if err == nil {
		events = loadEventsByIDs(idStrings, conn, userID)
	}

	loadEventReply := &MsgClubLoadEventsReply{}
	loadEventReply.Events = events
	cursor32 := int32(cursor + len(idStrings))
	loadEventReply.Cursor = &cursor32

	b, err := proto.Marshal(loadEventReply)
	if err != nil {
		log.Println("onLoadMyApplyEvent, marshal error:", err)
		sendGenericError(w, ClubOperError_CERR_Encode_Decode)
		return
	}

	sendMsgClubReply(w, ClubReplyCode_RCOperation, b)
}

func loadEventsByIDs(idStrings []string, conn redis.Conn, userID string) []*MsgClubEvent {
	conn.Send("MULTI")
	for _, idString := range idStrings {
		ids := strings.Split(idString, ",")
		clubID := ids[0]
		eventID := ids[1]
		log.Printf("clubID:%s, eventID:%s", clubID, eventID)
		conn.Send("HGET", gconst.LobbyClubEventTablePrefix+clubID, eventID)
	}

	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		log.Panicln("loadEventsByIDs failed:", err)
	}

	events := make([]*MsgClubEvent, 0, len(values))
	for _, v := range values {
		b, err := redis.Bytes(v, nil)
		if err != nil {
			log.Println("loadEventsByIDs, convert value to bytes failed:", err)
			continue
		}

		e := &MsgClubEvent{}
		err = proto.Unmarshal(b, e)
		if err != nil {
			log.Println("loadEventsByIDs, unmarshal bytes to event failed:", err)
			continue
		}

		events = append(events, e)
	}

	return events
}
