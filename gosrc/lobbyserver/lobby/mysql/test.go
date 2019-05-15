package mysql


import (
	// "database/sql"
	// "fmt"
	"lobbyserver/lobby"
	log "github.com/sirupsen/logrus"
)

func testGetOrGenerateUserID(account string) {
	userID, isNew := getOrGenerateUserID(account)
	log.Printf("userID:%s, isNew:%v", userID, isNew)
}

func testGetPassword(account string) {
	psw := getPasswordBy(account)
	log.Printf("password:%s", psw)
}

func testGetUserID(account string) {
	userID := getUserIDBy(account)
	log.Println("userID:", userID)
}

func testGetUserInfo(userID string) {
	userInfo := loadUserInfo(userID)
	log.Println("userInfo:", userInfo)
}

func testRegisterAccount() {
	clientInfo := &lobby.ClientInfo{}
	qMod := "qMod"
	modV := "modV"
	clientInfo.QMod = &qMod
	clientInfo.ModV = &modV

	userInfo := &lobby.UserInfo{}
	userID := "10000022"
	userInfo.UserID = &userID
	registerAccount("abc", "111111", userInfo, clientInfo)
}

func test() {
	// testGetOrGenerateUserID("abcddc")
	// testGetPassword("aa")
	// testGetUserID("aa")
	// testGetUserInfo("10000007")
	testRegisterAccount()

}