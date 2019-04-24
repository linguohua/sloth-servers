package mysql

import (
	// "database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
)

// 检查手机号是否已经注册过
func getUserIDBy(accout string) uint64 {
	query := fmt.Sprintf("select user_id from account where account = '%s'", accout)

	log.Println("query:", query)

	// return nil
	stmt, err := dbConn.Prepare(query)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	var userID uint64
	row := stmt.QueryRow()
	err = row.Scan(&userID)
	if err != nil {
		panic(err.Error())
	}

	return userID;
}