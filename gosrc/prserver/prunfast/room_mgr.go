package prunfast

import (
	"encoding/json"
	"fmt"
	"gconst"
	"gpubsub"
	"gscfg"
	"pokerface"
	"runtime/debug"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
)

var (
	pool                  *redis.Pool
	pubSubSequnce         = uint32(0)
	luaScriptForHandScore *redis.Script
)

// RoomMgr 房间管理
type RoomMgr struct {
	rooms map[string]*Room
	//rand  *rand.Rand
}

func (rm *RoomMgr) startup() {
	if gscfg.ServerID == "" {
		log.Panic("Must specify the server ID in config json")
		return
	}

	if gscfg.RoomServerID == "" {
		log.Panic("must spcify the RoomServerID in config json")
		return
	}

	// 初始化rooms
	rm.rooms = make(map[string]*Room)
	//rm.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	// 全局redis pool
	pool = newPool(gscfg.RedisServer)

	createLuaScript()

	// 往redis注册自己
	rm.register()

	rm.restoreRooms()

	// 新起一个goroutine去订阅redis
	gpubsub.Startup(pool, gscfg.ServerID, processNotifyMsg, processRequestMsg)

	go roomMonitor()
}

// newPool 新建redis连接池
func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

// processRequestMsg 处理redis publish过来的消息
func processRequestMsg(ssmsgBag *gconst.SSMsgBag) {
	var reqCode = ssmsgBag.GetRequestCode()
	var rm = roomMgr
	switch reqCode {
	case int32(gconst.SSMsgReqCode_CreateRoom):
		rm.onCreateRoomReq(ssmsgBag)
		break
	case int32(gconst.SSMsgReqCode_DeleteRoom):
		rm.onDeleteRoomReq(ssmsgBag)
		break
	case int32(gconst.SSMsgReqCode_UpdateLocation):
		rm.onUpdateLocation(ssmsgBag)
		break
	case int32(gconst.SSMsgReqCode_UpdatePropCfg):
		rm.onUpdatePropsCfg(ssmsgBag)
		break
	}
}

func processNotifyMsg(ssmsgBag *gconst.SSMsgBag) {
	log.Panic("processNotifyMsg game server not support notify type msg")
}

// onCreateRoomReq 房间管理服务器请求创建房间
func (rm *RoomMgr) onCreateRoomReq(msgBag *gconst.SSMsgBag) {
	log.Println("onCreateRoomReq begin-------------")
	roomParams := &gconst.SSMsgCreateRoom{}
	err := proto.Unmarshal(msgBag.Params, roomParams)

	if err != nil {
		log.Println(err)
		rm.replySSMsg(msgBag, gconst.SSMsgError_ErrDecode, nil)
		return
	}

	// 房间ID是否已经存在
	roomID := roomParams.GetRoomID()
	_, found := rm.rooms[roomID]
	if found {
		log.Println("room exists, id:", roomID)
		rm.replySSMsg(msgBag, gconst.SSMsgError_ErrRoomExist, nil)
		return
	}

	// 房间类型
	roomType := roomParams.GetRoomType()
	if roomType != int32(myRoomType) {
		log.Println("unsupport room tye:", roomType)
		rm.replySSMsg(msgBag, gconst.SSMsgError_ErrUnsupportRoomType, nil)
		return
	}

	roomConfigID := roomParams.GetRoomConfigID()
	log.Println("load room config, id:", roomConfigID)

	ownerID := roomParams.GetUserID()
	clubID := roomParams.GetClubID()

	// 获取redis链接，并退出函数时释放
	conn := pool.Get()
	defer conn.Close()

	cfgData, err := redis.Bytes(conn.Do("hget", gconst.LobbyRoomConfigTable, roomConfigID))
	if err != nil {
		log.Println("load room config failed:", err)
		rm.replySSMsg(msgBag, gconst.SSMsgError_ErrNoRoomConfig, nil)
		return
	}

	var configJSON = &RoomConfigJSON{}
	err = json.Unmarshal([]byte(cfgData), configJSON)
	if err != nil {
		log.Println("load room config failed:", err)
		rm.replySSMsg(msgBag, gconst.SSMsgError_ErrDecodeRoomConfig, nil)
		return
	}

	roomNumber := roomParams.GetRoomNumber()
	room := newRoomFromMgr(ownerID, clubID, roomID, roomConfigID, configJSON, roomNumber)
	rm.rooms[roomID] = room

	log.Printf("onCreateRoomReq success, UUID:%s, roomNumber:%s\n", roomID, room.roomNumber)

	// 创建房间完成
	rm.replySSMsg(msgBag, gconst.SSMsgError_ErrSuccess, nil)
}

// onDeleteRoomReq 房间管理服务器请求删除房间
func (rm *RoomMgr) onDeleteRoomReq(msgBag *gconst.SSMsgBag) {
	roomParams := &gconst.SSMsgDeleteRoom{}
	err := proto.Unmarshal(msgBag.Params, roomParams)

	if err != nil {
		log.Println(err)
		rm.replySSMsg(msgBag, gconst.SSMsgError_ErrDecode, nil)
		return
	}

	// 房间ID是否已经存在
	roomID := roomParams.GetRoomID()
	room, found := rm.rooms[roomID]
	if !found {
		log.Println("room not exists, id:", roomID)
		rm.replySSMsg(msgBag, gconst.SSMsgError_ErrRoomNotExist, nil)
		return
	}

	// 检查房间是否空闲
	if roomParams.GetOnlyEmpty() {
		if len(room.players) > 0 {
			log.Println("onDeleteRoomReq room is not empty, id:", roomID)
			rm.replySSMsg(msgBag, gconst.SSMsgError_ErrRoomIsNoEmpty, nil)
			return
		}
	}

	why := roomParams.GetWhy()
	log.Printf("onDeleteRoomReq, room id:%s, room number:%s, why:%d\n", roomID,
		room.roomNumber, why)

	// 这里不能使用forceDisbandRoom，因为他会反过来发送请求给房间管理服务器
	// 房间销毁，如果玩家在上面，则给玩家发送通知
	room.deleteReason = pokerface.RoomDeleteReason_DisbandByOwnerFromRMS
	delete(rm.rooms, roomID)
	room.destroy()

	// 完成删除房间
	rm.replySSMsg(msgBag, gconst.SSMsgError_ErrSuccess, nil)
}

// onUpdateLocation 管理服务器更新用户GPS信息
func (rm *RoomMgr) onUpdateLocation(msgBag *gconst.SSMsgBag) {
	log.Printf("onUpdateLocation")
	updateLocation := &gconst.SSMsgUpdateLocation{}
	err := proto.Unmarshal(msgBag.Params, updateLocation)
	if err != nil {
		log.Println("onUpdateLocation error:", err)
		return
	}

	userID := updateLocation.GetUserID()
	userMapItem, ok := usersMap[userID]
	if !ok {
		log.Printf("onUpdateLocation, user %s not online:", userID)
		return
	}

	user := userMapItem.user
	if user == nil {
		log.Printf("onUpdateLocation, user is nil ")
		return
	}

	room := user.getRoom()
	if room == nil {
		log.Printf("onUpdateLocation, can't get user room")
		return
	}

	room.updateUserLocation(userID, updateLocation.GetLocation())
}

// onUpdatePropCfg 更新牌局内的道具配置
func (rm *RoomMgr) onUpdatePropsCfg(msgBag *gconst.SSMsgBag) {
	log.Printf("onUpdatePropCfg")

	var cfgString = string(msgBag.Params)
	if cfgString == "" {
		log.Println(" cfgString is emtpy")
		return
	}

	msg := &pokerface.MsgUpdatePropCfg{}
	msg.PropCfg = &cfgString
	buf := formatGameMsg(msg, int32(pokerface.MessageCode_OPUpdatePropCfg))

	// 给所有用户发送道具配置
	for _, userMapItem := range usersMap {
		user := userMapItem.user
		if user != nil {
			user.send(buf)
		}
	}
}

// replySSMsg 给其他服务器回复请求完成
func (rm *RoomMgr) replySSMsg(msgBag *gconst.SSMsgBag, errCode gconst.SSMsgError, params []byte) {
	if msgBag.GetSourceURL() == "" {
		log.Println("replySSMsgError, no source URL")
		return
	}

	replyMsgBag := &gconst.SSMsgBag{}
	var msgType32 = int32(gconst.SSMsgType_Response)
	replyMsgBag.MsgType = &msgType32
	var seqNO32 = msgBag.GetSeqNO()
	replyMsgBag.SeqNO = &seqNO32
	var requestCode32 = msgBag.GetRequestCode()
	replyMsgBag.RequestCode = &requestCode32
	var status32 = int32(errCode)
	replyMsgBag.Status = &status32

	if params != nil {
		replyMsgBag.Params = params
	}

	gpubsub.PublishMsg(msgBag.GetSourceURL(), replyMsgBag)
}

// getRoom 获取房间
func (rm *RoomMgr) getRoom(ID string) *Room {
	room, ok := rm.rooms[ID]
	if ok {
		return room
	}
	return nil
}

func (rm *RoomMgr) getRoomByNumber(number string) *Room {
	for _, r := range rm.rooms {
		if r.roomNumber == number {
			return r
		}
	}

	return nil
}

// register 往redis上登记自己
func (rm *RoomMgr) register() {
	// 获取redis链接，并退出函数时释放
	conn := pool.Get()
	defer conn.Close()

	if rm.serverIDSubscriberExist(conn) {
		log.Panicln("The same UUID server instance exists, failed to startup, server ID:", gscfg.ServerID)
		return
	}

	hashKey := gconst.GameServerInstancePrefix + gscfg.ServerID
	conn.Send("MULTI")
	conn.Send("hmset", hashKey, "roomtype", int(myRoomType), "ver", versionCode, "p", gscfg.ServerPort)
	conn.Send("SADD", fmt.Sprintf("%s%d", gconst.GameServerInstancePrefix, int(myRoomType)), gscfg.ServerID)
	conn.Send("SADD", gconst.GameServerRoomTypeSet, int(myRoomType))
	// conn.Send("HSET", fmt.Sprintf("%s%d", gconst.RoomTypeKey, myRoomType), "type", 1)
	_, err := conn.Do("EXEC")
	if err != nil {
		log.Panicln("failed to register server to redis:", err)
	}

	if gscfg.EtcdServer != "" {
		// 如果服务器使用etcd，则需要去etcd登记自己
		gscfg.Regist2Etcd(versionCode, int(myRoomType))
	}
}

func (rm *RoomMgr) serverIDSubscriberExist(conn redis.Conn) bool {
	subCounts, err := redis.Int64Map(conn.Do("PUBSUB", "NUMSUB", gscfg.ServerID))
	if err != nil {
		log.Println("warning: serverIDSubscriberExist, redis err:", err)
	}

	count, _ := subCounts[gscfg.ServerID]
	if count > 0 {
		return true
	}

	return false
}

// 服务器重启后恢复房间
func (rm *RoomMgr) restoreRooms() {
	// 获取redis链接，并退出函数时释放
	conn := pool.Get()
	defer conn.Close()

	// 清空在线玩家数量
	var key = fmt.Sprintf("%s%d", gconst.GameServerOnlineUserNumPrefix, myRoomType)
	conn.Do("HSET", key, gscfg.ServerID, 0)

	// _, err := luaScriptForHandScore.Do(conn, gconst.LobbyPlayerTablePrefix+"1", -450)
	// if err != nil {
	// 	log.Println("luaScriptForHandScore err:", err)
	// }

	// 获取所有的房间ID
	roomIDs, err := redis.Strings(conn.Do("SMEMBERS", gconst.LobbyRoomTableSet))
	if err != nil {
		log.Println("restoreRooms, err:", err)
		return
	}

	log.Println("try to restore room, count:", len(roomIDs))

	if len(roomIDs) < 1 {
		return
	}

	conn.Send("MULTI")
	for _, roomID := range roomIDs {
		conn.Send("hmget", gconst.LobbyRoomTablePrefix+roomID, "ownerID", "roomConfigID", "gameServerID", "roomNumber", "clubID")
	}

	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		log.Println("restoreRooms err: ", err)
		return
	}

	type RoomInfo struct {
		ownerID      string
		roomConfigID string
		gameServerID string
		roomID       string
		roomNumber   string
		clubID       string
	}

	roomInfos := make([]*RoomInfo, 0, len(values))
	for i := 0; i < len(values); i++ {
		fields, err := redis.Strings(values[i], nil)
		if err != nil {
			continue
		}

		gameServerID := fields[2]
		if gameServerID != gscfg.ServerID {
			continue
		}

		roomInfo := &RoomInfo{}
		roomInfo.ownerID = fields[0]
		roomInfo.roomConfigID = fields[1]
		roomInfo.roomNumber = fields[3]
		roomInfo.clubID = fields[4]
		roomInfo.gameServerID = gameServerID
		roomInfo.roomID = roomIDs[i]
		roomInfos = append(roomInfos, roomInfo)
	}

	// 加载配置
	vs, err := redis.Values(conn.Do("hgetall", gconst.LobbyRoomConfigTable))
	if err != nil {
		log.Println("restoreRooms, get roomConfig err:", err)
		return
	}

	configs := make(map[string]*RoomConfigJSON)
	for i := 0; i < len(vs); i = i + 2 {
		roomConfigID, err := redis.String(vs[i], nil)
		buf, err := redis.Bytes(vs[i+1], nil)

		roomConfig := &RoomConfigJSON{}
		err = json.Unmarshal(buf, roomConfig)
		if err != nil {
			continue
		}

		configs[roomConfigID] = roomConfig
	}

	// 创建房间
	for _, roomInfo := range roomInfos {
		roomConfig, ok := configs[roomInfo.roomConfigID]
		if !ok {
			continue
		}

		roomID := roomInfo.roomID
		ownerID := roomInfo.ownerID
		// 暂时去掉,避免牌友群进不来
		// clubID := roomInfo.clubID
		clubID := ""

		room := newRoomFromMgr(ownerID, clubID, roomID, roomInfo.roomConfigID,
			roomConfig, roomInfo.roomNumber)

		room.readHandInfoFromRedis4Restore(conn)
		room.restorePlayersWhen()
		// 房间恢复，需要重置房间玩家信息，以及发送通知给房间管理服务器
		room.writeOnlinePlayerList2Redis(conn)

		log.Println("room restore ok, room ID:", roomID)

		rm.rooms[roomID] = room
	}
}

func roomMonitor() {
	defer func() {
		if r := recover(); r != nil {
			roomExceptionCount++
			debug.PrintStack()
			log.Printf("-----This RoomMonitor GR will die, Recovered in roomMonitor:%v\n", r)
		}
	}()

	for {
		// 每间隔10分钟唤醒一次，唤醒后检查每一个房间
		// 最后一个消息的接收时间
		time.Sleep(10 * time.Minute)
		now := time.Now()
		// 如果时间大于6小时,则认为房间空闲过久，直接关闭
		for _, r := range roomMgr.rooms {
			diff := now.Sub(r.lastReceivedTime)
			if diff > 6*time.Hour {
				log.Printf("room %s, owner:%s, has no message in pass 6 hours, disband it\n", r.ID, r.ownerID)
				// 执行房间解散流程
				roomMgr.forceDisbandRoom(r, pokerface.RoomDeleteReason_IdleTimeout)
			} else {
				// 处理黑名单超时
				if r.rbl != nil {
					r.rbl.unblockWhenTimePassed()
				}
			}
		}
	}
}

func (rm *RoomMgr) forceDisbandRoom(r *Room, reason pokerface.RoomDeleteReason) {
	r.deleteReason = reason
	r.forceDisband()
	delete(roomMgr.rooms, r.ID)
}

// pushNotify2RoomServer 给房间管理服务器发送通知
func pushNotify2RoomServer(reqCode gconst.SSMsgReqCode, params proto.Message) {
	var msg = &gconst.SSMsgBag{}
	var msgType32 = int32(gconst.SSMsgType_Notify)
	msg.MsgType = &msgType32
	var seqNo = uint32(0)
	msg.SeqNO = &seqNo
	var requestCode32 = int32(reqCode)
	msg.RequestCode = &requestCode32
	var status32 = int32(0)
	msg.Status = &status32

	if params != nil {
		bytes, err := proto.Marshal(params)

		if err != nil {
			log.Panic("marshal params failed:", err)
			return
		}
		msg.Params = bytes
	}

	gpubsub.PublishMsg(gscfg.RoomServerID, msg)
}

//lua脚本
func createLuaScript() {
	script2 := `local value = redis.call('hget', KEYS[1], 'dfHMW')
		if type(value) ~= 'string' then
			value = 0
		else
			value = tonumber(value)
		end
		if value < tonumber(KEYS[2]) then
			redis.call('HSET', KEYS[1], 'dfHMW', KEYS[2])
		end
		value = redis.call('hget', KEYS[1], 'dfHML')
		if type(value) ~= 'string' then
			value = 0
		else
			value = tonumber(value)
		end
		if value > tonumber(KEYS[2]) then
			redis.call('HSET', KEYS[1], 'dfHML', KEYS[2])
		end
		return 0`

	luaScriptForHandScore = redis.NewScript(2, script2)
}
