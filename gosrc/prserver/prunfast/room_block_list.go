package prunfast

import (
	"container/list"
	"time"
)

const (
	blockDuration  = 10 * time.Minute
	blockMaxNumber = 10
)

// RoomBlockUser forbit user to enter
type RoomBlockUser struct {
	userID       string
	blockBeginAt time.Time
}

// RoomBlockList block list
type RoomBlockList struct {
	room *Room
	bl   *list.List
}

func newRoomBlockList(room *Room) *RoomBlockList {
	rbl := &RoomBlockList{}
	rbl.room = room
	rbl.bl = list.New()

	return rbl
}

func (rbl *RoomBlockList) blockUser(userID string) {
	if rbl.isFulled() {
		rbl.bl.Remove(rbl.bl.Front())
	}

	if rbl.has(userID) {
		return
	}

	rbu := &RoomBlockUser{}
	rbu.userID = userID
	rbu.blockBeginAt = time.Now()

	rbl.bl.PushBack(rbu)
}

func (rbl *RoomBlockList) unblockWhenTimePassed() {
	now := time.Now()
	for e := rbl.bl.Front(); e != nil; {
		rbu := (e.Value).(*RoomBlockUser)

		diff := now.Sub(rbu.blockBeginAt)
		if diff > blockDuration {
			var x = e
			e = e.Next()
			rbl.bl.Remove(x)
		} else {
			break
		}
	}
}

func (rbl *RoomBlockList) isFulled() bool {
	return rbl.bl.Len() >= blockMaxNumber
}

func (rbl *RoomBlockList) has(userID string) bool {
	for e := rbl.bl.Front(); e != nil; e = e.Next() {
		rbu := (e.Value).(*RoomBlockUser)
		if rbu.userID == userID {
			return true
		}
	}

	return false
}
