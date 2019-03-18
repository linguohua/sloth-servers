package lobby

import "time"
import (
	log "github.com/sirupsen/logrus"
)

const (
	keeperAwakeTime = 15 * time.Second
)

func startAliveKeeper() {
	go doAliveKeep()
}

func doAliveKeep() {
	for {
		// 每间隔keeperAwakeTime唤醒一次，唤醒后检查usersMap中的websocket
		// 最后一个消息的接收时间
		time.Sleep(keeperAwakeTime)
		now := time.Now()
		// 如果时间大于90s，则认为客户端已经断开，直接关闭websocket
		for _, v := range userMgr.users {
			diff := now.Sub(v.lastReceivedTime)
			if diff > 90*time.Second {
				log.Println("user not response exceed 90s, close its ws:", v.uID)
				v.ws.Close()
			} else if diff >= 30*time.Second {
				// 如果时间大于30s，则发送一个ping消息
				diff = now.Sub(v.lastPingTime)
				if diff >= 20*time.Second {
					v.sendPing()
					v.lastPingTime = now
				}
			}
		}
	}
}
