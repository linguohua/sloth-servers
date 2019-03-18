package lobby

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"encoding/json"
)

func handleQueryUserRoomInfo(w http.ResponseWriter, r *http.Request, userID string) {
	log.Println("handleQueryRoomState, userID:", userID)
	roomInfo := loadLastRoomInfo(userID)
	if roomInfo == nil {
		w.WriteHeader(404)
		w.Write([]byte("User not in Room!"))
		return
	}

	type Reply struct {
		RoomID string `json:"roomID"`
		ArenaID string `json:"arenaID"`
		RoomCfg string `json:"roomCfg"`
	}

	reply := &Reply{}
	reply.RoomID = roomInfo.GetRoomID()
	reply.ArenaID = roomInfo.GetArenaID()
	reply.RoomCfg = roomInfo.GetConfig()

	buf, err := json.Marshal(reply)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(buf)
	return
}
