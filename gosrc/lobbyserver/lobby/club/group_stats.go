package club

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gconst"
	"time"
	"encoding/json"
)

func statsGroupBigWiner(groupID string, roomType string, roomID string, playerStats []*gconst.SSMsgGameOverPlayerStat, finishHand int32) {
	buf, _ := json.Marshal(playerStats)

	log.Printf("statsGroupBigWiner, groupID:%s, roomType:%s, roomID:%s, playerStats:%s",groupID, roomType, roomID, string(buf))
	log.Println("statsGroupBigWiner")
	if playerStats == nil {
		log.Println("playerStats == nil")
		return
	}

	if len(playerStats) == 0 {
		log.Println("playerStats == nil")
		return
	}

	buf ,er := json.Marshal(playerStats)
	if er != nil {
		log.Println("statsGroupBigWiner, error:", er)
	}

	log.Println("statsGroupBigWiner, buf:", string(buf))

	var bigScore = getBigScore(playerStats)
	var isAllScoreSame = isScoreAllSame(playerStats)

	// 慌庄，分数为0的过滤掉
	if isAllScoreSame && finishHand == 0 {
		log.Println("isAllScoreSame && finishHand == 0")
		return
	}
	// if bigScore == 0 {
	// 	log.Println("bigScore == 0")
	// 	return
	// }

	t := time.Now()
	nowTimeUnix := t.Unix()
	// nowTime := t.Format("20060102150405")
	dd := t.Format("20060102")
	var keyGroupBigWinnerStats = fmt.Sprintf(gconst.GroupBigWinnerStats, groupID, dd)
	var keyGroupBigHandStats = fmt.Sprintf(gconst.GroupBigHandStats, groupID, dd)
	var keyGroupSpecificRoomBigWinnerStats = fmt.Sprintf(gconst.GroupSpecificRoomBigWinnerStats, groupID, roomType, dd)
	var keyGroupSpecificRoomBigHandStats = fmt.Sprintf(gconst.GroupSpecificRoomBigHandStats, groupID, roomType, dd)
	var keyGroupStatsUpdateTime = fmt.Sprintf(gconst.GroupStatsUpdateTime, groupID, dd)

	tOfTomorrow5 := time.Now().AddDate(0, 0, 5)
	tOfTomorrow5 = time.Date(tOfTomorrow5.Year(), tOfTomorrow5.Month(), tOfTomorrow5.Day(),
		0, 0, 0, 0, time.Local)
	tOfTomorrow5 = tOfTomorrow5.Add(5 * time.Hour)
	expireat := tOfTomorrow5.Unix()

	conn := pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	for _, playerStat := range playerStats {
		if playerStat.GetScore() == bigScore && !isAllScoreSame {
			conn.Send("ZINCRBY", keyGroupBigWinnerStats, 1, playerStat.GetUserID())
			conn.Send("ZINCRBY", keyGroupSpecificRoomBigWinnerStats, 1, playerStat.GetUserID())
		}

		conn.Send("ZINCRBY", keyGroupBigHandStats, 1, playerStat.GetUserID())
		conn.Send("ZINCRBY", keyGroupSpecificRoomBigHandStats, 1, playerStat.GetUserID())
		conn.Send("HSET", keyGroupStatsUpdateTime, playerStat.GetUserID(), nowTimeUnix)
	}

	conn.Send("EXPIREAT", keyGroupBigWinnerStats, expireat)
	conn.Send("EXPIREAT", keyGroupBigHandStats, expireat)
	conn.Send("EXPIREAT", keyGroupSpecificRoomBigWinnerStats, expireat)
	conn.Send("EXPIREAT", keyGroupSpecificRoomBigHandStats, expireat)
	conn.Send("EXPIREAT", keyGroupStatsUpdateTime, expireat)

	_, err := conn.Do("EXEC")
	if err != nil {
		log.Println("statsGroupBigWiner err:", err)
		return
	}


}

func getBigScore(playerStats []*gconst.SSMsgGameOverPlayerStat) int32 {
	var bigScore = int32(0)
	for _, playerStat := range playerStats {
		if playerStat.GetScore() >= bigScore {
			bigScore = playerStat.GetScore()
		}
	}

	return bigScore;
}

func isScoreAllSame(playerStats []*gconst.SSMsgGameOverPlayerStat) bool{
	if len(playerStats) == 0{
		return true
	}

	score := playerStats[0].GetScore()
	for _, playerStat := range playerStats {
		if playerStat.GetScore() != score {
			return false
		}
	}

	return true
}