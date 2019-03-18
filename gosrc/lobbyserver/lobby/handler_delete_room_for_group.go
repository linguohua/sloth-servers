package lobby

import (
	"fmt"
	"gconst"
	"lobbyserver/config"
	"net/http"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	proto "github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

func handlerDeleteRoomForGroup(w http.ResponseWriter, r *http.Request, userID string) {
	groupID := r.URL.Query().Get("groupID")
	if groupID == "" {
		log.Println("handlerDeleteRoomForGroup, groupID is empty")
		var errCode = int32(MsgError_ErrRequestInvalidParam)
		replyDeleteRoom(w, errCode, ErrorString[errCode])
		return
	}

	roomNumber := r.URL.Query().Get("roomID")
	if roomNumber == "" {
		log.Println("handlerDeleteRoomForGroup, roomID is empty")
		var errCode = int32(MsgError_ErrRequestInvalidParam)
		replyDeleteRoom(w, errCode, ErrorString[errCode])
		return
	}

	// TODO: 检查是否是牌友群群主
	masterID := getMasterID(groupID)
	if userID != masterID {
		log.Println("handlerDeleteRoomForGroup, userID != masterID")
		var errCode = int32(MsgError_ErrRequestInvalidParam)
		replyDeleteRoom(w, errCode, ErrorString[errCode])
		return
	}

	log.Printf("handlerDeleteRoom, userID:%s,groupID:%s, roomID:%s", userID, groupID, roomNumber)

	conn := pool.Get()
	defer conn.Close()

	roomID, err := redis.String(conn.Do("HGET", gconst.RoomNumberTable+roomNumber, "roomID"))
	if err != nil && err != redis.ErrNil {
		log.Println("onMessageRequestRoomInfo get roomID err: ", err)
		replyRequestRoomInfo(w, int32(MsgError_ErrDatabase), nil)
		return
	}

	if roomID == "" {
		log.Println("roomNumber not exist")
		replyRequestRoomInfo(w, int32(MsgError_ErrRoomNumberNotExist), nil)
		return
	}

	exist, err := redis.Int(conn.Do("EXISTS", gconst.RoomTablePrefix+roomID))
	if exist == 0 {
		var errCode = int32(MsgError_ErrRoomNotExist)
		replyDeleteRoom(w, errCode, ErrorString[errCode])
		return
	}

	fields, err := redis.Strings(conn.Do("HMGET", gconst.RoomTablePrefix+roomID, "ownerID", "roomType", "roomConfigID"))
	if err == redis.ErrNil {
		var errCode = int32(MsgError_ErrRoomNotExist)
		replyDeleteRoom(w, errCode, ErrorString[errCode])
		return
	}

	if err != nil {
		var errCode = int32(MsgError_ErrDatabase)
		replyDeleteRoom(w, errCode, ErrorString[errCode])
		return
	}

	//检查房间的拥有者
	var ownerID = fields[0]
	// if ownerID != userID {
	// 	log.Printf("onMessageDeleteRoom, %s not room creator,cant delete room, owner is %s", userID, ownerID)
	// 	var errCode = int32(MsgError_ErrNotRoomCreater)
	// 	replyDeleteError(w, errCode, ErrorString[errCode])
	// 	return
	// }
	roomConfigID := fields[2]

	var roomTypeStr = fields[1]
	roomType, err := strconv.Atoi(roomTypeStr)
	if err != nil {
		var errCode = int32(MsgError_ErrDecode)
		replyDeleteRoom(w, errCode, ErrorString[errCode])
		return
	}

	//请求游戏服务器删除房间
	var msgDeleteRoom = &gconst.SSMsgDeleteRoom{}
	msgDeleteRoom.RoomID = &roomID

	msgDeleteRoomBuf, err := proto.Marshal(msgDeleteRoom)
	if err != nil {
		log.Println("parse roomConfig err： ", err)
		var errCode = int32(MsgError_ErrEncode)
		replyDeleteRoom(w, errCode, ErrorString[errCode])
		return
	}

	msgType := int32(gconst.SSMsgType_Request)
	requestCode := int32(gconst.SSMsgReqCode_DeleteRoom)
	status := int32(gconst.SSMsgError_ErrSuccess)

	msgBag := &gconst.SSMsgBag{}
	msgBag.MsgType = &msgType
	var sn = generateSn()
	msgBag.SeqNO = &sn
	msgBag.RequestCode = &requestCode
	msgBag.Status = &status
	var url = config.ServerID
	msgBag.SourceURL = &url
	msgBag.Params = msgDeleteRoomBuf

	// log.Println("roomType:", roomType)
	var gameServerID = getGameServerID(int(roomType))

	succeed, msgBagReply := sendAndWait(gameServerID, msgBag, time.Second)

	if succeed {
		errCode := msgBagReply.GetStatus()
		if errCode != 0 {

			errCode = converGameServerErrCode2AccServerErrCode(errCode)
			replyDeleteRoom(w, errCode, ErrorString[errCode])
			return
		}

		roomConfigJSON := GetRoomConfig(roomConfigID)
		refundAndStatsGroup(msgBagReply, groupID, roomType, roomConfigJSON)

		deleteRoomInfoFromRedis(roomID, ownerID)

		// 通知罗行的俱乐部解散房间
		publishRoomChangeMessage2Group(groupID, roomID, DeleteClubRoom)

		replyDeleteRoom(w, int32(MsgError_ErrSuccess), "ok")
	} else {
		var errCode = int32(MsgError_ErrRequestGameServerTimeOut)
		replyDeleteRoom(w, errCode, ErrorString[errCode])
	}
}

func refundAndStatsGroup(msgBag *gconst.SSMsgBag, groupID string, roomType int, roomCnofigJSON *RoomConfigJSON) {
	log.Println("refundAndStatsGroup")
	var gameServer2RoomMgrServerDisbandRoom = &gconst.SSMsgGameServer2RoomMgrServerDisbandRoom{}
	err := proto.Unmarshal(msgBag.GetParams(), gameServer2RoomMgrServerDisbandRoom)
	if err != nil {
		log.Println("refundAndStatsGroup, Unmarshal msg SSMsgGameServer2RoomMgrServerDisbandRoom err:", err)
		replySSMsg(msgBag, gconst.SSMsgError_ErrDecode, nil)
		return
	}

	var roomID = gameServer2RoomMgrServerDisbandRoom.GetRoomID()
	var startHand = gameServer2RoomMgrServerDisbandRoom.GetHandStart()
	var finishHand = gameServer2RoomMgrServerDisbandRoom.GetHandFinished()
	var userIDs = gameServer2RoomMgrServerDisbandRoom.GetPlayerUserIDs()

	log.Printf("refundAndStatsGroup, roomID:%s, startHand:%d, finishHand:%d, userIDs:%v", roomID, startHand, finishHand, userIDs)

	var payType = roomCnofigJSON.PayType

	var orders = make([]*OrderRecord, 0)

	// 群主支付， 不用管房间里面有多少人
	if groupID != "" && payType == groupPay {
		userIDs = make([]string, 0)
	}

	orders = refund2Users(roomID, int(startHand), userIDs)

	if orders == nil || len(orders) == 0 {
		log.Println("refund diamond failed")
	}

	log.Printf("groupID:%s, payType:%d, startHand:%d", groupID, payType, startHand)

	//统计茶馆的大赢家
	if startHand > 0 {
		roomTypeStr := fmt.Sprintf("%d", roomType)
		go statsGroupBigWiner(groupID, roomTypeStr, roomID, gameServer2RoomMgrServerDisbandRoom.PlayerStats, finishHand)
	}

	// TODO: llwant mysql
	// webdata.UpdateUsersExp(finishHand, userIDs)
}
