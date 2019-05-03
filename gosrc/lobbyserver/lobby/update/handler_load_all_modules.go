package update

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func handlerLoadAllModules(w http.ResponseWriter, r *http.Request) {
	ml := make([]*ModuleCfg, 0, 16)

	// 把所有的配置下发到客户端
	for _, m := range mmgr.moduels {
		for _, mc := range m.cfgs {
			ml = append(ml, mc)
		}
	}

	bytes, err := json.Marshal(ml)
	if err != nil {
		log.Panicln("update.handlerLoadAllModules json marshal error:", err)
	}

	w.Write(bytes)
}
