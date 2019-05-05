package dfmahjong

import (
	"mahjong"
)

func onMessageKongConcealed(s *SPlaying, player *PlayerHolder, msg *mahjong.MsgPlayerAction) {
	var tileID = int(msg.GetTile())

	if s.taskPlayerAction == nil {
		s.cl.Println("OnMessageKongConcealed, TaskPlayerAction must not be null")
		return
	}

	if s.taskPlayerAction.player != player {
		s.cl.Println("OnMessageKongConcealed, not expected player")
		return
	}

	// 检查杠牌合法性
	if !s.tileMgr.kongConcealedAble(player, tileID) {
		s.cl.Printf("OnMessageKongConcealed, userId:%s, not kong-concealed able", player.userID())
		return
	}

	s.taskPlayerAction.takeAction(player, int(mahjong.ActionType_enumActionType_KONG_Concealed), tileID)

}
