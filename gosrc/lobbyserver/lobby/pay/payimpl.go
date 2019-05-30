package pay

import (
	"fmt"
	"gconst"
	"lobbyserver/lobby"
	"lobbyserver/pricecfg"
	"math"

	log "github.com/sirupsen/logrus"

	"encoding/json"

	"github.com/garyburd/redigo/redis"
)

const (
	dbServerID  = "8def07dc-a53f-4851-a88d-9d45d7db126a"
	dbServerURL = "tservice.dafeng.xy.qianz.com:9101"
)

const (
	// AApay aa pay
	AApay = 1
	// FundPay fund pay
	FundPay = 2
	// GroupPay group pay
	GroupPay = 2

	errorParams             = 1
	errorPayNoEnoughDiamond = 2
	errorPayOrderNotExist   = 3
	errorRefundParams       = 1
	errorRefundRepeat       = 2
)

// Order 保存用户支付记录
type Order struct {
	OrderID string `json:"orderID"`
	// 若创建房间失败，则roomConfigID不会保存到room表，因此保存到订单里面
	RoomConfigID string `json:"roomConfigID"`
	Cost         int    `json:"cost"`
}

func getPayDiamondNum(payType int, playerNumAcquired int, handNum int, roomType int) (pay int, errCode int32) {
	log.Printf("getPayDiamondNum, payType:%d, playerNumAcquired:%d, handNum:%d, roomType:%d", payType, playerNumAcquired, handNum, roomType)

	var payKey = fmt.Sprintf("ownerPay:%d:%d", playerNumAcquired, handNum)
	if payType == AApay {
		payKey = fmt.Sprintf("aaPay:%d:%d", playerNumAcquired, handNum)
	} else if payType == GroupPay {
		payKey = fmt.Sprintf("GroupPay:%d:%d", playerNumAcquired, handNum)
	}

	cfg := pricecfg.GetPriceCfg(roomType)
	if cfg == nil {
		log.Error("getPayDiamondNum, no priceCfg found for roomType:", roomType)
		return 0, int32(lobby.MsgError_ErrRoomPriceCfgNotExist)
	}

	// 必现要有原价表
	if cfg.OriginalPriceCfg == nil {
		log.Errorf("getPayDiamondNum, roomType %d OriginalPriceCfg is nil", roomType)
		return 0, int32(lobby.MsgError_ErrRoomPriceCfgNotExist)
	}

	if cfg.ActivityPriceCfg != nil && cfg.ActivityPriceCfg.DiscountCfg != nil && payType != FundPay {
		var discountCfg = cfg.ActivityPriceCfg.DiscountCfg
		pay, ok := discountCfg[payKey]
		if ok {
			return pay, 0
		}

	}

	pay, ok := cfg.OriginalPriceCfg[payKey]
	if !ok {
		log.Errorf("getPayDiamondNum, no price found for pay key %s", payKey)
		return 0, int32(lobby.MsgError_ErrRoomPriceCfgNotExist)
	}

	return pay, 0
}

func saveOrder(userID string, roomID string, remainDiamond int64, order *Order) {
	buf, err := json.Marshal(order)
	if err != nil {
		log.Error("doPayWith, marsh order error:", err)
	}

	key := fmt.Sprintf("%s%s", gconst.LobbyUserTablePrefix, userID)

	conn := lobby.Pool().Get()
	defer conn.Close()

	// 更新用户钻石
	// 保存订单
	conn.Send("MULTI")
	conn.Send("HSET", key, "diamond", remainDiamond)
	conn.Send("HSET", gconst.LobbyPayOrderPrefix+roomID, userID, buf)
	_, err = conn.Do("EXEC")
	if err != nil {
		log.Println("saveOrder err: ", err)
	}
}

// doPayAndSave2RedisWith TODO：需要返回扣钱失败的具体原因，如余额不足，io失败等
func doPayWith(roomConfigID string, roomID string, userID string) (errCode int32) {
	log.Printf("payAndSave2RedisWith, roomConfigID:%s, roomID:%s, userID:%s", roomConfigID, roomID, userID)

	roomConfig := lobby.GetRoomConfig(roomConfigID)
	if roomConfig == nil {
		// return int32(gconst.SSMsgError_ErrNoRoomConfig)
		return int32(lobby.MsgError_ErrNoRoomConfig)
	}

	pay, errCode := getPayDiamondNum(roomConfig.PayType, roomConfig.PlayerNumAcquired, roomConfig.HandNum, roomConfig.RoomType)
	if errCode != 0 {
		return errCode
	}

	mySQLUtil := lobby.MySQLUtil()
	result, lastNum, orderID := mySQLUtil.PayForRoom(userID, pay, roomID)
	if result == errorParams {
		log.Panicf("Pay error params, userID:%s, pay:%d, roomID:%s", userID, pay, roomID)
	} else if result == errorPayNoEnoughDiamond {
		log.Errorf("Pay no enough diamond, userID:%s, pay:%d, current diamond:%d", userID, pay, lastNum)
		return int32(lobby.MsgError_ErrTakeoffDiamondFailedNotEnough)
	}

	order := &Order{}
	order.OrderID = orderID
	order.Cost = -pay
	order.RoomConfigID = roomConfigID

	saveOrder(userID, roomID, lastNum, order)

	sessionMgr := lobby.SessionMgr()
	sessionMgr.UpdateUserDiamond(userID, uint64(lastNum))

	return 0
}

func loadUsersInRoom(roomID string, conn redis.Conn) []string {
	// 获取房间内的用户列表
	vs, err := redis.Values(conn.Do("HGETALL", gconst.LobbyPayOrderPrefix+roomID))
	if err != nil {
		log.Println("readUserIDListInRoom, get room players failed:", err)
		return []string{}
	}

	userIDs := make([]string, 0)
	for i := 0; i < len(vs); i = i + 2 {
		userID, _ := redis.String(vs[i], nil)
		userIDs = append(userIDs, userID)
	}

	return userIDs
}

func doPayForEnterRoom(roomID string, userID string) (errCode int32) {
	log.Printf("payAndSave2Redis, roomID:%s, userID:%s", roomID, userID)
	conn := lobby.Pool().Get()
	defer conn.Close()

	roomConfigID, err := redis.String(conn.Do("HGET", gconst.LobbyRoomTablePrefix+roomID, "roomConfigID"))
	if err != nil {
		log.Println("payAndSave2Redis, get roomConfigID failed, err:", err)
		return int32(gconst.SSMsgError_ErrNoRoomConfig)
	}

	return doPayWith(roomConfigID, roomID, userID)

}

func getRetrunDiamond(pay int, handNum int, handFinish int) int {
	log.Printf("getRetrunDiamond, pay:%d, handNum:%d, handFinish:%d", pay, handNum, handFinish)

	var unPlayHand = handNum - handFinish
	var refundDiamond = float64(unPlayHand) / float64(handNum) * float64(pay)
	return int(math.Floor(refundDiamond))
}

// refund2UserAndSave2Redis refund money to user
func refund2UserWith(roomID string, userID string, handFinish int) (errCode int32) {
	log.Printf("refund2UserAndSave2Redis, roomID:%s, userID:%s, handFinish:%d", roomID, userID, handFinish)
	conn := lobby.Pool().Get()
	defer conn.Close()

	orderBuf, err := redis.Bytes(conn.Do("HGET", gconst.LobbyPayOrderPrefix+roomID, userID))
	if err != nil {
		log.Error("refund2UserWith, load order from redis err:", err)
		return 0
	}

	order := &Order{}
	err = json.Unmarshal(orderBuf, order)
	if err != nil {
		log.Error("refund2UserWith, Unmarshal order error:", err)
		return 0
	}

	roomConfig := lobby.GetRoomConfig(order.RoomConfigID)
	if roomConfig == nil {
		log.Error("refund2UserWith, Can not find room config, roomConfigID:", order.RoomConfigID)
		return 0
	}

	refund := getRetrunDiamond(order.Cost, roomConfig.HandNum, handFinish)
	if refund < 0 {
		log.Error("refund2UserAndSave2Redis, refundDiamond < 0")
		return
	}

	mySQLUtil := lobby.MySQLUtil()
	result, lastNum := mySQLUtil.RefundForRoom(userID, refund, order.OrderID)
	if result == errorParams {
		log.Panicf("RefundForRoom failed, error params, userID:%s, refund:%d, orderID:%s", userID, refund, order.OrderID)
	}

	if result == 0 {
		// 更新用户钻石
		// 保存订单
		key := fmt.Sprintf("%s%s", gconst.LobbyUserTablePrefix, userID)

		conn.Send("MULTI")
		conn.Send("HSET", key, "diamond", lastNum)
		conn.Send("HDEL", gconst.LobbyPayOrderPrefix+roomID, userID)
		_, err = conn.Do("EXEC")
		if err != nil {
			log.Println("refund2UserWith, redis err: ", err)
		}

		sessionMgr := lobby.SessionMgr()
		sessionMgr.UpdateUserDiamond(userID, uint64(lastNum))
	}

	return int32(result)
}

func isUserExist(userID string, userIDs []string) bool {
	for _, id := range userIDs {
		if userID == id {
			return true
		}
	}

	return false
}

// refund2Users inGameUserIDs 是房间里面游戏的用户，不是要返还钻石的用户，只用检查作返还的用户是否是游戏里面的用户
func refund2Users(roomID string, handFinish int, inGameUserIDs []string) bool {
	log.Printf("refund, roomID:%s,  handFinish:%d", roomID, handFinish)
	conn := lobby.Pool().Get()
	defer conn.Close()

	var userIDs = loadUsersInRoom(roomID, conn)
	log.Println("userIDs:", userIDs)
	result := true
	for _, userID := range userIDs {
		var finish = handFinish
		// 如果用户已经离开了房间，还没还钱，则全部返还给用户
		if len(inGameUserIDs) > 0 && !isUserExist(userID, inGameUserIDs) {
			finish = 0
		}

		errCode := refund2UserWith(roomID, userID, finish)
		if errCode != 0 {
			result = false
			log.Errorf("refund2Users failed, errCode:%d, userID:%s, roomID:%s, finish:%d", errCode, userID, roomID, finish)
		}

	}

	return result
}
