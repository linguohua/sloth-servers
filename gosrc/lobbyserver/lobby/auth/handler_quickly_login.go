package auth

import (
	"net/http"
	"lobbyserver/lobby"
	"fmt"
	// "crypto/md5"
	uuid "github.com/satori/go.uuid"
	"github.com/golang/protobuf/proto"
	"log"
)

func replyQuicklyLogin(w http.ResponseWriter, loginReply *lobby.MsgQuicklyLoginReply) {
	buf, err := proto.Marshal(loginReply)
	if err != nil {
		log.Println("replyQuicklyLogin, Marshal err:", err)
		return
	}

	w.Write(buf)
}

// 客户端发用户ID上来
// 如果用户不发用户ID上来，则生成一个新的账号
// 如果用户存在，则下发用户信息回去给用户
func handlerQuicklyLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("handlerQuicklyLogin")
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

	mySQLUtil := lobby.MySQLUtil()

	// var uint64UserID uint64
	// isNew := false
	// account := ""
	if account == "" {
		// 生成新账号
		uid, _ := uuid.NewV4()
		account = fmt.Sprintf("%v", uid)
	}

	loginReply := &lobby.MsgQuicklyLoginReply{}

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

		mySQLUtil.RegisterAccount(userID, account, "", "", clientInfo)
	} else {
		// 要校检是否是快速登录账号，快速登录没有密码
		myPassword := mySQLUtil.GetPasswordBy(account)
		if myPassword != "" {
			errCode := int32(lobby.LoginError_ErrPasswordNotMatch)
			loginReply.Result = &errCode
			replyQuicklyLogin(w, loginReply)

			return
		}
	}



	userInfo := loadUserInfoFromRedis(userID)
	if userInfo == nil {
		userInfo = &lobby.UserInfo{}
		userInfo.UserID = &userID
	}

	log.Println("userInfo:", userInfo)

	// 生成token给客户端
	tk := lobby.GenTK(fmt.Sprintf("%d", userID))

	errCode := int32(lobby.LoginError_ErrLoginSuccess)

	loginReply.Result = &errCode
	loginReply.Token = &tk
	loginReply.UserInfo = userInfo
	loginReply.Account = &account
	replyQuicklyLogin(w, loginReply)

}
