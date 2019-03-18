package lobby

import(
	"gconst"
)
var (

	activityCfgs = make(map[int]*ActivityCfg)
	// payCfg = make(map[int]priceCfg)
)

// ActivityCfg 折扣表
type ActivityCfg struct {
	IsEnable bool `json:"isEnable"`
	StartTime int `json:"startTime"`
	EndTime int `json:"endTime"`
	DiscountCfg map[string]int `json:"discountCfg"`
}

func init() {
	//原价表，需要保持到redis
	// map的解析：
	// 如：payCfg["o:4:8"] = 48, o表示房主支付,4表示4人场,8表示8局，48表示扣取钻石
	// 如：payCfg["a:4:8"] = 12, a表示AA支付,4表示4人场,8表示8局，12表示扣取钻石
	initDaFengDiscountCfg()
	initGZDisscountCfg()
	initNADisscountCfg()
}

// 初始化大丰麻将价格表
func initDaFengDiscountCfg() {

	var discountCfg = make(map[string]int)

	discountCfg["o:4:4"] = 32
	discountCfg["o:4:8"] = 48
	discountCfg["o:4:16"] = 88

	discountCfg["o:3:4"] = 24
	discountCfg["o:3:8"] = 36
	discountCfg["o:3:16"] = 66

	discountCfg["o:2:4"] = 16
	discountCfg["o:2:8"] = 24
	discountCfg["o:2:16"] = 44

	discountCfg["a:4:4"] = 8
	discountCfg["a:4:8"] = 12
	discountCfg["a:4:16"] = 22

	discountCfg["a:3:4"] = 8
	discountCfg["a:3:8"] = 12
	discountCfg["a:3:16"] = 22

	discountCfg["a:2:4"] = 8
	discountCfg["a:2:8"] = 12
	discountCfg["a:2:16"] = 22
	var roomType = int(gconst.RoomType_DafengMJ)
	var activityCfg  = &ActivityCfg{}
	activityCfg.StartTime = 0
	activityCfg.EndTime = 0
	activityCfg.DiscountCfg = discountCfg
	activityCfg.IsEnable = false
	activityCfgs[roomType] = activityCfg
	// discountCfgs[roomType] = discountCfg
}

// 初始化大丰关张价格表
func initGZDisscountCfg() {
	var discountCfg = make(map[string]int)
	discountCfg["o:4:8"] = 66
	discountCfg["o:4:16"] = 99
	discountCfg["o:4:32"] = 165

	discountCfg["o:3:8"] = 66
	discountCfg["o:3:16"] = 99
	discountCfg["o:3:32"] = 165

	discountCfg["o:2:8"] = 66
	discountCfg["o:2:16"] = 99
	discountCfg["o:2:32"] = 165

	discountCfg["a:4:8"] = 22
	discountCfg["a:4:16"] = 33
	discountCfg["a:4:32"] = 55

	discountCfg["a:3:8"] = 22
	discountCfg["a:3:16"] = 33
	discountCfg["a:3:32"] = 55

	discountCfg["a:2:8"] = 22
	discountCfg["a:2:16"] = 33
	discountCfg["a:2:32"] = 55
	var roomType = int(gconst.RoomType_DafengGZ)
	var activityCfg  = &ActivityCfg{}
	activityCfg.StartTime = 0
	activityCfg.EndTime = 0
	activityCfg.DiscountCfg = discountCfg
	activityCfg.IsEnable = false
	activityCfgs[roomType] = activityCfg
}

// 初始化宁安麻将价格表
func initNADisscountCfg() {
	var discountCfg = make(map[string]int)
	discountCfg["o:4:4"] = 10
	discountCfg["o:4:8"] = 20
	discountCfg["o:4:16"] = 40

	discountCfg["o:3:4"] = 10
	discountCfg["o:3:8"] = 20
	discountCfg["o:3:16"] = 40

	discountCfg["o:2:4"] = 10
	discountCfg["o:2:8"] = 20
	discountCfg["o:2:16"] = 40

	discountCfg["a:4:4"] = 3
	discountCfg["a:4:8"] = 5
	discountCfg["a:4:16"] = 10

	discountCfg["a:3:4"] = 3
	discountCfg["a:3:8"] = 5
	discountCfg["a:3:16"] = 10

	discountCfg["a:2:4"] = 3
	discountCfg["a:2:8"] = 5
	discountCfg["a:2:16"] = 10
	var roomType = int(gconst.RoomType_NingAnMJ)
	var activityCfg  = &ActivityCfg{}
	activityCfg.StartTime = 0
	activityCfg.EndTime = 0
	activityCfg.DiscountCfg = discountCfg
	activityCfg.IsEnable = true
	activityCfgs[roomType] = activityCfg
}