package lobby

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func handleUpdateMoney(w http.ResponseWriter, r *http.Request) {
	log.Println("user offline")
	var userID = r.URL.Query().Get("userID")
	if userID == "" {
		w.WriteHeader(404)
		w.Write([]byte("User id is empty !"))
		return
	}
	// 更新用户钻石
	// user := userMgr.getUserByID(userID)
	// if user == nil {
	// 	log.Println("user offline")
	// 	w.Write([]byte("User offline !"))
	// 	return
	// }

	// TODO: llwant mysql
	// diamond, err := webdata.QueryDiamond(userID)
	// if err != nil {
	// 	var msg = fmt.Sprintf("Query user %s diamond failed: %v", userID, err)
	// 	log.Println(msg)
	// 	w.Write([]byte(msg))
	// 	return
	// }

	// user.updateMoney(uint32(0))

	w.Write([]byte("Update user money success !"))
}
