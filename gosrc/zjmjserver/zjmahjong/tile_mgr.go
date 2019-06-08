package zjmahjong

import (
	"mahjong"
	"math/rand"

	"github.com/sirupsen/logrus"
)

// TileMgr 牌管理器
// 主要是发牌、抽牌
type TileMgr struct {
	room       *Room
	players    []*PlayerHolder
	wallTiles  []*Tile
	rand       *rand.Rand
	customDraw []int
	cl         *logrus.Entry
}

// newTileMgr 创建一个TileMgr对象
func newTileMgr(room *Room, players []*PlayerHolder) *TileMgr {
	tm := TileMgr{}
	tm.room = room
	tm.players = players
	tm.rand = room.rand
	tm.cl = room.cl

	/*
		台灣16張麻將共144張，分六大類：萬子牌、索子牌、筒子牌、字牌、三元牌、花牌
		萬子牌：一萬至九萬各四張，共計三十六張
		索子牌： 一索至九索各四張，共計三十六張
		筒子牌：一筒至九筒各四張，共計三十六張
		三元牌：中、發、白各四張，共計十二張
		字 牌：東、南、西、北各四張，共計十六張
		花 牌：春、夏、秋、冬、梅、蘭、菊、竹各一張，總計八張

		廣東13張麻將共136張，分五大類：四風牌、三元牌（--合稱為「字牌」）跟萬子牌、索子牌、筒子牌（--合稱為「數牌」）共五類
	*/
	// 标准麻将136张牌: 3*(9*4) + 7*4 = 136
	// 索子，万子，筒子的数量各自是：9*4
	// 东南西北中发白：7*4
	var wallTiles = make([]*Tile, 136)
	cnt := 0
	for i := MAN; i < PLUM; i++ {
		for j := 0; j < 4; j++ {
			wallTiles[cnt] = &Tile{tileID: i}
			cnt++
		}
	}

	tm.wallTiles = shuffleArray(wallTiles, tm.rand)
	return &tm
}

// Implementing Fisher–Yates shuffle
func shuffleArray(ar []*Tile, rnd *rand.Rand) []*Tile {
	// If running on Java 6 or older, use `new Random()` on RHS here
	for i := len(ar) - 1; i > 0; i-- {
		index := rnd.Intn(i + 1)
		// Simple swap
		a := ar[index]
		ar[index] = ar[i]
		ar[i] = a
	}

	return ar
}

// tileRemainInHandOrWall 还剩下多少张在其他玩家手上以及牌墙上
func (tm *TileMgr) tileRemainInHandOrWall() []int {
	slots := make([]int, TILEMAX)
	for _, tile := range tm.wallTiles {
		slots[tile.tileID]++
	}

	for _, p := range tm.players {
		pt := p.tiles
		for e := pt.hand.Front(); e != nil; e = e.Next() {
			t := e.Value.(*Tile)
			slots[t.tileID]++
		}

		// 如果暗杠此时对手有人是不可见的（大丰的暗杠，如果有吃椪，需要明牌暗杠）
		if pt.chowPongExposedKongCount() == 0 {
			// 则把暗杠也放在候选听牌列表上
			for e := pt.fixedMelds.Front(); e != nil; e = e.Next() {
				m := e.Value.(*Meld)
				if m.isConcealedKong() {
					slots[m.t1.tileID] += 4
				}
			}
		}
	}

	return slots
}

// discardAble 是否可以打出该张牌
func (tm *TileMgr) discardAble(player *PlayerHolder, msgTile int) bool {
	tiles := player.tiles
	if !tiles.hasTileInHand(msgTile) {
		tm.cl.Printf("discardAble, no tile:%d\n", msgTile)
		return false
	}

	// 如果手牌数量大于2，而且刚刚吃碰，那么就有限制
	// if tiles.tileCountInHand() > 2 &&
	// 	player.hStatis.latestChowPongTileLocked.tileID == msgTile {
	// 	return false
	// }

	// if player.hStatis.isRichi || tm.wallEmpty() {
	// 	// 报听情况下只能打最后一张牌
	// 	if msgTile != tiles.latestHandTile().tileID {
	// 		return false
	// 	}
	// }
	return true
}

// triplet2KongAble 是否可以加杠
func (tm *TileMgr) triplet2KongAble(player *PlayerHolder, tileID int) bool {
	var tiles = player.tiles
	if !tiles.hasTileInHand(tileID) {
		tm.cl.Printf("Triplet2KongAble, no tile %d\n", tileID)
		return false
	}

	return nil != tiles.tripletMeldWith(tileID)
}

// triplet2Kong 加杠操作
func (tm *TileMgr) triplet2Kong(player *PlayerHolder, tileID int) *Tile {
	tiles := player.tiles
	tile := tiles.removeTileInHand(tileID)

	if tile == nil {
		tm.cl.Panic("Triplet2Kong tile must not be null")
		return nil
	}

	var meld = tiles.tripletMeldWith(tileID)
	if meld == nil {
		tm.cl.Panic("Triplet2Kong meld must not be null")
		return nil
	}

	meld.triplet2Kong(tile)

	return tile
}

// winSelfDrawAble 是否可以胡牌
func (tm *TileMgr) winSelfDrawAble(player *PlayerHolder) bool {
	return player.tiles.winAble()
}

// kongConcealed 暗杠
func (tm *TileMgr) kongConcealed(player *PlayerHolder, msgTile int) {
	// 从player hands列表移除对应4个牌
	var tiles = make([]*Tile, 4)
	for i := 0; i < 4; i++ {
		tiles[i] = player.tiles.removeTileInHand(msgTile)
	}

	// 加入player的meld列表
	var meld = &Meld{
		t1: tiles[0],
		t2: tiles[1],
		t3: tiles[2],
		t4: tiles[3],
		mt: mahjong.MeldType_enumMeldTypeConcealedKong}
	player.tiles.addMeld(meld)
}

// kongConcealedAble 是否可以暗杠
func (tm *TileMgr) kongConcealedAble(player *PlayerHolder, msgTile int) bool {
	return player.tiles.concealedKongAbleWith(msgTile)
}

// winChuck 放铳胡牌
func (tm *TileMgr) winChuck(player *PlayerHolder, dq *TaskPlayerReAction) {

	var discardPlayer = dq.actionPlayer
	var targetTile = dq.actionTile
	if dq.isForRobKong {
		discardPlayer.tiles.triplet2KongRollback(targetTile)
	} else {
		// 从discard player移除
		// 对于一炮多响的情形，tile会多次从Discarded队列中移除
		discardPlayer.tiles.removeTileInDiscarded(targetTile)
	}

	// 加入到手牌列表
	// 对于一炮多响情形，tile会加入到几个胡牌玩家的手牌列表中
	player.tiles.addHandTile(targetTile)
}

// winChuckAble 是否可以放铳胡牌
func (tm *TileMgr) winChuckAble(player *PlayerHolder, dq *TaskPlayerReAction) bool {
	return player.tiles.winAbleWith(dq.actionTile)
}

// kongExposed 明杠
func (tm *TileMgr) kongExposed(player *PlayerHolder, dq *TaskPlayerReAction) *Meld {
	// 从discard player移除
	var discardPlayer = dq.actionPlayer
	var targetTile = dq.actionTile

	discardPlayer.tiles.removeTileInDiscarded(targetTile)

	// 从player hands列表移除对应3个牌
	var tiles = make([]*Tile, 4)
	for i := 0; i < 3; i++ {
		tiles[i] = player.tiles.removeTileInHand(targetTile.tileID)
	}
	// 第4个牌
	tiles[3] = targetTile

	// 加入player的meld列表
	var meld = &Meld{
		t1: tiles[0],
		t2: tiles[1],
		t3: tiles[2],
		t4: tiles[3],
		mt: mahjong.MeldType_enumMeldTypeExposedKong}
	player.tiles.addMeld(meld)

	return meld
}

// kongExposedAble 是否可以明杠
func (tm *TileMgr) kongExposedAble(player *PlayerHolder, dq *TaskPlayerReAction, msgActionMeld *mahjong.MsgMeldTile) bool {
	if msgActionMeld.GetMeldType() != int32(mahjong.MeldType_enumMeldTypeExposedKong) {
		tm.cl.Println("KongExposedAble require meld type ExposedKong")
		return false
	}

	var targetTile = dq.actionTile
	if msgActionMeld.GetTile1() != int32(targetTile.tileID) {
		tm.cl.Printf("KongExposedAble require same tile:%d\n", targetTile.tileID)
		return false
	}

	var count = player.tiles.tileCountInHandOf(targetTile.tileID)
	if count == 3 {
		return true
	}

	tm.cl.Printf("KongExposedAble require 3 tile:%d\n", targetTile.tileID)
	return false
}

// pong 碰牌
func (tm *TileMgr) pong(player *PlayerHolder, dq *TaskPlayerReAction) *Meld {
	// 从discard player移除
	var discardPlayer = dq.actionPlayer
	var targetTile = dq.actionTile
	discardPlayer.tiles.removeTileInDiscarded(targetTile)

	// 从player hands列表移除对应2个牌
	var tiles = make([]*Tile, 3)
	for i := 0; i < 2; i++ {
		tiles[i] = player.tiles.removeTileInHand(targetTile.tileID)
	}
	// 第3个牌
	tiles[2] = targetTile

	// 加入player的meld列表
	var meld = &Meld{t1: tiles[0], t2: tiles[1], t3: tiles[2], t4: InvalidTile, mt: mahjong.MeldType_enumMeldTypeTriplet}
	player.tiles.addMeld(meld)

	return meld
}

// pongAble 是否可以碰牌
func (tm *TileMgr) pongAble(player *PlayerHolder, dq *TaskPlayerReAction, msgActionMeld *mahjong.MsgMeldTile) bool {
	if msgActionMeld.GetMeldType() != int32(mahjong.MeldType_enumMeldTypeTriplet) {
		tm.cl.Println("Ponable require meld type Triplet")
		return false
	}

	var targetTile = dq.actionTile
	if msgActionMeld.GetTile1() != int32(targetTile.tileID) {
		tm.cl.Printf("Ponable require same tile:%d\n", targetTile.tileID)
		return false
	}

	var count = player.tiles.tileCountInHandOf(targetTile.tileID)
	if count >= 2 {
		return true
	}
	tm.cl.Printf("Ponable require at least 2 tile:%d\n", targetTile.tileID)
	return false
}

// chow 吃牌
func (tm *TileMgr) chow(player *PlayerHolder, dq *TaskPlayerReAction) *Meld {
	// 从discard player移除
	var discardPlayer = dq.actionPlayer
	var targetTile = dq.actionTile
	discardPlayer.tiles.removeTileInDiscarded(targetTile)

	var msgActionMeld = dq.actionMeld()

	// 从player hands列表移除对应2个牌
	var t1Id = int(msgActionMeld.GetTile1())
	var tileIds = []int{t1Id, t1Id + 1, t1Id + 2}

	var tiles = make([]*Tile, 3)
	var i = 0
	for _, tid := range tileIds {
		if tid != targetTile.tileID {
			tiles[i] = player.tiles.removeTileInHand(tid)
		} else {
			tiles[i] = targetTile
		}
		i++
	}

	for j, t := range tiles {
		if t == nil {
			tm.cl.Panic("chow, t is nil:", tileIds[j])
		}
	}

	// 加入player的meld列表
	var meld = &Meld{t1: tiles[0], t2: tiles[1], t3: tiles[2], t4: InvalidTile, mt: mahjong.MeldType_enumMeldTypeSequence}
	player.tiles.addMeld(meld)

	return meld
}

// chowAble 是否可以吃牌
func (tm *TileMgr) chowAble(player *PlayerHolder, dq *TaskPlayerReAction, msgActionMeld *mahjong.MsgMeldTile) bool {
	if msgActionMeld.GetMeldType() != int32(mahjong.MeldType_enumMeldTypeSequence) {
		tm.cl.Println("chowAble require meld type Sequence")
		return false
	}

	var tiles = player.tiles
	var msgTile = dq.actionTile.tileID

	if tm.nextPlayerImpl(dq.actionPlayer) != player {
		tm.cl.Println("ChowAbleWith, player not right opponent")
		return false
	}

	var tile1 = int(msgActionMeld.GetTile1())
	if tile1 != msgTile && !tiles.hasTileInHand(tile1) {
		tm.cl.Printf("ChowAbleWith, player has no tile:%d\n", tile1)
		return false
	}

	var tile2 = tile1 + 1
	if tile2 != msgTile && !tiles.hasTileInHand(tile2) {
		tm.cl.Printf("ChowAbleWith, player has no tile:%d\n", tile2)
		return false
	}

	var tile3 = tile1 + 2
	if tile3 != msgTile && !tiles.hasTileInHand(tile3) {
		tm.cl.Printf("ChowAbleWith, player has no tile:%d\n", tile3)
		return false
	}

	return true
}

// getContributor 获得牌组贡献者
func (tm *TileMgr) getContributor(owner *PlayerHolder, m *Meld) *PlayerHolder {
	myUserID := owner.userID()
	ts := []*Tile{m.t1, m.t2, m.t3, m.t4}
	for _, t := range ts {
		if t.drawBy != myUserID {
			return tm.room.getPlayerByUserID(t.drawBy)
		}
	}

	return owner
}

// nextPlayreImpl 下一个玩家
func (tm *TileMgr) nextPlayerImpl(player *PlayerHolder) *PlayerHolder {
	var players = tm.players
	var length = len(players)
	for i := 0; i < length; i++ {
		if players[i] == player {
			return players[(i+1)%length]
		}
	}

	return nil
}

// prevPlayerImpl 上一个玩家
func (tm *TileMgr) prevPlayerImpl(player *PlayerHolder) *PlayerHolder {
	var players = tm.players
	var length = len(players)
	for i := 0; i < length; i++ {
		if players[i] == player {
			return players[(i-1+length)%length]
		}
	}

	return nil
}

// rightOpponent 下家
func (tm *TileMgr) rightOpponent(curPlayer *PlayerHolder) *PlayerHolder {
	return tm.nextPlayerImpl(curPlayer)
}

// leftOpponent 上家
func (tm *TileMgr) leftOpponent(curPlayer *PlayerHolder) *PlayerHolder {
	return tm.prevPlayerImpl(curPlayer)
}

// getOrderPlayers 依据逆时针获得下家，下下家，下下下家
func (tm *TileMgr) getOrderPlayers(curPlayer *PlayerHolder) []*PlayerHolder {
	var length = len(tm.players)
	var orderPlayers = make([]*PlayerHolder, length-1)

	var idx = -1
	for i := 0; i < length; i++ {
		if tm.players[i] != curPlayer {
			continue
		}

		idx = i
		break
	}

	if idx < 0 {
		return nil
	}

	idx++
	for i := 0; i < (length - 1); i++ {
		orderPlayers[i] = tm.players[(i+idx)%length]
	}

	return orderPlayers
}

// actionForDiscardPlayer 计算出牌玩家可以执行的动作集合
func (tm *TileMgr) actionForDiscardPlayer(discarder *PlayerHolder, newDraw bool) int {
	var tiles = discarder.tiles
	var winable = tiles.winAble()

	var action mahjong.ActionType
	// 自摸胡牌
	if !discarder.hStatis.isWinAbleLocked && winable && newDraw {
		action |= mahjong.ActionType_enumActionType_WIN_SelfDrawn
	}

	if discarder.hStatis.actionCounter == 0 && discarder == tm.room.bankerPlayer() {
		if action&(mahjong.ActionType_enumActionType_WIN_SelfDrawn) != 0 {
			// 天胡，不允许过，客户端只能胡
			action = (mahjong.ActionType_enumActionType_WIN_SelfDrawn)
			return int(action)
		}
	}

	// 让客户端显示一个过按钮
	action = action | mahjong.ActionType_enumActionType_DISCARD | mahjong.ActionType_enumActionType_SKIP

	// 起手听牌者只能胡，或者出牌，而且只能出抽到的那张
	if discarder.hStatis.isRichi {
		return int(action)
	}

	// 暗杠
	if tiles.concealedKongAble() && newDraw {
		action |= mahjong.ActionType_enumActionType_KONG_Concealed
	}

	// 加杠
	// var tile = tiles.latestHandTile()
	// if newDraw && tiles.triplet2KongAbleWith(tile.tileID) {
	// 	action |= ActionType_enumActionType_KONG_Triplet2
	// }

	// 加杠
	if newDraw {
		pongMelds := tiles.pongMelds()
		for _, pm := range pongMelds {
			pongTileID := pm.t1.tileID

			if tiles.tileCountInHandOf(pongTileID) > 0 {
				action |= mahjong.ActionType_enumActionType_KONG_Triplet2
				break
			}
		}
	}

	return int(action)
}

// drawForPlayer 为玩家抽牌
func (tm *TileMgr) drawForPlayer(player *PlayerHolder, needDrawData bool, reserveLast bool) (ok bool, handTile *Tile, newFlowers []*Tile) {
	handTile = nil
	newFlowers = nil

	reserved := 0
	if reserveLast {
		reserved = 1
	}

	if len(tm.wallTiles) <= reserved {
		ok = false
		return
	}

	ok, handTile, newFlowers = tm.drawNonFlower(player, needDrawData, reserved)
	return
}

// drawForPlayer 为玩家抽到一张非花牌的牌
func (tm *TileMgr) drawNonFlower(player *PlayerHolder, needDrawData bool, reserved int) (ok bool, handTile *Tile, newFlowers []*Tile) {
	handTile = nil
	if len(tm.wallTiles) <= reserved {
		tm.cl.Panic("wall tiles empty")
		ok = false
		return
	}

	ok = false
	flowerCnt := 0
	for len(tm.wallTiles) > reserved {
		var t = tm.drawOne()
		nt := &Tile{drawBy: player.userID(), tileID: t.tileID}
		// 普通牌，停止抽牌
		player.tiles.addHandTile(nt)
		if needDrawData {
			handTile = nt
		}
		ok = true
		break
	}

	if needDrawData {
		newFlowers = newFlowers[:flowerCnt]
	}

	return
}

// removeTileFromWall 从牌墙中移除一张牌
func (tm *TileMgr) removeTileFromWall(tileID int) *Tile {
	var length = len(tm.wallTiles)
	for i, v := range tm.wallTiles {
		if v.tileID == tileID {
			// 删除一个元素，通过交换尾部的元素来填补空隙
			if i != length-1 {
				tm.wallTiles[i] = tm.wallTiles[length-1]
			}
			tm.wallTiles = tm.wallTiles[:length-1]

			return v
		}
	}

	return nil
}

// drawOne 抽取一张牌
func (tm *TileMgr) drawOne() *Tile {
	// monkey测试如果配置了抽牌序列则按照配置来抽牌
	if len(tm.customDraw) > 0 {
		tileID := tm.customDraw[0]
		tm.customDraw = tm.customDraw[1:]
		t := tm.removeTileFromWall(tileID)
		//Debug.Assert(t != null, "custom draw failed")
		if t == nil {
			tm.cl.Println("custom draw failed:", tileID)
		} else {
			return t
		}
	}

	if len(tm.wallTiles) < 1 {
		tm.cl.Panic("wallTiles is empty")
		return nil
	}

	var length = len(tm.wallTiles)
	var next = tm.rand.Intn(length)
	t := tm.wallTiles[next]
	if next != length-1 {
		tm.wallTiles[next] = tm.wallTiles[length-1]
	}

	tm.wallTiles = tm.wallTiles[:length-1]
	return t
}

// drawForMonkeys 为测试构建发牌牌表
func (tm *TileMgr) drawForMonkeys() {
	tm.cl.Println("draw for monkeys, room:", tm.room.ID)
	var bankerPlayer = tm.room.bankerPlayer()
	var cfg = tm.room.monkeyCfg

	var tcfg *MonkeyUserTilesCfg
	// 抽庄家的牌
	tcfg = cfg.monkeyUserTilesCfgList[0]
	tm.fillFor(bankerPlayer, tcfg)

	var orderPlayers = tm.getOrderPlayers(bankerPlayer)
	// 按照顺序为其他玩家抽牌
	var i = 1
	for _, player := range orderPlayers {
		tcfg = cfg.monkeyUserTilesCfgList[i]
		tm.fillFor(player, tcfg)
		i++
	}

	// 为不足够牌的玩家补牌
	for _, player := range tm.players {
		tm.padPlayerTiles(player)
	}

	// 如果配置了抽牌系列，则保存一下抽牌系列，以便后面按照这个系列来抽牌
	// 不能直接使用cfg里面的draw数组，因为那样会修改它，下次就不能用了
	if len(cfg.draws) > 0 {
		var customDraw = make([]int, len(cfg.draws))
		copy(customDraw, cfg.draws)
		tm.customDraw = customDraw
	} else {
		tm.customDraw = nil
	}
}

// 抽取马牌
func (tm *TileMgr) drawHorseTiles(horseTileCount int) []*Tile {
	remain := tm.tileCountInWall()

	tm.cl.Printf("drawHorseTiles, horseTileCount:%d, remain:%d\n", horseTileCount, remain)

	if horseTileCount > remain {
		horseTileCount = remain
	}

	tiles := make([]*Tile, 0, horseTileCount)

	for i := 0; i < horseTileCount; i++ {
		var t = tm.drawOne()
		nt := &Tile{drawBy: "", tileID: t.tileID}
		tiles = append(tiles, nt)
	}

	return tiles
}

// padPlayerTiles 如果玩家的手牌不足够13张则为其抽牌补足
func (tm *TileMgr) padPlayerTiles(player *PlayerHolder) {
	var total = 13
	if player == tm.room.bankerPlayer() {
		total++
	}

	var reamin = total - player.tiles.tileCountInHand()
	for i := 0; i < reamin; i++ {
		tm.drawNonFlower(player, false, 1)
	}
}

// fillFor 为player填充花牌列表以及手牌列表
func (tm *TileMgr) fillFor(player *PlayerHolder, cfgUserTiles *MonkeyUserTilesCfg) {
	var tiles = player.tiles

	if len(cfgUserTiles.handTiles) > 0 {
		tm.fillTiles(player, cfgUserTiles.handTiles, tiles)
	}
}

// fillTiles 根据tileIDs为player填充牌表
func (tm *TileMgr) fillTiles(player *PlayerHolder, tileIDs []int, tiles *PlayerTiles) {
	for _, tileID := range tileIDs {
		var t = tm.drawWith(tileID)
		nt := &Tile{drawBy: player.userID(), tileID: t.tileID}

		tiles.addHandTile(nt)
	}
}

// drawWith 从牌墙中抽取指定的牌
func (tm *TileMgr) drawWith(tileID int) *Tile {
	var tile = tm.removeTileFromWall(tileID)

	if nil == tile {
		tm.cl.Panic("DrawWith, no tile remain")
		return nil
	}

	return tile
}

// drawForAll 为所有人发牌
func (tm *TileMgr) drawForAll() {
	for _, player := range tm.players {
		if !player.tiles.isEmpty() {
			tm.cl.Panic("Player tile list should be empty")
			player.tiles.clear()
		}
	}

	if tm.room.monkeyCfg != nil {
		tm.drawForMonkeys()
		return
	}

	// 抽取13张牌
	for i := 0; i < 13; i++ {
		for _, player := range tm.players {
			// 不会出现无牌可抽情况
			tm.drawNonFlower(player, false, 1)
		}
	}

	tm.drawNonFlower(tm.room.bankerPlayer(), false, 1)
}

// tileCountInWall 牌墙中剩余的牌张数
func (tm *TileMgr) tileCountInWall() int {
	return len(tm.wallTiles)
}

// playerDiscard 处理玩家出牌
func (tm *TileMgr) playerDiscard(player *PlayerHolder, tileID int) *Tile {
	tiles := player.tiles
	if !tiles.hasTileInHand(tileID) {
		tm.cl.Panic("Player Discard failed, no such card in hand list:", tileID)
		return nil
	}

	t := tiles.removeTileInHand(tileID)
	tiles.addDiscardedTile(t)

	return t
}

// wallEmpty 牌墙是否已经空
func (tm *TileMgr) wallEmpty() bool {
	return len(tm.wallTiles) == 0
}
