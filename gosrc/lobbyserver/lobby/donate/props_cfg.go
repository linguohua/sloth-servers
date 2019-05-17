package donate

import (
	"encoding/json"
	"fmt"
	"gconst"
	"gpubsub"
	"lobbyserver/config"
	"lobbyserver/lobby"

	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
)

const (
	// 道具在房间里面对应的序列
	// 1：鲜花    2：啤酒    3：鸡蛋    4：拖鞋
	// 8：献吻    7：红酒    6：大便    5：拳头
	flower   = 1
	beer     = 2
	egg      = 3
	slippers = 4
	punch    = 5
	faeces   = 6
	redWine  = 7
	kiss     = 8

	// 道具对应在道具表中的ID
	flowerID   = 1
	beerID     = 2
	eggID      = 3
	slippersID = 4
	punchID    = 5
	faecesID   = 6
	redWineID  = 7
	kissID     = 8
)

// propsDiamond[flower] = 1
// propsDiamond[beer] = 3
// propsDiamond[egg] = 1
// propsDiamond[slippers] = 3
// propsDiamond[punch] = 10
// propsDiamond[faeces] = 5
// propsDiamond[redWine] = 10
// propsDiamond[kiss] = 5
// propsCharm[flower] = 1
// propsCharm[beer] = 3
// propsCharm[egg] = -1
// propsCharm[slippers] = -3
// propsCharm[punch] = -10
// propsCharm[faeces] = -5
// propsCharm[redWine] = 10
// propsCharm[kiss] = 5

// PropCfgMap 道具配置映射表
type PropCfgMap map[int]*Prop

var (
	// ClientPropCfgsMap key index, value Prop
	// 下发给客户端显示
	clientPropCfgsMap = make(map[int]PropCfgMap)
	// key propID, value Prop
	// 服务器扣钻加魅力值用
	serverPropCfgsMap = make(map[int]PropCfgMap)

	// props = make(map[uint32]Prop)
)

// Prop 道具
type Prop struct {
	// 对应道具表中的ID
	PropID int `json:"propID"`
	// 消耗一个道具所需要的钻石
	Diamond int `json:"diamond"`
	// 消耗一个道具，对方增加的魅力值
	Charm int `json:"charm"`
	// // 当前拥有的此道具的数量
	// Have int `json:"have"`
}

func initGamePropCfgs() {
	loadAllRoomPropCfgs()
}

func defaultPropsCfg() string {
	// 以第一项为例：
	// "1"表示道具卡槽位置
	// "propID"表示道具ID
	// "diamon"表示消耗道具花费的钻石
	// "charm"表示消耗道具后，对方获取的魅力值
	//
	var propsCfg = `{
		"1":{
			"propID":1002,
			"diamond":1,
			"charm":1
		},
		"2":{
			"propID":1003,
			"diamond":3,
			"charm":3
		},
		"3":{
			"propID":1004,
			"diamond":1,
			"charm":-1
		},
		"4":{
			"propID":1005,
			"diamond":3,
			"charm":-3
		},
		"5":{
			"propID":1009,
			"diamond":10,
			"charm":-10
		},
		"6":{
			"propID":1008,
			"diamond":5,
			"charm":-5
		},
		"7":{
			"propID":1007,
			"diamond":10,
			"charm":10
		},
		"8":{
			"propID":1006,
			"diamond":5,
			"charm":5
		}
	}`

	return propsCfg
}

func parsePropsCfgJSON(roomType int, confgJSON string) {
	var propMap = make(map[int]interface{})
	err := json.Unmarshal([]byte(confgJSON), &propMap)
	if err != nil {
		log.Println("parsePropsCfgJSON, Unmarshal error:", err)
		return
	}

	// 以道具槽位ID为key
	var clientPropCfgMap = make(map[int]*Prop)
	// 以道具ID为key
	var serverPropCfgMap = make(map[int]*Prop)

	for k, v := range propMap {
		buf, err := json.Marshal(v)
		if err != nil {
			log.Println("parsePropsCfgJSON, Marshal error:", err)
			continue
		}

		var prop = Prop{}
		err = json.Unmarshal(buf, &prop)
		if err != nil {
			log.Println("parsePropsCfgJSON, Unmarshal error:", err)
			continue
		}

		clientPropCfgMap[k] = &prop
		serverPropCfgMap[prop.PropID] = &prop
	}

	clientPropCfgsMap[roomType] = clientPropCfgMap
	serverPropCfgsMap[roomType] = serverPropCfgMap
}

func loadAllRoomPropCfgs() {
	log.Println("loadAllPropCfg")
	conn := lobby.Pool().Get()
	defer conn.Close()

	gameRoomTypes, err := redis.Ints(conn.Do("SMEMBERS", gconst.GameServerRoomTypeSet))
	if err != nil {
		log.Println("loadAllRoomPropCfgs, err:", err)
		return
	}

	roomTypes := make([]int32, 0, len(gconst.RoomType_value))

	conn.Send("MULTI")
	for _, roomType := range gconst.RoomType_value {
		conn.Do("HGET", gconst.LobbyPropsCfgTable, roomType)
		roomTypes = append(roomTypes, roomType)
	}

	for _, roomType := range gameRoomTypes {
		roomTypeInt32 := int32(roomType)
		_, ok := gconst.RoomType_name[roomTypeInt32]
		if !ok {
			conn.Do("HGET", gconst.LobbyPropsCfgTable, roomType)
			roomTypes = append(roomTypes, roomTypeInt32)
		}
	}

	cfgStrings, err := redis.Strings(conn.Do("EXEC"))
	if err != nil {
		log.Println("loadAllRoomPropCfgs error:", err)
	}

	var fromRedis = 0
	for index, cfgString := range cfgStrings {
		var roomType = int(roomTypes[index])
		if cfgString != "" {
			fromRedis++
			parsePropsCfgJSON(roomType, cfgString)
		} else {
			var defaultCfg = defaultPropsCfg()
			parsePropsCfgJSON(roomType, defaultCfg)
		}
	}

	log.Printf("loadAllRoomPropCfgs, from redis:%d, default:%d", fromRedis, len(clientPropCfgsMap)-fromRedis)
}

// GetAllRoomPropCfgs 导出给web
func GetAllRoomPropCfgs() interface{} {
	return clientPropCfgsMap
}

// PP 道具属性，db那边预先定义好
type PP struct {
	CoinTypeID int    `json:"CoinTypeID"`
	CoinName   string `json:"CoinName"`
	Remark     string `json:"Remark"`
}

func loadPropTable(conn redis.Conn) map[int]PP {

	ppMap := make(map[int]PP)
	// 拉取配置表
	ppString, err := redis.String(conn.Do("GET", gconst.LobbyPropsTable))
	if err != nil {
		log.Println("err:", err)
		return ppMap
	}
	log.Println("loadPropTable, ppString:", ppString)

	pps := make([]PP, 0)
	err = json.Unmarshal([]byte(ppString), &pps)
	if err != nil {
		log.Println("err:", err)
		return ppMap
	}

	for _, pp := range pps {
		ppMap[pp.CoinTypeID] = pp
	}
	return ppMap
}

// UpdateRoomPropsCfg web更新道具配置
func UpdateRoomPropsCfg(JSONString string) error {
	if JSONString == "" {
		return fmt.Errorf("Config is emtpy")
	}

	type Cfg struct {
		PropID  int `json:"propID"`
		Diamond int `json:"diamond"`
		Charm   int `json:"charm"`
		Index   int `json:"id"`
	}

	type PropsCfg struct {
		RoomType int   `json:"roomType"`
		Cfgs     []Cfg `json:"propsCfg"`
	}

	var propsCfg = &PropsCfg{}
	err := json.Unmarshal([]byte(JSONString), propsCfg)
	if err != nil {
		log.Println("UpdateRoomPropsCfg, Unmarshal:", err)
		return err
	}

	log.Println("PropsCfg:")
	log.Println(propsCfg)

	if propsCfg.RoomType == 0 {
		return fmt.Errorf("Room type can't be 0")
	}

	if len(propsCfg.Cfgs) == 0 {
		return fmt.Errorf("Room type can't be 0")
	}

	// 拉取配置表
	conn := lobby.Pool().Get()
	defer conn.Close()

	ppMap := loadPropTable(conn)

	var notExistPropID = 0
	var repeatPropID = 0
	var isPropIDNotExist = false
	var isPropIDRepeat = false
	propCfgMap := make(PropCfgMap)
	checkPropRepeat := make(map[int]Cfg)
	for _, propCfg := range propsCfg.Cfgs {
		// 检查这个道具是否已经在数据库中配置
		_, ok := ppMap[propCfg.PropID]
		if !ok {
			notExistPropID = propCfg.PropID
			isPropIDNotExist = true
			break
		}

		// 检查这个道具是否重复配置
		_, ok = checkPropRepeat[propCfg.PropID]
		if ok {
			repeatPropID = propCfg.PropID
			isPropIDRepeat = true
			break
		}
		checkPropRepeat[propCfg.PropID] = propCfg

		var prop = &Prop{}
		prop.Charm = propCfg.Charm
		prop.Diamond = propCfg.Diamond
		prop.PropID = propCfg.PropID
		propCfgMap[propCfg.Index] = prop
	}

	if isPropIDNotExist {
		return fmt.Errorf("Prop %d not exist", notExistPropID)
	}

	if isPropIDRepeat {
		return fmt.Errorf("Prop %d repeat config", repeatPropID)
	}

	buf, err := json.Marshal(propCfgMap)
	if err != nil {
		return err
	}

	clientPropCfgsMap[propsCfg.RoomType] = propCfgMap

	key := fmt.Sprintf("%s%d", gconst.GameServerInstancePrefix, propsCfg.RoomType)
	conn.Send("MULTI")
	conn.Send("HSET", gconst.LobbyPropsCfgTable, propsCfg.RoomType, string(buf))
	conn.Send("SMEMBERS", key)

	vs, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		return err
	}

	serverIDs, err := redis.Strings(vs[1], nil)
	if err != nil {
		return err
	}

	for _, serverID := range serverIDs {
		sendPropCfg2GameServer(string(buf), serverID)
	}
	// log.Println(string(buf))

	return nil
}

func sendPropCfg2GameServer(propCfgString string, serverID string) {
	log.Printf("sendPropCfg2GameServer,serverID:%s, propCfgString:%s", serverID, propCfgString)

	msgType := int32(gconst.SSMsgType_Request)
	requestCode := int32(gconst.SSMsgReqCode_UpdatePropCfg)
	status := int32(gconst.SSMsgError_ErrSuccess)

	msgBag := &gconst.SSMsgBag{}
	msgBag.MsgType = &msgType
	var sn = lobby.GenerateSn()
	msgBag.SeqNO = &sn
	msgBag.RequestCode = &requestCode
	msgBag.Status = &status
	var url = config.ServerID
	msgBag.SourceURL = &url
	msgBag.Params = []byte(propCfgString)

	gpubsub.PublishMsg(serverID, msgBag)
}
