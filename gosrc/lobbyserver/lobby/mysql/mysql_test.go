package mysql
import (
	"database/sql"
	"testing"
	"log"
)


// TestSomething 测试用例
func TestSomething(t *testing.T) {
	dbConn = testDBConnect()
	defer dbConn.Close()

	testLoadUserIDByAccount()
	testLoadUserDiamond()
	testLoadPasswordByAccount()
}

func testDBConnect() *sql.DB{
	gameDBCon, err := newDbConnect("localhost", 3306, "localTest", "12345678", "game")
	if err != nil {
		log.Println("err:", err)
	}

	return gameDBCon

	// defer dbConn.Close()
}

func testLoadUserIDByAccount() {
	userID := loadUserIDByAccount("022008f3-d970-41e0-a9-8ca8e6283a4e")
	log.Println("userID:", userID)
}

func testLoadUserDiamond() {
	diamond := loadUserDiamond("11")
	log.Println("diamond:", diamond)
}

func testLoadPasswordByAccount() {
	password := loadPasswordByAccount("022008f3-d970-41e0-af39-8e6283a4e")
	log.Println("password:", password)
}