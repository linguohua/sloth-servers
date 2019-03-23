package lobby

/*
func isUserInABTest(userID string) bool {
	conn := pool.Get()
	defer conn.Close()

	result, _ := redis.Int(conn.Do("EXISTS", gconst.ABTestUser+userID))
	if result == 1 {
		return true
	}

	return false
}

func getGameServerLowBound(roomType int) (lowBound int, isInControl bool) {
	lowBound = 0
	isInControl = false

	conn := pool.Get()
	defer conn.Close()

	var key = fmt.Sprintf("%s%d", gconst.ABTestController, roomType)
	values, err := redis.Strings(conn.Do("HMGET", key, "isEnable", "lowBoundVersion"))
	if err != nil {
		log.Println("getGameServerLowBound, error:", err)
		return 0, false
	}

	var isEnableStr = values[0]
	var lowBoundVersionStr = values[1]
	isEnable, err := strconv.ParseBool(isEnableStr)
	if err != nil {
		return 0, false
	}

	lowBoundVersion, err := strconv.Atoi(lowBoundVersionStr)
	if err != nil {
		return 0, false
	}

	return lowBoundVersion, isEnable
}

func getGameServerInfosByRoomType(myRoomType int) []*GameServerInfo {
	conn := pool.Get()
	defer conn.Close()

	if myRoomType == 0 {
		myRoomType = int(gconst.RoomType_DafengMJ)
	}

	var setkey = fmt.Sprintf("%s%d", gconst.GameServerInstancePrefix, myRoomType)
	log.Println("setkey:", setkey)
	gameServerIDs, err := redis.Strings(conn.Do("SMEMBERS", setkey))
	if err != nil {
		log.Println("get game server keys from redis err: ", err)
		return nil
	}

	conn.Send("MULTI")
	for _, gameServerID := range gameServerIDs {
		conn.Send("HMGET", gconst.GameServerInstancePrefix+gameServerID, "roomtype", "ver")
	}

	values, err := redis.Values(conn.Do("EXEC"))

	var gameServerInfos = make([]*GameServerInfo, 0, len(values))
	for index, value := range values {
		fileds, err := redis.Ints(value, nil)
		if err != nil {
			continue
		}

		var roomType = fileds[0]
		var ver = fileds[1]
		if roomType != myRoomType {
			continue
		}

		serverID := gameServerIDs[index]

		var gsi = &GameServerInfo{}
		gsi.roomType = roomType
		gsi.serverID = serverID
		gsi.version = ver
		gameServerInfos = append(gameServerInfos, gsi)
	}

	sortGameServer(gameServerInfos)

	return gameServerInfos
}

func getMaxVersionGameServer(myRoomType int) *GameServerInfo {
	var gameServerInfos = getGameServerInfosByRoomType(myRoomType)

	var length = len(gameServerInfos)
	if length > 0 {
		return gameServerInfos[length-1]
	}

	return nil
}

func getGameServerOnBound(myRoomType int, lowBoundVersion int) *GameServerInfo {
	var gameServerInfos = getGameServerInfosByRoomType(myRoomType)

	var length = len(gameServerInfos)
	if length > 0 && gameServerInfos[length-1].version > lowBoundVersion {
		return gameServerInfos[length-1]
	}
	return nil
}

func getGameServerUnderBound(myRoomType int, lowBoundVersion int) *GameServerInfo {
	var gameServerInfos = getGameServerInfosByRoomType(myRoomType)
	var length = len(gameServerInfos)
	for i := 0; i < length; i++ {
		if i > 0 && gameServerInfos[i].version > lowBoundVersion {
			// 返回小于或者等于临界值的服务器信息
			return gameServerInfos[i-1]
		}
	}

	return nil
}

func getMaxVersionGameServerOnBound(roomType int) *GameServerInfo {
	lowBound, ok := getGameServerLowBound(roomType)
	if !ok {
		return nil
	}

	return getGameServerOnBound(roomType, lowBound)
}

func getMaxVersionGameServerUnderBound(roomType int) *GameServerInfo {
	lowBound, ok := getGameServerLowBound(roomType)
	if !ok {
		return getMaxVersionGameServer(roomType)
	}

	return getGameServerUnderBound(roomType, lowBound)
}

// 获取游戏服务器url
// 1.判断用户是否在灰度列表中
// 2.若在灰度列表则获取临界值之上最大版本游戏服务器
// 3. 若不在灰度列表则获取临界值之下最大版本游戏服务器
func getGamerServerInfo(userID string, roomType int) *GameServerInfo {
	if isUserInABTest(userID) {
		return getMaxVersionGameServerOnBound(roomType)
	}
	return getMaxVersionGameServerUnderBound(roomType)

}
*/
