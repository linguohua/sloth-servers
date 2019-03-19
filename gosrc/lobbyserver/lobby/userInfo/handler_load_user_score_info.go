package userInfo

import (
	log "github.com/sirupsen/logrus"
	"gconst"
	"net/http"
	"github.com/garyburd/redigo/redis"
	//"github.com/golang/protobuf/proto"
)
// User 表示一个用户

// type UserScoreInfo struct {
// 	maxWinScore              *int32 // 约牌房单局最大得分
// 	customCount              *int32 // 约牌房总局数
// 	maxWinMoney              *int32 // 金币房单局最大赢取金币
// 	coinCount 				 *int32 // 金币房总局数
// }


func replyRequestUserScoreInfo(w http.ResponseWriter, errorCode int32, msgRequestUserScoreInfoRsp *MsgRequestUserScoreInfoRsp) {


	if errorCode != int32(MsgError_ErrSuccess){
		var errString = ErrorString[errorCode]
		msgRequestUserScoreInfoRsp.RetMsg = &errString
	}

	reply(w, msgRequestUserScoreInfoRsp, int32(MessageCode_OPRequestUserScoreInfo))
}

func handleLoadUserScoreInfo(w http.ResponseWriter, r *http.Request, userID string) {
	log.Println("handleLoadUserScoreInfo call, userID:", userID)
	// 1. 从请求中获取房间6位数字ID
	// 2. 根据用户ID读取出数值
	// 3. 返回结果給客户端

	conn := pool.Get()
	defer conn.Close()

	values, err := redis.Ints(conn.Do("HMGET", gconst.PlayerTablePrefix+userID, "dfHands", "dfHMW"))

	if err != nil && err != redis.ErrNil {
		log.Println("handleLoadUserScoreInfo get user score info err: ", err)
		replyRequestUserScoreInfo(w, int32(MsgError_ErrDatabase), nil)
		return
	}

	var customCountAddr = int32(values[0])
	var maxWinScoreAddr = int32(values[1])

	log.Println("customCountAddr ",  customCountAddr)
	log.Println("maxWinScoreAddr ",  maxWinScoreAddr)

	var msgRequestUserScoreInfoRsp = &MsgRequestUserScoreInfoRsp{}

	msgRequestUserScoreInfoRsp.MaxWinScore = &maxWinScoreAddr
	msgRequestUserScoreInfoRsp.CustomCount = &customCountAddr

	replyRequestUserScoreInfo(w, int32(MsgError_ErrSuccess),msgRequestUserScoreInfoRsp)

}
