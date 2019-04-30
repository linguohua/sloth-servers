package auth

import (
	"fmt"
	"gconst"
	log "github.com/sirupsen/logrus"
	"lobbyserver/lobby"
	"lobbyserver/wechat"
	"net/http"
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
	sex := userInfo.GetSex()
	provice := userInfo.GetProvince()
	city := userInfo.GetCity()
	country := userInfo.GetCountry()
	headImgURL := userInfo.GetHeadImgUrl()

	conn.Do("HMSET", key, "userID", userID, "openID", openID, "nickName", nickName, "sex", sex,
		"provice", provice, "city", city, "country", country, "headImgURL", headImgURL)
}

func updateWxUserInfo(userInfo *lobby.UserInfo, clientInfo *lobby.ClientInfo) {
	mySQLUtil := lobby.MySQLUtil()
	mySQLUtil.UpdateWxUserInfo(userInfo, clientInfo)

	// 保存到redis
	saveUserInfo2Redis(userInfo)
}

func loadUserInfoFromWeChatServer(wechatCode string) (*lobby.UserInfo, error) {
	// 根据code去微信服务器拉取access_token和openID
	accessTokenReply, err := wechat.LoadAccessTokenFromWeChatServer(wechatCode)
	if err != nil {
		log.Panicln("loadAccessTokenFromWeChatServer err:", err)
		return nil, err
	}

	// 检查结果
	if accessTokenReply.ErrorCode != 0 {
		log.Panicf("loadAccessTokenFromWeChatServer, wechat server reply error code:%d, msg:%s\n",
			accessTokenReply.ErrorCode, accessTokenReply.ErrorMsg)
		return nil, fmt.Errorf("load access token from wechat server failed, error code:%d", accessTokenReply.ErrorCode)
	}

	// 根据access token拉取用户信息
	userInfoReply, err := wechat.LoadUserInfoFromWeChatServer(accessTokenReply.AcessToken, accessTokenReply.OpenID)
	if err != nil {
		log.Panicln("loadUserInfoFromWeChatServer err:", err)
		return nil, err
	}

	// 检查结果
	if userInfoReply.ErrorCode != 0 {
		log.Panicf("loadUserInfoFromWeChatServer, wechat server reply error code:%d, msg:%s\n",
			userInfoReply.ErrorCode, userInfoReply.ErrorMsg)
		return nil, fmt.Errorf("load userInfo from wechat server failed, error code:%d", accessTokenReply.ErrorCode)
	}

	userInfo := &lobby.UserInfo{}
	userInfo.OpenID = &userInfoReply.OpenID
	userInfo.NickName = &userInfoReply.NickName
	sexUint32 := uint32(userInfoReply.Sex)
	userInfo.Sex = &sexUint32
	userInfo.Province = &userInfoReply.Province
	userInfo.City = &userInfoReply.City
	userInfo.Country = &userInfoReply.Country
	userInfo.HeadImgUrl = &userInfoReply.HeadImgURL

	return userInfo, nil
}

func handlerWxLogin(w http.ResponseWriter, r *http.Request) {
	loginReply := &lobby.MsgLoginReply{}

	wechatCode := r.URL.Query().Get("code")
	if wechatCode == "" {
		errCode := int32(lobby.LoginError_ErrParamWechatCodeIsEmpty)
		loginReply.Result = &errCode
		replyWxLogin(w, loginReply)

		return
	}

	userInfo, err := loadUserInfoFromWeChatServer(wechatCode)
	if err != nil {
		errCode := int32(lobby.LoginError_ErrLoadWechatUserInfoFailed)
		loginReply.Result = &errCode
		replyWxLogin(w, loginReply)

		return
	}

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
	userID, _ := mySQLUtil.GetOrGenerateUserID(userInfo.GetOpenID())
	userInfo.UserID = &userID
	// 保存用户信息
	updateWxUserInfo(userInfo, clientInfo)

	// 生成token给客户端
	tk := lobby.GenTK(userID)

	errCode := int32(lobby.LoginError_ErrLoginSuccess)

	loginReply.Result = &errCode
	loginReply.Token = &tk
	loginReply.UserInfo = userInfo
	replyWxLogin(w, loginReply)
}
