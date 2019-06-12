package pddz

import (
	"pokerface"
)

func onMessageCallDouble(s *SPlaying, player *PlayerHolder, msg *pokerface.MsgPlayerAction) {
	if s.taskCallDoulbe == nil {
		s.cl.Panicln("onMessageCall error, taskCallDoulbe is nil, chair:", player.chairID)
		return
	}

	// s.taskPlayerAction.flags = int(msg.GetFlags())
	if msg.GetFlags() != 0 {
		player.hStatis.isCallDouble = true
	}

	s.taskCallDoulbe.takeAction(player, int(ActionType_enumActionType_CallDouble), int(msg.GetFlags()))

	var msgActionNotifyResult = serializeMsgActionResultNotifyForNoCard(ActionType_enumActionType_CallDouble, player)
	cardInWall32 := int32(0) // 0表示玩家不加注
	if player.hStatis.isCallDouble {
		cardInWall32 = int32(100) // 100表示玩家加注
	}
	msgActionNotifyResult.CardsInWall = &cardInWall32

	// 发送结果给所有其他用户
	for _, p := range s.players {
		p.sendActionResultNotify(msgActionNotifyResult)
	}
}
