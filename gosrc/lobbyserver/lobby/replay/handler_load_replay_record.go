package replay

import (
	"bytes"
	"compress/gzip"
	"gconst"
	"lobbyserver/lobby"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/golang/protobuf/proto"

	//"mssql"
	"github.com/garyburd/redigo/redis"
)

// func loadPlayerHeadIconURI(players []*gconst.SRMsgPlayerInfo) {
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

func loadReplayRecordFromSQLServer(w http.ResponseWriter, r *http.Request, recordID string) {
	// TODO: llwant mysql
	// conn, err := mssql.StartMssql(config.DbIP, config.DbPort, config.DbUser, config.DbPassword, config.DbName)
	// if err != nil {
	// 	log.Println("handleLoadReplayRecord, StartMssql err:", err)
	// 	return
	// }

	// defer conn.Close()

	// var grcRecord = mssql.LoadGRCRcordFromSQLServer(recordID, conn)

	// var msgHandRecorder = &MsgAccLoadReplayRecord{}
	// msgHandRecorder.ReplayRecordBytes = grcRecord.RecordData
	// roomConfigID := grcRecord.RoomConfigID
	// var roomConfig = roomConfigs[roomConfigID]
	// msgHandRecorder.RoomJSONConfig = &roomConfig

	// //loadPlayerHeadIconURI(msgHandRecorder.GetPlayers())

	// bytesArray, err := proto.Marshal(msgHandRecorder)
	// if err != nil {
	// 	log.Println("handleLoadReplayRecord, marshal err:", err)
	// 	return
	// }

	// writeHTTPBodyWithGzip(w, r, bytesArray)
}

func handleLoadReplayRecord(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	log.Println("handleLoadReplayRecord call, userID:", userID)

	replayType := r.URL.Query().Get("rt")
	if replayType != "1" {
		log.Println("handleLoadReplayRecord,not support replay type:", replayType)
		return
	}

	// 获取redis链接，并退出函数时释放
	conn := lobby.Pool().Get()
	defer conn.Close()

	recordID := r.URL.Query().Get("rid")
	if recordID == "" {
		log.Println("handleLoadReplayRecord, no recordID, now try to use sid")

		sid := r.URL.Query().Get("sid")
		if sid == "" {
			log.Println("handleLoadReplayRecord, no recordID, no sid, can't load")
			return
		}

		recordID, _ = redis.String(conn.Do("HGET", gconst.GameServerMJRecorderTablePrefix+sid, "rid"))
		// 新的代码已经把sharedID放在MJRecorderShareIDTable哈希表中
		if recordID == "" {
			recordID, _ = redis.String(conn.Do("HGET", gconst.GameServerMJRecorderShareIDTable, sid))
			if recordID == "" {
				log.Println("handleLoadReplayRecord, no recordID found with sid:", sid)
				return
			}
		}
	}

	// "d" 二进制数据， "cid" 房间配置id
	values, err := redis.Values(conn.Do("HMGET", gconst.GameServerMJRecorderTablePrefix+recordID, "d", "cid"))
	if err != nil {
		log.Println("handleLoadReplayRecord, HMGET err:", err)
		return
	}

	// TODO: 日光，如果err == redis.NilErr，则需要去持久化数据库拉取数据
	bytesArray, err := redis.Bytes(values[0], nil)
	if err != nil {
		log.Println("handleLoadReplayRecord, load from record table err:", err)
		if err == redis.ErrNil {
			loadReplayRecordFromSQLServer(w, r, recordID)
		}
		return
	}

	cid, err := redis.String(values[1], nil)
	if err != nil {
		log.Println("handleLoadReplayRecord, load from record table err:", err)
		return
	}

	var msgHandRecorder = &lobby.MsgAccLoadReplayRecord{}
	// err = proto.Unmarshal(bytesArray, msgHandRecorder)
	// if err != nil {
	// 	log.Println("handleLoadReplayRecord, unmarshal err:", err)
	// 	return
	// }

	msgHandRecorder.ReplayRecordBytes = bytesArray
	roomConfigID := cid
	var roomConfig = lobby.RoomConfigs[roomConfigID]
	msgHandRecorder.RoomJSONConfig = &roomConfig

	//loadPlayerHeadIconURI(msgHandRecorder.GetPlayers())

	bytesArray, err = proto.Marshal(msgHandRecorder)
	if err != nil {
		log.Println("handleLoadReplayRecord, marshal err:", err)
		return
	}

	writeHTTPBodyWithGzip(w, r, bytesArray)
}

func writeHTTPBodyWithGzip(w http.ResponseWriter, r *http.Request, bytesArray []byte) {
	gzipSupport := false

	acceptContentEncodeStr := r.Header.Get("Accept-Encoding")
	if strings.Contains(acceptContentEncodeStr, "gzip") {
		log.Println("client support gzip")
		gzipSupport = true
	}

	if gzipSupport {
		var buf bytes.Buffer
		g := gzip.NewWriter(&buf)
		if _, err := g.Write(bytesArray); err != nil {
			log.Println("writeHTTPBodyWithGzip, write gzip err:", err)
			return
		}
		if err := g.Close(); err != nil {
			log.Println("writeHTTPBodyWithGzip, close gzip err:", err)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Add("Vary", "Accept-Encoding")
		bytesCompressed := buf.Bytes()
		log.Printf("COMPRESS, before:%d, after:%d\n", len(bytesArray), len(bytesCompressed))
		w.Write(bytesCompressed)
	} else {
		w.Write(bytesArray)
	}
}
