package lobby

import (
	"encoding/json"
	"fmt"
	"lobbyserver/config"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func setOrderInfoForTTL(order *OrderInfo, ttl int) {
	con := pool.Get()
	defer con.Close()

	data, err := json.Marshal(order)
	if err != nil {
		log.Printf("Marshal err, reason:%s", err)
		return
	}
	StrOrder := fmt.Sprintf("%s", data)

	con.Do("SETEX", redisKeyOrderInfo+order.ID, ttl, StrOrder)
}

// OnCreateOrderForIOS 客户端ios支付下单商品
func OnCreateOrderForIOS(w http.ResponseWriter, r *http.Request, userID string) {
	log.Println("OnCreateOrderForIOS")

	defer r.Body.Close()
	r.ParseForm()
	// var token = r.FormValue("token")
	// UserID, ok := parseToken(token)
	// if ok != nil {
	// 	common.Logger.Errorf("ParseToken err, reason:%s", ok)
	// 	w.WriteHeader(500)
	// 	w.Write(CreateRsp(-1, "invalid token"))
	// 	return
	// }

	commoID := r.FormValue("id")
	if commoID == "" {
		log.Println("ParseToken err, commodity id is null")
		w.WriteHeader(500)
		w.Write(CreateRsp(-1, "commodity id is null"))
		return
	}

	log.Printf("OnCreateOrderForIOS,userID:%s, commoID:%s", userID, commoID)
	// 检查客户端的sign
	// vals := r.URL.Query()
	// success := common.VerifySign(&vals, w, UserID, CommoID)
	// if !success {
	// 	common.Logger.Error("VerifySign err")
	// 	return
	// }

	// log.Println(addr)
	// 根据id拿到对应的商品
	Commo := getCommodityByID(commoID)
	if Commo == nil {
		log.Println("result err, not find commodity")
		w.WriteHeader(500)
		w.Write(CreateRsp(-1, "not find commodity"))
		return
	}

	// 生成单号
	OrderID := getOrderID(userID, "0", commoID)

	// 保存订单信息
	order := OrderInfo{}
	order.ID = OrderID
	order.Userid = userID
	order.Appid = config.PayAPPID
	order.Commo.ID = Commo.ID
	order.Commo.Price = Commo.Price
	order.Commo.Value = Commo.Value
	order.Commo.Discount = Commo.Discount
	order.Commo.Extravalue = Commo.Extravalue
	order.Commo.Active = Commo.Active
	order.Commo.IOSTag = Commo.IOSTag
	order.Status = OrderStatusWaitPay

	setOrderInfoForTTL(&order, IOSOrderTTL)

	log.Println("success return")
	// 返回订单号给客户端
	w.WriteHeader(200)
	w.Write(CreateRsp(0, order.ID))
}
