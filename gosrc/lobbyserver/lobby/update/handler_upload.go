package update

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	// "lobbyserver/lobby"
	// "github.com/golang/protobuf/proto"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func handlerUpload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 从body中读取json数据
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var mc ModuleCfg
	err = json.Unmarshal(b, &mc)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	mmgr.addNewModuleCfg(&mc)

	w.Write([]byte("ok"))

	log.Info("update.handlerUpload ok, module name:", mc.Name)
}
