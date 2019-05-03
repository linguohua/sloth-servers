package lobby

import (
	"math/rand"

	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

var (
	// SessionMgr mgr
	sessionMgr ISessionMgr
	// RoomUtil room helper functions
	roomUtil IRoomUtil

	payUtil IPayUtil

	mySQLUtil IMySQLUtil

	// RandGenerator rand generator
	RandGenerator *rand.Rand
)

// PayUtil room utility
func PayUtil() IPayUtil {
	if payUtil == nil {
		log.Panic("payUtil is null, maybe not mount pay package yet")
	}

	return payUtil
}

// SetPayUtil set room utility
func SetPayUtil(obj IPayUtil) {
	payUtil = obj
}

// SessionMgr session manager
func SessionMgr() ISessionMgr {
	if sessionMgr == nil {
		log.Panic("sessionMgr is null, maybe not mount sessions package yet")
	}

	return sessionMgr
}

// SetSessionMgr set session manager
func SetSessionMgr(sMgr ISessionMgr) {
	sessionMgr = sMgr
}

// RoomUtil room utility
func RoomUtil() IRoomUtil {
	if roomUtil == nil {
		log.Panic("roomUtil is null, maybe not mount room package yet")
	}

	return roomUtil
}

// SetRoomUtil set room utility
func SetRoomUtil(obj IRoomUtil) {
	roomUtil = obj
}

// MySQLUtil mysql utility
func MySQLUtil() IMySQLUtil {
	if mySQLUtil == nil {
		log.Panic("mySQLUtil is null, maybe not mount room package yet")
	}

	return mySQLUtil
}

// SetMySQLUtil set sql utility
func SetMySQLUtil(obj IMySQLUtil) {
	mySQLUtil = obj
}

// ISessionMgr websocket mgr
type ISessionMgr interface {
	Broacast(msg []byte)
	SendTo(id string, msg []byte) bool
	SendProtoMsgTo(userID string, protoMsg proto.Message, opcode int32) bool
	UserCount() int
}

// IRoomUtil room helper
type IRoomUtil interface {
	LoadLastRoomInfo(userID string) *RoomInfo
	LoadUserLastEnterRoomID(userID string) string
	DeleteRoomInfoFromRedis(roomID string, userID string)
}

// IPayUtil pay
type IPayUtil interface {
	Refund2Users(roomID string, handFinish int, inGameUserIDs []string) bool
	DoPayAndSave2RedisWith(roomType int, roomConfigID string,
		roomID string, userID string) (remainDiamond int, errCode int32)

	Refund2UserAndSave2Redis(roomID string, userID string, handFinish int) (remainDiamond int, err error)

	DoPayAndSave2Redis(roomID string, userID string) (remainDiamond int, errCode int32)
}

// IMySQLUtil sql utility
type IMySQLUtil interface {
	// StartMySQL(ip string, port int, user string, password string, gameDB string)
	UpdateWxUserInfo(UserInfo *UserInfo, clientInfo *ClientInfo) error
	UpdateAccountUserInfo(account string, clientInfo *ClientInfo) error
	GetUserIDBy(account string) string
	GetPasswordBy(account string) string
	GetOrGenerateUserID(account string) (userID string, isNew bool)
	RegisterAccount(account string, passwd string, phone string, userInfo *UserInfo, clientInfo *ClientInfo) error
	LoadUserInfo(userID string) *UserInfo
}
