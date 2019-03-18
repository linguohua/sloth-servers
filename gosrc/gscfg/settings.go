package gscfg

import (
	"context"
	"encoding/json"
	"fmt"
	"gconst"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/DisposaBoy/JsonConfigReader"
	"github.com/coreos/etcd/client"
)

// make a copy of this file, rename to settings.go
// then set the correct value for these follow variables
var (
	monitorEstablished = false
	ServerPort         = 3001
	// LogFile              = ""
	Daemon               = "yes"
	RedisServer          = ":6379"
	ServerID             = ""
	RequiredAppModuleVer = 0
	EtcdServer           = ""
	RoomServerID         = ""

	DbIP       = "localhost"
	DbPort     = 1433
	DbUser     = "abc"
	DbPassword = "ab"
	DbName     = "gamedb"

	RoomTypeName string
)

var (
	loadedCfgFilePath = ""
)

// ReLoadConfigFile 重新加载配置
func ReLoadConfigFile() bool {
	log.Println("ReLoadConfigFile-------------------")
	if loadedCfgFilePath == "" {
		log.Println("ReLoadConfigFile-------cfg file path is empty, try load from etcd----")
		if EtcdServer != "" {
			log.Println("ReLoadConfigFile-----------From ETCD--------:", EtcdServer)
			if !LoadConfigFromEtcd() {
				log.Println("ReLoadConfigFile-------------------FAILED")
				return false
			}

			log.Println("ReLoadConfigFile-------------------OK")
			return true
		}

		log.Println("ReLoadConfigFile----FAILED:---neigther cfg file path or etcd is valid")
		return false
	}

	log.Println("ReLoadConfigFile-----------From File--------:", loadedCfgFilePath)
	if !ParseConfigFile(loadedCfgFilePath) {
		log.Println("ReLoadConfigFile-------------------FAILED")
		return false
	}

	log.Println("ReLoadConfigFile-------------------OK")
	return true
}

// ParseConfigFile 解析配置
func ParseConfigFile(filepath string) bool {
	type Params struct {
		ServerPort int `json:"port"`
		// LogFile           string `json:"log_file"`
		Daemon      string `json:"daemon"`
		RedisServer string `json:"redis_server"`
		ServreID    string `json:"guid"`
		// URL         string `json:"url"`

		EtcdServer string `json:"etcd"`

		RequiredAppModuleVer int `json:"requiredAppModuleVer"`

		RoomServerID string `json:"roomServerID"`

		DbIP       string `json:"dbIP"`
		DbPort     int    `json:"dbPort"`
		DbPassword string `json:"dbPassword"`
		DbUser     string `json:"dbUser"`
		DbName     string `json:"dbName"`

		RoomTypeName string `json:"roomTypeName"`
	}

	loadedCfgFilePath = filepath

	var params = &Params{}

	f, err := os.Open(filepath)
	if err != nil {
		log.Println("failed to open config file:", filepath)
		return false
	}

	// wrap our reader before passing it to the json decoder
	r := JsonConfigReader.New(f)
	err = json.NewDecoder(r).Decode(params)

	if err != nil {
		log.Println("json un-marshal error:", err)
		return false
	}

	log.Println("-------------------Configure params are:-------------------")
	log.Printf("%+v\n", params)

	// if params.LogFile != "" {
	// 	LogFile = params.LogFile
	// }

	if params.Daemon != "" {
		Daemon = params.Daemon
	}

	if params.ServerPort != 0 {
		ServerPort = params.ServerPort
	}

	if params.RedisServer != "" {
		RedisServer = params.RedisServer
	}

	if params.ServreID != "" {
		ServerID = params.ServreID
	}

	// if params.URL != "" {
	// 	URL = params.URL
	// }

	RoomTypeName = params.RoomTypeName

	if params.RequiredAppModuleVer > 0 {
		RequiredAppModuleVer = params.RequiredAppModuleVer
	}

	if params.RoomServerID != "" {
		RoomServerID = params.RoomServerID
	}

	if params.EtcdServer != "" {
		EtcdServer = params.EtcdServer
	}

	if params.DbIP != "" {
		DbIP = params.DbIP
	}

	if params.DbUser != "" {
		DbUser = params.DbUser
	}

	if params.DbPassword != "" {
		DbPassword = params.DbPassword
	}

	if params.DbName != "" {
		DbName = params.DbName
	}

	if params.DbPort != 0 {
		DbPort = params.DbPort
	}

	if ServerID == "" {
		log.Println("Server id 'guid' must not be empty!")
		return false
	}

	if EtcdServer != "" {
		if !LoadConfigFromEtcd() {
			return false
		}
	}

	if RoomServerID == "" {
		log.Println("room server id  must not be empty!")
		return false
	}

	if RedisServer == "" {
		log.Println("redis server id  must not be empty!")
		return false
	}

	return true
}

// LoadConfigFromEtcd 从etcd加载配置
func LoadConfigFromEtcd() bool {
	etcdServers := strings.Split(EtcdServer, ",")
	cfg := client.Config{
		Endpoints: etcdServers,
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)
	if err != nil {
		log.Println("LoadConfigFromEtcd, error:", err)
		return false
	}

	kapi := client.NewKeysAPI(c)

	resp, err := kapi.Get(context.Background(), gconst.EtcdRedisServerHost, nil)
	if err == nil {

		// // print common key info
		// log.Printf("Get is done. Metadata is %q\n", resp)
		// // print value
		// log.Printf("%q key has %q value\n", resp.Node.Key, resp.Node.Value)
		RedisServer = resp.Node.Value
		log.Println("load etcd config, redis server:", RedisServer)
	} else {
		log.Println("kapi get:", err)
		return false
	}

	// 确保acc 目录存在
	sopt := &client.SetOptions{}
	sopt.Dir = true
	sopt.PrevExist = client.PrevNoExist
	sopt.NoValueOnSuccess = true
	kapi.Set(context.Background(), gconst.EtcdAccInstanceDir, "", sopt)

	// 选择一个最高版本的acc
	selectHighestAcc()

	gop := &client.GetOptions{}
	gop.Recursive = true
	instanceKey := fmt.Sprintf(gconst.EtcdGameInstancesFormat, ServerID)
	resp, err = kapi.Get(context.Background(), instanceKey, gop)
	if err == nil {
		// // print common key info
		log.Printf("Get is done. Metadata is %q\n", resp)
		// // print value
		// log.Printf("%q key has %q value\n", resp.Node.Key, resp.Node.Value)
		for _, node := range resp.Node.Nodes {
			if node.Key == (instanceKey + "/port") {
				ServerPort, _ = strconv.Atoi(node.Value)
				log.Println("load etc config, game server listen port:", ServerPort)
			}
		}

	} else {
		log.Println("kapi get:", err)
		log.Println("no instance specific config found, use default")
		// return false
	}

	// resp, err = kapi.Get(context.Background(), gconst.EtcdProductURL, nil)
	// if err == nil {
	// 	// print common key info
	// 	log.Printf("Get is done. Metadata is %q\n", resp)
	// 	// // print value
	// 	log.Printf("%q key has %q value\n", resp.Node.Key, resp.Node.Value)

	// 	base, err := url.Parse(resp.Node.Value)
	// 	if err == nil {
	// 		base.Path = path.Join(base.Path, URL)
	// 		URL = base.String()
	// 		log.Println("load etcd config, URL now is:", URL)
	// 	} else {
	// 		log.Println("url.Parse node value error:", err)
	// 	}
	// } else {
	// 	log.Println("kapi get:", err)
	// 	log.Println("no system baseurl config found, use default")
	// 	//return false
	// }

	resp, err = kapi.Get(context.Background(), gconst.EtcdDBGamePassword, nil)
	if err == nil {
		// print common key info
		log.Printf("Get is done. Metadata is %q\n", resp)
		// print value
		log.Printf("%q key has %q value\n", resp.Node.Key, resp.Node.Value)

		DbPassword = resp.Node.Value

	} else {
		log.Println("kapi get:", err)
		return false
	}

	resp, err = kapi.Get(context.Background(), gconst.EtcdDBGameUser, nil)
	if err == nil {
		// print common key info
		log.Printf("Get is done. Metadata is %q\n", resp)
		// print value
		log.Printf("%q key has %q value\n", resp.Node.Key, resp.Node.Value)

		DbUser = resp.Node.Value

	} else {
		log.Println("kapi get:", err)
		return false
	}

	resp, err = kapi.Get(context.Background(), gconst.EtcdDBGameHost, nil)
	if err == nil {
		// print common key info
		log.Printf("Get is done. Metadata is %q\n", resp)
		// print value
		log.Printf("%q key has %q value\n", resp.Node.Key, resp.Node.Value)

		addressPair := strings.Split(resp.Node.Value, ":")
		if len(addressPair) < 2 {
			return false
		}

		DbIP = addressPair[0]
		DbPort, _ = strconv.Atoi(addressPair[1])
	} else {
		log.Println("kapi get:", err)
		return false
	}

	resp, err = kapi.Get(context.Background(), gconst.EtcdDBGameName, nil)
	if err == nil {
		// print common key info
		log.Printf("Get is done. Metadata is %q\n", resp)
		// print value
		log.Printf("%q key has %q value\n", resp.Node.Key, resp.Node.Value)

		DbName = resp.Node.Value
	} else {
		log.Println("kapi get:", err)
		return false
	}

	if !monitorEstablished {
		monitorEstablished = true
		go MonitorAccUUIDChanged()
	}

	return true
}

// MonitorAccUUIDChanged 监控acc变化
func MonitorAccUUIDChanged() {
	for {
		etcdServers := strings.Split(EtcdServer, ",")
		cfg := client.Config{
			Endpoints: etcdServers,
			Transport: client.DefaultTransport,
			// set timeout per request to fail fast when the target endpoint is unavailable
			HeaderTimeoutPerRequest: time.Second,
		}

		c, err := client.New(cfg)
		if err != nil {
			log.Println("monitorAccUUIDChanged:", err)
			time.Sleep(time.Second * 2)
			continue
		}

		kapi := client.NewKeysAPI(c)

		log.Printf("monitorAccUUIDChanged begin to watch %s key change\n",
			gconst.EtcdAccInstanceDir)

		watchOp := &client.WatcherOptions{}
		watchOp.Recursive = true
		watcher := kapi.Watcher(gconst.EtcdAccInstanceDir, watchOp)

		for {
			resp, err := watcher.Next(context.Background())
			if err != nil {
				log.Println("monitorAccUUIDChanged Next:", err)
				break
			}

			log.Printf("Next is done. Metadata is %q\n", resp)
			selectHighestAcc()
		}
	}
}

// selectHighestAcc 选择一个最高版本的ACC服务器，以便游戏服务器可以请求ACC服务器扣费之类
func selectHighestAcc() {
	etcdServers := strings.Split(EtcdServer, ",")
	cfg := client.Config{
		Endpoints: etcdServers,
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)
	if err != nil {
		log.Println("selectHighestAcc, err:", err)
		return
	}

	kapi := client.NewKeysAPI(c)

	gop := &client.GetOptions{}
	gop.Recursive = true
	resp, err := kapi.Get(context.Background(), gconst.EtcdAccInstanceDir, gop)
	if err != nil {
		log.Println("kapi get:", err)
	} else {
		// print common key info
		log.Printf("Get is done. Metadata is %q\n", resp)
		// print value
		log.Printf("%q key has %q value\n", resp.Node.Key, resp.Node.Value)

		var highestVersion = 0
		var highestUUID = ""

		for _, n := range resp.Node.Nodes {
			// 每一个instance
			hv := parseInstanceVersion(n)
			if hv > highestVersion {
				highestVersion = hv
				highestUUID = path.Base(n.Key)
			}
		}

		log.Printf("ACC highestVersion:%d, highestUUID:%s\n", highestVersion,
			highestUUID)

		if highestUUID != "" {
			RoomServerID = highestUUID
		}
	}
}

func parseInstanceVersion(node *client.Node) (ver int) {
	for _, n := range node.Nodes {
		if strings.Contains(n.Key, "version") {
			ver, _ = strconv.Atoi(n.Value)
			break
		}
	}

	return
}

// Regist2Etcd 向etcd注册版本号
func Regist2Etcd(version int, roomType int) {
	etcdServers := strings.Split(EtcdServer, ",")
	cfg := client.Config{
		Endpoints: etcdServers,
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)
	if err != nil {
		log.Println("Regist2Etcd, error:", err)
		return
	}

	kapi := client.NewKeysAPI(c)

	instanceKey := fmt.Sprintf(gconst.EtcdGameInstancesFormat, ServerID)

	_, err = kapi.Set(context.Background(), instanceKey+"/version", fmt.Sprintf("%d", version), nil)
	if err != nil {
		log.Println("Regist2Etcd kapi Set:", err)
		return
	}

	_, err = kapi.Set(context.Background(), instanceKey+"/roomtype", fmt.Sprintf("%d", roomType), nil)
	if err != nil {
		log.Println("Regist2Etcd kapi Set:", err)
		return
	}
}
