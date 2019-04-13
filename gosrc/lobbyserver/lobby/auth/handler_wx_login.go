package auth

import (
	"lobbyserver/wechat"
	"net/http"
	log "github.com/sirupsen/logrus"
	"fmt"
	"lobbyserver/lobby"
	"github.com/golang/protobuf/proto"
)

var (
	errCodeWechatCodeIsEmpty = uint32(1)
	errCodeLoadUserInfoFailed = uint32(2)
)

func replyWxLogin(w http.ResponseWriter, loginReply *lobby.MsgWxLoginReply) {
	buf, err := proto.Marshal(loginReply)
	if err != nil {
		log.Println("replyWxLogin, Marshal err:", err)
		return
	}

	w.Write(buf)
}

func saveUserInfo(userInfoReply *lobby.WxUserInfo, clientInfo *lobby.ClientInfo) {
	mySQLUtil := lobby.MySQLUtil()
	mySQLUtil.UpdateWxUserInfo(userInfoReply, clientInfo)
}

func loadUserInfoFromWeChatServer(wechatCode string) (*lobby.WxUserInfo, error) {
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

	userInfo := &lobby.WxUserInfo{}
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
	wechatCode := r.URL.Query().Get("code")
	if wechatCode == "" {
		loginReply := &lobby.MsgWxLoginReply{}
		loginReply.Result = &errCodeWechatCodeIsEmpty
		replyWxLogin(w, loginReply)

		return;
	}

	userInfo, err := loadUserInfoFromWeChatServer(wechatCode)
	if err != nil {
		loginReply := &lobby.MsgWxLoginReply{}
		loginReply.Result = &errCodeLoadUserInfoFailed
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

	// 保存用户信息
	saveUserInfo(userInfo, clientInfo)

	// TODO: 需要用userId生成token
	// 生成token给客户端
	tk := lobby.GenTK(userInfo.GetOpenID())

	var result = uint32(0);

	loginReply := &lobby.MsgWxLoginReply{}
	loginReply.Result = &result
	loginReply.Token = &tk
	replyWxLogin(w, loginReply)
}