package mysql

// 检查手机号是否已经注册过
func countUserClub(userID string) (count int) {
	stmt, err := dbConn.Prepare("select count(*) from user_club where user_id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	row := stmt.QueryRow(userID)
	err = row.Scan(&count)
	if err != nil {
		panic(err.Error())
	}

	return
}
