package pddz

import (
	"pokerface"
)

func onMessageSkip(s *SPlaying, player *PlayerHolder, msg *pokerface.MsgPlayerAction) {
	if s.taskDiscardReAction == nil {
		s.cl.Panicln("onMessageSkip error, taskDiscardReAction is nil, chair:", player.chairID)
		return
	}

	if !s.taskDiscardReAction.isExpectedPlayerAction(player, ActionType_enumActionType_SKIP) {
		s.cl.Panicln("onMessageSkip error, not expected player or action, chair:", player.chairID)
		return
	}

	s.taskDiscardReAction.takeAction(player, ActionType_enumActionType_SKIP, nil)
}
