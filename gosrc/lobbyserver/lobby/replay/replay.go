package replay

import (
	"lobbyserver/lobby"
)

// InitWith init
func InitWith() {
	lobby.MainRouter.HandleFunc("/lrproom", handleLoadReplayRooms)
	lobby.MainRouter.HandleFunc("/lrprecord", handleLoadReplayRecord)
}
