package auth

import (
	"fmt"
	"gconst"
	"lobbyserver/lobby"
	"net/http"
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
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

func loadUserInfoFromRedis(userID uint64) *lobby.UserInfo {
	conn := lobby.Pool().Get()
	defer conn.Close()

	key := fmt.Sprintf("%s%d", gconst.LobbyUserTablePrefix, userID)

	fields, err := redis.Strings(conn.Do("HMGET", key, "openID", "nickName", "sex", "provice", "city", "country", "headImgURL", "phone"))
	if err != nil {
		log.Println("loadUserInfoFromRedis, error", err)
		return nil
	}

	openID := fields[0]
	nickName := fields[1]
	sex, _ := strconv.Atoi(fields[2])
	provice := fields[3]
	city := fields[4]
	country := fields[5]
	headImgURL := fields[6]
	phone := fields[7]

	userInfo := &lobby.UserInfo{}
	userInfo.UserID = &userID
	userInfo.OpenID = &openID
	userInfo.NickName = &nickName
	sexUint32 := uint32(sex)
	userInfo.Sex = &sexUint32
	userInfo.Province = &provice
	userInfo.City = &city
	userInfo.Country = &country
	userInfo.HeadImgUrl = &headImgURL
	userInfo.Phone = &phone

	return nil
}

func handlerAccountLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	// TODO: 更新客户端信息

	userInfo := loadUserInfoFromRedis(userID)

	// 生成token给客户端
	tk := lobby.GenTK(fmt.Sprintf("%d", userID))

	errCode := int32(lobby.LoginError_ErrLoginSuccess)

	loginReply.Result = &errCode
	loginReply.Token = &tk
	loginReply.UserInfo = userInfo
	replyAccountLogin(w, loginReply)
}
