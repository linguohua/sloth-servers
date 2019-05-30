package pay

import (
	"lobbyserver/lobby"
)

var (
	payUtil = &myPayUtil{}
)

// myPayUtil implements IPayUtil
type myPayUtil struct {
}

func (*myPayUtil) Refund2Users(roomID string, handFinish int, inGameUserIDs []string) bool {
	return refund2Users(roomID, handFinish, inGameUserIDs)
}

// 创建房间扣钱
func (*myPayUtil) DoPayForCreateRoom(roomConfigID string,
	roomID string, userID string) (errCode int32) {
	return doPayWith(roomConfigID, roomID, userID)
}

// 进入房间扣钱
func (*myPayUtil) DoPayForEnterRoom(roomID string, userID string) ( errCode int32) {
	return doPayForEnterRoom(roomID, userID)
}

func (*myPayUtil) Refund2UserWith(roomID string, userID string, handFinish int) (errCode int32) {
	return refund2UserWith(roomID, userID, handFinish)
}

// InitWith init
func InitWith() {
	lobby.SetPayUtil(payUtil)
	lobby.RegHTTPHandle("GET", "/loadPrices", handleLoadPrices)
}
