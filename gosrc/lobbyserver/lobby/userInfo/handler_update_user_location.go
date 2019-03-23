package userinfo

import (
	"fmt"
	"gconst"
	"lobbyserver/lobby"
	"net/http"

	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

func handleUpdateUserLocation(w http.ResponseWriter, r *http.Request, userID string) {
	log.Println("handleUpdateUserLocation, userID:", userID)

	accessoryMessage, errCode := lobby.ParseAccessoryMessage(r)
	if errCode != int32(lobby.MsgError_ErrSuccess) {
		var msg = fmt.Sprintf("Update user location error, code:%d", errCode)
		w.WriteHeader(404)
		w.Write([]byte(msg))
		return
	}

	var buf = accessoryMessage.GetData()
	var updateUserInfo = &lobby.MsgUpdateUserInfo{}
	err := proto.Unmarshal(buf, updateUserInfo)
	if err != nil {
		log.Println("handleUpdateUserLocation, decode error:", err)
		var msg = fmt.Sprintf("Decode error:%v", err)
		w.WriteHeader(404)
		w.Write([]byte(msg))
		return
	}

	var location = updateUserInfo.GetLocation()
	conn := lobby.Pool().Get()
	defer conn.Close()
	conn.Do("HSET", gconst.LobbyUserTablePrefix+userID, "location", location)

	sendLocation2GameServer(location, userID)
}

func sendLocation2GameServer(location string, userID string) {
	log.Printf("sendLocation2GameServer, userID:%s", userID)
	// enterRoomID := loadUserLastEnterRoomID(userID)
	// if enterRoomID == "" {
	// 	log.Println("sendLocation2GameServer, enterRoomID is nil")
	// 	return
	// }

	// conn := lobby.Pool().Get()
	// defer conn.Close()

	// serverID, err := redis.String(conn.Do("HGET", gconst.LobbyRoomTablePrefix+enterRoomID, "gameServerID"))
	// if err != nil {
	// 	log.Println("load gameServerID error:", err)
	// 	return
	// }
	// log.Println("sendLocation2GameServer, serverID:", serverID)
	// if serverID == "" {
	// 	log.Println("sendLocation2GameServer, can't get serverID")
	// 	return
	// }

	// updateLocation := &gconst.SSMsgUpdateLocation{}
	// updateLocation.UserID = &userID
	// updateLocation.Location = &location

	// buf, err := proto.Marshal(updateLocation)
	// if err != nil {
	// 	log.Println("Marshal SSMsgUpdateLocation error:", err)
	// 	return
	// }

	// msgType := int32(gconst.SSMsgType_Request)
	// requestCode := int32(gconst.SSMsgReqCode_UpdateLocation)
	// status := int32(gconst.SSMsgError_ErrSuccess)

	// msgBag := &gconst.SSMsgBag{}
	// msgBag.MsgType = &msgType
	// var sn = lobby.GenerateSn()
	// msgBag.SeqNO = &sn
	// msgBag.RequestCode = &requestCode
	// msgBag.Status = &status
	// var url = config.ServerID
	// msgBag.SourceURL = &url
	// msgBag.Params = buf

	// lobby.PublishMsg(serverID, msgBag)
}
