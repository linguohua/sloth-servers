package lobby

import (
	"fmt"
	"gconst"
	"gpubsub"
	"lobbyserver/config"
	"math/rand"
	"net/http"
	"path"
	"time"

	log "github.com/sirupsen/logrus"

	"lobbyserver/pricecfg"

	"github.com/garyburd/redigo/redis"

	"github.com/gorilla/mux"
)

const (
	authServer = "http://wjr.u8.login.qianz.com:8080/user/verifyUser?userID=%s&token=%s&sign=%s"
	appKey     = "253c8d16bc73d85ac7066dcae0e478fe"
)

type accUserIDHTTPHandler func(w http.ResponseWriter, r *http.Request, userID string)
type accRawHTTPHandler func(w http.ResponseWriter, r *http.Request)

var (
	// 根router，只有http server看到
	rootRouter = mux.NewRouter()

	accSysExceptionCount int // 异常计数
)

func loadCharm(userID string) int32 {
	conn := pool.Get()
	defer conn.Close()

	charm, _ := redis.Int(conn.Do("HGET", gconst.LobbyUserTablePrefix+userID, "charm"))
	return int32(charm)
}

func echoVersion(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("version:%d", versionCode)))
}

func authorizedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := "/" + path.Base(r.URL.Path)
		h2, ok := AccUserIDHTTPHandlers[requestPath]
		if ok {
			userID, rt := verifyTokenByQuery(r)
			if rt {
				h2(w, r, userID)
			} else {
				w.WriteHeader(404)
				w.Write([]byte("oh, no valid token for handler:" + r.URL.Path))
			}
		} else {
			w.WriteHeader(404)
			w.Write([]byte("oh, can't found handler for:" + r.URL.Path))
		}
	})
}

func trustHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := "/" + path.Base(r.URL.Path)
		h, ok := AccRawHTTPHandlers[requestPath]
		if ok {
			h(w, r)
		} else {
			w.WriteHeader(404)
			w.Write([]byte("oh, can't found handler for:" + r.URL.Path))
		}
	})
}

// CreateHTTPServer 启动服务器
func CreateHTTPServer() {
	RandGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))
	// TODO: just for test, please remove later
	log.Println("For cub test:" + genTK("10024063"))

	startRedisClient()

	initFileServer();

	//loadAllRoomConfigFromRedis()

	pricecfg.LoadAllPriceCfg(pool)

	// 所有模块看到的mainRouter
	// 外部访问需要形如/prunfast/uuid/pok
	var mainRouter = rootRouter.PathPrefix("/lobby/uuid/").Subrouter()
	//mainRouter.HandleFunc("/ws/{wtype}", acceptWebsocket)
	mainRouter.HandleFunc("/version", echoVersion)
	mainRouter.PathPrefix("/untrust").Handler(authorizedHandler())

	mainRouter.PathPrefix("/trust").Handler(trustHandler())
	// log.Println("AccRawHTTPHandlers:", AccRawHTTPHandlers)
	// for k, v := range AccRawHTTPHandlers {
	// 	log.Println("path:", k)
	// 	uhRouter.HandleFunc(k, v)
	// }

	MainRouter = mainRouter

	gpubsub.Startup(pool, config.ServerID, onNotifyMessage, onGameServerRequest)

	go acceptHTTPRequest()
}

func acceptHTTPRequest() {
	portStr := fmt.Sprintf(":%d", config.AccessoryServerPort)
	s := &http.Server{
		Addr:    portStr,
		Handler: rootRouter,
		// ReadTimeout:    10 * time.Second,
		//WriteTimeout:   120 * time.Second,
		MaxHeaderBytes: 1 << 8,
	}

	log.Printf("Http server listen at:%d\n", config.AccessoryServerPort)

	err := s.ListenAndServe()
	if err != nil {
		log.Fatalf("Http server ListenAndServe %d failed:%v", config.AccessoryServerPort, err)
	}

}

func initFileServer() {
		// 文件服务器
		var gameServerHandler = http.StripPrefix("/lobby/upgrade/download/", http.FileServer(http.Dir(config.FileServerPath)))
		rootRouter.PathPrefix("/lobby/upgrade/download/").Handler(gameServerHandler)
}