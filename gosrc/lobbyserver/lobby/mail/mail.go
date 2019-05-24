package mail

import (
	"lobbyserver/lobby"
)

var (
	myMailUtil = &mailUtil{}
)

type mailUtil struct {
}

func (*mailUtil) SendMail(userID string, content string, title string) {
	sendMail(userID, content, title)
}

// InitWith init
func InitWith() {

	lobby.SetMailUtil(myMailUtil)
	lobby.RegHTTPHandle("POST", "/sendMail", handlerSendMail)
	lobby.RegHTTPHandle("GET", "/loadMails", handlerLoadMail)
	lobby.RegHTTPHandle("GET", "/setMailRead", handlerSetMsgRead)
	lobby.RegHTTPHandle("GET", "/deleteMail", handlerSendMail)
	lobby.RegHTTPHandle("GET", "/receiveAttachment", handlerReceiveAttahment)
}
