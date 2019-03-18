package lobby

import (
	log "github.com/sirupsen/logrus"
	"gconst"
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
)

func replyLoadRooms(user *User, msgLoadRoomListRsp *MsgLoadRoomListRsp) {
	user.sendMsg(msgLoadRoomListRsp, int32(MessageCode_OPLoadRooms))
}

func loadUsersProfile(userIDs []string) []*UserProfile {
	userProfiles := make([]*UserProfile, 0, len(userIDs))

	conn := pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	for _, userID := range userIDs {
		conn.Send("HMGET", gconst.AsUserTablePrefix+userID, "userName", "userNick")
	}

	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		log.Println("loadUsersProfile err: ", err)
		return userProfiles
	}

	for index, v := range values {
		fileds, err := redis.Strings(v, nil)
		if err != nil {
			continue
		}
		userName := fileds[0]
		nickName := fileds[1]
		userID := userIDs[index]
		userProfile := &UserProfile{}
		userProfile.UserID = &userID
		userProfile.UserName = &userName
		userProfile.NickName = &nickName

		userProfiles = append(userProfiles, userProfile)
	}

	return userProfiles
}

func loadRoomInfos(userIDString string) []*RoomInfo {
	conn := pool.Get()
	defer conn.Close()

	bytes, err := redis.Bytes(conn.Do("HGET", gconst.AsUserTablePrefix+userIDString, "rooms"))
	if err != nil {
		log.Println("loadRoomInfos, err:", err)
		return make([]*RoomInfo, 0)
	}

	var roomIDList = &RoomIDList{}
	if bytes != nil {
		err := proto.Unmarshal(bytes, roomIDList)
		if err != nil {
			log.Println("loadRoomInfos, err:", err)
			return make([]*RoomInfo, 0)
		}
	}

	// var msgLoadRoomList = &MsgLoadRoomList{}
	var roomIDs = roomIDList.GetRoomIDs()
	if len(roomIDs) == 0 {
		log.Println("room ids is empty")
		return make([]*RoomInfo, 0)
	}

	var roomInfos = make([]*RoomInfo, 0, len(roomIDs))

	conn.Send("MULTI")
	for _, roomID := range roomIDs {
		conn.Send("HMGET", gconst.RoomTablePrefix+roomID, "roomNumber", "gameServerID", "roomConfigID", "timeStamp", "lastActiveTime")
		conn.Send("HMGET", gconst.GsRoomTablePrefix+roomID, "players", "state", "hrStartted")
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

		var gameServerURL = getGameServerURL(gameServerID)
		if gameServerURL == "" {
			continue
		}

		var roomInfo = &RoomInfo{}
		roomInfo.RoomID = &roomID
		roomInfo.RoomNumber = &roomNunmber
		roomInfo.GameServerURL = &gameServerURL
		roomInfo.TimeStamp = &timeStamp

		lastActiveTimeInt32, _ := strconv.Atoi(lastActiveTimeStr)
		lastActiveTimeUint32 := uint32(lastActiveTimeInt32)
		roomInfo.LastActiveTime = &lastActiveTimeUint32

		roomConfig, ok := roomConfigs[roomConfigID]
		if ok {
			roomInfo.Config = &roomConfig
		}

		roomInfos = append(roomInfos, roomInfo)

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

	// roomID2RoomConfig := loadRoomConfig(roomConfigID2RoomID)

	return roomInfos
}

// 1.获取对应房间配置ID、房间号、游戏服务器ID
// 2.用房间配置ID获取房间配置
// 3.用游戏服务器ID获取游戏服务器url
func onMessageGetRooms(user *User, accessoryMessage *AccessoryMessage) {
	conn := pool.Get()
	defer conn.Close()

	var roomInfos = loadRoomInfos(user.userID())

	var msgLoadRoomListRsp = &MsgLoadRoomListRsp{}
	msgLoadRoomListRsp.RoomInfos = roomInfos
	var result = int32(MsgError_ErrSuccess)
	msgLoadRoomListRsp.Result = &result
	replyLoadRooms(user, msgLoadRoomListRsp)
}
