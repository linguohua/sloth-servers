package prunfast

import (
	"pokerface"
)

func onMessageDiscard(s *SPlaying, player *PlayerHolder, msg *pokerface.MsgPlayerAction) {
	if s.taskPlayerAction != nil || s.taskDiscardReAction != nil {
		cards := msg.Cards
		var msgCardHand *pokerface.MsgCardHand
		msgCardHand, ok := s.cardMgr.cardsDiscardAble(player, cards)

		if !ok {
			s.cl.Panicln("onMessageDiscard error, cards not discardAble, player chair:", player.chairID)
			return
		}

		if s.taskPlayerAction != nil {
			s.taskPlayerAction.takeAction(player, ActionType_enumActionType_DISCARD, msgCardHand)
		} else if s.taskDiscardReAction != nil {
			prevCardHand := s.taskDiscardReAction.prevCardHand
			if !isMsgCardHandGreatThan(prevCardHand, msgCardHand) {
				s.cl.Panicln("onMessageDiscard error, msgCardHand should great than prev-cardhand, player chair:", player.chairID)
				return
			}

			if prevCardHand.ht == CardHandType_Single && prevCardHand.cards[0].cardID/4 == (AH/4) {
				if player.cards.hasCardInHand(R2H) {
					if len(cards) != 1 || cards[0] != int32(R2H) {
						// 上手打出ACE，本玩家手上有2，必须打2，不能打炸弹
						s.cl.Panicln("onMessageDiscard error, must discard R2H, player chair:", player.chairID)
						return
					}
				}
			}

			s.taskDiscardReAction.takeAction(player, ActionType_enumActionType_DISCARD, msgCardHand)
		}
	} else {
		s.cl.Panicln("onMessageDiscard error, taskPlayerAction and taskDiscardReAction both are nil, player chair:", player.chairID)
	}
}
