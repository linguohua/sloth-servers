package zjmahjong

import (
	"mahjong"
)

// SIdle 房间处于空闲状态
type SIdle struct {
	room *Room
}

// onMessage IDLE状态下处理用户消息
func (s *SIdle) onMessage(iu IUser, gmsg *mahjong.GameMessage) {
	s.room.cl.Panic("SIdle should not process message")
}

// onStateEnter 进入IDLE状态
func (s *SIdle) onStateEnter() {
	s.room.cl.Println("room enter idle state")
}

// onStateLeave 离开IDLE状态
func (s *SIdle) onStateLeave() {
	// DO nothing!
	s.room.cl.Println("room leave idle state")
}

// getStateConst IDLE状态标志
func (s *SIdle) getStateConst() mahjong.RoomState {
	return mahjong.RoomState_SRoomIdle
}

// onPlayerEnter IDLE状态处理用户进入：立即转入等待状态
func (s *SIdle) onPlayerEnter(player *PlayerHolder) {
	// 只要有一个用户进来，立即转到waiting状态
	s.room.state2(s, mahjong.RoomState_SRoomWaiting)

	// 让新状态继续处理玩家进入
	s.room.state.onPlayerEnter(player)
}

// onPlayerLeave IDLE状态处理用户离开
func (s *SIdle) onPlayerLeave(player *PlayerHolder) {
	s.room.cl.Panic("SIdle should not process onPlayerLeave")
}

// onPlayerReEnter IDLE状态处理用户重入
func (s *SIdle) onPlayerReEnter(player *PlayerHolder) {
	s.room.cl.Panic("SIdle should not process onPlayerReEnter")
}

func (s *SIdle) getStateName() string {
	return "SIdle"
}
