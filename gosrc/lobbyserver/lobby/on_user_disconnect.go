package lobby

import (
	"gconst"

	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
)

func subscriberUserDisConnectEvent() {
	events := map[string]func(string, []byte){
		gconst.EventPlayerDisconnectAcc: onDisConnect,
	}

	// cache.Sub(handlers)

	go func() {
		conn := pool.Get()
		defer conn.Close()

		channels := redis.Args{}
		for key := range events {
			channels = channels.Add(key)
		}

		psc := redis.PubSubConn{Conn: conn}
		psc.Subscribe(channels...)
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				log.Printf("sub %s message: %s\n", v.Channel, v.Data)
				handler := events[v.Channel]
				if handler != nil {
					value, exists := getPubData(string(v.Data))
					if exists {
						handler(string(v.Data), value)
					}
				}
			case redis.Subscription:
				log.Printf("sub %s %s %d\n", v.Channel, v.Kind, v.Count)
			case redis.PMessage:
				log.Printf("sub %s %s %s\n", v.Channel, v.Pattern, v.Data)
			case error:
				log.Panic(v)
			}
		}
	}()
}

func onDisConnect(msg string, data []byte) {
	playerID := string(data)

	log.Println("onDisConnect, userID:", playerID)
	//chost.presentNotifyFunc(playerID, false)
}
