package zjmahjong

import "mahjong"

// IState 房间状态机接口
type IState interface {
	onMessage(iu IUser, gmsg *mahjong.GameMessage)
	onStateEnter()
	onStateLeave()
	getStateConst() mahjong.RoomState
	getStateName() string
	onPlayerEnter(player *PlayerHolder)
	onPlayerLeave(player *PlayerHolder)
	onPlayerReEnter(player *PlayerHolder)
}
