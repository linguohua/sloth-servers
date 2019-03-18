package lobby

import (
	"gconst"
	"lobbyserver/config"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/garyburd/redigo/redis"
	proto "github.com/golang/protobuf/proto"
)

func deleteRoomDataForClub(roomID string) {
	log.Printf("deleteRoomDataForClub, roomID:%s", roomID)

	conn := pool.Get()
	defer conn.Close()

	roomNumberString, err := redis.String(conn.Do("HGET", gconst.RoomTablePrefix+roomID, "roomNumber"))
	if err != nil {
		log.Println("deleteRoomDataForClub, error:", err)
		return
	}

	conn.Send("MULTI")
	conn.Send("DEL", gconst.RoomTablePrefix+roomID)
	conn.Send("DEL", gconst.RoomNumberTable+roomNumberString)
	conn.Send("SREM", gconst.RoomTableACCSet, roomID)
	_, err = conn.Do("EXEC")
	if err != nil {
		log.Println("deleteRoomDataForClub err:", err)
		return
	}
}

func deleteRoomForClub(roomID string, onlyEmpty bool, why int32) {
	log.Printf("deleteRoomForClub, roomID:%s, onlyEmpty:%v, why:%d", roomID, onlyEmpty, why)
	conn := pool.Get()
	defer conn.Close()

	roomType, err := redis.Int(conn.Do("HGET", gconst.RoomTablePrefix+roomID, "roomType"))
	if err != nil {
		log.Println("deleteRoomForClub, error:", err)
		return
	}

	if roomType == 0 {
		log.Println("deleteRoomForClub, roomType == 0")
		return
	}

	var msgDeleteRoom = &gconst.SSMsgDeleteRoom{}
	msgDeleteRoom.RoomID = &roomID
	msgDeleteRoom.OnlyEmpty = &onlyEmpty
	msgDeleteRoom.Why = &why

	msgDeleteRoomBuf, err := proto.Marshal(msgDeleteRoom)
	if err != nil {
		log.Println("deleteRoomForClub Marshal err： ", err)
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

	var gameServerID = getGameServerID(int(roomType))

	succeed, msgBagReply := sendAndWait(gameServerID, msgBag, time.Second)
	if succeed {
		errCode := msgBagReply.GetStatus()
		if errCode != 0 {
			errCode = converGameServerErrCode2AccServerErrCode(errCode)
			log.Println("deleteRoomForClub failed, errCode:", errCode)
			return
		}

		clubID, err := redis.String(conn.Do("HGET", gconst.RoomTablePrefix+roomID, "clubID"))
		if err != nil {
			log.Println("")
		}
		order := refund2ClubAndSave2Redis(roomID, clubID, 0)
		if order.Refund == nil || order.Refund.Result != 0 {
			log.Println("deleteRoomDataForClub error:", order)
			return
		}

		log.Println("deleteRoomDataForClub Refund:", order.Refund.Refund)

		deleteRoomDataForClub(roomID)

		//chost.clubRoomsListener.OnClubRoomDestroy(clubID, roomID)

		// 用go routing 以防死锁
		go notifyClubFundAddByRoom(order.Refund.Refund, order.Refund.RemainDiamond, "", clubID)

	} else {
		log.Println("Request game server time out, errCode:", int32(MsgError_ErrRequestGameServerTimeOut))
	}
}
