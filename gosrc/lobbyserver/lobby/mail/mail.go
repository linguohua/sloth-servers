package mail

import (
	"lobbyserver/lobby"
)

// InitWith init
func InitWith() {
	lobby.RegHTTPHandle("POST", "/sendMail", handlerSendMail)
	lobby.RegHTTPHandle("GET", "/loadMails", handlerLoadMail)
	lobby.RegHTTPHandle("GET", "/setMailRead", handlerSetMsgRead)
	lobby.RegHTTPHandle("GET", "/deleteMail", handlerSendMail)
	lobby.RegHTTPHandle("GET", "/receiveAttachment", handlerReceiveAttahment)
}
