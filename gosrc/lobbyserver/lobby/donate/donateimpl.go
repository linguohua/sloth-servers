package donate

// "webdata"
import (
	"fmt"
	"gconst"
	"lobbyserver/lobby"

	"github.com/garyburd/redigo/redis"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

func generateDonateUUID() string {
	uid, _ := uuid.NewV4()
	donateID := fmt.Sprintf("%s", uid)
	return donateID
}

func getPropByRoomType(roomType int, propsType int) *Prop {
	propCfgMap := clientPropCfgsMap[roomType]
	if propCfgMap == nil {
		return nil
	}

	prop := propCfgMap[propsType]

	return prop

}

// TODO: 需要检查用户是否有道具，如果有道具，从道具那里消耗，不扣钻
func donate(propsType uint32, from string, to string, roomType int) (result *gconst.SSMsgDonateRsp, errCode int32) {
	log.Printf("donate, propsType:%d, from:%s, to:%s, roomType:%d", propsType, from, to, roomType)

	var prop = getPropByRoomType(roomType, int(propsType))
	if prop == nil {
		var errMsg = fmt.Sprintf("RoomType:%d not exist propsType:%d", roomType, propsType)
		log.Panicln(errMsg)
		return
	}

	var propID = uint32(prop.PropID)
	if isUserHaveProp(propID, from) {
		rsp, errCode := consumeUserProp(prop, from, to, roomType)
		if errCode != int32(gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough) {
			return rsp, errCode
		}

		log.Printf("User %s Prop %d number in redis not same as database", from, propID)
	}

	var costDiamond = prop.Diamond
	var charm = prop.Charm
	var cost = int64(costDiamond)

	mySQLUtil := lobby.MySQLUtil()
	remainDiamond, errCode := mySQLUtil.UpdateDiamond(from, -cost)
	// TODO: 在lobby中写成常量
	if errCode == 2 {
		log.Error("donate error, diamond not enough, remainDiamond:", remainDiamond)
		errCode = int32(gconst.SSMsgError_ErrTakeoffDiamondFailedNotEnough)
		return nil, errCode
	} else if errCode != 0 {
		log.Errorf("donate, UpdateDiamond unknow error code:%d", errCode)
	}

	conn := lobby.Pool().Get()
	defer conn.Close()

	originCharm, err := redis.Int(conn.Do("HGET", gconst.LobbyUserTablePrefix+to, "charm"))
	if err != nil && err != redis.ErrNil {
		var errMsg = fmt.Sprintf("load user %s charm failed, err:%v", to, err)
		log.Panicln(errMsg)
	}

	var newCharm = originCharm + charm

	conn.Send("MULTI")
	conn.Send("HSET", gconst.LobbyUserTablePrefix+from, "diamond", remainDiamond)
	conn.Send("HSET", gconst.LobbyUserTablePrefix+to, "charm", newCharm)

	_, err = conn.Do("EXEC")
	if err != nil {
		log.Panicln("save donate err:", err)
	}

	// TODO: 在数据库中添加魅力值

	var msgDonateRsp = &gconst.SSMsgDonateRsp{}
	var int32Diamond = int32(remainDiamond)
	msgDonateRsp.Diamond = &int32Diamond
	var int32Charm = int32(newCharm)
	msgDonateRsp.Charm = &int32Charm

	return msgDonateRsp, int32(gconst.SSMsgError_ErrSuccess)
}

func isUserHaveProp(propID uint32, userID string) bool {
	log.Printf("isUserHaveProp, propID:%d, userID:%s", propID, userID)
	conn := lobby.Pool().Get()
	defer conn.Close()

	propNum, err := redis.Int(conn.Do("HGET", gconst.LobbyUserDonatePrefix+userID, propID))
	if err != nil {
		log.Println("getUserPropNum error:", err)
		return false
	}

	if propNum > 0 {
		return true
	}
	return false
}

func consumeUserProp(prop *Prop, from string, to string, roomType int) (result *gconst.SSMsgDonateRsp, errCode int32) {
	log.Printf("consumeUserProp,propID:%d,from:%s, to:%s, roomType:%d", prop.PropID, from, to, roomType)

	var charm = prop.Charm
	// var propID = prop.PropID

	// TODO: 在数据库中扣取道具

	conn := lobby.Pool().Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Do("HGET", gconst.LobbyUserTablePrefix+to, "charm")
	conn.Do("HGET", gconst.LobbyUserTablePrefix+from, "diamond")
	vs, err := redis.Ints(conn.Do("EXEC"))
	if err != nil {
		log.Println("consumeUserProp get charm and diamond error:", err)
	}

	var originCharm = vs[0]
	var diamond = vs[1]
	var newCharm = originCharm + charm

	_, err = conn.Do("HSET", gconst.LobbyUserTablePrefix+to, "charm", newCharm)
	if err != nil {
		log.Panicln("save donate err:", err)
	}

	// TODO: 在数据中添加魅力值

	var msgDonateRsp = &gconst.SSMsgDonateRsp{}
	var int32Diamond = int32(diamond)
	msgDonateRsp.Diamond = &int32Diamond
	var int32Charm = int32(newCharm)
	msgDonateRsp.Charm = &int32Charm

	return msgDonateRsp, int32(gconst.SSMsgError_ErrSuccess)
}
