package lobby

import (
	log "github.com/sirupsen/logrus"
)

func createRoomForClub(clubID string, roomTypeStr string, roomRuleJSON string) {
	log.Printf("createRoomForClub, clubID:%s, roomTypeStr:%s, roomRuleJSON:%s", clubID, roomTypeStr, roomRuleJSON)
	// TODO: llwant mysql
	// conn := pool.Get()
	// defer conn.Close()
	// userID, err := redis.String(conn.Do("HGET", gconst.ClubTablePrefix+clubID, "owner"))
	// if err != nil {
	// 	log.Println("createRoomForClub error:", err)
	// }

	// log.Println("createRoomForClub configString:", roomRuleJSON)
	// //保存配置
	// roomConfigID, errCode := saveRoomConfigIfNotExist(roomRuleJSON)
	// if errCode != int32(MsgError_ErrSuccess) {
	// 	log.Println("save room config error, errCode:", errCode)
	// 	return
	// }

	// // 分配房间ID
	// uid, _ := uuid.NewV4()
	// roomIDString := fmt.Sprintf("%s", uid)

	// roomTypeInt, err := strconv.Atoi(roomTypeStr)
	// if err != nil {
	// 	log.Println("parse RoomType error:", err)
	// 	return
	// }

	// var roomType = int32(roomTypeInt)

	// var gameServerID = getGameServerID(int(roomType))
	// if gameServerID == "" {
	// 	log.Println("GameServerId not exist, maybe GamerServer not start")
	// 	return
	// }

	// // gameNo为数据库生成房间唯一ID
	// roomNumber, gameNo, err := webdata.GenerateRoomNum(userID)
	// if err != nil {
	// 	log.Println("GenerateRoomNum faile err:", err)
	// 	return
	// }
	// log.Println("createRoomForClub, roomNumber:", roomNumber)

	// roomConfig := parseRoomConfigFromString(roomRuleJSON)
	// if roomConfig == nil {
	// 	log.Println("parse room config error")
	// 	return
	// }

	// diamond, payDiamond, errCode := clubPayAndSave2Redis(roomTypeInt, roomConfigID, roomIDString, clubID)
	// if errCode == int32(gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough) {
	// 	// TODO: 给俱乐部发邮件
	// 	fields, err := redis.Strings(conn.Do("HMGET", gconst.ClubTablePrefix+clubID, "isSendDiamondNotEnough", "clubNumber"))
	// 	if err != nil {
	// 		log.Println("createRoomForClub error:", err)
	// 	}

	// 	var isSendDiamondNotEnough = fields[0]
	// 	var clubNumber = fields[1]

	// 	if isSendDiamondNotEnough != "true" {
	// 		log.Println("Have been send diamond not enough to club")
	// 		var msg = clubNumber + "俱乐部基金不足，自动开房功能已关闭！"
	// 		sendClubMail(msg, userID)
	// 		conn.Do("HSET", gconst.ClubTablePrefix+clubID, "isSendDiamondNotEnough", "true")
	// 		return
	// 	}
	// }

	// if errCode != int32(gconst.SSMsgError_ErrSuccess) {
	// 	log.Println("payAndSave2RedisWith faile err:", err)
	// 	return
	// }

	// var roomNumberString = fmt.Sprintf("%d", roomNumber)

	// msgCreateRoom := &gconst.SSMsgCreateRoom{}
	// msgCreateRoom.RoomID = &roomIDString
	// msgCreateRoom.RoomConfigID = &roomConfigID
	// msgCreateRoom.RoomType = &roomType
	// msgCreateRoom.UserID = &userID
	// msgCreateRoom.RoomNumber = &roomNumberString
	// msgCreateRoom.ClubID = &clubID

	// msgCreateRoomBuf, err := proto.Marshal(msgCreateRoom)
	// if err != nil {
	// 	log.Println("Marshal SSMsgCreateRoom  err： ", err)
	// 	return
	// }

	// msgType := int32(gconst.SSMsgType_Request)
	// requestCode := int32(gconst.SSMsgReqCode_CreateRoom)
	// status := int32(gconst.SSMsgError_ErrSuccess)

	// msgBag := &gconst.SSMsgBag{}
	// msgBag.MsgType = &msgType
	// var sn = generateSn()
	// msgBag.SeqNO = &sn
	// msgBag.RequestCode = &requestCode
	// msgBag.Status = &status
	// var url = config.ServerID
	// msgBag.SourceURL = &url
	// msgBag.Params = msgCreateRoomBuf

	// //等待游戏服务器的回应
	// log.Printf("createRoomForClub, request gameServer create room userID:%s, roomNumber:%s, roomID:%s, gameServerID:%s", userID, roomNumberString, roomIDString, gameServerID)
	// succeed, msgBagReply := sendAndWait(gameServerID, msgBag, time.Second)

	// if succeed {
	// 	errCode := msgBagReply.GetStatus()
	// 	if errCode != 0 {
	// 		log.Println("request game server error:, errCode:", errCode)
	// 		// 创建房间失败，返还钻石
	// 		refund2ClubAndSave2Redis(roomIDString, clubID, 0)

	// 		errCode = converGameServerErrCode2AccServerErrCode(errCode)
	// 		log.Println("Create room for club error, errCode:", errCode)
	// 		return
	// 	}

	// 	// 删除邮件发送标志
	// 	conn.Do("HDEL", gconst.ClubTablePrefix+clubID, "isSendDiamondNotEnough")

	// 	t := time.Now().UTC()
	// 	timeStampInSecond := t.UnixNano() / int64(time.Second)
	// 	saveRoomInfo(msgCreateRoom, gameServerID, roomNumberString, timeStampInSecond, gameNo)

	// 	roomTypeString := fmt.Sprintf("%d", roomType)
	// 	chost.clubRoomsListener.OnClubRoomCreate(clubID, roomIDString, roomTypeString)
	// 	// 使用go routing以防死锁
	// 	go notifyClubFundReduceByRoom(payDiamond, diamond, "", clubID)
	// } else {
	// 	// 创建房间失败，返还钻石
	// 	log.Printf("createRoomForClub, user %s create room failed, request game server timeout", userID)
	// 	refund2ClubAndSave2Redis(roomIDString, clubID, 0)
	// }
}
