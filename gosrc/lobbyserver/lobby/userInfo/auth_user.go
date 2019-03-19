package userInfo

// UserInfo 授权成功返回结果
type UserInfo struct {
	UserID      int64  `json:"userID"`
	UserName    string `json:"userName"`
	SdkUserName string `json:"sdkUserName"`
	SdkUserNick string `json:"sdkUserNick"`
	SdkUserSex  string `json:"sdkUserSex"`
	SdkUserLogo string `json:"sdkUserLogo"`
}

// ResultJSON 授权结果JSON
type ResultJSON struct {
	State int32     `json:"state"`
	Data  *UserInfo `json:"data"`
}