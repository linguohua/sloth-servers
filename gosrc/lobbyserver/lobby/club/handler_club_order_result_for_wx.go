package club

import (
	"fmt"
	"lobbyserver/config"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// OnOrderResultForWX 客户端支付后的回调
func OnOrderResultForWX(w http.ResponseWriter, r *http.Request) {
	log.Println("OnOrderResultForWX")
	log.Println(r)

	defer r.Body.Close()

	r.ParseForm()
	srcsign := r.FormValue("sign")
	orderid := r.FormValue("orderid")
	price := r.FormValue("price")

	// 参数校验
	if srcsign == "" || orderid == "" || price == "" {
		log.Println("err, some param is null")
		w.WriteHeader(500)
		w.Write(CreateRsp(-1, "some param is null"))
		return
	}

	// token校验
	// vals := r.URL.Query()
	// success := common.VerifySign(&vals, w, orderid, price)
	// if !success {
	// 	common.Logger.Error("err, VerifySign err")
	// 	return
	// }

	// 检查一波商品对不对
	orderinfo := getOrderInfo(orderid)
	if orderinfo == nil {
		log.Println("GetOrderInfo err")
		w.WriteHeader(500)
		w.Write(CreateRsp(-1, "not find order"))
		return
	}

	if orderinfo.Status != OrderStatusWaitPay {
		var msg = fmt.Sprintf("Order status is %s", orderinfo.Status)
		w.WriteHeader(500)
		w.Write(CreateRsp(-1, msg))
		return
	}

	var localprice string
	if config.OrderTestMode == "true" {
		localprice = "1"
	} else {
		localprice = orderinfo.Commo.GetRealPrice()
	}

	if localprice == price {
		log.Println("success return payserver")
		w.WriteHeader(200)
		w.Write(CreateRsp(0, "success"))
		// 改变订单状态
		orderinfo.Status = OrderStatusWaitDeliver
		setOrderInfo(orderinfo)

		// TODO: 保存到数据库
		// ret := DeliverCommodity(orderinfo)
		// if !ret {
		// 	common.Logger.Errorf("deliver commodity fail, orderinfo:%v", orderinfo)
		// }
		// 发货，包含了写数据
		deliverCommodity(orderinfo)
	} else {
		log.Println("orader pay fail, need to examine")
		w.WriteHeader(500)
		w.Write(CreateRsp(-1, "fail"))
		orderinfo.Status = OrderStatusError
		setOrderInfo(orderinfo)
	}
}

// DeliverCommodity 发货，包含了写数据
func deliverCommodity(orderinfo *OrderInfo) {
	addValue, _ := strconv.Atoi(orderinfo.Commo.Value)

	addDiamon2ClubAndSendNotify(orderinfo.Userid, orderinfo.Clubid, addValue)

	orderinfo.Status = OrderStatusFinish
	setOrderInfo(orderinfo)

	writeOrderToDB(orderinfo)
}

func updateDiamond() {
	// TODO： 更新俱乐部数据库中的钻石
}
