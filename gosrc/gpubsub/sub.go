package gpubsub

import (
	"fmt"
	"gconst"
	"runtime/debug"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"

	log "github.com/sirupsen/logrus"
)

// redisSubscriber 订阅redis频道
func redisSubscriber() {
	for {
		conn := pool.Get()

		psc := redis.PubSubConn{Conn: conn}
		psc.Subscribe(myServerID)
		keep := true
		fmt.Println("begin to wait redis publish msg")
		for keep {
			switch v := psc.Receive().(type) {
			case redis.Message:
				// fmt.Printf("sub %s: message: %s\n", v.Channel, v.Data)
				// 因为只订阅一个主题，因此忽略change参数
				// 同时不可能是
				processRedisPublish(v.Data)
			case redis.Subscription:
				fmt.Printf("sub %s: %s %d\n", v.Channel, v.Kind, v.Count)
			case redis.PMessage:
				fmt.Printf("sub %s: %s %s\n", v.Channel, v.Pattern, v.Data)
			case error:
				log.Println("gpubsub redisSubscriber redis error:", v)
				conn.Close()
				keep = false
				time.Sleep(2 * time.Second)
				break
			}
		}
	}
}

func processRedisPublish(data []byte) {
	loadMsgsAndDispatch()
}

func loadMsgsAndDispatch() {
	if isInProcessState {
		return
	}

	isInProcessState = true

	conn := pool.Get()
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			log.Printf("-----Recovered in processRedisPublish:%v\n", r)
		}

		conn.Close()
		isInProcessState = false
	}()

	var myMsgListID = gconst.LobbyMsgListPrefix + myServerID
	for {
		// 10 msg per-batch
		values, err := redis.Values(conn.Do("LRANGE", myMsgListID, 0, 10))
		if err != nil {
			log.Panic("loadMsgsAndDispatch error:", err)
			break
		}

		valuesCount := 0
		for _, v := range values {
			bytes, err := redis.Bytes(v, nil)
			if err != nil || bytes == nil {
				continue
			}

			ssmsgBag := &gconst.SSMsgBag{}
			err = proto.Unmarshal(bytes, ssmsgBag)
			if err != nil {
				log.Panic("loadMsgsAndDispatch, decode error:", err)
				break
			}

			valuesCount++

			var msgType = ssmsgBag.GetMsgType()
			switch int32(msgType) {
			case int32(gconst.SSMsgType_Notify):
				notifyMsgDispatcher(ssmsgBag)
				break
			case int32(gconst.SSMsgType_Request):
				go requestMsgDispatcher(ssmsgBag)
				break
			case int32(gconst.SSMsgType_Response):
				onPeerServerRespone(ssmsgBag)
				break
			default:
				log.Panicf("No handler for this type %d message", int32(msgType))
			}
		}

		if valuesCount > 0 {
			_, err = conn.Do("LTRIM", myMsgListID, valuesCount, -1)
			if err != nil {
				log.Panic("loadMsgsAndDispatch, LTRIM error:", err)
			}
		}

		if len(values) < 1 {
			break
		}
	}
}

func onPeerServerRespone(ssmsgBag *gconst.SSMsgBag) {
	sn := int(ssmsgBag.GetSeqNO())
	wait, ok := waitingMap[sn]
	if ok {
		delete(waitingMap, sn)
		wait.rspMsg = ssmsgBag
		wait.waitChan <- true
	} else {
		log.Println("processRedisPublish, can't find message sn")
	}
}
