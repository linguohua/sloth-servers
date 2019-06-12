package pddz

import (
	"pokerface"
)

// CardHand 一手牌
type CardHand struct {
	ht    CardHandType
	cards []*Card
}

func newCardHand(ht CardHandType, cards []*Card) *CardHand {
	cardHand := &CardHand{}
	cardHand.ht = ht
	cardHand.cards = cards

	return cardHand
}

// cardHand2MsgCardHand server内部的cardHand转换为MsgMeldCard，以便发送给客户端
func (ch *CardHand) cardHand2MsgCardHand() *pokerface.MsgCardHand {
	msgCardHand := &pokerface.MsgCardHand{}
	var cardHandType = int32(ch.ht)
	msgCardHand.CardHandType = &cardHandType

	var cards = make([]int32, len(ch.cards))
	for i, c := range ch.cards {
		cards[i] = int32(c.cardID)
	}

	msgCardHand.Cards = cards
	return msgCardHand
}
