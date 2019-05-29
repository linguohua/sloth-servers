package mysql

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

// 检查手机号是否已经注册过
func loadUserIDByAccount(account string) string {
	stmt, err := dbConn.Prepare("select user_id from account where account = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	var userID sql.NullString
	row := stmt.QueryRow(account)
	err = row.Scan(&userID)
	if err == sql.ErrNoRows {
		return ""
	}

	if err != nil {
		panic(err.Error())
	}

	if userID.Valid {
		return userID.String
	}

	log.Panicf("loadUserIDByAccount %s, can't convert userID to string", account)

	return ""
}
