package chat

import (
	"lobbyserver/config"
	"lobbyserver/lobby"
)

// InitWith init
func InitWith() {
	loadSensitiveWordDictionary(config.SensitiveWordFilePath)
	lobby.AccUserIDHTTPHandlers["/chat"] = handlerChat
}
