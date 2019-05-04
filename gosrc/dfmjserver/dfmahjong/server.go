package dfmahjong

import (
	"fmt"
	"gconst"
	"gscfg"
	"mahjong"
	"net/http"
	"path"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"

	"github.com/julienschmidt/httprouter"
)

const (
	wsReadLimit       = 1024 // 每个websocket的接收数据包长度限制
	wsReadBufferSize  = 2048 // 每个websocket的接收缓冲限制
	wsWriteBufferSize = 4096 // 每个websocket的发送缓冲限制

	myRoomType = gconst.RoomType_DafengMJ
)

var (
	upgrader = websocket.Upgrader{ReadBufferSize: wsReadBufferSize, WriteBufferSize: wsWriteBufferSize}
	// 根router，只有http server看到
	rootRouter = httprouter.New()

	roomMgr            = &RoomMgr{}                    // 房间管理
	monkeyMgr          = &MonkeyMgr{}                  // monkey
	usersMap           = make(map[string]*UserMapItem) // 所有接入玩家
	roomExceptionCount int                             // 房间异常计数
)

// UserMapItem 表格item，用于usersMap
type UserMapItem struct {
	user             IUser
	wg               sync.WaitGroup
	lastReceivedTime time.Time
	lastPingTime     time.Time
}

// 在线玩家数量加1
func incrOnlinePlayerNum() {
	conn := pool.Get()
	defer conn.Close()

	var key = fmt.Sprintf("%s%d", gconst.GameServerOnlineUserNumPrefix, myRoomType)
	conn.Do("HINCRBY", key, gscfg.ServerID, 1)
}

// 在线玩家数量减1
func decrOnlinePlayerNum() {
	conn := pool.Get()
	defer conn.Close()

	var key = fmt.Sprintf("%s%d", gconst.GameServerOnlineUserNumPrefix, myRoomType)
	conn.Do("HINCRBY", key, gscfg.ServerID, -1)
}

// waitWebsocketMessage 接收和分发游戏玩家的websocket消息
func waitWebsocketMessage(ws *websocket.Conn, userMapItem *UserMapItem, r *http.Request) {
	user := userMapItem.user

	ws.SetPongHandler(func(msg string) error {
		//log.Printf("websocket recv ping msg:%s, size:%d\n", msg, len(msg))
		userMapItem.lastReceivedTime = time.Now()
		return nil
	})

	ws.SetPingHandler(func(msg string) error {
		//log.Printf("websocket recv ping msg size:%d\n", len(msg))
		userMapItem.lastReceivedTime = time.Now()
		user.sendPong(msg)
		return nil
	})

	// 确保无论出任何情况都会调用onWebsocketClosed，以便房间可以做对玩家做离线处理
	defer func() {
		user.onWebsocketClosed(ws)
	}()

	log.Printf("wait ws msg, userId: %s, peer: %s", user.userID(), r.RemoteAddr)
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println(" websocket receive error:", err)
			ws.Close()
			break
		}

		userMapItem.lastReceivedTime = time.Now()

		// 只处理BinaryMessage，其他的忽略
		if message != nil && len(message) > 0 && mt == websocket.BinaryMessage {
			user.onWebsocketMessage(ws, message)
		}

		//log.Printf("receive from user %d message:%s", user.userID(), message)
	}
	log.Printf("ws closed, userId %s, peer:%s", user.userID(), r.RemoteAddr)
}

// tryAcceptGameUser 游戏玩家接入
func tryAcceptGameUser(userID string, roomIDString string, ws *websocket.Conn, r *http.Request) {

	// 查找房间，如果房间不存在则非法
	var room = roomMgr.getRoom(roomIDString)
	if room == nil {
		log.Printf("invalid room ID:%s, Peer:%s\n", roomIDString, r.RemoteAddr)
		// 给客户端发送错误
		sendEnterRoomError(ws, userID, mahjong.EnterRoomStatus_RoomNotExist)
		return
	}

	oldUserItem, ok := usersMap[userID]
	if ok && oldUserItem.user.getRoom() != room {
		// 给客户端发送错误
		sendEnterRoomError(ws, userID, mahjong.EnterRoomStatus_InAnotherRoom)
		return
	}

	if ok {
		// 等待老websocket的关闭
		oldUserItem.user.closeWebsocket()
		oldUserItem.wg.Wait()
		log.Println("wait old ws ok:", userID)
	}

	user := room.userTryEnter(ws, userID)

	if user != nil {
		userMapItem := &UserMapItem{user: user,
			lastReceivedTime: time.Now(), lastPingTime: time.Now()}
		userMapItem.wg.Add(1)
		usersMap[userID] = userMapItem

		defer func() {
			delete(usersMap, userID)
			userMapItem.wg.Done()
			decrOnlinePlayerNum()
		}()

		incrOnlinePlayerNum()

		waitWebsocketMessage(ws, userMapItem, r)
	}
}

// acceptWebsocket 把http请求转换为websocket
func acceptWebsocket(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var requestPath = r.URL.Path
	requestPath = path.Base(requestPath)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// 接收限制
	ws.SetReadLimit(wsReadLimit)

	// 确保 websocket 关闭
	defer ws.Close()

	log.Println("accept websocket:", r.URL)
	switch requestPath {
	case "play":
		var token = r.URL.Query().Get("tk")
		userID, ok := parseTK(token)
		if !ok {
			log.Printf("invalid token, Peer: %s", r.RemoteAddr)
			return
		}

		if gscfg.RequiredAppModuleVer > 0 {
			appModuleVer, err := strconv.Atoi(r.URL.Query().Get("amv"))
			if err != nil || appModuleVer < gscfg.RequiredAppModuleVer {
				log.Printf("app module too old, ID:%s, Peer:%s\n", userID, r.RemoteAddr)
				// 给客户端发送错误
				sendEnterRoomError(ws, userID, mahjong.EnterRoomStatus_AppModuleNeedUpgrade)
				return
			}
		}

		// 房间ID要合法
		var roomIDString = r.URL.Query().Get("roomID")
		tryAcceptGameUser(userID, roomIDString, ws, r)
		break
	case "monkey":
		var userID = r.URL.Query().Get("userID")
		// 房间ID要合法
		var roomIDString = r.URL.Query().Get("roomID")
		if roomIDString == "" {
			// 此时monkey传上来的可能是号码
			var roomNumber = r.URL.Query().Get("roomNumber")
			if roomNumber == "" {
				log.Println("monkey has no roomID and roomNumber")
				return
			}

			room := roomMgr.getRoomByNumber(roomNumber)
			if room != nil {
				roomIDString = room.ID
			} else {
				log.Println("no room found for roomNumber:", roomNumber)
				return
			}
		}
		tryAcceptGameUser(userID, roomIDString, ws, r)
		break
	}
}

func echoVersion(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte(fmt.Sprintf("version:%d", versionCode)))
}

func monkeyHTTPHandle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var spName = ps.ByName("sp")

	log.Println("monkey support handler call:", spName)
	if monkeyAccountVerify(w, r) {
		h, ok := monkeySupportHandlers[spName]
		if ok {
			h(w, r)
		} else {
			log.Println("no monkey support handler found:", spName)
		}
	} else {
		var msg = "no authorization for call monkey handler:" + spName
		log.Println(msg)
		w.Write([]byte(msg))
	}
}

// CreateHTTPServer 启动服务器
func CreateHTTPServer() {
	roomMgr.startup()
	monkeyMgr.start()
	startAliveKeeper()

	// 所有模块看到的mainRouter
	// 外部访问需要形如/game/uuid/play
	rootRouter.Handle("GET", "/game/:uuid/ws/:wtype", acceptWebsocket)
	rootRouter.Handle("GET", "/game/:uuid/version", echoVersion)

	// POST和GET都要订阅
	rootRouter.Handle("GET", "/game/:uuid/support/*sp", monkeyHTTPHandle)
	rootRouter.Handle("POST", "/game/:uuid/support/*sp", monkeyHTTPHandle)

	go acceptHTTPRequest()
}

// acceptHTTPRequest 监听和接受HTTP
func acceptHTTPRequest() {
	portStr := fmt.Sprintf(":%d", gscfg.ServerPort)
	s := &http.Server{
		Addr:    portStr,
		Handler: rootRouter,
		// ReadTimeout:    10 * time.Second,
		//WriteTimeout:   120 * time.Second,
		MaxHeaderBytes: 1 << 8,
	}

	log.Printf("Http server listen at:%d\n", gscfg.ServerPort)

	err := s.ListenAndServe()
	if err != nil {
		log.Fatalf("Http server ListenAndServe %d failed:%s\n", gscfg.ServerPort, err)
	}

}
