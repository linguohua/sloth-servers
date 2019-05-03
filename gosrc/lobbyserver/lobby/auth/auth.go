package auth

import (
	"lobbyserver/lobby"
)

// InitWith init
func InitWith() {
	lobby.MainRouter.HandleFunc("/wxLogin", handlerWxLogin)
	lobby.MainRouter.HandleFunc("/accountLogin", handlerAccountLogin)
	lobby.MainRouter.HandleFunc("/quicklyLogin", handlerQuicklyLogin)
	lobby.MainRouter.HandleFunc("/register", handlerRegister)
	lobby.MainRouter.HandleFunc("/test", handlerTest)
}
