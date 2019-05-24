package update

import (
	"encoding/json"
	"lobbyserver/lobby"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var (
	updateUtil = &myUpdateUtil{}
)

// myRoomUtil implements IRoomUtil
type myUpdateUtil struct {
}

// 获取模块配置
func (*myUpdateUtil) GetModuleCfg(r *http.Request) string {
	ctx := parseFromHTTPReq(r)
	cfg := mmgr.findModuleCfg(ctx)
	log.Println("GetModuleCfg, cfg:", cfg)
	if cfg == nil {
		cfg, _ = mmgr.getDefaultCfg(ctx)
	}

	if cfg == nil {
		return ""
	}

	buf, err := json.Marshal(cfg)
	if err != nil {
		log.Error("GetModuleCfg error:", err)
		return ""
	}

	return string(buf)

}

// InitWith init
func InitWith() {
	lobby.SetUpdateUtil(updateUtil)
	lobby.RegHTTPHandle("GET", "/upgradeQuery", handlerUpgradeQuery)
	lobby.RegHTTPHandle("POST", "/webapi/update/uploadModule", handlerUpload)
	lobby.RegHTTPHandle("GET", "/webpi/update/loadAllModules", handlerLoadAllModules)
	lobby.RegHTTPHandle("POST", "/webapi/update/deleteModules", handlerDeleteModules)

	initConditionVariableCfg()

	initModulesMgr()
}
