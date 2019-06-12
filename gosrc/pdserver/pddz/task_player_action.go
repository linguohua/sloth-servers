package pddz

import (
	"pokerface"
)

// TaskPlayerAction 等待一个玩家操作
type TaskPlayerAction struct {
	player     *PlayerHolder
	flags      int
	allowFlags int

	replyAction   ActionType
	replyCardHand *pokerface.MsgCardHand

	applyCardHand *CardHand

	waitActions int

	isFinished bool
	chanWait   chan bool
	s          *SPlaying

	actionMsgForRestore *pokerface.MsgAllowPlayerAction
}

// newTaskPlayerAction 新建等待任务
func newTaskPlayerAction(player *PlayerHolder, waitActions int) *TaskPlayerAction {
	t := &TaskPlayerAction{}
	t.flags = 0
	t.player = player

	player.hStatis.lastExpectedActions = waitActions
	player.hStatis.lastExpectedType = 0

	t.chanWait = make(chan bool, 1) // buffered channel,1 slots
	t.waitActions = waitActions

	return t
}

// takeAction 玩家做了选择
func (tpa *TaskPlayerAction) takeAction(player *PlayerHolder, action ActionType, msgCardHand *pokerface.MsgCardHand) {
	if tpa.player != player {
		return
	}

	if tpa.isFinished {
		return
	}

	tpa.replyAction = action
	tpa.replyCardHand = msgCardHand

	// 增加动作计数器
	player.hStatis.actionCounter++
	tpa.completed(true)
}

// completed 完成等待
func (tpa *TaskPlayerAction) completed(result bool) {
	if tpa.isFinished {
		return
	}

	tpa.isFinished = true
	if tpa.chanWait == nil {
		return
	}

	tpa.chanWait <- result
}

// wait 等待
func (tpa *TaskPlayerAction) wait() bool {
	if tpa.isFinished {
		return false
	}

	result := <-tpa.chanWait

	if result == false {
		return result
	}

	// 如果房间正在解散处理，等待解散结果
	if tpa.s.room.disband != nil {
		result = <-tpa.s.room.disband.chany
	}

	return result
}

// cancel 取消等待
func (tpa *TaskPlayerAction) cancel() {
	tpa.completed(false)
}

// onPlayerRestore 玩家重入恢复
func (tpa *TaskPlayerAction) onPlayerRestore(player *PlayerHolder) {
	qaIndex := tpa.s.room.qaIndex

	if player == tpa.player {

		if tpa.actionMsgForRestore != nil {
			player.sendActoinAllowedMsg(tpa.actionMsgForRestore)
			return
		}

		actions := tpa.waitActions
		player.expectedAction = actions
		msgAllowPlayerAction := serializeMsgAllowedForDiscard(tpa.s, player, actions, qaIndex)
		player.sendActoinAllowedMsg(msgAllowPlayerAction)
	} else {
		msgAllowPlayerAction2 := serializeMsgAllowedForDiscard2Opponent(tpa.player, qaIndex)
		player.sendActoinAllowedMsg(msgAllowPlayerAction2)
	}
}
