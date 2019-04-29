package lobby

import (
	"context"
	"fmt"
	"lobbyserver/config"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/coreos/etcd/client"
	"github.com/garyburd/redigo/redis"
)

var (
	pool *redis.Pool

	// LuaScript lua script
	LuaScript *redis.Script
)

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

//lua脚本
// KEYS[1] 表前缀stateless.LobbyRoomNumberTablePrefix
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

	LuaScript = redis.NewScript(3, script)
}
