package share

import (
	"gconst"
	"lobbyserver/lobby"
	"net/http"
	"strconv"
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/julienschmidt/httprouter"

	log "github.com/sirupsen/logrus"

	"github.com/golang/protobuf/proto"
)

func replyGetShareInfo(w http.ResponseWriter, shareInfo *lobby.MsgShareInfo) {
	bytes, err := proto.Marshal(shareInfo)
	if err != nil {
		log.Panic("reply msg, marshal msg failed")
		return
	}

	w.Write(bytes)
}

func getMultimedia(sence int, mediaType int) string {
	conn := lobby.Pool().Get()
	defer conn.Close()

	shareServer, err := redis.String(conn.Do("HGET", gconst.LobbyConfig, "shareServer"))
	if err != nil {
		log.Error("getMultimedia, load shereServer url error:", err)
		return ""
	}

	key := fmt.Sprintf(gconst.LobbyShareMedia, sence, mediaType)

	mediaName, err := redis.String(conn.Do("SRANDMEMBER", key))
	if err != nil {
		log.Error("getMultimedia, load mediaName error:", err)
		return ""
	}

	url := fmt.Sprintf("https://%s/%d/%d/%s", shareServer, sence, mediaType, mediaName)

	return url
}

func getText(sence int) string {
	conn := lobby.Pool().Get()
	defer conn.Close()

	key := fmt.Sprintf(gconst.LobbyShareText, sence)
	text, err := redis.String(conn.Do("SRANDMEMBER", key))
	if err != nil {
		log.Error("getText, load text error:", err)
		return ""
	}

	return text
}

func handlerGetShareInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("userID")

	log.Println("handlerGetShareInfo, userID:", userID)
	// 1-游戏分享
	senceStr := r.URL.Query().Get("sence")
	// 1-图片 2-小视频 3-动图
	mediaTypeStr := r.URL.Query().Get("mediaType")
	// // 1-朋友圈 2-好友
	shareTypeStr := r.URL.Query().Get("shareType")

	sence, _ := strconv.Atoi(senceStr)
	mediaType, _ := strconv.Atoi(mediaTypeStr)
	shareType, _ := strconv.Atoi(shareTypeStr)

	result := int32(0)
	shareInfo := &lobby.MsgShareInfo{}

	if sence == 0 {
		result = int32(lobby.MsgError_ErrRequestInvalidParam)
		shareInfo.Result = &result
		replyGetShareInfo(w, shareInfo)

		return
	}

	if mediaType == 0 {
		result = int32(lobby.MsgError_ErrRequestInvalidParam)
		shareInfo.Result = &result
		replyGetShareInfo(w, shareInfo)

		return
	}

	if shareType == 0 {
		result = int32(lobby.MsgError_ErrRequestInvalidParam)
		shareInfo.Result = &result
		replyGetShareInfo(w, shareInfo)

		return
	}

	multimedia := getMultimedia(sence, mediaType)
	if multimedia == "" {
		log.Errorf("handlerGetShareInfo multimedia is empty for sence:%d, mediaType:%d, shareType:%d")
	}

	text := getText(sence)
	if text == "" {
		log.Errorf("handlerGetShareInfo text is emtpy")
	}

	shareInfo.Result = &result
	shareInfo.Text = &text
	shareInfo.Multimedia = &multimedia

	replyGetShareInfo(w, shareInfo)
}
