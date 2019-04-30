package update

import (
	"lobbyserver/lobby"
)

// InitWith init
func InitWith() {
	lobby.MainRouter.HandleFunc("/upgradeQuery", handlerUpgradeQuery)
	lobby.MainRouter.HandleFunc("/upload", handlerUpload)

	initConditionVariableCfg()

	initModulesMgr()
}
