package lobby

import (
	"encoding/json"
	"fmt"
	"gconst"
	"math"
	"strconv"

	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
)

var (
	redisKeyCOMMList        = fmt.Sprintf("%s:%s", gconst.ClubShopPrefix, "commlist")
	redisKeyOrderAutoIncrID = fmt.Sprintf("%s:%s", gconst.ClubShopPrefix, "orderautoincrid")
	// RedisKeyOrderInfo 订单信息
	redisKeyOrderInfo = fmt.Sprintf("%s:%s", gconst.ClubShopPrefix, "orderinfo:")

	// OrderStatusError 异常订单
	OrderStatusError = "-1"
	// OrderStatusWaitPay 支付中
	OrderStatusWaitPay = "0"
	// OrderStatusWaitDeliver 支付完成准备发货
	OrderStatusWaitDeliver = "1"
	// OrderStatusFinish 发货完毕
	OrderStatusFinish = "2"
	// IOSOrderTTL iOS订单过期时间
	IOSOrderTTL = 600
)

// Commodity 商品
type Commodity struct {
	ID         string            `json:"id"`
	Index      string            `json:"index,omitempty"`
	Price      string            `json:"price"`
	Realprece  string            `json:"realprice"` // 单位为分
	Value      string            `json:"value"`
	Discount   string            `json:"discount"`
	Name       string            `json:"name"`
	Icon       string            `json:"icon"`
	Extravalue map[string]string `json:"extravalue"`
	Active     bool              `json:"active"`
	IOSTag     string            `json:"iostag"`
}

// OrderInfo 订单信息
type OrderInfo struct {
	ID     string    `json:"id"`
	Userid string    `json:"userid"`
	Clubid string    `json:"clubid"`
	Appid  string    `json:"appid"`
	Commo  Commodity `json:"commodity"`
	Status string    `json:"status"`
}

// GetRealPrice 获取真实价格 单位为分
func (Commo *Commodity) GetRealPrice() string {
	Price, _ := strconv.ParseFloat(Commo.Price, 32)
	DisCount, _ := strconv.ParseFloat(Commo.Discount, 32)

	Price32 := float32(Price)
	DisCount32 := float32(DisCount)

	result := Price32 * DisCount32
	real := math.Ceil(float64(result))
	RealPrice := fmt.Sprintf("%.0f", real*100)

	return RealPrice
}

// Rsp  结果处理
type Rsp struct {
	Ret int    `json:"result"`
	Msg string `json:"Message"`
}

// CreateRsp httprsp结构
func CreateRsp(ret int, msg string) []byte {
	var rsp Rsp
	rsp.Ret = ret
	rsp.Msg = msg

	returnstr, err := json.Marshal(rsp)
	if err != nil {
		log.Println(err.Error())
	}

	return returnstr
}

// VerifySign 验证签名
// func verifySign(rawVals *url.Values, w http.ResponseWriter, vals ...string) bool {
// 	if !config.Param.SignEnabled {
// 		return true
// 	}

// 	sign, success := QuerySign(rawVals, w)
// 	if !success {
// 		return false
// 	}

// 	localSign := GenSign(vals...)

// 	if sign != localSign {
// 		msg := fmt.Sprintf("client sign:%s  localsign:%s", sign, localSign)
// 		mlog.Logger.Error(msg)
// 		SendFailed(w, 99, "sign error")
// 	}

// 	return sign == localSign
// }

func getAutoIncrID() string {
	conn := pool.Get()
	defer conn.Close()

	value, _ := redis.Int(conn.Do("GET", redisKeyOrderAutoIncrID))
	value++
	conn.Do("SET", redisKeyOrderAutoIncrID, value)

	return fmt.Sprintf("%d", value)
}

// GetAllCommodity 获取所有的Commoddity
// filtervalid为true代表只返回有效的，为false表示全部返回
func GetAllCommodity(filtervalid bool) ([]string, []string) {
	// 获得redis连接
	conn := pool.Get()
	defer conn.Close()

	var storeCommo []string
	var storeIndex []string

	// Commos, _ := redis.Values(conn.Do("ZRANGE", config.RedisKeyCOMMList, 0, -1, "withscores"))
	Commos, _ := redis.Values(conn.Do("ZRANGE", redisKeyCOMMList, 0, -1, "withscores"))

	var nowCommodityStr string
	var nowIndex string
	var Commo Commodity
	for i, v := range Commos {
		if i%2 == 0 {
			nowCommodityStr = fmt.Sprintf("%s", v.([]byte))
			json.Unmarshal(v.([]byte), &Commo)
		} else {
			nowIndex = fmt.Sprintf("%s", v.([]byte))
		}

		if nowIndex != "" {
			if !filtervalid || (filtervalid && Commo.Active) {
				storeCommo = append(storeCommo, nowCommodityStr)
				storeIndex = append(storeIndex, nowIndex)
			}
			nowIndex = ""
		}
	}

	return storeCommo, storeIndex
}

// GetCommodityByID 根据id拿到商品信息
func getCommodityByID(id string) *Commodity {
	CommoList, _ := GetAllCommodity(true)

	var Commo Commodity
	for _, v := range CommoList {
		json.Unmarshal([]byte(v), &Commo)

		if Commo.ID == id {
			log.Printf("success return:%v", Commo)
			return &Commo
		}
	}

	log.Println("[GetCommodityByID] return nil")
	return nil
}

// AddCommodity 增加一个Commodity
func AddCommodity(Commo *Commodity) bool {
	con := pool.Get()
	defer con.Close()

	CommoIndex := Commo.Index
	Commo.Index = ""
	Commo.Realprece = Commo.GetRealPrice()

	if Commo.Realprece == "0" {
		log.Println("err, realprice can't eq zero")
		return false
	}

	data, err := json.Marshal(Commo)
	if err != nil {
		log.Printf("Marshal err, reason:%s", err)
		return false
	}
	StrCommo := fmt.Sprintf("%s", data)

	_, storeIndex := GetAllCommodity(true)
	if inslice(storeIndex, CommoIndex) {
		log.Printf("index:%s already in storeIndex", CommoIndex)
		return false
	}

	result, err := con.Do("ZADD", redisKeyCOMMList, CommoIndex, StrCommo)
	if result == nil {
		log.Printf("zadd err, reason:%s", err)
		return false
	}

	return true
}

// DelCommodity 删除一个Commodity
func DelCommodity(Commo *Commodity) bool {
	con := pool.Get()
	defer con.Close()

	Commo.Index = ""
	data, err := json.Marshal(Commo)
	if err != nil {
		log.Printf("Marshal err, reason:%s", err)
	}
	StrCommodity := fmt.Sprintf("%s", data)

	con.Do("ZREM", redisKeyCOMMList, StrCommodity)

	return true
}

// DelAllCommodity 删除所有Commodity
func DelAllCommodity() bool {
	con := pool.Get()
	defer con.Close()

	con.Do("DEL", redisKeyCOMMList)

	return true
}

// ModifyCommodity 修改一个Commodity的属性 -- 删除Commo，添加AfterCommo
func modifyCommodity(Commo *Commodity, AfterCommo *Commodity) bool {
	DelCommodity(Commo)
	AddCommodity(AfterCommo)

	return true
}

// MakeCommodityValid 启用指定Commodity
func MakeCommodityValid(Commo *Commodity, Index string) bool {
	_, storeIndex := GetAllCommodity(true)
	if inslice(storeIndex, Index) {
		log.Printf("index:%s already in storeIndex", Index)
		return false
	}

	AfterCommo := *Commo
	AfterCommo.Active = true
	AfterCommo.Index = Index

	return modifyCommodity(Commo, &AfterCommo)
}

// MakeCommodityInvalid 让一个Commodity无效
func MakeCommodityInvalid(Commo *Commodity) bool {
	AfterCommo := *Commo
	AfterCommo.Active = false
	AfterCommo.Index = "9999"

	return modifyCommodity(Commo, &AfterCommo)
}

func inslice(s []string, value string) bool {
	for _, v := range s {
		if v == value {
			return true
		}
	}
	return false
}

// SetOrderInfo 存储订单信息
func setOrderInfo(order *OrderInfo) {
	conn := pool.Get()
	defer conn.Close()

	data, err := json.Marshal(order)
	if err != nil {
		log.Printf("Marshal err, reason:%s", err)
		return
	}
	StrOrder := fmt.Sprintf("%s", data)

	conn.Do("SET", redisKeyOrderInfo+order.ID, StrOrder)
}

// GetOrderInfo 取出订单信息
func getOrderInfo(orderid string) *OrderInfo {
	conn := pool.Get()
	defer conn.Close()

	strorder, _ := redis.String(conn.Do("GET", redisKeyOrderInfo+orderid))
	if strorder == "" {
		log.Println("strorder err")
		return nil
	}

	orderinfo := OrderInfo{}
	err := json.Unmarshal([]byte(strorder), &orderinfo)
	if err != nil {
		log.Printf("Unmarshal err, reason:%s", err)
		return nil
	}

	return &orderinfo
}

// addDiamon2ClubAndSendNotify 充值成功后给俱乐部基金加钻石, 然后通知俱乐部
func addDiamon2ClubAndSendNotify(userID string, clubID string, addDiamon int) {
	log.Printf("AddDiamon2Club, clubID:%s, addDiamon:%d", clubID, addDiamon)

	if addDiamon < 0 {
		log.Println("AddDiamon2Club, addDiamon < 0")
		return
	}

	conn := pool.Get()
	defer conn.Close()

	diamon, err := redis.Int(conn.Do("HGET", gconst.ClubTablePrefix+clubID, "diamond"))
	if err != nil {
		log.Printf("Can't get Club %s diamon, init as 0", clubID)
	}

	diamon = diamon + addDiamon

	conn.Do("HSET", gconst.ClubTablePrefix+clubID, "diamond", diamon)

	notifyClubFundAddByShop(addDiamon, diamon, userID, clubID)

	refreshClubFun(uint32(diamon), userID)
}

//RefreshDiamond 通知客户端刷新玩家钻石
func refreshClubFun(diamond uint32, userID string) {
	log.Printf("refreshClubFun, diamond:%d, userID:%s", diamond, userID)
	type UpdateClubFun struct {
		Fund uint32 `json:"fund"`
	}
	updateClubFun := &UpdateClubFun{}
	updateClubFun.Fund = diamond

	buffer, err := json.Marshal(updateClubFun)
	if err != nil {
		log.Printf("player %s RefreshDiamond Marshal err:%v", userID, err)
		return
	}

	push(int32(MessageCode_OPUpdateClubFun), buffer, userID)
}

// WriteOrderToDB 订单数据写到db
func writeOrderToDB(order *OrderInfo) {
	log.Printf("enter, orderinfo:%v", order)

	// // 获取写入时间
	// timestamp := time.Now().UnixNano()
	// tm := time.Unix(0, timestamp)
	// date := strings.Replace(tm.Format("2006-01-02 15:04:05"), ".", "", -1)

	// // 是否是绑推广码首充
	// var firstpay = 0

	// price, _ := strconv.Atoi(order.Commo.GetRealPrice())
	// value, _ := strconv.Atoi(order.Commo.Value)
	// webdata.WriteFriendChargeLog(order.Userid, order.ID, price, value, date, firstpay, 1, 1)
	// TODO: llwant mysql
}
