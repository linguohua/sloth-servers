package mysql

import (
	log "github.com/sirupsen/logrus"
)

// 开房扣钱
func payForRoom(userID string, pay int, roomID string) (errCode int, lastNum int64, orderID string) {
	/* 储存过程
	pay_for_room`(<{IN userId INT(11)}>, <{IN pay INT(11)}>, <{IN roomId varchar(32)}>, <{IN roomNum varchar(6)}>);
	*/
	log.Printf("userID:%s, pay:%d, roomID:%s", userID, pay, roomID)
	stmt, err := dbConn.Prepare("Call pay_for_room(?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	row := stmt.QueryRow(userID, pay, roomID)

	// errCode 1 参数pay不能小于0，2 钻石不足
	err = row.Scan(&orderID, &lastNum, &errCode)
	if err != nil {
		panic(err.Error())
	}

	return
}
