package lobby

// func push(ops int32, data []byte, uid string) {
// 	conn := pool.Get()
// 	defer conn.Close()

// 	//加载用户推送服务器ID
// 	serverID, err := redis.String(conn.Do("HGET", gconst.AsUserTablePrefix+uid, "PushServerID"))
// 	if err != nil {
// 		log.Printf("player %s getPushServerID err:%v", uid, err)
// 		return
// 	}

// 	id, _ := uuid.NewV4()
// 	key := fmt.Sprintf("client:%s:%d:%s", uid, ops, id)

// 	conn.Send("MULTI")
// 	conn.Send("SET", key, data)
// 	conn.Send("PUBLISH", serverID, key)
// 	_, err = conn.Do("EXEC")
// 	if err != nil {
// 		log.Printf("player %s push cmd %d err %v", uid, ops, err)
// 		return
// 	}

// 	log.Printf("push player %s push event %s", uid, key)
// }

// // PushAllPlayer 推送数据给所有用户
// func pushAllPlayer(ops int32, data []byte) {
// 	conn := pool.Get()
// 	defer conn.Close()

// 	list, err := redis.Strings(conn.Do("SUNION", gconst.PushServerList))
// 	if err != nil {
// 		log.Println("GetPushServerList err:", err)
// 		return
// 	}

// 	for _, v := range list {
// 		id, _ := uuid.NewV4()
// 		key := fmt.Sprintf("client:%d:%d:%s", 0, ops, id)

// 		conn.Send("MULTI")
// 		conn.Send("SET", key, data)
// 		conn.Send("PUBLISH", v, key)
// 		_, err := conn.Do("EXEC")
// 		if err != nil {
// 			log.Printf("push allplayer server:%s event %s err:%v", v, key, err)
// 			continue
// 		}
// 		log.Printf("push allplayer server:%s push event %s", v, key)
// 	}
// }

// // PushServerEvent 推送给推送服务器事件
// func PushServerEvent(event string, content interface{}) {
// 	conn := pool.Get()
// 	defer conn.Close()

// 	id, _ := uuid.NewV4()
// 	key := fmt.Sprintf("%s:%s", event, id)

// 	conn.Send("MULTI")
// 	conn.Send("SET", key, content)
// 	conn.Send("PUBLISH", event, key)
// 	_, err := conn.Do("EXEC")
// 	if err != nil {
// 		log.Printf("push server event %s err %v", event, err)
// 		return
// 	}

// 	log.Printf("push server event %s", event)
// }

// // PushEvent 推送事件
// func pushEvent(event string, content interface{}, userID string) {
// 	conn := pool.Get()
// 	defer conn.Close()

// 	id, _ := uuid.NewV4()
// 	key := fmt.Sprintf("%s:%s", event, id)

// 	conn.Send("MULTI")
// 	conn.Send("SET", key, content)
// 	conn.Send("PUBLISH", event, key)
// 	_, err := conn.Do("EXEC")
// 	if err != nil {
// 		log.Printf("player %s push event %s err %v", userID, event, err)
// 		return
// 	}
// 	log.Printf("push player %s event %s", userID, event)
// }

// func getPubData(key string) ([]byte, bool) {
// 	conn := pool.Get()
// 	defer conn.Close()

// 	conn.Send("MULTI")
// 	conn.Send("EXISTS", key)
// 	conn.Send("GET", key)
// 	conn.Send("DEL", key)
// 	data, err := redis.Values(conn.Do("EXEC"))
// 	if err != nil {
// 		log.Println("GetPubData err:", err)
// 		return nil, false
// 	}

// 	if data[1] == nil {
// 		return nil, data[0].(int64) == 1
// 	}

// 	return data[1].([]byte), data[0].(int64) == 1
// }
