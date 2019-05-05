package dfmahjong

// TaskFirstReadyHand 等待玩家起手听牌
type TaskFirstReadyHand struct {
	waitQueue []*TFActionQueueItem

	chanWait   chan bool
	isFinished bool

	s *SPlaying
}

// TFActionQueueItem 队列对象
type TFActionQueueItem struct {
	player  *PlayerHolder
	isRichi bool
	isReply bool
}

// newTaskFirstReadyHand 构造一个task用于等待玩家起手听牌
func newTaskFirstReadyHand(players []*PlayerHolder) *TaskFirstReadyHand {
	t := &TaskFirstReadyHand{}
	t.chanWait = make(chan bool, 1) // buffered channel,2 slots
	t.waitQueue = make([]*TFActionQueueItem, len(players))

	for i := range t.waitQueue {
		t.waitQueue[i] = &TFActionQueueItem{player: players[i], isRichi: false, isReply: false}
	}

	return t
}

// findWiByPlayer 根据玩家对象找到队列项
func (tf *TaskFirstReadyHand) findWiByPlayer(player *PlayerHolder) *TFActionQueueItem {
	for _, wi := range tf.waitQueue {
		if wi.player == player {
			return wi
		}
	}
	return nil
}

// takeAction 玩家做了选择
func (tf *TaskFirstReadyHand) takeAction(player *PlayerHolder, richi bool) {
	wi := tf.findWiByPlayer(player)

	if wi == nil || wi.isReply {
		return
	}

	wi.isReply = true
	wi.isRichi = richi
	// 增加动作计数器
	wi.player.hStatis.actionCounter++
	var isFinished = true
	for _, wi := range tf.waitQueue {
		if !wi.isReply {
			isFinished = false
			break
		}
	}

	if isFinished {
		tf.completed(true)
	}
}

// completed 完成等待
func (tf *TaskFirstReadyHand) completed(result bool) {
	if tf.isFinished {
		return
	}

	tf.isFinished = true
	if tf.chanWait == nil {
		return
	}

	tf.chanWait <- result
}

// wait 等待玩家回复
func (tf *TaskFirstReadyHand) wait() bool {
	if tf.isFinished {
		return false
	}

	result := <-tf.chanWait

	if result == false {
		return result
	}

	// 如果房间正在解散处理，等待解散结果
	if tf.s.room.disband != nil {
		result = <-tf.s.room.disband.chany
	}

	return result
}

// cancel 取消等待，例如房间强制游戏循环结束
func (tf *TaskFirstReadyHand) cancel() {
	tf.completed(false)
}

// onPlayerRestore 玩家掉线恢复时，如果其位于起手听列表，而且尚未回复，那么重新给他发起收听请求
func (tf *TaskFirstReadyHand) onPlayerRestore(player *PlayerHolder) {
	wi := tf.findWiByPlayer(player)

	if wi == nil || wi.isReply {
		return
	}

	qaIndex := tf.s.room.qaIndex
	msgAllowedAction := serializeMsgAllowedForRichi(tf.s, player, qaIndex)
	player.expectedAction = int(msgAllowedAction.GetAllowedActions())
	player.sendActoinAllowedMsg(msgAllowedAction)
}
