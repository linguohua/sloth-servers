package support

// 用于支持monkey账户修改、查询服务器
import (
	"fmt"
	"gconst"
	"lobbyserver/lobby"
	"net/http"

	//"webdata"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

// accSupportVerify 检查monkey用户接入合法
func accSupportVerify(w http.ResponseWriter, r *http.Request) bool {
	var account = r.URL.Query().Get("account")
	var password = r.URL.Query().Get("password")
	// log.Printf("monkey access, account:%s, password:%s\n", account, password)
	conn := lobby.Pool().Get()
	defer conn.Close()

	tableName := gconst.AccMonkeyAccountTalbe
	pass, e := redis.String(conn.Do("HGET", tableName, account))
	if e != nil || password != pass {
		return false
	}

	return true
}

func onQueryOnlineUser(w http.ResponseWriter, r *http.Request) {
	var resultString = `{"count":%d}`
	w.Write([]byte(fmt.Sprintf(resultString, lobby.UserCount())))
}

// // AddUser2Blacklist 添加用户到黑名单
// func AddUser2Blacklist(userID string) error {
// 	return addUser2Blacklist(userID)
// }

// // RemoveUserFromBlacklist 从黑名单中剔除用户
// func RemoveUserFromBlacklist(userID string) error {
// 	return removeUserFromBlacklist(userID)
// }

// // LoadBlacklist 加载黑名单列表
// func LoadBlacklist() []string {
// 	return loadBlacklist()
// }

// AddDiamond2User 给用户添加钻石
func AddDiamond2User(userID string, diamond int) (remainDiamond int, errCode int32) {
	// TODO: llwant mysql
	// remainDiamondInt64, err := webdata.ModifyDiamond(userID, addDiamondFromBackend, int64(diamond), "给用户添加钻石", 0, "", "")
	// if err != nil {
	// 	var errString = fmt.Sprintf("%v", err)
	// 	if strings.Contains(errString, diamondNotEnoughMsg) {
	// 		errCode = int32(gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough)
	// 	} else {
	// 		errCode = int32(gconst.SSMsgError_ErrTakeoffDiamondFailedIO)
	// 	}
	// 	return
	// }
	return int(0), 0
}

// QueryUserDiamond 查询用户钻石
func QueryUserDiamond(userID string) (remainDiamond int, errCode int32) {
	// TODO: llwant mysql
	// remainDiamondInt64, err := webdata.QueryDiamond(userID)
	// if err != nil {
	// 	errCode = int32(gconst.SSMsgError_ErrTakeoffDiamondFailedIO)
	// 	return
	// }
	return int(0), 0
}

func supportMiddleware(old http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if accSupportVerify(w, r) {
			old.ServeHTTP(w, r)
		} else {
			w.WriteHeader(404)
			w.Write([]byte("oh, no auth"))
		}
	})
}

// InitWith init
func InitWith(mainRouter *mux.Router) {
	var support = mainRouter.PathPrefix("/support").Subrouter()
	support.Use(supportMiddleware)

	support.HandleFunc("/onlineUser", onQueryOnlineUser)
}
