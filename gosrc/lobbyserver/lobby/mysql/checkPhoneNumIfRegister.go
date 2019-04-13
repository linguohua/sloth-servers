package mysql

import (
	// "database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
)

// 检查手机号是否已经注册过
func checkPhoneNumIfRegister(phoneNum string) bool {
	query := fmt.Sprintf("select exists(select * from phone_account where phone_num = %s)", phoneNum)

	log.Println("query:", query)

	// return nil
	stmt, err := dbConn.Prepare(query)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	var result int
	row := stmt.QueryRow()
	err = row.Scan(&result)
	if err != nil {
		panic(err.Error())
	}

	if result == 1 {
		return true
	}

	return false
}