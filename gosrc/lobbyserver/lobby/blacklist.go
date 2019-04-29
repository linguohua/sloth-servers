package lobby

import (
	"fmt"
	"gconst"

	log "github.com/sirupsen/logrus"

	"github.com/garyburd/redigo/redis"
)

func isUserInBlacklist(userID string) bool {

	conn := pool.Get()
	defer conn.Close()

	exist, _ := redis.Int(conn.Do("SISMEMBER", gconst.LobbyUserBlacklistSet, userID))
	if exist == 1 {
		return true
	}

	return false

}

func addUser2Blacklist(userID string) error {
	conn := pool.Get()
	defer conn.Close()

	exist, _ := redis.Int(conn.Do("EXISTS", gconst.LobbyUserTablePrefix+userID))
	if exist == 0 {
		return fmt.Errorf("User %s not exist", userID)
	}

	_, err := conn.Do("SADD", gconst.LobbyUserBlacklistSet, userID)
	if err != nil {
		log.Println("addUser2Blacklist err:", err)
		return fmt.Errorf("redis error %v", err)
	}

	return nil

}

func removeUserFromBlacklist(userID string) error {
	conn := pool.Get()
	defer conn.Close()

	exist, _ := redis.Int(conn.Do("SISMEMBER", gconst.LobbyUserBlacklistSet, userID))
	if exist == 0 {
		return fmt.Errorf("User %s not in blacklist", userID)
	}

	_, err := conn.Do("SREM", gconst.LobbyUserBlacklistSet, userID)
	if err != nil {
		log.Println("removeUserFromBlacklist err:", err)
		return fmt.Errorf("redis error %v", err)
	}

	return nil
}

func loadBlacklist() []string {
	conn := pool.Get()
	defer conn.Close()

	userIDs, err := redis.Strings(conn.Do("SMEMBERS", gconst.LobbyUserBlacklistSet))
	if err != nil {
		log.Println("removeUserFromBlacklist err:", err)
	}

	return userIDs
}
