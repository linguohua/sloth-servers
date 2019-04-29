package chat

import (
// log "github.com/sirupsen/logrus"
)

// onMessageChat 处理聊天消息
// func onMessageVoiceChat(user *User, message []byte) {
// 	log.Println("onMessageVoiceChat, userID:", user.uID)
// 	userIDList := readUserIDListInRoom(user.userID())

// 	var selfSent = false
// 	for _, uID := range userIDList {
// 		var toUser = userMgr.getUserByID(uID)
// 		if toUser != nil {
// 			toUser.send(message)
// 			if uID == user.uID {
// 				selfSent = true
// 			}
// 		}
// 	}

// 	// 如果碰巧列表中没有自己，那就给自己发一份
// 	if !selfSent {
// 		user.send(message)
// 	}
// }
