package wechat

import (
	"fmt"
	"log"
)

// LoadJSSDKTicket 从微信服务器拉取js ticket
func LoadJSSDKTicket(accessToken string) (*GetTicketReply, error) {
	urlGetTicket := fmt.Sprintf(urlWeChatGetJSTicket, accessToken)

	log.Println("LoadJSSDKTicket, full url:", urlGetTicket)

	reply := &GetTicketReply{}
	err := loadDataUseHTTPGet(urlGetTicket, reply)
	if err != nil {
		return nil, err
	}

	log.Printf("LoadJSSDKTicket, reply json:%+v\n", reply)
	return reply, nil
}

// LoadSysAccessToken 加载公众号access token
func LoadSysAccessToken() (*LoadSysAccessTokenReply, error) {
	urlGetTicket := fmt.Sprintf(urlWeChatSysAccessToken, weChatAPPID, weChatAPPSecret)

	log.Println("LoadSysAccessToken, full url:", urlGetTicket)

	reply := &LoadSysAccessTokenReply{}
	err := loadDataUseHTTPGet(urlGetTicket, reply)
	if err != nil {
		return nil, err
	}

	log.Printf("LoadSysAccessToken, reply json:%+v\n", reply)
	return reply, nil
}
