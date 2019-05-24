package mysql

import (
	"database/sql"
)

// 拉取牌友群信息
func loadClubIDByNumber(clubNumber string) string {
	stmt, err := dbConn.Prepare("select club_id from club_number where club_num = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	var myClubID sql.NullString

	rows, err := stmt.Query(clubNumber)
	if err != nil {
		panic(err.Error())
	}

	if !rows.Next() {
		return ""
	}

	err = rows.Scan(&myClubID)
	if err != nil {
		panic(err.Error())
	}

	if myClubID.Valid {
		return myClubID.String
	}

	return ""
}
