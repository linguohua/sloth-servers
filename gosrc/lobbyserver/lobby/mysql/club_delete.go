package mysql

// 删除牌友群
func deleteClub(clubID string) (errCode int32) {
	// delete_club`(IN clubId varchar(64))
	stmt, err := dbConn.Prepare("call delete_club(?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	row := stmt.QueryRow(clubID)
	err = row.Scan(&errCode)
	if err != nil {
		panic(err.Error())
	}

	return
}
