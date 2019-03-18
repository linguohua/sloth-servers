package prunfast

import "pokerface"

// IState 房间状态机接口
type IState interface {
	onMessage(iu IUser, gmsg *pokerface.GameMessage)
	onStateEnter()
	onStateLeave()
	getStateConst() pokerface.RoomState
	getStateName() string
	onPlayerEnter(player *PlayerHolder)
	onPlayerLeave(player *PlayerHolder)
	onPlayerReEnter(player *PlayerHolder)
}
