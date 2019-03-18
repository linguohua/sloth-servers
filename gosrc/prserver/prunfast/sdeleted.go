package prunfast

import (
	"pokerface"
)

// SDeleted 房间处于空闲状态
type SDeleted struct {
	room *Room
}

func (s *SDeleted) onMessage(iu IUser, gmsg *pokerface.GameMessage) {
	s.room.cl.Panic("SDeleted should not process message")
}

func (s *SDeleted) onStateEnter() {
	s.room.cl.Println("room enter deleted state")
}

func (s *SDeleted) onStateLeave() {
	// DO nothing!
	s.room.cl.Panic("SDeleted should not leave")
}

func (s *SDeleted) getStateConst() pokerface.RoomState {
	return pokerface.RoomState_SRoomDeleted
}

func (s *SDeleted) getStateName() string {
	return "SDeleted"
}

func (s *SDeleted) onPlayerEnter(player *PlayerHolder) {
	s.room.cl.Panic("SDeleted should not process onPlayerEnter")
}

func (s *SDeleted) onPlayerLeave(player *PlayerHolder) {
	s.room.cl.Panic("SDeleted should not process onPlayerLeave")
}

func (s *SDeleted) onPlayerReEnter(player *PlayerHolder) {
	s.room.cl.Panic("SDeleted should not process onPlayerReEnter")
}
