package replay

import (
	"gconst"
	"lobbyserver/lobby"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"

	log "github.com/sirupsen/logrus"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
)

// func loadReplayPlayerHeadIconURI(players []*gconst.MsgReplayPlayerInfo) {
// 	if players == nil || len(players) == 0 {
// 		return
// 	}

// 	conn := pool.Get()
// 	defer conn.Close()

// 	conn.Send("MULTI")

// 	for _, player := range players {
// 		conn.Send("HMGET", gconst.LobbyUserTablePrefix+player.GetUserID(), "userSex", "userLogo")
// 	}

// 	values, err := redis.Values(conn.Do("EXEC"))
// 	if err != nil {
// 		log.Println("handleLoadUserHeadIconURI， load user head icon from redis err:", err)
// 		return
// 	}

// 	for index, v := range values {
// 		fileds, _ := redis.Strings(v, nil)
// 		var sexString = fileds[0]
// 		var headIconURI = fileds[1]

// 		sexUint64, _ := strconv.ParseUint(sexString, 10, 32)
// 		var sex = uint32(sexUint64)

// 		var player = players[index]
// 		player.Sex = &sex
// 		player.HeadIconURI = &headIconURI
// 	}
// }

func handleLoadReplayRooms(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := r.Context().Value("userID").(string)
	log.Println("handleLoadReplayRooms call, userID:", userID)
	replayType := r.URL.Query().Get("rt")

	if replayType != "1" {
		log.Println("handleLoadReplayRooms,not support replay type:", replayType)
		return
	}

	// 获取redis链接，并退出函数时释放
	conn := lobby.Pool().Get()
	defer conn.Close()

	// 取出玩家的回播房间列表
	replayRoomsStr, err := redis.String(conn.Do("HGET", gconst.LobbyPlayerTablePrefix+userID, "rr"))
	if err != nil {
		log.Println("handleLoadReplayRooms, err:", err)
		return
	}

	replayRooms := strings.Split(replayRoomsStr, ",")
	if len(replayRooms) < 1 {
		log.Println("handleLoadReplayRooms, replay room list is empty")
		return
	}

	bytes := loadReplayRoomsByIDs(replayRooms, conn)

	writeHTTPBodyWithGzip(w, r, bytes)
}

func loadReplayRoomsByIDs(replayRoomIDs []string, conn redis.Conn) []byte {

	// 加载所有回播房间记录概要
	conn.Send("MULTI")
	for _, rr := range replayRoomIDs {
		// "d" 二进制数据， "rrt" 回播房间类型：1是大丰，2是东台，3是盐城等等，具体看game_replay.proto定义
		conn.Send("HMGET", gconst.GameServerMJReplayRoomTablePrefix+rr, "d", "rrt")
	}
	values, err := redis.Values(conn.Do("EXEC"))

	if err != nil {
		log.Println("handleLoadReplayRooms, values error:", err)
		return nil
	}

	log.Println("handleLoadReplayRooms, values length:", len(values))

	msgLoadReplayRoomReply := &lobby.MsgAccLoadReplayRoomsReply{}
	msgReplayRooms := make([]*lobby.MsgAccReplayRoom, 0, len(values))
	for i := 0; i < len(values); i++ {
		vx, err := redis.Values(values[i], nil)
		if err != nil {
			log.Println("handleLoadReplayRooms, Values err:", err)
			continue
		}

		bytes, err := redis.Bytes(vx[0], nil)
		if err != nil {
			log.Println("handleLoadReplayRooms, bytes err:", err)
			continue
		}

		rrtInt, err := redis.Int(vx[1], nil)
		if err != nil {
			log.Println("handleLoadReplayRooms, int err:", err)
			continue
		}

		var replayRoom = &lobby.MsgReplayRoom{}
		err = proto.Unmarshal(bytes, replayRoom)
		if err != nil {
			log.Println("handleLoadReplayRooms, parser replay room error :", err)
			continue
		}

		// var roomConfigID = replayRoom.GetRoomConfigID()
		// var roomConfigJSON = lobby.GetRoomConfig(roomConfigID)
		// if roomConfigJSON != nil && roomConfigJSON.Race == 1 {
		// 	log.Println("Not need to load Race room")
		// 	continue
		// }

		t := time.Now().Local()
		var startTime = time.Date(t.Year(), t.Month(), t.Day()-1, 0, 0, 0, 0, t.Location())
		var startTimeUTC = uint32(startTime.Unix() / 60)
		var endTime = replayRoom.GetEndTime()
		if endTime < startTimeUTC {
			log.Println("replay room is data out")
			continue
		}

		msgReplayRoom := &lobby.MsgAccReplayRoom{}
		msgReplayRoom.ReplayRoomBytes = bytes
		var rrt32 = int32(rrtInt)
		msgReplayRoom.RecordRoomType = &rrt32

		// loadReplayPlayerHeadIconURI(msgReplayRoom.GetPlayers())
		msgReplayRooms = append(msgReplayRooms, msgReplayRoom)
	}

	msgLoadReplayRoomReply.ReplayRooms = msgReplayRooms

	bytes, err := proto.Marshal(msgLoadReplayRoomReply)
	if err != nil {
		log.Println("handleLoadReplayRooms, Marshal err:", err)
		return nil
	}

	return bytes
}
