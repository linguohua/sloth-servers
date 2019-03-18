package lobby

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lobbyserver/config"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var (

// 这里面的配置全部搬到配置中
// orderTestMode = true
// payBackURL = "test.games.dfppl.xy.qianz.com:81/acc"
// payURL = "http://pay.wechat.qianz.com/release/WebPage/OAuthPay/MoreH5PayOrderNew"
// signKey = "EE7a1c5bc548e542GBFc340c531657F4"
// payAPPID        = "10009"      // 支付标识
)

// Rsp  结果处理
// type Rsp struct {
// 	Ret int    `json:"result"`
// 	Msg string `json:"Message"`
// }

// // CreateRsp httprsp结构
// func CreateRsp(ret int, msg string) []byte {
// 	var rsp Rsp
// 	rsp.Ret = ret
// 	rsp.Msg = msg

// 	returnstr, err := json.Marshal(rsp)
// 	if err != nil {
// 		log.Println(err.Error())
// 	}

// 	return returnstr
// }

// GenSign 生成签名
func genSign(vals ...string) string {
	if len(vals) == 0 {
		return ""
	}

	var str string
	for _, v := range vals {
		str += fmt.Sprintf("%s+", v)
	}

	str += config.SignKey

	data := []byte(str)
	hash := md5.Sum(data)
	msg := fmt.Sprintf("GenSign:%s MD5:%s", str, strings.ToLower(fmt.Sprintf("%x", hash)))
	log.Println(msg)
	return strings.ToLower(fmt.Sprintf("%x", hash))
}

// GetOrderID 生成订单号规则: 日期+userid+gameid+commodityid+自增流水号
func getOrderID(userid string, gameid string, commodityid string) string {
	timestamp := time.Now().UnixNano()
	tm := time.Unix(0, timestamp)
	date := strings.Replace(tm.Format("20060102150405.000"), ".", "", -1)

	incrid := getAutoIncrID()

	return date + gameid + commodityid + incrid
}

// OnCreateClubOrderForWX 客户端俱乐部支付下单商品
func OnCreateClubOrderForWX(w http.ResponseWriter, r *http.Request, userID string) {
	// common.Logger.Debug(r)
	log.Printf("OnCreateClubOrderForWX")

	defer r.Body.Close()

	r.ParseForm()
	// var clubID = r.FormValue("clubID")
	// userID, ok := common.ParseToken(token)
	// if ok != nil {
	// 	common.Logger.Errorf("ParseToken err, reason:%s", ok)
	// 	w.WriteHeader(500)
	// 	w.Write(CreateRsp(-1, "invalid token"))
	// 	return

	// }

	CommoID := r.FormValue("id")
	if CommoID == "" {
		log.Println("OnCreateClubOrderForWX, CommoID id is null")
		w.WriteHeader(500)
		w.Write(CreateRsp(-1, "commodity id is null"))
		return
	}

	clubID := r.FormValue("clubID")
	if clubID == "" {
		log.Println("[OnCreateOrder]  clubID id is null")
		w.WriteHeader(500)
		w.Write(CreateRsp(-1, "club id is null"))
		return
	}

	log.Printf("OnCreateClubOrderForWX,userID:%s commoID:%s, clubID:%s", userID, CommoID, clubID)
	// 检查客户端的sign
	// vals := r.URL.Query()
	// success := common.VerifySign(&vals, w, userID, clubID, CommoID)
	// if !success {
	// 	log.Println("VerifySign err")
	// 	return
	// }

	var addr string
	if r.Header.Get("X-Forwarded-For") == "" {
		addr = strings.Split(r.RemoteAddr, ":")[0]
	} else {
		addr = r.Header.Get("X-Forwarded-For")
	}

	// log.Println(addr)
	// 根据id拿到对应的商品
	Commo := getCommodityByID(CommoID)
	if Commo == nil {
		log.Println("result err, not find commodity")
		w.WriteHeader(500)
		w.Write(CreateRsp(-1, "not find commodity"))
		return
	}

	// 生成单号 --
	OrderID := getOrderID(userID, "0", CommoID)
	// 取出价格
	var price string
	if config.OrderTestMode == "true" {
		price = "1"
	} else {
		price = Commo.GetRealPrice()
	}

	var payBackURL = config.PayBackURL
	// 生成payserver Sgin值
	Sign := genSign(userID, config.PayAPPID, OrderID, price, Commo.ID, payBackURL)

	CommoName := Commo.Name
	if CommoName == "" {
		CommoName = "wjr"
	}

	// 拼接下单URL
	orderURL := config.PayURL + "?appUserId=" + userID +
		"&appId=" + config.PayAPPID +
		"&cpOrderId=" + OrderID +
		"&price=" + price +
		"&waresId=" + CommoID +
		"&waresName=" + CommoName +
		"&sign=" + Sign +
		"&ip=" + addr +
		"&backurl=" + payBackURL
	log.Printf("req:%s", orderURL)

	// 请求下单
	resp, err := http.Post(orderURL, "application/x-www-form-urlencoded", strings.NewReader(""))
	if err != nil {
		log.Printf("Post err, reason:%s", err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	// 结果处理
	var result struct {
		Message   string `json:"Message"`
		Successed bool   `json:"Successed"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("Unmarshal err, reason:%s", err)
		return
	}

	log.Printf("rsp:%v", result)
	if result.Successed {
		order := OrderInfo{}
		order.ID = OrderID
		order.Userid = userID
		order.Clubid = clubID
		order.Appid = config.PayAPPID
		order.Commo.ID = Commo.ID
		order.Commo.Price = Commo.Price
		order.Commo.Value = Commo.Value
		order.Commo.Discount = Commo.Discount
		order.Commo.Extravalue = Commo.Extravalue
		order.Commo.Active = Commo.Active
		order.Status = OrderStatusWaitPay

		setOrderInfo(&order)

		log.Println("success return")
		// 返回支付地址给客户端
		w.WriteHeader(200)
		w.Write(CreateRsp(0, result.Message))
	} else {
		log.Println("result err")
		w.WriteHeader(500)
		w.Write(CreateRsp(-1, result.Message))
	}
}
