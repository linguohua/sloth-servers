package lobby

import (
	"gconst"
	"runtime/debug"
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

func onNotifyMessage(msgBag *gconst.SSMsgBag) {
	var requestCode = msgBag.GetRequestCode()
	switch requestCode {
	case int32(gconst.SSMsgReqCode_RoomStateNotify):
		onRoomStateNotify(msgBag)
		break
	case int32(gconst.SSMsgReqCode_AAExitRoomNotify):
		onReturnDiamondNotify(msgBag)
		break
	case int32(gconst.SSMsgReqCode_HandBeginNotify):
		onHandBeginNotify(msgBag)
		break
	default:
		log.Println("not handle for request code:", requestCode)
	}
}

// 通知用户房间已满
func notifyUserRoomIsFull(roomNum string, userIDs []string) {
	log.Printf("notifyUserRoomIsFull, roomNum:%s, userIDs:%v", roomNum, userIDs)

	// var msgString = fmt.Sprintf(`{"roomNumber":%s}`, roomNum)
	// var data = []byte(msgString)

	// for _, usreID := range userIDs {
	// 	push(int32(MessageCode_OPNotifyUserRoomIsFull), data, usreID)
	// }
}

func onRoomStateNotify(msgBag *gconst.SSMsgBag) {
	var roomStateNotify = &gconst.SSMsgRoomStateNotify{}
	err := proto.Unmarshal(msgBag.GetParams(), roomStateNotify)
	if err != nil {
		log.Println("onRoomStateNotify, error:", err)
		return
	}

	var roomID = roomStateNotify.GetRoomID()
	var handStartted32 = roomStateNotify.GetHandStartted()
	var userIDs = roomStateNotify.GetUserIDs()
	var roomState = roomStateNotify.GetState()

	log.Printf("onRoomStateNotify, state:%d, roomID:%s, HandStartted:%d, LastActiveTime:%d, UserIDs:%v", roomStateNotify.GetState(),
		roomID, handStartted32, roomStateNotify.GetLastActiveTime(), roomStateNotify.GetUserIDs())

	conn := pool.Get()
	defer conn.Close()

	// 下面的代码通知俱乐部
	strs, _ := redis.Strings(conn.Do("HMGET", gconst.LobbyRoomTablePrefix+roomID, "groupID", "roomType", "ownerID", "roomConfigID", "roomNumber", "clubID"))
	groupID := strs[0]
	roomTypeStr := strs[1]
	ownerID := strs[2]
	roomConfigID := strs[3]
	roomNumber := strs[4]
	clubID := strs[5]

	cfgString, ok := RoomConfigs[roomConfigID]
	if !ok {
		return
	}

	roomConfigJSON := ParseRoomConfigFromString(cfgString)
	rquirePlayerNum := roomConfigJSON.PlayerNumAcquired

	if len(userIDs) == rquirePlayerNum && groupID == "" && roomState == int32(gconst.RoomState_SRoomWaiting) {
		notifyUserRoomIsFull(roomNumber, userIDs)
	}

	log.Println("onRoomStateNotify groupID ID:", groupID)

	// 发通知给牌友群
	if groupID != "" {
		// 如果房主不在房间内，也要加上房主
		var playerNum = len(userIDs)
		for _, userID := range userIDs {
			if userID == ownerID {
				playerNum = playerNum - 1
				break
			}
		}

		playerNum = playerNum + 1

		roomTypeInt, _ := strconv.Atoi(roomTypeStr)
		// TODO：台安的暂时特殊处理
		if roomTypeInt == int(gconst.RoomType_TacnMJ) || roomTypeInt == int(gconst.RoomType_TacnPok) || roomTypeInt == int(gconst.RoomType_DDZ) {
			playerNum = len(userIDs)
		}

		// publishRoomStateChange2Group(groupID, roomID, ClubStateChange, playerNum, userIDs)
	}

	// 发通知给俱乐部
	if clubID != "" {
		// var hasStart = false
		// if handStartted32 > 0 {
		// 	hasStart = true
		// }
		//chost.clubRoomsListener.OnClubRoomStateChanged(clubID, roomID, roomTypeStr, hasStart)
	}
}

func onHandBeginNotify(msgBag *gconst.SSMsgBag) {
	// var handBeginNotify = &gconst.SSMsgHandBeginNotify{}
	// err := proto.Unmarshal(msgBag.GetParams(), handBeginNotify)
	// if err != nil {
	// 	log.Println("onRoomStateNotify,  error:", err)
	// 	return
	// }

	// var roomID = handBeginNotify.GetRoomID()
	// var handStartted32 = handBeginNotify.GetHandStartted()

	// conn := pool.Get()
	// defer conn.Close()

	// roomConfigID, err := redis.String(conn.Do("HGET", gconst.LobbyRoomTablePrefix+roomID, "roomConfigID"))
	// if err != nil {
	// 	log.Println("onHandBeginNotify err:", err)
	// 	return
	// }

	// cfgString, ok := RoomConfigs[roomConfigID]
	// if !ok {
	// 	log.Println("onHandBeginNotify cant' find room config for roomConfigID:", roomConfigID)
	// 	return
	// }

	// roomConfigJSON := ParseRoomConfigFromString(cfgString)
	// if roomConfigJSON.Race == 1 {
	// 	//publishHandBegin2Arena(roomID, int(handStartted32))
	// }
}

func onReturnDiamondNotify(msgBag *gconst.SSMsgBag) {
	log.Println("onReturnDiamondNotify")
	// var msgUpdateBalance = &gconst.SSMsgUpdateBalance{}
	// err := proto.Unmarshal(msgBag.GetParams(), msgUpdateBalance)
	// if err != nil {
	// 	log.Println("onReturnDiamondNotify, err:", err)
	// 	return
	// }

	// var roomID = msgUpdateBalance.GetRoomID()
	// var userID = msgUpdateBalance.GetUserID()
	// log.Printf("onReturnDiamondNotify, roomID:%s, userID:%s", roomID, userID)

	// order := refund2UserAndSave2Redis(roomID, userID, 0)
	// if order != nil && order.Refund != nil {
	// 	updateMoney(uint32(order.Refund.RemainDiamond), userID)
	// }
}

func updateMoney(diamond uint32, userID string) {
	var updateUserMoney = &MsgUpdateUserMoney{}
	var userDiamond = diamond
	updateUserMoney.Diamond = &userDiamond
	SessionMgr().SendProtoMsgTo(userID, updateUserMoney, int32(MessageCode_OPUpdateUserMoney))
}

func onGameServerRequest(msgBag *gconst.SSMsgBag) {
	defer func() {
		if r := recover(); r != nil {
			accSysExceptionCount++
			debug.PrintStack()
			log.Printf("-----Recovered in processRedisPublish:%v\n", r)
		}
	}()

	var requestCode = msgBag.GetRequestCode()
	log.Println("onGameServerRequest, requestCode:", requestCode)
	switch requestCode {
	case int32(gconst.SSMsgReqCode_DeleteRoom):
		onDeleteRoomRequest(msgBag)
		break
	case int32(gconst.SSMsgReqCode_AAEnterRoom):
		onAAEnterRoomRequest(msgBag)
		break
	case int32(gconst.SSMsgReqCode_Donate):
		onDonateRequest(msgBag)
		break
	default:
		log.Println("not handle for request code:", requestCode)
	}
}

func onDeleteRoomRequest(msgBag *gconst.SSMsgBag) {
	log.Println("onDeleteRoomRequest")
	// var gameServer2RoomMgrServerDisbandRoom = &gconst.SSMsgGameServer2RoomMgrServerDisbandRoom{}
	// err := proto.Unmarshal(msgBag.GetParams(), gameServer2RoomMgrServerDisbandRoom)
	// if err != nil {
	// 	log.Println("onDeleteRoomRequest, Unmarshal msg SSMsgGameServer2RoomMgrServerDisbandRoom err:", err)
	// 	replySSMsg(msgBag, gconst.SSMsgError_ErrDecode, nil)
	// 	return
	// }

	// var roomID = gameServer2RoomMgrServerDisbandRoom.GetRoomID()
	// var startHand = gameServer2RoomMgrServerDisbandRoom.GetHandStart()
	// var finishHand = gameServer2RoomMgrServerDisbandRoom.GetHandFinished()
	// var userIDs = gameServer2RoomMgrServerDisbandRoom.GetPlayerUserIDs()

	// log.Printf("onDeleteRoomRequest, roomID:%s, startHand:%d, userIDs:%v", roomID, startHand, userIDs)

	// conn := pool.Get()
	// defer conn.Close()

	// fields, err := redis.Strings(conn.Do("HMGET", gconst.LobbyRoomTablePrefix+roomID, "ownerID", "clubID", "roomConfigID", "groupID", "roomType"))
	// if err == redis.ErrNil {
	// 	log.Printf("onDeleteRoomRequest room %s not exit", roomID)
	// 	replySSMsg(msgBag, gconst.SSMsgError_ErrRoomNotExist, nil)
	// 	return
	// }

	// var onwerID = fields[0]
	// var clubID = fields[1]
	// var roomConfigID = fields[2]
	// var groupID = fields[3]
	// var roomType = fields[4]

	// var roomConfig = GetRoomConfig(roomConfigID)
	// if roomConfig == nil {
	// 	log.Printf("Can't get config,  room:%s,configID:%s", roomID, roomConfigID)
	// 	replySSMsg(msgBag, gconst.SSMsgError_ErrRoomNotExist, nil)
	// 	return
	// }

	// var payType = roomConfig.PayType

	// var orders = make([]*OrderRecord, 0)
	// if clubID != "" && payType == ClubFundPay {
	// 	// 返还钻石给俱乐部, 这是旧的俱乐部扣钻，已经弃用
	// 	var order = refund2Club(roomID, int(startHand))
	// 	if order != nil {
	// 		if order.Refund.Refund != 0 {
	// 			notifyClubFundAddByRoom(order.Refund.Refund, order.Refund.RemainDiamond, "", clubID)
	// 		}

	// 		orders = append(orders, order)
	// 	}
	// } else {
	// 	// 群主支付， 不用管房间里面有多少人
	// 	if groupID != "" && payType == groupPay {
	// 		userIDs = make([]string, 0)
	// 	}

	// 	orders = refund2Users(roomID, int(startHand), userIDs)
	// }

	// if orders == nil || len(orders) == 0 {
	// 	log.Println("refund diamond failed")
	// }

	// deleteRoomInfoFromRedis(roomID, onwerID)

	// if clubID != "" {
	// 	//chost.clubRoomsListener.OnClubRoomDestroy(clubID, roomID)
	// }

	// log.Printf("groupID:%s, payType:%d, startHand:%d", groupID, payType, startHand)
	// if groupID != "" {
	// 	// 通知罗行的俱乐部解散房间
	// 	// publishRoomChangeMessage2Group(groupID, roomID, DeleteClubRoom)

	// 	//统计茶馆的大赢家
	// 	// if payType == groupPay && startHand > 0 {
	// 	// 	go statsGroupBigWiner(groupID, roomType, roomID, gameServer2RoomMgrServerDisbandRoom.PlayerStats)
	// 	// }
	// 	if startHand > 0 {
	// 		// go statsGroupBigWiner(groupID, roomType, roomID, gameServer2RoomMgrServerDisbandRoom.PlayerStats, finishHand)
	// 	}
	// }

	// if roomConfig.Race == 1 {
	// 	var playerStats = gameServer2RoomMgrServerDisbandRoom.GetPlayerStats()
	// 	//publishGameOver2Arena(roomID, int(startHand), playerStats)
	// }

	// //webdata.UpdateUsersExp(finishHand, userIDs)

	// // 回复游戏服务器
	// replySSMsg(msgBag, gconst.SSMsgError_ErrSuccess, nil)
}

// AA制进入房间扣钱请求
func onAAEnterRoomRequest(msgBag *gconst.SSMsgBag) {
	log.Println("onAAEnterRoomRequest")
	// var msgUpdateBalance = &gconst.SSMsgUpdateBalance{}
	// err := proto.Unmarshal(msgBag.GetParams(), msgUpdateBalance)
	// if err != nil {
	// 	log.Println("onAAEnterRoomRequest, Unmarshal msg SSMsgUpdateBalance err:", err)
	// 	replySSMsg(msgBag, gconst.SSMsgError_ErrDecode, nil)
	// 	return
	// }

	// var roomID = msgUpdateBalance.GetRoomID()
	// var userID = msgUpdateBalance.GetUserID()

	// log.Printf("onAAEnterRoomRequest, roomID:%s, userID:%s", roomID, userID)
	// // roomType := int(gconst.RoomType_DafengMJ)
	// diamond, result := payAndSave2Redis(roomID, userID)
	// if result != int32(gconst.SSMsgError_ErrSuccess) {
	// 	var errCode gconst.SSMsgError
	// 	switch result {
	// 	case int32(gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough):
	// 		errCode = gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough
	// 		break
	// 	case int32(gconst.SSMsgError_ErrTakeoffDiamondFailedIO):
	// 		errCode = gconst.SSMsgError_ErrTakeoffDiamondFailedIO
	// 		break
	// 	case int32(gconst.SSMsgError_ErrNoRoomConfig):
	// 		errCode = gconst.SSMsgError_ErrNoRoomConfig
	// 		break
	// 	case int32(gconst.SSMsgError_ErrTakeoffDiamondFailedRepeat):
	// 		// 如果已经扣取钻石，则直接返回成功，让用户再次进入房间
	// 		errCode = gconst.SSMsgError_ErrSuccess
	// 		break
	// 	default:
	// 		log.Panicln("costMoney, unknow errCode:", result)
	// 		break
	// 	}

	// 	replySSMsg(msgBag, errCode, nil)

	// 	log.Printf("onAAEnterRoomRequest, pay failed reply game server, roomID:%s, userID:%s,remaind diamond:%d", roomID, userID, diamond)
	// 	return
	// }

	// replySSMsg(msgBag, gconst.SSMsgError_ErrSuccess, nil)

	// log.Printf("onAAEnterRoomRequest, pay successed reply game server, roomID:%s, userID:%s,remaind diamond:%d", roomID, userID, diamond)
}

func onDonateRequest(msgBag *gconst.SSMsgBag) {
	log.Println("onDonateRequest")
	// TODO: llwant mysql
	// var gameServerID = msgBag.GetSourceURL()
	// var msgDonate = &gconst.SSMsgDonate{}
	// err := proto.Unmarshal(msgBag.GetParams(), msgDonate)
	// if err != nil {
	// 	log.Panicln("Unmarshal SSMsgDonate err:", err)
	// 	return
	// }

	// var from = msgDonate.GetFrom()
	// var to = msgDonate.GetTo()
	// var propsType = msgDonate.GetPropsType()
	// if from == "" {
	// 	log.Panicln("request params from can't be empty")
	// 	return
	// }

	// if to == "" {
	// 	log.Panicln("request params from can't be empty")
	// 	return
	// }

	// if propsType == 0 {
	// 	log.Panicln("request params propsType can't be 0")
	// 	return
	// }

	// if gameServerID == "" {
	// 	log.Panicln("request params gameServerID can't be emtpy")
	// 	return
	// }

	// var roomType = getRoomTypeWithServerID(gameServerID)
	// msgDonateRsp, errCode := donate(uint32(propsType), from, to, roomType)
	// if errCode != int32(gconst.SSMsgError_ErrSuccess) {
	// 	var msgError = gconst.SSMsgError_ErrTakeoffDiamondFailedIO
	// 	if errCode == int32(gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough) {
	// 		msgError = gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough
	// 	}

	// 	replySSMsg(msgBag, msgError, nil)
	// 	return
	// }

	// // 通过房间服务器更新用户钻石
	// user := userMgr.getUserByID(from)
	// if user != nil {
	// 	var diamond = msgDonateRsp.GetDiamond()
	// 	user.updateMoney(uint32(diamond))
	// }

	// msgDonateRspBuf, err := proto.Marshal(msgDonateRsp)
	// if err != nil {
	// 	log.Panicln("Marshal msgDonateRsp err:", err)
	// 	return
	// }

	// // 通过游戏服务器更新用户钻石与魅力
	// replySSMsg(msgBag, gconst.SSMsgError_ErrSuccess, msgDonateRspBuf)

}

// replySSMsg 给其他服务器回复请求完成
func replySSMsg(msgBag *gconst.SSMsgBag, errCode gconst.SSMsgError, params []byte) {
	if msgBag.GetSourceURL() == "" {
		log.Println("replySSMsgError, no source URL")
		return
	}

	replyMsgBag := &gconst.SSMsgBag{}
	var msgType32 = int32(gconst.SSMsgType_Response)
	replyMsgBag.MsgType = &msgType32
	var seqNO32 = msgBag.GetSeqNO()
	replyMsgBag.SeqNO = &seqNO32
	var requestCode32 = msgBag.GetRequestCode()
	replyMsgBag.RequestCode = &requestCode32
	var status32 = int32(errCode)
	replyMsgBag.Status = &status32

	if params != nil {
		replyMsgBag.Params = params
	}

	bytes, err := proto.Marshal(replyMsgBag)
	if err != nil {
		log.Println(err)
		return
	}

	// 获取redis链接，并退出函数时释放
	conn := pool.Get()
	defer conn.Close()

	log.Println("publish message, game server url:", msgBag.GetSourceURL())
	conn.Do("PUBLISH", msgBag.GetSourceURL(), bytes)
}

func updateUserRoomList(userID string) {
	// user := userMgr.getUserByID(userID)
	// if user == nil {
	// 	log.Println("update user roomList failed, user is nil")
	// 	return
	// }

	// var roomInfos = loadRoomInfos(userID)

	// var msgUpdateRoomList = &MsgUpdateRoomList{}
	// msgUpdateRoomList.RoomInfos = roomInfos
	// SessionMgr.SendProtoMsgTo(userID, msgUpdateRoomList, int32(MessageCode_OPUpdateRoomList))
}

func getRoomTypeWithServerID(gameServerID string) int {
	conn := pool.Get()
	defer conn.Close()
	roomType, err := redis.Int(conn.Do("HGET", gconst.GameServerInstancePrefix+gameServerID, "roomtype"))
	if err != nil {
		log.Println("getRoomTypeWithServerID, error:", err)
		return 0
	}
	return roomType
}

// // redisSubscriber 订阅redis频道
// func redisSubscriber() {
// 	for {
// 		conn := pool.Get()

// 		psc := redis.PubSubConn{Conn: conn}
// 		psc.Subscribe(config.ServerID)
// 		keep := true
// 		fmt.Println("begin to wait redis publish msg")
// 		for keep {
// 			switch v := psc.Receive().(type) {
// 			case redis.Message:
// 				// fmt.Printf("sub %s: message: %s\n", v.Channel, v.Data)
// 				// 因为只订阅一个主题，因此忽略change参数
// 				// 同时不可能是
// 				processRedisPublish(v.Data)
// 			case redis.Subscription:
// 				fmt.Printf("sub %s: %s %d\n", v.Channel, v.Kind, v.Count)
// 			case redis.PMessage:
// 				fmt.Printf("sub %s: %s %s\n", v.Channel, v.Pattern, v.Data)
// 			case error:
// 				log.Println("RoomMgr redisSubscriber redis error:", v)
// 				conn.Close()
// 				keep = false
// 				time.Sleep(2 * time.Second)
// 				break
// 			}
// 		}
// 	}
// }

// func processRedisPublish(data []byte) {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			accSysExceptionCount++
// 			debug.PrintStack()
// 			log.Printf("-----Recovered in processRedisPublish:%v\n", r)
// 		}
// 	}()

// 	ssmsgBag := &gconst.SSMsgBag{}
// 	err := proto.Unmarshal(data, ssmsgBag)
// 	if err != nil {
// 		log.Println("processRedisPublish, decode error:", err)
// 		return
// 	}

// 	var msgType = ssmsgBag.GetMsgType()
// 	switch int32(msgType) {
// 	case int32(gconst.SSMsgType_Notify):
// 		onNotifyMessage(ssmsgBag)
// 		break
// 	case int32(gconst.SSMsgType_Request):
// 		go onGameServerRequest(ssmsgBag)
// 		break
// 	case int32(gconst.SSMsgType_Response):
// 		onGameServerRespone(ssmsgBag)
// 		break
// 	default:
// 		log.Printf("No handler for this type %d message", int32(msgType))
// 	}
// }
