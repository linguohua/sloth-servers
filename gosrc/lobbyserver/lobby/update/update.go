package update

import (
	"lobbyserver/lobby"
)

// InitWith init
func InitWith() {
	lobby.RegHTTPHandle("GET", "/upgradeQuery", handlerUpgradeQuery)
	lobby.RegHTTPHandle("POST", "/webapi/update/uploadModule", handlerUpload)
	lobby.RegHTTPHandle("GET", "/webpi/update/loadAllModules", handlerLoadAllModules)
	lobby.RegHTTPHandle("POST", "/webapi/update/deleteModules", handlerDeleteModules)

	initConditionVariableCfg()

	initModulesMgr()
}
