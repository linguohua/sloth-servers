package donate

import(
	"lobbyserver/lobby"
	"gconst"
)

var (
	donateUtil *myDonateUtil
)

type myDonateUtil struct {

}

func (*myDonateUtil) GetRoomPropsCfg(roomType int) string {
	return getRoomPropsCfg(roomType)
}


func (*myDonateUtil) DoDoante(propsType uint32, from string, to string, roomType int) (result *gconst.SSMsgDonateRsp, errCode int32) {
	return donate(propsType, from, to, roomType)
}


// InitWith init
func InitWith() {
	initGamePropCfgs()

	lobby.SetDonateUtil(donateUtil)
}
