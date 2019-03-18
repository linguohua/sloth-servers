package lobby

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"gconst"
	"time"
	"strconv"
	"fmt"

	"github.com/golang/protobuf/proto"

	"github.com/garyburd/redigo/redis"
)

func handleLoadGroupReplayRooms(w http.ResponseWriter, r *http.Request, userID string) {
	groupID := r.URL.Query().Get("groupID")
	cursorStr := r.URL.Query().Get("cursor")
	roomType := r.URL.Query().Get("roomType")
	loadCountStr := r.URL.Query().Get("loadCount")
	isSingleUser := r.URL.Query().Get("isSingleUser")
	log.Printf("handleLoadReplayRooms call, userID:%s, groupID:%s, cursorStr:%s,roomType:%s, loadCountStr:%s, isSingleUser:%s", userID, groupID, cursorStr, roomType, loadCountStr, isSingleUser)

	if groupID == "" {
		log.Println("handleLoadGroupReplayRooms, groupID is empty")
		var errCode = int32(MsgError_ErrRequestInvalidParam)
		replyDeleteRoom(w, errCode, ErrorString[errCode])
		return
	}

	cursor := 0
	if cursorStr != "" {
		cursor, _ = strconv.Atoi(cursorStr)
	}

	loadCount := 10
	if loadCountStr != "" {
		count, err := strconv.Atoi(loadCountStr)
		if err != nil {
			log.Println("handleLoadGroupReplayRooms, err:", err)
		}

		loadCount = count
	}

	isAll := true
	if isSingleUser == "true" {
		isAll = false
	}

	t := time.Now().Local()
	var startTime = time.Date(t.Year(), t.Month(), t.Day()-2, 0, 0, 0, 0, t.Location())
	var startTimeUTC = int32(startTime.Unix())
	var endTimeUTC = int32(t.Unix())

	log.Printf("startTimeUTC:%d, endTimeUTC:%d", startTimeUTC, endTimeUTC)
	// 获取redis链接，并退出函数时释放
	conn := pool.Get()
	defer conn.Close()

	totalCount, err := redis.Int(conn.Do("ZCOUNT", gconst.GroupRoomSortSetPrefix + groupID, startTimeUTC , endTimeUTC))
	if err != nil {
		log.Println("handleLoadGroupReplayRooms, parame error", err)
		var errCode = int32(MsgError_ErrDatabase)
		replyDeleteRoom(w, errCode, ErrorString[errCode])
		return
	}

	// 返回空列表
	if cursor > totalCount {
		log.Printf("cursor %d > totalCount %d", cursor, totalCount)
		replayRoomsReply := &MsgAccLoadReplayRoomsReply{}
		cursorInt32 := int32(0)
		replayRoomsReply.Cursor = &cursorInt32
		replayRoomsReply.ReplayRooms = make([]*MsgAccReplayRoom, 0)

		bytes, err := proto.Marshal(replayRoomsReply)
		if err != nil {
			log.Println("handleLoadGroupReplayRooms, Marshal err:", err)
			return
		}

		writeHTTPBodyWithGzip(w, r, bytes)
		return
	}

	// if cursor + loadCount > totalCount {
	// 	loadCount = totalCount - cursor
	// }

	roomIDs, err := redis.Strings(conn.Do("ZREVRANGE", gconst.GroupRoomSortSetPrefix+ groupID, cursor, cursor+loadCount))
	if err != nil {
		log.Println("handleLoadGroupReplayRooms, error:", err)
		var errCode = int32(MsgError_ErrDatabase)
		replyDeleteRoom(w, errCode, ErrorString[errCode])
		return
	}

	log.Println("handleLoadGroupReplayRooms, roomIDs:", roomIDs)

	var reply *MsgAccLoadReplayRoomsReply
	if isAll {
		reply = loadGroupReplayRoomsByIDs(roomIDs, roomType, conn)
	} else {
		reply = loadGroupReplayRoomsByIDAndUserID(roomIDs, roomType, userID, conn )
	}

	cursorInt32 := int32(cursor)
	if len(roomIDs) < loadCount {
		cursorInt32 = 0
	} else {
		cursorInt32 = int32(cursor + len(roomIDs))
	}

	// totalCursorInt := int32(totalCount)

	reply.Cursor = &cursorInt32
	// reply.TotalCursor = &totalCursorInt

	bytes, err := proto.Marshal(reply)
	if err != nil {
		log.Println("handleLoadGroupReplayRooms, Marshal err:", err)
		return
	}

	writeHTTPBodyWithGzip(w, r, bytes)
}


func loadGroupReplayRoomsByIDs(replayRoomIDs []string, roomType string, conn redis.Conn) *MsgAccLoadReplayRoomsReply {

	// 加载所有回播房间记录概要
	conn.Send("MULTI")
	for _, rr := range replayRoomIDs {
		// "d" 二进制数据， "rrt" 回播房间类型：1是大丰，2是东台，3是盐城等等，具体看game_replay.proto定义
		conn.Send("HMGET", gconst.MJReplayRoomTablePrefix+rr, "d", "rrt")
	}
	values, err := redis.Values(conn.Do("EXEC"))

	if err != nil {
		log.Println("handleLoadGroupReplayRooms, values error:", err)
		return nil
	}

	log.Println("handleLoadGroupReplayRooms, values length:", len(values))

	msgLoadReplayRoomReply := &MsgAccLoadReplayRoomsReply{}
	msgReplayRooms := make([]*MsgAccReplayRoom, 0, len(values))
	for i := 0; i < len(values); i++ {
		vx, err := redis.Values(values[i], nil)
		if err != nil {
			log.Println("handleLoadGroupReplayRooms, Values err:", err)
			continue
		}

		bytes, err := redis.Bytes(vx[0], nil)
		if err != nil {
			log.Printf("loadGroupReplayRoomsByIDs, room %s bytes err:%v", replayRoomIDs[i], err)
			continue
		}

		rrtInt, err := redis.Int(vx[1], nil)
		if err != nil {
			log.Println("handleLoadGroupReplayRooms, int err:", err)
			continue
		}

		rrtStr := fmt.Sprintf("%d", rrtInt)
		if roomType != "" &&  roomType != rrtStr {
			log.Println("roomType != rrtInt")
			continue
		}


		var replayRoom = &MsgReplayRoom{}
		err = proto.Unmarshal(bytes, replayRoom)
		if err != nil {
			log.Println("handleLoadReplayRooms, parser replay room error :", err)
			continue
		}

		t := time.Now().Local()
		var startTime = time.Date(t.Year(), t.Month(), t.Day()-1, 0, 0, 0, 0, t.Location())
		var startTimeUTC = uint32(startTime.Unix() / 60)
		var endTime = replayRoom.GetEndTime()
		if endTime < startTimeUTC {
			log.Println("replay room is data out")
			continue
		}


		msgReplayRoom := &MsgAccReplayRoom{}
		msgReplayRoom.ReplayRoomBytes = bytes
		var rrt32 = int32(rrtInt)
		msgReplayRoom.RecordRoomType = &rrt32

		// loadReplayPlayerHeadIconURI(msgReplayRoom.GetPlayers())
		msgReplayRooms = append(msgReplayRooms, msgReplayRoom)
	}

	msgLoadReplayRoomReply.ReplayRooms = msgReplayRooms
	return msgLoadReplayRoomReply
}


func loadGroupReplayRoomsByIDAndUserID(replayRoomIDs []string, roomType string, userID string, conn redis.Conn) *MsgAccLoadReplayRoomsReply {

		// 加载所有回播房间记录概要
		conn.Send("MULTI")
		for _, rr := range replayRoomIDs {
			// "d" 二进制数据， "rrt" 回播房间类型：1是大丰，2是东台，3是盐城等等，具体看game_replay.proto定义
			conn.Send("HMGET", gconst.MJReplayRoomTablePrefix+rr, "d", "rrt")
		}
		values, err := redis.Values(conn.Do("EXEC"))

		if err != nil {
			log.Println("handleLoadGroupReplayRooms, values error:", err)
			return nil
		}

		log.Println("handleLoadGroupReplayRooms, values length:", len(values))

		msgLoadReplayRoomReply := &MsgAccLoadReplayRoomsReply{}
		msgReplayRooms := make([]*MsgAccReplayRoom, 0, len(values))
		for i := 0; i < len(values); i++ {
			vx, err := redis.Values(values[i], nil)
			if err != nil {
				log.Println("handleLoadGroupReplayRooms, Values err:", err)
				continue
			}

			bytes, err := redis.Bytes(vx[0], nil)
			if err != nil {
				log.Println("loadGroupReplayRoomsByIDAndUserID, bytes err:", err)
				continue
			}

			rrtInt, err := redis.Int(vx[1], nil)
			if err != nil {
				log.Println("handleLoadGroupReplayRooms, int err:", err)
				continue
			}

			rrtStr := fmt.Sprintf("%d", rrtInt)
			if roomType != "" &&  roomType != rrtStr {
				log.Printf("roomType != rrtInt, roomType:%s, rrtStr:%s", roomType, rrtStr)
				continue
			}


			var replayRoom = &MsgReplayRoom{}
			err = proto.Unmarshal(bytes, replayRoom)
			if err != nil {
				log.Println("handleLoadGroupReplayRooms, parser replay room error :", err)
				continue
			}

			if !isUserInRoom(userID, replayRoom.GetPlayers()) {
				continue
			}

			t := time.Now().Local()
			var startTime = time.Date(t.Year(), t.Month(), t.Day()-1, 0, 0, 0, 0, t.Location())
			var startTimeUTC = uint32(startTime.Unix() / 60)
			var endTime = replayRoom.GetEndTime()
			if endTime < startTimeUTC {
				log.Println("replay room is data out")
				continue
			}


			msgReplayRoom := &MsgAccReplayRoom{}
			msgReplayRoom.ReplayRoomBytes = bytes
			var rrt32 = int32(rrtInt)
			msgReplayRoom.RecordRoomType = &rrt32

			// loadReplayPlayerHeadIconURI(msgReplayRoom.GetPlayers())
			msgReplayRooms = append(msgReplayRooms, msgReplayRoom)
		}

		msgLoadReplayRoomReply.ReplayRooms = msgReplayRooms

		return msgLoadReplayRoomReply
}

func isUserInRoom(userID string, playerInfos []*MsgReplayPlayerInfo) bool{
	for _, playerInfo := range playerInfos {
		if playerInfo.GetUserID() == userID {
			return true
		}
	}

	return false

}
