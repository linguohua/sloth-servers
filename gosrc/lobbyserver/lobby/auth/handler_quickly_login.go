package auth

import (
	"net/http"
	"lobbyserver/lobby"
	"fmt"
	"crypto/md5"
	uuid "github.com/satori/go.uuid"
)

func replyQuicklyLogin(w http.ResponseWriter, loginReply *lobby.MsgLoginReply) {
	replyLogin(w, loginReply)
}

func handlerQuicklyLogin(w http.ResponseWriter, r *http.Request) {
	qMod := r.URL.Query().Get("qMod")
	modV := r.URL.Query().Get("modV")
	csVer := r.URL.Query().Get("csVer")
	lobbyVer := r.URL.Query().Get("lobbyVer")
	operatingSystem := r.URL.Query().Get("operatingSystem")
	operatingSystemFamily := r.URL.Query().Get("operatingSystemFamily")
	deviceUniqueIdentifier := r.URL.Query().Get("deviceUniqueIdentifier")
	deviceName := r.URL.Query().Get("deviceName")
	deviceModel := r.URL.Query().Get("deviceModel")
	network := r.URL.Query().Get("network")

	account := r.URL.Query().Get("account")
	password := r.URL.Query().Get("password")

	if account == "" {
		// TODO: 生成新账号，新密码
		uid, _ := uuid.NewV4()
		account = fmt.Sprintf("%v", uid)

		uid, _ = uuid.NewV4()
		password = fmt.Sprintf("%v", uid)
	}

	loginReply := &lobby.MsgLoginReply{}

	mySQLUtil := lobby.MySQLUtil()
	userID, isNew := mySQLUtil.GetOrGenerateUserID(account)
	if isNew {
		clientInfo := &lobby.ClientInfo{}
		clientInfo.QMod = &qMod
		clientInfo.ModV = &modV
		clientInfo.CsVer = &csVer
		clientInfo.LobbyVer = &lobbyVer
		clientInfo.OperatingSystem = &operatingSystem
		clientInfo.OperatingSystemFamily = &operatingSystemFamily
		clientInfo.DeviceUniqueIdentifier = &deviceUniqueIdentifier
		clientInfo.DeviceName = &deviceName
		clientInfo.DeviceModel = &deviceModel
		clientInfo.Network = &network

		data := []byte(password)
		passwdMD5 := fmt.Sprintf("%x", md5.Sum(data))

		mySQLUtil.RegisterAccount(userID, account, passwdMD5, clientInfo)
	} else {
		// 旧的一定要做密码校检
		// myPassword := mySQLUtil.GetPasswordBy(account)
		// if myPassword == "" {
		// 	errCode := int32(lobby.LoginError_ErrAccountNotSetPassword)
		// 	loginReply.Result = &errCode
		// 	replyQuicklyLogin(w, loginReply)

		// 	return
		// }

		// if password != myPassword {
		// 	errCode := int32(lobby.LoginError_ErrPasswordNotMatch)
		// 	loginReply.Result = &errCode
		// 	replyQuicklyLogin(w, loginReply)

		// 	return
		// }

	}

	uint64UserID := uint64(userID)
	userInfo := &lobby.UserInfo{}
	userInfo.UserID = &uint64UserID

	// 生成token给客户端
	tk := lobby.GenTK(fmt.Sprintf("%d", userID))

	errCode := int32(lobby.LoginError_ErrLoginSuccess)

	loginReply.Result = &errCode
	loginReply.Token = &tk
	loginReply.UserInfo = userInfo
	replyQuicklyLogin(w, loginReply)

}
