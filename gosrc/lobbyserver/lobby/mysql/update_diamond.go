package mysql

// 更新钻石
func updateDiamond(userID string, changeNum int64) (lastNum int64, errCode int32) {
	/* 储存过程
	update_diamond`(IN userId INT(11), IN changeNum INT(11))
	*/

	stmt, err := dbConn.Prepare("Call update_diamond(?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	row := stmt.QueryRow(userID, changeNum)

	err = row.Scan(&lastNum, &errCode)
	if err != nil {
		panic(err.Error())
	}

	return
}
