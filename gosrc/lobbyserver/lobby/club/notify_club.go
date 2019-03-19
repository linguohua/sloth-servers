package club

import (
	// "club"
	"encoding/json"
	"fmt"
	"gconst"
	"strconv"

	"github.com/garyburd/redigo/redis"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

// CFET_None = 0; // 无效事件
// CFET_Add_By_Shop = 1; // 商城充值
// CFET_Award_By_System = 3; // 系统奖励
// CFET_Gift_By_System = 4; // 系统赠送
// CFET_Reduce_By_Room = 5; // 开房扣除
// CFET_Add_By_Room = 6; // 开房返还
const (
	EventTimeOut    = 60
	CreateClubRoom  = 1
	DeleteClubRoom  = 2
	ClubStateChange = 3
)

// 通知俱乐部，钻石改变
func notifyClubFundAddByShop(amount int, total int, userID string, clubID string) {
	// club.OnFundEvent(et ClubFundEventType, userID string, clubID string, amount int, total int)
	//chost.clubRoomsListener.OnFundEvent(club.ClubFundEventType_CFET_Add_By_Shop, userID, clubID, amount, total)
}

func notifyClubFundReduceByRoom(amount int, total int, userID string, clubID string) {
	// club.OnFundEvent(et ClubFundEventType, userID string, clubID string, amount int, total int)
	//chost.clubRoomsListener.OnFundEvent(club.ClubFundEventType_CFET_Reduce_By_Room, userID, clubID, amount, total)
}

func notifyClubFundAddByRoom(amount int, total int, userID string, clubID string) {
	// club.OnFundEvent(et ClubFundEventType, userID string, clubID string, amount int, total int)
	//chost.clubRoomsListener.OnFundEvent(club.ClubFundEventType_CFET_Add_By_Room, userID, clubID, amount, total)
}

// 创建、删除俱乐部房间发通知给牌友群
func publishRoomChangeMessage2Group(clubID string, roomID string, operateType int) {
	conn := pool.Get()
	defer conn.Close()

	fields, err := redis.Strings(conn.Do("HMGET", gconst.RoomTablePrefix+roomID, "ownerID", "roomType"))
	if err != nil {
		log.Println("publishRoomChangeMessage2Group, load ownerID and roomType error:", err)
	}

	ownerID := fields[0]
	roomType, _ := strconv.Atoi(fields[1])
	var playerIDs = "[]"
	// 台安的用户不在房间内
	if ownerID != "" && roomType != int(gconst.RoomType_TacnMJ) && roomType != int(gconst.RoomType_TacnPok) &&
		roomType != int(gconst.RoomType_DDZ) {
		playerIDs = fmt.Sprintf(`["%s"]`, ownerID)
	}

	var content = `{"clubid":"%s", "roomid":"%s", "type":%d, "playernum":%d, "playerids":%s}`

	// 创建俱乐部房间，房主默认在房间内, 台安的用户不在房间内，过滤掉
	if operateType == CreateClubRoom && roomType != int(gconst.RoomType_TacnMJ) && roomType != int(gconst.RoomType_TacnPok) && roomType != int(gconst.RoomType_DDZ) {
		content = fmt.Sprintf(content, clubID, roomID, operateType, 1, playerIDs)
	} else {
		content = fmt.Sprintf(content, clubID, roomID, operateType, 0, playerIDs)
	}

	var event = gconst.RedisKeyEventRoomInfoChange

	id, _ := uuid.NewV4()
	key := fmt.Sprintf("%s:%s", event, id)

	conn.Send("MULTI")
	conn.Send("SET", key, content, "EX", EventTimeOut)
	conn.Send("PUBLISH", event, key)
	_, err = conn.Do("EXEC")
	if err != nil {
		log.Printf("publishRoomChangeMessage2Group err %v", err)
		return
	}

	log.Println("publishRoomChangeMessage2Group:", content)
}

// 房间状态改变，发通知给牌友群
func publishRoomStateChange2Group(clubID string, roomID string, operateType int, playerNum int, playerIDs []string) {
	conn := pool.Get()
	defer conn.Close()

	buf, err := json.Marshal(playerIDs)
	if err != nil {
		log.Println("publishRoomStateChange2Group error:", err)
		return
	}

	var content = `{"clubid":"%s", "roomid":"%s", "type":%d, "playernum":%d, "playerids":%s}`
	content = fmt.Sprintf(content, clubID, roomID, operateType, playerNum, string(buf))

	var event = gconst.RedisKeyEventRoomInfoChange

	id, _ := uuid.NewV4()
	key := fmt.Sprintf("%s:%s", event, id)

	conn.Send("MULTI")
	conn.Send("SET", key, content, "EX", EventTimeOut)
	conn.Send("PUBLISH", event, key)
	_, err = conn.Do("EXEC")
	if err != nil {
		log.Printf("publishRoomStateChange2Group err %v", err)
		return
	}

	log.Println("publishRoomStateChange2Group:", content)
}

func publishGameOver2Arena(roomID string, handStarted int, gameOverPlayerStats []*gconst.SSMsgGameOverPlayerStat) {
	var gameOver = 2
	publishRoomState2Arena(roomID, handStarted, gameOver, gameOverPlayerStats)
}

func publishHandBegin2Arena(roomID string, handStarted int) {
	var playerStats = make([]*gconst.SSMsgGameOverPlayerStat, 0)
	var handBegin = 1
	publishRoomState2Arena(roomID, handStarted, handBegin, playerStats)
}

func publishRoomState2Arena(roomID string, handStarted int, roomState int, gameOverPlayerStats []*gconst.SSMsgGameOverPlayerStat) {
	log.Printf("publishRoomState2Arena, roomID:%s, handStarted:%d, roomState:%d, gameOverPlayerStats:%v", roomID, handStarted, roomState, gameOverPlayerStats)
	//玩家积分：
	type PlayerStat struct {
		UserID string `json:"userID"` //玩家ID
		Score  int32  `json:"score"`  // 玩家积分
	}

	playerStats := make([]*PlayerStat, len(gameOverPlayerStats))
	for i, stat := range gameOverPlayerStats {
		playerStat := &PlayerStat{}
		playerStat.Score = stat.GetScore()
		playerStat.UserID = stat.GetUserID()
		playerStats[i] = playerStat
	}

	// 房间状态，目前只在开小局和解散房间发送
	type RoomState struct {
		RoomID       string        `json:"roomID"`
		HandStarted  int           `json:"handStarted"`
		RoomState    int           `json:"roomState"`
		PlayerScores []*PlayerStat `json:"playerScores"`
	}

	var state = &RoomState{}
	state.RoomID = roomID
	state.HandStarted = handStarted
	state.RoomState = roomState
	state.PlayerScores = playerStats

	buf, err := json.Marshal(state)
	if err != nil {
		log.Println("publishRoomState2Arena, error:", err)
		return
	}

	conn := pool.Get()
	defer conn.Close()

	var content = string(buf)

	var event = gconst.RaceRoomStateChange

	// id, _ := uuid.NewV4()
	// key := fmt.Sprintf("%s:%s", event, id)

	// conn.Send("MULTI")
	// conn.Send("SET", key, content, "EX", EventTimeOut)
	// conn.Send("PUBLISH", event, key)
	// _, err = conn.Do("EXEC")
	// if err != nil {
	// 	log.Printf("publishRoomState2Arena err %v", err)
	// 	return
	// }
	conn.Do("PUBLISH", event, content)

	log.Println("publishRoomState2Arena:", content)
}
