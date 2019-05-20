package mysql
import (
	"testing"
	"log"
)


// TestSomething 测试用例
func TestSomething(t *testing.T) {
	testDBConnect()

	testLoadClubInfoByNumber("237910")
}

func testDBConnect() {
	gameDBCon, err := newDbConnect("localhost", 3306, "localTest", "12345678", "game")
	if err != nil {
		log.Println("err:", err)
	}

	dbConn = gameDBCon

	// defer dbConn.Close()
}


func testLoadClubInfoByNumber(number string) {
	clubInfo := loadClubInfoByNumber(number)
	log.Println("clubInfo:", clubInfo)
}