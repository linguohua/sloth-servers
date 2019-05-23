package room

import (
	"fmt"
	"gconst"
	"gpubsub"
	"lobbyserver/config"
	"lobbyserver/lobby"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	"io/ioutil"

	"github.com/golang/protobuf/proto"
)

func handlerCreateClubRoom(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("userID")
	isForceUpgrade := r.URL.Query().Get("forceUpgrade")
	clubID := r.URL.Query().Get("clubID")

	log.Printf("handlerCreateRoom call, userID:%s, isForceUpgrade:%s, clubID:%s", userID, isForceUpgrade, clubID)

	updatUtil := lobby.UpdateUtil()
	moduleCfg := updatUtil.GetModuleCfg(r)
	if isForceUpgrade == "true" && moduleCfg != "" {
		replayCreateRoom(w, nil, int32(lobby.MsgError_ErrIsNeedUpdate), 0)
		return
	}

	// 1. 判断牌友圈是否存在
	// 2. 判断用户是否是管理员或者群主
	clubMgr := lobby.ClubMgr()
	club := clubMgr.GetClub(clubID)
	if club == nil {
		log.Printf("handlerCreateClubRoom, no club found for %s", clubID)
		replayCreateRoom(w, nil, int32(lobby.MsgError_ErrRequestInvalidParam), 0)
		return
	}

	if !clubMgr.IsUserPermisionCreateRoom(userID, clubID) {
		log.Printf("handlerCreateClubRoom, user %s not allow create room in club %s", userID, clubID)
		replayCreateRoom(w, nil, int32(lobby.MsgError_ErrOnlyClubCreatorOrManagerAllowCreateRoom), 0)
		return
	}

	uid, _ := uuid.NewV4()
	roomIDString := fmt.Sprintf("%s", uid)

	if isUserCreateRoomLock(userID, roomIDString) {
		log.Println("User crate room is lock !")
		replayCreateRoom(w, nil, int32(lobby.MsgError_ErrUserCreateRoomLock), 0)
		return
	}

	// 退出函数，则清除锁
	defer func() {
		removeUserCreateRoomLock(userID)
	}()

	var gameServerID = r.URL.Query().Get("gsid")

	// 检查是否已经在房间里面
	// var lastRoomInfo = loadLastRoomInfo(userID)

	// if lastRoomInfo != nil {
	// 	log.Printf("handlerCreateRoom, User %s in other room, roomNumber: %s, roomId:%s",
	// 		userID, lastRoomInfo.GetRoomNumber(), lastRoomInfo.GetRoomID())
	// 	// reply(w, msgCreateRoomRsp, int32(lobby.MessageCode_OPCreateRoom))
	// 	replayCreateRoom(w, lastRoomInfo, int32(lobby.MsgError_ErrUserInOtherRoom), 0)
	// 	return
	// }

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("handlerCreateRoom error:", err)
		return
	}

	msg := &lobby.MsgCreateRoomReq{}
	err = proto.Unmarshal(body, msg)
	if err != nil {
		log.Println("onMessageCreateRoom, Unmarshal err:", err)
		replayCreateRoom(w, nil, int32(lobby.MsgError_ErrDecode), 0)
		return
	}

	configString := msg.GetConfig()
	if configString == "" {
		log.Println("room config is not available")
		replayCreateRoom(w, nil, int32(lobby.MsgError_ErrNoRoomConfig), 0)
		return
	}

	log.Println("configString:", configString)
	//保存配置
	roomConfigID, errCode := saveRoomConfigIfNotExist(configString)
	if errCode != int32(lobby.MsgError_ErrSuccess) {
		log.Println("save room config error, errCode:", errCode)
		replayCreateRoom(w, nil, errCode, 0)
		return
	}

	roomConfig := lobby.ParseRoomConfigFromString(configString)
	var roomType = roomConfig.RoomType

	if gameServerID == "" {
		gameServerID = loadLatestGameServer(int(roomType))
	}

	if gameServerID == "" {
		log.Println("GameServerId not exist, maybe GamerServer not start")
		replayCreateRoom(w, nil, int32(lobby.MsgError_ErrGameServerIDNotExist), 0)
		return
	}

	log.Println("handlerCreateRoom, gameServerID:", gameServerID)

	roomNumberString := randomRoomNumber(roomIDString)
	if roomNumberString == "" {
		log.Println("handlerCreateRoom, GenerateRoomNum faile err:", err)
		replayCreateRoom(w, nil, int32(lobby.MsgError_ErrGenerateRoomNumber), 0)
		return
	}

	var diamond = 0
	diamond, errCode = lobby.PayUtil().DoPayForCreateRoom(roomConfigID, roomIDString, userID)
	if errCode != int32(gconst.SSMsgError_ErrSuccess) {
		log.Println("payAndSave2RedisWith faile err:", err)
		replayCreateRoom(w, nil, errCode, int32(diamond))
		return
	}

	roomTypeInt32 := int32(roomType)
	msgCreateRoom := &gconst.SSMsgCreateRoom{}
	msgCreateRoom.RoomID = &roomIDString
	msgCreateRoom.RoomConfigID = &roomConfigID
	msgCreateRoom.RoomType = &roomTypeInt32
	msgCreateRoom.UserID = &userID
	msgCreateRoom.RoomNumber = &roomNumberString
	msgCreateRoom.ClubID = &clubID

	msgCreateRoomBuf, err := proto.Marshal(msgCreateRoom)
	if err != nil {
		log.Println("parse roomConfig err： ", err)
		replayCreateRoom(w, nil, int32(lobby.MsgError_ErrEncode), 0)
		return
	}

	msgType := int32(gconst.SSMsgType_Request)
	requestCode := int32(gconst.SSMsgReqCode_CreateRoom)
	status := int32(gconst.SSMsgError_ErrSuccess)

	msgBag := &gconst.SSMsgBag{}
	msgBag.MsgType = &msgType
	var sn = lobby.GenerateSn()
	msgBag.SeqNO = &sn
	msgBag.RequestCode = &requestCode
	msgBag.Status = &status
	var url = config.ServerID
	msgBag.SourceURL = &url
	msgBag.Params = msgCreateRoomBuf

	//等待游戏服务器的回应
	log.Printf("handlerCreateRoom, request gameServer create room userID:%s, roomNumber:%s, roomID:%s, gameServerID:%s",
		userID, roomNumberString, roomIDString, gameServerID)
	succeed, msgBagReply := gpubsub.SendAndWait(gameServerID, msgBag, 10*time.Second)

	if succeed {
		errCode := msgBagReply.GetStatus()
		if errCode != 0 {
			log.Println("request game server error:, errCode:", errCode)
			// 创建房间失败，返还钻石
			lobby.PayUtil().Refund2UserWith(roomIDString, userID, 0)

			errCode = converGameServerErrCode2AccServerErrCode(errCode)
			replayCreateRoom(w, nil, errCode, 0)
			return
		}

		t := time.Now().UTC()
		timeStampInSecond := t.UnixNano() / int64(time.Second)

		saveRoomInfo(msgCreateRoom, gameServerID, roomNumberString, timeStampInSecond)

		roomInfo := &lobby.RoomInfo{}
		roomInfo.RoomID = &roomIDString
		roomInfo.RoomNumber = &roomNumberString
		var timeStampString = fmt.Sprintf("%d", timeStampInSecond)
		roomInfo.TimeStamp = &timeStampString
		var lastActiveTime = uint32(timeStampInSecond / 60)
		roomInfo.LastActiveTime = &lastActiveTime
		roomInfo.Config = &configString
		roomInfo.GameServerID = &gameServerID
		var propCfg = getPropCfg(int(roomType))
		roomInfo.PropCfg = &propCfg

		//writeGameStartRecord(int(roomType), roomNumberString, userID, roomConfig.PayType+1, roomConfig.HandNum, configString, webdata.GTFriends, "0", subType, gameNo)
		log.Printf("handlerCreateRoom, user %s create room success", userID)
		replayCreateRoom(w, roomInfo, int32(lobby.MsgError_ErrSuccess), int32(diamond))
	} else {
		// 创建房间失败，返还钻石
		log.Printf("handlerCreateRoom, user %s create room failed, request game server timeout", userID)
		lobby.PayUtil().Refund2UserWith(roomIDString, userID, 0)

		replayCreateRoom(w, nil, int32(lobby.MsgError_ErrRequestGameServerTimeOut), 0)
	}
}
