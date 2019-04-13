package auth
import (
	"net/http"
	"lobbyserver/lobby"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

func replyRegister(w http.ResponseWriter, registerReply *lobby.MsgRegisterReply) {
	buf, err := proto.Marshal(registerReply)
	if err != nil {
		log.Println("replyRegister, Marshal err:", err)
		return
	}

	w.Write(buf)
}

func handlerRegister(w http.ResponseWriter, r *http.Request) {
	qMod := r.URL.Query().Get("qMod")
	modV := r.URL.Query().Get("modV")
	csVer := r.URL.Query().Get("csVer")
	lobbyVer := r.URL.Query().Get("lobbyVer")
	operatingSystem := r.URL.Query().Get("operatingSystem")
	operatingSystemFamily := r.URL.Query().Get("operatingSystemFamily")
	deviceUniqueIdentifier := r.URL.Query().Get("deviceUniqueIdentifier")
	deviceName := r.URL.Query().Get("deviceName")
	deviceModel := r.URL.Query().Get("deviceModel")
	network := r.URL.Query().Get("network")

	phoneNum := r.URL.Query().Get("phoneNum")

	if phoneNum == "" {
		// TODO: 返回参数错误给客户端
		reply := &lobby.MsgRegisterReply{}
		replyRegister(w,reply)
		return
	}

	// TODO: 检查手机号是否已经注册过, 如果已经注册过，返回错误
	// 如果没注册过，则生成个新用户
	mySQLUtil := lobby.MySQLUtil()
	isRegister := mySQLUtil.CheckPhoneNumIfRegister(phoneNum)
	if isRegister {
		// TODO: 返回错误码给客户端
		reply := &lobby.MsgRegisterReply{}
		replyRegister(w,reply)
		return
	}

	clientInfo := &lobby.ClientInfo{}
	clientInfo.QMod = &qMod
	clientInfo.ModV = &modV
	clientInfo.CsVer = &csVer
	clientInfo.LobbyVer = &lobbyVer
	clientInfo.OperatingSystem = &operatingSystem
	clientInfo.OperatingSystemFamily = &operatingSystemFamily
	clientInfo.DeviceUniqueIdentifier = &deviceUniqueIdentifier
	clientInfo.DeviceName = &deviceName
	clientInfo.DeviceModel = &deviceModel
	clientInfo.Network = &network

	mySQLUtil.UpdateAccountUserInfo(phoneNum, clientInfo)

	// TODO: 返回token给客户端
	reply := &lobby.MsgRegisterReply{}
	replyRegister(w,reply)

}