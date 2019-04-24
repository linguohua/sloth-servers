package sessions

import (
	"encoding/binary"
	"encoding/hex"
	"gconst"
	"lobbyserver/lobby"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	proto "github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
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
}

// newUser 新建用户对象
func newUser(ws *websocket.Conn, userID string) *User {
	u := &User{}
	u.uID = userID
	u.ws = ws
	u.lastPingTime = time.Now()
	u.lastReceivedTime = time.Now()
	u.wsLock = &sync.Mutex{}
	u.sqllock = &sync.Mutex{}

	u.wg.Add(1)

	return u
}

func (u *User) sendPing() {
	if u.ws != nil {
		u.wsLock.Lock()
		u.ws.WriteMessage(websocket.PingMessage, []byte("ka"))
		u.wsLock.Unlock()
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
	accessoryMessage := &lobby.AccessoryMessage{}
	accessoryMessage.Ops = &ops

	if pb != nil {
		bytes, err := proto.Marshal(pb)

		if err != nil {
			log.Panic("sendMsg, marshal msg failed")
			return
		}
		accessoryMessage.Data = bytes
	}

	bytes, err := proto.Marshal(accessoryMessage)
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
	accessoryMessage := &lobby.AccessoryMessage{}
	err := proto.Unmarshal(message, accessoryMessage)
	if err != nil {
		log.Println(err)
		return
	}

	var msgCode = lobby.MessageCode(accessoryMessage.GetOps())

	switch msgCode {
	// case MessageCode_OPCreateRoom:
	// 	onMessageCreateRoom(u, accessoryMessage)
	// 	break
	// case MessageCode_OPRequestRoomInfo:
	// 	onMessageRequestRoomInfo(u, accessoryMessage)
	// 	break
	case lobby.MessageCode_OPDeleteRoom:
		// onMessageDeleteRoom(u, accessoryMessage)
		break
	case lobby.MessageCode_OPChat:
		// onMessageChat(u, accessoryMessage)
		break
	case lobby.MessageCode_OPUpdateUserInfo:
		onMessageUpdateUserInfo(u, accessoryMessage)
		break
	case lobby.MessageCode_OPVoiceData:
		// onMessageVoiceChat(u, message)
		break
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

func (u *User) saveAuthInfo(userInfo *UserInfo, realIP string) {
	log.Println("saveAuthInfo")
	if userInfo == nil {
		userInfo = &UserInfo{}
	}

	var ws = u.ws
	remoteAddr := ws.RemoteAddr().String()
	addrs := strings.Split(remoteAddr, ":")
	var ip = realIP
	if ip == "" && len(addrs) == 2 {
		ip = addrs[0]
	}

	// 查询钻石
	// TODO: llwant mysql
	// diamond, _ := webdata.QueryDiamond(u.uID)
	// log.Println("ip:", ip)
	diamond := 0
	conn := lobby.Pool().Get()
	defer conn.Close()

	conn.Do("HMSET", gconst.LobbyUserTablePrefix+u.uID, "userName", userInfo.SdkUserName,
		"userNick", userInfo.SdkUserNick, "userSex", userInfo.SdkUserSex, "userLogo", userInfo.SdkUserLogo, "ip", ip, "diaomond", diamond)
}

func (u *User) detach() {
	if u.ws != nil {
		u.ws.Close()
		u.ws = nil
	}
}

func onMessageUpdateUserInfo(user *User, accessoryMessage *lobby.AccessoryMessage) {
	log.Println("onMessageUpdateUserInfo")
	var buf = accessoryMessage.GetData()
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
