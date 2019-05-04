package dfmahjong

import (
	"mahjong"
)

func onMessagePong(s *SPlaying, player *PlayerHolder, msg *mahjong.MsgPlayerAction) {
	var tileID = int(msg.GetTile())

	if s.taskDiscardReAction == nil {
		s.cl.Printf("onMessagePong, userId:%s, taskDiscardReAction is null\n", player.userID())
		return
	}

	if s.taskDiscardReAction.actionTile.tileID != tileID {
		s.cl.Printf("onMessagePong, userId:%s, not expected target tile %d\n", player.userID(), tileID)
		return
	}

	if !s.taskDiscardReAction.isExpectedPlayerAction(player, int(mahjong.ActionType_enumActionType_PONG)) {
		s.cl.Printf("onMessagePong, userId:%s, not expected player\n", player.userID())
		return
	}

	msgMeld := &mahjong.MsgMeldTile{}
	var meldTyp32 = msg.GetMeldType()
	var meldTile132 = msg.GetMeldTile1()
	msgMeld.Tile1 = &meldTile132
	msgMeld.MeldType = &meldTyp32

	//  检查碰牌合法性
	if !s.tileMgr.pongAble(player, s.taskDiscardReAction, msgMeld) {
		s.cl.Printf("onMessagePong, userId:%s, not pongAble\n", player.userID())
		return
	}

	// 异步等待队列完成
	var taskDiscardReAction = s.taskDiscardReAction
	taskDiscardReAction.takeAction(player, int(mahjong.ActionType_enumActionType_PONG), msgMeld)
}
