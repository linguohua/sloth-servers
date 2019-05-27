package lobby

import (
	"fmt"
	"gconst"
	"gpubsub"
	"lobbyserver/config"
	"math/rand"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"lobbyserver/pricecfg"

	"github.com/garyburd/redigo/redis"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

const (
	authServer = "http://wjr.u8.login.qianz.com:8080/user/verifyUser?userID=%s&token=%s&sign=%s"
	appKey     = "253c8d16bc73d85ac7066dcae0e478fe"
)

type accUserIDHTTPHandler func(w http.ResponseWriter, r *http.Request, userID string)
type accRawHTTPHandler func(w http.ResponseWriter, r *http.Request)

var (
	// 根router，只有http server看到
	rootRouter           = httprouter.New()
	accSysExceptionCount int // 异常计数
)

func loadCharm(userID string) int32 {
	conn := pool.Get()
	defer conn.Close()

	charm, _ := redis.Int(conn.Do("HGET", gconst.LobbyUserTablePrefix+userID, "charm"))
	return int32(charm)
}

func echoVersion(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte(fmt.Sprintf("version:%d", versionCode)))
}

// RegHTTPHandle 注册HTTP handler
func RegHTTPHandle(method string, path string, handle httprouter.Handle) {
	rootRouter.Handle(method, "/lobby/:uuid"+path, func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		log.Println("RegHTTPHandle")

		var query = r.URL.Query()
		var tk = query.Get("tk")

		if tk != "" {
			userID, result := parseTK(tk)
			if result {
				var p = httprouter.Param{}
				p.Key = "userID"
				p.Value = userID
				ps = append(ps, p)
			}
		}

		handle(w, r, ps)
	})
}

// CreateHTTPServer 启动服务器
func CreateHTTPServer() {
	RandGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))

	startRedisClient()

	loadAllRoomConfigFromRedis()

	pricecfg.LoadAllPriceCfg(pool)

	RegHTTPHandle("GET", "/version", echoVersion)

	// 注册一个文件服务器，以程序当前目录下的web作为根目录
	rootRouter.ServeFiles("/webax/*filepath", http.Dir("./web/dist"))
	gpubsub.Startup(pool, config.ServerID, onNotifyMessage, onGameServerRequest)

	go acceptHTTPRequest()
}

func acceptHTTPRequest() {
	portStr := fmt.Sprintf(":%d", config.AccessoryServerPort)
	s := &http.Server{
		Addr:    portStr,
		Handler: cors.Default().Handler(rootRouter),
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

// func initFileServer() {
// 	// 文件服务器
// 	var gameServerHandler = http.StripPrefix("/lobby/upgrade/download/", http.FileServer(http.Dir(config.FileServerPath)))
// 	rootRouter.PathPrefix("/lobby/upgrade/download/").Handler(gameServerHandler)
// }
