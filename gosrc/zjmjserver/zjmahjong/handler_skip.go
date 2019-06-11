package zjmahjong

import (
	"mahjong"
)

func onMessageSkip(s *SPlaying, player *PlayerHolder, msg *mahjong.MsgPlayerAction) {
	if s.taskPlayerAction != nil {
		if s.taskPlayerAction.player != player {
			s.cl.Println("onMessageSkip, not expected player")
			return
		}

		s.taskPlayerAction.takeAction(player, int(mahjong.ActionType_enumActionType_SKIP), TILEMAX)
	} else {
		if s.taskDiscardReAction == nil {
			// 此时可以最高优先级的玩家已经回复，因此taskDiscardReAction可能为空
			s.cl.Println("OnMessageSkip, taskDiscardReAction is null,userId:", player.userID())
			return
		}

		if !s.taskDiscardReAction.isExpectedPlayerAction(player, int(mahjong.ActionType_enumActionType_SKIP)) {
			s.cl.Panic("OnMessageSkip, not expected player, userId:", player.userID())
			return
		}

		// 异步等待队列完成
		var taskDiscardReAction = s.taskDiscardReAction
		taskDiscardReAction.takeAction(player, int(mahjong.ActionType_enumActionType_SKIP), nil)
	}
}
