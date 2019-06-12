package pddz

var (
	// InvalidCard 无效的Card对象
	InvalidCard *Card
	// EmptyCard 一个特殊的牌，用于发牌时表示无牌可抽
	EmptyCard *Card
)

// Card 麻将牌对象
type Card struct {
	cardID int
	drawBy string
}

func init() {
	t := Card{cardID: CARDMAX}
	InvalidCard = &t

	EmptyCard = &Card{cardID: (CARDMAX + 1)}
}

func cardRank(cardID int) int {
	return cardID / 4
}

func cardSuit(cardID int) int {
	return cardID % 4
}
