package pddz

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

func TestSingle1(t *testing.T) {
	var hai = []int{
		R2H, JOR,
	}

	var hai2 = []int{
		JOB,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)

	if msgCardHand.GetCardHandType() != int32(CardHandType_Single) {
		t.Error("test single failed")
	}

	if !pc.hasSingleGreatThan(msgCardHand) {
		t.Error("test single failed")
	}
}

func TestPair1(t *testing.T) {
	var hai = []int{
		R2H, R2D,
	}

	var hai2 = []int{
		R3H, R3H,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)

	if msgCardHand.GetCardHandType() != int32(CardHandType_Pair) {
		t.Error("test pair failed")
	}

	if !pc.hasPairGreatThan(msgCardHand) {
		t.Error("test pair failed")
	}
}

func TestJOY1(t *testing.T) {
	// var hai = []int{
	// 	R2H, R2D,
	// }

	var hai2 = []int{
		JOB, JOR,
	}

	// pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)

	if msgCardHand.GetCardHandType() != int32(CardHandType_Roket) {
		t.Error("test pair failed")
	}

	// if !pc.hasPairGreatThan(msgCardHand) {
	// 	t.Error("test pair failed")
	// }
}

func TestTriplet1(t *testing.T) {
	var hai = []int{
		R2H, R2D, R2C,
	}

	var hai2 = []int{
		KH, KD, KC,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)

	if msgCardHand.GetCardHandType() != int32(CardHandType_Triplet) {
		t.Error("test triplet failed")
	}

	if !pc.hasTripletGreatThan(msgCardHand) {
		t.Error("test triplet failed")
	}
}

func TestBomb1(t *testing.T) {
	var hai = []int{
		AC, AD, AH, AS,
	}
	// var hai = []int{
	// 	R2H, R2D, R2C, R2S,
	// }
	var hai2 = []int{

		KH, KD, KC, KS,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)

	if msgCardHand.GetCardHandType() != int32(CardHandType_Bomb) {
		t.Error("test bomb failed")
	}

	if !pc.hasBombGreatThan(msgCardHand) {
		t.Error("test bomb failed")
	}
}

func TestTripletSingle1(t *testing.T) {
	var hai = []int{
		R2H, R2D, R2C, JOR,
	}

	var hai2 = []int{

		KH, KD, KC, QS,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)

	if msgCardHand.GetCardHandType() != int32(CardHandType_TripletSingle) {
		t.Error("test TripletSingle failed")
	}

	if !pc.hasTripletSingleGreatThan(msgCardHand) {
		t.Error("test TripletSingle failed")
	}
}

func TestTripletPair1(t *testing.T) {
	var hai = []int{
		R2H, R2D, R2C, R3H, R3C,
	}

	var hai2 = []int{

		KH, KD, KC, QS, QD,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)

	if msgCardHand.GetCardHandType() != int32(CardHandType_TripletPair) {
		t.Error("test TripletPair failed")
	}

	if !pc.hasTripletPairGreatThan(msgCardHand) {
		t.Error("test TripletPair failed")
	}
}

func TestFlush1(t *testing.T) {
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

	if !pc.hasFlushGreatThan(msgCardHand) {
		t.Error("test flush failed")
	}
}

func TestFlush2(t *testing.T) {
	var hai = []int{
		R8H,
		R9H,
		R10D,
		JH, KH,
		AH,
	}

	var hai2 = []int{
		R8H,
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

	if pc.hasFlushGreatThan(msgCardHand) {
		t.Error("test flush failed")
	}
}

func TestPair3X1(t *testing.T) {
	var hai = []int{
		R9H, R9D,
		R10H, R10D,
		JH, JD,
		QH, QD,
		KH, KD,
		AH, AD,
	}

	var hai2 = []int{
		R9H, R9D,
		R10H, R10D,
		JH, JD,
		QH, QD,
		KH, KD,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)
	if msgCardHand.GetCardHandType() != int32(CardHandType_Pair3X) {
		t.Error("test Pair3X failed")
	}

	if !pc.hasPair3XGreatThan(msgCardHand) {
		t.Error("test Pair3X failed")
	}
}

func TestTriplet2X1(t *testing.T) {
	var hai = []int{
		R9H, R9D, R9C,
		R10H, R10D, R10C,
		JH, JD, JC,
		QH, QD, QC,
		KH, KD, KC,
		AH, AD, AC,
	}

	var hai2 = []int{
		R10H, R10D, R10C,
		JH, JD, JC,
		QH, QD, QC,
		KH, KD, KC,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)
	if msgCardHand.GetCardHandType() != int32(CardHandType_Triplet2X) {
		t.Error("test Triplet2X failed")
		return
	}

	if !pc.hasTriplet2XGreatThan(msgCardHand) {
		t.Error("test Triplet2X failed")
	}
}

func TestTriplet2X2Single2(t *testing.T) {
	var hai = []int{
		R9H, R9D, R9C,
		R10H, R10D, R10C,
		JH, JD, JC,
		QH, QD, QC,
		KH, KD, KC,
		AH, AD, AC,
	}

	var hai2 = []int{
		R9H, R9D, R9C,
		JH, JD, JC,
		QH, QD, QC,
		KH, KD, KC,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)
	if msgCardHand.GetCardHandType() != int32(CardHandType_Triplet2X2Single) {
		t.Error("test Triplet2X failed")
		return
	}

	if !pc.hasTriplet2X2SingleGreatThan(msgCardHand) {
		t.Error("test Triplet2X failed")
	}
}

func TestTriplet2X2Single1(t *testing.T) {
	var hai = []int{
		JOB,
		R10H,
		JH,
		QH, QD, QC,
		KH, KD, KC,
		AH, AD, AC,
	}

	var hai2 = []int{
		R6H,
		R9H,
		R10H,
		JH, JD, JC,
		QH, QD, QC,
		KH, KD, KC,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)
	if msgCardHand.GetCardHandType() != int32(CardHandType_Triplet2X2Single) {
		t.Error("test Triplet2X2Single failed")
	}

	if !pc.hasTriplet2X2SingleGreatThan(msgCardHand) {
		t.Error("test Triplet2X2Single failed")
	}
}

func TestTriplet2X2Pair1(t *testing.T) {
	var hai = []int{
		R9H, R9D,
		R10H,
		JH, JD,
		QH, QD, QC,
		KH, KD, KC,
		AH, AD, AC,
	}

	var hai2 = []int{
		R7H, R7D,
		R9H,
		JH, JD, JC,
		QH, QD, QC,
		KH, KD, KC,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)
	if msgCardHand.GetCardHandType() != int32(CardHandType_Triplet2X2Single) {
		t.Error("test Triplet2X2Pair failed")
		return
	}

	if !pc.hasTriplet2X2SingleGreatThan(msgCardHand) {
		t.Error("test Triplet2X2Pair failed")
	}
}

func TestFourX2Pair1(t *testing.T) {
	var hai = []int{
		R2H, R2D,
		R9H, R9D,
		// JOB, JOR,
		AH, AD, AC, AS,
	}

	var hai2 = []int{

		R9H, R9D,
		R10H, R10D,

		KH, KD, KC, KS,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)
	if msgCardHand.GetCardHandType() != int32(CardHandType_FourX2Pair) {
		t.Error("test FourX2Pair failed")
	}

	if !pc.hasFourX2PairGreatThan(msgCardHand) {
		t.Error("test FourX2Pair failed")
	}
}

func TestFourX2Single1(t *testing.T) {
	var hai = []int{
		R2H, JOR,

		AH, AD, AC, AS,
	}

	var hai2 = []int{

		R9H,
		R10H,

		KH, KD, KC, KS,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)
	if msgCardHand.GetCardHandType() != int32(CardHandType_FourX2Single) {
		t.Error("test FourX2Single failed")
	}

	if !pc.hasFourX2SingleGreatThan(msgCardHand) {
		t.Error("test FourX2Single failed")
	}
}

func TestFourX2Single2(t *testing.T) {
	var hai = []int{
		KH, KD, KC, KS,
	}

	var hai2 = []int{

		R9H,
		R10H,

		KH, KD, KC, KS,
	}

	pc := testCreatePlayerCards(hai)
	msgCardHand := testCreateMsgCardHand(hai2)
	if msgCardHand.GetCardHandType() != int32(CardHandType_FourX2Single) {
		t.Error("test FourX2Single failed")
	}

	if !pc.hasMsgCardHandGreatThan(msgCardHand) {
		t.Error("test FourX2Single failed")
	}
}
