package mysql

func createClub(clubName string, creator string, isLeague int, wanka int, candy int, maxMember int) (clubID string, clubNumber string, errCode int) {
	// create_club`(IN clubName VARCHAR(32), IN creator VARCHAR(64), IN isLeague INT, IN wanka INT(11), IN candy INT(11), IN maxMember INT(4)

	stmt, err := dbConn.Prepare("call create_club(?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	row := stmt.QueryRow(clubName, creator, isLeague, wanka, candy, maxMember)
	err = row.Scan(&clubID, &clubNumber, &errCode)
	if err != nil {
		panic(err.Error())
	}

	return
}
