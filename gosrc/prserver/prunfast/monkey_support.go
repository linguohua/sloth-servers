package prunfast

// 用于支持monkey账户修改、查询服务器
import (
	"bytes"
	"compress/gzip"
	fmt "fmt"
	"gconst"
	"net/http"
	"pokerface"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/garyburd/redigo/redis"
)

type monkeySupportHandler func(w http.ResponseWriter, r *http.Request)

var (
	monkeySupportHandlers = make(map[string]monkeySupportHandler)
)

// monkeyAccountVerify 检查monkey用户接入合法
func monkeyAccountVerify(w http.ResponseWriter, r *http.Request) bool {
	var account = r.URL.Query().Get("account")
	var password = r.URL.Query().Get("password")
	// log.Printf("monkey access, account:%s, password:%s\n", account, password)
	conn := pool.Get()
	defer conn.Close()

	tableName := fmt.Sprintf("%s%d", gconst.GameServerMonkeyAccountTablePrefix, gconst.RoomType_DafengGZ)
	pass, e := redis.String(conn.Do("HGET", tableName, account))
	if e != nil || password != pass {
		return false
	}

	return true
}

func onExportRoomOperations(w http.ResponseWriter, r *http.Request) {
	var userID = r.URL.Query().Get("userID")
	var recordSID = ""
	if userID == "" {
		recordSID = r.URL.Query().Get("recordSID")
		if recordSID == "" {
			w.WriteHeader(404)
			w.Write([]byte("must supply userID or recordSID"))
			return
		}
	}

	if userID != "" {
		exportRoomOperationsByUserID(w, r, userID)
	} else {
		exportRoomOperationsByRecordSID(w, r, recordSID)
	}
}

func onExportRoomReplayRecordsSIDs(w http.ResponseWriter, r *http.Request) {
	recordSID := r.URL.Query().Get("recordSID")
	if recordSID == "" {
		w.WriteHeader(404)
		w.Write([]byte("must supply userID or recordSID"))
		return
	}

	conn := pool.Get()
	defer conn.Close()

	recordID, err := redis.String(conn.Do("HGET", gconst.GameServerMJRecorderTablePrefix+recordSID, "rid"))
	if err != nil && err != redis.ErrNil {
		log.Println("can't found rid for sid:", recordSID)
		w.Write([]byte("no mj record found for record:" + recordSID))
		return
	}

	// 新的代码已经把sharedID放在MJRecorderShareIDTable哈希表中
	if recordID == "" {
		recordID, err = redis.String(conn.Do("HGET", gconst.GameServerMJRecorderShareIDTable, recordSID))
		if err != nil {
			log.Println("can't found rid for sid:", recordSID)
			w.Write([]byte("no mj record found for record:" + recordSID))
			return
		}
	}

	w.WriteHeader(404)
	w.Write([]byte("onExportRoomReplayRecordsSIDs has removed"))
}

func onExportRoomReplayRecord(w http.ResponseWriter, r *http.Request) {
	recordID := r.URL.Query().Get("recordID")

	conn := pool.Get()
	defer conn.Close()
	buf := loadMJRecord(conn, recordID)

	if buf == nil {
		w.WriteHeader(404)
		w.Write([]byte("no mj record found for record:" + recordID))
		return
	}

	writeHTTPBodyWithGzip(w, r, buf)
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

func exportRoomOperationsByRecordSID(w http.ResponseWriter, r *http.Request, recordSID string) {
	conn := pool.Get()
	defer conn.Close()

	recordID, err := redis.String(conn.Do("HGET", gconst.GameServerMJRecorderTablePrefix+recordSID, "rid"))
	if err != nil && err != redis.ErrNil {
		log.Println("can't found rid for sid:", recordSID)
		w.Write([]byte("no mj record found for record:" + recordSID))
		return
	}

	// 新的代码已经把sharedID放在MJRecorderShareIDTable哈希表中
	if recordID == "" {
		recordID, err = redis.String(conn.Do("HGET", gconst.GameServerMJRecorderShareIDTable, recordSID))
		if err != nil {
			log.Println("can't found rid for sid:", recordSID)
			w.Write([]byte("no mj record found for record:" + recordSID))
			return
		}
	}

	buf := loadMJRecord(conn, recordID)

	if buf == nil {
		w.WriteHeader(404)
		w.Write([]byte("no mj record found for record:" + recordSID))
		return
	}

	w.Write(buf)
}

func exportRoomOperationsByUserID(w http.ResponseWriter, r *http.Request, userID string) {
	var buf []byte
	user, ok := usersMap[userID]
	if ok {
		// 先尝试加载其所在的房间的操作列表
		var room = user.user.getRoom()

		if room != nil {
			log.Println("user in server, room ID:", room.ID)
			switch room.state.(type) {
			case *SPlaying:
				s := room.state.(*SPlaying)
				ctx := s.lctx
				if ctx != nil {
					log.Println("found active ctx in room:", room.ID)
					buf = ctx.toByteArray()
				}
				break
			default:
				break
			}
		}
	}

	if buf == nil {
		log.Println("can't found active record, try to load from redis")
		buf = loadMJLastRecordForUser(userID)
	}

	if buf == nil {
		w.WriteHeader(404)
		w.Write([]byte("no mj record found for user:" + userID))
		return
	}

	w.Write(buf)
}

func onExportRoomCfg(w http.ResponseWriter, r *http.Request) {
	var roomConfigID = r.URL.Query().Get("roomConfigID")
	if roomConfigID == "" {
		w.WriteHeader(404)
		w.Write([]byte("must supply roomConfigID"))
		return
	}

	buf := loadRoomConfigFromRedis(roomConfigID)
	if buf == nil {
		w.WriteHeader(404)
		w.Write([]byte("failed to load config for:" + roomConfigID))
		return
	}

	w.Write(buf)
}

func onRoomKickAll(w http.ResponseWriter, r *http.Request) {
	var roomNumber = r.URL.Query().Get("roomNumber")
	if roomNumber == "" {
		w.WriteHeader(404)
		w.Write([]byte("must supply roomNumber"))
		return
	}

	room := roomMgr.getRoomByNumber(roomNumber)
	if room == nil {
		w.WriteHeader(404)
		w.Write([]byte("no room for :" + roomNumber))
		return
	}

	room.kickAll()

	w.Write([]byte("OK, kick out all in room:" + roomNumber))
}

func onRoomReset(w http.ResponseWriter, r *http.Request) {
	var roomNumber = r.URL.Query().Get("roomNumber")
	if roomNumber == "" {
		w.WriteHeader(404)
		w.Write([]byte("must supply roomNumber"))
		return
	}

	room := roomMgr.getRoomByNumber(roomNumber)
	if room == nil {
		w.WriteHeader(404)
		w.Write([]byte("no room for :" + roomNumber))
		return
	}

	room.reset()

	w.Write([]byte("OK, reset room:" + roomNumber))
}

func onRoomDisband(w http.ResponseWriter, r *http.Request) {
	log.Println("monkey try to disband room...")
	var roomNumber = r.URL.Query().Get("roomNumber")
	if roomNumber == "" {
		w.WriteHeader(404)
		w.Write([]byte("must supply roomNumber"))
		return
	}

	room := roomMgr.getRoomByNumber(roomNumber)
	if room == nil {
		w.WriteHeader(404)
		w.Write([]byte("no room for :" + roomNumber))
		return
	}

	roomMgr.forceDisbandRoom(room, pokerface.RoomDeleteReason_DisbandBySystem)

	w.Write([]byte("OK, disband room:" + roomNumber))
}

func onExportUserLastRecord(w http.ResponseWriter, r *http.Request) {
	var userID = r.URL.Query().Get("userID")
	buf := loadMJLastRecordForUser(userID)
	if buf == nil {
		w.WriteHeader(404)
		w.Write([]byte("failed to load record for userID:" + userID))
		return
	}

	w.Write(buf)
}

func onQueryRoomCount(w http.ResponseWriter, r *http.Request) {
	roomCount := len(roomMgr.rooms)
	roomIdle := 0
	roomWaiting := 0
	roomPlaying := 0

	for _, r := range roomMgr.rooms {
		stateConst := r.state.getStateConst()
		switch stateConst {
		case pokerface.RoomState_SRoomIdle:
			roomIdle++
			break
		case pokerface.RoomState_SRoomWaiting:
			roomWaiting++
			break
		case pokerface.RoomState_SRoomPlaying:
			roomPlaying++
			break
		}
	}

	w.Write([]byte(fmt.Sprintf("room count:%d, idle:%d, wait:%d, play:%d", roomCount, roomIdle, roomWaiting, roomPlaying)))
}

func onQueryUserCount(w http.ResponseWriter, r *http.Request) {
	userCount := len(usersMap)
	w.Write([]byte(strconv.Itoa(userCount)))
}

func onUploadCfgs(w http.ResponseWriter, r *http.Request) {
	body := make([]byte, r.ContentLength)
	n, _ := r.Body.Read(body)
	if n != int(r.ContentLength) {
		msg := "upload cfg error, read message length not match content length"
		log.Println(msg)
		return
	}

	monkeyMgr.doUploadCfgs(w, r, string(body))
}

func onCreateMonkeyRoom(w http.ResponseWriter, r *http.Request) {
	monkeyMgr.createMonkeyRoom(w, r)
}

func onDestroyMonkeyRoom(w http.ResponseWriter, r *http.Request) {
	monkeyMgr.destroyMonkeyRoom(w, r)
}

func onAttachDealCfg2Room(w http.ResponseWriter, r *http.Request) {
	body := make([]byte, r.ContentLength)
	n, _ := r.Body.Read(body)
	if n != int(r.ContentLength) {
		log.Println("attach deal cfg error, read message length not match content length")
		return
	}

	monkeyMgr.attachDealCfg2Room(w, r, string(body))
}

func onAttachRoomCfg2Room(w http.ResponseWriter, r *http.Request) {
	body := make([]byte, r.ContentLength)
	n, _ := r.Body.Read(body)
	if n != int(r.ContentLength) {
		log.Println("attach room cfg error, read message length not match content length")
		return
	}

	monkeyMgr.attachRoomCfg2Room(w, r, string(body))
}

func onKickUser(w http.ResponseWriter, r *http.Request) {
	var userID = r.URL.Query().Get("userID")
	userItem, ok := usersMap[userID]
	if ok {
		userItem.user.closeWebsocket()
		w.Write([]byte("kickout ok, ID:" + userID))
	} else {
		w.Write([]byte("kickout falied, not found ID:" + userID))
	}
}

func onQueryRoomExceptionCount(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("%d", roomExceptionCount)))
}

func onClearRoomExceptionCount(w http.ResponseWriter, r *http.Request) {
	roomExceptionCount = 0
}

func onULimitRound(w http.ResponseWriter, r *http.Request) {
	log.Println("monkey try to ulimit room...")
	var roomNumber = r.URL.Query().Get("roomNumber")
	if roomNumber == "" {
		w.WriteHeader(404)
		w.Write([]byte("must supply roomNumber"))
		return
	}

	room := roomMgr.getRoomByNumber(roomNumber)
	if room == nil {
		w.WriteHeader(404)
		w.Write([]byte("no room for :" + roomNumber))
		return
	}

	room.isUlimitRound = true

	w.Write([]byte("OK:" + roomNumber))
}

func init() {
	monkeySupportHandlers["/uploadCfgs"] = onUploadCfgs
	monkeySupportHandlers["/createMonkeyRoom"] = onCreateMonkeyRoom
	monkeySupportHandlers["/destroyMonkeyRoom"] = onDestroyMonkeyRoom
	monkeySupportHandlers["/attachDealCfg"] = onAttachDealCfg2Room
	monkeySupportHandlers["/attachRoomCfg"] = onAttachRoomCfg2Room
	monkeySupportHandlers["/kickUser"] = onKickUser
	monkeySupportHandlers["/exportRoomOps"] = onExportRoomOperations
	monkeySupportHandlers["/exportRoomCfg"] = onExportRoomCfg
	monkeySupportHandlers["/exportUserLastRecord"] = onExportUserLastRecord
	monkeySupportHandlers["/kickAll"] = onRoomKickAll
	monkeySupportHandlers["/resetRoom"] = onRoomReset
	monkeySupportHandlers["/disbandRoom"] = onRoomDisband
	monkeySupportHandlers["/roomCount"] = onQueryRoomCount
	monkeySupportHandlers["/userCount"] = onQueryUserCount
	monkeySupportHandlers["/roomException"] = onQueryRoomExceptionCount
	monkeySupportHandlers["/clearRoomException"] = onClearRoomExceptionCount
	monkeySupportHandlers["/ulimitRound"] = onULimitRound
	monkeySupportHandlers["/exportRoomSIDs"] = onExportRoomReplayRecordsSIDs
	monkeySupportHandlers["/exportRR"] = onExportRoomReplayRecord
}
