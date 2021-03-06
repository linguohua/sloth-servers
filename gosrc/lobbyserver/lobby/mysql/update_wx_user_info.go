package mysql

import (
	// "database/sql"
	"fmt"
	"lobbyserver/lobby"

	log "github.com/sirupsen/logrus"
)

// SaveGRCRecord2SqlServer 保存牌局记录到数据库
func updateWxUserInfo(userInfo *lobby.UserInfo, clientInfo *lobby.ClientInfo) error {
	/* Description:	更新微信用户信息
	`update_wx_user`(
		in userId int(11),
		in openId varchar(32),
		in userName varchar(64) ,
		in nickName varchar(32),
		in gender int(1) ,
		in provice varchar(32),
		in city varchar(32),
		in country varchar(32),
		in headImgUrl varchar(128),
		in modName varchar(32),
		in modVersion varchar(32),
		in coreVersion varchar(32),
		in lobbyVersion varchar(32),
		in operatingSystem varchar(32),
		in systemFamily varchar(32),
		in deviceId varchar(32),
		in deviceName varchar(32),
		in deviceMode varchar(32),
		in networkType varchar(32))
	*/

	// return nil
	stmt, err := dbConn.Prepare("Call update_wx_user(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("updateWxUserInfo Prepare1 Err:", err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		userInfo.GetUserID(),
		userInfo.GetOpenID(),
		userInfo.GetNickName(),
		userInfo.GetNickName(),
		userInfo.GetGender(),
		userInfo.GetProvince(),
		userInfo.GetCity(),
		userInfo.GetCountry(),
		userInfo.GetHeadImgUrl(),

		clientInfo.GetQMod(),
		clientInfo.GetModV(),
		clientInfo.GetCsVer(),
		clientInfo.GetLobbyVer(),
		clientInfo.GetOperatingSystem(),
		clientInfo.GetOperatingSystemFamily(),
		clientInfo.GetDeviceUniqueIdentifier(),
		clientInfo.GetDeviceName(),
		clientInfo.GetDeviceModel(),
		clientInfo.GetNetwork())
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(result)

	return nil

}
