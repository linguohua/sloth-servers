package lobby

import (
	"encoding/json"
	"fmt"
	"net/http"
	"gconst"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	// "fmt"

	// "github.com/golang/protobuf/proto"

	"github.com/garyburd/redigo/redis"
)

// conn.Send("EXPIREAT", keyGroupBigWinnerStats, expireat)
// conn.Send("EXPIREAT", keyGroupBigHandStats, expireat)
// conn.Send("EXPIREAT", keyGroupSpecificRoomBigWinnerStats, expireat)
// conn.Send("EXPIREAT", keyGroupSpecificRoomBigHandStats, expireat)
// conn.Send("EXPIREAT", keyGroupStatsUpdateTime, expireat)

// BigWinner 大赢家
type BigWinner struct {
	UserID       string `json:"userID"`
	DateTime     int64  `json:"dateTime"`
	Nick         string `json:"nick"`
	BigHandCount int    `json:"bigHandCount"`
	BigWinCount  int    `json:"bigWinCount"`
	IsConfirm    int    `json:"isConfirm"`
}

// ReplyBigWinnerList 大赢家列表
type ReplyBigWinnerList struct {
	ErrorCode     int          `json:"errorCode"`
	ErrorMsg      string       `json:"errorMsg"`
	BigWinnerList []*BigWinner `json:"bigWinnerList"`
	Cursor        int          `json:"cursor"`
}

func replyLoadBigWinner(w http.ResponseWriter, errCode int32, errMsg string) {
	reply := &ReplyBigWinnerList{}
	reply.ErrorCode = int(errCode)
	reply.ErrorMsg = errMsg
	reply.BigWinnerList = nil

	b, err := json.Marshal(reply)
	if err != nil {
		log.Panicln("genericReply, json marshal error:", err)
		return
	}

	w.Write(b)
}

func handleLoadGroupBigWinner(w http.ResponseWriter, r *http.Request, userID string) {
	log.Println("handleLoadReplayRooms call, userID:", userID)
	groupID := r.URL.Query().Get("groupID")
	cursorStr := r.URL.Query().Get("cursor")
	roomType := r.URL.Query().Get("roomType")
	loadCountStr := r.URL.Query().Get("loadCount")
	dayStr := r.URL.Query().Get("dd")
	log.Printf("handleLoadGroupBigWinner, groupID:%s, cursorStr:%s, roomType:%s, loadCountStr:%s, dayStr:%s", groupID, cursorStr, roomType, loadCountStr, dayStr)

	if groupID == "" {
		log.Println("handleLoadGroupBigWinner, groupID is empty")
		var errCode = int32(MsgError_ErrRequestInvalidParam)
		replyLoadBigWinner(w, errCode, ErrorString[errCode])
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
			log.Println("handleLoadGroupBigWinner, err:", err)
		}

		loadCount = count
	}

	day := 0
	if dayStr != "" {
		day, _ = strconv.Atoi(dayStr)
	}

	t := time.Now().Local()
	var dayTime = time.Date(t.Year(), t.Month(), t.Day()-day, 0, 0, 0, 0, t.Location())
	var dd = dayTime.Format("20060102")

	var keyStats = fmt.Sprintf(gconst.GroupBigWinnerStats, groupID, dd)
	if roomType != "" {
		keyStats = fmt.Sprintf(gconst.GroupSpecificRoomBigWinnerStats, groupID, roomType, dd)
	}

	log.Println("keyStats:", keyStats)
	// 获取redis链接，并退出函数时释放
	conn := pool.Get()
	defer conn.Close()

	totalCount, err := redis.Int(conn.Do("ZCARD", keyStats))
	if err != nil {
		log.Println("handleLoadGroupBigWinner, parame error", err)
		var errCode = int32(MsgError_ErrDatabase)
		replyLoadBigWinner(w, errCode, ErrorString[errCode])
		return
	}

	// 如果没有， 返回空列表
	if cursor > totalCount {
		log.Printf("handleLoadGroupBigWinner cursor %d > totalCount %d", cursor, totalCount)
		reply := &ReplyBigWinnerList{}
		reply.ErrorCode = int(0)
		reply.ErrorMsg = "ok"
		reply.Cursor = 0
		reply.BigWinnerList = make([]*BigWinner, 0)

		b, err := json.Marshal(reply)
		if err != nil {
			log.Panicln("genericReply, json marshal error:", err)
			return
		}

		w.Write(b)
		return
	}

	// if cursor + loadCount > totalCount {
	// 	loadCount = totalCount - cursor
	// }

	userIDs, err := redis.Strings(conn.Do("ZREVRANGE", keyStats, cursor, cursor+loadCount))
	if err != nil {
		log.Println("handleLoadGroupBigWinner, error:", err)
		var errCode = int32(MsgError_ErrDatabase)
		replyLoadBigWinner(w, errCode, ErrorString[errCode])
		return
	}

	if len(userIDs) < loadCount {
		cursor = 0
	} else {
		cursor = cursor + len(userIDs)
	}

	keyGroupStatsUpdateTime := fmt.Sprintf(gconst.GroupStatsUpdateTime, groupID, dd)
	keyGroupStatsConfirm := fmt.Sprintf(gconst.GroupStatsConfirm, groupID, dd)

	var key = fmt.Sprintf(gconst.GroupBigHandStats, groupID, dd)
	if roomType != "" {
		key = fmt.Sprintf(gconst.GroupSpecificRoomBigHandStats, groupID, roomType, dd)
		keyGroupStatsConfirm = fmt.Sprintf(gconst.GroupStatsSpecificeRoomConfirm, groupID, roomType, dd)
	}

	log.Println("userIDs:", userIDs)

	conn.Send("MULTI")
	for _, userID := range userIDs {
		conn.Send("ZSCORE", keyStats, userID)
		conn.Send("ZSCORE", key, userID)
		conn.Send("HGET", gconst.AsUserTablePrefix+userID, "Nick")
		conn.Send("HGET", keyGroupStatsUpdateTime, userID)
		conn.Send("HGET", keyGroupStatsConfirm, userID)
	}

	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		log.Println("handleLoadGroupBigWinner err:", err)
		var errCode = int32(MsgError_ErrDatabase)
		replyLoadBigWinner(w, errCode, ErrorString[errCode])
		return
	}

	log.Println("values:", values)

	bigWinnerList := make([]*BigWinner, 0, len(userIDs))
	for i, userID := range userIDs {
		bigWinCount, err := redis.Int(values[i*5], nil)
		if err != nil {
			log.Println("get bigWinCount error:", err)
			continue
		}

		bigHandCount, err := redis.Int(values[i*5+1], nil)
		if err != nil {
			log.Println("get bigHandCount error:", err)
			continue
		}

		nick, err := redis.String(values[i*5+2], nil)
		if err != nil {
			log.Println("get nick error:", err)
			continue
		}

		dateTime, err := redis.Int64(values[i*5+3], nil)
		if err != nil {
			log.Println("get dateTime error:", err)
			continue
		}

		isConfirm, _ := redis.Int(values[i*5+4], nil)

		bigWinner := &BigWinner{}
		bigWinner.BigWinCount = bigWinCount
		bigWinner.BigHandCount = bigHandCount
		bigWinner.Nick = nick
		bigWinner.UserID = userID
		bigWinner.DateTime = dateTime
		bigWinner.IsConfirm = isConfirm

		bigWinnerList = append(bigWinnerList, bigWinner)
	}

	reply := &ReplyBigWinnerList{}
	reply.ErrorCode = 0
	reply.ErrorMsg = "ok"
	reply.BigWinnerList = bigWinnerList
	reply.Cursor = cursor

	b, err := json.Marshal(reply)
	if err != nil {
		log.Panicln("genericReply, json marshal error:", err)
		return
	}

	w.Write(b)
}
