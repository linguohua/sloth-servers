// 模块更新配置存放于redis的结构：
// 其中，有一个set，用于保存当前所有的模块名字
// 另外，每一个模块名字对应一个hash表，里面的key是版本号，value是整个ModuleCfg的json字符串

package update

import (
	"encoding/json"
	"gconst"
	"lobbyserver/lobby"
	"sort"

	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
)

func loadModuleCfgFromRedis(conn redis.Conn, name string) {
	values, err := redis.Strings(conn.Do("HGETALL", gconst.LobbyUpgradeModuleHashPrefix+name))
	if err != nil {
		log.Error("loadModuleCfgFromRedis, redis error:", err)
		return
	}

	m, ok := mmgr.moduels[name]
	if !ok {
		m = &Module{}
		m.cfgs = make([]*ModuleCfg, 0, 16)
		mmgr.moduels[name] = m
	}

	total := len(values)
	for i := 0; i < total; {
		versionStr := values[i]
		i++
		cfgJSON := values[i]
		i++

		log.Printf("loadModuleCfgFromRedis, name:%s, versionStr:%s\n", name, versionStr)
		mc := &ModuleCfg{}
		err = json.Unmarshal([]byte(cfgJSON), mc)
		if err != nil {
			log.Error("loadModuleCfgFromRedis, json Unmarshal error:", err)
			continue
		}

		// 把需要转换为integer的字符串转换一下，方便后续比较
		mc.strings2Int()
		m.cfgs = append(m.cfgs, mc)
	}

	if len(m.cfgs) > 0 {
		// 排序，按照版本号由高到底排序
		sort.Sort(ByVersion(m.cfgs))

		// 检查有没有默认配置
		for _, mc := range m.cfgs {
			if mc.IsDefault {
				m.defaultCfg = mc
				break
			}
		}
	}
}

func initModulesMgr() {
	// 初始化模块管理器
	mmgr = &ModulesMgr{}
	mmgr.moduels = make(map[string]*Module)

	conn := lobby.Pool().Get()
	defer conn.Close()

	moduleNames, err := redis.Strings(conn.Do("members", gconst.LobbyUpgradeModuleSet))
	if err != nil {
		log.Error("initModulesMgr, redis error:", err)
		return
	}

	for _, mName := range moduleNames {
		loadModuleCfgFromRedis(conn, mName)
	}
}
