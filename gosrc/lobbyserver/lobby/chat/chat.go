package chat

import (
	"lobbyserver/lobby"
	"lobbyserver/config"
)


// InitWith init
func InitWith() {
	loadSensitiveWordDictionary(config.SensitiveWordFilePath)
	lobby.AccUserIDHTTPHandlers["/chat"] = handlerChat
}
