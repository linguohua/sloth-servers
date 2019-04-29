package prunfast

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokerface"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/DisposaBoy/JsonConfigReader"
)

const (
	monkeyRoomName = "monkey-room"
)

// MonkeyMgr 测试用
type MonkeyMgr struct {
	// currentCfg *MonkeyCfg
	room       *Room
	cfgs       []*MonkeyCfg
	roomConfig *RoomConfig
}

// start 启动monkeyMgr
func (mgr *MonkeyMgr) start() {

}

// createMonkeyRoom 通过mjtest测试工具创建monkey房间，monkey房间是一个特殊的房间，专门用于测试
func (mgr *MonkeyMgr) createMonkeyRoom(w http.ResponseWriter, r *http.Request) {

	if mgr.room != nil {
		w.WriteHeader(404)
		w.Write([]byte("monkey room exists"))
		return
	}

	var roomConfig = newRoomConfig()

	query := r.URL.Query()
	roomID := query.Get("roomID")
	if roomID == "" {
		roomID = monkeyRoomName
	}

	_, exist := roomMgr.rooms[roomID]
	if exist {
		w.WriteHeader(404)
		w.Write([]byte(fmt.Sprintf("monkey room %s exists", roomID)))
		return
	}

	// 创建一个monkey房间
	room := newRoomForMonkey("1", roomID, roomConfig)
	if roomID == monkeyRoomName {
		mgr.room = room
	}

	log.Printf("create monkey room with ID:%s\n", roomID)
	roomMgr.rooms[roomID] = room
}

// destroyMonkeyRoom 通过mjtest销毁房间
func (mgr *MonkeyMgr) destroyMonkeyRoom(w http.ResponseWriter, r *http.Request) {
	if mgr.room == nil {
		w.WriteHeader(404)
		w.Write([]byte("no monkey room"))
		return
	}

	_, ok := roomMgr.rooms[monkeyRoomName]
	if ok {
		delete(roomMgr.rooms, monkeyRoomName)
		mgr.room.destroy()
	}

	mgr.room = nil

	w.Write([]byte("delete OK"))
}

// attachRoomCfg2Room 通过mjtest测试工具，附加房间配置到某一个房间
func (mgr *MonkeyMgr) attachRoomCfg2Room(w http.ResponseWriter, r *http.Request, body string) {
	var roomNumber = r.URL.Query().Get("roomNumber")
	if roomNumber == "" {
		w.WriteHeader(404)
		w.Write([]byte("must supply roomNumber"))
		return
	}

	room := roomMgr.getRoomByNumber(roomNumber)
	if room == nil {
		w.WriteHeader(404)
		w.Write([]byte("no room found for room number:" + roomNumber))
		return
	}

	var cfg = &RoomConfigJSON{}

	// wrap our reader before passing it to the json decoder
	reader := JsonConfigReader.New(strings.NewReader(body))
	err := json.NewDecoder(reader).Decode(cfg)
	if err != nil {
		log.Println("json un-marshal error:", err)
		w.WriteHeader(404)
		w.Write([]byte("json un-marshal error"))
		return
	}

	roomConfig := newRoomConfigFromJSON(cfg)
	// if room.config.playerNumAcquired != roomConfig.playerNumAcquired {
	// 	w.WriteHeader(404)
	// 	w.Write([]byte(fmt.Sprintf("player number not equal,old:%d, your cfg:%d\n", room.config.playerNumAcquired,
	// 		roomConfig.playerNumAcquired)))
	// 	return
	// }

	log.Printf("bind room cfg to room:%s, cfg:%+v\n", roomNumber, cfg)
	room.config = roomConfig
}

// attachDealCfg2Room 通过mjtest测试工具附加发牌配置到某一个房间
func (mgr *MonkeyMgr) attachDealCfg2Room(w http.ResponseWriter, r *http.Request, body string) {
	var roomNumber = r.URL.Query().Get("roomNumber")
	if roomNumber == "" {
		w.WriteHeader(404)
		w.Write([]byte("must supply roomNumber"))
		return
	}

	room := roomMgr.getRoomByNumber(roomNumber)
	if room == nil {
		w.WriteHeader(404)
		w.Write([]byte("no room found for room number:" + roomNumber))
		return
	}

	var cfgs, errS = mgr.onUploadCfgs(body)

	if errS != nil {
		w.WriteHeader(404)
		w.Write([]byte(errS.Error()))
		return
	}

	if len(cfgs) < 1 {
		w.WriteHeader(404)
		w.Write([]byte("\n no cfg uploaded"))
		return
	}

	var stateConst = room.state.getStateConst()
	var cfg = cfgs[0]

	// 要求的玩家数量不一致
	if cfg.playerCount() != room.config.playerNumAcquired {
		w.WriteHeader(404)
		w.Write([]byte(fmt.Sprintf("player number not the same, room require:%d, but cfg:%d", room.config.playerNumAcquired,
			cfg.playerCount())))
		return
	}

	// 如果是强制要求顺序，则必须所有的userID不能为空
	if cfg.isForceConsistent {
		for i, tc := range cfg.monkeyUserCardsCfgList {
			if tc.userID == "" {
				w.WriteHeader(404)
				w.Write([]byte(fmt.Sprintf("player %d not supply userID", i)))
				return
			}
		}
	}

	err := false

	switch stateConst {
	case pokerface.RoomState_SRoomIdle:
		break
	case pokerface.RoomState_SRoomPlaying:
		w.WriteHeader(404)
		w.Write([]byte("room is playing state, not allow to attach deal cfg"))
		err = true
		break
	case pokerface.RoomState_SRoomWaiting:

		if cfg.playerCount() < len(room.players) {
			w.WriteHeader(404)
			w.Write([]byte("current player number is large than config's"))
			err = true
			break
		}

		if cfg.isForceConsistent {
			for i, p := range room.players {
				if p.userID() != cfg.monkeyUserCardsCfgList[i].userID {
					w.WriteHeader(404)
					w.Write([]byte("player userID or sequence not match"))
					err = true
					break
				}
			}
		}

		break
	}

	if !err {
		log.Printf("bind deal cfg to room:%s, cfg:%s\n", roomNumber, cfg.name)
		room.monkeyCfg = cfg
		// 如果需要，重设一下风牌
		// if cfg.windID != "" {
		// 	room.forceWind(dict[cfg.windID])
		// }

		// 重设置一下庄家ID
		if cfg.monkeyUserCardsCfgList[0].userID != "" {
			room.bankerUserID = cfg.monkeyUserCardsCfgList[0].userID
		}

	}
}

// doUploadCfgs 通过mjtest上传monkey房间配置
func (mgr *MonkeyMgr) doUploadCfgs(w http.ResponseWriter, r *http.Request, body string) {
	var cfgs, err = mgr.onUploadCfgs(body)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}

	mgr.onUpdateCfgs(cfgs)

	w.WriteHeader(200)

	type Objx struct {
		Name           string `json:"name"`
		PlayerRequired int    `json:"playerRequired"`
	}

	type Reply struct {
		Cfgs []Objx `json:"cfgs"`
	}

	var array = make([]Objx, 0, len(cfgs))
	for _, cfg := range cfgs {
		objx := Objx{}
		objx.Name = cfg.name
		objx.PlayerRequired = cfg.playerCount()
		array = append(array, objx)
	}

	reply := Reply{}
	reply.Cfgs = array

	buf, err := json.Marshal(reply)
	if err != nil {
		log.Println("decode err: ", err)
	}

	w.Write(buf)

}

// onUpdateCfgs 更新monkey房间配置
func (mgr *MonkeyMgr) onUpdateCfgs(cfgs []*MonkeyCfg) {
	mgr.cfgs = make([]*MonkeyCfg, 0, len(cfgs))
	mgr.cfgs = append(mgr.cfgs, cfgs...)
	// if len(cfgs) > 0 {
	// 	mgr.currentCfg = cfgs[0]
	// }
}

// getCfgByName 根据名字获得对应的配置
func (mgr *MonkeyMgr) getCfgByName(cfgName string) *MonkeyCfg {
	for _, cfg := range mgr.cfgs {
		if cfg.name == cfgName {
			return cfg
		}
	}

	return nil
}

// verifyHeader 检查上传的csv文件头部是否有效
func (mgr *MonkeyMgr) verifyHeader(record []string) error {
	headers := []string{"名称", "类型", "庄家userID", "庄家手牌", "庄家动作提示", "userID2", "手牌", "动作提示", "userID3", "手牌",
		"动作提示", "userID4", "手牌", "动作提示", "强制一致", "房间配置ID", "是否连庄"}

	if len(headers) != len(record) {
		return fmt.Errorf("csv file not match, maybe you use old versoin pokerface test client, input header length:%d, require:%d", len(record), len(headers))
	}

	for i, h := range headers {
		if h != record[i] {
			return fmt.Errorf("csv file not match, maybe you use old versoin pokerface test client, %s != %s", h, record[i])
		}
	}

	return nil
}

// extractUserCardsCfg 从csv文件里面抽取玩家发牌配置
func (mgr *MonkeyMgr) extractUserCardsCfg(record []string, userIndex int) *MonkeyUserCardsCfg {
	beginIdx := userIndex*3 + 2
	tuc := newMonkeyUserCardsCfg(userIndex == 0, userIndex)

	// userID
	userID := record[beginIdx]
	// if userID == "" {
	// 	return nil
	// }

	tuc.userID = userID
	// 手牌
	handcomps := record[beginIdx+1]
	handcomps = strings.Trim(handcomps, " \t")
	if handcomps == "" {
		// 必须有手牌配置
		return nil
	}
	handcomps = strings.Replace(handcomps, "，", ",", -1)
	var hands = strings.Split(handcomps, ",")
	mgr.trimLeftRight(hands)
	if len(hands) >= 13 {
		tuc.setHandCards(hands)
	}

	// 花牌
	// flowercomps := record[beginIdx+2]
	// flowercomps = strings.Trim(flowercomps, " \t")
	// if flowercomps != "" {
	// 	flowercomps = strings.Replace(flowercomps, "，", ",", -1)
	// 	var flowers = strings.Split(flowercomps, ",")
	// 	if len(flowers) > 0 {
	// 		mgr.trimLeftRight(flowers)
	// 		tuc.setFlowerCards(flowers)
	// 	}
	// }

	// 动作提示
	tipscomps := record[beginIdx+2]
	tipscomps = strings.Trim(tipscomps, " \t")
	if tipscomps != "" {
		tipscomps = strings.Replace(tipscomps, "，", ",", -1)
		var tips = strings.Split(tipscomps, ",")
		if len(tips) > 0 {
			mgr.trimLeftRight(tips)
			tuc.setActionTips(tips)
		}
	}

	return tuc
}

// onUploadCfgs 分析配置
func (mgr *MonkeyMgr) onUploadCfgs(body string) ([]*MonkeyCfg, error) {
	var cfgs = make([]*MonkeyCfg, 0, 10)
	textReader := strings.NewReader(body)
	reader := csv.NewReader(textReader)

	// 跳过头部第一行
	record, err := reader.Read()
	err = mgr.verifyHeader(record)
	if err != nil {
		log.Println("verifyHeader:", err.Error())
		return nil, err
	}

	for {
		record, err = reader.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}

		if len(record) == 0 {
			continue
		}

		var name = record[0]
		if name == "" {
			continue
		}

		var ttype = record[1]
		if ttype != "大丰关张" {
			continue
		}

		//log.Println("cfg name: ", name)
		var cfg = newMonkeyCfg(name)

		// 4 个玩家的手牌，花牌，和动作提示
		for i := 0; i < 4; i++ {
			tuc := mgr.extractUserCardsCfg(record, i)
			if tuc != nil {
				cfg.monkeyUserCardsCfgList = append(cfg.monkeyUserCardsCfgList, tuc)
			}
		}

		var forceConsistent = strings.Trim(record[14], "\t ")
		if forceConsistent == "1" {
			cfg.isForceConsistent = true
		}

		var isContinuousBanker = strings.Trim(record[16], "\t ")
		if isContinuousBanker == "1" {
			cfg.isContinuousBanker = true
		}

		if cfg.isValid() {
			cfgs = append(cfgs, cfg)
		} else {
			err = fmt.Errorf("config %s invalid", name)
			return nil, err
		}
	}

	return cfgs, nil
}

// trimLeftRight 从左，从右裁剪所有数组中的字符串
func (mgr *MonkeyMgr) trimLeftRight(draws []string) {
	for index, draw := range draws {
		draws[index] = strings.Trim(draw, " \t")
	}
}
