package mysql

import (
	// "database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
)

// 检查手机号是否已经注册过
func getOrGenerateUserID(account string) (userID uint64, isNew bool) {
	/* 储存过程
	get_or_generate_user_id`(in account varchar(32), out userId int(11), out isNew boolean)
	*/

	query := fmt.Sprintf("Call get_or_generate_user_id('%s', @out_userId, @out_isNew)", account)

	log.Println("query:", query)

	// return nil
	stmt, err := dbConn.Prepare(query)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	result, err := stmt.Exec()
    if err != nil {
        panic(err.Error())
	}

	fmt.Println(result)


	query = "SELECT @out_userId, @out_isNew;"
	stmt, err = dbConn.Prepare(query)
	if err != nil {
		panic(err.Error())
	}

	row := stmt.QueryRow()
	err = row.Scan(&userID, &isNew)
	if err != nil {
		panic(err.Error())
	}

	return
}