package lobby

import (
	"fmt"
	"gconst"
	"lobbyserver/config"
	"net/http"
	"path"
	"time"

	log "github.com/sirupsen/logrus"

	"lobbyserver/pricecfg"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"

	"github.com/gorilla/mux"
)

const (
	wsReadLimit       = 64 * 1024 // 64K
	wsReadBufferSize  = 4 * 1024
	wsWriteBufferSize = 4 * 1024
	authServer        = "http://wjr.u8.login.qianz.com:8080/user/verifyUser?userID=%s&token=%s&sign=%s"
	appKey            = "253c8d16bc73d85ac7066dcae0e478fe"
)

type accUserIDHTTPHandler func(w http.ResponseWriter, r *http.Request, userID string)
type accRawHTTPHandler func(w http.ResponseWriter, r *http.Request)

var (
	upgrader = websocket.Upgrader{ReadBufferSize: wsReadBufferSize, WriteBufferSize: wsWriteBufferSize}
	// 根router，只有http server看到
	rootRouter = mux.NewRouter()

	// MainRouter main-router
	MainRouter *mux.Router

	userMgr              *UserMgr
	accSysExceptionCount int // 异常计数
	// chost                = &clubHost{}

	accRawHTTPHandlers    = make(map[string]accRawHTTPHandler)
	accUserIDHTTPHandlers = make(map[string]accUserIDHTTPHandler)
)

// UserCount user count
func UserCount() int {
	return len(userMgr.users)
}

func waitWebsocketMessage(ws *websocket.Conn, user *User, r *http.Request) {
	log.Printf("wait ws msg, userId: %s, peer: %s", user.userID(), r.RemoteAddr)

	ws.SetPongHandler(func(msg string) error {
		user.lastReceivedTime = time.Now()
		return nil
	})

	ws.SetPingHandler(func(msg string) error {
		user.lastReceivedTime = time.Now()
		user.sendPong(msg)
		return nil
	})

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println(" websocket receive error:", err)
			ws.Close()
			user.onWebsocketClosed(ws)
			break
		}

		user.lastReceivedTime = time.Now()

		if message != nil && len(message) > 0 {
			user.onWebsocketMessage(ws, message)
		}

		//log.Printf("receive from user %d message:%s", user.userID(), message)
	}
	log.Printf("ws closed, userId %s, peer:%s", user.userID(), r.RemoteAddr)
}

func loadCharm(userID string) int32 {
	conn := pool.Get()
	defer conn.Close()

	charm, _ := redis.Int(conn.Do("HGET", gconst.AsUserTablePrefix+userID, "charm"))
	return int32(charm)
}

func replyLoginError(ws *websocket.Conn, errCode int32) {
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

func loginReply(ws *websocket.Conn, userID string) {
	var msgLoginReply = &MsgLoginReply{}
	var errCode = int32(LoginState_Success)
	msgLoginReply.Result = &errCode

	var tk = genTK(userID)
	msgLoginReply.Token = &tk

	var lastRoomInfo = loadLastRoomInfo(userID)
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

func tryAcceptUser(ws *websocket.Conn, r *http.Request) {
	userID, ok := verifyToken(r)
	if !ok {
		log.Println("verifyUser failed")
		replyLoginError(ws, int32(LoginState_ParseTokenError))
		return
	}

	var user = newUser(ws, userID)

	oldUser := userMgr.getUserByID(user.uID)
	if oldUser != nil {
		oldUser.detach()
		oldUser.wg.Wait()
	}

	userMgr.addUser(user)

	loginReply(ws, userID)

	defer func() {
		userMgr.removeUser(user)
		user.wg.Done()
	}()

	waitWebsocketMessage(ws, user, r)
}

func acceptWebsocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// 接收限制
	ws.SetReadLimit(wsReadLimit)

	defer ws.Close()

	log.Println("accept websocket:", r.URL)
	tryAcceptUser(ws, r)
}

func initAccUserIDHTTPHandlers() {
	accUserIDHTTPHandlers["/lrproom"] = handleLoadReplayRooms
	accUserIDHTTPHandlers["/lrprecord"] = handleLoadReplayRecord
	accUserIDHTTPHandlers["/createRoom"] = handlerCreateRoom
	accUserIDHTTPHandlers["/requestRoomInfo"] = handlerRequestRoomInfo
	accUserIDHTTPHandlers["/loadUserScoreInfo"] = handleLoadUserScoreInfo
	accUserIDHTTPHandlers["/loadUserHeadIconURI"] = handleLoadUserHeadIconURI
	accUserIDHTTPHandlers["/uploadLogFile"] = handleUploadLogFile
	accUserIDHTTPHandlers["/updateUserLocation"] = handleUpdateUserLocation
	accUserIDHTTPHandlers["/loadPrices"] = handleLoadPrices
	accUserIDHTTPHandlers["/clubCreateOrder"] = OnCreateClubOrderForWX       // 创建俱乐部基金订单
	accUserIDHTTPHandlers["/addAgentInfo"] = OnAddAgentInfo                  // 添加代理信息到后台
	accUserIDHTTPHandlers["/getClubShopConfig"] = OnGetClubShopConfig        // 获取俱乐部商城配置
	accUserIDHTTPHandlers["/loadLastRoomInfo"] = handlerLoadLastRoomInfo     // 拉取用户最后所在的房间
	accUserIDHTTPHandlers["/lgrouprproom"] = handleLoadGroupReplayRooms      // 拉取茶馆战绩
	accUserIDHTTPHandlers["/lgroupbw"] = handleLoadGroupBigWinner            // 拉取大赢家
	accUserIDHTTPHandlers["/queryUserRoomInfo"] = handleQueryUserRoomInfo    // 查询房间状态
	accUserIDHTTPHandlers["/deleteRoom"] = handlerDeleteRoom                 // 删除房间
	accUserIDHTTPHandlers["/deleteRoomForGroup"] = handlerDeleteRoomForGroup // 删除房间
	accUserIDHTTPHandlers["/confirmBigWinner"] = handleConfirmGroupBigWinner // 茶馆确认是否是大赢家
}

func initAccRawHTTPHandlers() {
	// accRawHTTPHandlers["/acc"] = acceptWebsocket
	accRawHTTPHandlers["/updateDiamond"] = handleUpdateDiamond
	accRawHTTPHandlers["/orderresult"] = OnOrderResultForWX // 这个是url是易洋web服务器定的，与大厅商店的url一样
}

func echoVersion(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("version:%d", versionCode)))
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
		h2, ok := accUserIDHTTPHandlers[requestPath]
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
	// TODO: just for test, please remove later
	log.Println("For cub test:" + genTK("10024063"))

	startRedisClient()

	loadRoomTypeFromRedis()

	initGamePropCfgs()

	// subscriberUserConnectEvent()
	// subscriberUserDisConnectEvent()

	loadAllRoomConfigFromRedis()

	pricecfg.LoadAllPriceCfg(pool)

	loadSensitiveWordDictionary(config.SensitiveWordFilePath)

	userMgr = newUserMgr()
	startAliveKeeper()

	initAccUserIDHTTPHandlers()
	initAccRawHTTPHandlers()

	// 所有模块看到的mainRouter
	// 外部访问需要形如/prunfast/uuid/pok
	var mainRouter = rootRouter.PathPrefix("/lobby/{uuid}/").Subrouter()
	mainRouter.HandleFunc("/ws/{wtype}", acceptWebsocket)
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
