package mysql

import (
	// "database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"lobbyserver/lobby"
)

// SaveGRCRecord2SqlServer 保存牌局记录到数据库
func registerAccount(account string, passwd string, phone string, userInfo *lobby.UserInfo, clientInfo *lobby.ClientInfo) error {
	/* Description:	更新微信用户信息
		`register_account`(
		in userId varchar(64),
		in acc varchar(64),
		in passwd varchar(64),
		in phone varchar(11),
		in openId varchar(32),
		in nickName varchar(64),
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

	// log.Println("query:", query)

	// return nil
	stmt, err := dbConn.Prepare("Call register_account(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("updateWxUserInfo Prepare1 Err:", err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		 userInfo.GetUserID(),
		 account,
		 passwd,
		 phone,

		 userInfo.GetOpenID(),
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
	if err != nil {
		return err
	}
	fmt.Println(result)

	return nil

}
