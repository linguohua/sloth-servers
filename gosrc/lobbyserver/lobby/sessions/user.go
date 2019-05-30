package sessions

import (
	"encoding/binary"
	"encoding/hex"
	"gconst"
	"lobbyserver/lobby"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	proto "github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

const (
	websocketWriteDeadLine = 5 * time.Second
)

// UserInfo 用户信息
type UserInfo struct {
	UserID      int64  `json:"userID"`
	UserName    string `json:"userName"`
	SdkUserName string `json:"sdkUserName"`
	SdkUserNick string `json:"sdkUserNick"`
	SdkUserSex  string `json:"sdkUserSex"`
	SdkUserLogo string `json:"sdkUserLogo"`
}

// User 表示一个用户
type User struct {
	uID              string          // 用户唯一ID
	ws               *websocket.Conn // websocket 连接对象
	wg               sync.WaitGroup
	lastReceivedTime time.Time
	lastPingTime     time.Time

	wsLock *sync.Mutex // websocket并发写锁

	sqllock *sync.Mutex

	isFromWeb bool
}

// newUser 新建用户对象
func newUser(ws *websocket.Conn, userID string, isFromWeb bool) *User {
	u := &User{}
	u.uID = userID
	u.ws = ws
	u.lastPingTime = time.Now()
	u.lastReceivedTime = time.Now()
	u.wsLock = &sync.Mutex{}
	u.sqllock = &sync.Mutex{}
	u.isFromWeb = isFromWeb

	u.wg.Add(1)

	return u
}

func (u *User) sendPing() {
	if u.ws != nil {
		// u.wsLock.Lock()
		// u.ws.WriteMessage(websocket.PingMessage, []byte("ka"))
		// u.wsLock.Unlock()

		u.wsLock.Lock()
		defer u.wsLock.Unlock()

		u.ws.SetWriteDeadline(time.Now().Add(websocketWriteDeadLine))

		var err error
		if u.isFromWeb {
			buf := formatMsgByData([]byte("ka"), int32(lobby.MessageCode_OPPing))
			u.ws.WriteMessage(websocket.BinaryMessage, buf)
		} else {
			err = u.ws.WriteMessage(websocket.PingMessage, []byte("ka"))
		}

		if err != nil {
			log.Printf("user %s ws write err:", err)
			u.ws.Close()
		}
	}
}

func (u *User) sendPong(msg string) {
	if u.ws != nil {
		u.wsLock.Lock()
		u.ws.WriteMessage(websocket.PongMessage, []byte(msg))
		u.wsLock.Unlock()
	}
}

func (u *User) userID() string {
	return u.uID
}

// reBind 重新绑定websocket
func (u *User) reBind(ws *websocket.Conn) {
	if u.ws != nil {
		u.ws.Close()
	}

	u.ws = ws
}

func (u *User) send(bytes []byte) {
	if u.ws != nil {
		// log.Println(string(bytes))
		u.wsLock.Lock()
		u.ws.WriteMessage(websocket.BinaryMessage, bytes)
		u.wsLock.Unlock()
	}
}

func (u *User) sendMsg(pb proto.Message, ops int32) {
	lobbyMessage := &lobby.LobbyMessage{}
	lobbyMessage.Ops = &ops

	if pb != nil {
		bytes, err := proto.Marshal(pb)

		if err != nil {
			log.Panic("sendMsg, marshal msg failed")
			return
		}
		lobbyMessage.Data = bytes
	}

	bytes, err := proto.Marshal(lobbyMessage)
	if err != nil {
		log.Panic("sendMsg, marshal msg failed")
		return
	}

	u.send(bytes)
}

func (u *User) onWebsocketClosed(ws *websocket.Conn) {
	// if u.ws == ws {
	// 	userMgr.removeUser(u)
	// }
}

func (u *User) onWebsocketMessage(ws *websocket.Conn, message []byte) {
	lobbyMessage := &lobby.LobbyMessage{}
	err := proto.Unmarshal(message, lobbyMessage)
	if err != nil {
		log.Println(err)
		return
	}

	var msgCode = lobby.MessageCode(lobbyMessage.GetOps())

	switch msgCode {
	case lobby.MessageCode_OPPing:
		u.lastReceivedTime = time.Now()
		buf := formatMsgByData(lobbyMessage.GetData(), int32(lobby.MessageCode_OPPong))
		u.send(buf)
		break
	case lobby.MessageCode_OPPong:
		u.lastReceivedTime = time.Now()
		break
	// case lobby.MessageCode_OPUpdateUserInfo:
	// 	onMessageUpdateUserInfo(u, accessoryMessage)
	// 	break
	default:
		log.Println("onMessage unsupported msgCode:", msgCode)
		break
	}
}

func userID2Str(userID int64) string {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(userID))
	return hex.EncodeToString(b)
}

func (u *User) detach() {
	if u.ws != nil {
		u.ws.Close()
		u.ws = nil
	}
}

func onMessageUpdateUserInfo(user *User, lobbyMessage *lobby.LobbyMessage) {
	log.Println("onMessageUpdateUserInfo")
	var buf = lobbyMessage.GetData()
	var updateUserInfo = &lobby.MsgUpdateUserInfo{}
	err := proto.Unmarshal(buf, updateUserInfo)
	if err != nil {
		log.Println("onMessageUpdateUserInfo, decode error:", err)
		return
	}

	var userIDstring = user.userID()
	var location = updateUserInfo.GetLocation()
	conn := lobby.Pool().Get()
	defer conn.Close()
	conn.Do("HSET", gconst.LobbyUserTablePrefix+userIDstring, "location", location)
}

func formatMsgByData(data []byte, msgCode int32) []byte {
	lobbyMessage := &lobby.LobbyMessage{}
	lobbyMessage.Ops = &msgCode
	lobbyMessage.Data = data

	buf, err := proto.Marshal(lobbyMessage)
	if err != nil {
		log.Println(err)
		return []byte{}
	}

	return buf

}
