package dfmahjong

import (
	"mahjong"
)

func onMessageChow(s *SPlaying, player *PlayerHolder, msg *mahjong.MsgPlayerAction) {
	// var tileID = int(msg.GetTile())

	if s.taskDiscardReAction == nil {
		s.cl.Printf("onMessageChow, userId:%s, taskDiscardReAction is null\n", player.userID())
		return
	}

	// if s.taskDiscardReAction.latestDiscardedTile.tileID != tileID {
	// 	s.cl.Printf("onMessageChow, userId:%d, not expected target tile %d\n", player.userID(), tileID)
	// 	return
	// }

	if !s.taskDiscardReAction.isExpectedPlayerAction(player, int(mahjong.ActionType_enumActionType_CHOW)) {
		s.cl.Printf("onMessageChow, userId:%s, not expected player\n", player.userID())
		return
	}

	msgMeld := &mahjong.MsgMeldTile{}
	var meldTyp32 = msg.GetMeldType()
	var meldTile132 = msg.GetMeldTile1()
	msgMeld.Tile1 = &meldTile132
	msgMeld.MeldType = &meldTyp32

	// 检查吃牌合法性
	if !s.tileMgr.chowAble(player, s.taskDiscardReAction, msgMeld) {
		s.cl.Printf("onMessageChow, userId:%s, not chowAble\n", player.userID())
		return
	}

	// 异步等待队列完成
	var taskDiscardReAction = s.taskDiscardReAction
	taskDiscardReAction.takeAction(player, int(mahjong.ActionType_enumActionType_CHOW), msgMeld)
}
