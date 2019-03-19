package pay

import(
	"gconst"
)
var (
	roomPriceCfg = make(map[int]map[string]int)
	// payCfg = make(map[int]priceCfg)
)


func init() {
	//原价表，需要保持到redis
	// map的解析：
	// 如：payCfg["o:4:8"] = 48, o表示房主支付,4表示4人场,8表示8局，48表示扣取钻石
	// 如：payCfg["a:4:8"] = 12, a表示AA支付,4表示4人场,8表示8局，12表示扣取钻石
	initDaFengPayCfg()
	initGZPayCfg()
	initNAPayCfg()
	initXiJiangPayCfg()
}

// 初始化大丰麻将价格表
func initDaFengPayCfg() {
	var payCfg = make(map[string]int)
	payCfg["o:4:4"] = 32
	payCfg["o:4:8"] = 48
	payCfg["o:4:16"] = 88

	payCfg["o:3:4"] = 24
	payCfg["o:3:8"] = 36
	payCfg["o:3:16"] = 66

	payCfg["o:2:4"] = 16
	payCfg["o:2:8"] = 24
	payCfg["o:2:16"] = 44

	payCfg["a:4:4"] = 8
	payCfg["a:4:8"] = 12
	payCfg["a:4:16"] = 22

	payCfg["a:3:4"] = 8
	payCfg["a:3:8"] = 12
	payCfg["a:3:16"] = 22

	payCfg["a:2:4"] = 8
	payCfg["a:2:8"] = 12
	payCfg["a:2:16"] = 22
	var roomType = int(gconst.RoomType_DafengMJ)
	roomPriceCfg[roomType] = payCfg
}

// 初始化大丰关张价格表
func initGZPayCfg() {
	var payCfg = make(map[string]int)
	payCfg["o:4:8"] = 66
	payCfg["o:4:16"] = 99
	payCfg["o:4:32"] = 165

	payCfg["o:3:8"] = 66
	payCfg["o:3:16"] = 99
	payCfg["o:3:32"] = 165

	payCfg["o:2:8"] = 66
	payCfg["o:2:16"] = 99
	payCfg["o:2:32"] = 165

	payCfg["a:4:8"] = 22
	payCfg["a:4:16"] = 33
	payCfg["a:4:32"] = 55

	payCfg["a:3:8"] = 22
	payCfg["a:3:16"] = 33
	payCfg["a:3:32"] = 55

	payCfg["a:2:8"] = 22
	payCfg["a:2:16"] = 33
	payCfg["a:2:32"] = 55
	var roomType = int(gconst.RoomType_DafengGZ)
	roomPriceCfg[roomType] = payCfg
}

// 初始化宁安麻将价格表
func initNAPayCfg() {
	var payCfg = make(map[string]int)
	payCfg["o:4:4"] = 20
	payCfg["o:4:8"] = 40
	payCfg["o:4:16"] = 80

	payCfg["o:3:4"] = 20
	payCfg["o:3:8"] = 40
	payCfg["o:3:16"] = 80

	payCfg["o:2:4"] = 20
	payCfg["o:2:8"] = 40
	payCfg["o:2:16"] = 80

	payCfg["a:4:4"] = 5
	payCfg["a:4:8"] = 10
	payCfg["a:4:16"] = 20

	payCfg["a:3:4"] = 5
	payCfg["a:3:8"] = 10
	payCfg["a:3:16"] = 20

	payCfg["a:2:4"] = 5
	payCfg["a:2:8"] = 10
	payCfg["a:2:16"] = 20
	var roomType = int(gconst.RoomType_NingAnMJ)
	roomPriceCfg[roomType] = payCfg
}

// 初始化新疆杠后价格表
func initXiJiangPayCfg() {
	var payCfg = make(map[string]int)
	payCfg["o:4:4"] = 32
	payCfg["o:4:8"] = 48
	payCfg["o:4:16"] = 88

	payCfg["o:3:4"] = 24
	payCfg["o:3:8"] = 36
	payCfg["o:3:16"] = 66

	payCfg["o:2:4"] = 16
	payCfg["o:2:8"] = 24
	payCfg["o:2:16"] = 44

	payCfg["a:4:4"] = 8
	payCfg["a:4:8"] = 12
	payCfg["a:4:16"] = 22

	payCfg["a:3:4"] = 8
	payCfg["a:3:8"] = 12
	payCfg["a:3:16"] = 22

	payCfg["a:2:4"] = 8
	payCfg["a:2:8"] = 12
	payCfg["a:2:16"] = 22
	var roomType = int(gconst.RoomType_XinJiangGH)
	roomPriceCfg[roomType] = payCfg
}