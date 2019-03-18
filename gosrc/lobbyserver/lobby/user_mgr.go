package lobby

import (
	"math/rand"
	"time"
)

// UserMgr 用户管理
type UserMgr struct {
	users map[string]*User
	rand  *rand.Rand
}

func newUserMgr() *UserMgr {
	um := &UserMgr{}
	um.users = make(map[string]*User)
	um.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
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
