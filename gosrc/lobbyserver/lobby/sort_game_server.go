package lobby

import (
	"sort"
)

// GameServerInfo 保存游戏服务器信息
type GameServerInfo struct {
	serverID string
	version  int
	roomType int
}

// byServerVersion 根据座位ID排序
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

// sortPlayers 根据座位ID排序
func sortGameServer(gameServerInfos []*GameServerInfo) {
	sort.Sort(byServerVersion(gameServerInfos))
}
