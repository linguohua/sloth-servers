package main

import (
	"lobbyserver/config"
	"lobbyserver/lobby"
	"lobbyserver/lobby/auth"
	"lobbyserver/lobby/chat"
	"lobbyserver/lobby/mysql"
	"lobbyserver/lobby/pay"
	"lobbyserver/lobby/replay"
	"lobbyserver/lobby/room"
	"lobbyserver/lobby/update"
	"lobbyserver/lobby/mail"
	"lobbyserver/lobby/club"
	"lobbyserver/lobby/donate"
	"lobbyserver/lobby/share"
	"lobbyserver/wechat"

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

	lobby.CreateHTTPServer()
	log.Println("start lobbyserver ok!")

	// 加载各个子模块
	support.InitWith()
	sessions.InitWith()
	room.InitWith()
	replay.InitWith()
	pay.InitWith()
	mysql.InitWith()
	auth.InitWith()
	update.InitWith()
	chat.InitWith()
	mail.InitWith()
	club.InitWith()
	donate.InitWith()
	wechat.InitWechat()
	share.InitWith()

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
