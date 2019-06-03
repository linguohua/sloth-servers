package zjmahjong

import (
	"mahjong"
)

func onMessageShouldFinalDraw(s *SPlaying, player *PlayerHolder, msg *mahjong.MsgPlayerAction) {
	if s.taskPlayerAction == nil {
		s.cl.Panic("onMessageShouldFinalDraw, taskPlayerAction is null,userId:", player.userID())
		return
	}

	if s.taskPlayerAction.player != player {
		s.cl.Println("onMessageShouldFinalDraw, not expected player")
		return
	}

	// 异步等待队列完成
	var taskPlayerAction = s.taskPlayerAction
	taskPlayerAction.takeAction(player, int(mahjong.ActionType_enumActionType_CustomB), TILEMAX)
}
