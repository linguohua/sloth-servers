package donate

import(
	"lobbyserver/lobby"
	"encoding/json"
)

var (
	donateUtil *myDonateUtil
)

type myDonateUtil struct {

}

func (*myDonateUtil) GetRoomPropsCfg(roomType int) string {
	clientPropCfgMap, ok := clientPropCfgsMap[roomType]
	if !ok {
		return ""
	}

	buf, err := json.Marshal(clientPropCfgMap)
	if err != nil {
		return ""
	}

	return string(buf)
}

// InitWith init
func InitWith() {
	initGamePropCfgs()

	lobby.SetDonateUtil(donateUtil)
}
