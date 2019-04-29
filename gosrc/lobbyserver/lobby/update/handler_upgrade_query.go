package update

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	queryErrorParamQModIsNull  = 1
	queryErrorParamModVIsNull  = 2
	queryErrorModuleNotExist   = 3
	queryErrorUnmarshalCfg     = 4
	queryErrorNeedUpgradeCS    = 5
	queryErrorNeedUpgradeLobby = 6
)

// UpgradeQueryReply 更新查询回复
type UpgradeQueryReply struct {
	Code      int               `json:"code"`
	ABValid   bool              `json:"abValid"`
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	ABList    []AssetsBundleCfg `json:"abList"`
	DlRootURL string            `json:"rootURL"` // 客户端根据该地址拼接完整的资源下载地址
}

func replyUpgradeQuery(w http.ResponseWriter, reply *UpgradeQueryReply) {
	buf, err := json.Marshal(reply)
	if err != nil {
		log.Println("replyWxLogin, Marshal err:", err)
		return
	}

	w.Write(buf)
}

func handlerUpgradeQuery(w http.ResponseWriter, r *http.Request) {
	reply := &UpgradeQueryReply{}
	reply.Code = 0
	reply.ABValid = false

	// 根据request构造查询上下文
	fctx := parseFromHTTPReq(r)
	// 寻找可用的更新配置
	cfg := mmgr.findModuleCfg(fctx)
	if cfg == nil {
		// 没有可用的更新配置
		// 检查是否存在强制更新，也即是模块是否存在默认模块
		var code int
		cfg, code = mmgr.getDefaultCfg(fctx)
		reply.Code = code
	}

	if cfg != nil {
		// 发现更新
		reply.ABValid = true
		reply.Name = cfg.Name
		reply.Version = cfg.Version
		reply.ABList = cfg.AbList
	}

	replyUpgradeQuery(w, reply)
}
