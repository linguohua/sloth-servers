package update

import (
	"lobbyserver/lobby"
	"net/http"
)

var (
	updateUtil = &myUpdateUtil{}
)

// myRoomUtil implements IRoomUtil
type myUpdateUtil struct {
}

func (*myUpdateUtil) CheckUpdate(r *http.Request) bool {
	ctx := parseFromHTTPReq(r)
	cfg := mmgr.findModuleCfg(ctx)
	if cfg == nil {
		cfg, _ = mmgr.getDefaultCfg(ctx)
	}

	if cfg != nil {
		return true
	}

	return false
}

// InitWith init
func InitWith() {
	lobby.SetUpdateUtil(updateUtil)
	lobby.RegHTTPHandle("GET", "/upgradeQuery", handlerUpgradeQuery)
	lobby.RegHTTPHandle("POST", "/webapi/update/uploadModule", handlerUpload)
	lobby.RegHTTPHandle("GET", "/webpi/update/loadAllModules", handlerLoadAllModules)
	lobby.RegHTTPHandle("POST", "/webapi/update/deleteModules", handlerDeleteModules)

	initConditionVariableCfg()

	initModulesMgr()
}
