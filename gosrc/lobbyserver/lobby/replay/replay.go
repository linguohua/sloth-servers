package replay

import (
	"lobbyserver/lobby"
)

// InitWith init
func InitWith() {
	lobby.AccUserIDHTTPHandlers["/lrproom"] = handleLoadReplayRooms
	lobby.AccUserIDHTTPHandlers["/lrprecord"] = handleLoadReplayRecord
}
