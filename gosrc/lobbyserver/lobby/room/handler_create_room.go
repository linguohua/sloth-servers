package room

import (
	"crypto/md5"
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
		replayCreateRoom(w, nil, int32(lobby.MsgError_ErrTakeoffDiamondFailedIO), 0)
		break
	case int32(gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough):
		replayCreateRoom(w, nil, int32(lobby.MsgError_ErrTakeoffDiamondFailedNotEnough), 0)
		break
	case int32(gconst.SSMsgError_ErrNoRoomConfig):
		replayCreateRoom(w, nil, int32(lobby.MsgError_ErrTakeoffDiamondFailedNotEnough), 0)
		break
	case int32(gconst.SSMsgError_ErrTakeoffDiamondFailedRepeat):
		replayCreateRoom(w, nil, int32(lobby.MsgError_ErrTakeoffDiamondFailedRepeat), 0)
		break
	default:
		log.Println("unknow pay error:", payError)
	}
	return
}

func replayCreateRoom(w http.ResponseWriter, roomInfo *lobby.RoomInfo, errorCode int32, remainDiamond int32) {
	msgCreateRoomRsp := &lobby.MsgCreateRoomRsp{}
	msgCreateRoomRsp.Result = proto.Int32(errorCode)
	var errString = lobby.ErrorString[errorCode]
	msgCreateRoomRsp.RetMsg = &errString
	msgCreateRoomRsp.RoomInfo = roomInfo
	msgCreateRoomRsp.Diamond = &remainDiamond

	bytes, err := proto.Marshal(msgCreateRoomRsp)
	if err != nil {
		log.Panic("reply msg, marshal msg failed")
		return
	}

	w.Write(bytes)
}

func saveRoomInfo(msgCreateRoom *gconst.SSMsgCreateRoom, gameServerID string, roomNumberString string, timeStampInSecond int64) {
	var userID = msgCreateRoom.GetUserID()
	var roomID = msgCreateRoom.GetRoomID()
	var roomConfigID = msgCreateRoom.GetRoomConfigID()
	var roomType = msgCreateRoom.GetRoomType()

	conn := lobby.Pool().Get()
	defer conn.Close()

	lastActiveTime := timeStampInSecond / 60

	conn.Send("MULTI")
	conn.Send("HSET", gconst.LobbyUserTablePrefix+userID, "roomID", roomID)
	conn.Send("HSET", gconst.LobbyRoomNumberTablePrefix+roomNumberString, "roomID", roomID)
	conn.Send("HMSET", gconst.LobbyRoomTablePrefix+roomID, "ownerID", userID, "roomConfigID",
		roomConfigID, "gameServerID", gameServerID, "roomNumber", roomNumberString, "timeStamp", timeStampInSecond,
		"lastActiveTime", lastActiveTime, "roomType", roomType)
	conn.Send("SADD", gconst.LobbyRoomTableSet, roomID)

	_, err := conn.Do("EXEC")
	if err != nil {
		log.Println("saveRoomInfo err: ", err)
	}

}

func strArray2Comma(ss []string) string {
	result := ""
	for i := 0; i < len(ss)-1; i++ {
		result = result + ss[i] + ","
	}

	result = result + ss[len(ss)-1]

	return result
}

// 1.检查数据库是否已经存在随机数
// 2.若不存在，则保存到数据库，然后返回这个随机数
func validRandNumber(roomNumbers string, roomID string) string {
	conn := lobby.Pool().Get()
	defer conn.Close()
	// luaScript 在startRedis中创建
	randNumber, err := redis.String(lobby.LuaScript.Do(conn, gconst.LobbyRoomNumberTablePrefix, roomID, roomNumbers))
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

	result, _ := redis.Int(conn.Do("HEXISTS", gconst.LobbyRoomConfigTable, roomConfigID))
	if result != 1 {
		_, err := conn.Do("HSET", gconst.LobbyRoomConfigTable, roomConfigID, bytes)
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

func removeUserCreateRoomLock(userID string) {
	conn := lobby.Pool().Get()
	defer conn.Close()

	conn.Do("DEL", gconst.LobbyUserCreatRoomLockPrefix+userID)
}

func isUserCreateRoomLock(userID string, roomID string) bool {
	log.Println("isUserCreateRoomLock, userID:", userID)
	conn := lobby.Pool().Get()
	defer conn.Close()

	// 10秒后自动清除
	var lockTime = 10
	var key = fmt.Sprintf("%s%s", gconst.LobbyUserCreatRoomLockPrefix, userID)
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

func handlerCreateRoom(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := r.Context().Value("userID").(string)
	log.Println("handlerCreateRoom call, userID:", userID)

	// 分配房间ID
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
	var lastRoomInfo = loadLastRoomInfo(userID)

	if lastRoomInfo != nil {
		log.Printf("handlerCreateRoom, User %s in other room, roomNumber: %s, roomId:%s",
			userID, lastRoomInfo.GetRoomNumber(), lastRoomInfo.GetRoomID())
		// reply(w, msgCreateRoomRsp, int32(lobby.MessageCode_OPCreateRoom))
		replayCreateRoom(w, lastRoomInfo, int32(lobby.MsgError_ErrUserInOtherRoom), 0)
		return
	}

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
	diamond, errCode = lobby.PayUtil().DoPayAndSave2RedisWith(int(roomType), roomConfigID, roomIDString, userID)

	// 如果是钻石不足，获取最新的钻石返回给客户端
	if errCode == int32(gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough) {
		log.Println("handlerCreateRoom faile err:", err)
		// TODO: llwant mysql
		var currentDiamond = 0 // webdata.QueryDiamond(userID)
		replayCreateRoom(w, nil, int32(lobby.MsgError_ErrTakeoffDiamondFailedNotEnough), int32(currentDiamond))
		return
	}

	if errCode != int32(gconst.SSMsgError_ErrSuccess) {
		log.Println("payAndSave2RedisWith faile err:", err)
		replyPayError(w, errCode)
		return
	}

	roomTypeInt32 := int32(roomType)
	msgCreateRoom := &gconst.SSMsgCreateRoom{}
	msgCreateRoom.RoomID = &roomIDString
	msgCreateRoom.RoomConfigID = &roomConfigID
	msgCreateRoom.RoomType = &roomTypeInt32
	msgCreateRoom.UserID = &userID
	msgCreateRoom.RoomNumber = &roomNumberString

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
			lobby.PayUtil().Refund2UserAndSave2Redis(roomIDString, userID, 0)

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
		lobby.PayUtil().Refund2UserAndSave2Redis(roomIDString, userID, 0)

		replayCreateRoom(w, nil, int32(lobby.MsgError_ErrRequestGameServerTimeOut), 0)
	}
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
