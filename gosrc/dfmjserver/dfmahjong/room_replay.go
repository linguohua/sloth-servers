package dfmahjong

import (
	"fmt"
	"gconst"
	"mahjong"

	log "github.com/sirupsen/logrus"

	"github.com/garyburd/redigo/redis"

	"github.com/golang/protobuf/proto"

	uuid "github.com/satori/go.uuid"
)

// dump2Redis 输出到redis，注意此时不能再引用room中的players
func (lc *LoopContext) dump2Redis(s *SPlaying) {
	if lc.s.room.isUlimitRound {
		return
	}

	if lc.recorder.Actions == nil {
		lc.actionList2Actions()
	}

	newUUID, err := uuid.NewV4()
	if err != nil {
		log.Panicln("dump2Redis failed, new uuid error:", err)
	}

	recordID := fmt.Sprintf("%s", newUUID)

	buf, err := proto.Marshal(lc.recorder)
	if err != nil {
		log.Panicln(err)
	}

	// 获取redis链接，并退出函数时释放
	conn := pool.Get()
	defer conn.Close()

	// 记录回播记录
	conn.Send("MULTI")
	conn.Send("HMSET", gconst.GameServerMJRecorderTablePrefix+recordID, "d", buf, "r", s.room.ID, "cid", s.room.configID)
	// key 48小时后过期
	conn.Send("EXPIRE", gconst.GameServerMJRecorderTablePrefix+recordID, 48*60*60)

	_, err = conn.Do("EXEC")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("dump2Redis, new record id:", recordID)

	// 记录回播房间
	// 如果对象为空，则新建一个
	if s.room.msgReplayRoom == nil {
		s.room.msgReplayRoom = lc.createReplayRoom(s.room)
	}

	lc.updateReplayRoom(conn, s.room, recordID)
	roomID := s.room.ID

	// 为每一个玩家添加回播房间记录
	for _, p := range s.room.msgReplayRoom.Players {
		// 检查用户是否位于房间的回播中
		userID := p.GetUserID()
		exist, err := redis.Int(conn.Do("sismember", gconst.GameServerReplayRoomsReferenceSetPrefix+roomID, userID))
		if err != nil {
			continue
		}

		if exist == 1 {
			continue
		}

		// 不存在则增加到list中
		conn.Send("MULTI")
		conn.Send("RPUSH", gconst.GameServerMJReplayRoomListPrefix+userID, roomID)
		// 确保最多50个
		conn.Send("LTRIM", gconst.GameServerMJReplayRoomListPrefix+userID, -1, -50)

		conn.Send("SADD", gconst.GameServerReplayRoomsReferenceSetPrefix+roomID, userID)
		conn.Do("EXEC")
	}

	// key 48小时后过期
	conn.Do("EXPIRE", gconst.GameServerReplayRoomsReferenceSetPrefix+roomID, 48*60*60)
}

// updateReplayRoom 更新回播房间记录
func (lc *LoopContext) updateReplayRoom(conn redis.Conn, room *Room, recordID string) {
	// 房间持有replayRoom对象
	msgReplayRoom := room.msgReplayRoom

	var replayRecordSummary = lc.replayRecordSummary
	replayRecordSummary.RecordUUID = &recordID

	msgReplayRoom.Records = append(msgReplayRoom.Records, replayRecordSummary)

	var endTime = unixTimeInMinutes()
	msgReplayRoom.EndTime = &endTime

	buf, err := proto.Marshal(msgReplayRoom)
	if err != nil {
		log.Println("updateReplayRoom, marshal error:", err)
		return
	}

	unixTimeInMinutes32 := unixTimeInMinutes()

	userIDs := make([]string, len(msgReplayRoom.Players))
	for i, p := range msgReplayRoom.Players {
		userIDs[i] = p.GetUserID()
	}

	roomID := room.ID
	// 写redis
	// 更新回播房间信息
	conn.Send("MULTI")
	conn.Send("HMSET", gconst.GameServerMJReplayRoomTablePrefix+roomID,
		"d", buf, "rrt", int(myRoomType), "date", unixTimeInMinutes32)
	// key 48小时后过期
	conn.Send("EXPIRE", gconst.GameServerMJReplayRoomTablePrefix+roomID, 48*60*60)

	// for _, u := range userIDs {
	// 	conn.Send("SADD", gconst.GameServerReplayRoomsReferenceSetPrefix+roomID, u)
	// }
	// // key 48小时后过期
	// conn.Send("EXPIRE", gconst.GameServerReplayRoomsReferenceSetPrefix+roomID, 48*60*60)
	conn.Do("EXEC")
}

func (lc *LoopContext) createReplayRoom(room *Room) *mahjong.MsgReplayRoom {
	var msgReplayRoom = &mahjong.MsgReplayRoom{}
	var recordRoomType32 = int32(myRoomType)
	msgReplayRoom.RecordRoomType = &recordRoomType32
	msgReplayRoom.RoomNumber = &room.roomNumber
	var startTime = unixTimeInMinutes()
	msgReplayRoom.StartTime = &startTime
	msgReplayRoom.OwnerUserID = &room.ownerID

	// 玩家列表
	var replayPlayers = make([]*mahjong.MsgReplayPlayerInfo, len(room.players))
	for i, p := range room.players {
		rp := &mahjong.MsgReplayPlayerInfo{}
		var chairID32 = int32(p.chairID)
		rp.ChairID = &chairID32
		var userID = p.userID()
		rp.UserID = &userID
		var nick = p.user.getInfo().nick
		rp.Nick = &nick
		var sex = p.user.getInfo().sex
		rp.Sex = &sex
		var headIconURL = p.user.getInfo().headIconURI
		rp.HeadIconURI = &headIconURL

		var totalScore32 = int32(p.gStatis.roundScore)
		rp.TotalScore = &totalScore32

		var avatarID = int32(p.user.getInfo().avatarID)
		rp.AvatarID = &avatarID

		replayPlayers[i] = rp
	}

	log.Println("createReplayRoom, players count:", len(replayPlayers))
	msgReplayRoom.Players = replayPlayers

	return msgReplayRoom
}
