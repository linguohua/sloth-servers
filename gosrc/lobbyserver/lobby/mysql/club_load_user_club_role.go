package mysql

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

// 拉取牌友群的成员ID，user_club表，user_id和club_id都是索引，可以从user_id查club_id,也可以从club_id查user_id
func loadUserClubRole(userID string, clubID string) (role int32) {
	log.Printf("loadUserClubRole, userID:%s, clubID:%s", userID, clubID)
	stmt, err := dbConn.Prepare("select club_role from user_club where user_id = ? and club_id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	var myRole sql.NullInt64
	rows, err := stmt.Query(userID, clubID)
	if err != nil {
		panic(err.Error())
	}

	if rows.Next() {
		err = rows.Scan(&myRole)
		if err != nil {
			panic(err.Error())
		}
	}

	if myRole.Valid {
		return int32(myRole.Int64)
	}

	return 0
}
