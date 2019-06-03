package zjmahjong

import (
	"mahjong"
)

func onMessageWinChuck(s *SPlaying, player *PlayerHolder, msg *mahjong.MsgPlayerAction) {
	var tileID = int(msg.GetTile())

	if s.taskDiscardReAction == nil {
		s.cl.Printf("OnMessageWinChuck, userId:%s, taskDiscardReAction is null\n", player.userID())
		return
	}

	if s.taskDiscardReAction.actionTile.tileID != tileID {
		s.cl.Printf("OnMessageWinChuck, userId:%s, not expected target tile %d\n", player.userID(), tileID)
		return
	}

	if !s.taskDiscardReAction.isExpectedPlayerAction(player, int(mahjong.ActionType_enumActionType_WIN_Chuck)) {
		s.cl.Printf("OnMessageWinChuck, userId:%s, not expected player\n", player.userID())
		return
	}

	// 检查胡牌合法性
	if !s.tileMgr.winChuckAble(player, s.taskDiscardReAction) {
		s.cl.Printf("OnMessageWinChuck, userId:%s, not winable\n", player.userID())
		return
	}

	// 异步等待队列完成
	var taskDiscardReAction = s.taskDiscardReAction
	taskDiscardReAction.takeAction(player, int(mahjong.ActionType_enumActionType_WIN_Chuck), nil)

}
