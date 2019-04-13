

package auth

import (
	"lobbyserver/lobby"
)

// InitWith init
func InitWith() {
	lobby.AccRawHTTPHandlers["/wxLogin"] = handlerWxLogin
	lobby.AccRawHTTPHandlers["/accountLogin"] = handlerAccountLogin
	lobby.AccRawHTTPHandlers["/quicklyLogin"] = handlerQuicklyLogin
	lobby.AccRawHTTPHandlers["/register"] = handlerRegister
	lobby.AccRawHTTPHandlers["/test"] = handlerTest
}