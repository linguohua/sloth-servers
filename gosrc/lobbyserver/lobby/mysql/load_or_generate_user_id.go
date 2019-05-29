package mysql

import (
)

// 检查手机号是否已经注册过
func loadOrGenerateUserID(account string) (userID string, isNew bool) {
	/* 储存过程
	get_or_generate_user_id`(in account varchar(32), out userId int(11), out isNew boolean)
	*/

	stmt, err := dbConn.Prepare("Call get_or_generate_user_id(?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	row := stmt.QueryRow(account)
	err = row.Scan(&userID, &isNew)
	if err != nil {
		panic(err.Error())
	}

	return
}
