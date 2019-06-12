package pddz

import (
	"pokerface"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

const (
	websocketWriteDeadLine = 5 * time.Second
)

// GUser 表示一个游戏用户
type GUser struct {
	uID       string          // 用户唯一ID
	ws        *websocket.Conn // websocket 连接对象
	room      *Room
	info      *UserInfo
	isfromWeb bool
	wsLock    *sync.Mutex // websocket并发写锁
}

func newGUser(userID string, ws *websocket.Conn, room *Room) *GUser {
	gu := &GUser{}
	gu.uID = userID
	gu.room = room
	gu.ws = ws
	gu.wsLock = &sync.Mutex{}
	if room.isForMonkey {
		gu.info = &UserInfo{nick: "", gender: 0, headIconURI: ""}
	} else {
		gu.info = loadUserInfoFromRedis(userID)
	}

	return gu
}

func (gu *GUser) userID() string {
	return gu.uID
}

func (gu *GUser) getRoom() *Room {
	return gu.room
}

func (gu *GUser) send(bytes []byte) {
	ws := gu.ws
	if ws != nil {
		gu.wsLock.Lock()
		defer gu.wsLock.Unlock()

		ws.SetWriteDeadline(time.Now().Add(websocketWriteDeadLine))
		err := ws.WriteMessage(websocket.BinaryMessage, bytes)
		if err != nil {
			ws.Close()
			log.Printf("user %s ws write err:", err)
		}
	}
}

func (gu *GUser) onWebsocketClosed(ws *websocket.Conn) {
	if gu.ws == ws {
		if gu.room != nil {
			gu.ws = nil
			gu.room.onUserOffline(gu, true)
		}
	}
}

func (gu *GUser) onWebsocketMessage(ws *websocket.Conn, message []byte) {
	if gu.ws == ws {
		if gu.room != nil {
			gu.room.onUserMessage(gu, message)
		}
	}
}

func (gu *GUser) sendPing() {
	ws := gu.ws
	if ws != nil {
		gu.wsLock.Lock()
		defer gu.wsLock.Unlock()

		ws.SetWriteDeadline(time.Now().Add(websocketWriteDeadLine))

		var err error
		if gu.isfromWeb {
			buf := formatGameMsgByData([]byte("ka"), int32(pokerface.MessageCode_OPPing))
			ws.WriteMessage(websocket.BinaryMessage, buf)
		} else {
			err = ws.WriteMessage(websocket.PingMessage, []byte("ka"))
		}

		if err != nil {
			log.Printf("user %s ws write err:", err)
			ws.Close()
		}
	}
}

func (gu *GUser) sendPong(msg string) {
	ws := gu.ws
	if ws != nil {
		gu.wsLock.Lock()
		defer gu.wsLock.Unlock()

		if len(msg) == 0 {
			msg = "kr"
		}

		ws.SetWriteDeadline(time.Now().Add(websocketWriteDeadLine))
		err := ws.WriteMessage(websocket.PongMessage, []byte(msg))
		if err != nil {
			log.Printf("user %s ws write err:", err)
			ws.Close()
		}
	}
}

func (gu *GUser) rebind(ws *websocket.Conn) {
	room := gu.room
	gu.detach()

	gu.ws = ws
	gu.room = room
}

func (gu *GUser) detach() {
	gu.room = nil

	if gu.ws != nil {
		ws := gu.ws
		gu.ws = nil
		ws.Close()
	}
}

func (gu *GUser) getInfo() *UserInfo {
	return gu.info
}

// 只有在用户重连进入房间的情况下才从Redis重新拉取信息
func (gu *GUser) updateInfo() {
	userInfo := loadUserInfoFromRedis(gu.uID)
	gu.info = userInfo
}

func (gu *GUser) closeWebsocket() {
	if gu.ws != nil {
		// 仅仅close websocket，不要赋值为nil
		// 否则不能进入room.onUserMessage(gu, message)
		gu.ws.Close()
	}
}

func (gu *GUser) setFromWeb(isFromWeb bool) {
	gu.isfromWeb = isFromWeb
}

func (gu *GUser) isFromWeb() bool {
	return gu.isfromWeb
}
