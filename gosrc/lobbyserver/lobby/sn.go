package lobby

var (
	//Sn 发消息的系列号
	sn = uint32(0)
)

func generateSn() uint32 {
	sn++
	maxInt32 := uint32(1<<32 - 1)
	if sn >= maxInt32 {
		sn = 1
	}

	return (sn)
}

// GenerateSn 生成一个系列号
func GenerateSn() uint32 {
	return generateSn()
}
