package lobby

// func onMessageUpdateUserInfo(user *User, accessoryMessage *AccessoryMessage) {
// 	log.Println("onMessageUpdateUserInfo")
// 	var buf = accessoryMessage.GetData()
// 	var updateUserInfo = &MsgUpdateUserInfo{}
// 	err := proto.Unmarshal(buf, updateUserInfo)
// 	if err != nil {
// 		log.Println("onMessageUpdateUserInfo, decode error:", err)
// 		return
// 	}

// 	var userIDstring = user.userID()
// 	var location = updateUserInfo.GetLocation()
// 	conn := pool.Get()
// 	defer conn.Close()
// 	conn.Do("HSET", gconst.AsUserTablePrefix+userIDstring, "location", location)
// }
