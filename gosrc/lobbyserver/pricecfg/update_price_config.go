package pricecfg

import (
	"encoding/json"
	"fmt"
	"gconst"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/garyburd/redigo/redis"
	uuid "github.com/satori/go.uuid"
)

var (
	updatePriceCfg = 1
	deletePriceCfg = 2
)

// EnablePriceCfg 从redis中取出原价表，然后格式化成map
func EnablePriceCfg(roomType int, priceCfgType string) error {
	log.Printf("EnablePriceCfg, roomType:%d, priceCfgType:%s", roomType, priceCfgType)
	// var pool = accessory.GetRedisPool()
	conn := pool.Get()
	defer conn.Close()

	var key = fmt.Sprintf("%s%d", gconst.PriceConfigDisable, roomType)
	pricecfg, err := redis.String(conn.Do("hget", key, priceCfgType))
	if err != nil {
		log.Println("EnablePriceCfg, get price cfg error:", err)
		return err
	}

	// var key = fmt.Sprintf("%s%d", gconst.PriceConfig, roomType)
	// pricecfg, err := redis.String(conn.Do("HGET", key, priceCfgType))
	// if err != nil {
	// 	log.Println("EnablePriceCfg failed:", err)
	// 	return fmt.Errorf("EnablePriceCfg failed:%v", err)
	// }

	cfg, ok := cfgs[roomType]
	if !ok {
		cfg = &Cfg{}
	}

	// log.Println("pricecfg:", pricecfg)
	if priceCfgType == "originalPrice" {
		originalPriceCfg, _ := parserOriginalPriceCfgStr2Map(pricecfg)
		cfg.OriginalPriceCfg = originalPriceCfg
		cfgs[roomType] = cfg
		notifyAllUserPriceChange(roomType, cfg, updatePriceCfg)
		return nil

	}

	activityPriceCfg := parserActivityPriceCfgStr2Map(pricecfg)
	if activityPriceCfg == nil {
		log.Println("activityCfg == nil")
		return fmt.Errorf("Can't parse activity price cfg")
	}

	if activityPriceCfg.StartTime > activityPriceCfg.EndTime {
		log.Printf("Invalid confg StartTime:%s > EndTime:%s", convertUnixTime2TimeString(activityPriceCfg.StartTime), convertUnixTime2TimeString(activityPriceCfg.EndTime))
		return fmt.Errorf("Invalid confg StartTime:%s > EndTime:%s", convertUnixTime2TimeString(activityPriceCfg.StartTime), convertUnixTime2TimeString(activityPriceCfg.EndTime))
	}

	var nowTime = time.Now().Unix()
	// 如果活动已经过期，则丢弃
	if activityPriceCfg.EndTime <= nowTime {
		log.Println("activityCfg.EndTime <= nowTime")
		return fmt.Errorf("当前活动配置已经过期，请配置正确的结束时间")
	}

	// TODO：需要检查配置是否合法
	// 1.折扣表中的key,在原价表中是否存在
	// 2.折扣价格是否比原价高

	// 如果活动处于执行时间，则放在执行列表中
	if activityPriceCfg.StartTime <= nowTime {
		cfg, ok := cfgs[activityPriceCfg.RoomType]
		if !ok {
			log.Printf("Can't get room type:%d Original price cfg", activityPriceCfg.RoomType)
			return fmt.Errorf("Can't find original price cfg, please upload original price first")
		}

		cfg.ActivityPriceCfg = activityPriceCfg
		notifyAllUserPriceChange(activityPriceCfg.RoomType, cfg, updatePriceCfg)
		schedule2StopTask(cfg.ActivityPriceCfg)
		return nil
	}

	schedule2StartTask(activityPriceCfg)
	schedule2StopTask(activityPriceCfg)

	return nil
}

func notifyAllUserPriceChange(roomType int, pricecfg *Cfg, action int) {
	type NotifyPriceChange struct {
		RoomType int  `json:"roomType"`
		PriceCfg *Cfg `json:"pricecfg"`
		Action   int  `json:"acton"` // 1 添加， 2删除
	}

	var notifyPriceChange = &NotifyPriceChange{}
	notifyPriceChange.RoomType = roomType
	notifyPriceChange.PriceCfg = pricecfg
	notifyPriceChange.Action = action

	buf, err := json.Marshal(notifyPriceChange)
	if err != nil {
		log.Println("err:", err)
		return
	}

	var OPUpdateRoomPriceCfg = 20
	pushAllPlayer(int32(OPUpdateRoomPriceCfg), buf)
}

// PushAllPlayer 推送数据给所有用户
func pushAllPlayer(ops int32, data []byte) {
	conn := pool.Get()
	defer conn.Close()

	list, err := redis.Strings(conn.Do("SUNION", gconst.PushServerList))
	if err != nil {
		log.Println("GetPushServerList err:", err)
		return
	}

	for _, v := range list {
		newUUID, err := uuid.NewV4()
		if err != nil {
			log.Println("onCreateClub, new uuid error:", err)
			return
		}
		key := fmt.Sprintf("client:%d:%d:%s", 0, ops, newUUID)

		conn.Send("MULTI")
		conn.Send("SET", key, data)
		conn.Send("PUBLISH", v, key)
		_, err = conn.Do("EXEC")
		if err != nil {
			log.Printf("push allplayer server:%s event %s err:%v", v, key, err)
			continue
		}
		log.Printf("push allplayer server:%s push event %s", v, key)
	}
}
