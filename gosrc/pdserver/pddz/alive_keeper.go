package pddz

import (
	log "github.com/sirupsen/logrus"
	"runtime/debug"
	"time"
)

const (
	keeperAwakeTime = 5 * time.Second
	diff2Close      = 90
)

func startAliveKeeper() {
	go doAliveKeep()
}

func doAliveKeep() {
	defer func() {
		if r := recover(); r != nil {
			roomExceptionCount++
			debug.PrintStack()
			log.Printf("-----This DoAliveKeep GR will die, Recovered in doAliveKeep:%v\n", r)
		}
	}()
	for {
		// 每间隔keeperAwakeTime唤醒一次，唤醒后检查usersMap中的websocket
		// 最后一个消息的接收时间
		time.Sleep(keeperAwakeTime)
		now := time.Now()
		// 如果时间大于90s，则认为客户端已经断开，直接关闭websocket
		for _, v := range usersMap {
			diff := now.Sub(v.lastReceivedTime)
			if diff > diff2Close*time.Second {
				log.Printf("user not response exceed %ds, close its ws:%s\n", diff2Close, v.user.userID())
				v.user.closeWebsocket()
			} else if diff >= diff2Close/2*time.Second {
				// 如果时间大于30s，则发送一个ping消息
				diff = now.Sub(v.lastPingTime)
				if diff >= diff2Close/3*time.Second {
					v.user.sendPing()
					v.lastPingTime = now
				}
			}
		}
	}
}
