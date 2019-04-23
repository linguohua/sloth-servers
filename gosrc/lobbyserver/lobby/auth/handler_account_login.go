package auth

import (
	"net/http"
	"lobbyserver/lobby"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"fmt"
)

func replyLogin(w http.ResponseWriter, loginReply *lobby.MsgLoginReply) {
	buf, err := proto.Marshal(loginReply)
	if err != nil {
		log.Println("replyWxLogin, Marshal err:", err)
		return
	}

	w.Write(buf)
}

func replyAccountLogin(w http.ResponseWriter, loginReply *lobby.MsgLoginReply) {
	replyLogin(w, loginReply)
}

func handlerAccountLogin(w http.ResponseWriter, r *http.Request) {
	phoneNum := r.URL.Query().Get("phoneNum")
	password := r.URL.Query().Get("password")

	loginReply := &lobby.MsgLoginReply{}

	if phoneNum == "" {
		errCode := int32(lobby.LoginError_ErrParamAccountIsEmpty)
		loginReply.Result = &errCode
		replyAccountLogin(w, loginReply)

		return
	}

	if password == "" {
		errCode := int32(lobby.LoginError_ErrParamPasswordIsEmpty)
		loginReply.Result = &errCode
		replyAccountLogin(w, loginReply)

		return
	}

	mySQLUtil := lobby.MySQLUtil()

	userID := mySQLUtil.GetUserIDBy(phoneNum)
	if userID == 0 {
		errCode := int32(lobby.LoginError_ErrAccountNotExist)
		loginReply.Result = &errCode
		replyAccountLogin(w, loginReply)

		return
	}

	myPassword := mySQLUtil.GetPasswordBy(phoneNum)
	if myPassword == "" {
		errCode := int32(lobby.LoginError_ErrAccountNotSetPassword)
		loginReply.Result = &errCode
		replyAccountLogin(w, loginReply)

		return
	}

	if password != myPassword {
		errCode := int32(lobby.LoginError_ErrPasswordNotMatch)
		loginReply.Result = &errCode
		replyAccountLogin(w, loginReply)

		return
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
	replyAccountLogin(w, loginReply)
}