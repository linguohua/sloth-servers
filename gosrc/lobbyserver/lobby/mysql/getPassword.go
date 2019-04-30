package mysql

import (
	"database/sql"
)

// 检查手机号是否已经注册过
func getPasswordBy(account string) string {
	stmt, err := dbConn.Prepare("select password from account where account = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	var password sql.NullString
	row := stmt.QueryRow(account)
	err = row.Scan(&password)
	if err != nil {
		panic(err.Error())
	}

	if password.Valid {
		return password.String
	}

	return ""
}
