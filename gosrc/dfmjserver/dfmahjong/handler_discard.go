package dfmahjong

import (
	"mahjong"
)

func onMessageDiscard(s *SPlaying, player *PlayerHolder, msg *mahjong.MsgPlayerAction) {
	var tileID = int(msg.GetTile())

	if s.taskPlayerAction == nil {
		s.cl.Println("OnMessageDiscard, TaskPlayerAction must not be null")
		return
	}

	if s.taskPlayerAction.player != player {
		s.cl.Println("OnMessageDiscard, not expected player")
	}

	if !s.tileMgr.discardAble(player, tileID) {
		s.cl.Printf("OnMessageDiscard, userId:%s, not discardAble\n", player.userID())
		return
	}

	// 保存flags，因为庄家出第一个牌时可能选择起手听牌
	s.taskPlayerAction.flags = int(msg.GetFlags())

	s.taskPlayerAction.takeAction(player, int(mahjong.ActionType_enumActionType_DISCARD), tileID)
}
