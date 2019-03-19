package race

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func createRoomForRace(w http.ResponseWriter, msg *MsgCreateRoomReq, userID string, roomIDString string) {
	log.Printf("createRoomForRace call,arenaID:%s,waitTime:%d, userids:%v, onwerID:%s, configString:%s", msg.GetArenaID(), msg.GetWaitTime(), msg.GetUserIDs(), userID, msg.GetConfig())
	// TODO: llwant mysql
	// // 判断是否俱乐部创建房间
	// var arenaID = msg.GetArenaID()
	// if arenaID == "" {
	// 	log.Printf("createRoomForRace, arenaID is empty")
	// 	replyCreateRoomError(w, int32(MsgError_ErrClubIDIsEmtpy))
	// 	return
	// }

	// configString := msg.GetConfig()
	// if configString == "" {
	// 	log.Println("createRoomForRace, room config is not available")
	// 	replyCreateRoomError(w, int32(MsgError_ErrNoRoomConfig))
	// 	return
	// }

	// //保存配置
	// roomConfigID, errCode := saveRoomConfigIfNotExist(configString)
	// if errCode != int32(MsgError_ErrSuccess) {
	// 	log.Println("createRoomForRace, save room config error, errCode:", errCode)
	// 	replyCreateRoomError(w, errCode)
	// 	return
	// }

	// var roomType = msg.GetRoomType()
	// var gameServerID = getGameServerID(int(roomType))
	// if gameServerID == "" {
	// 	log.Println("createRoomForRace, GameServerId not exist, maybe GamerServer not start")
	// 	replyCreateRoomError(w, int32(MsgError_ErrGameServerIDNotExist))
	// 	return
	// }

	// log.Println("createRoomForRace, gameServerID:", gameServerID)

	// roomNumberString := randomRoomNumber(roomIDString)
	// if roomNumberString == "" {
	// 	log.Println("createRoomForRace, randomRoomNumber faile ")
	// 	replyCreateRoomError(w, int32(MsgError_ErrGenerateRoomNumber))
	// 	return
	// }

	// gameNo, err := generateGameNo()
	// if err != nil {
	// 	log.Println("createRoomForRace, generateGameNo error:", err)
	// 	replyCreateRoomError(w, int32(MsgError_ErrGenerateRoomNumber))
	// 	return
	// }

	// log.Printf("createRoomForRace, roomNumber:%s, gameNo:%d", roomNumberString, gameNo)

	// msgCreateRoom := &gconst.SSMsgCreateRoom{}
	// msgCreateRoom.RoomID = &roomIDString
	// msgCreateRoom.RoomConfigID = &roomConfigID
	// msgCreateRoom.RoomType = &roomType
	// msgCreateRoom.UserID = &userID
	// msgCreateRoom.RoomNumber = &roomNumberString
	// msgCreateRoom.ArenaID = &arenaID
	// msgCreateRoom.UserIDs = msg.GetUserIDs()
	// var raceTemplateID = msg.GetRaceTemplateID()
	// msgCreateRoom.RaceTemplateID = &raceTemplateID

	// var waitTime = msg.GetWaitTime()
	// msgCreateRoom.WaitTime = &waitTime

	// msgCreateRoomBuf, err := proto.Marshal(msgCreateRoom)
	// if err != nil {
	// 	log.Println("createRoomForRace, parse roomConfig err： ", err)
	// 	replyCreateRoomError(w, int32(MsgError_ErrEncode))
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
	// log.Printf("createRoomForGroup, request gameServer create room userID:%s, roomNumber:%s, roomID:%s, gameServerID:%s", userID, roomNumberString, roomIDString, gameServerID)
	// succeed, msgBagReply := sendAndWait(gameServerID, msgBag, 10*time.Second)

	// if succeed {
	// 	errCode := msgBagReply.GetStatus()
	// 	if errCode != 0 {
	// 		log.Println("createRoomForGroup, request game server error:, errCode:", errCode)
	// 		// 创建房间失败，返还钻石
	// 		refund2UserAndSave2Redis(roomIDString, userID, 0)

	// 		errCode = converGameServerErrCode2AccServerErrCode(errCode)
	// 		replyCreateRoomError(w, errCode)
	// 		return
	// 	}

	// 	// appendRoom2UserRoomList(msgCreateRoom)

	// 	t := time.Now().UTC()
	// 	timeStampInSecond := t.UnixNano() / int64(time.Second)
	// 	saveRoomInfo(msgCreateRoom, gameServerID, roomNumberString, timeStampInSecond, gameNo)

	// 	roomInfo := &RoomInfo{}
	// 	roomInfo.RoomID = &roomIDString
	// 	roomInfo.RoomNumber = &roomNumberString
	// 	var timeStampString = fmt.Sprintf("%d", timeStampInSecond)
	// 	roomInfo.TimeStamp = &timeStampString
	// 	var lastActiveTime = uint32(timeStampInSecond / 60)
	// 	roomInfo.LastActiveTime = &lastActiveTime

	// 	roomInfo.Config = &configString
	// 	var gameServerURL = getGameServerURL(gameServerID)
	// 	roomInfo.GameServerURL = &gameServerURL

	// 	var propCfg = getPropCfg(int(roomType))
	// 	roomInfo.PropCfg = &propCfg

	// 	var arenaID = msg.GetArenaID()
	// 	roomInfo.ArenaID = &arenaID

	// 	var raceTemplateID = msg.GetRaceTemplateID()
	// 	roomInfo.RaceTemplateID = &raceTemplateID

	// 	var openType = msg.GetOpenType()

	// 	var subType = "0"
	// 	if roomType == int32(gconst.RoomType_TacnMJ) {
	// 		subType = "1"
	// 	}

	// 	roomConfig := parseRoomConfigFromString(configString)
	// 	writeGameStartRecord(int(roomType), roomNumberString, userID, roomConfig.PayType+1, roomConfig.HandNum, configString, webdata.GTMatch, arenaID, subType, gameNo)
	// 	log.Printf("createRoomForRace, user %s create room success", userID)

	// 	replyCreateRoomSuccess(w, roomInfo, openType, 0)

	// } else {
	// 	replyCreateRoomError(w, int32(MsgError_ErrRequestGameServerTimeOut))
	// }
}
