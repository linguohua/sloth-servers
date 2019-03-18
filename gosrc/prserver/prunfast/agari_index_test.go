package prunfast

import "testing"
import "pokerface"

func testCreatePlayerCards(hai []int) *PlayerCards {
	pc := newPlayerCards(nil)

	for _, h := range hai {
		c := &Card{}
		c.cardID = h
		pc.addHandCard(c)
	}

	return pc
}

func testCreateMsgCardHand(hai []int) *pokerface.MsgCardHand {
	hai32 := make([]int32, len(hai))
	for i, v := range hai {
		hai32[i] = int32(v)
	}

	return agariConvertMsgCardHand(hai32)
}

func TestSingle(t *testing.T) {
	var hai = []int{
		R2H, R2D,
		R3H, R4H,
		R5H, R6H,
		R7H, R8H,
		R9H, R10H,
	}

	var hai2 = []int{
		R3H,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)

	if msgCardHand.GetCardHandType() != int32(CardHandType_Single) {
		t.Error("test single failed")
	}

	if !pc.hasMsgCardHandGreatThan(msgCardHand) {
		t.Error("test single failed")
	}
}

func TestFlush(t *testing.T) {
	var hai = []int{
		R2H, R2D,
		R3H, R4H,
		R5H, R6H,
		R7H, R8H,
		R9H, R10D,
		JH, QH, KH,
		AH,
	}

	var hai2 = []int{
		R9H,
		R10H,
		JH,
		QH,
		KH,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)
	if msgCardHand.GetCardHandType() != int32(CardHandType_Flush) {
		t.Error("test flush failed")
	}

	if !pc.hasMsgCardHandGreatThan(msgCardHand) {
		t.Error("test flush failed")
	}
}

// func TestFlush4(t *testing.T) {
// 	var hai = []int{
// 		R2H, R2D,
// 		R3H, R4H,
// 		R5H, R6H,
// 		R7H, R8H,
// 		R9H, R10D,
// 		JH, QH, KH,
// 		AH,
// 	}

// 	var hai2 = []int{
// 		AH,
// 		R2H,
// 		R3H, R4H,
// 		R5H,
// 	}

// 	pc := testCreatePlayerCards(hai)
// 	msgCardHand := testCreateMsgCardHand(hai2)
// 	if msgCardHand.GetCardHandType() != int32(CardHandType_Flush) {
// 		t.Error("test flush failed")
// 	}

// 	println("first flush card:", msgCardHand.Cards[0])
// 	if !pc.hasMsgCardHandGreatThan(msgCardHand) {
// 		t.Error("test flush failed")
// 	}
// }

func TestFlush2(t *testing.T) {
	var hai = []int{
		R2H, R2D,
		R3H, R4H,
		R5H, R6H,
		R7H, R8H,
		R9H, R9D,
		JH, QH, KH,
		AH,
	}

	var hai2 = []int{
		R9H,
		R10H,
		JH,
		QH,
		KH,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)
	if msgCardHand.GetCardHandType() != int32(CardHandType_Flush) {
		t.Error("test flush2 failed")
	}

	if pc.hasMsgCardHandGreatThan(msgCardHand) {
		t.Error("test flush2 failed")
	}
}

func TestFlush3(t *testing.T) {
	var hai = []int{
		R2H, R2D,
		R3H, R4H,
		R5H, R5S,
		R7H, R8H,
		R9H, R10D,
		JH, QH, KH,
		AH,
	}

	var hai2 = []int{
		R3H, R4H,
		R5H, R6S,
		R7H,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)

	if msgCardHand.GetCardHandType() != int32(CardHandType_Flush) {
		t.Error("test flush failed")
	}

	if !pc.hasMsgCardHandGreatThan(msgCardHand) {
		t.Error("test flush failed")
	}
}

func TestTriplet(t *testing.T) {
	var hai = []int{
		R2H, R2D,
		R3H, R4H,
		R5H, R5S,
		R7H, R8H,
		R9H, R10D,
		JH, QH, KH,
		AH,
	}

	var hai2 = []int{
		R3H,
		R3D, R3S,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)
	if msgCardHand.GetCardHandType() != int32(CardHandType_Triplet) {
		t.Error("test triplet failed")
	}

	if pc.hasMsgCardHandGreatThan(msgCardHand) {
		t.Error("test triplet failed")
	}
}

func TestTriplet2(t *testing.T) {
	var hai = []int{
		R2H, R2D,
		R3H, R4H,
		R5H, R5S,
		R5D, R8H,
		R9H, R10D,
		JH, QH, KH,
		AH,
	}

	var hai2 = []int{
		R3H,
		R3D, R3S,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)
	if msgCardHand.GetCardHandType() != int32(CardHandType_Triplet) {
		t.Error("test triplet failed")
	}

	if !pc.hasMsgCardHandGreatThan(msgCardHand) {
		t.Error("test triplet failed")
	}
}

func TestTriplet2x(t *testing.T) {
	var hai = []int{
		R2H, R2D,
		R3H, R4H,
		R5H, R5S,
		R5D, R8H,
		R9H, R10D,
		JH, QH, KH,
		AH,
	}

	var hai2 = []int{
		R3H,
		R3D, R3S,
		R4H,
		R4D, R4S,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)
	if msgCardHand.GetCardHandType() != int32(CardHandType_Triplet2X) {
		t.Error("test TestTriplet2x failed")
	}

	if pc.hasMsgCardHandGreatThan(msgCardHand) {
		t.Error("test TestTriplet2x failed")
	}
}

func TestTriplet2x2(t *testing.T) {
	var hai = []int{
		R2H, R2D,
		R3H, R4H,
		R5H, R5S,
		R5D, R6H,
		R6S, R6D,
		JH, QH, KH,
		AH,
	}

	var hai2 = []int{
		R3H,
		R3D, R3S,
		R4H,
		R4D, R4S,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)
	if msgCardHand.GetCardHandType() != int32(CardHandType_Triplet2X) {
		t.Error("test TestTriplet2x2 failed")
	}
	// println("ct:", msgCardHand.GetCardHandType())

	if !pc.hasMsgCardHandGreatThan(msgCardHand) {
		t.Error("test TestTriplet2x2 failed")
	}
}

func TestTripletPair(t *testing.T) {
	var hai = []int{
		R2H, R2D,
		R3H, R4H,
		R5H, R5S,
		R5D, R7H,
		R8S, R9D,
		JH, QH, KH,
		AH,
	}

	var hai2 = []int{
		R3H,
		R3D, R3S,
		R4H,
		R4D,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)
	if msgCardHand.GetCardHandType() != int32(CardHandType_TripletPair) {
		t.Error("test TripletPair failed")
	}

	if !pc.hasMsgCardHandGreatThan(msgCardHand) {
		t.Error("test TripletPair failed")
	}
}

func TestTripletPair2(t *testing.T) {
	var hai = []int{
		R2H,
		R3H, R4H,
		R5H, R5S,
		R5D, R5H,
		R5S, R9D,
		JH, QH, KH,
		AH,
	}

	var hai2 = []int{
		R3H,
		R3D, R3S,
		R4H,
		R4D,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)
	if msgCardHand.GetCardHandType() != int32(CardHandType_TripletPair) {
		t.Error("test TripletPair2 failed")
	}

	if !pc.hasMsgCardHandGreatThan(msgCardHand) {
		t.Error("test TripletPair2 failed")
	}
}

func TestTriplet2X2Pair(t *testing.T) {
	var hai = []int{
		R5H, R5S,
		R5D, R6H,
		R6S, R6D,
		QH, QC, QH,
		KH, KC, KS,
		AH,
	}

	var hai2 = []int{
		R3H,
		R3D, R3S,
		R4H,
		R4D, R4S,
		AH, AC,
		QH, QC,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)
	if msgCardHand.GetCardHandType() != int32(CardHandType_Triplet2X2Pair) {
		t.Error("test Triplet2x2Pair failed")
	}

	if !pc.hasMsgCardHandGreatThan(msgCardHand) {
		t.Error("test Triplet2x2Pair failed")
	}
}

func TestTriplet2X2Pair2(t *testing.T) {
	var hai = []int{
		R5H, R5S,
		R6H,
		R6S,
		R7H,
		R7S, R7D,
		JH, JD, JS,
		QH, QC, QH,
		KH, KC, KS,
	}

	var hai2 = []int{
		R3H,
		R3D, R3S,
		R4H,
		R4D, R4S,
		R5H,
		R5D, R5S,
		AH, AC,
		QH, QC,
		KH, KC,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)
	if msgCardHand.GetCardHandType() != int32(CardHandType_Triplet2X2Pair) {
		t.Error("test Triplet2x2Pair2 failed")
	}

	if !pc.hasMsgCardHandGreatThan(msgCardHand) {
		t.Error("test Triplet2x2Pair2 failed")
	}
}

func TestBomb(t *testing.T) {
	var hai = []int{
		JH, QH, KH,
		AH, AD, AS,
	}

	var hai2 = []int{
		R4H, R4D, R4S, R4C,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)
	if msgCardHand.GetCardHandType() != int32(CardHandType_Bomb) {
		t.Error("test TestBomb failed")
	}

	if !pc.hasMsgCardHandGreatThan(msgCardHand) {
		t.Error("test TestBomb failed")
	}
}

func TestBomb3(t *testing.T) {
	var hai = []int{
		AH, AD, AS,
	}

	var hai2 = []int{
		R4H, R4D, R4S, R4C,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)
	if msgCardHand.GetCardHandType() != int32(CardHandType_Bomb) {
		t.Error("test TestBomb failed")
	}

	if !pc.hasMsgCardHandGreatThan(msgCardHand) {
		t.Error("test TestBomb failed")
	}
}

func TestBomb5(t *testing.T) {
	var hai = []int{
		R4H, R4D, R4S, R4C,
	}

	var hai2 = []int{
		AH, AD, AS,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)
	if msgCardHand.GetCardHandType() != int32(CardHandType_Bomb) {
		t.Error("test TestBomb failed")
	}

	if pc.hasMsgCardHandGreatThan(msgCardHand) {
		t.Error("test TestBomb failed")
	}
}

func TestBomb6(t *testing.T) {
	var hai = []int{
		R3H, R3D, R3S, R3C,
	}

	var hai2 = []int{
		AH, AD,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)
	if msgCardHand.GetCardHandType() != int32(CardHandType_Pair) {
		t.Error("test TestBomb failed")
	}

	if !pc.hasMsgCardHandGreatThan(msgCardHand) {
		t.Error("test TestBomb failed")
	}
}

func TestBomb7(t *testing.T) {
	var hai = []int{
		R3H, R3D, R3S,
	}

	var hai2 = []int{
		AH, AD,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)
	if msgCardHand.GetCardHandType() != int32(CardHandType_Pair) {
		t.Error("test TestBomb failed")
	}

	if pc.hasMsgCardHandGreatThan(msgCardHand) {
		t.Error("test TestBomb failed")
	}
}

func TestBomb8(t *testing.T) {
	var hai = []int{
		R3C, R3D, R3S,
	}

	var hai2 = []int{
		AH, AD,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)
	if msgCardHand.GetCardHandType() != int32(CardHandType_Pair) {
		t.Error("test TestBomb failed")
	}

	if !pc.hasMsgCardHandGreatThan(msgCardHand) {
		t.Error("test TestBomb failed")
	}
}

func TestBomb2(t *testing.T) {
	var hai = []int{
		JH, QH, KH,
	}

	var hai2 = []int{
		R4H, R4D, R4S, R4C,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)
	if msgCardHand.GetCardHandType() != int32(CardHandType_Bomb) {
		t.Error("test TestBomb2 failed")
	}

	if pc.hasMsgCardHandGreatThan(msgCardHand) {
		t.Error("test TestBomb2 failed")
	}
}
