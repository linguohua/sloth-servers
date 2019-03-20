package logupload

import (
	"bytes"
	"fmt"
	"gconst"
	"io/ioutil"
	"lobbyserver/lobby"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/garyburd/redigo/redis"
)

const (
	// 5分钟
	maxTime = 60 * 5
)

func handleUploadLogFile(w http.ResponseWriter, r *http.Request, userID string) {
	log.Println("handleUploadLogFile, userID:", userID)
	conn := lobby.Pool().Get()
	defer conn.Close()

	saveLogTimeStr, err := redis.String(conn.Do("HGET", gconst.AsUserTablePrefix+userID, "saveLogTime"))
	if err != nil {
		saveLogTimeStr = "0"
	}

	saveLogTimeInt64, err := strconv.ParseUint(saveLogTimeStr, 10, 32)
	if err != nil {
		saveLogTimeInt64 = 0
	}

	// 检查时间是否超过5分钟
	nowTime := time.Now()
	timeStampInSecond := nowTime.UTC().UnixNano() / int64(time.Second)
	var diff = timeStampInSecond - int64(saveLogTimeInt64)
	if diff < maxTime {
		log.Println("handleUploadLogFile can't upload log file frequently")
		w.Write([]byte(`{"error":1}`))
		return
	}

	if r.ContentLength < 1 {
		log.Println("handleUploadLogFile failed, content length is zero")
		w.Write([]byte(`{"error":1}`))
		return
	}

	buf := bytes.NewBuffer(make([]byte, 0, r.ContentLength))
	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		log.Println("handleUploadLogFile, can't read from body")
		w.Write([]byte(`{"error":1}`))
		return
	}

	body := buf.Bytes()

	var ft = nowTime.Local().Format("2006010215")
	var logPath = fmt.Sprintf("/var/log/%s-%s.log", userID, ft)
	err = ioutil.WriteFile(logPath, body, 0644)
	if err != nil {
		w.Write([]byte(`{"error":1}`))
		log.Println("write log file err:", err)
		return
	}

	conn.Do("HSET", gconst.AsUserTablePrefix+userID, "saveLogTime", timeStampInSecond)

	w.Write([]byte(`{"error":0}`))
}
