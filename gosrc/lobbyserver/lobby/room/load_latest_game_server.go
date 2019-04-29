package room

import (
	"fmt"
	"gconst"
	"github.com/garyburd/redigo/redis"
	"lobbyserver/lobby"
	"log"
	"sort"
)

// GameServerInfo 保存游戏服务器信息
type GameServerInfo struct {
	serverID string
	version  int
	roomType int
}

// byServerVersion 根据游戏服版本号排序
type byServerVersion []*GameServerInfo

func (s byServerVersion) Len() int {
	return len(s)
}
func (s byServerVersion) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byServerVersion) Less(i, j int) bool {
	return s[i].version > s[j].version
}

// SortGameServer 根据游戏服版本号排序
func sortGameServer(gameServerInfos []*GameServerInfo) {
	sort.Sort(byServerVersion(gameServerInfos))
}

// LoadLatestGameServer 拉取最新版本服务器ID
func loadLatestGameServer(myRoomType int) string {
	log.Println("LoadGameServerID, myRoomType:", myRoomType)
	conn := lobby.Pool().Get()
	defer conn.Close()

	var setkey = fmt.Sprintf("%s%d", gconst.GameServerInstancePrefix, myRoomType)
	log.Println("setkey:", setkey)
	gameServerIDs, err := redis.Strings(conn.Do("SMEMBERS", setkey))
	if err != nil {
		log.Println("get game server keys from redis err: ", err)
		return ""
	}

	log.Println("gameServerIDs:", gameServerIDs)

	conn.Send("MULTI")
	for _, key := range gameServerIDs {
		var gameServerKey = fmt.Sprintf("%s%s", gconst.GameServerInstancePrefix, key)
		log.Println("gameServerKey:", gameServerKey)
		conn.Send("HGET", gameServerKey, "ver")
	}

	values, err := redis.Ints(conn.Do("EXEC"))

	var gameServerInfos = make([]*GameServerInfo, 0, len(values))

	for index, value := range values {
		var ver = value

		serverID := gameServerIDs[index]

		var gsi = &GameServerInfo{}
		gsi.roomType = myRoomType
		gsi.serverID = serverID
		gsi.version = ver
		gameServerInfos = append(gameServerInfos, gsi)
	}

	sortGameServer(gameServerInfos)

	if len(gameServerInfos) > 0 {
		return gameServerInfos[0].serverID
	}
	return ""
}
