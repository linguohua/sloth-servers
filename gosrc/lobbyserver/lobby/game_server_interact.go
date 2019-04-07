package lobby

import (
	"gconst"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

func onNotifyMessage(msgBag *gconst.SSMsgBag) {
	var requestCode = msgBag.GetRequestCode()
	switch requestCode {
	case int32(gconst.SSMsgReqCode_RoomStateNotify):
		onRoomStateNotify(msgBag)
		break
	case int32(gconst.SSMsgReqCode_AAExitRoomNotify):
		onReturnDiamondNotify(msgBag)
		break
	case int32(gconst.SSMsgReqCode_HandBeginNotify):
		break
	default:
		log.Println("not handle for request code:", requestCode)
	}
}

func onRoomStateNotify(msgBag *gconst.SSMsgBag) {

}

func onReturnDiamondNotify(msgBag *gconst.SSMsgBag) {
	log.Println("onReturnDiamondNotify")
	var msgUpdateBalance = &gconst.SSMsgUpdateBalance{}
	err := proto.Unmarshal(msgBag.GetParams(), msgUpdateBalance)
	if err != nil {
		log.Println("onReturnDiamondNotify, err:", err)
		return
	}

	var roomID = msgUpdateBalance.GetRoomID()
	var userID = msgUpdateBalance.GetUserID()
	log.Printf("onReturnDiamondNotify, roomID:%s, userID:%s", roomID, userID)

	remainDiamond, err := PayUtil().Refund2UserAndSave2Redis(roomID, userID, 0)
	if err == nil {
		updateMoney(uint32(remainDiamond), userID)
	}
}

func updateMoney(diamond uint32, userID string) {
	var updateUserMoney = &MsgUpdateUserMoney{}
	var userDiamond = diamond
	updateUserMoney.Diamond = &userDiamond
	SessionMgr().SendProtoMsgTo(userID, updateUserMoney, int32(MessageCode_OPUpdateUserMoney))
}

func onGameServerRequest(msgBag *gconst.SSMsgBag) {
	var requestCode = msgBag.GetRequestCode()
	log.Println("onGameServerRequest, requestCode:", requestCode)
	switch requestCode {
	case int32(gconst.SSMsgReqCode_DeleteRoom):
		onDeleteRoomRequest(msgBag)
		break
	case int32(gconst.SSMsgReqCode_AAEnterRoom):
		onAAEnterRoomRequest(msgBag)
		break
	case int32(gconst.SSMsgReqCode_Donate):
		onDonateRequest(msgBag)
		break
	default:
		log.Println("not handle for request code:", requestCode)
	}
}

func onDeleteRoomRequest(msgBag *gconst.SSMsgBag) {
	log.Println("onDeleteRoomRequest")
	var gameServer2RoomMgrServerDisbandRoom = &gconst.SSMsgGameServer2RoomMgrServerDisbandRoom{}
	err := proto.Unmarshal(msgBag.GetParams(), gameServer2RoomMgrServerDisbandRoom)
	if err != nil {
		log.Println("onDeleteRoomRequest, Unmarshal msg SSMsgGameServer2RoomMgrServerDisbandRoom err:", err)
		replySSMsg(msgBag, gconst.SSMsgError_ErrDecode, nil)
		return
	}

	var roomID = gameServer2RoomMgrServerDisbandRoom.GetRoomID()
	var startHand = gameServer2RoomMgrServerDisbandRoom.GetHandStart()
	// var finishHand = gameServer2RoomMgrServerDisbandRoom.GetHandFinished()
	var userIDs = gameServer2RoomMgrServerDisbandRoom.GetPlayerUserIDs()

	log.Printf("onDeleteRoomRequest, roomID:%s, startHand:%d, userIDs:%v", roomID, startHand, userIDs)

	conn := pool.Get()
	defer conn.Close()

	fields, err := redis.Strings(conn.Do("HMGET", gconst.LobbyRoomTablePrefix+roomID, "ownerID","roomConfigID"))
	if err == redis.ErrNil {
		log.Printf("onDeleteRoomRequest room %s not exit", roomID)
		replySSMsg(msgBag, gconst.SSMsgError_ErrRoomNotExist, nil)
		return
	}

	var onwerID = fields[0]
	var roomConfigID = fields[1]

	var roomConfig = GetRoomConfig(roomConfigID)
	if roomConfig == nil {
		log.Printf("Can't get config,  room:%s,configID:%s", roomID, roomConfigID)
		replySSMsg(msgBag, gconst.SSMsgError_ErrRoomNotExist, nil)
		return
	}

	var payType = roomConfig.PayType

	if !PayUtil().Refund2Users(roomID, int(startHand), userIDs) {
		log.Println("refund diamond failed")
	}

	RoomUtil().DeleteRoomInfoFromRedis(roomID, onwerID)

	log.Printf("onDeleteRoomRequest payType:%d, startHand:%d", payType, startHand)

	// 回复游戏服务器
	replySSMsg(msgBag, gconst.SSMsgError_ErrSuccess, nil)
}

// AA制进入房间扣钱请求
func onAAEnterRoomRequest(msgBag *gconst.SSMsgBag) {
	log.Println("onAAEnterRoomRequest")
	var msgUpdateBalance = &gconst.SSMsgUpdateBalance{}
	err := proto.Unmarshal(msgBag.GetParams(), msgUpdateBalance)
	if err != nil {
		log.Println("onAAEnterRoomRequest, Unmarshal msg SSMsgUpdateBalance err:", err)
		replySSMsg(msgBag, gconst.SSMsgError_ErrDecode, nil)
		return
	}

	var roomID = msgUpdateBalance.GetRoomID()
	var userID = msgUpdateBalance.GetUserID()

	log.Printf("onAAEnterRoomRequest, roomID:%s, userID:%s", roomID, userID)
	// roomType := int(gconst.RoomType_DafengMJ)
	diamond, result := PayUtil().DoPayAndSave2Redis(roomID, userID)
	if result != int32(gconst.SSMsgError_ErrSuccess) {
		var errCode gconst.SSMsgError
		switch result {
		case int32(gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough):
			errCode = gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough
			break
		case int32(gconst.SSMsgError_ErrTakeoffDiamondFailedIO):
			errCode = gconst.SSMsgError_ErrTakeoffDiamondFailedIO
			break
		case int32(gconst.SSMsgError_ErrNoRoomConfig):
			errCode = gconst.SSMsgError_ErrNoRoomConfig
			break
		case int32(gconst.SSMsgError_ErrTakeoffDiamondFailedRepeat):
			// 如果已经扣取钻石，则直接返回成功，让用户再次进入房间
			errCode = gconst.SSMsgError_ErrSuccess
			break
		default:
			log.Panicln("costMoney, unknow errCode:", result)
			break
		}

		replySSMsg(msgBag, errCode, nil)

		log.Printf("onAAEnterRoomRequest, pay failed reply game server, roomID:%s, userID:%s,remaind diamond:%d", roomID, userID, diamond)
		return
	}

	replySSMsg(msgBag, gconst.SSMsgError_ErrSuccess, nil)

	log.Printf("onAAEnterRoomRequest, pay successed reply game server, roomID:%s, userID:%s,remaind diamond:%d", roomID, userID, diamond)
}

func onDonateRequest(msgBag *gconst.SSMsgBag) {
	log.Println("onDonateRequest")
	// TODO: llwant mysql
	// var gameServerID = msgBag.GetSourceURL()
	// var msgDonate = &gconst.SSMsgDonate{}
	// err := proto.Unmarshal(msgBag.GetParams(), msgDonate)
	// if err != nil {
	// 	log.Panicln("Unmarshal SSMsgDonate err:", err)
	// 	return
	// }

	// var from = msgDonate.GetFrom()
	// var to = msgDonate.GetTo()
	// var propsType = msgDonate.GetPropsType()
	// if from == "" {
	// 	log.Panicln("request params from can't be empty")
	// 	return
	// }

	// if to == "" {
	// 	log.Panicln("request params from can't be empty")
	// 	return
	// }

	// if propsType == 0 {
	// 	log.Panicln("request params propsType can't be 0")
	// 	return
	// }

	// if gameServerID == "" {
	// 	log.Panicln("request params gameServerID can't be emtpy")
	// 	return
	// }

	// var roomType = getRoomTypeWithServerID(gameServerID)
	// msgDonateRsp, errCode := donate(uint32(propsType), from, to, roomType)
	// if errCode != int32(gconst.SSMsgError_ErrSuccess) {
	// 	var msgError = gconst.SSMsgError_ErrTakeoffDiamondFailedIO
	// 	if errCode == int32(gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough) {
	// 		msgError = gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough
	// 	}

	// 	replySSMsg(msgBag, msgError, nil)
	// 	return
	// }

	// // 通过房间服务器更新用户钻石
	// user := userMgr.getUserByID(from)
	// if user != nil {
	// 	var diamond = msgDonateRsp.GetDiamond()
	// 	user.updateMoney(uint32(diamond))
	// }

	// msgDonateRspBuf, err := proto.Marshal(msgDonateRsp)
	// if err != nil {
	// 	log.Panicln("Marshal msgDonateRsp err:", err)
	// 	return
	// }

	// // 通过游戏服务器更新用户钻石与魅力
	// replySSMsg(msgBag, gconst.SSMsgError_ErrSuccess, msgDonateRspBuf)
}

// replySSMsg 给其他服务器回复请求完成
func replySSMsg(msgBag *gconst.SSMsgBag, errCode gconst.SSMsgError, params []byte) {
	if msgBag.GetSourceURL() == "" {
		log.Println("replySSMsgError, no source URL")
		return
	}

	replyMsgBag := &gconst.SSMsgBag{}
	var msgType32 = int32(gconst.SSMsgType_Response)
	replyMsgBag.MsgType = &msgType32
	var seqNO32 = msgBag.GetSeqNO()
	replyMsgBag.SeqNO = &seqNO32
	var requestCode32 = msgBag.GetRequestCode()
	replyMsgBag.RequestCode = &requestCode32
	var status32 = int32(errCode)
	replyMsgBag.Status = &status32

	if params != nil {
		replyMsgBag.Params = params
	}

	bytes, err := proto.Marshal(replyMsgBag)
	if err != nil {
		log.Println(err)
		return
	}

	// 获取redis链接，并退出函数时释放
	conn := pool.Get()
	defer conn.Close()

	log.Println("publish message, game server url:", msgBag.GetSourceURL())
	conn.Do("PUBLISH", msgBag.GetSourceURL(), bytes)
}

func getRoomTypeWithServerID(gameServerID string) int {
	conn := pool.Get()
	defer conn.Close()
	roomType, err := redis.Int(conn.Do("HGET", gconst.GameServerInstancePrefix+gameServerID, "roomtype"))
	if err != nil {
		log.Println("getRoomTypeWithServerID, error:", err)
		return 0
	}
	return roomType
}
