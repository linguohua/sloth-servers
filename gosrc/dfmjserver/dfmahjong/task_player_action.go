package dfmahjong

import (
	"mahjong"
)

// TaskPlayerAction 等待一个玩家操作
type TaskPlayerAction struct {
	player      *PlayerHolder
	flags       int
	action      int
	waitActions int
	tileID      int
	// tile           *Tile
	isFinished     bool
	isForFinalDraw bool
	chanWait       chan bool
	s              *SPlaying
}

// newTaskPlayerAction 新建等待任务
func newTaskPlayerAction(player *PlayerHolder, waitActions int) *TaskPlayerAction {
	t := &TaskPlayerAction{}
	t.flags = 0
	t.player = player
	t.tileID = TILEMAX
	t.chanWait = make(chan bool, 1) // buffered channel,1 slots
	t.waitActions = waitActions
	return t
}

// takeAction 玩家做了选择
func (tpa *TaskPlayerAction) takeAction(player *PlayerHolder, action int, tileID int) {
	if tpa.player != player {
		return
	}

	if tpa.isFinished {
		return
	}
	tpa.action = action
	tpa.tileID = tileID

	if 0 != (tpa.waitActions&int(mahjong.ActionType_enumActionType_WIN_SelfDrawn)) &&
		action != int(mahjong.ActionType_enumActionType_WIN_SelfDrawn) {
		// 可以胡牌却选择不胡，在本人重新出牌之前不可以再胡其他人的牌
		// 需求变更：过手胡也包括自己自摸而不胡的情形
		player.hStatis.isWinAbleLocked = true

		// 需求变更，起手听不胡，只能自摸胡牌了
		if player.hStatis.isRichi {
			player.hStatis.isRichiWinableOnlySelfDrawn = true
		}
	}

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
		// var newDraw = tpa.s.lctx.isSelfDraw(player)
		// if player == tpa.s.room.bankerPlayer() && player.hStatis.actionCounter == 0 {
		// 	newDraw = true
		// }
		var qaIndex = tpa.s.room.nextQAIndex()
		actions := tpa.waitActions

		player.expectedAction = actions
		if tpa.isForFinalDraw {
			msgAllowPlayerAction2 := serializeMsgAllowedForDiscard2Opponent(player, qaIndex, actions)
			player.sendActoinAllowedMsg(msgAllowPlayerAction2)
		} else {

			msgAllowPlayerAction := serializeMsgAllowedForDiscard(tpa.s, player, actions, qaIndex)
			player.sendActoinAllowedMsg(msgAllowPlayerAction)
		}

	} else {
		msgAllowPlayerAction2 := serializeMsgAllowedForDiscard2Opponent(tpa.player, qaIndex,
			int(mahjong.ActionType_enumActionType_DISCARD))
		player.sendActoinAllowedMsg(msgAllowPlayerAction2)
	}
}
