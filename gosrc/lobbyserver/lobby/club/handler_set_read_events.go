package club

import (
	"log"
	"net/http"
	"gconst"
	"strings"
	"lobbyserver/lobby"

	"github.com/garyburd/redigo/redis"
)

// onSetReadEvents 剔除某个成员
func onSetReadEvents(w http.ResponseWriter, r *http.Request, userID string) {
	log.Println("onSetReadEvents, userID:", userID)

	var query = r.URL.Query()
	// 事件id列表
	eventIDStrs := query.Get("eIDs")
	clubID := query.Get("clubID")

	if eventIDStrs == "" {
		log.Println("onSetReadEvents, error, need eIDs")
		sendGenericError(w, ClubOperError_CERR_Invalid_Input_Parameter)
		return
	}

	if clubID == "" {
		log.Println("onSetReadEvents, error, need clubID")
		sendGenericError(w, ClubOperError_CERR_Invalid_Input_Parameter)
		return
	}

	eventIDs := strings.Split(eventIDStrs, ",")

	_, ok := clubMgr.clubs[clubID]
	if !ok {
		log.Printf("onSetReadEvents, no such club:%s\n", clubID)
		sendGenericError(w, ClubOperError_CERR_Club_Not_Exist)
		return
	}

	// 请求redis获得该玩家的所有俱乐部数量
	conn := lobby.Pool().Get()
	defer conn.Close()

	// 检查事件是否可以直接被设置为read
	conn.Send("MULTI")
	for _, eID := range eventIDs {
		conn.Send("HGET", gconst.LobbyClubNeedHandledTablePrefix+clubID, eID)
	}

	targets, err := redis.Strings(conn.Do("EXEC"))
	if err != nil {
		log.Println("onSetReadEvents, redis err:", err)
		sendGenericError(w, ClubOperError_CERR_Database_IO)
		return
	}

	for index, t := range targets {
		if t == userID {
			log.Println("onSetReadEvents, event can't be set to read, event id:", eventIDs[index])
			sendGenericError(w, ClubOperError_CERR_Invalid_Input_Parameter)
			return
		}
	}

	// 检查是否所有事件都是未读
	conn.Send("MULTI")
	for _, eID := range eventIDs {
		conn.Send("SISMEMBER", gconst.LobbyClubUnReadEventUserSetPrefix+clubID+":"+userID, eID)
	}

	ints, err := redis.Ints(conn.Do("EXEC"))
	if err != nil {
		log.Println("onSetReadEvents, redis err:", err)
		sendGenericError(w, ClubOperError_CERR_Database_IO)
		return
	}

	for index, t := range ints {
		if t < 1 {
			log.Println("onSetReadEvents, event not unread, event id:", eventIDs[index])
			sendGenericError(w, ClubOperError_CERR_Invalid_Input_Parameter)
			return
		}
	}

	// 未读事件清理
	conn.Send("MULTI")
	for _, eID := range eventIDs {
		conn.Send("LREM", gconst.LobbyClubUnReadEventUserListPrefix+clubID+":"+userID, 1, eID)
		conn.Send("SREM", gconst.LobbyClubUnReadEventUserSetPrefix+clubID+":"+userID, eID)
	}
	conn.Do("EXEC")

	// 操作成功
	sendGenericError(w, ClubOperError_CERR_OK)
}
