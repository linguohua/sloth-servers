package wechat

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// LoadAccessTokenFromWeChatServer 从wechat服务器加载access token
func LoadAccessTokenFromWeChatServer(wechatCode string) (*LoadAccessTokenReply, error) {
	urlGetAccessToken := fmt.Sprintf(urlWeChatGetAccessToken, weChatAPPID, weChatAPPSecret, wechatCode)
	log.Println("loadAccessTokenFromWeChatServer, full url:", urlGetAccessToken)

	reply := &LoadAccessTokenReply{}
	err := loadDataUseHTTPGet(urlGetAccessToken, reply)
	if err != nil {
		return nil, err
	}

	log.Printf("loadAccessTokenFromWeChatServer, reply json:%+v\n", reply)
	return reply, nil
}

// LoadUserInfoFromWeChatServer 从wechat服务器加载用户信息
func LoadUserInfoFromWeChatServer(accessToken string, openID string) (*LoadUserInfoReply, error) {
	urlGetUserInfo := fmt.Sprintf(urlWeChatGetUserInfo, accessToken, openID)

	log.Println("loadUserInfoFromWeChatServer, full url:", urlGetUserInfo)

	reply := &LoadUserInfoReply{}
	err := loadDataUseHTTPGet(urlGetUserInfo, reply)
	if err != nil {
		return nil, err
	}

	log.Printf("loadUserInfoFromWeChatServer, reply json:%+v\n", reply)
	return reply, nil
}

func loadDataUseHTTPGet(url string, jsonStruct interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	// 确保body关闭
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(jsonStruct)
	return err
}
