package pay

import (
	"encoding/json"
	"lobbyserver/pricecfg"
	"net/http"

	"github.com/julienschmidt/httprouter"

	log "github.com/sirupsen/logrus"
)

func handleLoadPrices(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("userID")
	log.Printf("handleLoadPrices, user %s request load prices", userID)

	priceCfgs := pricecfg.GetAllPriceCfgs()

	buf, err := json.Marshal(priceCfgs)
	if err != nil {
		log.Error("Marshal json error:", err)
		return
	}

	w.Write(buf)
	return
}
