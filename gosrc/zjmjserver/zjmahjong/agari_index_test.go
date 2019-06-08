package zjmahjong

import "testing"

var (
	mj = [][]int{
		{4, 4, 4, 4, 4, 4, 4, 4, 4},
		{4, 4, 4, 4, 4, 4, 4, 4, 4},
		{4, 4, 4, 4, 4, 4, 4, 4, 4},
		{4, 4, 4, 4, 4, 4, 4, 0, 0},
	}

	testCount = 0
)

// TestWinKot 测试胡牌刻子个数
// func TestWinKot(t *testing.T) {
// 	var hai = []int{
// 		MAN1, MAN2, MAN3,
// 		NAN, NAN, NAN,
// 		PIN1, PIN1, PIN1,
// 		PIN3, PIN4, PIN5,
// 		PIN4, PIN4}

// 	slots := make([]int, 42)
// 	for _, tileID := range hai {
// 		slots[tileID]++
// 	}

// 	if winAbleKotCount(slots) != 2 {
// 		t.Errorf("winable test failed")
// 	}

// 	hai = []int{
// 		MAN1, MAN2, MAN3,
// 		PIN4, PIN4}

// 	slots = make([]int, 42)
// 	for _, tileID := range hai {
// 		slots[tileID]++
// 	}

// 	if !isWinable(slots) || winAbleKotCount(slots, TILEMAX) != 0 {
// 		t.Errorf("winable test failed")
// 	}

// 	hai = []int{
// 		MAN, MAN, MAN,
// 		PIN4, PIN4}

// 	slots = make([]int, 42)
// 	for _, tileID := range hai {
// 		slots[tileID]++
// 	}

// 	if !isWinable(slots) || winAbleKotCount(slots) != 1 {
// 		t.Errorf("winable test failed")
// 	}

// 	hai = []int{
// 		PIN4, PIN4}

// 	slots = make([]int, 42)
// 	for _, tileID := range hai {
// 		slots[tileID]++
// 	}

// 	if !isWinable(slots) || winAbleKotCount(slots) != 0 {
// 		t.Errorf("winable test failed")
// 	}
// }

// TestWinable 测试胡牌判断
func TestWinable(t *testing.T) {
	var hai = []int{
		MAN1, MAN2, MAN3,
		NAN, NAN, NAN,
		PIN1, PIN1, PIN1,
		PIN3, PIN3, PIN3,
		PIN4, PIN4}

	slots := make([]int, 42)
	for _, tileID := range hai {
		slots[tileID]++
	}

	if !isWinable(slots) {
		t.Errorf("winable test failed")
	}
}
func TestWinableOfPair(t *testing.T) {
	var hai = []int{
		MAN6, MAN6}

	slots := make([]int, 42)
	for _, tileID := range hai {
		slots[tileID]++
	}

	if !isWinable(slots) {
		t.Errorf("TestWinableOfPair test failed")
	}
}
func TestWinable7Pair(t *testing.T) {
	var hai = []int{
		MAN1, MAN1,
		MAN2, MAN2,
		MAN3, MAN3,
		MAN4, MAN4,
		MAN5, MAN5,
		MAN6, MAN6,
		MAN7, MAN7}

	slots := make([]int, 42)
	for _, tileID := range hai {
		slots[tileID]++
	}

	if !isWinable(slots) {
		t.Errorf("TestWinable7Pair test failed")
	}
}

func isWinableTest(hai []int) bool {
	slots := make([]int, 42)
	for _, tileID := range hai {
		slots[tileID]++
	}

	return isWinable(slots)
}

// 测试7对子，没有重复对子
func TestPair7NoDuplicate(t *testing.T) {
	_arr := make([]int, 14)

	for a := 0; a < 34; a++ {
		_arr[0] = a
		_arr[1] = a

		for b := a + 1; b < 34; b++ {
			_arr[2] = b
			_arr[3] = b

			for c := b + 1; c < 34; c++ {
				_arr[4] = c
				_arr[5] = c

				for d := c + 1; d < 34; d++ {
					_arr[6] = d
					_arr[7] = d

					for e := d + 1; e < 34; e++ {
						_arr[8] = e
						_arr[9] = e
						for f := e + 1; f < 34; f++ {
							_arr[10] = f
							_arr[11] = f
							for g := f + 1; g < 34; g++ {
								_arr[12] = g
								_arr[13] = g

								if !isWinableTest(_arr) {
									t.Errorf("TestPair7NoDuplicate test failed")
								}
							}
						}
					}
				}
			}
		}
	}
}

// 测试7对子，有一对重复
func TestPair7With1Duplicate(t *testing.T) {
	_arr := make([]int, 14)

	for a := 0; a < 34; a++ {
		_arr[0] = a
		_arr[1] = a
		_arr[2] = a
		_arr[3] = a

		mj[a/9][a%9] -= 4

		for c := 0; c < 34; c++ {
			if mj[c/9][c%9] < 2 {
				continue
			}

			_arr[4] = c
			_arr[5] = c

			for d := c + 1; d < 34; d++ {
				if mj[d/9][d%9] < 2 {
					continue
				}

				_arr[6] = d
				_arr[7] = d

				for e := d + 1; e < 34; e++ {
					if mj[e/9][e%9] < 2 {
						continue
					}
					_arr[8] = e
					_arr[9] = e
					for f := e + 1; f < 34; f++ {
						if mj[f/9][f%9] < 2 {
							continue
						}
						_arr[10] = f
						_arr[11] = f
						for g := f + 1; g < 34; g++ {
							if mj[g/9][g%9] < 2 {
								continue
							}

							_arr[12] = g
							_arr[13] = g

							if !isWinableTest(_arr) {
								t.Errorf("TestPair7With1Duplicate test failed")
							}
						}
					}
				}
			}
		}

		mj[a/9][a%9] += 4
	}
}

// 测试7对子，有两对重复
func TestPair7With2Duplicate(t *testing.T) {
	_arr := make([]int, 14)
	for a := 0; a < 34; a++ {
		_arr[0] = a
		_arr[1] = a
		_arr[2] = a
		_arr[3] = a

		mj[a/9][a%9] -= 4

		for c := a + 1; c < 34; c++ {
			_arr[4] = c
			_arr[5] = c
			_arr[6] = c
			_arr[7] = c

			mj[c/9][c%9] -= 4

			for d := 0; d < 34; d++ {
				if mj[d/9][d%9] < 2 {
					continue
				}

				_arr[8] = d
				_arr[9] = d

				for e := d + 1; e < 34; e++ {
					if mj[e/9][e%9] < 2 {
						continue
					}
					_arr[10] = e
					_arr[11] = e
					for f := e + 1; f < 34; f++ {
						if mj[f/9][f%9] < 2 {
							continue
						}

						_arr[12] = f
						_arr[13] = f

						if !isWinableTest(_arr) {
							t.Errorf("TestPair7With2Duplicate test failed")
						}
					}
				}
			}

			mj[c/9][c%9] += 4
		}

		mj[a/9][a%9] += 4
	}
}

// 测试7对子，有三对重复
func TestPair7With3Duplicate(t *testing.T) {
	_arr := make([]int, 14)
	for a := 0; a < 34; a++ {
		_arr[0] = a
		_arr[1] = a
		_arr[2] = a
		_arr[3] = a

		mj[a/9][a%9] -= 4

		for c := a + 1; c < 34; c++ {
			_arr[4] = c
			_arr[5] = c
			_arr[6] = c
			_arr[7] = c

			mj[c/9][c%9] -= 4

			for d := c + 1; d < 34; d++ {
				_arr[8] = d
				_arr[9] = d
				_arr[10] = d
				_arr[11] = d

				mj[d/9][d%9] -= 4

				for e := 0; e < 34; e++ {
					if mj[e/9][e%9] < 2 {
						continue
					}

					_arr[12] = e
					_arr[13] = e

					if !isWinableTest(_arr) {
						t.Errorf("TestPair7With3Duplicate test failed")
					}
				}

				mj[d/9][d%9] += 4
			}

			mj[c/9][c%9] += 4
		}

		mj[a/9][a%9] += 4
	}
}

// 测试4个面子牌组+一对将
// 注意有大量重复，如果希望过滤重复则需要用一个哈希表来存储已枚举的组合
func Test4Melds(t *testing.T) {
	testCount = 0
	_arr := make([]int, 14)
	for j := 0; j < 4; j++ {
		for i := 0; i < 9; i++ {
			if mj[j][i] < 2 {
				// System.Diagnostics.Debug.Assert(j == 3 && (i == 7 || i == 8));
				continue
			}

			// remove head elements
			mj[j][i] = mj[j][i] - 2
			_arr[0] = 9*j + i
			_arr[1] = 9*j + i

			if !isWinableTest(_arr[:2]) {
				t.Errorf("Test4Melds test failed")
			}

			Four3x(4, _arr, t)

			//_set.Clear()
			// restore head elements
			mj[j][i] = mj[j][i] + 2
		}
	}

	t.Logf("Test4Melds test count:%d\n", testCount)
}

func Four3x(n int, _arr []int, t *testing.T) {

	if !isWinableTest(_arr[:(14 - n*3)]) {
		t.Errorf("Test4Melds test failed")
	}

	// 顺子
	for j := 0; j < 3; j++ {
		for i := 0; i < 7; i++ {
			if mj[j][i] > 0 && mj[j][i+1] > 0 && mj[j][i+2] > 0 {
				_arr[2+(4-n)*3] = 9*j + i
				_arr[2+(4-n)*3+1] = 9*j + i + 1
				_arr[2+(4-n)*3+2] = 9*j + i + 2

				if n-1 > 0 {
					// remove 3x elements
					mj[j][i]--
					mj[j][i+1]--
					mj[j][i+2]--

					Four3x(n-1, _arr, t)
					// restore 3x elements
					mj[j][i]++
					mj[j][i+1]++
					mj[j][i+2]++
				} else {
					testCount++

					//Array.Sort(_arr);
					//System.Diagnostics.Debug.Assert(jap.CanHu(_arr));
					//SaveCalcKey(_arr);
					//_set.Add(new SetItem(_arr));
					//Keep6(_arr);
					if !isWinableTest(_arr) {
						t.Errorf("Test4Melds test failed")
					}
				}
			}
		}
	}

	// 刻子
	for j := 0; j < 4; j++ {
		for i := 0; i < 9; i++ {
			if mj[j][i] > 2 {
				_arr[2+(4-n)*3] = 9*j + i
				_arr[2+(4-n)*3+1] = 9*j + i
				_arr[2+(4-n)*3+2] = 9*j + i
				if n-1 > 0 {
					// remove 3x elements
					mj[j][i] -= 3

					Four3x(n-1, _arr, t)

					// restore 3x elements
					mj[j][i] += 3
				} else {
					testCount++
					//Array.Sort(_arr);
					//System.Diagnostics.Debug.Assert(jap.CanHu(_arr));
					//SaveCalcKey(_arr);
					//_set.Add(new SetItem(_arr));
					//Keep6(_arr);
					if !isWinableTest(_arr) {
						t.Errorf("Test4Melds test failed")
					}
				}

			}
		}
	}
}
