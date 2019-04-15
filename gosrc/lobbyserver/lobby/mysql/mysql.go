package mysql

import (
	"lobbyserver/lobby"
	"lobbyserver/config"
	"database/sql"
	"log"
)

var (
	sqlUtil = &mySQLUtil{}
	dbConn *sql.DB
)

// myRoomUtil implements IRoomUtil
type mySQLUtil struct {

}

func (*mySQLUtil) UpdateWxUserInfo(wxUserInfo *lobby.WxUserInfo, clientInfo *lobby.ClientInfo) error {
	return updateWxUserInfo(wxUserInfo, clientInfo)
}

func (*mySQLUtil) UpdateAccountUserInfo(account string, clientInfo *lobby.ClientInfo) error {
	return updateAccountUserInfo(account, clientInfo)
}

func (*mySQLUtil) GetUserIDBy(account string) int {
	return getUserIDBy(account)
}

func (*mySQLUtil) GetOrGenerateUserID(account string) (userID string, isNew bool) {
	return getOrGenerateUserID(account)
}

// InitWith init
func InitWith() {
	lobby.SetMySQLUtil(sqlUtil)

	conn, err := startMySQL(config.DbIP, config.DbPort, config.DbUser, config.DbPassword, config.DbName)
	if err != nil {
		log.Println("StartMssql error ", err)
	}
	dbConn = conn
}