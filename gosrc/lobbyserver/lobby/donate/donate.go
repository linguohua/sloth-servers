package donate

// "webdata"

// func generateDonateUUID() string {
// 	uid, _ := uuid.NewV4()
// 	donateID := fmt.Sprintf("%s", uid)
// 	return donateID
// }

// func getSubGameIDByRoomType(roomType int) int {
// 	var roomTypeStr = fmt.Sprintf("%d", roomType)
// 	gameID, ok := config.SubGameIDs[roomTypeStr]
// 	if !ok {
// 		gameID = 0
// 	}

// 	return gameID
// }

// func getPropByRoomType(roomType int, propsType int) *Prop {
// 	propCfgMap := clientPropCfgsMap[roomType]
// 	if propCfgMap == nil {
// 		return nil
// 	}

// 	prop := propCfgMap[propsType]

// 	return prop

// }

// // TODO: 需要检查用户是否有道具，如果有道具，从道具那里消耗，不扣钻
// func donate(propsType uint32, from string, to string, roomType int) (result *gconst.SSMsgDonateRsp, errCode int32) {
// 	log.Printf("donate, propsType:%d, from:%s, to:%s, roomType:%d", propsType, from, to, roomType)

// 	var prop = getPropByRoomType(roomType, int(propsType))
// 	if prop == nil {
// 		var errMsg = fmt.Sprintf("RoomType:%d not exist propsType:%d", roomType, propsType)
// 		log.Panicln(errMsg)
// 		return
// 	}

// 	var propID = uint32(prop.PropID)
// 	if isUserHaveProp(propID, from) {
// 		rsp, errCode := consumeUserProp(prop, from, to, roomType)
// 		if errCode != int32(gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough) {
// 			return rsp, errCode
// 		}

// 		log.Printf("User %s Prop %d number in redis not same as database", from, propID)
// 	}

// 	var costDiamond = prop.Diamond
// 	var charm = prop.Charm
// 	var cost = int64(costDiamond)
// 	var gameID = getSubGameIDByRoomType(roomType)

// 	remainDiamond, err := webdata.ModifyDiamond(from, modDiamondDonate, -cost, "道具消费", gameID, "", "")
// 	if err != nil {
// 		result = nil
// 		errCode = int32(gconst.SSMsgError_ErrTakeoffDiamondFailedIO)
// 		var errString = fmt.Sprintf("%v", err)
// 		if strings.Contains(errString, diamondNotEnoughMsg) {
// 			errCode = int32(gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough)
// 		}
// 		return
// 	}

// 	conn := pool.Get()
// 	defer conn.Close()

// 	fields, err := redis.Strings(conn.Do("HMGET", gconst.AsUserTablePrefix+to, "charm", "isSyncCharm"))
// 	if err != nil && err != redis.ErrNil {
// 		var errMsg = fmt.Sprintf("load user %s charm failed, err:%v", to, err)
// 		log.Panicln(errMsg)
// 	}

// 	originCharm, err := strconv.Atoi(fields[0])
// 	if err != nil {
// 		originCharm = 0
// 	}

// 	var isSyncCharm = fields[1]

// 	var newCharm = originCharm + charm
// 	var donateID = generateDonateUUID()
// 	var propIDfield = fmt.Sprintf("ExtendCoin:%d", propID)

// 	conn.Send("MULTI")
// 	conn.Send("HSET", gconst.AsUserTablePrefix+from, "diamond", remainDiamond)
// 	conn.Send("HDEL", gconst.AsUserTablePrefix+from, propIDfield)
// 	conn.Send("HSET", gconst.AsUserTablePrefix+to, "charm", newCharm)
// 	conn.Send("HMSET", gconst.DonateTablePrefix+donateID, "from", from, "to", to, "propsType", propsType, "costDiamond", costDiamond, "charm", charm)
// 	conn.Send("LPUSH", gconst.UserDonatePrefix+from, donateID)
// 	conn.Send("LPUSH", gconst.UserDonatePrefix+to, donateID)

// 	_, err = conn.Do("EXEC")
// 	if err != nil {
// 		log.Panicln("save donate err:", err)
// 	}

// 	// 如果是第一次同步DB中的魅力值
// 	if isSyncCharm == "" {
// 		charmInDB := webdata.GetCharmValue(to)
// 		addCharm := newCharm - int(charmInDB)
// 		_, err = webdata.UpdateCharmValue(to, int(propsType), addCharm, "同步魅力值")
// 		if err != nil {
// 			log.Println("UpdateCharmValue error:", err)
// 		} else {
// 			conn.Do("HSET", gconst.AsUserTablePrefix+to, "isSyncCharm", "true")
// 		}

// 	} else {
// 		var remark = fmt.Sprintf("%d游戏牌局内发送道具", gameID)
// 		_, err = webdata.UpdateCharmValue(to, int(propsType), charm, remark)
// 		if err != nil {
// 			log.Println("UpdateCharmValue error:", err)
// 		}
// 	}

// 	var msgDonateRsp = &gconst.SSMsgDonateRsp{}
// 	var int32Diamond = int32(remainDiamond)
// 	msgDonateRsp.Diamond = &int32Diamond
// 	var int32Charm = int32(newCharm)
// 	msgDonateRsp.Charm = &int32Charm

// 	return msgDonateRsp, int32(gconst.SSMsgError_ErrSuccess)
// }

// func isUserHaveProp(propID uint32, userID string) bool {
// 	log.Printf("isUserHaveProp, propID:%d, userID:%s", propID, userID)
// 	conn := pool.Get()
// 	defer conn.Close()

// 	var field = fmt.Sprintf("ExtendCoin:%d", propID)
// 	log.Println("field:", field)
// 	propNum, err := redis.Int(conn.Do("HGET", gconst.AsUserTablePrefix+userID, field))
// 	if err != nil {
// 		log.Println("getUserPropNum error:", err)
// 		return false
// 	}

// 	if propNum > 0 {
// 		return true
// 	}
// 	return false
// }

// func consumeUserProp(prop *Prop, from string, to string, roomType int) (result *gconst.SSMsgDonateRsp, errCode int32) {
// 	log.Printf("consumeUserProp,propID:%d,from:%s, to:%s, roomType:%d", prop.PropID, from, to, roomType)

// 	var charm = prop.Charm
// 	var propID = prop.PropID
// 	var gameID = getSubGameIDByRoomType(roomType)

// 	_, err := webdata.ModifyProps(from, modDiamondDonate, -1, "道具消费", gameID, propID)
// 	if err != nil {
// 		result = nil
// 		errCode = int32(gconst.SSMsgError_ErrTakeoffDiamondFailedIO)
// 		var errString = fmt.Sprintf("%v", err)
// 		if strings.Contains(errString, diamondNotEnoughMsg) {
// 			errCode = int32(gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough)
// 		}
// 		return
// 	}

// 	conn := pool.Get()
// 	defer conn.Close()

// 	conn.Send("MULTI")
// 	conn.Do("HGET", gconst.AsUserTablePrefix+to, "charm")
// 	conn.Do("HGET", gconst.AsUserTablePrefix+from, "diamond")
// 	vs, err := redis.Ints(conn.Do("EXEC"))
// 	if err != nil {
// 		log.Println("consumeUserProp get charm and diamond error:", err)
// 	}

// 	var originCharm = vs[0]
// 	var diamond = vs[1]

// 	var newCharm = originCharm + charm
// 	var donateID = generateDonateUUID()

// 	conn.Send("MULTI")
// 	conn.Send("HSET", gconst.AsUserTablePrefix+to, "charm", newCharm)
// 	conn.Send("HMSET", gconst.DonateTablePrefix+donateID, "from", from, "to", to, "propsType", propID, "costDiamond", 0, "charm", charm)
// 	conn.Send("LPUSH", gconst.UserDonatePrefix+from, donateID)
// 	conn.Send("LPUSH", gconst.UserDonatePrefix+to, donateID)

// 	_, err = conn.Do("EXEC")
// 	if err != nil {
// 		log.Panicln("save donate err:", err)
// 	}

// 	var remark = fmt.Sprintf("%d游戏牌局内发送道具", gameID)
// 	_, err = webdata.UpdateCharmValue(to, propID, charm, remark)
// 	if err != nil {
// 		log.Println("consumeUserProp error:", err)
// 	}

// 	var msgDonateRsp = &gconst.SSMsgDonateRsp{}
// 	var int32Diamond = int32(diamond)
// 	msgDonateRsp.Diamond = &int32Diamond
// 	var int32Charm = int32(newCharm)
// 	msgDonateRsp.Charm = &int32Charm

// 	return msgDonateRsp, int32(gconst.SSMsgError_ErrSuccess)
// }
