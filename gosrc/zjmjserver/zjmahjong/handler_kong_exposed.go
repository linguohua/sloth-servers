package zjmahjong

import (
	"mahjong"
)

func onMessageKongExposed(s *SPlaying, player *PlayerHolder, msg *mahjong.MsgPlayerAction) {
	var tileID = int(msg.GetTile())

	if s.taskDiscardReAction == nil {
		s.cl.Printf("onMessageKongExposed, userId:%s, taskDiscardReAction is null\n", player.userID())
		return
	}

	if s.taskDiscardReAction.actionTile.tileID != tileID {
		s.cl.Printf("onMessageKongExposed, userId:%s, not expected target tile %d\n", player.userID(), tileID)
		return
	}

	if !s.taskDiscardReAction.isExpectedPlayerAction(player, int(mahjong.ActionType_enumActionType_KONG_Exposed)) {
		s.cl.Printf("onMessageKongExposed, userId:%s, not expected player\n", player.userID())
		return
	}

	msgMeld := &mahjong.MsgMeldTile{}
	var meldTyp32 = msg.GetMeldType()
	var meldTile132 = msg.GetMeldTile1()
	msgMeld.Tile1 = &meldTile132
	msgMeld.MeldType = &meldTyp32

	//  检查杠牌合法性
	if !s.tileMgr.kongExposedAble(player, s.taskDiscardReAction, msgMeld) {
		s.cl.Printf("onMessageKongExposed, userId:%s, not kongExposedAble\n", player.userID())
		return
	}

	// 异步等待队列完成
	var taskDiscardReAction = s.taskDiscardReAction
	taskDiscardReAction.takeAction(player, int(mahjong.ActionType_enumActionType_KONG_Exposed), msgMeld)
}
