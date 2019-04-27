package wechat

var (
	weChatAPPID             = ""
	weChatAPPSecret         = ""
	urlWeChatGetAccessToken = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	urlWeChatGetUserInfo    = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN"
	urlWeChatGetJSTicket    = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
	urlWeChatRefreshToken   = "https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s"
	urlWeChatSysAccessToken = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
)

// InitWechat 初始化wechat
func InitWechat(myAppID string, myAppSecret string) {
	weChatAPPID = myAppID
	weChatAPPSecret = myAppSecret
}

// GetAppID 获得APP ID
func GetAppID() string {
	return weChatAPPID
}
