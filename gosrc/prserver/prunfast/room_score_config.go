package prunfast

// RoomScoreConfig 分数配置
type RoomScoreConfig struct {
	miniWinBasicScore int
}

func newRoomScoreConfig() *RoomScoreConfig {
	rsc := &RoomScoreConfig{}
	rsc.miniWinBasicScore = 1

	return rsc
}
