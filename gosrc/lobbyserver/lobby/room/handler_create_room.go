package room

import (
	"crypto/md5"
	"fmt"
	"gconst"
	"gpubsub"
	"lobbyserver/config"
	"lobbyserver/lobby"
	"lobbyserver/lobby/pay"
	"net/http"
	"sort"
	"time"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
)

const (
	//MinRoomNum 最小房间号
	MinRoomNum = 100000

	//MaxRoomNum 最大房间号
	MaxRoomNum = 199999

	// ClubFundPay 基金支付
	ClubFundPay = 2
)

var (
	//MaxRoomCount 用户可以创建房间的最多个数
	MaxRoomCount = 20
)

// 注意支付函数返回的错误码都是是stateless里面的
// 需要转成客户端对应的错误
func replyPayError(w http.ResponseWriter, payError int32) {
	switch payError {
	case int32(gconst.SSMsgError_ErrTakeoffDiamondFailedIO):
		replyCreateRoomError(w, int32(lobby.MsgError_ErrTakeoffDiamondFailedIO))
		break
	case int32(gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough):
		replyCreateRoomError(w, int32(lobby.MsgError_ErrTakeoffDiamondFailedNotEnough))
		break
	case int32(gconst.SSMsgError_ErrNoRoomConfig):
		replyCreateRoomError(w, int32(lobby.MsgError_ErrTakeoffDiamondFailedNotEnough))
		break
	case int32(gconst.SSMsgError_ErrTakeoffDiamondFailedRepeat):
		replyCreateRoomError(w, int32(lobby.MsgError_ErrTakeoffDiamondFailedRepeat))
		break
	default:
		log.Println("unknow pay error:", payError)
	}
	return
}

func replyCreateRoomError(w http.ResponseWriter, errorCode int32) {
	msgCreateRoomRsp := &lobby.MsgCreateRoomRsp{}
	msgCreateRoomRsp.Result = proto.Int32(errorCode)
	var errString = lobby.ErrorString[errorCode]
	msgCreateRoomRsp.RetMsg = &errString
	log.Println("errString:", errString)
	reply(w, msgCreateRoomRsp, int32(lobby.MessageCode_OPCreateRoom))
}

func replyCreateRoomErrorAndLastDiamond(w http.ResponseWriter, errorCode int32, diamond int32) {
	msgCreateRoomRsp := &lobby.MsgCreateRoomRsp{}
	msgCreateRoomRsp.Result = proto.Int32(errorCode)
	var errString = lobby.ErrorString[errorCode]
	msgCreateRoomRsp.RetMsg = &errString
	msgCreateRoomRsp.Diamond = proto.Int32(diamond)
	log.Printf("errorCode:%d, errString:%s", errorCode, errString)
	reply(w, msgCreateRoomRsp, int32(lobby.MessageCode_OPCreateRoom))
}

func replyCreateRoomSuccess(w http.ResponseWriter, roomInfo *lobby.RoomInfo, openType int32, diamond int32) {

	msgCreateRoomRsp := &lobby.MsgCreateRoomRsp{}
	var errorCode = int32(lobby.MsgError_ErrSuccess)
	msgCreateRoomRsp.Result = &errorCode
	var errString = lobby.ErrorString[errorCode]
	msgCreateRoomRsp.RetMsg = &errString
	msgCreateRoomRsp.RoomInfo = roomInfo
	msgCreateRoomRsp.OpenType = proto.Int32(openType)

	msgCreateRoomRsp.RoomType = proto.Int32(openType)
	msgCreateRoomRsp.Diamond = proto.Int32(diamond)

	reply(w, msgCreateRoomRsp, int32(lobby.MessageCode_OPCreateRoom))
}

func appendRoom2UserRoomList(msgCreateRoom *gconst.SSMsgCreateRoom) {
	conn := lobby.Pool().Get()
	defer conn.Close()

	var userIDString = msgCreateRoom.GetUserID()

	buf, err := redis.Bytes(conn.Do("HGET", gconst.AsUserTablePrefix+userIDString, "rooms"))
	if err != nil && err != redis.ErrNil {
		log.Println("get user rooms err:", err)
		return
	}

	var roomIDList = &lobby.RoomIDList{}
	if buf != nil {
		err := proto.Unmarshal(buf, roomIDList)
		if err != nil {
			log.Println(err)
			return
		}
	}

	var roomID = msgCreateRoom.GetRoomID()
	roomIDList.RoomIDs = append(roomIDList.RoomIDs, roomID)

	bytes, err := proto.Marshal(roomIDList)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = conn.Do("HSET", gconst.AsUserTablePrefix+userIDString, "rooms", bytes)
	if err != nil {
		log.Println("appendRoom2UserRoomList err:", err)
	}
}

func saveRoomInfo(msgCreateRoom *gconst.SSMsgCreateRoom, gameServerID string, roomNumberString string, timeStampInSecond int64, gameNo int64) {
	var userID = msgCreateRoom.GetUserID()
	var roomID = msgCreateRoom.GetRoomID()
	var roomConfigID = msgCreateRoom.GetRoomConfigID()
	var roomType = msgCreateRoom.GetRoomType()
	var clubID = msgCreateRoom.GetClubID()
	var groupID = msgCreateRoom.GetGroupID()
	var arenaID = msgCreateRoom.GetArenaID()
	var raceTemplateID = msgCreateRoom.GetRaceTemplateID()

	conn := lobby.Pool().Get()
	defer conn.Close()

	lastActiveTime := timeStampInSecond / 60
	// "userID, roomID, configID"
	var record = userID + "," + roomID + "," + roomConfigID + "," + groupID
	conn.Send("MULTI")
	conn.Send("HSET", gconst.RoomGameNo, gameNo, record)
	conn.Send("HSET", gconst.AsUserTablePrefix+userID, "roomID", roomID)
	conn.Send("HSET", gconst.RoomNumberTable+roomNumberString, "roomID", roomID)
	conn.Send("HMSET", gconst.RoomTablePrefix+roomID, "ownerID", userID, "roomConfigID",
		roomConfigID, "gameServerID", gameServerID, "roomNumber", roomNumberString, "timeStamp", timeStampInSecond,
		"lastActiveTime", lastActiveTime, "roomType", roomType, "clubID", clubID, "groupID", groupID,
		"arenaID", arenaID, "raceTemplateID", raceTemplateID, "gameNo", gameNo)
	conn.Send("SADD", gconst.RoomTableACCSet, roomID)

	if groupID != "" {

		// TODO: 记住在解散牌友群的时候把这个对应的数据删除
		conn.Send("SADD", gconst.GroupRoomsSetPrefix+groupID, roomID)

		groupMemberRoomsSetKey := fmt.Sprintf(gconst.GroupMemberRoomsSet, groupID, userID)
		conn.Send("SADD", groupMemberRoomsSetKey, roomID)

	}
	_, err := conn.Do("EXEC")
	if err != nil {
		log.Println("saveRoomInfo err: ", err)
	}

}

func getGameServerURL(gameServerID string) string {
	conn := lobby.Pool().Get()
	defer conn.Close()
	url, _ := redis.String(conn.Do("HGET", gconst.GameServerKeyPrefix+gameServerID, "url"))
	if url == "" {
		url = config.GameServerURL
	}
	return url
}

func strArray2Comma(ss []string) string {
	result := ""
	for i := 0; i < len(ss)-1; i++ {
		result = result + ss[i] + ","
	}

	result = result + ss[len(ss)-1]

	return result
}

// 这个是为了写牌局记录到sql, sql保存牌局需要这个GameNo
func generateGameNo() (int64, error) {
	conn := lobby.Pool().Get()
	defer conn.Close()

	gameNo, err := redis.Int64(conn.Do("INCR", gconst.CurrentRoomGameNo))
	if err != nil {
		log.Println("generateGameNo error:", err)
		return 0, err
	}

	// TODO: llwant mysql

	// 为了续上数据库里面的GameNo
	// if gameNo <= 1 {
	// 	_, gameNo, err = webdata.GenerateRoomNum("10000000")
	// 	if err != nil {
	// 		return 0, err
	// 	}
	// }

	return gameNo, nil
}

// 1.检查数据库是否已经存在随机数
// 2.若不存在，则保存到数据库，然后返回这个随机数
func validRandNumber(roomNumbers string, roomID string) string {
	conn := lobby.Pool().Get()
	defer conn.Close()
	// luaScript 在startRedis中创建
	randNumber, err := redis.String(lobby.LuaScript.Do(conn, gconst.RoomNumberTable, roomID, roomNumbers))
	if err != nil {
		log.Printf("randromNumber error, roomNumbers %s, roomID %s, error:%v ", roomNumbers, roomID, err)
	}

	return randNumber
}

// 生成房间号随机数
func randomRoomNumber(roomID string) string {
	randNumberArray := make([]string, 10)
	for i := 0; i < 10; i++ {
		randNumber := lobby.RandGenerator.Intn(MaxRoomNum-MinRoomNum) + MinRoomNum
		randNumberStr := fmt.Sprintf("%d", randNumber)
		randNumberArray[i] = randNumberStr
	}

	randNumberStrs := strArray2Comma(randNumberArray)
	return validRandNumber(randNumberStrs, roomID)
}

func saveRoomConfigIfNotExist(roomConfig string) (roomConfigID string, errorCode int32) {
	conn := lobby.Pool().Get()
	defer conn.Close()

	bytes := []byte(roomConfig)
	md5Value := md5.Sum(bytes)
	roomConfigID = fmt.Sprintf("%x", md5Value)

	result, _ := redis.Int(conn.Do("HEXISTS", gconst.RoomConfigTable, roomConfigID))
	if result != 1 {
		_, err := conn.Do("HSET", gconst.RoomConfigTable, roomConfigID, bytes)
		if err != nil {
			log.Println("save room config err:", err)
			errorCode = int32(lobby.MsgError_ErrDatabase)
			return
		}
	}

	_, exist := lobby.RoomConfigs[roomConfigID]
	if !exist {
		lobby.RoomConfigs[roomConfigID] = roomConfig
	}

	errorCode = int32(lobby.MsgError_ErrSuccess)
	return
}

// GameServerInfo 保存游戏服务器信息
type GameServerInfo struct {
	serverID string
	version  int
	roomType int
}

// byServerVersion 根据座位ID排序
type byServerVersion []*GameServerInfo

func (s byServerVersion) Len() int {
	return len(s)
}
func (s byServerVersion) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byServerVersion) Less(i, j int) bool {
	return s[i].version > s[j].version
}

// sortPlayers 根据座位ID排序
func sortGameServer(gameServerInfos []*GameServerInfo) {
	sort.Sort(byServerVersion(gameServerInfos))
}

func getGameServerID(myRoomType int) string {
	log.Println("getGameServerID, myRoomType:", myRoomType)
	conn := lobby.Pool().Get()
	defer conn.Close()

	if myRoomType == 0 {
		myRoomType = int(gconst.RoomType_DafengMJ)
	}

	var setkey = fmt.Sprintf("%s%d", gconst.GameServerKeyPrefix, myRoomType)
	log.Println("setkey:", setkey)
	gameServerIDs, err := redis.Strings(conn.Do("SMEMBERS", setkey))
	if err != nil {
		log.Println("get game server keys from redis err: ", err)
		return ""
	}

	log.Println("gameServerIDs:", gameServerIDs)

	conn.Send("MULTI")
	for _, key := range gameServerIDs {
		var gameServerKey = fmt.Sprintf("%s%s", gconst.GameServerKeyPrefix, key)
		log.Println("gameServerKey:", gameServerKey)
		conn.Send("HGET", gameServerKey, "ver")
	}

	values, err := redis.Ints(conn.Do("EXEC"))

	var gameServerInfos = make([]*GameServerInfo, 0, len(values))

	for index, value := range values {
		var ver = value

		serverID := gameServerIDs[index]

		var gsi = &GameServerInfo{}
		gsi.roomType = myRoomType
		gsi.serverID = serverID
		gsi.version = ver
		gameServerInfos = append(gameServerInfos, gsi)
	}

	sortGameServer(gameServerInfos)

	if len(gameServerInfos) > 0 {
		return gameServerInfos[0].serverID
	}
	return ""
}

func checkRoomLimit(userID string) (errCode int32) {
	conn := lobby.Pool().Get()
	defer conn.Close()

	bytes, err := redis.Bytes(conn.Do("HGET", gconst.AsUserTablePrefix+userID, "rooms"))
	if err != nil && err != redis.ErrNil {
		log.Println("get rooms err:", err)
		errCode = int32(lobby.MsgError_ErrDatabase)
		return
	}

	var roomIDList = &lobby.RoomIDList{}
	if bytes != nil {
		err := proto.Unmarshal(bytes, roomIDList)
		if err != nil {
			log.Println(err)
			errCode = int32(lobby.MsgError_ErrDecode)
			return
		}

	}

	if len(roomIDList.GetRoomIDs()) < MaxRoomCount {
		errCode = int32(lobby.MsgError_ErrSuccess)
		return
	}

	errCode = int32(lobby.MsgError_ErrRoomCountIsOutOfLimit)
	return
}

func ifNotExistInGameServerAndCleanRoom(roomInfo *lobby.RoomInfo, owner string) bool {
	if roomInfo == nil {
		log.Println("ifNotExistInGameServerAndCleanRoom, roomInfo == nil")
		return false
	}

	var roomID = roomInfo.GetRoomID()
	if roomID == "" {
		log.Println("ifNotExistInGameServerAndCleanRoom, roomID is empty")
		return false
	}

	conn := lobby.Pool().Get()
	defer conn.Close()

	// startHand, err := redis.Int(conn.Send("HGET",  gconst.GameRoomStatistics+roomID,  "hrStartted"))
	conn.Send("MULTI")
	conn.Send("HGET", gconst.GameRoomStatistics+roomID, "hrStartted")
	conn.Send("EXISTS", gconst.GsRoomTablePrefix+roomID)
	values, err := redis.Ints(conn.Do("EXEC"))
	if err != nil {
		log.Println("ifNotExistInGameServerAndCleanRoom, err:", err)
		return false
	}

	var handStart = values[0]
	var exist = values[1]

	if exist == 1 {
		return false
	}

	var result = lobby.PayUtil().Refund2Users(roomID, int(handStart), []string{owner})
	if result {
		deleteRoomInfoFromRedis(roomID, owner)
		log.Printf("ifNotExistInGameServerAndCleanRoom, room %s, not exist in game server, clean room info", roomID)
		return true
	}

	return false
}

func loadUserRoom(userID string) *lobby.RoomInfo {
	conn := lobby.Pool().Get()
	defer conn.Close()

	bytes, err := redis.Bytes(conn.Do("HGET", gconst.AsUserTablePrefix+userID, "rooms"))
	if err != nil && err != redis.ErrNil {
		return nil
	}

	if bytes == nil {
		return nil
	}

	var roomIDList = &lobby.RoomIDList{}
	err = proto.Unmarshal(bytes, roomIDList)
	if err != nil {
		return nil
	}

	var roomIDs = roomIDList.GetRoomIDs()
	if len(roomIDs) == 0 {
		return nil
	}

	var roomID = roomIDs[0]
	if roomID == "" {
		return nil
	}

	log.Println("loadUserRoom:", roomID)

	values, err := redis.Strings(conn.Do("HMGET", gconst.RoomTablePrefix+roomID, "roomNumber", "roomConfigID", "gameServerID"))
	if err != nil {
		log.Println("load room info err:", err)
		return nil
	}

	var roomNumber = values[0]
	var roomConfigID = values[1]
	var gameServerID = values[2]

	if roomNumber == "" || roomConfigID == "" || gameServerID == "" {
		log.Printf("loadLastRoom, roomID:%s, roomNumber:%s, roomConfigID:%s, gameServerID:%s\n", roomID, roomNumber, roomConfigID, gameServerID)
		return nil
	}

	var gameServerURL = getGameServerURL(gameServerID)
	if gameServerURL == "" {
		log.Printf("loadLastRoom, roomID:%s, gameServerURL is nil\n", roomID)
		return nil
	}

	var roomInfo = &lobby.RoomInfo{}
	roomInfo.RoomID = &roomID
	roomInfo.RoomNumber = &roomNumber
	roomInfo.GameServerURL = &gameServerURL

	//log.Println("loadLastRoom, gserverURL:", gameServerURL)
	roomConfig, ok := lobby.RoomConfigs[roomConfigID]
	if ok {
		roomInfo.Config = &roomConfig
	}

	return roomInfo
}

func removeUserCreateRoomLock(userID string) {
	conn := lobby.Pool().Get()
	defer conn.Close()

	conn.Do("DEL", gconst.UserCreatRoomLock+userID)
}

func isUserCreateRoomLock(userID string, roomID string) bool {
	log.Println("isUserCreateRoomLock, userID:", userID)
	conn := lobby.Pool().Get()
	defer conn.Close()

	// 10秒后自动清除
	var lockTime = 10
	var key = fmt.Sprintf("%s%s", gconst.UserCreatRoomLock, userID)
	result, err := conn.Do("set", key, roomID, "ex", lockTime, "nx")
	if err != nil {
		log.Println("isUserCreateRoomLock, err:", err)
		return false
	}

	if result != nil {
		return false
	}

	return true
}

func handlerCreateRoom(w http.ResponseWriter, r *http.Request, userID string) {
	log.Println("handlerCreateRoom call, userID:", userID)

	// 分配房间ID
	uid, _ := uuid.NewV4()
	roomIDString := fmt.Sprintf("%s", uid)

	if isUserCreateRoomLock(userID, roomIDString) {
		log.Println("User crate room is lock !")
		replyCreateRoomError(w, int32(lobby.MsgError_ErrUserCreateRoomLock))
		return
	}

	// 退出函数，则清除锁
	defer func() {
		removeUserCreateRoomLock(userID)
	}()

	// if isUserInBlacklist(userID) {
	// 	log.Println("handlerCreateRoom,user in black list")
	// 	replyCreateRoomError(w, int32(lobby.MsgError_ErrUserInBlacklist))
	// 	return
	// }

	var gameServerID = r.URL.Query().Get("gsid")

	// 检查是否已经在房间里面
	var lastRoomInfo = loadLastRoomInfo(userID)
	if lastRoomInfo == nil {
		// 如果用户不在最后的房间，拉取用户创建的房间，因为可以开多房引起
		lastRoomInfo = loadUserRoom(userID)
		if ifNotExistInGameServerAndCleanRoom(lastRoomInfo, userID) {
			lastRoomInfo = nil
		}
	}

	if lastRoomInfo != nil {
		msgCreateRoomRsp := &lobby.MsgCreateRoomRsp{}
		var errorCode = int32(lobby.MsgError_ErrUserInOtherRoom)
		msgCreateRoomRsp.Result = &errorCode
		var errString = lobby.ErrorString[errorCode]
		msgCreateRoomRsp.RetMsg = &errString
		msgCreateRoomRsp.RoomInfo = lastRoomInfo
		log.Printf("handlerCreateRoom, User %s in other room, roomNumber: %s, roomId:%s", userID, lastRoomInfo.GetRoomNumber(), lastRoomInfo.GetRoomID())
		reply(w, msgCreateRoomRsp, int32(lobby.MessageCode_OPCreateRoom))
		return
	}

	accessoryMessage, errCode := parseAccessoryMessage(r)
	if errCode != int32(lobby.MsgError_ErrSuccess) {
		replyCreateRoomError(w, errCode)
	}

	bytes := accessoryMessage.GetData()

	msg := &lobby.MsgCreateRoomReq{}
	err := proto.Unmarshal(bytes, msg)
	if err != nil {
		log.Println("onMessageCreateRoom, Unmarshal err:", err)
		replyCreateRoomError(w, int32(lobby.MsgError_ErrDecode))
		return
	}

	// 牌友群创建的房间
	// if msg.GetClubID() != "" {
	// 	createRoomForGroup(w, msg, userID, roomIDString)
	// 	return
	// }

	// 创建比赛房
	// if msg.GetArenaID() != "" {
	// 	createRoomForRace(w, msg, userID, roomIDString)
	// 	return
	// }

	// log.Println("msg:", msg)
	//检查有效性
	//TODO: 检查他的房间数量是否达到他账户类型的限制
	// 例如VIP等级不同，可以开的房间数量也不同
	// 创建房间有不同种类的房间，不仅游戏类型不一样，而且房间的生命周期
	// 扣费等等都可以不一样创建房间有不同种类的房间，不仅游戏类型不一样，而且房间的生命周期
	// 扣费等等都可以不一样
	errCode = checkRoomLimit(userID)
	if errCode != int32(lobby.MsgError_ErrSuccess) {
		log.Println("checkRoomLimit, errCode: ", errCode)
		replyCreateRoomError(w, errCode)
		return
	}

	configString := msg.GetConfig()
	if configString == "" {
		log.Println("room config is not available")
		replyCreateRoomError(w, int32(lobby.MsgError_ErrNoRoomConfig))
		return
	}

	log.Println("configString:", configString)
	//保存配置
	roomConfigID, errCode := saveRoomConfigIfNotExist(configString)
	if errCode != int32(lobby.MsgError_ErrSuccess) {
		log.Println("save room config error, errCode:", errCode)
		replyCreateRoomError(w, errCode)
		return
	}

	var roomType = msg.GetRoomType()

	if gameServerID == "" {
		gameServerID = getGameServerID(int(roomType))
	}

	if gameServerID == "" {
		log.Println("GameServerId not exist, maybe GamerServer not start")
		replyCreateRoomError(w, int32(lobby.MsgError_ErrGameServerIDNotExist))
		return
	}

	log.Println("handlerCreateRoom, gameServerID:", gameServerID)

	// gameNo为数据库生成房间唯一ID
	// roomNumber, gameNo, err := webdata.GenerateRoomNum(userID)
	// if err != nil {
	// 	log.Println("GenerateRoomNum faile err:", err)
	// 	replyCreateRoomError(w, int32(MsgError_ErrRoomNumberNotExist))
	// 	return
	// }
	roomNumberString := randomRoomNumber(roomIDString)
	if roomNumberString == "" {
		log.Println("handlerCreateRoom, GenerateRoomNum faile err:", err)
		replyCreateRoomError(w, int32(lobby.MsgError_ErrGenerateRoomNumber))
		return
	}

	gameNo, err := generateGameNo()
	if err != nil {
		log.Println("handlerCreateRoom, generateGameNo error:", err)
		replyCreateRoomError(w, int32(lobby.MsgError_ErrGenerateRoomNumber))
		return
	}

	log.Printf("handlerCreateRoom, roomNumber:%s, gameNo:%d", roomNumberString, gameNo)

	// var roomNumberString = fmt.Sprintf("%d", roomNumber)

	roomConfig := lobby.ParseRoomConfigFromString(configString)

	gameNoString := fmt.Sprintf("%d", gameNo)
	var diamond = 0
	diamond, errCode = lobby.PayUtil().DoPayAndSave2RedisWith(int(roomType), roomConfigID, roomIDString, userID, gameNoString)

	// 如果是钻石不足，获取最新的钻石返回给客户端
	if errCode == int32(gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough) {
		log.Println("handlerCreateRoom faile err:", err)
		// TODO: llwant mysql
		var currentDiamond = 0 // webdata.QueryDiamond(userID)
		replyCreateRoomErrorAndLastDiamond(w, int32(lobby.MsgError_ErrTakeoffDiamondFailedNotEnough), int32(currentDiamond))
		return
	}

	if errCode != int32(gconst.SSMsgError_ErrSuccess) {
		log.Println("payAndSave2RedisWith faile err:", err)
		replyPayError(w, errCode)
		return
	}

	msgCreateRoom := &gconst.SSMsgCreateRoom{}
	msgCreateRoom.RoomID = &roomIDString
	msgCreateRoom.RoomConfigID = &roomConfigID
	msgCreateRoom.RoomType = &roomType
	msgCreateRoom.UserID = &userID
	msgCreateRoom.RoomNumber = &roomNumberString

	msgCreateRoomBuf, err := proto.Marshal(msgCreateRoom)
	if err != nil {
		log.Println("parse roomConfig err： ", err)
		replyCreateRoomError(w, int32(lobby.MsgError_ErrEncode))
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
			lobby.PayUtil().Refund2UserAndSave2Redis(roomIDString, userID, 0)

			errCode = converGameServerErrCode2AccServerErrCode(errCode)
			replyCreateRoomError(w, errCode)
			return
		}

		// roomConfig := parseRoomConfigFromString(configJSON)
		if roomConfig != nil && roomConfig.PayType != pay.FundPay {
			appendRoom2UserRoomList(msgCreateRoom)
		}

		t := time.Now().UTC()
		timeStampInSecond := t.UnixNano() / int64(time.Second)

		saveRoomInfo(msgCreateRoom, gameServerID, roomNumberString, timeStampInSecond, gameNo)

		roomInfo := &lobby.RoomInfo{}
		roomInfo.RoomID = &roomIDString
		roomInfo.RoomNumber = &roomNumberString
		var timeStampString = fmt.Sprintf("%d", timeStampInSecond)
		roomInfo.TimeStamp = &timeStampString
		var lastActiveTime = uint32(timeStampInSecond / 60)
		roomInfo.LastActiveTime = &lastActiveTime
		roomInfo.Config = &configString
		var gameServerURL = getGameServerURL(gameServerID)
		roomInfo.GameServerURL = &gameServerURL
		var propCfg = getPropCfg(int(roomType))
		roomInfo.PropCfg = &propCfg

		var openType = msg.GetOpenType()

		// var subType = "0"
		// if roomType == int32(gconst.RoomType_TacnMJ) {
		// 	subType = "1"
		// }

		//writeGameStartRecord(int(roomType), roomNumberString, userID, roomConfig.PayType+1, roomConfig.HandNum, configString, webdata.GTFriends, "0", subType, gameNo)
		log.Printf("handlerCreateRoom, user %s create room success", userID)
		replyCreateRoomSuccess(w, roomInfo, openType, int32(diamond))
	} else {
		// 创建房间失败，返还钻石
		log.Printf("handlerCreateRoom, user %s create room failed, request game server timeout", userID)
		lobby.PayUtil().Refund2UserAndSave2Redis(roomIDString, userID, 0)

		replyCreateRoomError(w, int32(lobby.MsgError_ErrRequestGameServerTimeOut))
	}
}

func writeGameStartRecord(roomType int, roomNumber string, ownerID string, roundType int, totalRoound int, gameRoule string, gameType int, mainType string, subType string, gameNo int64) {
	// TODO: llwant mysql
	// gameStartRecord := &webdata.GameStartRecord{}
	// gameStartRecord.GameID = getSubGameIDByRoomType(roomType)
	// gameStartRecord.GameType = gameType // 通过这个接口创建的房间是亲友房
	// gameStartRecord.MainType = mainType // 亲友房不写
	// gameStartRecord.SubType = subType   // 1局模式，2圈模式，3锅模式
	// gameStartRecord.RoomNo = roomNumber
	// gameStartRecord.OwnerPlayerID = ownerID
	// gameStartRecord.RoundType = roundType
	// gameStartRecord.TotalRound = totalRoound
	// gameStartRecord.GameRoule = gameRoule
	// gameStartRecord.WriteTime = time.Now().Local().Format("2006-01-02 15:04:05")
	// gameStartRecord.GameNo = gameNo

	// ok := webdata.WriteGameStartRecord(gameStartRecord)
	// if !ok {
	// 	log.Println("writeGameStartRecord failed!")
	// }
}

func reply(w http.ResponseWriter, pb proto.Message, ops int32) {
	accessoryMessage := &lobby.AccessoryMessage{}
	accessoryMessage.Ops = &ops

	if pb != nil {
		bytes, err := proto.Marshal(pb)

		if err != nil {
			log.Panic("reply msg, marshal msg failed")
			return
		}
		accessoryMessage.Data = bytes
	}

	bytes, err := proto.Marshal(accessoryMessage)
	if err != nil {
		log.Panic("reply msg, marshal msg failed")
		return
	}

	w.Write(bytes)
}

func parseAccessoryMessage(r *http.Request) (accMsg *lobby.AccessoryMessage, errCode int32) {
	if r.ContentLength < 1 {
		log.Println("parseAccessoryMessage failed, content length is zero")
		errCode = int32(lobby.MsgError_ErrRequestInvalidParam)
		return
	}

	message := make([]byte, r.ContentLength)
	n, _ := r.Body.Read(message)
	if n != int(r.ContentLength) {
		log.Println("parseAccessoryMessage failed, can't read request body")
		errCode = int32(lobby.MsgError_ErrRequestInvalidParam)
		return
	}

	accessoryMessage := &lobby.AccessoryMessage{}
	err := proto.Unmarshal(message, accessoryMessage)
	if err != nil {
		log.Println("parseAccessoryMessage failed, Unmarshal msg error:", err)
		errCode = int32(lobby.MsgError_ErrDecode)
		return
	}

	accMsg = accessoryMessage
	errCode = int32(lobby.MsgError_ErrSuccess)
	return
}

func converGameServerErrCode2AccServerErrCode(gameServerErrCode int32) int32 {
	var errCode = gameServerErrCode
	if errCode == int32(gconst.SSMsgError_ErrEncode) {
		errCode = int32(lobby.MsgError_ErrEncode)
	} else if errCode == int32(gconst.SSMsgError_ErrDecode) {
		errCode = int32(lobby.MsgError_ErrDecode)
	} else if errCode == int32(gconst.SSMsgError_ErrRoomExist) {
		errCode = int32(lobby.MsgError_ErrGameServerRoomExist)
	} else if errCode == int32(gconst.SSMsgError_ErrNoRoomConfig) {
		errCode = int32(lobby.MsgError_ErrGameServerNoRoomConfig)
	} else if errCode == int32(gconst.SSMsgError_ErrUnsupportRoomType) {
		errCode = int32(lobby.MsgError_ErrGameServerUnsupportRoomType)
	} else if errCode == int32(gconst.SSMsgError_ErrDecodeRoomConfig) {
		errCode = int32(lobby.MsgError_ErrGameServerDecodeRoomConfig)
	} else if errCode == int32(gconst.SSMsgError_ErrRoomNotExist) {
		errCode = int32(lobby.MsgError_ErrGameServerRoomNotExist)
	}

	return errCode

}
