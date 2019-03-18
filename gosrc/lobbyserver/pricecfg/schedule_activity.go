package pricecfg

import (
	"time"

	log "github.com/sirupsen/logrus"
)

var (
// taskList
)

func convertUnixTime2TimeString(unixTime int64) string {
	t := time.Unix(unixTime, 0)
	t.Format("2006-01-02 15:04:05")
	return t.Format("2006-01-02 15:04:05")
}

// 计划关闭折扣活动
func startScheduleActivity(activityCfg *ActivityCfg) {
	if activityCfg == nil {
		log.Println("activityCfg == nil")
		return
	}

	if activityCfg.RoomType == 0 {
		log.Println("activityCfg.RoomType == 0")
		return
	}

	if activityCfg.StartTime > activityCfg.EndTime {
		log.Printf("Invalid confg StartTime:%d > EndTime:%d", activityCfg.StartTime, activityCfg.EndTime)
		return
	}

	var nowTime = time.Now().Unix()
	// 如果活动已经过期，则丢弃
	if activityCfg.EndTime <= nowTime {
		log.Println("activityCfg.EndTime <= nowTime")
		return
	}

	// 如果活动处于执行时间，则放在执行列表中
	if activityCfg.StartTime <= nowTime {
		cfg, ok := cfgs[activityCfg.RoomType]
		if !ok {
			log.Printf("Can't get room type:%d price cfg", activityCfg.RoomType)
			return
		}

		cfg.ActivityPriceCfg = activityCfg
		schedule2StopTask(activityCfg)
		return
	}

	// 剩下的还没到时间执行，则放在等待队列中
	schedule2StartTask(activityCfg)
	schedule2StopTask(activityCfg)
}

// 计划激活任务
func schedule2StartTask(activityCfg *ActivityCfg) {
	log.Printf("schedule2StartTask, roomType:%d, startTime:%s, endTime:%s", activityCfg.RoomType, convertUnixTime2TimeString(activityCfg.StartTime), convertUnixTime2TimeString(activityCfg.EndTime))
	nowTime := time.Now().Unix()
	diff := activityCfg.StartTime - nowTime
	timer2 := time.NewTimer(time.Duration(diff) * time.Second)
	go func() {
		<-timer2.C
		// 启动活动
		cfg, ok := cfgs[activityCfg.RoomType]
		if !ok {
			log.Printf("Schedule start task, but can't get room type:%d price cfg", activityCfg.RoomType)
			return
		}

		cfg.ActivityPriceCfg = activityCfg
		notifyAllUserPriceChange(activityCfg.RoomType, cfg, updatePriceCfg)

		log.Printf("startTask, roomType:%d, endTime:%s", activityCfg.RoomType, convertUnixTime2TimeString(activityCfg.EndTime))
	}()
}

// 计划关闭任务
func schedule2StopTask(activityCfg *ActivityCfg) {
	log.Printf("schedule2StopTask, roomType:%d, startTime:%s, endTime:%s", activityCfg.RoomType, convertUnixTime2TimeString(activityCfg.StartTime), convertUnixTime2TimeString(activityCfg.EndTime))
	nowTime := time.Now().Unix()
	diff := activityCfg.EndTime - nowTime
	timer2 := time.NewTimer(time.Duration(diff) * time.Second)
	go func() {
		<-timer2.C
		cfg, ok := cfgs[activityCfg.RoomType]
		if !ok {
			log.Printf("Schedule start task, but can't get room type:%d price cfg", activityCfg.RoomType)
			return
		}

		// 停止活动
		cfg.ActivityPriceCfg = nil
		notifyAllUserPriceChange(activityCfg.RoomType, cfg, updatePriceCfg)

		log.Printf("StopTask, roomType:%d", activityCfg.RoomType)

	}()
}
