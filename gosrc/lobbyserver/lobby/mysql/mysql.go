package mysql

import (
	"database/sql"
	"lobbyserver/lobby"

	log "github.com/sirupsen/logrus"
)

var (
	sqlUtil = &mySQLUtil{}
	dbConn  *sql.DB
)

// myRoomUtil implements IRoomUtil
type mySQLUtil struct {
}

func (*mySQLUtil) UpdateWxUserInfo(userInfo *lobby.UserInfo, clientInfo *lobby.ClientInfo) error {
	return updateWxUserInfo(userInfo, clientInfo)
}

func (*mySQLUtil) UpdateAccountUserInfo(account string, clientInfo *lobby.ClientInfo) error {
	return updateAccountUserInfo(account, clientInfo)
}

func (*mySQLUtil) GetUserIDBy(account string) string {
	return getUserIDBy(account)
}

func (*mySQLUtil) GetPasswordBy(account string) string {
	return getPasswordBy(account)
}

func (*mySQLUtil) GetOrGenerateUserID(account string) (userID string, isNew bool) {
	return getOrGenerateUserID(account)
}

func (*mySQLUtil) RegisterAccount(account string, passwd string, phone string, userInfo *lobby.UserInfo, clientInfo *lobby.ClientInfo) error {
	return registerAccount(account, passwd, phone, userInfo, clientInfo)
}

func (*mySQLUtil) LoadUserInfo(userID string) *lobby.UserInfo {
	return loadUserInfo(userID)
}

func (*mySQLUtil) PayForRoom(userID string, pay int, roomID string) (errCode int, lastNum int64, orderID string) {
	return payForRoom(userID, pay, roomID)
}

func (*mySQLUtil) RefundForRoom(userID string, refund int, orderID string) (errCode int, lastNum int64) {
	return refundForRoom(userID, refund, orderID)
}

// InitWith init
func InitWith() {
	lobby.SetMySQLUtil(sqlUtil)

	conn, err := startMySQL()
	if err != nil {
		// log.Panic("StartMssql error ", err)
		log.Warn("StartMssql error ", err)
	}
	dbConn = conn

	test()
}
