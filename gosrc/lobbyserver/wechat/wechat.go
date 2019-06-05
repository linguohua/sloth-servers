package wechat

import (
	"gconst"
	"lobbyserver/lobby"

	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

var (
	weChatAPPID     = ""
	weChatAPPSecret = ""
	// urlWeChatGetAccessToken = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	urlWeChatGetAccessToken = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	// urlWeChatGetUserInfo    = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN"
	// urlWeChatGetJSTicket    = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
	// urlWeChatRefreshToken   = "https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s"
	// urlWeChatSysAccessToken = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
)

// MyWeChatConfig 微信配置
type MyWeChatConfig struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

func loadWeChatConfig() *MyWeChatConfig {
	conn := lobby.Pool().Get()
	defer conn.Close()

	fields, err := redis.Strings(conn.Do("HMGET", gconst.LobbyWeChatConfig, "appID", "appAcret"))
	if err != nil {
		log.Println("loadUserInfoFromRedis, error", err)
		return nil
	}

	wechatCfg := &MyWeChatConfig{}
	wechatCfg.AppID = fields[0]
	wechatCfg.AppSecret = fields[1]

	return wechatCfg
}

// InitWechat 初始化wechat
func InitWechat() {
	weChatConfig := loadWeChatConfig()

	weChatAPPID = weChatConfig.AppID
	weChatAPPSecret = weChatConfig.AppSecret
}

// GetAppID 获得APP ID
func GetAppID() string {
	return weChatAPPID
}
