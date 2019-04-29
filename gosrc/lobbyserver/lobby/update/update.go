package update

import (
	"lobbyserver/lobby"
)

// InitWith init
func InitWith() {
	lobby.AccRawHTTPHandlers["/upgradeQuery"] = handlerUpgradeQuery
	lobby.AccRawHTTPHandlers["/upload"] = handlerUpload

	initConditionVariableCfg()
}
