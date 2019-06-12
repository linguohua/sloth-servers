package pddz

import (
	"pokerface"
)

func onMessageRob(s *SPlaying, player *PlayerHolder, msg *pokerface.MsgPlayerAction) {
	if s.taskPlayerAction == nil {
		s.cl.Panicln("onMessageRob error, taskPlayerAction is nil, chair:", player.chairID)
		return
	}

	s.taskPlayerAction.flags = int(msg.GetFlags())

	s.taskPlayerAction.takeAction(player, ActionType_enumActionType_Rob, nil)
}
