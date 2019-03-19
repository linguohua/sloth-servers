package pay

import (
	"encoding/json"
	"fmt"
	"gconst"
	"lobbyserver/pricecfg"
	"math"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/garyburd/redigo/redis"
	proto "github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
)

const (
	dbServerID  = "8def07dc-a53f-4851-a88d-9d45d7db126a"
	dbServerURL = "tservice.dafeng.xy.qianz.com:9101"
)

// UserCost 保存用户支付记录
type UserCost struct {
	Cost          int   `json:"cost"`
	Result        int32 `json:"resultCode"`
	RemainDiamond int   `json:"remainDiamond"` // 只有扣钱成功，返回的剩余钻石才有效
	TimeStamp     int64 `json:"timeStamp"`
}

// UserRefund 保存还钱记录
type UserRefund struct {
	Refund        int   `json:"refund"`
	Result        int32 `json:"resultCode"`
	RemainDiamond int   `json:"remainDiamond"` // 只有退款成功，返回的剩余钻石才有效
	TimeStamp     int64 `json:"timeStamp"`
}

// OrderRecord 保存一个完整的订单
type OrderRecord struct {
	UserID     string      `json:"userID"`
	RoomID     string      `json:"roomID"`
	PayType    int         `json:"payType"`
	HandNum    int         `json:"handNum"`
	HandFinish int         `json:"handFinish"`
	PlayerNum  int         `json:"playerNum"`
	Cost       *UserCost   `json:"cost"`
	Refund     *UserRefund `json:"refund"`
}

func loadGameNoAndGroupID(roomID string) (gameNo string, groupID string) {
	groupRoomInfo := groupRoomInfoMap[roomID]
	if groupRoomInfo != nil {
		gameNo = fmt.Sprintf("%d", groupRoomInfo.GameNo)
		groupID = groupRoomInfo.ClubID
	}
	return
}

func getPayDiamondNum(payType int, playerNumAcquired int, handNum int, roomType int) (int, error) {
	log.Printf("getPayDiamondNum, payType:%d, playerNumAcquired:%d, handNum:%d, roomType:%d", payType, playerNumAcquired, handNum, roomType)

	var payKey = fmt.Sprintf("ownerPay:%d:%d", playerNumAcquired, handNum)
	if payType == aapay {
		payKey = fmt.Sprintf("aaPay:%d:%d", playerNumAcquired, handNum)
	} else if payType == groupPay {
		payKey = fmt.Sprintf("groupPay:%d:%d", playerNumAcquired, handNum)
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

	if cfg.ActivityPriceCfg != nil && cfg.ActivityPriceCfg.DiscountCfg != nil && payType != fundPay {
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
	log.Println(cfg.OriginalPriceCfg)
	return pay, nil

	// payCfg, ok := roomPriceCfg[roomType]
	// if !ok {
	// 	return 0
	// }

	// return payCfg[payKey]
}

// 判断用户是否支付过
// 用户对于同个房间只能有一个正在处理的订单
func isUserHavePay(roomID string, userID string) bool {
	log.Printf("isUserHavePay, roomID:%s, userID:%s", roomID, userID)
	conn := pool.Get()
	defer conn.Close()

	exist, _ := redis.Int(conn.Do("HEXISTS", gconst.RoomUnRefund+userID, roomID))
	if exist == 1 {
		return true
	}

	return false
}

// TODO：需要返回扣钱失败的具体原因，如余额不足，io失败等
func payAndSave2RedisWith(roomType int, roomConfigID string, roomID string, userID string, gameNo string) (remainDiamond int, errCode int32) {
	log.Printf("payAndSave2RedisWith, roomType:%d, roomConfigID:%s, roomID:%s, userID:%s", roomType, roomConfigID, roomID, userID)
	// 如果用户已经支付过，则不用再次支付
	var isPay = isUserHavePay(roomID, userID)
	if isPay {
		log.Printf("payAndSave2RedisWith, user %s have been pay for room %s", userID, roomID)
		remainDiamond = 0
		errCode = int32(gconst.SSMsgError_ErrTakeoffDiamondFailedRepeat)
		return
	}

	roomConfigString, ok := roomConfigs[roomConfigID]
	if !ok {
		remainDiamond = 0
		errCode = int32(gconst.SSMsgError_ErrNoRoomConfig)
		return
	}

	roomConfig := parseRoomConfigFromString(roomConfigString)

	// var pay = 0
	// if roomType != int(gconst.RoomType_GuanDang) {
	// 	needPay, err := getPayDiamondNum(roomConfig.PayType, roomConfig.PlayerNumAcquired, roomConfig.HandNum, roomType)
	// 	if err != nil {
	// 		log.Println("getPayDiamondNum error:", err)
	// 		return
	// 	}
	// 	pay = needPay
	// }

	// var subGameID = 0 // getSubGameIDByRoomType(roomType)

	var groupID = ""
	if roomConfig.IsGroup {
		_, groupID = loadGameNoAndGroupID(roomID)
		log.Printf("payAndSave2RedisWith, group pay,,groupID:%s", groupID)
	}

	// 扣钻类型
	var modDiamondType = ownerModDiamondCreateRoom
	if roomConfig.PayType == aapay {
		modDiamondType = aaModDiamondCreateRoom
		if roomConfig.IsGroup {
			modDiamondType = aaModDiamondCreateRoomForGroup
		}
	} else if roomConfig.PayType == groupPay {
		modDiamondType = masterModDiamondCreateRoomForGroup
	} else {
		if roomConfig.IsGroup {
			modDiamondType = ownerModDiamondCreateRoomForGroup
		}
	}

	log.Println("payAndSave2RedisWith modDiamondType:", modDiamondType)

	var result int32
	// money, err := webdata.ModifyDiamond(userID, modDiamondType, int64(-pay), "创建房间扣钱", subGameID, gameNo, groupID)
	// if err != nil {
	// 	log.Println("ModifyDiamond err:", err)
	// 	result = int32(gconst.SSMsgError_ErrTakeoffDiamondFailedIO)
	// 	var errString = fmt.Sprintf("%v", err)
	// 	if strings.Contains(errString, diamondNotEnoughMsg) {
	// 		result = int32(gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough)
	// 	}
	// } else {
	// 	result = int32(gconst.SSMsgError_ErrSuccess)
	// }

	// savePay2Redis(roomConfig, roomID, userID, pay, int(money), result)

	remainDiamond = int(0)
	errCode = result

	return
}

func loadUsersInRoom(roomID string, conn redis.Conn) []string {
	// 获取房间内的用户列表
	buf, err := redis.Bytes(conn.Do("HGET", gconst.RoomPayUsers+roomID, "users"))
	if err != nil {
		log.Println("readUserIDListInRoom, get room players failed:", err)
		return []string{}
	}

	userIDList := &gconst.SSMsgUserIDList{}
	err = proto.Unmarshal(buf, userIDList)
	if err != nil {
		log.Println("readUserIDListInRoom, unmarshal failed:", err)
		return []string{}
	}

	return userIDList.GetUserIDs()
}

// 生成一个新的订单保存到redis
func savePay2Redis(roomConfig *RoomConfigJSON, roomID string, userID string, cost int, remainDiamond int, result int32) {
	log.Printf("savePay2Redis,payType:%d, roomID:%s, userID:%s, cost:%d, remainDiamond:%d, result:%d", roomConfig.PayType, roomID, userID, cost, remainDiamond, result)

	conn := pool.Get()
	defer conn.Close()

	var userIDs = loadUsersInRoom(roomID, conn)
	userIDs = append(userIDs, userID)

	var userIDList = &gconst.SSMsgUserIDList{}
	userIDList.UserIDs = userIDs
	userIDListBuf, err := proto.Marshal(userIDList)
	if err != nil {
		log.Panicln("savePay2Redis, Marshal userIDListBuf err: ", err)
	}

	log.Println("savePay2Redis, userIDList:", userIDList)

	t := time.Now().UTC()

	var uCost = &UserCost{}
	uCost.Cost = cost
	uCost.Result = result
	uCost.RemainDiamond = remainDiamond
	uCost.TimeStamp = t.UnixNano()

	var orderRecord = &OrderRecord{}
	orderRecord.UserID = userID
	orderRecord.Cost = uCost
	orderRecord.RoomID = roomID
	orderRecord.PayType = roomConfig.PayType
	orderRecord.HandNum = roomConfig.HandNum
	orderRecord.HandFinish = 0
	orderRecord.PlayerNum = roomConfig.PlayerNumAcquired

	uid, _ := uuid.NewV4()
	orderID := fmt.Sprintf("%s", uid)

	buf, err := json.Marshal(orderRecord)
	if err != nil {
		log.Panicln("savePay2Redis, Marshal orderRecord err: ", err)
	}

	conn.Send("MULTI")
	conn.Send("HSET", gconst.UserOrderRecord+userID, orderID, string(buf))

	if result == int32(gconst.SSMsgError_ErrSuccess) {
		conn.Send("HSET", gconst.RoomUnRefund+userID, roomID, orderID)
		conn.Send("HSET", gconst.RoomPayUsers+roomID, "users", userIDListBuf)
		conn.Send("HSET", gconst.AsUserTablePrefix+userID, "diamond", remainDiamond)
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		log.Panicln("savePay2Redis err:", err)
	}

}

func payAndSave2Redis(roomID string, userID string) (remainDiamond int, errCode int32) {
	log.Printf("payAndSave2Redis, roomID:%s, userID:%s", roomID, userID)
	conn := pool.Get()
	defer conn.Close()

	vs, err := redis.Strings(conn.Do("HMGET", gconst.RoomTablePrefix+roomID, "roomConfigID", "roomType", "gameNo"))
	if err != nil {
		log.Println("payAndSave2Redis, get roomConfigID failed, err:", err)
		remainDiamond = 0
		errCode = int32(gconst.SSMsgError_ErrNoRoomConfig)
		return
	}
	var roomConfigID = vs[0]
	var roomType = vs[1]
	var gameNo = vs[2]
	log.Println("payAndSave2Redis, roomType:", roomType)
	if roomConfigID == "" {
		log.Println("payAndSave2Redis, roomConfig not exist")
		remainDiamond = 0
		errCode = int32(gconst.SSMsgError_ErrNoRoomConfig)
		return
	}

	if roomType == "" {
		log.Println("payAndSave2Redis, roomType not exist")
		errCode = int32(gconst.SSMsgError_ErrUnsupportRoomType)
		return
	}

	rooTypeInt, err := strconv.Atoi(roomType)
	if err != nil {
		log.Println("payAndSave2Redis, roomType not exist")
		errCode = int32(gconst.SSMsgError_ErrUnsupportRoomType)
		return
	}

	return payAndSave2RedisWith(rooTypeInt, roomConfigID, roomID, userID, gameNo)

}

func getRetrunDiamond(pay int, handNum int, handFinish int) int {
	log.Printf("getRetrunDiamond, pay:%d, handNum:%d, handFinish:%d", pay, handNum, handFinish)

	var unPlayHand = handNum - handFinish
	var refundDiamond = float64(unPlayHand) / float64(handNum) * float64(pay)
	return int(math.Floor(refundDiamond))
}

func saveRefund(orderRecord *OrderRecord, orderID string) {
	conn := pool.Get()
	defer conn.Close()

	// remove user from room
	var userIDs = loadUsersInRoom(orderRecord.RoomID, conn)
	for i, userID := range userIDs {
		if userID == orderRecord.UserID {
			userIDs = append(userIDs[:i], userIDs[i+1:]...)
			break
		}
	}

	userIDList := &gconst.SSMsgUserIDList{}
	userIDList.UserIDs = userIDs
	userIDListBuf, err := proto.Marshal(userIDList)
	if err != nil {
		log.Println("saveRefund, Marshal userIDList err: ", err)
		return
	}

	buf, err := json.Marshal(orderRecord)
	if err != nil {
		log.Println("saveRefund, Marshal UserRefund err: ", err)
		return
	}
	log.Println("saveRefund:", string(buf))
	conn.Send("MULTI")
	conn.Send("HSET", gconst.UserOrderRecord+orderRecord.UserID, orderID, string(buf))
	conn.Send("HDEL", gconst.RoomUnRefund+orderRecord.UserID, orderRecord.RoomID)

	if orderRecord.Refund.Result == int32(gconst.SSMsgError_ErrSuccess) {
		conn.Send("HSET", gconst.AsUserTablePrefix+orderRecord.UserID, "diamond", orderRecord.Refund.RemainDiamond)
	} else {
		conn.Send("HSET", gconst.RoomRefundFailed+orderRecord.UserID, orderRecord.RoomID, orderID)
	}

	if len(userIDs) > 0 {
		conn.Send("HSET", gconst.RoomPayUsers+orderRecord.RoomID, "users", userIDListBuf)
	} else {
		conn.Send("HDEL", gconst.RoomPayUsers+orderRecord.RoomID, "users")
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		log.Panicln("saveRefund err:", err)
	}

}

func refund2UserAndSave2Redis(roomID string, userID string, handFinish int) *OrderRecord {
	log.Printf("refund2UserAndSave2Redis, roomID:%s, userID:%s, handFinish:%d", roomID, userID, handFinish)
	conn := pool.Get()
	defer conn.Close()

	// 获取未返还的订单，同个用户一个房间只能有一个未返还的订单
	// orderID, err := redis.String(conn.Do("HGET", gconst.RoomUnRefund+userID, roomID))
	conn.Send("MULTI")
	conn.Send("HGET", gconst.RoomUnRefund+userID, roomID)
	conn.Send("HMGET", gconst.RoomTablePrefix+roomID, "roomType", "groupID", "gameNo")
	vs, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		log.Println("refund2UserAndSave2Redis, load UnRefund order err:", err)
		return nil
	}

	var orderID, _ = redis.String(vs[0], nil)
	if orderID == "" {
		log.Println("refund2UserAndSave2Redis, return failed, can't get orderID")
		return nil
	}

	fields, err := redis.Strings(vs[1], nil)
	if err != nil {
		log.Println("refund2UserAndSave2Redis, return failed, can't get roomType ")
		return nil
	}

	var roomType = fields[0]
	roomTypeInt, err := strconv.Atoi(roomType)
	if err != nil {
		roomTypeInt = 0
	}

	var groupID = fields[1]
	// var gameNo = fields[2]

	// var subGameID = 0 // getSubGameIDByRoomType(roomTypeInt)

	log.Printf("refund2UserAndSave2Redis, orderID:%s, roomTypeInt:%d, groupID:%s", orderID, roomTypeInt, groupID)
	// 获取用户的订单
	order, err := redis.String(conn.Do("HGET", gconst.UserOrderRecord+userID, orderID))
	if err != nil {
		log.Println("refund2UserAndSave2Redis, load user order err:", err)
		return nil
	}

	if order == "" {
		log.Println("refund2UserAndSave2Redis, user order not exist")
		return nil
	}

	var orderRecord = &OrderRecord{}
	err = json.Unmarshal([]byte(order), orderRecord)
	if err != nil {
		log.Panicln("refund2UserAndSave2Redis, Unmarshal orderRecord err: ", err)
		return nil
	}

	if orderRecord.Cost == nil {
		log.Println("refund2UserAndSave2Redis, orderRecord.Cost == nil")
		return nil
	}

	if orderRecord.HandNum == 0 {
		log.Println("refund2UserAndSave2Redis, orderRecord.HandNum == 0")
		return nil
	}

	refundDiamond := getRetrunDiamond(orderRecord.Cost.Cost, orderRecord.HandNum, handFinish)
	if refundDiamond < 0 {
		log.Println("refund2UserAndSave2Redis, refundDiamond < 0")
		return nil
	}

	// 哪种类型的返还
	var modDiamondType = ownerModDiamondReturn
	if orderRecord.PayType == aapay {
		modDiamondType = aaModDiamondReturn
		if groupID != "" {
			modDiamondType = aaModDiamondCreateRoomForGroupReturn
		}
	} else if orderRecord.PayType == groupPay {
		modDiamondType = masterModDiamondCreateRoomForGroupReturn
	} else {
		if groupID != "" {
			modDiamondType = ownerModDiamondCreateRoomForGroupReturn
		}
	}

	log.Println("refund2UserAndSave2Redis modDiamondType:", modDiamondType)

	var result int32
	// remainDiamond, err := webdata.ModifyDiamond(userID, modDiamondType, int64(refundDiamond), "解散房间退钱", subGameID, gameNo, groupID)
	// if err != nil {
	// 	result = int32(gconst.SSMsgError_ErrTakeoffDiamondFailedIO)
	// 	var errString = fmt.Sprintf("%v", err)
	// 	if strings.Contains(errString, diamondNotEnoughMsg) {
	// 		result = int32(gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough)
	// 	}

	// } else {
	// 	result = int32(gconst.SSMsgError_ErrSuccess)
	// }

	t := time.Now().UTC()
	var uRefund = &UserRefund{}
	uRefund.Refund = refundDiamond
	uRefund.RemainDiamond = int(0)
	uRefund.Result = result
	uRefund.TimeStamp = t.UnixNano()

	orderRecord.Refund = uRefund
	orderRecord.HandFinish = handFinish

	saveRefund(orderRecord, orderID)

	return orderRecord
}

func isUserExist(userID string, userIDs []string) bool {
	for _, id := range userIDs {
		if userID == id {
			return true
		}
	}

	return false
}

// inGameUserIDs 是房间里面游戏的用户，不是要返还钻石的用户，只用检查作返还的用户是否是游戏里面的用户
func refund2Users(roomID string, handFinish int, inGameUserIDs []string) []*OrderRecord {
	log.Printf("refund, roomID:%s,  handFinish:%d", roomID, handFinish)
	conn := pool.Get()
	defer conn.Close()

	var userIDs = loadUsersInRoom(roomID, conn)
	var orderRecords = make([]*OrderRecord, 0, len(userIDs))
	log.Println("userIDs:", userIDs)
	for _, userID := range userIDs {
		var finish = handFinish
		if len(inGameUserIDs) > 0 && !isUserExist(userID, inGameUserIDs) {
			finish = 0
		}

		orderRecord := refund2UserAndSave2Redis(roomID, userID, finish)
		if orderRecord != nil {
			orderRecords = append(orderRecords, orderRecord)
		}
	}

	return orderRecords
}

// 俱乐部部创建房间扣钱
func clubPayAndSave2Redis(roomType int, roomConfigID string, roomID string, clubID string) (remainDiamond int, payDiamond int, errCode int32) {
	log.Printf("clubPayAndSave2Redis, roomType:%d, roomConfigID:%s, roomID:%s, clubID:%s", roomType, roomConfigID, roomID, clubID)

	// 如果用户已经支付过，则不用再次支付
	var isPay = isUserHavePay(roomID, clubID)
	if isPay {
		log.Printf("clubPayAndSave2Redis, user %s have been pay for room %s", clubID, roomID)
		remainDiamond = 0
		errCode = int32(gconst.SSMsgError_ErrTakeoffDiamondFailedRepeat)
		return
	}

	roomConfigString, ok := roomConfigs[roomConfigID]
	if !ok {
		remainDiamond = 0
		errCode = int32(gconst.SSMsgError_ErrNoRoomConfig)
		return
	}

	roomConfig := parseRoomConfigFromString(roomConfigString)

	var pay = 0
	if roomType != int(gconst.RoomType_GuanDang) {
		needPay, err := getPayDiamondNum(roomConfig.PayType, roomConfig.PlayerNumAcquired, roomConfig.HandNum, roomType)
		if err != nil {
			log.Println("getPayDiamondNum error:", err)
			return
		}
		pay = needPay
	}

	var modDiamondType = ownerModDiamondCreateRoom
	if roomConfig.PayType == aapay {
		modDiamondType = aaModDiamondCreateRoom
	}
	log.Println("payAndSave2RedisWith modDiamondType:", modDiamondType)

	remainDiamond, errCode = modifyClubDiamond(clubID, int(-pay))

	savePay2Redis(roomConfig, roomID, clubID, pay, int(remainDiamond), errCode)

	payDiamond = pay
	return
}

func modifyClubDiamond(clubID string, pay int) (remaindDiamond int, errCode int32) {
	log.Printf("modifyClubDiamond, clubID:%s, pay:%d", clubID, pay)
	conn := pool.Get()
	defer conn.Close()

	diamond, err := redis.Int(conn.Do("HGET", gconst.ClubTablePrefix+clubID, "diamond"))
	if err != nil {
		log.Printf("Can't get Club %s diamon, init as 0", clubID)
		remaindDiamond = 0
		errCode = int32(gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough)
		return
	}

	diamond = diamond + pay
	if diamond < 0 {
		remaindDiamond = diamond
		errCode = int32(gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough)
		return
	}

	_, err = conn.Do("HSET", gconst.ClubTablePrefix+clubID, "diamond", diamond)
	if err != nil {
		log.Println("modifyClubDiamond, set club diamond error", err)
		remaindDiamond = diamond
		errCode = int32(gconst.SSMsgError_ErrTakeoffDiamondFailedIO)
		return
	}

	remaindDiamond = diamond
	return diamond, int32(gconst.SSMsgError_ErrSuccess)
}

func refund2Club(roomID string, handFinish int) *OrderRecord {
	log.Printf("refund2Club, roomID:%s, handFinish:%d", roomID, handFinish)
	conn := pool.Get()
	defer conn.Close()

	var userIDs = loadUsersInRoom(roomID, conn)
	if len(userIDs) != 1 {
		log.Printf("Can't refund to clubRoom %s", roomID)
		return nil
	}

	var clubID = userIDs[0]
	return refund2ClubAndSave2Redis(roomID, clubID, handFinish)
}

func refund2ClubAndSave2Redis(roomID string, clubID string, handFinish int) *OrderRecord {
	log.Printf("refund2ClubAndSave2Redis, roomID:%s, clubID:%s, handFinish:%d", roomID, clubID, handFinish)

	conn := pool.Get()
	defer conn.Close()

	// 获取未返还的订单，同个用户一个房间只能有一个未返还的订单
	orderID, err := redis.String(conn.Do("HGET", gconst.RoomUnRefund+clubID, roomID))
	if err != nil {
		log.Println("refund2ClubAndSave2Redis, load UnRefund order err:", err)
		return nil
	}

	if orderID == "" {
		log.Println("refund2ClubAndSave2Redis, user UnRefund order not exist")
		return nil
	}

	// 获取用户的订单
	order, err := redis.String(conn.Do("HGET", gconst.UserOrderRecord+clubID, orderID))
	if err != nil {
		log.Println("refund2ClubAndSave2Redis, load user order err:", err)
		return nil
	}

	if order == "" {
		log.Println("refund2ClubAndSave2Redis, user order not exist")
		return nil
	}

	var orderRecord = &OrderRecord{}
	err = json.Unmarshal([]byte(order), orderRecord)
	if err != nil {
		log.Panicln("refund2ClubAndSave2Redis, Unmarshal orderRecord err: ", err)
		return nil
	}

	if orderRecord.Cost == nil {
		log.Println("refund2ClubAndSave2Redis, orderRecord.Cost == nil")
		return nil
	}

	if orderRecord.HandNum == 0 {
		log.Println("refund2ClubAndSave2Redis, orderRecord.HandNum == 0")
		return nil
	}

	refundDiamond := getRetrunDiamond(orderRecord.Cost.Cost, orderRecord.HandNum, handFinish)
	if refundDiamond < 0 {
		log.Println("refund2ClubAndSave2Redis, refundDiamond < 0")
		return nil
	}

	remainDiamond, result := modifyClubDiamond(clubID, int(refundDiamond))

	t := time.Now().UTC()
	var uRefund = &UserRefund{}
	uRefund.Refund = refundDiamond
	uRefund.RemainDiamond = int(remainDiamond)
	uRefund.Result = result
	uRefund.TimeStamp = t.UnixNano()

	orderRecord.Refund = uRefund
	orderRecord.HandFinish = handFinish

	saveRefund(orderRecord, orderID)

	return orderRecord
}
