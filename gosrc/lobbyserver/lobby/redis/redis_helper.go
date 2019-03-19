package lobby

import (
	"context"
	"fmt"
	"gconst"
	"lobbyserver/config"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/coreos/etcd/client"
	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
)

var (
	pool *redis.Pool

	waitingMap = make(map[int]*WaitSubcriberRsp) // 正在等待的集合

	luaScript *redis.Script
)

// WaitSubcriberRsp 等待游戏服务器的返回
type WaitSubcriberRsp struct {
	waitChan chan bool
	rspMsg   *gconst.SSMsgBag
}

// Pool pool
func Pool() *redis.Pool {
	return pool
}

func startRedisClient() {
	pool = newPool(config.RedisServer)

	if config.ServerID == "" {
		log.Panic("Must provide a GUID in config json")
		return
	}

	createLuaScript()

	conn := pool.Get()
	if serverIDSubscriberExist(conn) {
		log.Panicln("The same UUID server instance exists, failed to startup, server ID:", config.ServerID)
		return
	}

	// 新起一个goroutine去订阅redis
	go redisSubscriber()

	// 如果etcd服务器地址确定，则写入自己的版本号
	if config.EtcdServer != "" {
		registerWithEtcd()
	}

}

// registerWithEtcd 往etcd注册自己
func registerWithEtcd() {
	etcdServers := strings.Split(config.EtcdServer, ",")
	cfg := client.Config{
		Endpoints: etcdServers,
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)
	if err != nil {
		log.Fatalln("registerWithEtcd fatal:", err)
	}
	kapi := client.NewKeysAPI(c)

	writeOp := &client.SetOptions{}
	key := fmt.Sprintf("/acc/instances/%s/version", config.ServerID)
	resp, err := kapi.Set(context.Background(), key, strconv.Itoa(GetVersion()), writeOp)
	if err != nil {
		log.Fatal("registerWithEtcd error:", err)
		return
	}

	log.Println("registerWithEtcd ok, resp:", resp)
}

func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

// redisSubscriber 订阅redis频道
func redisSubscriber() {
	for {
		conn := pool.Get()

		psc := redis.PubSubConn{Conn: conn}
		psc.Subscribe(config.ServerID)
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
				log.Println("RoomMgr redisSubscriber redis error:", v)
				conn.Close()
				keep = false
				time.Sleep(2 * time.Second)
				break
			}
		}
	}
}

func serverIDSubscriberExist(conn redis.Conn) bool {
	subCounts, err := redis.Int64Map(conn.Do("PUBSUB", "NUMSUB", config.ServerID))
	if err != nil {
		log.Println("warning: serverIDSubscriberExist, redis err:", err)
	}

	count, _ := subCounts[config.ServerID]
	if count > 0 {
		return true
	}

	return false
}

func processRedisPublish(data []byte) {
	defer func() {
		if r := recover(); r != nil {
			accSysExceptionCount++
			debug.PrintStack()
			log.Printf("-----Recovered in processRedisPublish:%v\n", r)
		}
	}()

	ssmsgBag := &gconst.SSMsgBag{}
	err := proto.Unmarshal(data, ssmsgBag)
	if err != nil {
		log.Println("processRedisPublish, decode error:", err)
		return
	}

	var msgType = ssmsgBag.GetMsgType()
	switch int32(msgType) {
	case int32(gconst.SSMsgType_Notify):
		onNotifyMessage(ssmsgBag)
		break
	case int32(gconst.SSMsgType_Request):
		go onGameServerRequest(ssmsgBag)
		break
	case int32(gconst.SSMsgType_Response):
		onGameServerRespone(ssmsgBag)
		break
	default:
		log.Printf("No handler for this type %d message", int32(msgType))
	}
}

func onGameServerRespone(ssmsgBag *gconst.SSMsgBag) {
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

// sendAndWait 给dst发送消息（通过redis推送），并等待回复，timeout 指定超时时间
func sendAndWait(dst string, msg *gconst.SSMsgBag, timeout time.Duration) (bool, *gconst.SSMsgBag) {
	if dst == "" {
		log.Panicln("publishMessage, need dst")
		return false, nil
	}

	if msg == nil {
		log.Panicln("publishMessage, msg == nil")
		return false, nil
	}

	// 填上源url，以便对方可以发回回复
	msg.SourceURL = &config.ServerID

	var wait = &WaitSubcriberRsp{}
	wait.waitChan = make(chan bool, 1)
	waitingMap[int(msg.GetSeqNO())] = wait

	publishMsg(dst, msg)

	var rspGot = false
	select {
	case <-wait.waitChan:
		rspGot = true
		break
	case <-time.After(timeout):
		break
	}

	// 任何情况都删除这个seqNo
	delete(waitingMap, int(msg.GetSeqNO()))
	return rspGot, wait.rspMsg
}

// publishMsg 往redis publish消息
func publishMsg(dst string, msg *gconst.SSMsgBag) {
	bytes, err := proto.Marshal(msg)
	if err != nil {
		log.Println(err)
		return
	}

	// 获取redis链接，并退出函数时释放
	conn := pool.Get()
	defer conn.Close()

	conn.Do("PUBLISH", dst, bytes)
}

//lua脚本
// KEYS[1] 表前缀stateless.RoomNumberTable
// KEYS[2] roomID
// KEYS[3] roomNumbers
func createLuaScript() {
	script := `for roomNumber in string.gmatch(KEYS[3], '%d+') do
					local value = redis.call('EXISTS', KEYS[1]..roomNumber)
					if value == 0 then
						redis.call('HSET', KEYS[1]..roomNumber, 'roomID', KEYS[2])
						return roomNumber
					end
				end`

	luaScript = redis.NewScript(3, script)
}

// GetRedisPool 导出redisPool
func GetRedisPool() *redis.Pool {
	return pool
}
