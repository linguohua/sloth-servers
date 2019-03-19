package pay

import (
	"encoding/json"
	"fmt"
	"lobbyserver/pricecfg"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func loadPricesReplyError(w http.ResponseWriter, errCode int32) {
	var reply = fmt.Sprintf(`{"error":%d}`, errCode)
	w.Write([]byte(reply))
}

func loadPricesReply(w http.ResponseWriter, priceCfgs string) {
	var reply = fmt.Sprintf(`{"priceCfgs":%s}`, priceCfgs)
	w.Write([]byte(reply))
}

func handleLoadPrices(w http.ResponseWriter, r *http.Request, userID string) {
	log.Printf("handleLoadPrices, user %s request load prices", userID)

	if r.ContentLength < 1 {
		log.Println("handleLoadPrices failed, content length is zero")
		loadPricesReplyError(w, int32(MsgError_ErrRequestInvalidParam))
		return
	}

	message := make([]byte, r.ContentLength)
	n, _ := r.Body.Read(message)
	if n != int(r.ContentLength) {
		log.Println("handleLoadPrices failed, content length is not match")
		loadPricesReplyError(w, int32(MsgError_ErrRequestInvalidParam))
		return
	}

	type LoadPriceRequest struct {
		RoomTypes []int `json:"roomTypes"`
	}

	var loadPriceRequest = &LoadPriceRequest{}
	err := json.Unmarshal(message, loadPriceRequest)
	if err != nil {
		log.Println("handleLoadPrices, Unmarshal error:", err)
		loadPricesReplyError(w, int32(MsgError_ErrRequestInvalidParam))
		return
	}

	if len(loadPriceRequest.RoomTypes) < 1 {
		log.Println("handleLoadPrices, loadPriceRequest params len(roomType) < 1")
		loadPricesReplyError(w, int32(MsgError_ErrRequestInvalidParam))
		return
	}

	log.Println("handleLoadPrices, loadPriceRequest.RoomTypes", loadPriceRequest.RoomTypes)

	var priceCfgs = make(map[int]*pricecfg.Cfg)
	for _, roomType := range loadPriceRequest.RoomTypes {
		var cfg = pricecfg.GetPriceCfg(roomType)
		priceCfgs[roomType] = cfg
	}

	buf, err := json.Marshal(priceCfgs)
	if err != nil {
		log.Println("handleLoadPrices, loadPriceRequest params len(roomType) < 1")
		loadPricesReplyError(w, int32(MsgError_ErrEncode))
		return
	}

	loadPricesReply(w, string(buf))
	return
}
