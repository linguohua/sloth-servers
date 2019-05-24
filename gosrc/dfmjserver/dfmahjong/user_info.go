package dfmahjong

import (
	"gconst"
	"log"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

// UserInfo 用户基本信息
type UserInfo struct {
	nick        string
	sex         uint32
	headIconURI string
	ip          string
	location    string
	dfHands     int
	diamond     int // 钻石数量
	charm       int // 魅力数量，可能是负数
	avatarID    int // 头像框ID
	clubIDs     []string
	dan         int // 段位

}

func loadUserInfoFromRedis(userID string) *UserInfo {
	conn := pool.Get()
	defer conn.Close()

	var userInfo = &UserInfo{}

	conn.Send("MULTI")
	conn.Send("HMGET", gconst.LobbyUserTablePrefix+userID, "Nick", "Sex", "Protrait", "Addr", "location", "diamond", "charm", "AvatarID", "DanID")
	conn.Send("HGET", gconst.LobbyPlayerTablePrefix+userID, "dfHands")
	// conn.Send("HGETALL", gconst.UserClubTablePrefix+userID)
	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		log.Println("loadUserInfoFromRedis error: ", err)
		return userInfo
	}

	fileds, err := redis.Strings(values[0], nil)
	if err != nil {
		log.Println("parser fileds error: ", err)
		return userInfo
	}

	userInfo.nick = fileds[0]

	sex, _ := strconv.ParseUint(fileds[1], 10, 32)
	userInfo.sex = uint32(sex)

	userInfo.headIconURI = fileds[2]
	userInfo.ip = fileds[3]
	userInfo.location = fileds[4]

	diamond, err := strconv.ParseInt(fileds[5], 10, 32)
	if err != nil {
		userInfo.diamond = 0
	} else {
		userInfo.diamond = int(diamond)
	}

	charm, err := strconv.ParseInt(fileds[6], 10, 32)
	if err != nil {
		userInfo.charm = 0
	} else {
		userInfo.charm = int(charm)
	}

	avatarID, err := strconv.ParseInt(fileds[7], 10, 32)
	if err != nil {
		userInfo.avatarID = 0
	} else {
		userInfo.avatarID = int(avatarID)
	}

	dan, err := strconv.ParseInt(fileds[8], 10, 32)
	if err != nil {
		userInfo.dan = 0
	} else {
		userInfo.dan = int(dan)
	}

	log.Printf("loadUserInfoFromRedis, user %s avatarID:%d", userID, avatarID)
	dfHands, err := redis.Int(values[1], nil)
	if err != nil {
		log.Println("parse int error: ", err)
		dfHands = 0
	}

	userInfo.dfHands = dfHands

	// vs, err := redis.Strings(values[2], nil)
	// if err != nil {
	// 	log.Println("parser fileds error: ", err)
	// 	return userInfo
	// }

	// var clubIDs = make([]string, 0, len(vs)/2)
	// for i := 0; i < len(vs); i = i + 2 {
	// 	clubID := vs[i]
	// 	if clubID != "location" {
	// 		clubIDs = append(clubIDs, clubID)
	// 	}
	// }

	// userInfo.clubIDs = clubIDs

	return userInfo
}
