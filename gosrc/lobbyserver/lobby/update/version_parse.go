package update
import (
	"strings"
	"strconv"
)


// Version 版本号结构
type Version struct {
	bigVer int    // 大的版本号
	middleVer int // 中版本号
	smallVer int  // 小版本号
}

// 只对版本号类似1.0.1的版本号进行解析
func parseVersionString(versionString string) Version {
	vers := strings.Split(versionString, ".")
	if len(vers) != 3 {
		panic("invalid version string")
	}

	bigVer, err := strconv.Atoi(vers[0])
	if err != nil {
		panic("invalid big version")
	}

	middleVer, err := strconv.Atoi(vers[1])
	if err != nil {
		panic("invalid big version")
	}

	smallVer, err := strconv.Atoi(vers[2])
	if err != nil {
		panic("invalid big version")
	}


	ver := Version{}
	ver.bigVer = bigVer
	ver.middleVer = middleVer
	ver.smallVer = smallVer

	return ver
}

func parseVersionConfig(config string) {

}