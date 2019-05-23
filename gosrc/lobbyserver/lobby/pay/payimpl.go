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

	errorPayParams          = 1
	errorPayNoEnoughDiamond = 2
	errorPayOrderNotExist   = 3
)

// Order 保存用户支付记录
type Order struct {
	OrderID string `json:"orderID"`
	Cost    int    `json:"cost"`
}

func getPayDiamondNum(payType int, playerNumAcquired int, handNum int, roomType int) (int, error) {
	log.Printf("getPayDiamondNum, payType:%d, playerNumAcquired:%d, handNum:%d, roomType:%d", payType, playerNumAcquired, handNum, roomType)

	var payKey = fmt.Sprintf("ownerPay:%d:%d", playerNumAcquired, handNum)
	if payType == AApay {
		payKey = fmt.Sprintf("aaPay:%d:%d", playerNumAcquired, handNum)
	} else if payType == GroupPay {
		payKey = fmt.Sprintf("GroupPay:%d:%d", playerNumAcquired, handNum)
	}

	// log.Println("payKey", payKey)
	cfg := pricecfg.GetPriceCfg(roomType)
	if cfg == nil {
		return 0, fmt.Errorf("Price config not exist")
	}

	// 必现要有原价表
	if cfg.OriginalPriceCfg == nil {
		return 0, fmt.Errorf("Original price config not exist")
	}

	if cfg.ActivityPriceCfg != nil && cfg.ActivityPriceCfg.DiscountCfg != nil && payType != FundPay {
		var discountCfg = cfg.ActivityPriceCfg.DiscountCfg
		pay, ok := discountCfg[payKey]
		if ok {
			log.Printf("discountCfg pay:%d, KEY:%s", pay, payKey)
			return pay, nil
		}

	}

	pay, ok := cfg.OriginalPriceCfg[payKey]
	if !ok {
		return 0, fmt.Errorf("OriginalPriceCfg %s not exist", payKey)
	}
	log.Printf("Original pay:%d, KEY:%s", pay, payKey)
	// log.Println(cfg.OriginalPriceCfg)
	return pay, nil
}

// doPayAndSave2RedisWith TODO：需要返回扣钱失败的具体原因，如余额不足，io失败等
func doPayWith(roomConfigID string, roomID string, userID string) (remainDiamond int, errCode int32) {
	log.Printf("payAndSave2RedisWith, roomConfigID:%s, roomID:%s, userID:%s", roomConfigID, roomID, userID)

	roomConfig := lobby.GetRoomConfig(roomConfigID)
	if roomConfig == nil {
		remainDiamond = 0
		errCode = int32(gconst.SSMsgError_ErrNoRoomConfig)
		return
	}

	pay, err := getPayDiamondNum(roomConfig.PayType, roomConfig.PlayerNumAcquired, roomConfig.HandNum, roomConfig.RoomType)
	if err != nil {
		log.Error("doPayAndSave2RedisWith, getPayDiamondNum error:", err)
		remainDiamond = 0
		errCode = 0
		return
	}

	mySQLUtil := lobby.MySQLUtil()
	result, lastNum, orderID := mySQLUtil.PayForRoom(userID, pay, roomID)
	if result == errorPayNoEnoughDiamond {
		log.Errorf("Pay no enough diamond, userID:%s, pay:%d, current diamond:%d", userID, pay, lastNum)
		errCode = int32(gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough)
	} else if result != 0 {
		log.Errorf("Pay for room failed, result:%d, userID:%s, pay:%d, roomID:%s", result, userID, pay, roomID)
		errCode = int32(gconst.SSMsgError_ErrTakeoffDiamondFailedIO)
	} else {
		errCode = 0
	}

	remainDiamond = int(lastNum)

	if result == 0 {
		order := &Order{}
		order.OrderID = orderID
		order.Cost = -pay

		buf, err := json.Marshal(order)
		if err != nil {
			log.Error("doPayWith, marsh order error:", err)
		}

		conn := lobby.Pool().Get()
		defer conn.Close()

		conn.Do("HSET", gconst.LobbyPayRoomPrefix+roomID, userID, buf)
	}

	return
}

func loadUsersInRoom(roomID string, conn redis.Conn) []string {
	// 获取房间内的用户列表
	vs, err := redis.Values(conn.Do("HGETALL", gconst.LobbyPayRoomPrefix+roomID))
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

func doPayForEnterRoom(roomID string, userID string) (remainDiamond int, errCode int32) {
	log.Printf("payAndSave2Redis, roomID:%s, userID:%s", roomID, userID)
	conn := lobby.Pool().Get()
	defer conn.Close()

	roomConfigID, err := redis.String(conn.Do("HGET", gconst.LobbyRoomTablePrefix+roomID, "roomConfigID"))
	if err != nil {
		log.Println("payAndSave2Redis, get roomConfigID failed, err:", err)
		remainDiamond = 0
		errCode = int32(gconst.SSMsgError_ErrNoRoomConfig)
		return
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
func refund2UserWith(roomID string, userID string, handFinish int) (remainDiamond int, errCode int32) {
	log.Printf("refund2UserAndSave2Redis, roomID:%s, userID:%s, handFinish:%d", roomID, userID, handFinish)
	conn := lobby.Pool().Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("HGET", gconst.LobbyPayRoomPrefix+roomID, userID)
	conn.Send("HGET", gconst.LobbyRoomTablePrefix+roomID, "config")
	vs, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		log.Error("refund2UserWith, load UnRefund order err:", err)
		return
	}

	var orderBuf, _ = redis.Bytes(vs[0], nil)
	if len(orderBuf) == 0 {
		log.Error("refund2UserWith, return failed, can't get order")
		return
	}

	order := &Order{}
	err = json.Unmarshal(orderBuf, order)
	if err != nil {
		log.Error("refund2UserWith, Unmarshal order error:", err)
		return
	}

	roomConfigID, err := redis.String(vs[0], nil)
	if err != nil {
		log.Error("refund2UserWith, return failed, can't get roomType ")
		return
	}

	roomConfig := lobby.GetRoomConfig(roomConfigID)
	if roomConfig == nil {
		log.Error("refund2UserWith, Can not find room config, roomConfigID:", roomConfigID)
		return
	}

	refund := getRetrunDiamond(order.Cost, roomConfig.HandNum, handFinish)
	if refund < 0 {
		log.Error("refund2UserAndSave2Redis, refundDiamond < 0")
		return
	}

	mySQLUtil := lobby.MySQLUtil()
	result, lastNum := mySQLUtil.RefundForRoom(userID, refund, order.OrderID)

	if result == 0 {
		conn.Do("HDEL", gconst.LobbyPayRoomPrefix+roomID, userID)
	}

	return int(lastNum), int32(result)
	// log.Errorf("refund2UserWith, result:%d, lastNum:%d", result, lastNum)
	// return false
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

		_, errCode := refund2UserWith(roomID, userID, finish)
		if errCode != 0 {
			result = false
			log.Errorf("refund2Users failed, errCode:%d, userID:%s, roomID:%s, finish:%d", errCode, userID, roomID, finish)
		}

	}

	return result
}
