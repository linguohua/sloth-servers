package mysql

import (
	// "database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"lobbyserver/lobby"
)

// 查询用户信息
func loadUserInfo(userID uint64) *lobby.UserInfo {
	query := fmt.Sprintf("select user_id,open_id, phone, nick_name, sex, provice, city, country, head_img_url from user where user_id = %d", userID)

	log.Println("loadUserInfo query:", query)

	// return nil
	stmt, err := dbConn.Prepare(query)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	var uint64UserID uint64
	var openID string
	var phone int
	var nickName string
	var sex uint32
	var provice string
	var city string
	var country string
	var headImgURL string

	row := stmt.QueryRow()
	err = row.Scan(&uint64UserID, &openID, &phone, &nickName, &sex, &provice, &city, &country, &headImgURL)
	if err != nil {
		panic(err.Error())
	}

	if uint64UserID == 0 {
		return nil
	}

	phontStr := fmt.Sprintf("%d", phone)

	userInfo := &lobby.UserInfo{}
	userInfo.UserID = &uint64UserID
	userInfo.OpenID = &openID
	userInfo.Phone = &phontStr
	userInfo.NickName = &nickName
	userInfo.Sex = &sex
	userInfo.Province = &provice
	userInfo.City = &city
	userInfo.Country = &country
	userInfo.HeadImgUrl = &headImgURL

	return userInfo
}
