package auth

import (
	"crypto/md5"
	"fmt"
	"gconst"
	"lobbyserver/config"
	"lobbyserver/lobby"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func replyRegister(w http.ResponseWriter, registerReply *lobby.MsgRegisterReply) {
	buf, err := proto.Marshal(registerReply)
	if err != nil {
		log.Println("replyRegister, Marshal err:", err)
		return
	}

	w.Write(buf)
}

func addUser2Set(userID string) {
	conn := lobby.Pool().Get()
	defer conn.Close()

	conn.Do("SADD", gconst.LobbyUserSet, userID)
}

func registerAccount(account string, passwdMD5 string, userInfo *lobby.UserInfo, clientInfo *lobby.ClientInfo) {
	mySQLUtil := lobby.MySQLUtil()
	err := mySQLUtil.RegisterAccount(account, passwdMD5, userInfo, clientInfo)
	if err != nil {
		log.Error("registerAccount error:", err)
	}

	// 注册账号的时候，送钻石，可以在配置中配
	lastDiamond, errCode := mySQLUtil.UpdateDiamond(userInfo.GetUserID(), int64(config.DefaultDiamond))
	if errCode != 0 {
		log.Error("registerAccount UpdateDiamond, errCode:", errCode)
	}

	userInfo.Diamond = &lastDiamond

	saveUserInfo2Redis(userInfo)

	addUser2Set(userInfo.GetUserID())
}

func handlerRegister(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	reply := &lobby.MsgRegisterReply{}

	if account == "" {
		errCode := int32(lobby.RegisterError_ErrAccountIsEmpty)
		reply.Result = &errCode
		replyRegister(w, reply)
		return
	}

	if password == "" {
		errCode := int32(lobby.RegisterError_ErrPasswordIsEmpty)
		reply.Result = &errCode
		replyRegister(w, reply)
		return
	}

	// 检查手机号是否已经注册过, 如果已经注册过，返回错误
	// 如果没注册过，则生成个新用户
	mySQLUtil := lobby.MySQLUtil()
	userID, isNew := mySQLUtil.LoadOrGenerateUserID(account)
	if !isNew {
		errCode := int32(lobby.RegisterError_ErrAccountExist)
		reply.Result = &errCode
		replyRegister(w, reply)
		return
	}

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

	userInfo := &lobby.UserInfo{}
	userInfo.UserID = &userID

	data := []byte(password)
	passwdMD5 := fmt.Sprintf("%x", md5.Sum(data))

	registerAccount(account, passwdMD5, userInfo, clientInfo)

	tk := lobby.GenTK(userID)
	reply.Token = &tk
	errCode := int32(0)
	reply.Result = &errCode
	replyRegister(w, reply)

}
