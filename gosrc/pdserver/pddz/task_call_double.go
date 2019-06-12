package pddz

import (
	"container/list"
)

// TaskCallDouble 加注
type TaskCallDouble struct {
	isFinished bool
	chanWait   chan bool
	s          *SPlaying

	waitQueue *list.List
}

// TaskCallDoubleQueueItem 等待队列项
type TaskCallDoubleQueueItem struct {
	player     *PlayerHolder
	reply      bool
	waitAction int
	replyFlags int
}

// newTaskCallDouble 新建等待任务
func newTaskCallDouble(s *SPlaying) *TaskCallDouble {
	t := &TaskCallDouble{}

	t.chanWait = make(chan bool, 1) // buffered channel,1 slots
	t.s = s

	t.waitQueue = list.New()
	for _, p := range s.players {
		ti := &TaskCallDoubleQueueItem{}
		ti.player = p
		ti.waitAction = int(ActionType_enumActionType_CallDouble)

		t.waitQueue.PushBack(ti)
	}

	return t
}

// findWaitQueueItem 根据player找到wait item
func (tes *TaskCallDouble) findWaitQueueItem(player *PlayerHolder) *TaskCallDoubleQueueItem {
	for e := tes.waitQueue.Front(); e != nil; e = e.Next() {
		qi := e.Value.(*TaskCallDoubleQueueItem)
		if qi.player == player {
			return qi
		}
	}
	return nil
}

// takeAction 玩家做了选择
func (tes *TaskCallDouble) takeAction(player *PlayerHolder, action int, replyFlags int) {
	wi := tes.findWaitQueueItem(player)
	if wi == nil {
		tes.s.cl.Printf("player %s not in TaskCallDouble queue", player.userID())
		return
	}

	if wi.reply {
		tes.s.cl.Printf("player %s ha already reply TaskCallDouble queue", player.userID())
		return
	}

	wi.reply = true
	wi.replyFlags = replyFlags

	// 增加动作计数器
	player.hStatis.actionCounter++

	allReply := true
	for e := tes.waitQueue.Front(); e != nil; e = e.Next() {
		qi := e.Value.(*TaskCallDoubleQueueItem)
		if !qi.reply {
			allReply = false
			break
		}
	}

	if allReply {
		tes.completed(true)
	}
}

// completed 完成等待
func (tes *TaskCallDouble) completed(result bool) {
	if tes.isFinished {
		return
	}

	tes.isFinished = true
	if tes.chanWait == nil {
		return
	}

	tes.chanWait <- result
}

// wait 等待
func (tes *TaskCallDouble) wait() bool {
	if tes.isFinished {
		return false
	}

	result := <-tes.chanWait

	if result == false {
		return result
	}

	// 如果房间正在解散处理，等待解散结果
	if tes.s.room.disband != nil {
		result = <-tes.s.room.disband.chany
	}

	return result
}

// cancel 取消等待
func (tes *TaskCallDouble) cancel() {
	tes.completed(false)
}

// onPlayerRestore 玩家重入恢复
func (tes *TaskCallDouble) onPlayerRestore(player *PlayerHolder) {
	var wi = tes.findWaitQueueItem(player)
	qaIndex := tes.s.room.qaIndex

	// 玩家不存在，或者已经回复了
	if wi == nil || wi.reply {
		return
	}

	player.expectedAction = wi.waitAction
	msgAllowedAction := serializeMsgAllowedForDiscard(tes.s, player, wi.waitAction, qaIndex)
	player.sendActoinAllowedMsg(msgAllowedAction)
}
