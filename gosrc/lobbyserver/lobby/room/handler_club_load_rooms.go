package room

import (
	"gconst"
	"lobbyserver/lobby"
	"net/http"
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func replyLoadClubRooms(w http.ResponseWriter, loadRoomListRsp *lobby.MsgLoadRoomListRsp) {
	bytes, err := proto.Marshal(loadRoomListRsp)
	if err != nil {
		log.Panic("reply msg, marshal msg failed")
		return
	}

	w.Write(bytes)

}

func loadClubRooms(clubID string) []*lobby.RoomInfo {
	// 1 拉取俱乐部的所有房间ID
	// 2 拉取房间数据
	// 3 返回房间

	conn := lobby.Pool().Get()
	defer conn.Close()
	// 现在限制每个牌友群只可创建50个牌友群，多了需要批量拉取
	roomIDs, err := redis.Strings(conn.Do("SMEMBERS", gconst.LobbyClubRoomSetPrefix+clubID))
	if err != nil {
		log.Error("Load club rooms from redis error:", err)
		return []*lobby.RoomInfo{}
	}

	var roomInfos = make([]*lobby.RoomInfo, 0, len(roomIDs))

	conn.Send("MULTI")
	for _, roomID := range roomIDs {
		conn.Send("HMGET", gconst.LobbyRoomTablePrefix+roomID, "roomNumber", "gameServerID", "roomConfigID", "timeStamp", "lastActiveTime")
		conn.Send("HMGET", gconst.GameServerRoomTablePrefix+roomID, "players", "state", "hrStartted")
	}

	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		log.Println("loadRoomInfos err: ", err)
		return roomInfos
	}

	for i := 0; i < len(values); i = i + 2 {
		fileds, err := redis.Strings(values[i], nil)
		if err != nil {
			log.Println("load roomInfo err:", err)
			continue
		}

		var roomID = roomIDs[i/2]
		var roomNunmber = fileds[0]
		var gameServerID = fileds[1]
		var roomConfigID = fileds[2]
		var timeStamp = fileds[3]
		var lastActiveTimeStr = fileds[4]
		if roomNunmber == "" && gameServerID == "" && roomConfigID == "" {
			continue
		}

		var roomInfo = &lobby.RoomInfo{}
		roomInfos = append(roomInfos, roomInfo)

		roomInfo.RoomID = &roomID
		roomInfo.RoomNumber = &roomNunmber
		roomInfo.GameServerID = &gameServerID
		roomInfo.TimeStamp = &timeStamp

		lastActiveTimeInt32, _ := strconv.Atoi(lastActiveTimeStr)
		lastActiveTimeUint32 := uint32(lastActiveTimeInt32)
		roomInfo.LastActiveTime = &lastActiveTimeUint32

		roomConfig, ok := lobby.RoomConfigs[roomConfigID]
		if ok {
			roomInfo.Config = &roomConfig
		}

		vs, err := redis.Values(values[i+1], nil)
		if err != nil {
			log.Println("load user for room err:", err)
			continue
		}

		state, _ := redis.Int(vs[1], nil)
		stateInt32 := int32(state)
		roomInfo.State = &stateInt32

		handStartted, _ := redis.Int(vs[2], nil)
		handStartted32 := int32(handStartted)
		roomInfo.HandStartted = &handStartted32
		log.Println("room's HandStartted:", handStartted32)

		buf, _ := redis.Bytes(vs[0], nil)
		var userIDList = &gconst.SSMsgUserIDList{}
		proto.Unmarshal(buf, userIDList)

		roomInfo.Users = loadUsersProfile(userIDList.UserIDs)

		log.Println("userIDList size:", len(roomInfo.Users))
	}

	return roomInfos
}

func handlerLoadClubRooms(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("userID")
	clubID := r.URL.Query().Get("clubID")

	log.Printf("handlerLoadClubRooms, userID:%s, clubID:%s", userID, clubID)

	var msgLoadRoomListRsp = &lobby.MsgLoadRoomListRsp{}

	clubMgr := lobby.ClubMgr()
	club := clubMgr.GetClub(clubID)
	if club == nil {
		log.Printf("handlerLoadClubRooms, no club found for %s", clubID)
		var result = int32(lobby.MsgError_ErrRequestInvalidParam)
		msgLoadRoomListRsp.Result = &result
		replyLoadClubRooms(w, msgLoadRoomListRsp)

		return
	}

	isClubMember := clubMgr.IsClubMember(userID, clubID)
	if !isClubMember {
		var result = int32(lobby.MsgError_ErrNotClubMember)
		msgLoadRoomListRsp.Result = &result
		replyLoadClubRooms(w, msgLoadRoomListRsp)

		return
	}

	roomInfos := loadClubRooms(clubID)

	msgLoadRoomListRsp.RoomInfos = roomInfos
	var result = int32(lobby.MsgError_ErrSuccess)
	msgLoadRoomListRsp.Result = &result
	replyLoadClubRooms(w, msgLoadRoomListRsp)

}
