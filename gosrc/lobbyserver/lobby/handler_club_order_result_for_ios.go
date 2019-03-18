package lobby
import(
	log "github.com/sirupsen/logrus"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
)

type receipt struct {
	OriginalPurchaseDatePst string `json:"original_purchase_date_pst"`
	PurchaseDateMs          string `json:"purchase_date_ms"`
	UniqueIdentifier        string `json:"unique_identifier"`
	OriginalTransactionID   string `json:"original_transaction_id"`
	Bvrs                    string `json:"bvrs"`
	TransactionID           string `json:"transaction_id"`
	Quantity                string `json:"quantity"`
	UniqueVendorIdentifier  string `json:"unique_vendor_identifier"`
	ItemID                  string `json:"item_id"`
	OriginalPurchaseDate    string `json:"original_purchase_date"`
	IsInIntroOfferPeriod    string `json:"is_in_intro_offer_period"`
	ProductID               string `json:"product_id"`
	PurchaseDate            string `json:"purchase_date"`
	IsTrialPeriod           string `json:"is_trial_period"`
	PurchaseDatePst         string `json:"purchase_date_pst"`
	Bid                     string `json:"bid"`
	OriginalPurchaseDateMs  string `json:"original_purchase_date_ms"`
}

type iosrsp struct {
	Status  int     `json:"status"`
	Receipt receipt `json:"receipt"`
}
// OnOrderResultForIOS -- IOS支付后的回调
func OnOrderResultForIOS(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)

	r.ParseForm()
	srcsign := r.FormValue("sign")
	orderid := r.FormValue("orderid")
	receipt := string(body)

	log.Printf("sign:%s, orderid:%s, receipt:%s", srcsign, orderid, string(body))

	// 参数校验
	if srcsign == "" || receipt == "" || orderid == "" {
		log.Printf("err, some param is null")
		w.WriteHeader(500)
		w.Write(CreateRsp(-1, "some param is null"))
		return
	}

	// token校验
	// vals := r.URL.Query()
	// success := common.VerifySign(&vals, w, orderid, receipt)
	// if !success {
	// 	common.Logger.Error("err, VerifySign err")
	// 	return
	// }

	// 检查一波客户端带的商品商品对不对
	orderinfo := getOrderInfo(orderid)
	if orderinfo == nil {
		log.Println("GetOrderInfo err")
		w.WriteHeader(500)
		w.Write(CreateRsp(-1, "not find order"))
		return
	}

	// 去苹果拉一把商品信息
	rsp := getAppleReceiptRsp(receipt)
	if rsp.Status != 0 {
		log.Printf("getAppleReceiptRsp err, status:%d", rsp.Status)
		return
	}

	log.Printf("status:%d, productid:%s", rsp.Status, rsp.Receipt.ProductID)

	// 判断商品是否合法
	if orderinfo.Commo.IOSTag == rsp.Receipt.ProductID {
		log.Printf("success return payserver")
		w.WriteHeader(200)
		w.Write(CreateRsp(0, "success"))
		// 改变订单状态
		orderinfo.Status = OrderStatusWaitDeliver
		setOrderInfo(orderinfo)
		// 发货
		// ret := DeliverCommodity(orderinfo)
		// if !ret {
		// 	log.Printf("deliver commodity fail, orderinfo:%v", orderinfo)
		// }
		deliverCommodity(orderinfo)
	} else {
		log.Printf("orader pay fail, need to examine")
		w.WriteHeader(500)
		w.Write(CreateRsp(-1, "fail"))
		orderinfo.Status = OrderStatusError
		setOrderInfo(orderinfo)
	}
}

// getAppleReceiptRsp 获取苹果支付订单信息
func getAppleReceiptRsp(token string) iosrsp {
	var rsp iosrsp
	result := verifyAppleReceipt("https://buy.itunes.apple.com/verifyReceipt", token)
	if result == nil {
		log.Println("verifyAppleReceipt err, rsp nil")
	} else {
		err := json.Unmarshal([]byte(result), &rsp)
		if err != nil || rsp.Status != 0 {
			log.Printf("verifyAppleReceipt one err, err:%s, status:%d", err, rsp.Status)

			result := verifyAppleReceipt("https://sandbox.itunes.apple.com/verifyReceipt", token)
			if result == nil {
				log.Println("verifyAppleReceipt err, rsp nil")
			} else {
				err := json.Unmarshal([]byte(result), &rsp)
				if err != nil {
					log.Printf("verifyAppleReceipt two err, err:%s", err)
				}
			}
		}
	}

	return rsp
}

// verifyAppleReceipt 苹果商店订单信息查询
func verifyAppleReceipt(verifyurl string, token string) []byte {
	var data struct {
		Token string `json:"receipt-data"`
	}

	data.Token = token
	jsondata, err := json.Marshal(data)
	if err != nil {
		log.Printf("Marshal err, reason:%s", err)
		return nil
	}

	resp, err := http.Post(verifyurl, "application/json;charset=utf-8", strings.NewReader(string(jsondata)))
	if err != nil {
		log.Printf("Post err, reason:%s", err)
		return nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ReadAll err, reason:%s", err)
		return nil
	}

	log.Printf("body:%s", string(body))
	return body
}
