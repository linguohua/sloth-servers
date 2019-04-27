package mysql

import (
	// "database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
)

// 检查手机号是否已经注册过
func getPasswordBy(accout string) string {
	query := fmt.Sprintf("select password from account where account = '%s'", accout)

	log.Println("query:", query)

	// return nil
	stmt, err := dbConn.Prepare(query)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	var password string
	row := stmt.QueryRow()
	err = row.Scan(&password)
	if err != nil {
		panic(err.Error())
	}

	return password;
}