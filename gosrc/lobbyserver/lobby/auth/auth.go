package auth

import (
	"lobbyserver/lobby"
)

// InitWith init
func InitWith() {
	lobby.RegHTTPHandle("GET", "/wxLogin", handlerWxLogin)
	lobby.RegHTTPHandle("GET", "/accountLogin", handlerAccountLogin)
	lobby.RegHTTPHandle("GET", "/quicklyLogin", handlerQuicklyLogin)
	lobby.RegHTTPHandle("GET", "/register", handlerRegister)
	lobby.RegHTTPHandle("GET", "/test", handlerTest)
}
