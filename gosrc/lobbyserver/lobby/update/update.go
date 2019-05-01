package update

import (
	"lobbyserver/lobby"
)

// InitWith init
func InitWith() {
	lobby.MainRouter.HandleFunc("/upgradeQuery", handlerUpgradeQuery)
	lobby.MainRouter.HandleFunc("/webapi/update/uploadModule", handlerUpload)
	lobby.MainRouter.HandleFunc("/webpi/update/loadAllModules", handlerLoadAllModules)
	lobby.MainRouter.HandleFunc("/webapi/update/deleteModules", handlerDeleteModules)

	initConditionVariableCfg()

	initModulesMgr()
}
