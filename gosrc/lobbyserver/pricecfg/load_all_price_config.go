package pricecfg

import (
	"encoding/json"
	"fmt"
	"gconst"
	"time"

	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
)

// 时间格式：YY-MM-DD HH:MM:SS
func convertTimeString2Unix(timeString string) int64 {
	const layOut = "2006-01-02 15:04:05"
	t, err := time.ParseInLocation(layOut, timeString, time.Now().Location())
	if err != nil {
		log.Println("err:", err)
		return 0
	}

	return t.Unix()
}

func parserOriginalPriceCfgStr2Map(origianglPriceCfgStr string) (map[string]int, int) {
	log.Println("parserOriginalPriceCfgStr2Map")

	if origianglPriceCfgStr == "" {
		return nil, 0
	}

	type ParseOriginalPriceCfg struct {
		RoomType int                    `json:"roomType"`
		Prices   map[string]interface{} `json:"prices"`
	}

	var parseOriginalPriceCfg = &ParseOriginalPriceCfg{}
	if err := json.Unmarshal([]byte(origianglPriceCfgStr), parseOriginalPriceCfg); err != nil {
		log.Println("error:", err)
	}

	var priceCfgMap = make(map[string]int)
	for k, v := range parseOriginalPriceCfg.Prices {
		var key = fmt.Sprintf("%s:", k)
		var value = v.(map[string]interface{})
		for k1, v1 := range value {
			var key = fmt.Sprintf("%s%s:", key, k1)
			var value1 = v1.(map[string]interface{})
			for k2, v2 := range value1 {
				var key = fmt.Sprintf("%s%s", key, k2)
				var value2 = v2.(float64)
				priceCfgMap[key] = int(value2)
			}
		}
	}

	// log.Println("priceCfgMap", priceCfgMap)
	return priceCfgMap, parseOriginalPriceCfg.RoomType
}

func parserActivityPriceCfgStr2Map(activityPriceCfgStr string) *ActivityCfg {
	log.Println("parserActivityPriceCfgStr2Map")

	if activityPriceCfgStr == "" {
		return nil
	}

	var activityCfg = &ActivityCfg{}
	type ParseActivityCfg struct {
		StartTime   string                 `json:"startTime"`
		EndTime     string                 `json:"endTime"`
		RoomType    int                    `json:"roomType"`
		DiscountCfg map[string]interface{} `json:"discountCfg"`
	}

	var parseActivityCfg = &ParseActivityCfg{}
	err := json.Unmarshal([]byte(activityPriceCfgStr), parseActivityCfg)
	if err != nil {
		log.Printf("parserActivityPriceCfgStr2Map, Unmarshal error:%v, priceCfgStr:%s", err, activityPriceCfgStr)
		return nil
	}

	var discountCfgMap = make(map[string]int)
	for k, v := range parseActivityCfg.DiscountCfg {
		var key = fmt.Sprintf("%s:", k)
		var value = v.(map[string]interface{})
		for k1, v1 := range value {
			var key = fmt.Sprintf("%s%s:", key, k1)
			var value1 = v1.(map[string]interface{})
			for k2, v2 := range value1 {
				var key = fmt.Sprintf("%s%s", key, k2)
				var value2 = v2.(float64)
				discountCfgMap[key] = int(value2)
			}
		}
	}

	activityCfg.DiscountCfg = discountCfgMap
	activityCfg.StartTime = convertTimeString2Unix(parseActivityCfg.StartTime)
	activityCfg.EndTime = convertTimeString2Unix(parseActivityCfg.EndTime)
	activityCfg.RoomType = parseActivityCfg.RoomType
	return activityCfg
}

// 从redis中拉取所有的价格配置
func loadAllPriceConfigFromRedis() {
	conn := pool.Get()
	defer conn.Close()
	gameRoomTypes, err := redis.Ints(conn.Do("SMEMBERS", gconst.RoomTypeSet))
	if err != nil {
		log.Println("loadAllRoomPropCfgs, err:", err)
		return
	}

	conn.Send("MULTI")
	for _, roomType := range gconst.RoomType_value {
		var k = fmt.Sprintf("%s%d", gconst.PriceConfig, roomType)
		// log.Println("k:", k)
		conn.Send("HMGET", k, "originalPrice", "activityPrice")
	}

	for _, roomType := range gameRoomTypes {
		roomTypeInt32 := int32(roomType)
		_, ok := gconst.RoomType_name[roomTypeInt32]
		if !ok {
			var k = fmt.Sprintf("%s%d", gconst.PriceConfig, roomType)
			conn.Send("HMGET", k, "originalPrice", "activityPrice")
		}
	}

	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		log.Println("loadAllPriceConfigFromRedis error: ", err)
		return
	}

	log.Println("values len:", len(values))
	for _, value := range values {
		fields, err := redis.Strings(value, nil)
		if err != nil {
			log.Println("loadAllPriceConfigFromRedis, error:", err)
			continue
		}

		var originalPriceCfgStr = fields[0]
		if originalPriceCfgStr == "" {
			continue
		}

		var cfg = &Cfg{}
		var originalPriceCfg, roomType = parserOriginalPriceCfgStr2Map(originalPriceCfgStr)
		if originalPriceCfg == nil {
			log.Println("loadAllPriceConfigFromRedis, originalPriceCfg == nil")
			continue
		}

		cfg.OriginalPriceCfg = originalPriceCfg

		if roomType == 0 {
			log.Println("loadAllPriceConfigFromRedis, roomType == 0")
			continue
		}

		cfgs[roomType] = cfg
		var activityPriceCfgStr = fields[1]
		var activityCfg = parserActivityPriceCfgStr2Map(activityPriceCfgStr)
		if activityCfg == nil {
			log.Println(" activityCfg ==  nil")
			continue
		}

		log.Println("try start schedule activity")
		startScheduleActivity(activityCfg)
	}

	// buf, err := json.Marshal(cfgs)
	// if err != nil {
	// 	log.Println("err:", err)
	// 	return
	// }
	// log.Println(string(buf))
}

// LoadAllPriceCfg 加载所以的价格配置
// func LoadAllPriceCfg(p *redis.Pool) {
// 	pool = p
// 	loadAllPriceConfigFromRedis()

// 	println("cfg len :", len(cfgs))
// }
