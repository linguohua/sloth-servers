package lobby

import (
	"fmt"
	"gconst"
	"lobbyserver/config"
	"math/rand"
	"net/http"
	"path"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"

	log "github.com/sirupsen/logrus"

	"lobbyserver/pricecfg"

	"github.com/garyburd/redigo/redis"

	"github.com/gorilla/mux"
)

const (
	authServer = "http://wjr.u8.login.qianz.com:8080/user/verifyUser?userID=%s&token=%s&sign=%s"
	appKey     = "253c8d16bc73d85ac7066dcae0e478fe"
)

type accUserIDHTTPHandler func(w http.ResponseWriter, r *http.Request, userID string)
type accRawHTTPHandler func(w http.ResponseWriter, r *http.Request)

var (
	// 根router，只有http server看到
	rootRouter = mux.NewRouter()

	// MainRouter main-router
	MainRouter *mux.Router

	accSysExceptionCount int // 异常计数
	// chost                = &clubHost{}

	accRawHTTPHandlers = make(map[string]accRawHTTPHandler)

	// AccUserIDHTTPHandlers trust handlers
	AccUserIDHTTPHandlers = make(map[string]accUserIDHTTPHandler)

	// SessionMgr mgr
	SessionMgr ISessionMgr
	// RoomUtil room helper functions
	RoomUtil IRoomUtil

	// RandGenerator rand generator
	RandGenerator *rand.Rand
)

func loadCharm(userID string) int32 {
	conn := pool.Get()
	defer conn.Close()

	charm, _ := redis.Int(conn.Do("HGET", gconst.AsUserTablePrefix+userID, "charm"))
	return int32(charm)
}

// func initAccUserIDHTTPHandlers() {
// 	accUserIDHTTPHandlers["/createRoom"] = handlerCreateRoom
// 	accUserIDHTTPHandlers["/requestRoomInfo"] = handlerRequestRoomInfo
// 	accUserIDHTTPHandlers["/loadUserScoreInfo"] = handleLoadUserScoreInfo
// 	accUserIDHTTPHandlers["/loadUserHeadIconURI"] = handleLoadUserHeadIconURI
// 	accUserIDHTTPHandlers["/uploadLogFile"] = handleUploadLogFile
// 	accUserIDHTTPHandlers["/updateUserLocation"] = handleUpdateUserLocation
// 	accUserIDHTTPHandlers["/loadPrices"] = handleLoadPrices
// 	accUserIDHTTPHandlers["/clubCreateOrder"] = OnCreateClubOrderForWX       // 创建俱乐部基金订单
// 	accUserIDHTTPHandlers["/addAgentInfo"] = OnAddAgentInfo                  // 添加代理信息到后台
// 	accUserIDHTTPHandlers["/getClubShopConfig"] = OnGetClubShopConfig        // 获取俱乐部商城配置
// 	accUserIDHTTPHandlers["/loadLastRoomInfo"] = handlerLoadLastRoomInfo     // 拉取用户最后所在的房间
// 	accUserIDHTTPHandlers["/lgrouprproom"] = handleLoadGroupReplayRooms      // 拉取茶馆战绩
// 	accUserIDHTTPHandlers["/lgroupbw"] = handleLoadGroupBigWinner            // 拉取大赢家
// 	accUserIDHTTPHandlers["/queryUserRoomInfo"] = handleQueryUserRoomInfo    // 查询房间状态
// 	accUserIDHTTPHandlers["/deleteRoom"] = handlerDeleteRoom                 // 删除房间
// 	accUserIDHTTPHandlers["/deleteRoomForGroup"] = handlerDeleteRoomForGroup // 删除房间
// 	accUserIDHTTPHandlers["/confirmBigWinner"] = handleConfirmGroupBigWinner // 茶馆确认是否是大赢家
// }

func initAccRawHTTPHandlers() {
	// accRawHTTPHandlers["/acc"] = acceptWebsocket
	// accRawHTTPHandlers["/updateDiamond"] = handleUpdateDiamond
	// accRawHTTPHandlers["/orderresult"] = OnOrderResultForWX // 这个是url是易洋web服务器定的，与大厅商店的url一样
}

func echoVersion(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("version:%d", versionCode)))
}

// ReplyLoginError login error
func ReplyLoginError(ws *websocket.Conn, errCode int32) {
	var msgLoginReply = &MsgLoginReply{}
	msgLoginReply.Result = proto.Int32(errCode)
	var errorString = LoginString[errCode]
	msgLoginReply.RetMsg = &errorString

	msgLoginReplyBuf, err := proto.Marshal(msgLoginReply)
	if err != nil {
		log.Println("replyLoginError, Marshal error:", err)
		return
	}

	var messageCode = int32(MessageCode_OPLoginReply)
	accessoryMessage := &AccessoryMessage{}
	accessoryMessage.Ops = &messageCode
	accessoryMessage.Data = msgLoginReplyBuf
	accessoryMessageBuf, err := proto.Marshal(accessoryMessage)
	if err != nil {
		log.Println("replyLoginError, Marshal error:", err)
		return
	}
	// log.Println(msgLoginReply)
	ws.WriteMessage(websocket.BinaryMessage, accessoryMessageBuf)

}

// LoginReply login reply
func LoginReply(ws *websocket.Conn, userID string) {
	var msgLoginReply = &MsgLoginReply{}
	var errCode = int32(LoginState_Success)
	msgLoginReply.Result = &errCode

	var tk = genTK(userID)
	msgLoginReply.Token = &tk

	var lastRoomInfo = RoomUtil.LoadLastRoomInfo(userID)
	if lastRoomInfo != nil {
		msgLoginReply.LastRoomInfo = lastRoomInfo
	}

	msgLoginReplyBuf, err := proto.Marshal(msgLoginReply)
	if err != nil {
		log.Println(err)
		return
	}

	var messageCode = int32(MessageCode_OPLoginReply)
	accessoryMessage := &AccessoryMessage{}
	accessoryMessage.Ops = &messageCode
	accessoryMessage.Data = msgLoginReplyBuf

	accessoryMessageBuf, err := proto.Marshal(accessoryMessage)
	if err != nil {
		log.Println(err)
		return
	}
	// log.Println(msgLoginReply)
	ws.WriteMessage(websocket.BinaryMessage, accessoryMessageBuf)
}

// func (mux *myHTTPServerMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

// 	var requestPath = r.URL.Path
// 	log.Println(requestPath)
// 	// log.Println("requestPath:", requestPath)
// 	requestPath = "/" + path.Base(requestPath)
// 	if requestPath == "/version" {
// 		w.Write([]byte(fmt.Sprintf("version:%d", versionCode)))
// 		return
// 	}

// 	h, ok := accRawHTTPHandlers[requestPath]
// 	if ok {
// 		h(w, r)
// 		return
// 	}

// 	h2, ok := accUserIDHTTPHandlers[requestPath]
// 	if ok {
// 		userID, rt := verifyTokenByQuery(r)
// 		if rt {
// 			h2(w, r, userID)
// 		} else {
// 			w.WriteHeader(404)
// 			w.Write([]byte("oh, no valid token"))
// 		}

// 		return
// 	}

// 	h3, ok := accSupportHandlers[requestPath]
// 	if ok {
// 		if accSupportVerify(w, r) {
// 			h3(w, r)
// 		} else {
// 			w.WriteHeader(404)
// 			w.Write([]byte("oh, no auth"))
// 		}
// 	} else {
// 		w.WriteHeader(404)
// 		w.Write([]byte("oh, no acc support handler"))
// 	}
// }

func authorizedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := "/" + path.Base(r.URL.Path)
		h2, ok := AccUserIDHTTPHandlers[requestPath]
		if ok {
			userID, rt := verifyTokenByQuery(r)
			if rt {
				h2(w, r, userID)
			} else {
				w.WriteHeader(404)
				w.Write([]byte("oh, no valid token for handler:" + r.URL.Path))
			}
		} else {
			w.WriteHeader(404)
			w.Write([]byte("oh, can't found handler for:" + r.URL.Path))
		}
	})
}

// CreateHTTPServer 启动服务器
func CreateHTTPServer() {
	RandGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))
	// TODO: just for test, please remove later
	log.Println("For cub test:" + genTK("10024063"))

	startRedisClient()

	loadRoomTypeFromRedis()

	//initGamePropCfgs()

	// subscriberUserConnectEvent()
	// subscriberUserDisConnectEvent()

	//loadAllRoomConfigFromRedis()

	pricecfg.LoadAllPriceCfg(pool)

	//initAccUserIDHTTPHandlers()
	initAccRawHTTPHandlers()

	// 所有模块看到的mainRouter
	// 外部访问需要形如/prunfast/uuid/pok
	var mainRouter = rootRouter.PathPrefix("/lobby/{uuid}/").Subrouter()
	//mainRouter.HandleFunc("/ws/{wtype}", acceptWebsocket)
	mainRouter.HandleFunc("/version", echoVersion)
	mainRouter.PathPrefix("/untrust").Handler(authorizedHandler())

	var uhRouter = mainRouter.PathPrefix("/trust").Subrouter()
	for k, v := range accRawHTTPHandlers {
		uhRouter.HandleFunc(k, v)
	}

	MainRouter = mainRouter

	// 挂载俱乐部
	// hostClub()

	// 恢复俱乐部房间
	//chost.clubRoomsListener.RestoreClubRoomsFromRedis()

	go acceptHTTPRequest()
}

func acceptHTTPRequest() {
	portStr := fmt.Sprintf(":%d", config.AccessoryServerPort)
	s := &http.Server{
		Addr:    portStr,
		Handler: rootRouter,
		// ReadTimeout:    10 * time.Second,
		//WriteTimeout:   120 * time.Second,
		MaxHeaderBytes: 1 << 8,
	}

	log.Printf("Http server listen at:%d\n", config.AccessoryServerPort)

	err := s.ListenAndServe()
	if err != nil {
		log.Fatalf("Http server ListenAndServe %d failed:%v", config.AccessoryServerPort, err)
	}

}

func loadRoomTypeFromRedis() {
	conn := pool.Get()
	defer conn.Close()

	roomTypes, err := redis.Ints(conn.Do("SMEMBERS", gconst.RoomTypeSet))
	if err != nil {
		log.Println("loadRoomTypeFromRedis, err:", err)
		return
	}

	conn.Send("MULTI")
	for _, roomType := range roomTypes {
		var key = fmt.Sprintf("%s%d", gconst.RoomTypeKey, roomType)
		conn.Send("HGET", key, "gameID")
	}
	gameIDs, err := redis.Ints(conn.Do("EXEC"))

	for i, roomType := range roomTypes {
		gameID := gameIDs[i]
		if gameID != 0 {
			var key = fmt.Sprintf("%d", roomType)
			config.SubGameIDs[key] = gameID
		}
	}
}

// UpdateRoomGameID web那边更新gameID
func UpdateRoomGameID(roomType int, gameID int) {
	var key = fmt.Sprintf("%d", roomType)
	config.SubGameIDs[key] = gameID

	log.Println("UpdateRoomGameID:", config.SubGameIDs)
}
