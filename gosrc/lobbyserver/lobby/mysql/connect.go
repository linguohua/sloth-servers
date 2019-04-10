package mysql

import (
	"fmt"
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql" //不能去掉，不然连接数据库的时候提示找不到mssql
)

func newDbConnect(ip string, port int, user string, password string, database string) (*sql.DB, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, ip, port, database)

	fmt.Printf("mysql connString:%s\n", connString)
	dbCon, err := sql.Open("mysql", connString)
	if err != nil {
		log.Println("Open mssql connection failed:", err.Error())
		return nil, err
	}

	err = dbCon.Ping()
	if err != nil {
		log.Println(database, "Cannot ping: ", err.Error())
		return nil, err
	}

	return dbCon, nil
}

// StartMssql 启动mssql,只能调用一次
func startMySQL(ip string, port int, user string, password string, gameDB string) (*sql.DB, error) {
	gameDBCon, err := newDbConnect(ip, port, user, password, gameDB)

	return gameDBCon, err
}
