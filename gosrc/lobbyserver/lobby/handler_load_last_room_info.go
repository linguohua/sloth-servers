package lobby

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func handlerLoadLastRoomInfo(w http.ResponseWriter, r *http.Request, userID string) {
	log.Println("handlerLoadLastRoomInfo call, userID:", userID)

	var lastRoomInfo = loadLastRoomInfo(userID)
	if lastRoomInfo != nil {
		log.Printf("handlerRequestRoomInfo, User %s in other room, roomNumber: %s, roomId:%s", userID, lastRoomInfo.GetRoomNumber(), lastRoomInfo.GetRoomID())
		replyRequestRoomInfo(w, 0, lastRoomInfo)
		return
	}

	replyRequestRoomInfo(w, int32(MsgError_ErrRoomNotExist), nil)
}
