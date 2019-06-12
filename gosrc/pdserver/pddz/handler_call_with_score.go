package pddz

import (
	"pokerface"
)

func onMessageCallWithScore(s *SPlaying, player *PlayerHolder, msg *pokerface.MsgPlayerAction) {
	if s.taskPlayerAction == nil {
		s.cl.Panicln("onMessageCallWithScore error, taskPlayerAction is nil, chair:", player.chairID)
		return
	}

	s.taskPlayerAction.flags = int(msg.GetFlags())

	flags := (1 << uint(s.taskPlayerAction.flags))
	if (flags & s.taskPlayerAction.allowFlags) == 0 {
		s.cl.Panicf("onMessageCallWithScore error, input flag %d not include in allowed-flags:%d, chair-id:%d\n",
			s.taskPlayerAction.flags, s.taskPlayerAction.allowFlags, player.chairID)
		return
	}

	s.taskPlayerAction.takeAction(player, ActionType_enumActionType_CallWithScore, nil)
}
