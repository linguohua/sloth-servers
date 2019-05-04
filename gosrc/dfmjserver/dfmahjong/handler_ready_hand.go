package dfmahjong

import (
	"mahjong"
)

func onMessageFirstReadyHand(s *SPlaying, player *PlayerHolder, msg *mahjong.MsgPlayerAction) {
	if s.taskFirstReadyHand == nil {
		s.cl.Printf("OnMessageFirstReadyHand, userId:%s null TaskFirstReadyHand\n", player.userID())
		return
	}

	var tiles = player.tiles
	if !tiles.readyHandAble() {
		s.cl.Printf("OnMessageFirstReadyHand, userId:%s, not ready hand able\n", player.userID())
		return
	}

	var isRichi = msg.GetFlags() == 1
	if isRichi {
		player.hStatis.isRichi = true
		// 发送起手听牌结果给其他用户
		var msgActionNotifyResult = serializeMsgActionResultNotifyForNoTile(int(mahjong.ActionType_enumActionType_FirstReadyHand), player)
		// 发送结果给所有用户
		for _, p := range s.players {
			p.sendActionResultNotify(msgActionNotifyResult)
		}
	}

	s.taskFirstReadyHand.takeAction(player, isRichi)
}
