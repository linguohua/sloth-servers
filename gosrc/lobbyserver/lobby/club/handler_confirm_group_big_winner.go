package club

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"gconst"
	"time"
	"strconv"
	"encoding/json"
)


func replyConfirmGroupBigWinner(w http.ResponseWriter, errCode int32, errMsg string) {
	type ReplyConfirmBigWinner struct {
		ErrorCode int32 `json:"errorCode"`
		ErrorMsg string `json:"errorMsg"`
	}

	reply := &ReplyConfirmBigWinner{}
	reply.ErrorCode = errCode
	reply.ErrorMsg = errMsg

	b, err := json.Marshal(reply)
	if err != nil {
		log.Panicln("genericReply, json marshal error:", err)
		return
	}

	w.Write(b)
}

func handleConfirmGroupBigWinner(w http.ResponseWriter, r *http.Request, userID string) {
	log.Println("handleLoadReplayRooms call, userID:", userID)
	groupID := r.URL.Query().Get("groupID")
	dayStr := r.URL.Query().Get("dd")
	isConfirmStr := r.URL.Query().Get("isConfirm")
	uID := r.URL.Query().Get("uID")
	roomTypeStr := r.URL.Query().Get("roomType")

	log.Printf("handleConfirmGroupBigWinner, groupID:%s, dayStr:%s, isConfirmStr:%s, uID:%s", groupID, dayStr, isConfirmStr, uID)

	if groupID == "" {
		log.Println("handleLoadGroupReplayRooms, groupID is empty")
		var errCode = int32(MsgError_ErrRequestInvalidParam)
		replyConfirmGroupBigWinner(w, errCode, ErrorString[errCode])
		return
	}

	if dayStr == "" {
		log.Println("handleLoadGroupReplayRooms, dayStr is empty")
		var errCode = int32(MsgError_ErrRequestInvalidParam)
		replyConfirmGroupBigWinner(w, errCode, ErrorString[errCode])
		return
	}

	if isConfirmStr == "" {
		log.Println("handleLoadGroupReplayRooms, isConfirmStr is empty")
		var errCode = int32(MsgError_ErrRequestInvalidParam)
		replyConfirmGroupBigWinner(w, errCode, ErrorString[errCode])
		return
	}

	if roomTypeStr == "" {
		log.Println("handleLoadGroupReplayRooms, roomTypeStr is empty")
		var errCode = int32(MsgError_ErrRequestInvalidParam)
		replyConfirmGroupBigWinner(w, errCode, ErrorString[errCode])
		return
	}

	var err error
	day := 0
	if dayStr != "" {
		day, err = strconv.Atoi(dayStr)
		if err != nil {
			log.Println("handleLoadGroupReplayRooms, parse dayStr error", err)
			var errCode = int32(MsgError_ErrRequestInvalidParam)
			replyConfirmGroupBigWinner(w, errCode, ErrorString[errCode])
			return
		}
	}


	isConfirm := 0
	if isConfirmStr != "" {
		isConfirm, err = strconv.Atoi(isConfirmStr)
		if err != nil {
			log.Println("handleLoadGroupReplayRooms, parse isConfirm error", err)
			var errCode = int32(MsgError_ErrRequestInvalidParam)
			replyConfirmGroupBigWinner(w, errCode, ErrorString[errCode])
			return
		}
	}

	roomType := 0
	if roomTypeStr != "" {
		roomType, err = strconv.Atoi(roomTypeStr)
		if err != nil {
			log.Println("handleLoadGroupReplayRooms, parse roomTypeStr error", err)
			var errCode = int32(MsgError_ErrRequestInvalidParam)
			replyConfirmGroupBigWinner(w, errCode, ErrorString[errCode])
			return
		}
	}


	t := time.Now().Local()
	var dayTime = time.Date(t.Year(), t.Month(), t.Day() - day, 0, 0, 0, 0, t.Location())
	var dd = dayTime.Format("20060102")

	var keyGroupStatsConfirm = fmt.Sprintf(gconst.GroupStatsConfirm, groupID, dd)
	if roomType != 0 {
		keyGroupStatsConfirm = fmt.Sprintf(gconst.GroupStatsSpecificeRoomConfirm, groupID, roomTypeStr, dd)
	}

	log.Println("keyGroupStatsConfirm:", keyGroupStatsConfirm)
	// 获取redis链接，并退出函数时释放
	conn := pool.Get()
	defer conn.Close()

	conn.Do("HSET", keyGroupStatsConfirm, uID, isConfirm)

	replyConfirmGroupBigWinner(w, int32(MsgError_ErrSuccess), "ok")
}