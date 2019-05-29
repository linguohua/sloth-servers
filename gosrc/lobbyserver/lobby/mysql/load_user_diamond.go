package mysql

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

// 检查手机号是否已经注册过
func loadUserDiamond(userID string) int64 {
	stmt, err := dbConn.Prepare("select num from diamond where user_id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	row := stmt.QueryRow(userID)

	var diamond sql.NullInt64
	err = row.Scan(&diamond)
	if err == sql.ErrNoRows {
		return 0
	}

	if err != nil {
		panic(err.Error())
	}

	if diamond.Valid {
		return diamond.Int64
	}

	log.Panicf("loadUserDiamond, userID:%s, can't convert diamond to int64", userID)

	return 0
}
