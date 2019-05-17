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

func (*mySQLUtil) RegisterAccount(account string, passwd string, userInfo *lobby.UserInfo, clientInfo *lobby.ClientInfo) error {
	return registerAccount(account, passwd, userInfo, clientInfo)
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

func (*mySQLUtil) UpdateDiamond(userID string, change int64) (lastNum int64, errCode int32) {
	return updateDiamond(userID, change)
}

func (*mySQLUtil) CountUserClubNumber(userID string) (count int) {
	return countUserClub(userID)
}

func (*mySQLUtil) CreateClub(clubName string, creator string, isLeague int, wanka int, candy int, maxMember int) (clubID string, clubNumber string, errCode int) {
	return createClub(clubName, creator, isLeague, wanka, candy, maxMember)
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
