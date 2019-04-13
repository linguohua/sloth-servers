package wechat

import (
	"fmt"
	"log"
)

// RefreshToken 从微信服务器刷新token
func RefreshToken(refreshToken string) (*RefreshTokenReply, error) {
	urlRefreshToken := fmt.Sprintf(urlWeChatRefreshToken, weChatAPPID, refreshToken)

	log.Println("RefreshToken, full url:", urlRefreshToken)

	reply := &RefreshTokenReply{}
	err := loadDataUseHTTPGet(urlRefreshToken, reply)
	if err != nil {
		return nil, err
	}

	log.Printf("RefreshToken, reply json:%+v\n", reply)
	return reply, nil
}
