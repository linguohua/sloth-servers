package sessions

import (
	"lobbyserver/lobby"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	log "github.com/sirupsen/logrus"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

const (
	wsReadLimit       = 64 * 1024 // 64K
	wsReadBufferSize  = 4 * 1024
	wsWriteBufferSize = 4 * 1024
)

var (
	upgrader = websocket.Upgrader{ReadBufferSize: wsReadBufferSize, WriteBufferSize: wsWriteBufferSize}

	userMgr *UserMgr
)

func replyConnectMsg(ws *websocket.Conn, errCode int32) {
	reply := &lobby.MsgWebsocketConnectReply{}
	reply.Result = &errCode

	ops := int32(lobby.MessageCode_OPConnectReply)

	lobbyMessage := &lobby.LobbyMessage{}
	lobbyMessage.Ops = &ops

	bytes, err := proto.Marshal(reply)
	if err != nil {
		log.Panic("replyConnectMsg, marshal msg failed")
		return
	}
	lobbyMessage.Data = bytes

	bytes, err = proto.Marshal(lobbyMessage)
	if err != nil {
		log.Panic("sendMsg, marshal msg failed")
		return
	}

	ws.WriteMessage(websocket.BinaryMessage, bytes)
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

func tryAcceptUser(ws *websocket.Conn, r *http.Request) {
	userID, ok := lobby.VerifyToken(r)
	if !ok {
		log.Println("verifyUser failed")
		replyConnectMsg(ws, int32(lobby.WebsocketConnectError_ParseTokenFailed))
		return
	}

	replyConnectMsg(ws, int32(lobby.WebsocketConnectError_ConnectSuccess))

	var user = newUser(ws, userID)

	oldUser := userMgr.getUserByID(user.uID)
	if oldUser != nil {
		oldUser.detach()
		oldUser.wg.Wait()
	}

	userMgr.addUser(user)

	// lobby.LoginReply(ws, userID)

	defer func() {
		userMgr.removeUser(user)
		user.wg.Done()
	}()

	waitWebsocketMessage(ws, user, r)
}

func acceptWebsocket(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

// InitWith init
func InitWith() {
	userMgr = newUserMgr()

	startAliveKeeper()

	lobby.SetSessionMgr(userMgr)

	lobby.RegHTTPHandle("GET", "/ws", acceptWebsocket)
}
