package zjmahjong

import (
	log "github.com/sirupsen/logrus"
)

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

/*
1码：1、5、9、东
2码：2、6、南、中
3码：3、7、西、发
4码：4、8、北、白
*/
func (t *Tile) horseType() int {
	if t.isSuit() {
		rank := t.tileID%9 + 1
		switch rank {
		case 1, 5, 9:
			return 1
		case 2, 6:
			return 2
		case 3, 7:
			return 3
		case 4, 8:
			return 4
		}
	}

	if t.isHonor() {
		switch t.tileID {
		case TON:
			return 1
		case NAN, CHU:
			return 2
		case SHA, HAT:
			return 3
		case PEI, HAK:
			return 4
		}
	}

	log.Panicln("horseType failed, invalid tileID:", t.tileID)

	// 无效的马类型
	return 0
}

func init() {
	t := Tile{tileID: TILEMAX}
	InvalidTile = &t

	EmptyTile = &Tile{tileID: (TILEMAX + 1)}
}
