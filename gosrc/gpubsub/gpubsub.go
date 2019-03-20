package gpubsub

import (
	"gconst"

	"github.com/garyburd/redigo/redis"
)

// MsgDispatchHandler msg dispatch handler
type MsgDispatchHandler func(*gconst.SSMsgBag)

var (
	pool       *redis.Pool
	myServerID string

	notifyMsgDispatcher  MsgDispatchHandler
	requestMsgDispatcher MsgDispatchHandler

	waitingMap = make(map[int]*waitSubcriberRsp) // 正在等待的集合
)

// Startup startup
func Startup(myPool *redis.Pool, myServerID1 string,
	notifyMsgDispatcher1 MsgDispatchHandler, requestMsgDispatcher1 MsgDispatchHandler) {
	pool = myPool
	myServerID = myServerID1
	notifyMsgDispatcher = notifyMsgDispatcher1
	requestMsgDispatcher = requestMsgDispatcher1

	go redisSubscriber()
}
