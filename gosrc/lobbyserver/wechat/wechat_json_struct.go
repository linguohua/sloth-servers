package wechat

// LoadAccessTokenReply 微信拉取access token回复
type LoadAccessTokenReply struct {
	AcessToken   string `json:"access_token"`
	ExpireIn     int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`

	ErrorCode int    `json:"errcode"`
	ErrorMsg  string `json:"errmsg"`
}

// LoadUserInfoReply 微信拉取user info回复
type LoadUserInfoReply struct {
	OpenID     string   `json:"openid"`
	NickName   string   `json:"nickname"`
	Gender        int      `json:"gender"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgURL string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	UnionID    string   `json:"unionid"`
	Tags       string   `json:"tags"`
	Thumb      int      `json:"thumb"`

	ErrorCode int    `json:"errcode"`
	ErrorMsg  string `json:"errmsg"`
}

// GetTicketReply 微信获得JS SDK ticket的回复
type GetTicketReply struct {
	Ticket    string `json:"ticket"`
	ExpireIn  int    `json:"expires_in"`
	ErrorCode int    `json:"errcode"`
	ErrorMsg  string `json:"errmsg"`
}

// RefreshTokenReply 微信获得刷新token的回复
type RefreshTokenReply struct {
	AcessToken   string `json:"access_token"`
	ExpireIn     int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`

	ErrorCode int    `json:"errcode"`
	ErrorMsg  string `json:"errmsg"`
}

// LoadSysAccessTokenReply 微信获得公众号刷新token的回复
type LoadSysAccessTokenReply struct {
	AcessToken string `json:"access_token"`
	ExpireIn   int    `json:"expires_in"`

	ErrorCode int    `json:"errcode"`
	ErrorMsg  string `json:"errmsg"`
}
