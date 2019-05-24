package mysql

// 拉取用户的牌友群ID，user_club表，user_id和club_id都是索引，可以从user_id查club_id,也可以从club_id查user_id
func loadUserClubIDs(userID string) (clubIDs []string) {
	stmt, err := dbConn.Prepare("select club_id from user_club where user_id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	myClubIDs := make([]string, 0)
	rows, err := stmt.Query(userID)
	for rows.Next() {
		var clubID string
		err = rows.Scan(&clubID)
		if err != nil {
			panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
		}

		myClubIDs = append(myClubIDs, clubID)
	}

	return myClubIDs
}
