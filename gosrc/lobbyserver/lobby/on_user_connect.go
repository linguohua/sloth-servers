package lobby

import (
	"gconst"

	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
)

func subscriberUserConnectEvent() {
	events := map[string]func(string, []byte){
		gconst.EventPlayerConnect: onConnect,
	}

	// cache.Sub(handlers)

	go func() {
		conn := pool.Get()
		defer conn.Close()

		channels := redis.Args{}
		for key := range events {
			channels = channels.Add(key)
		}

		psc := redis.PubSubConn{Conn: conn}
		psc.Subscribe(channels...)
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				log.Printf("sub %s message: %s\n", v.Channel, v.Data)
				handler := events[v.Channel]
				if handler != nil {
					value, exists := getPubData(string(v.Data))
					if exists {
						handler(string(v.Data), value)
					}
				}
			case redis.Subscription:
				log.Printf("sub %s %s %d\n", v.Channel, v.Kind, v.Count)
			case redis.PMessage:
				log.Printf("sub %s %s %s\n", v.Channel, v.Pattern, v.Data)
			case error:
				log.Panic(v)
			}
		}
	}()
}

// 用户连接
func onConnect(msg string, data []byte) {
	// playerID := string(data)
	// log.Println("onConnect, userID:", playerID)
	// //推送登录成功
	// if !pushLoginReply(playerID) {
	// 	log.Printf("push player %s info failed!", playerID)
	// 	return
	// }

	// //推送用户信息
	// if !pushPlayerInfo(playerID) {
	// 	log.Printf("push player %s info failed!", playerID)
	// 	return
	// }

	// //投递登录成功事件
	// pushEvent(gconst.EventPlayerLoginSuccess, playerID, playerID)

	//chost.presentNotifyFunc(playerID, true)
}

// PushLoginReply 发送登录成功
func pushLoginReply(playerID string) bool {
	// conn := pool.Get()
	// defer conn.Close()

	// token, _ := redis.String(conn.Do("HGET", gconst.AsUserTablePrefix+playerID, "Token"))
	// loginInfo := &MsgLoginReply{}
	// code := int32(LoginState_Success)
	// loginInfo.Result = &code
	// loginInfo.Token = &token
	// loginInfo.LastRoomInfo = loadLastRoomInfo(playerID)

	// loginbuffer, err := proto.Marshal(loginInfo)
	// if err != nil {
	// 	log.Printf("Login Marshal err:%v", err)
	// 	return false
	// }

	// push(int32(MessageCode_OPLoginReply), loginbuffer, playerID)
	// log.Printf("push player %s LoginReply", playerID)

	return true
}

// PushPlayerInfo 推送用户信息
func pushPlayerInfo(playerID string) bool {
	// TODO: llwant mysql
	// diamond, err := webdata.QueryDiamond(playerID)
	// if err != nil {
	// 	log.Printf("PushPlayerInfo QueryDiamond err:%v", err)
	// 	return false
	// }

	// gold, err := webdata.ReadGold(playerID)
	// if err != nil {
	// 	log.Printf("PushPlayerInfo ReadGold err:%v", err)
	// 	return false
	// }

	// conn := pool.Get()
	// defer conn.Close()

	// fields, err := redis.Strings(conn.Do("HMGET",gconst.AsUserTablePrefix+playerID, "Name", "Nick", "Sex", "Protrait", "Addr", "Token", "charm", "AvatarID", "isSyncCharm", "DanID"))
	// if err != nil {
	// 	log.Printf("pushPlayerInfo %s error:", err)
	// 	return false
	// }
	// var name = fields[0]
	// var nick = fields[1]
	// var sex = fields[2]
	// var protrait = fields[3]
	// var addr = fields[4]
	// var token = fields[5]
	// var charm = fields[6]

	// var dan = fields[9]
	// // var avatarID = fields[7]

	// sexInt64, _:= strconv.ParseInt(sex, 10, 32)
	// charmInt64, _:= strconv.ParseInt(charm, 10, 32)
	// danInt64, _:= strconv.ParseInt(dan, 10, 32)

	// info := &MsgUserInfo{}
	// info.Uid = &playerID
	// info.Name = &name
	// info.Sex = &sexInt64
	// info.Protrait = &protrait
	// info.Token = &token
	// info.Diamond = &diamond
	// info.Nick = &nick
	// info.Charm = &charmInt64
	// info.Addr = &addr
	// tmpdata := ""
	// info.Avatar = &tmpdata
	// info.Dan = &danInt64
	// info.Gold = &gold

	// if fields[7] != "" {
	// 	avatarBoxKey := gconst.SurrServerAvatarBox + playerID + ":" + fields[7]
	// 	expire, _ := redis.Int64(conn.Do("ttl", avatarBoxKey))
	// 	if expire > 0 {
	// 		info.Avatar = &fields[7]
	// 		log.Printf("user %s avatar box id %s.", playerID, *info.Avatar)
	// 	} else {
	// 		conn.Do("hdel", gconst.AsUserTablePrefix+playerID, "AvatarID")
	// 	}
	// }

	// // 是否已经同步魅力值到db的标志
	// var isSyncCharm = fields[8]
	// if isSyncCharm == "" {
	// 	_, err := webdata.UpdateCharmValue(playerID, 0, int(charmInt64), "同步数据")
	// 	if err == nil {
	// 		conn.Do("hset", gconst.AsUserTablePrefix+playerID, "isSyncCharm", "true")
	// 	}
	// }

	// buffer, err := proto.Marshal(info)
	// if err != nil {
	// 	log.Printf("player %s pushPlayerInfo Marshal err:%v", playerID, err)
	// 	return false
	// }

	// //推送用户信息
	// push(int32(MessageCode_OPSendUserInfo), buffer, playerID)

	// log.Printf("push player %s PlayerInfo", playerID)

	return true
}
