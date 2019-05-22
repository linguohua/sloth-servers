package mysql

// errCode = 1 参数错误， 2 已经在牌友群，3牌友群不存在
func removeUserFromClub(userID string, clubID string) (errCode int) {
	// dd_user_to_club`(IN userId varchar(64), IN clubId varchar(64))s

	stmt, err := dbConn.Prepare("call remove_user_from_club(?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	row := stmt.QueryRow(userID, clubID)
	err = row.Scan(&errCode)
	if err != nil {
		panic(err.Error())
	}

	return
}
