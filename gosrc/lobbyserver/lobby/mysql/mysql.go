package mysql

import (
	"lobbyserver/lobby"
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

func (*mySQLUtil) StartMySQL(ip string, port int, user string, password string, gameDB string) {
	conn, err := startMySQL(ip, port, user, password, gameDB)
	if err != nil {
		log.Println("StartMssql error ", err)
	}
	dbConn = conn
}

// InitWith init
func InitWith() {
	lobby.SetMySQLUtil(sqlUtil)
	ip := "127.0.0.1"
	port := 3306
	user := "root"
	password := "123456"
	gameDB := "test"

	sqlUtil.StartMySQL(ip, port, user, password, gameDB)
}