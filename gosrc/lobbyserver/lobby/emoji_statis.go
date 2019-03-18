package lobby

import (
	"fmt"
	"gconst"
	"time"

	log "github.com/sirupsen/logrus"
)

func emojiCount(roomType int, emojiIndex int) {
	conn := pool.Get()
	defer conn.Close()

	dd := time.Now().Format("20060102")

	var key = fmt.Sprintf(gconst.EmojiDailyStatisTablePrefix, dd, roomType)

	// key 30后过期
	time2Delete := time.Now().AddDate(0, 0, 1)
	time2Delete = time2Delete.Truncate(time.Hour * 24)
	time2Delete = time2Delete.Add(24 * 30 * time.Hour)

	conn.Send("MULTI")
	conn.Send("INCR", key)
	conn.Send("EXPIREAT", key, time2Delete.Unix())

	_, err := conn.Do("EXEC")
	if err != nil {
		log.Println("err:", err)
	}
}
