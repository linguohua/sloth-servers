package config

import (
	"context"
	"encoding/json"
	"gconst"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/DisposaBoy/JsonConfigReader"
	"github.com/coreos/etcd/client"
)

// make a copy of this file, rename to settings.go
// then set the correct value for these follow variables
var (
	AccessoryServerPort   = 4001
	LogFile               = ""
	Daemon                = "yes"
	RedisServer           = ":6379"
	ServerID              = "27522493-64c3-4899-9e8f-514233ee9f0a"
	SensitiveWordFilePath = "./sensitiveWord.txt"
	GameServerURL         = ""

	WebDataCfgFile = "webdata.json"
	SyncdataTime   = 1

	OrderTestMode = "false"
	PayBackURL    = "test.games.dfppl.xy.qianz.com:81/acc"
	PayURL        = "http://pay.wechat.qianz.com/release/WebPage/OAuthPay/MoreH5PayOrderNew"
	SignKey       = "EE7a1c5bc548e542GBFc340c531657F4"
	PayAPPID      = "10009"

	RoomPayCfgFile = ""

	DbIP       = "localhost"
	DbPort     = 3306
	DbUser     = "root"
	DbPassword = "123456"
	DbName     = "game"

	EtcdServer = ""

	FileServerPath = "./fileServer"
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
		AccessoryServerPort   int    `json:"accessory_server_port"`
		LogFile               string `json:"log_file"`
		Daemon                string `json:"daemon"`
		RedisServer           string `json:"redis_server"`
		ServerID              string `json:"guid"`
		GameServerURL         string `json:"game_server_url"`
		SensitiveWordFilePath string `json:"sensitive_word_file_path"`
		FileServerPath        string `json:"file_server_path"`

		WebDataURL   string `json:"web_data_url"`
		SyncdataTime int    `json:"syncdata_time"`

		OrderTestMode string `json:"order_test_mode"`
		PayBackURL    string `json:"pay_back_url"`
		PayURL        string `json:"pay_url"`
		SignKey       string `json:"sign_key"`
		PayAPPID      string `json:"pay_app_id"`

		RoomPayCfgFile string `json:"room_pay_cfg_file"`

		LobbyID int `json:"lobby_id"`

		EtcdServer string `json:"etcd"`

		DbIP       string `json:"dbIP"`
		DbPort     int    `json:"dbPort"`
		DbPassword string `json:"dbPassword"`
		DbUser     string `json:"dbUser"`
		DbName     string `json:"dbName"`
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
	log.Println(params)

	if params.LogFile != "" {
		LogFile = params.LogFile
	}

	if params.Daemon != "" {
		Daemon = params.Daemon
	}

	if params.AccessoryServerPort != 0 {
		AccessoryServerPort = params.AccessoryServerPort
	}

	if params.RedisServer != "" {
		RedisServer = params.RedisServer
	}

	if params.ServerID != "" {
		ServerID = params.ServerID
	}

	if params.SensitiveWordFilePath != "" {
		SensitiveWordFilePath = params.SensitiveWordFilePath
	}
	/*----------------------------------webserver----------------------------*/

	if params.WebDataURL != "" {
		WebDataCfgFile = params.WebDataURL
	}

	if params.SyncdataTime != 0 {
		SyncdataTime = params.SyncdataTime
	}

	if params.OrderTestMode != "" {
		OrderTestMode = params.OrderTestMode
	}

	if params.PayBackURL != "" {
		PayBackURL = params.PayBackURL
	}

	if params.PayURL != "" {
		PayURL = params.PayURL
	}

	if params.SignKey != "" {
		SignKey = params.SignKey
	}

	if params.PayAPPID != "" {
		PayAPPID = params.PayAPPID
	}

	if params.RoomPayCfgFile != "" {
		RoomPayCfgFile = params.RoomPayCfgFile
	}

	// if params.LobbyID != 0 {
	// 	LobbyID = params.LobbyID
	// }

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

	if params.FileServerPath != "" {
		FileServerPath = params.FileServerPath
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
		log.Println("LoadConfigFromEtcd error:", err)
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

	return true
}
