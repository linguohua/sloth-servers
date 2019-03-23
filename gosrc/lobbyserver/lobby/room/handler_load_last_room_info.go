package room

import (
	"lobbyserver/lobby"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func handlerLoadLastRoomInfo(w http.ResponseWriter, r *http.Request, userID string) {
	log.Println("handlerLoadLastRoomInfo call, userID:", userID)

	var lastRoomInfo = loadLastRoomInfo(userID)
	if lastRoomInfo != nil {
		log.Printf("handlerRequestRoomInfo, User %s in other room, roomNumber: %s, roomId:%s", userID, lastRoomInfo.GetRoomNumber(), lastRoomInfo.GetRoomID())
		replyRequestRoomInfo(w, 0, lastRoomInfo)
		return
	}

	replyRequestRoomInfo(w, int32(lobby.MsgError_ErrRoomNotExist), nil)
}
