package lobby

import (
	"encoding/json"
	"fmt"
	"gconst"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/garyburd/redigo/redis"
)

var (
	// 因为gameNo和clubID还没保存到redis，暂时做个缓存为支付写流水用，取出后立即删除
	groupRoomInfoMap = make(map[string]*GroupRoomInfo)
)

// GroupRoomInfo 牌友群相关信息，为写流水用
type GroupRoomInfo struct {
	GameNo int64
	ClubID string
}

func isOutOfMaxMemberCreateRoomNum(groupID string, userID string) bool {
	var groupMemberRoomsSetKey = fmt.Sprintf(gconst.GroupMemberRoomsSet, groupID, userID)
	// 拉取配置
	conn := pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("HGET", gconst.Clubserverconfig, "MaxMemberCreateRoomNum")
	conn.Send("SCARD", groupMemberRoomsSetKey)

	values, err := redis.Ints(conn.Do("EXEC"))
	if err != nil {
		log.Println("isOutOfMaxMemberCreateRoomNum, error:", err)
	}

	maxMemberCreateRoomNum := values[0]
	currentClubMemberRoomNum := values[1]

	if currentClubMemberRoomNum >= maxMemberCreateRoomNum {
		return true
	}

	return false
}

// TODO :io挂起会导致不同步问题,解决方法是用lua脚本
func isOutOfMaxGroupCreateRoomNum(groupID string) bool {
	// 拉取配置
	conn := pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("HGET", gconst.Clubserverconfig, "MaxClubCreateRoomNum")
	conn.Send("SCARD", gconst.GroupRoomsSetPrefix+groupID)

	values, err := redis.Ints(conn.Do("EXEC"))
	if err != nil {
		log.Println("isOutOfMaxGroupCreateRoomNum, error:", err)
	}

	maxClubCreateRoomNum := values[0]
	currentClubRoomNum := values[1]

	if currentClubRoomNum >= maxClubCreateRoomNum {
		return true
	}

	return false
}

func isSupportMasterPay(roomType int32, groupID string, userID string) bool {
	conn := pool.Get()
	defer conn.Close()

	// 检查受权开关是否打开，不打开，则直接可群主支付,如果打开，则需要判断用户是否已经授权
	masterPayMode, err := redis.Int(conn.Do("HGET", gconst.Clubserverconfig, "MasterPayMode"))
	if err != nil {
		log.Println("isSupportMasterPay, error:", err)
	}

	if masterPayMode == 0 {
		return true
	}

	log.Println("masterPayMode:", masterPayMode)

	configString, err := redis.String(conn.Do("HGET", gconst.ClubMembersOther+groupID, userID))
	if err != nil {
		log.Println("isSupportMasterPay, err:", err)
		return false
	}

	if configString == "" {
		log.Printf("Group %s get user %s master pay config is empty", groupID, userID)
		return false
	}

	type MasterPayCfg struct {
		UserID    string `json:"id"`
		MasterPay int    `json:"masterpay"`
	}

	masterPayCfg := &MasterPayCfg{}
	err = json.Unmarshal([]byte(configString), &masterPayCfg)
	if err != nil {
		log.Println("isSupportMasterPay, err:", err)
		return false
	}

	if masterPayCfg.MasterPay == 1 {
		return true
	}

	return false
}

func getMasterID(clubID string) string {
	log.Println("getMasterID, clubID:", clubID)
	conn := pool.Get()
	defer conn.Close()

	clubInfoJSON, err := redis.String(conn.Do("HGET", gconst.ClubListKey, clubID))
	if err != nil {
		log.Println("err:", err)
	}

	log.Println("getMasterID, clubInfoJSON:", clubInfoJSON)

	if clubInfoJSON == "" {
		return ""
	}

	type ClubInfo struct {
		Master string `json:"master"`
	}

	clubInfo := &ClubInfo{}
	err = json.Unmarshal([]byte(clubInfoJSON), clubInfo)
	if err != nil {
		log.Println("getMasterID error:", err)
		return ""
	}

	return clubInfo.Master
}

func createRoomForGroup(w http.ResponseWriter, msg *MsgCreateRoomReq, userID string, roomIDString string) {
	log.Println("createRoomForGroup call, userID:", userID)

	// TODO: llwant mysql

	// // 判断是否俱乐部创建房间
	// var groupID = msg.GetClubID()
	// if groupID == "" {
	// 	log.Printf("createRoomForGroup, User %s can't create room in group %s", userID, groupID)
	// 	replyCreateRoomError(w, int32(MsgError_ErrClubIDIsEmtpy))
	// 	return
	// }

	// masterID := getMasterID(groupID)
	// if masterID == "" {
	// 	log.Println("createRoomForGroup masterID is empty")
	// }

	// // 检查是否已经满员了
	// if isOutOfMaxGroupCreateRoomNum(groupID) {
	// 	replyCreateRoomError(w, int32(MsgError_ErrOutOfMaxClubCreateRoomNum))
	// 	return
	// }

	// if masterID != userID && isOutOfMaxMemberCreateRoomNum(groupID, userID) {
	// 	replyCreateRoomError(w, int32(MsgError_ErrOutOfMaxClubMemberCreateRoomNum))
	// 	return
	// }

	// configString := msg.GetConfig()
	// if configString == "" {
	// 	log.Println("createRoomForGroup, room config is not available")
	// 	replyCreateRoomError(w, int32(MsgError_ErrNoRoomConfig))
	// 	return
	// }

	// log.Println("createRoomForGroup, configString:", configString)
	// //保存配置
	// roomConfigID, errCode := saveRoomConfigIfNotExist(configString)
	// if errCode != int32(MsgError_ErrSuccess) {
	// 	log.Println("createRoomForGroup, save room config error, errCode:", errCode)
	// 	replyCreateRoomError(w, errCode)
	// 	return
	// }

	// var roomType = msg.GetRoomType()
	// var gameServerID = getGameServerID(int(roomType))
	// if gameServerID == "" {
	// 	log.Println("createRoomForGroup, GameServerId not exist, maybe GamerServer not start")
	// 	replyCreateRoomError(w, int32(MsgError_ErrGameServerIDNotExist))
	// 	return
	// }

	// log.Println("createRoomForGroup, gameServerID:", gameServerID)

	// roomNumberString := randomRoomNumber(roomIDString)
	// if roomNumberString == "" {
	// 	log.Println("createRoomForGroup, randomRoomNumber faile ")
	// 	replyCreateRoomError(w, int32(MsgError_ErrGenerateRoomNumber))
	// 	return
	// }

	// gameNo, err := generateGameNo()
	// if err != nil {
	// 	log.Println("createRoomForGroup, generateGameNo error:", err)
	// 	replyCreateRoomError(w, int32(MsgError_ErrGenerateRoomNumber))
	// 	return
	// }

	// log.Printf("createRoomForGroup, roomNumber:%s, gameNo:%d", roomNumberString, gameNo)

	// // TODO: 选择牌友群扣费
	// var whoPay = userID
	// var roomConfigJSON = parseRoomConfigFromString(configString)
	// if roomConfigJSON.PayType == groupPay {
	// 	if !isSupportMasterPay(roomType, groupID, userID) {
	// 		log.Println("Not support group master pay")
	// 		replyCreateRoomError(w, int32(MsgError_ErrNotSupportGroupMasterPay))
	// 		return
	// 	}

	// 	whoPay = masterID
	// 	log.Println("Master id:", whoPay)
	// }

	// log.Printf("createRoomForGroup, groupID:%s", groupID)

	// // 为了扣钱流水添加到map中
	// var groupRoomInfo = &GroupRoomInfo{}
	// groupRoomInfo.GameNo = gameNo
	// groupRoomInfo.ClubID = groupID
	// groupRoomInfoMap[roomIDString] = groupRoomInfo

	// gameNoString := fmt.Sprintf("%d", gameNo)
	// _, errCode = payAndSave2RedisWith(int(roomType), roomConfigID, roomIDString, whoPay, gameNoString)

	// // 使用完后立即删除
	// delete(groupRoomInfoMap, roomIDString)
	// // 如果是钻石不足，获取最新的钻石返回给客户端
	// if errCode == int32(gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough) {
	// 	log.Println("createRoomForGroup faile err:", err)
	// 	var currentDiamond, _ = webdata.QueryDiamond(userID)
	// 	if roomConfigJSON.PayType == groupPay {
	// 		replyCreateRoomError(w, int32(MsgError_ErrGroupPayMasterDiamondNotEnough))
	// 		return
	// 	}

	// 	replyCreateRoomErrorAndLastDiamond(w, int32(MsgError_ErrTakeoffDiamondFailedNotEnough), int32(currentDiamond))
	// 	return
	// }

	// if errCode != int32(gconst.SSMsgError_ErrSuccess) {
	// 	log.Println("createRoomForGroup, payAndSave2RedisWith faile err:", err)
	// 	replyPayError(w, errCode)
	// 	return
	// }

	// msgCreateRoom := &gconst.SSMsgCreateRoom{}
	// msgCreateRoom.RoomID = &roomIDString
	// msgCreateRoom.RoomConfigID = &roomConfigID
	// msgCreateRoom.RoomType = &roomType
	// msgCreateRoom.UserID = &userID
	// msgCreateRoom.RoomNumber = &roomNumberString
	// msgCreateRoom.GroupID = &groupID

	// msgCreateRoomBuf, err := proto.Marshal(msgCreateRoom)
	// if err != nil {
	// 	log.Println("createRoomForGroup, parse roomConfig err： ", err)
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

	// 	if roomType != int32(gconst.RoomType_TacnMJ) && roomType != int32(gconst.RoomType_TacnPok) && roomType != int32(gconst.RoomType_DDZ) {
	// 		appendRoom2UserRoomList(msgCreateRoom)
	// 	}

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

	// 	var openType = msg.GetOpenType()

	// 	log.Printf("createRoomForGroup, user %s create room success", userID)

	// 	publishRoomChangeMessage2Group(groupID, roomIDString, CreateClubRoom)

	// 	var subType = "0"
	// 	if roomType == int32(gconst.RoomType_TacnMJ) {
	// 		subType = "1"
	// 	}

	// 	roomConfig := parseRoomConfigFromString(configString)
	// 	writeGameStartRecord(int(roomType), roomNumberString, userID, roomConfig.PayType+1, roomConfig.HandNum, configString, webdata.GTClub, groupID, subType, gameNo)

	// 	var currentDiamond, _ = webdata.QueryDiamond(userID)
	// 	replyCreateRoomSuccess(w, roomInfo, openType, int32(currentDiamond))

	// } else {
	// 	// 创建房间失败，返还钻石
	// 	log.Printf("createRoomForGroup, user %s create room failed, request game server timeout", userID)
	// 	refund2UserAndSave2Redis(roomIDString, userID, 0)

	// 	replyCreateRoomError(w, int32(MsgError_ErrRequestGameServerTimeOut))
	// }
}
