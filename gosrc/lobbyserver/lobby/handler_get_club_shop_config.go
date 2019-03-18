package lobby

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// 从这里开始是跟外部的交互 ---

// OnGetClubShopConfig 客户端拉取商城配置
func OnGetClubShopConfig(w http.ResponseWriter, r *http.Request, userID string) {
	log.Println(r)

	defer r.Body.Close()

	// var token = r.URL.Query().Get("token")
	// playerID, ok := ParseTK(token)
	// if !ok {
	// 	log.Println("ParseToken err")
	// 	w.WriteHeader(500)
	// 	w.Write(CreateRsp(-1, "invalid token"))
	// 	return
	// }

	// vals := r.URL.Query()
	// success := common.VerifySign(&vals, w, playerID)
	// if !success {
	// 	common.Logger.Error("VerifySign err")
	// 	return
	// }

	mapshopconfig, _ := GetAllCommodity(true)
	shopconfig, err := json.Marshal(mapshopconfig)
	if err != nil {
		log.Printf("Marshal err, reason:%s", err)
		w.WriteHeader(500)
		w.Write(CreateRsp(-1, err.Error()))
	}

	w.WriteHeader(200)
	w.Write(CreateRsp(0, fmt.Sprintf("%s", shopconfig)))
}
