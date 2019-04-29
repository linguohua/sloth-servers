package mysql

import (
	// "database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"lobbyserver/lobby"
)

// SaveGRCRecord2SqlServer 保存牌局记录到数据库
func updateWxUserInfo(userInfo *lobby.UserInfo, clientInfo *lobby.ClientInfo) error {
	/* Description:	更新微信用户信息
	`update_wx_user`(
		in userId int(11),
		in openId varchar(32),
		in userName varchar(64) ,
		in nickName varchar(32),
		in sex int(1) ,
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
	query := fmt.Sprintf("Call update_wx_user('%d', '%s', '%s', '%s', %d, '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')",
		userInfo.GetUserID(),
		userInfo.GetOpenID(),
		userInfo.GetNickName(),
		userInfo.GetNickName(),
		userInfo.GetSex(),
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

	log.Println("query:", query)

	// return nil
	stmt, err := dbConn.Prepare(query)
	if err != nil {
		log.Println("updateWxUserInfo Prepare1 Err:", err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(result)

	return nil

}
