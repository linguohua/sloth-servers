package sessions

import (
	"lobbyserver/lobby"

	"github.com/golang/protobuf/proto"
)

// UserMgr 用户管理
type UserMgr struct {
	users map[string]*User
}

func newUserMgr() *UserMgr {
	um := &UserMgr{}
	um.users = make(map[string]*User)
	return um
}

func (um *UserMgr) getUserByID(userID string) *User {
	user, ok := um.users[userID]
	if ok {
		return user
	}

	return nil
}

func (um *UserMgr) removeUser(user *User) {
	delete(um.users, user.uID)
}

func (um *UserMgr) addUser(user *User) {
	um.users[user.uID] = user
}

// Broacast 广播
func (um *UserMgr) Broacast(msg []byte) {
	for _, u := range um.users {
		u.send(msg)
	}
}

// SendTo 发送消息
func (um *UserMgr) SendTo(userID string, msg []byte) bool {
	u, ok := um.users[userID]
	if ok {
		u.send(msg)

		return true
	}

	return false
}

// SendProtoMsgTo 发送proto msg
func (um *UserMgr) SendProtoMsgTo(userID string, protoMsg proto.Message, opcode int32) bool {
	u, ok := um.users[userID]
	if ok {
		u.sendMsg(protoMsg, opcode)

		return true
	}

	return false
}

// UserCount 当前会话数量
func (um *UserMgr) UserCount() int {
	return len(um.users)
}

// UpdateUserDiamond 通过websocket 更新用户钻石
func (um *UserMgr) UpdateUserDiamond(userID string, diamond uint64) {
	var updateUserDiamond = &lobby.MsgUpdateUserDiamond{}
	updateUserDiamond.Diamond = &diamond
	um.SendProtoMsgTo(userID, updateUserDiamond, int32(lobby.MessageCode_OPUpdateDiamond))
}
