package club

import (
	"lobbyserver/lobby"
)
var (
	clubMgr *MyClubMgr
)

// InitWith init
func InitWith() {
	clubMgr := newClubMgr()

	lobby.SetClubMgr(clubMgr)

	lobby.RegHTTPHandle("GET", "/createClub", onCreateClub)
}
