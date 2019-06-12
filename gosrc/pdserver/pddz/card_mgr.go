package pddz

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"pokerface"
)

// CardMgr 扑克牌管理器
// 主要是发牌、抽牌
type CardMgr struct {
	room       *Room
	players    []*PlayerHolder
	wallCards  []*Card
	rand       *rand.Rand
	customDraw []int
}

// newCardMgr 创建一个CardMgr对象
func newCardMgr(room *Room, players []*PlayerHolder) *CardMgr {
	tm := CardMgr{}
	tm.room = room
	tm.players = players
	tm.rand = room.rand

	//
	maxCardCount := 54

	var wallCards = make([]*Card, maxCardCount)
	cnt := 0
	for i := R2H; i < CARDMAX; i++ {
		wallCards[cnt] = &Card{cardID: i}
		cnt++
	}

	tm.wallCards = shuffleArray(wallCards[0:cnt], tm.rand)
	return &tm
}

// Implementing Fisher–Yates shuffle
func shuffleArray(ar []*Card, rnd *rand.Rand) []*Card {
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

// nextPlayreImpl 下一个玩家
func (tm *CardMgr) nextPlayerImpl(player *PlayerHolder) *PlayerHolder {
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
func (tm *CardMgr) prevPlayerImpl(player *PlayerHolder) *PlayerHolder {
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
func (tm *CardMgr) rightOpponent(curPlayer *PlayerHolder) *PlayerHolder {
	return tm.nextPlayerImpl(curPlayer)
}

// leftOpponent 上家
func (tm *CardMgr) leftOpponent(curPlayer *PlayerHolder) *PlayerHolder {
	return tm.prevPlayerImpl(curPlayer)
}

// getOrderPlayers 依据逆时针获得下家，下下家，下下下家
func (tm *CardMgr) getOrderPlayers(curPlayer *PlayerHolder) []*PlayerHolder {
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

// getOrderPlayersWithFirst 依据逆时针获得下家，下下家，下下下家
func (tm *CardMgr) getOrderPlayersWithFirst(curPlayer *PlayerHolder) []*PlayerHolder {
	var length = len(tm.players)
	var orderPlayers = make([]*PlayerHolder, length)
	orderPlayers[0] = curPlayer

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

	for i := 1; i < (length); i++ {
		orderPlayers[i] = tm.players[(i+idx)%length]
	}

	return orderPlayers
}

func (tm *CardMgr) getOrderPlayersIndex(first *PlayerHolder, pos *PlayerHolder) int {
	if first == pos {
		return 0
	}

	orderPlayers := tm.getOrderPlayersWithFirst(first)

	for i, p := range orderPlayers {
		if p == pos {
			return i
		}
	}

	log.Panicf("getOrderPlayersIndex failed, first:%d, pos:%d\n", first.chairID, pos.chairID)
	return 0
}

// drawForPlayer 为玩家抽牌
func (tm *CardMgr) drawForPlayer(player *PlayerHolder, reverse bool) (ok bool, handCard *Card) {
	handCard = nil

	reserved := 0
	if len(tm.wallCards) <= reserved {
		log.Println("wall card less than:", reserved)
		ok = false
		return
	}

	ok, handCard = tm.drawNonFlower(player, reverse)
	return
}

// drawForPlayer 为玩家抽到一张非花牌的牌
func (tm *CardMgr) drawNonFlower(player *PlayerHolder, reverse bool) (ok bool, handCard *Card) {
	handCard = nil

	if len(tm.wallCards) < 1 {
		log.Panic("wall cards empty")
		ok = false
		return
	}

	ok = false

	for len(tm.wallCards) > 0 {
		var t *Card
		if !reverse {
			t = tm.drawOne()
		} else {
			t = tm.drawOneReverse()
		}

		nt := &Card{drawBy: player.userID(), cardID: t.cardID}

		// 普通牌，停止抽牌
		player.cards.addHandCard(nt)
		handCard = nt
		ok = true
		break
		// }
	}
	return
}

// removeCardFromWall 从牌墙中移除一张牌
func (tm *CardMgr) removeCardFromWall(cardID int) *Card {

	for i, v := range tm.wallCards {
		if v.cardID == cardID {
			// 删除一个元素
			wt := tm.wallCards[0:i]
			rm := tm.wallCards[i+1:]
			tm.wallCards = append(wt, rm...)

			return v
		}
	}

	return nil
}

// drawOne 抽取一张牌
func (tm *CardMgr) drawOne() *Card {
	// monkey测试如果配置了抽牌序列则按照配置来抽牌
	if len(tm.customDraw) > 0 {
		cardID := tm.customDraw[0]
		tm.customDraw = tm.customDraw[1:]
		t := tm.removeCardFromWall(cardID)
		//Debug.Assert(t != null, "custom draw failed")
		if t == nil {
			log.Println("custom draw failed:", cardID)
		} else {
			return t
		}
	}

	if len(tm.wallCards) < 1 {
		log.Panic("wallCards is empty")
		return nil
	}

	t := tm.wallCards[0]
	// 如果此时wallCards为1长度，[1:]则使得新数组长度为0
	tm.wallCards = tm.wallCards[1:]

	return t
}

// drawOneReverse 从尾部抽取
func (tm *CardMgr) drawOneReverse() *Card {
	lll := len(tm.wallCards)
	if lll < 1 {
		log.Panic("wallCards is empty")
		return nil
	}

	// 取尾部一个
	t := tm.wallCards[lll-1]
	// 如果此时wallCards为1长度，[1:]则使得新数组长度为0
	tm.wallCards = tm.wallCards[0 : lll-1]

	return t
}

// drawForMonkeys 为测试构建发牌牌表
func (tm *CardMgr) drawForMonkeys() {
	log.Println("draw for monkeys, room:", tm.room.ID)
	var bankerPlayer = tm.room.bankerPlayer()
	var cfg = tm.room.monkeyCfg

	var tcfg *MonkeyUserCardsCfg

	// 杠后牌
	// if len(cfg.kongDraws) > 0 {
	// 	m := len(cfg.kongDraws)
	// 	j := len(tm.wallCards) - 1
	// 	for i := 0; i < m; i++ {
	// 		tid := cfg.kongDraws[i]
	// 		// log.Println("drawKongX-monkey:", tid)
	// 		if tm.wallCards[j].cardID == tid {
	// 			j--
	// 			continue
	// 		} else {
	// 			// find tid and swap
	// 			var k int
	// 			found := false
	// 			for q := 0; q < j; q++ {
	// 				if tm.wallCards[q].cardID == tid {
	// 					k = q
	// 					found = true
	// 					break
	// 				}
	// 			}

	// 			if !found {
	// 				log.Panicln("can't found kongDraws card in draw seq:", tid)
	// 			}

	// 			t := tm.wallCards[j]
	// 			tm.wallCards[j] = tm.wallCards[k]
	// 			tm.wallCards[k] = t

	// 			j--
	// 		}

	// 	}

	// 	// tm.dumpWallCards()
	// }

	// 抽庄家的牌
	tcfg = cfg.monkeyUserCardsCfgList[0]
	tm.fillFor(bankerPlayer, tcfg)

	var orderPlayers = tm.getOrderPlayers(bankerPlayer)

	// 如果配置了抽牌系列，则保存一下抽牌系列，以便后面按照这个系列来抽牌
	// 不能直接使用cfg里面的draw数组，因为那样会修改它，下次就不能用了
	if len(cfg.draws) > 0 {
		var customDraw = make([]int, len(cfg.draws))
		copy(customDraw, cfg.draws)
		tm.customDraw = customDraw
	} else {
		tm.customDraw = nil
	}

	// 按照顺序为其他玩家抽牌
	var i = 1
	for _, player := range orderPlayers {
		tcfg = cfg.monkeyUserCardsCfgList[i]
		tm.fillFor(player, tcfg)
		i++
	}

	// 为不足够牌的玩家补牌
	for _, player := range tm.players {
		tm.padPlayerCards(player)
	}
}

// padPlayerCards 如果玩家的手牌不足够13张则为其抽牌补足
func (tm *CardMgr) padPlayerCards(player *PlayerHolder) {
	var total = 16
	// if player == tm.room.bankerPlayer() {
	// 	total++
	// }

	var reamin = total - player.cards.cardCountInHand()
	for i := 0; i < reamin; i++ {
		tm.drawNonFlower(player, false)
	}
}

// fillFor 为player填充手牌列表
func (tm *CardMgr) fillFor(player *PlayerHolder, cfgUserCards *MonkeyUserCardsCfg) {
	var cards = player.cards

	if len(cfgUserCards.handCards) > 0 {
		tm.fillCards(player, cfgUserCards.handCards, cards)
	}
}

// fillCards 根据cardIDs为player填充牌表
func (tm *CardMgr) fillCards(player *PlayerHolder, cardIDs []int, cards *PlayerCards) {
	for _, cardID := range cardIDs {
		var t = tm.drawWith(cardID)
		nt := &Card{drawBy: player.userID(), cardID: t.cardID}

		// if cardID == tm.room.pseudoFlowerCardID {
		// 	cards.addPesudoFlowerCard(nt)
		// } else if nt.isFlower() {
		// 	cards.addFlowerCard(nt)
		// } else {
		cards.addHandCard(nt)
		// }
	}

	// tm.dumpWallCards()
}

// dumpWallCards 打印牌墙
func (tm *CardMgr) dumpWallCards() {
	buf := bytes.NewBufferString("wall Cards:")
	for _, t := range tm.wallCards {
		buf.WriteString(dictName[t.cardID] + ",")
	}
	log.Println(buf.String())
}

// drawWith 从牌墙中抽取指定的牌
func (tm *CardMgr) drawWith(cardID int) *Card {
	var card = tm.removeCardFromWall(cardID)

	if nil == card {
		log.Panic("DrawWith, no card remain")
		return nil
	}

	return card
}

// drawForAll 为所有人发牌
func (tm *CardMgr) drawForAll() {
	for _, player := range tm.players {
		if !player.cards.isEmpty() {
			log.Panic("Player card list should be empty")
			player.cards.clear()
		}
	}

	if tm.room.monkeyCfg != nil {
		tm.drawForMonkeys()
		return
	}

	// 抽取17张牌
	for i := 0; i < 17; i++ {
		for _, player := range tm.players {
			// 不会出现无牌可抽情况
			tm.drawNonFlower(player, false)
		}
	}

	// tm.drawNonFlower(tm.room.bankerPlayer(), false)
}

// cardCountInWall 牌墙中剩余的牌张数
func (tm *CardMgr) cardCountInWall() int {
	return len(tm.wallCards)
}

// playerDiscard 处理玩家出牌
func (tm *CardMgr) playerDiscard(player *PlayerHolder, cardHand *pokerface.MsgCardHand) *CardHand {
	cards := player.cards
	if !cards.hasMsgCardHandOnHand(cardHand) {
		log.Panic("Player Discard failed, no such card in hand list")
		return nil
	}

	t := cards.removeMsgCardHandFromHand(cardHand)
	cards.addDiscardedCardHand(t)

	return t
}

// wallEmpty 牌墙是否已经空
func (tm *CardMgr) wallEmpty() bool {
	return len(tm.wallCards) == 0
}

// cardsDiscardAble 是否可以打牌
func (tm *CardMgr) cardsDiscardAble(player *PlayerHolder, cards []int32) (*pokerface.MsgCardHand, bool) {
	msgCardHand := player.cards.convertMsgCardHand(cards)
	return msgCardHand, msgCardHand != nil
}
