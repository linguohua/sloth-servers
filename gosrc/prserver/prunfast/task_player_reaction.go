package prunfast

import (
	"pokerface"
)

// TaskPlayerReAction 等待其他玩家完成操作
type TaskPlayerReAction struct {
	prevPlayer   *PlayerHolder
	prevCardHand *CardHand

	waitPlayer  *PlayerHolder
	waitActions int

	replyAction      ActionType
	replyMsgCardHand *pokerface.MsgCardHand

	applyCardHand *CardHand

	isFinished        bool
	chanWait          chan bool
	forceWaitAllReply bool

	s *SPlaying
}

func analyseTaskPlayerReAction(s *SPlaying, player *PlayerHolder, prevDiscardedPlayer *PlayerHolder, prevDiscardedCardHand *CardHand) *TaskPlayerReAction {
	cards := player.cards
	var actions = int(ActionType_enumActionType_SKIP)
	if cards.hasCardHandGreatThan(prevDiscardedCardHand) {
		actions |= int(ActionType_enumActionType_DISCARD)

		// sombodyRemainLessThan4 := false

		// // 只要有一个人的牌少于4张，强制出牌
		// for _, p := range s.players {
		// 	if p != player && p.cards.cardCountInHand() < 4 {
		// 		sombodyRemainLessThan4 = true
		// 		break
		// 	}
		// }

		// if !sombodyRemainLessThan4 {
		// 	// 如果上一手是ACE，本手有2必须打2，这个过滤在serializeMsgAllowedForDiscardResponse中处理
		// 	action |= int(ActionType_enumActionType_SKIP)
		// }

		// 需求变更：不管报警与否，玩家都可以选择过
	}

	player.expectedAction = actions

	player.hStatis.lastExpectedActions = actions
	player.hStatis.lastExpectedType = 1

	tdr := &TaskPlayerReAction{}
	tdr.prevPlayer = prevDiscardedPlayer
	tdr.prevCardHand = prevDiscardedCardHand
	tdr.waitPlayer = player
	tdr.waitActions = actions
	tdr.chanWait = make(chan bool, 1) // buffered channel,1 slots

	return tdr
}

// isExpectedPlayerAction 是否正在等待的玩家以及动作
func (tdr *TaskPlayerReAction) isExpectedPlayerAction(player *PlayerHolder, action ActionType) bool {
	if tdr.isFinished {
		return false
	}

	if tdr.waitPlayer != player {
		return false
	}

	return (tdr.waitActions & int(action)) != 0
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
func (tdr *TaskPlayerReAction) takeAction(player *PlayerHolder, action ActionType, msgMeldCardHand *pokerface.MsgCardHand) {
	if tdr.isFinished {
		return
	}

	// 增加动作计数器
	player.hStatis.actionCounter++

	tdr.replyAction = action
	tdr.replyMsgCardHand = msgMeldCardHand

	// 玩家是当前优先级最高的玩家，而且没有可胡牌的玩家在等待
	tdr.completed(true)
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
	qaIndex := tdr.s.room.qaIndex

	// 玩家不存在，或者已经回复了
	if tdr.waitPlayer != player {
		somebody := tdr.waitPlayer
		msgAllowedAction := serializeMsgAllowedForDiscardResponseRestore(somebody, qaIndex, tdr.prevCardHand, tdr.prevPlayer)
		player.sendReActoinAllowedMsg(msgAllowedAction)
		return
	}

	player.expectedAction = tdr.waitActions
	msgAllowedAction := serializeMsgAllowedForDiscardResponse(player, qaIndex, tdr.prevCardHand, tdr.prevPlayer)
	player.sendReActoinAllowedMsg(msgAllowedAction)
}
