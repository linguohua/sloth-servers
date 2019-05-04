package dfmahjong

import (
	"mahjong"
)

func onMessageWinSelfDraw(s *SPlaying, player *PlayerHolder, msg *mahjong.MsgPlayerAction) {
	if s.taskPlayerAction == nil {
		s.cl.Println("OnMessageWinSelfDraw, TaskPlayerAction must not be null")
		return
	}

	if s.taskPlayerAction.player != player {
		s.cl.Println("OnMessageWinSelfDraw, not expected player")
		return
	}

	// 检查胡牌合法性
	if !s.tileMgr.winSelfDrawAble(player) {
		s.cl.Printf("OnMessageWinSelfDraw, userId:%s, not win able", player.userID())
		return
	}

	s.taskPlayerAction.takeAction(player, int(mahjong.ActionType_enumActionType_WIN_SelfDrawn), TILEMAX)
}
