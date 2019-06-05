package wechat

// LoadAccessTokenReply 微信拉取access token回复
// type LoadAccessTokenReply struct {
// 	AcessToken   string `json:"access_token"`
// 	ExpireIn     int    `json:"expires_in"`
// 	RefreshToken string `json:"refresh_token"`
// 	OpenID       string `json:"openid"`
// 	Scope        string `json:"scope"`

// 	ErrorCode int    `json:"errcode"`
// 	ErrorMsg  string `json:"errmsg"`
// }

// // LoadUserInfoReply 微信拉取user info回复
// type LoadUserInfoReply struct {
// 	OpenID     string   `json:"openid"`
// 	NickName   string   `json:"nickname"`
// 	Gender     int      `json:"gender"`
// 	Province   string   `json:"province"`
// 	City       string   `json:"city"`
// 	Country    string   `json:"country"`
// 	HeadImgURL string   `json:"headimgurl"`
// 	Privilege  []string `json:"privilege"`
// 	UnionID    string   `json:"unionid"`
// 	Tags       string   `json:"tags"`
// 	Thumb      int      `json:"thumb"`

// 	ErrorCode int    `json:"errcode"`
// 	ErrorMsg  string `json:"errmsg"`
// }

// LoadAccessTokenReply 微信拉取access token回复
type LoadAccessTokenReply struct {
	SessionKey string `json:"session_key"`
	OpenID     string `json:"openid"`
	UnionID    string `json:"unionid"`

	ErrorCode int    `json:"errcode"`
	ErrorMsg  string `json:"errmsg"`
}

// WeiXinUserPlusInfo 微信用户信息
type WeiXinUserPlusInfo struct {
	OpenID    string `json:"openId"`
	NickName  string `json:"nickName"`
	Gender    int32  `json:"gender"`
	AvatarURL string `json:"avatarUrl"`
	UnionID   string `json:"unionId"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Country   string `json:"country"`
}
