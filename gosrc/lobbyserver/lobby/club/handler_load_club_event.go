package club

import (
	"net/http"
	"gconst"
	"strconv"
	"lobbyserver/lobby"

	"github.com/garyburd/redigo/redis"
	"github.com/julienschmidt/httprouter"
	proto "github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

// onLoadEvents 加载事件
func onLoadEvents(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("userID")

	log.Println("onLoadEvents, userID:", userID)

	var query = r.URL.Query()
	// 俱乐部ID
	clubID := query.Get("clubID")
	if clubID == "" {
		log.Println("onLoadEvents, need clubID")
		sendGenericError(w, ClubOperError_CERR_Invalid_Input_Parameter)
		return
	}

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
	eventIDs, err := redis.Strings(conn.Do("LRANGE", gconst.LobbyClubEventListPrefix+clubID, cursor,
		cursor+maxLoad-1))
	if err != nil && err != redis.ErrNil {
		log.Println("onLoadEvents, redis error:", err)
		sendGenericError(w, ClubOperError_CERR_Database_IO)
		return
	}

	events := make([]*MsgClubEvent, 0, maxLoad)
	if err == nil {
		events = loadEventsByEventIDs(eventIDs, clubID, conn, userID)
	}

	loadEventReply := &MsgClubLoadEventsReply{}
	loadEventReply.Events = events
	cursor32 := int32(cursor + len(eventIDs))
	loadEventReply.Cursor = &cursor32

	b, err := proto.Marshal(loadEventReply)
	if err != nil {
		log.Println("onLoadEvents, marshal error:", err)
		sendGenericError(w, ClubOperError_CERR_Encode_Decode)
		return
	}

	sendMsgClubReply(w, ClubReplyCode_RCOperation, b)
}

// loadEventsByEventIDs 根据事件ID加载事件，由于需要把buffer转换为proto对象，有一定的cpu压力
func loadEventsByEventIDs(eventIDs []string, clubID string, conn redis.Conn, userID string) []*MsgClubEvent {
	conn.Send("MULTI")
	for _, eID := range eventIDs {
		conn.Send("HGET", gconst.LobbyClubEventTablePrefix+clubID, eID)
	}

	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		log.Panicln("loadEventsByEventIDs failed:", err)
	}

	events := make([]*MsgClubEvent, 0, len(values))
	for i, v := range values {
		b, err := redis.Bytes(v, nil)
		if err != nil {
			log.Println("loadEventsByEventIDs, convert value to bytes failed:", err)
			continue
		}

		e := &MsgClubEvent{}
		err = proto.Unmarshal(b, e)
		if err != nil {
			log.Println("loadEventsByEventIDs, unmarshal bytes to event failed:", err)
			continue
		}

		unReadInt, _ := redis.Int(conn.Do("SISMEMBER", gconst.LobbyClubUnReadEventUserSetPrefix+clubID+":"+userID, eventIDs[i]))
		// 全部未读
		unRead := false
		if unReadInt > 0 {
			unRead = true
		}

		e.Unread = &unRead

		events = append(events, e)
	}

	if len(events) > 0 {
		loadDisplayInfoForEvents(events, conn)
	}

	return events
}

func loadDisplayInfoForEvents(events []*MsgClubEvent, conn redis.Conn) {
	for _, ev := range events {
		evType := ClubEventType(ev.GetEvtType())

		userID := ev.GetUserID1()
		if userID == "" {
			continue
		}

		switch evType {
		case ClubEventType_CEVT_NewApplicant, ClubEventType_CEVT_Join, ClubEventType_CEVT_Quit:
			ev.DisplayInfo1 = loadDisplayInfoByUserID(userID, conn)
			break
		}
	}
}

func loadDisplayInfoByUserID(userID string, conn redis.Conn) *MsgClubDisplayInfo {
	strValues, err := redis.Strings(conn.Do("HMGET", gconst.LobbyUserTablePrefix+userID, "nick", "gender", "avatarUrl", "avatarID"))
	if err != nil {
		log.Panicln("loadDisplayInfoByUserID, redis err:", err)
	}

	displayInfo := &MsgClubDisplayInfo{}
	nick := strValues[0]
	displayInfo.Nick = &nick
	gender, _ := strconv.Atoi(strValues[1])
	sex32 := uint32(gender)
	displayInfo.Gender = &sex32
	headIconURL := strValues[2]
	displayInfo.HeadIconURL = &headIconURL
	avatarID, _ := strconv.Atoi(strValues[3])
	avatarID32 := int32(avatarID)
	displayInfo.AvatarID = &avatarID32

	return displayInfo
}
