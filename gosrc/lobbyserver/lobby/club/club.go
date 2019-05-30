package club

import (
	"lobbyserver/lobby"
	"time"
	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
)
var (
	clubMgr *MyClubMgr

	// 2010，作为俱乐部所有计算时间的参考点
	time2010, _ = time.Parse("2006-Jan-02", "2010-Jan-01")

	luaScriptInsertNewEvent *redis.Script
	luaScriptRemoveMemberEventList *redis.Script
)

// createLuaScript lua脚本，用脚本的主要目的是如果把数据拉回到golang端判断
// 可能导致巨大的流量压力，因此用lua脚本在redis端处理后再把结果弄回来。lua脚本执行速度很慢。
func createClubLuaScript() {
    // KEYS[1] list prefix
    // KEYS[2] set prefix
    // KEYS[3] member-set key
    // KEYS[4] eventID
    script := `local prefix = KEYS[1]
		local sprefix = KEYS[2]
		local eventID = KEYS[4]
		local members = redis.call('SMEMBERS', KEYS[3])
		for _,m in pairs(members) do
			redis.call('LPUSH', prefix .. m, eventID)
			redis.call('SADD', sprefix .. m, eventID)
		end`

	luaScriptInsertNewEvent = redis.NewScript(4, script)


	script3 := `local prefix = KEYS[1]
		local sprefix = KEYS[2]
		local members = redis.call('SMEMBERS', KEYS[3])
		for _,m in pairs(members) do
			redis.call('DEL', prefix .. m)
			redis.call('DEL', sprefix .. m)
		end`

	luaScriptRemoveMemberEventList = redis.NewScript(3, script3)
}

func loadAllClub() {
	log.Println("loading club from database")
	mySQLUtil := lobby.MySQLUtil()
	cursor := 0
	count := 100
	for ; ; {
		infos := mySQLUtil.LoadClubInfos(cursor, count)
		clubInfos := infos.([]*MsgClubInfo)
		for _, clubInfo := range clubInfos {
			clubID := clubInfo.GetBaseInfo().GetClubID()
			club, ok := clubMgr.clubs[clubID]
			if !ok {
				club = newBaseClub(clubInfo, clubID)
				clubMgr.clubs[clubID] = club
			}
		}

		if count > len(clubInfos) {
			break
		}

		cursor = cursor + count
	}

}

// InitWith init
func InitWith() {
	clubMgr = newClubMgr()

	loadAllClub()

	createClubLuaScript()

	lobby.SetClubMgr(clubMgr)

	lobby.RegHTTPHandle("GET", "/createClub", onCreateClub)
	lobby.RegHTTPHandle("GET", "/loadMyClubs", onLoadMyClubs)
	lobby.RegHTTPHandle("GET", "/disbandClub", onDisbandClub)
	lobby.RegHTTPHandle("GET", "/loadClubMembers", onLoadClubMembers)
	lobby.RegHTTPHandle("GET", "/joinClub", onJoinClub)
	lobby.RegHTTPHandle("GET", "/joinApproval", onJoinApprove)
	lobby.RegHTTPHandle("GET", "/loadClubEvents", onLoadEvents)
	lobby.RegHTTPHandle("GET", "/quitClub", onQuit)
	lobby.RegHTTPHandle("GET", "/loadMyApplyEvent", onLoadMyApplyEvent)
}
