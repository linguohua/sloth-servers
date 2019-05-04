package dfmahjong

import (
	"container/list"
	"log"
	"mahjong"
)

// TaskPlayerReAction 等待其他玩家完成操作
type TaskPlayerReAction struct {
	actionPlayer *PlayerHolder
	actionTile   *Tile

	orderPlayers []*PlayerHolder

	waitQueue         *list.List
	isFinished        bool
	chanWait          chan bool
	forceWaitAllReply bool
	isForRobKong      bool

	s *SPlaying
}

// ReActionQueueItem 等待队列项
type ReActionQueueItem struct {
	player      *PlayerHolder
	actions     int
	replyAction int
	msgMeldTile *mahjong.MsgMeldTile
}

// analyseTaskDiscardReaction 分析其他玩家是否可以针对本次出牌进行动作，如果可以则返回一个TaskDiscardReAction用于等待
// 否则返回nil
func analyseTaskDiscardReaction(s *SPlaying, orderPlayers []*PlayerHolder, latestDiscardedTile *Tile, discardPlayer *PlayerHolder) *TaskPlayerReAction {
	for _, player := range orderPlayers {
		var action = 0
		var tiles = player.tiles

		if !player.hStatis.isRichiWinableOnlySelfDrawn &&
			tiles.winAbleWith(latestDiscardedTile) {

			// 如果玩家处于过手胡锁定状态，则不允许胡牌
			// 如果玩家处于听牌状态，且错过了一次胡牌机会，只能自摸胡牌而不能吃铳胡牌
			winLocked := (player.hStatis.isWinAbleLocked)
			winLockedForWinType := (player.hStatis.isWinLockedForGreatWin)

			var isGreatWin = calcWinChuckGreatWinType(s, player, latestDiscardedTile)
			player.hStatis.isWinLockedForGreatWin = isGreatWin

			// 2019-3-10 需求变更：如果过手胡锁定，但是胡的牌型（大胡、小胡）不同时，取消过手胡
			// 例如，本人是A，下一家B打出一个牌我可以胡成小胡而不胡，然后C打出一个牌我可以胡成小胡，此时我不可以
			// 胡C打出的这张牌，因为我之前没有胡B的小胡；但是如果C打出的牌我可以胡成大胡，则此时我可以胡C的这张牌,
			// 因为此时胡出的牌的大胡小胡类型发生了变化，就取消了过手胡的锁定
			if !winLocked || winLockedForWinType != isGreatWin {
				// 如果胡牌是小胡情形，而且自己有吃椪杠，那么就得看对方是不是补花后出的牌或者杠后出的牌
				// 否者不能算小胡
				if isGreatWin {
					// 大胡可以胡
					action |= int(mahjong.ActionType_enumActionType_WIN_Chuck)
				} else {
					// 小胡情形
					tiles := player.tiles
					// exposedMeldCount 不考虑暗杠，也即是说暗杠不影响小门清
					if tiles.isSecondClearFrontAble() {
						// 小门清可以胡对方的牌
						action |= int(mahjong.ActionType_enumActionType_WIN_Chuck)
					} else {
						// 有吃椪杠的情况下，对方得是补花后出的牌，或者杠后出的牌，才可以胡它
						if s.lctx.isMyDrawAndGotFlower2Discarded(discardPlayer) || s.lctx.isMyKong2Discarded(discardPlayer) {
							// 则只有对方补花后出的牌或者杠后出的牌，才可以小胡
							action |= int(mahjong.ActionType_enumActionType_WIN_Chuck)
						}
					}
				}
			} else {
				log.Printf("player %s can win, but in winlocked, winLockedForWinType:%+v, currentWinType:%+v\n",
					player.userID(), winLockedForWinType, isGreatWin)
			}
		}

		// 报听后不能吃椪杠，只能胡和摸牌，打牌的话只能打摸到的牌
		// 需求变更：必须牌墙还有牌，才可以吃椪杠
		if !player.hStatis.isRichi && s.tileMgr.tileCountInWall() > 0 {
			if tiles.exposedKongAbleWith(latestDiscardedTile) {
				action |= int(mahjong.ActionType_enumActionType_KONG_Exposed)
			}

			// 如果在再次摸牌/出牌之前已经不碰同一样牌，则现在也不能碰：
			// 例如C打出B可以碰的牌，B不碰，然后D打出同样的牌，B此时不能碰
			if latestDiscardedTile.tileID != player.hStatis.pongAbleTileLocked && tiles.pongAbleWith(latestDiscardedTile) {
				action |= int(mahjong.ActionType_enumActionType_PONG)
			}
		}

		player.expectedAction = action
	}

	// 只有下家才可以吃牌
	// 报听后不能吃椪杠，只能胡和摸牌，打牌的话只能打摸到的牌
	// 最后一次打出去的牌不能跟本次要吃的牌相同
	// 如果在再次摸牌/出牌之前已经不吃同一样牌，则现在也不能吃：
	// 例如A打出B可以吃的牌，B不吃，D碰，轮到A继续出同样的牌，B此时不能吃
	// 需求变更：必须牌墙还有牌，才可以吃椪杠
	var p = orderPlayers[0]
	if !p.hStatis.isRichi && s.tileMgr.tileCountInWall() > 0 && p.tiles.chowAbleWith(latestDiscardedTile) {
		p.expectedAction |= int(mahjong.ActionType_enumActionType_CHOW)
	}

	var found = false
	for _, v := range orderPlayers {
		if v.expectedAction != 0 {
			found = true
			break
		}
	}

	if !found {
		// 没有任何玩家可以操作
		return nil
	}

	// 增加skip选项,每一个可以操作的玩家，都可以选择“过”
	for _, v := range orderPlayers {
		if v.expectedAction != 0 {
			v.expectedAction |= int(mahjong.ActionType_enumActionType_SKIP)
		}
	}

	var tdr = &TaskPlayerReAction{}
	tdr.orderPlayers = orderPlayers
	tdr.actionTile = latestDiscardedTile
	tdr.actionPlayer = discardPlayer
	tdr.waitQueue = list.New()
	tdr.chanWait = make(chan bool, 1) // buffered channel,1 slots

	// 依优先级次序，把可以动作的玩家增加到等待列表
	// 复制一份orderPlayer，以免add2WaitQueue修改原始数组
	var orderPlayers2 = make([]*PlayerHolder, len(orderPlayers))
	copy(orderPlayers2, orderPlayers)

	add2WaitQueue(orderPlayers2, int(mahjong.ActionType_enumActionType_WIN_Chuck), tdr)
	add2WaitQueue(orderPlayers2, int(mahjong.ActionType_enumActionType_KONG_Exposed), tdr)
	add2WaitQueue(orderPlayers2, int(mahjong.ActionType_enumActionType_PONG), tdr)
	add2WaitQueue(orderPlayers2, int(mahjong.ActionType_enumActionType_CHOW), tdr)

	return tdr
}

// analyseTaskTriplet2KongReaction 加杠后看其他玩家是否可以抢杠
func analyseTaskTriplet2KongReaction(s *SPlaying, orderPlayers []*PlayerHolder, latestDiscardedTile *Tile, discardPlayer *PlayerHolder) *TaskPlayerReAction {
	for _, player := range orderPlayers {
		var action = 0
		var tiles = player.tiles

		// 如果玩家处于过手胡锁定状态，则不允许胡牌
		// 如果玩家处于听牌状态，且错过了一次胡牌机会，只能自摸胡牌而不能吃铳胡牌
		winLocked := player.hStatis.isWinAbleLocked
		if !winLocked && !player.hStatis.isRichiWinableOnlySelfDrawn &&
			tiles.winAbleWith(latestDiscardedTile) {
			// 抢杠胡必是大胡，因为抢杠胡+1辣子
			action |= int(mahjong.ActionType_enumActionType_WIN_Chuck)
		}

		player.expectedAction = action
	}

	var found = false
	for _, v := range orderPlayers {
		if v.expectedAction != 0 {
			found = true
			break
		}
	}

	if !found {
		// 没有任何玩家可以操作
		return nil
	}

	// 增加skip选项,每一个可以操作的玩家，都可以选择“过”
	for _, v := range orderPlayers {
		if v.expectedAction != 0 {
			v.expectedAction |= int(mahjong.ActionType_enumActionType_SKIP)
		}
	}

	var tdr = &TaskPlayerReAction{}
	tdr.orderPlayers = orderPlayers
	tdr.actionTile = latestDiscardedTile
	tdr.actionPlayer = discardPlayer
	tdr.isForRobKong = true
	tdr.waitQueue = list.New()
	tdr.chanWait = make(chan bool, 1) // buffered channel,1 slots

	// 依优先级次序，把可以动作的玩家增加到等待列表
	// 复制一份orderPlayer，以免add2WaitQueue修改原始数组
	var orderPlayers2 = make([]*PlayerHolder, len(orderPlayers))
	copy(orderPlayers2, orderPlayers)

	add2WaitQueue(orderPlayers2, int(mahjong.ActionType_enumActionType_WIN_Chuck), tdr)

	return tdr
}

// add2WaitQueue 把玩家按照顺序加入等待队列中
func add2WaitQueue(orderPlayers2 []*PlayerHolder, filter int, tdr *TaskPlayerReAction) {
	for i, pl := range orderPlayers2 {
		if pl == nil {
			continue
		}
		if (pl.expectedAction & filter) != 0 {
			a := &ReActionQueueItem{}
			a.actions = pl.expectedAction
			a.player = pl

			tdr.waitQueue.PushBack(a)
			orderPlayers2[i] = nil
		}
	}
}

// findWaitQueueItem 根据player找到wait item
func (tdr *TaskPlayerReAction) findWaitQueueItem(player *PlayerHolder) *ReActionQueueItem {
	for e := tdr.waitQueue.Front(); e != nil; e = e.Next() {
		qi := e.Value.(*ReActionQueueItem)
		if qi.player == player {
			return qi
		}
	}
	return nil
}

// isExpectedPlayerAction 是否正在等待的玩家以及动作
func (tdr *TaskPlayerReAction) isExpectedPlayerAction(player *PlayerHolder, action int) bool {
	if tdr.waitQueue == nil || tdr.waitQueue.Len() == 0 || tdr.isFinished {
		return false
	}

	var item = tdr.findWaitQueueItem(player)
	if item == nil {
		return false
	}

	return (item.actions & action) != 0
}

// removeWaitQueueItem 从等待队列中移除一个项
func (tdr *TaskPlayerReAction) removeWaitQueueItem(wi *ReActionQueueItem) {
	for e := tdr.waitQueue.Front(); e != nil; e = e.Next() {
		w := e.Value.(*ReActionQueueItem)
		if w == wi {
			tdr.waitQueue.Remove(e)
			break
		}
	}
}

// replayAction 获得最高优先级的玩家的回复
func (tdr *TaskPlayerReAction) replayAction() int {
	if tdr.waitQueue == nil || tdr.waitQueue.Len() < 1 {
		return int(mahjong.ActionType_enumActionType_SKIP)
	}

	// 返回优先级最高者
	// 对于多人胡牌的情况，需要额外处理
	wi := tdr.waitQueue.Front().Value.(*ReActionQueueItem)
	return wi.replyAction
}

// who 获得最高优先级的玩家
func (tdr *TaskPlayerReAction) who() *PlayerHolder {
	if tdr.waitQueue == nil || tdr.waitQueue.Len() < 1 {
		return nil
	}

	wi := tdr.waitQueue.Front().Value.(*ReActionQueueItem)
	return wi.player
}

// who 获得最高优先级的玩家
func (tdr *TaskPlayerReAction) whoWI() *ReActionQueueItem {
	if tdr.waitQueue == nil || tdr.waitQueue.Len() < 1 {
		return nil
	}

	wi := tdr.waitQueue.Front().Value.(*ReActionQueueItem)
	return wi
}

// actionMeld 获得最高优先级的玩家的操作的MsgMeldTile对象
func (tdr *TaskPlayerReAction) actionMeld() *mahjong.MsgMeldTile {
	if tdr.waitQueue == nil || tdr.waitQueue.Len() < 1 {
		return nil
	}

	wi := tdr.waitQueue.Front().Value.(*ReActionQueueItem)
	return wi.msgMeldTile
}

// completed 完成等待
func (tdr *TaskPlayerReAction) completed(result bool) {
	if tdr.isFinished {
		return
	}

	tdr.isFinished = true

	if tdr.chanWait == nil {
		return
	}

	tdr.chanWait <- result
}

// takeAction 玩家做了选择
func (tdr *TaskPlayerReAction) takeAction(player *PlayerHolder, action int, msgMeldTile *mahjong.MsgMeldTile) {
	var wi = tdr.findWaitQueueItem(player)

	// 玩家不存在
	if wi == nil {
		return
	}

	if tdr.isFinished {
		return
	}

	// 需求变更，动作不是选择了过，那么不考虑过手胡了
	// 备注：2018-3-15 代理杨总和大冯哥都确认，只有选择了“过”才需要过手胡锁定
	if 0 != (wi.actions&int(mahjong.ActionType_enumActionType_WIN_Chuck)) &&
		action == int(mahjong.ActionType_enumActionType_SKIP) {
		// 可以胡牌却选择不胡，在本人重新出牌之前不可以再胡其他人的牌
		player.hStatis.isWinAbleLocked = true

		log.Printf("player %s can win, but skip, now winLockedForWinType:%+v\n",
			player.userID(), player.hStatis.isWinLockedForGreatWin)

		// 需求变更，起手听不胡，只能自摸胡牌了
		if player.hStatis.isRichi {
			player.hStatis.isRichiWinableOnlySelfDrawn = true
		}
	}

	if 0 != (wi.actions&int(mahjong.ActionType_enumActionType_PONG)) &&
		action == int(mahjong.ActionType_enumActionType_SKIP) {
		// 只考虑过手碰/漏碰的情况，吃是不考虑的
		// 可以碰却选择过的人，在本人重新出牌之前不可以再碰其他人的牌
		player.hStatis.pongAbleTileLocked = tdr.actionTile.tileID
	}

	// 增加动作计数器
	player.hStatis.actionCounter++

	// 检查是否需要调整玩家位于队列中的优先级
	// 因为玩家的选择动作可能会导致其优先级变更
	// 例如：如果玩家可以胡牌和碰牌，此时他选择碰牌而不是胡牌
	// 如果队列中有其他玩家也可以胡牌，那么他的优先级就要比其他可以胡牌的玩家的低
	var actionPriority = action2Priority(action)
	var mostPriority = actionsMostPriority(wi.actions)
	if mostPriority > actionPriority {
		// 用户的选择并不是他可以选择的最高优先级的操作，而是低优先级的，因此需要修改其
		// 于队列中的位置
		tdr.changePlayerPosition(wi, actionPriority)
	}

	wi.replyAction = action
	wi.msgMeldTile = msgMeldTile

	// 还有更高级的玩家在等待
	firstWi := tdr.waitQueue.Front().Value.(*ReActionQueueItem)
	if wi != firstWi && firstWi.replyAction == 0 {
		return
	}

	// 没有可胡的等待玩家了，这里需要循环检查，是由于假如多个可以胡，那么虽然位于队首的玩家选择了胡
	// 也得等待其他可胡的玩家做选择后，才可以继续玩下走
	for e := tdr.waitQueue.Front(); e != nil; e = e.Next() {
		wi3 := e.Value.(*ReActionQueueItem)
		if (wi3.actions&int(mahjong.ActionType_enumActionType_WIN_Chuck)) != 0 && wi3.replyAction == 0 {
			log.Println("continue wait winchuck:", wi.player.chairID)
			return
		}
	}

	if tdr.forceWaitAllReply {
		//log.Println("reAction need to wait all players to reply")
		for e := tdr.waitQueue.Front(); e != nil; e = e.Next() {
			wi3 := e.Value.(*ReActionQueueItem)
			if wi3.replyAction == 0 {
				return
			}
		}
	}

	// 玩家是当前优先级最高的玩家，而且没有可胡牌的玩家在等待
	tdr.completed(true)
}

// action2Priority 得到action的优先级
func action2Priority(action int) int {
	var priority = 0
	switch action {
	case int(mahjong.ActionType_enumActionType_CHOW):
		priority = 1
		break
	case int(mahjong.ActionType_enumActionType_PONG):
		priority = 2
		break
	case int(mahjong.ActionType_enumActionType_KONG_Exposed):
		priority = 3
		break
	case int(mahjong.ActionType_enumActionType_WIN_Chuck):
		priority = 4
		break
	}

	return priority
}

// actionsMostPriority 寻找actions集合的最高优先级
func actionsMostPriority(actions int) int {
	if 0 != (actions & int(mahjong.ActionType_enumActionType_WIN_Chuck)) {
		return 4
	}

	if 0 != (actions & int(mahjong.ActionType_enumActionType_KONG_Exposed)) {
		return 3
	}

	if 0 != (actions & int(mahjong.ActionType_enumActionType_PONG)) {
		return 2
	}

	if 0 != (actions & int(mahjong.ActionType_enumActionType_CHOW)) {
		return 1
	}

	return 0
}

// changePlayerPosition 修正玩家在等待队列的位置，例如他本来是可以胡牌的，那么他一开始在队列前面
// 但是由于他选择了过，那么就得把它排到后面，以便让其他玩家有更高的优先级
func (tdr *TaskPlayerReAction) changePlayerPosition(wi *ReActionQueueItem, actionPriority int) {
	// 先移除
	tdr.removeWaitQueueItem(wi)

	var ewi *list.Element

	// 寻找一个优先级比当action优先级小的位置插入
	for e := tdr.waitQueue.Front(); e != nil; e = e.Next() {
		item := e.Value.(*ReActionQueueItem)
		if item.replyAction == 0 {
			if actionsMostPriority(item.actions) <= actionPriority {
				ewi = e
				break
			}
		} else {
			if action2Priority(item.replyAction) <= actionPriority {
				ewi = e
				break
			}
		}
	}

	if ewi == nil {
		tdr.waitQueue.PushBack(wi)
	} else {
		tdr.waitQueue.InsertBefore(wi, ewi)
	}
}

// wait 等待玩家回复
func (tdr *TaskPlayerReAction) wait() bool {
	if tdr.isFinished {
		return false
	}

	result := <-tdr.chanWait

	if result == false {
		return result
	}

	// 如果房间正在解散处理，等待解散结果
	if tdr.s.room.disband != nil {
		result = <-tdr.s.room.disband.chany
	}

	return result
}

// cancel 取消等待，例如房间强制游戏循环结束
func (tdr *TaskPlayerReAction) cancel() {
	tdr.completed(false)
}

func (tdr *TaskPlayerReAction) onPlayerRestore(player *PlayerHolder) {
	var wi = tdr.findWaitQueueItem(player)

	// 玩家不存在
	if wi == nil {
		return
	}

	if wi.replyAction != 0 {
		return
	}

	qaIndex := tdr.s.room.qaIndex
	player.expectedAction = wi.actions

	msgAllowedAction := serializeMsgAllowedForDiscardResponse(player, qaIndex, tdr.actionTile, tdr.actionPlayer)
	player.sendReActoinAllowedMsg(msgAllowedAction)
}
