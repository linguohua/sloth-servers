package dfmahjong

import (
	"mahjong"
)

func onMessageTriplet2Kong(s *SPlaying, player *PlayerHolder, msg *mahjong.MsgPlayerAction) {
	var tileID = int(msg.GetTile())

	if s.taskPlayerAction == nil {
		s.cl.Println("OnMessageTriplet2Kong, TaskPlayerAction must not be null")
		return
	}

	if s.taskPlayerAction.player != player {
		s.cl.Println("OnMessageTriplet2Kong, not expected player")
		return
	}

	// 检查加杠合法性
	if !s.tileMgr.triplet2KongAble(player, tileID) {
		s.cl.Printf("OnMessageTriplet2Kong, userId:%s, not win able", player.userID())
		return
	}

	s.taskPlayerAction.takeAction(player, int(mahjong.ActionType_enumActionType_KONG_Triplet2), tileID)
}
