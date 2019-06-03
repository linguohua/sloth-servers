package zjmahjong

var (
	// InvalidTile 无效的Tile对象
	InvalidTile *Tile
	// EmptyTile 一个特殊的牌，用于发牌时表示无牌可抽
	EmptyTile *Tile
)

// Tile 麻将牌对象
type Tile struct {
	tileID int
	drawBy string
}

// isFlower 是否花牌
func (t *Tile) isFlower() bool {
	return t.tileID >= FlowerBegin &&
		t.tileID <= WINTER
}

// isWind 是否风牌
func (t *Tile) isWind() bool {
	return t.tileID >= TON &&
		t.tileID <= PEI
}

// isDragon 是否箭牌
func (t *Tile) isDragon() bool {
	return t.tileID >= HAK &&
		t.tileID <= CHU
}

// isHonor 是否字牌
func (t *Tile) isHonor() bool {
	return t.tileID >= TON &&
		t.tileID <= CHU
}

// isSuit 是否数牌
func (t *Tile) isSuit() bool {
	return t.tileID >= MAN &&
		t.tileID <= SOU9
}

// suitType 数牌类型，必须调用isSuit
// 确保牌是数牌
func (t *Tile) suitType() int {
	if t.tileID <= MAN9 {
		return MAN
	} else if t.tileID <= PIN9 {
		return PIN
	} else if t.tileID <= SOU9 {
		return SOU
	}

	return 0
}

func init() {
	t := Tile{tileID: TILEMAX}
	InvalidTile = &t

	EmptyTile = &Tile{tileID: (TILEMAX + 1)}
}
