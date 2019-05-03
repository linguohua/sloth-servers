package replay

import (
	"lobbyserver/lobby"
)

// InitWith init
func InitWith() {
	lobby.RegHTTPHandle("GET", "/lrproom", handleLoadReplayRooms)
	lobby.RegHTTPHandle("GET", "/lrprecord", handleLoadReplayRecord)
}
