package mysql

import (
	// "database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"lobbyserver/lobby"
)

// SaveGRCRecord2SqlServer 保存牌局记录到数据库
func updateAccountUserInfo(phoneNum string, clientInfo *lobby.ClientInfo) error {
	/* Description:	更新微信用户信息
		`update_account_user`(
		in phoneNum int(11),
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
	query := fmt.Sprintf("Call update_account_user('%s', '%s', '%s', %s, '%s', '%s', '%s', '%s', '%s', '%s', '%s')",
		phoneNum,
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
