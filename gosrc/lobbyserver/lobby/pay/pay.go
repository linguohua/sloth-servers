package pay

import (
	"fmt"
	"lobbyserver/lobby"
)

var (
	payUtil = &myPayUtil{}
)

// myPayUtil implements IPayUtil
type myPayUtil struct {
}

func (*myPayUtil) Refund2Users(roomID string, handFinish int, inGameUserIDs []string) bool {
	orders := refund2Users(roomID, handFinish, inGameUserIDs)
	return len(orders) > 0
}

func (*myPayUtil) DoPayAndSave2RedisWith(roomType int, roomConfigID string,
	roomID string, userID string) (remainDiamond int, errCode int32) {
	return doPayAndSave2RedisWith(roomType, roomConfigID, roomID, userID)
}

func (*myPayUtil) Refund2UserAndSave2Redis(roomID string, userID string, handFinish int) (remainDiamond int, err error) {
	order := refund2UserAndSave2Redis(roomID, userID, handFinish)
	if order == nil {
		return 0, fmt.Errorf("Refund failed, order == nil")
	}

	if order.Refund != nil && order.Refund.Result == 0 {
		return order.Refund.RemainDiamond, nil
	}

	if order.Refund == nil {
		return 0, fmt.Errorf("%s", "Not refund")
	}

	return 0, fmt.Errorf("Refund failed, error code %d", order.Refund.Result)
}

func (*myPayUtil) DoPayAndSave2Redis(roomID string, userID string) (remainDiamond int, errCode int32) {
	return payAndSave2Redis(roomID, userID)
}

// InitWith init
func InitWith() {
	lobby.SetPayUtil(payUtil)
	lobby.RegHTTPHandle("GET", "/loadPrices", handleLoadPrices)
}
