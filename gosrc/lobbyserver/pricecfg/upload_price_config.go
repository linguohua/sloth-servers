package pricecfg

// "stateless"
// "fmt"
// log "github.com/sirupsen/logrus"
// "github.com/garyburd/redigo/redis"

// func checkPriceConfig(cfgFileBody string) bool {
// 	return true
// }

// func saveOriginalPriceConfig2Redis(cfgFileBody string, roomType int) {
// 	// var pool = accessory.GetRedisPool()
// 	conn := pool.Get()
// 	defer conn.Close()

// 	var key = fmt.Sprintf("%s%d", stateless.PriceConfig, roomType)
// 	conn.Do("HSET", key, "originalPrice", cfgFileBody)
// }

// func saveActivityPriceConfig2Redis(cfgFileBody string, roomType int) {
// 	// var pool = accessory.GetRedisPool()
// 	conn := pool.Get()
// 	defer conn.Close()

// 	var key = fmt.Sprintf("%s%d", stateless.PriceConfig, roomType)
// 	conn.Do("HSET", key, "activityPrice", cfgFileBody)
// }

// func loadPriceConfigFromRedis(roomType int) (string, string) {
// 	// var pool = accessory.GetRedisPool()
// 	conn := pool.Get()
// 	defer conn.Close()

// 	var key = fmt.Sprintf("%s%d", stateless.PriceConfig, roomType)
// 	priceCfgs, err := redis.Strings(conn.Do("HMGET", key, "originPrice", "activityPrice"))
// 	if err != nil {
// 		log.Println("loadPriceConfigFromRedis error:", err)
// 	}

// 	return priceCfgs[0], priceCfgs[1]
// }

// TestPrice 测试用
func TestPrice() {
}
