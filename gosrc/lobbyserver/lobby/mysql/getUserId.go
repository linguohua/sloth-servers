package mysql

import (
	"database/sql"
)

// 检查手机号是否已经注册过
func getUserIDBy(accout string) string {
	stmt, err := dbConn.Prepare("select user_id from account where account = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	var userID sql.NullString
	row := stmt.QueryRow(accout)
	err = row.Scan(&userID)
	if err != nil {
		panic(err.Error())
	}

	if userID.Valid {
		return userID.String
	}

	return ""
}
