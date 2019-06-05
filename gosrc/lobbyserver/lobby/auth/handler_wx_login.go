package auth

import (
	"fmt"
	"gconst"
	"io/ioutil"
	"lobbyserver/lobby"
	"lobbyserver/wechat"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func replyWxLogin(w http.ResponseWriter, loginReply *lobby.MsgLoginReply) {
	replyLogin(w, loginReply)
}

func saveUserInfo2Redis(userInfo *lobby.UserInfo) {
	// 获取redis链接，并退出函数时释放
	conn := lobby.Pool().Get()
	defer conn.Close()

	key := fmt.Sprintf("%s%s", gconst.LobbyUserTablePrefix, userInfo.GetUserID())

	userID := userInfo.GetUserID()
	openID := userInfo.GetOpenID()
	nickName := userInfo.GetNickName()
	gender := userInfo.GetGender()
	provice := userInfo.GetProvince()
	city := userInfo.GetCity()
	country := userInfo.GetCountry()
	headImgURL := userInfo.GetHeadImgUrl()
	diamond := userInfo.GetDiamond()

	conn.Do("HMSET", key, "userID", userID, "openID", openID, "nickName", nickName, "gender", gender,
		"provice", provice, "city", city, "country", country, "headImgURL", headImgURL, "diamond", diamond)
}

func updateWxUserInfo(userInfo *lobby.UserInfo, clientInfo *lobby.ClientInfo) {
	mySQLUtil := lobby.MySQLUtil()
	mySQLUtil.UpdateWxUserInfo(userInfo, clientInfo)

	// 更新用户redis中信息和钻石
	diamond := mySQLUtil.LoadUserDiamond(userInfo.GetUserID())
	userInfo.Diamond = &diamond

	saveUserInfo2Redis(userInfo)
}

func weiXinPlusUserInfo2UserInof(wxUserInfo *wechat.WeiXinUserPlusInfo) *lobby.UserInfo {
	userInfo := &lobby.UserInfo{}
	userInfo.OpenID = &wxUserInfo.OpenID
	userInfo.NickName = &wxUserInfo.NickName
	gender := uint32(wxUserInfo.Gender)
	userInfo.Gender = &gender
	userInfo.HeadImgUrl = &wxUserInfo.AvatarURL

	return userInfo
}

func handlerWxLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("handlerCreateRoom error:", err)
		return
	}

	loginReply := &lobby.MsgLoginReply{}
	wxLogin := &lobby.MsgWxLogin{}
	err = proto.Unmarshal(body, wxLogin)
	if err != nil {
		log.Println("handlerWxLogin, Unmarshal err:", err)
		errCode := int32(lobby.LoginError_ErrParamDecode)
		loginReply.Result = &errCode
		replyWxLogin(w, loginReply)
		return
	}

	if wxLogin.GetCode() == "" {
		errCode := int32(lobby.LoginError_ErrParamInvalidCode)
		loginReply.Result = &errCode
		replyWxLogin(w, loginReply)
		return
	}

	if wxLogin.GetEncrypteddata() == "" {
		errCode := int32(lobby.LoginError_ErrParamInvalidEncrypteddata)
		loginReply.Result = &errCode
		replyWxLogin(w, loginReply)
		return
	}

	if wxLogin.GetIv() == "" {
		errCode := int32(lobby.LoginError_ErrParamInvalidIv)
		loginReply.Result = &errCode
		replyWxLogin(w, loginReply)
		return
	}

	accessTokenReply, err := wechat.LoadAccessTokenFromWeChatServer(wxLogin.GetCode())
	if err != nil {
		log.Panicln("handlerWxLogin loadAccessTokenFromWeChatServer err:", err)
		errCode := int32(lobby.LoginError_ErrWxAuthFailed)
		loginReply.Result = &errCode
		replyWxLogin(w, loginReply)

		return
	}

	// 检查结果
	if accessTokenReply.ErrorCode != 0 {
		log.Errorf("loadAccessTokenFromWeChatServer, wechat server reply error code:%d, msg:%s\n", accessTokenReply.ErrorCode, accessTokenReply.ErrorMsg)
		errCode := int32(lobby.LoginError_ErrWxAuthFailed)
		loginReply.Result = &errCode
		replyWxLogin(w, loginReply)

		return
	}

	userInfoReply, err := wechat.GetWeiXinPlusUserInfo(accessTokenReply.SessionKey, wxLogin.GetEncrypteddata(), wxLogin.GetIv())
	if err != nil {
		log.Error("loadUserInfoFromWeChatServer err:", err)
		errCode := int32(lobby.LoginError_ErrDecodeUserInfoFailed)
		loginReply.Result = &errCode
		replyWxLogin(w, loginReply)

		return
	}

	userInfo := weiXinPlusUserInfo2UserInof(userInfoReply)

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

	mySQLUtil := lobby.MySQLUtil()
	userID, isNew := mySQLUtil.LoadOrGenerateUserID(userInfo.GetOpenID())
	userInfo.UserID = &userID

	if isNew {
		// TODO注册账号
		registerAccount(userInfo.GetOpenID(), "", userInfo, clientInfo)
	} else {
		// 更新用户信息
		updateWxUserInfo(userInfo, clientInfo)
	}

	// 生成token给客户端
	tk := lobby.GenTK(userID)

	errCode := int32(lobby.LoginError_ErrLoginSuccess)

	loginReply.Result = &errCode
	loginReply.Token = &tk
	loginReply.UserInfo = userInfo
	replyWxLogin(w, loginReply)
}
