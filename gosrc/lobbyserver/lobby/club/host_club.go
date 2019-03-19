package club

// type clubHost struct {
// 	presentNotifyFunc func(string, bool)

// 	clubRoomsListener club.IClubRoomsHubMgr
// }

// func (ch *clubHost) GetRedisPool() *redis.Pool {
// 	return pool
// }

// func (ch *clubHost) RegisterHTTPHandler(path string, h func(http.ResponseWriter, *http.Request, string)) {
// 	_, exist := accUserIDHTTPHandlers[path]
// 	if exist {
// 		log.Panicf("path %s already register\n", path)
// 		return
// 	}

// 	accUserIDHTTPHandlers[path] = h
// }

// func (ch *clubHost) IsMemberOnline(userID string) bool {
// 	conn := pool.Get()
// 	defer conn.Close()

// 	result, err := redis.Int(conn.Do("SISMEMBER", gconst.LobbyOnlinePlayerList, userID))
// 	if err != nil {
// 		return false
// 	}

// 	if result == 1 {
// 		return true
// 	}

// 	return false
// }

// func (ch *clubHost) RegisterPresentNotify(f func(string, bool)) {
// 	if ch.presentNotifyFunc != nil {
// 		log.Println("presentNotifyFunc not nil, will be replaced")
// 	}

// 	ch.presentNotifyFunc = f
// }

// func (ch *clubHost) RegisterClubRoomListener(listener club.IClubRoomsHubMgr) {
// 	if ch.clubRoomsListener != nil {
// 		log.Println("clubRoomListener not nil, will be replaced")
// 	}

// 	ch.clubRoomsListener = listener
// }

// func (ch *clubHost) GetRoomConfigJSON(configID string) string {
// 	config, ok := roomConfigs[configID]
// 	if ok {
// 		return config
// 	}

// 	return "{}"
// }

// func (ch *clubHost) OnLoadReplayRoomsByIDs(replayRoomIDs []string) []byte {
// 	conn := pool.Get()
// 	defer conn.Close()
// 	return loadReplayRoomsByIDs(replayRoomIDs, conn)
// }

// func (ch *clubHost) CreateRoomForClub(clubID string, roomTypeStr string, roomRuleJSON string) {
// 	createRoomForClub(clubID, roomTypeStr, roomRuleJSON)

// 	// 1. 必要时检查roomTypeStr和roomRuleJSON两个参数看看是否ok
// 	// 2. 检查俱乐部基金是否足够扣除
// 	// 3. 扣钱，扣俱乐部基金
// 	// 4. 创建房间
// 	// 5. 通知俱乐部系统基金变动
// 	// 6. 通知俱乐部系统房间创建：房间状态变化

// }

// func (ch *clubHost) DestroyIdleRoom(roomID string) {
// 	var why = int32(1)
// 	deleteRoomForClub(roomID, true, why)
// 	// TODO: 日光帮我实现一下这个销毁房间的函数
// 	// 其中请求游戏服务器删除房间的函数是以前那个解散房间函数
// 	// 但是增加一个
// 	// 1. 先请求游戏服务器关闭房间
// 	// 2. 然后还钱给俱乐部
// }

// func (ch *clubHost) ForceDestroyClubRoom(roomID string) {
// 	var why = int32(2)
// 	deleteRoomForClub(roomID, false, why)
// 	// TODO: 日光帮我实现一下这个销毁房间的函数
// 	// 其中请求游戏服务器删除房间的函数是以前那个解散房间函数
// 	// 但是增加一个
// 	// 1. 先请求游戏服务器关闭房间
// 	// 2. 然后还钱给俱乐部
// }

// func (ch *clubHost) SendMail(msg string, userID string) {
// 	sendClubMail(msg, userID)
// }

// func hostClub() {
// 	club.HostClub(chost)
// }
