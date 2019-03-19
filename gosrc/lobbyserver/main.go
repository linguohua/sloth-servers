package main

import (
	"lobbyserver/config"
	"lobbyserver/lobby"

	//"accwebserver"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"

	sessions "lobbyserver/lobby/sessions"
	support "lobbyserver/lobby/support"

	log "github.com/sirupsen/logrus"

	// "syncdata"
	"syscall"
	// "webdata"
)

var (
	cfgFilepath   = ""
	etcdServerURL = ""
	serverUUID    = ""
)

func init() {
	flag.StringVar(&cfgFilepath, "c", "", "specify the config file path name")
	flag.StringVar(&etcdServerURL, "e", "", "specify the etcd server URL")
	flag.StringVar(&serverUUID, "u", "", "specify the server UUID")
}

func main() {
	runtime.GOMAXPROCS(1)

	version := flag.Bool("v", false, "show version")

	flag.Parse()

	if *version {
		fmt.Println(lobby.GetVersion())
		os.Exit(0)
	}

	if etcdServerURL != "" {
		config.EtcdServer = etcdServerURL
	}

	if serverUUID != "" {
		config.ServerID = serverUUID
	}

	if cfgFilepath == "" {
		// 如果没有配置json文件，则必须提供uuid以及etcd地址
		if etcdServerURL == "" || serverUUID == "" {
			log.Fatal("must provide etcd and uuid when json config file is omit")
		}
	}

	if cfgFilepath != "" {
		r := config.ParseConfigFile(cfgFilepath)
		if r != true {
			log.Fatal("can't parse configure file:", cfgFilepath)
		}
	} else {
		r := config.LoadConfigFromEtcd()
		if r != true {
			log.Fatal("can't load config from etcd:", etcdServerURL)
		}
	}

	log.Println("try to start lobbyserver...")

	if config.Daemon == "yes" && config.LogFile != "" {
		sighup := make(chan os.Signal, 1)
		signal.Notify(sighup, syscall.SIGHUP)
	}

	// common.Startup(config.CommonCfgPath)

	// cache.Startup(config.CacheCfgPath)
	//初始化webdata 数据
	// webdata.Startup(config.WebDataCfgFile)

	lobby.CreateHTTPServer()
	log.Println("start lobbyserver ok!")

	support.InitWith(lobby.MainRouter)
	sessions.InitWith(lobby.MainRouter)

	// go syncdata.SyncRedisData2DB()
	// go syncdata.StartSyncDataSchedule()

	// go accwebserver.StartWebServer()
	// log.Println("start accWebServer ok!")

	if config.Daemon == "yes" {
		waitForSignal()
	} else {
		waitInput()
	}
	return
}

func waitInput() {
	var cmd string
	for {
		_, err := fmt.Scanf("%s\n", &cmd)
		if err != nil {
			//log.Println("Scanf err:", err)
			continue
		}

		switch cmd {
		case "exit", "quit":
			log.Println("exit by user")
			return
		case "gr":
			log.Println("current goroutine count:", runtime.NumGoroutine())
			break
		case "gd":
			pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
			break
		default:
			break
		}
	}
}

func dumpGoRoutinesInfo() {
	log.Println("current goroutine count:", runtime.NumGoroutine())
	// use DEBUG=2, to dump stack like golang dying due to an unrecovered panic.
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 2)
}
