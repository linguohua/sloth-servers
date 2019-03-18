package pricecfg

import (
	"encoding/json"
	"fmt"
	"gconst"
	"io/ioutil"
	"lobbyserver/config"
	"os"

	"github.com/DisposaBoy/JsonConfigReader"
	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
)

var (
	cfgs = make(map[int]*Cfg)
	pool *redis.Pool
)

// ActivityCfg 折扣表
type ActivityCfg struct {
	// 活动起始时间
	StartTime int64 `json:"startTime"`
	// 活动结束时间
	EndTime int64 `json:"endTime"`
	//房间类型
	RoomType int `json:"roomType"`
	// 折扣后的价格表
	DiscountCfg map[string]int `json:"discountCfg"`
}

// Cfg 价格配置
type Cfg struct {
	// 原价配置
	OriginalPriceCfg map[string]int `json:"originalPriceCfg"`
	// 折扣配置
	ActivityPriceCfg *ActivityCfg `json:"activityPriceCfg"`
}

// GetPriceCfg 获取对应类型房间的价格配置
func GetPriceCfg(roomType int) *Cfg {
	return cfgs[roomType]
}

// LoadAllPriceCfg 加载所以的价格配置
// 先读去取redis里面的配置
// 若redis没有配置，则读取文件中的配置
func LoadAllPriceCfg(p *redis.Pool) {
	pool = p

	loadAllPriceConfigFromRedis()

	originalPriceCfgs := loadOriginalPriceCfgFromFile(config.RoomPayCfgFile)
	for roomType, cfgString := range originalPriceCfgs {
		_, ok := cfgs[roomType]
		if !ok {
			originalCfgMap, _ := parserOriginalPriceCfgStr2Map(cfgString)
			var cfg = &Cfg{}
			cfg.OriginalPriceCfg = originalCfgMap
			cfgs[roomType] = cfg
			savePriceCfg2Redis(roomType, cfgString)
		}
	}

	println("Pricecfg len :", len(cfgs))
	// log.Println(originalPriceCfgs)
}

func loadOriginalPriceCfgFromFile(filePath string) map[int]string {
	var priceCfgMap = make(map[int]string)

	f, err := os.Open(filePath)
	if err != nil {
		log.Println("failed to open config file:", filePath)
		return priceCfgMap
	}

	// wrap our reader before passing it to the json decoder
	r := JsonConfigReader.New(f)
	// err = json.NewDecoder(r).Decode(params)

	// readFile := JsonConfigReader.New(infile)

	jsonBody, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println(err.Error())
		return priceCfgMap
	}

	// log.Println("jsonBody:", string(jsonBody))
	var priceCfgStr = string(jsonBody)

	return parserPriceCfgs(priceCfgStr)
}

func parserPriceCfgs(priceCfgStr string) map[int]string {

	type ParseCfg struct {
		RoomType int `json:"roomType"`
	}

	var priceCfgs []interface{}
	if err := json.Unmarshal([]byte(priceCfgStr), &priceCfgs); err != nil {
		log.Println("error:", err)
	}

	var priceCfgMap = make(map[int]string)

	for _, v := range priceCfgs {
		buf, err := json.Marshal(v)
		if err != nil {
			log.Println("parse config file err:", err)
			continue
		}

		var parseCfg = &ParseCfg{}
		err = json.Unmarshal(buf, parseCfg)
		if err != nil {
			log.Println("parse config file err:", err)
			continue
		}

		if parseCfg.RoomType == 0 {
			continue
		}

		priceCfgMap[parseCfg.RoomType] = string(buf)

	}

	return priceCfgMap
}

func savePriceCfg2Redis(roomType int, cfgString string) {
	conn := pool.Get()
	defer conn.Close()

	key := fmt.Sprintf("%s%d", gconst.PriceConfig, roomType)

	conn.Do("hset", key, "originalPrice", cfgString)
}
