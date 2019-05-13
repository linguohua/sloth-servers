package mysql

// 解散房间还钱
func refundForRoom(userID string, refund int, orderID string) (errCode int, lastNum int64) {
	/* 储存过程
	refund_for_room`(IN orderId varchar(32), IN userId INT(11), IN refund INT(11));
	*/

	stmt, err := dbConn.Prepare("Call refund_for_room(?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	row := stmt.QueryRow(orderID, userID, refund)

	// errCode 1 参数refund 不能小于1，并且不能大于扣的钱
	err = row.Scan(&lastNum, &errCode)
	if err != nil {
		panic(err.Error())
	}

	return
}
