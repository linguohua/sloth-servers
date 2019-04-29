package auth

import (
	"crypto/md5"
	"fmt"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"lobbyserver/lobby"
	"net/http"
)

func replyRegister(w http.ResponseWriter, registerReply *lobby.MsgRegisterReply) {
	buf, err := proto.Marshal(registerReply)
	if err != nil {
		log.Println("replyRegister, Marshal err:", err)
		return
	}

	w.Write(buf)
}

func handlerRegister(w http.ResponseWriter, r *http.Request) {
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
	userID, isNew := mySQLUtil.GetOrGenerateUserID(account)
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

	data := []byte(password)
	passwdMD5 := fmt.Sprintf("%x", md5.Sum(data))

	err := mySQLUtil.RegisterAccount(userID, account, passwdMD5, "", clientInfo)
	if err != nil {
		errCode := int32(lobby.RegisterError_ErrWriteDatabaseFailed)
		reply.Result = &errCode
		replyRegister(w, reply)
		return
	}

	// TODO: 需要保存到redis

	tk := lobby.GenTK(fmt.Sprintf("%d", userID))
	reply.Token = &tk
	errCode := int32(0)
	reply.Result = &errCode
	replyRegister(w, reply)

}
