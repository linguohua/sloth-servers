package mysql

// 拉取牌友群的成员ID，user_club表，user_id和club_id都是索引，可以从user_id查club_id,也可以从club_id查user_id
func loadClubUserIDs(clubID string) (userIDs []string) {
	stmt, err := dbConn.Prepare("select user_id from user_club where club_id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	myUserIDs := make([]string, 0)
	rows, err := stmt.Query(clubID)
	for rows.Next() {
		var userID string
		err = rows.Scan(&userID)
		if err != nil {
			panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		}

		myUserIDs = append(myUserIDs, userID)
	}

	return myUserIDs
}
