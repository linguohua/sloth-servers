package mysql

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

// 检查手机号是否已经注册过
func loadPasswordByAccount(account string) string {
	stmt, err := dbConn.Prepare("select password from account where account = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	var password sql.NullString
	row := stmt.QueryRow(account)
	err = row.Scan(&password)
	if err == sql.ErrNoRows {
		return ""
	}

	if err != nil {
		panic(err.Error())
	}

	if password.Valid {
		return password.String
	}

	log.Panicf("loadPasswordByAccount %s, can't convert password to string", account)

	return ""
}
