package share

import (
	"lobbyserver/lobby"
)

// InitWith init
func InitWith() {
	lobby.RegHTTPHandle("GET", "/getShareInfo", handlerGetShareInfo) // 获取分享内容
}
