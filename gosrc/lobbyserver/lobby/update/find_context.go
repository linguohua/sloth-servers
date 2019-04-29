package update

import (
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// findContext 用于查找可用更新模块时的上下文
type findContext struct {
	intDict    map[string]int
	stringDict map[string]string
}

// getInt 获得某个整数类型变量的值
func (ctx *findContext) getInt(name string) int {
	v, ok := ctx.intDict[name]
	if !ok {
		log.Panicln("findContext.getInt can't found name:", name)
	}
	return v
}

// getString 获得某个字符类型变量的值
func (ctx *findContext) getString(name string) string {
	v, ok := ctx.stringDict[name]
	if !ok {
		log.Panicln("findContext.getString can't found name:", name)
	}

	return v
}

// version2Int 把v1.0.0类似的版本号，转为一个int，方便比较大小
func version2Int(versionStr string) int {
	// 第一个字符是'v'，去除
	var versionStrDot = versionStr[1:]
	versionArray := strings.Split(versionStrDot, ".")
	var major, err1 = strconv.Atoi(versionArray[0])
	var minor, err2 = strconv.Atoi(versionArray[1])
	var patch, err3 = strconv.Atoi(versionArray[2])

	if err1 != nil {
		log.Panicln("version2Int, Atoi error:", err1)
	}

	if err2 != nil {
		log.Panicln("version2Int, Atoi error:", err2)
	}

	if err3 != nil {
		log.Panicln("version2Int, Atoi error:", err3)
	}

	return (major << 24) | (minor << 16) | (patch)
}

// parseFromHTTPReq 从request中提取参数构造一个find上下文
func parseFromHTTPReq(r *http.Request) *findContext {
	ctx := &findContext{}
	ctx.intDict = make(map[string]int)
	ctx.stringDict = make(map[string]string)

	var query = r.URL.Query()
	var v = query.Get("qMod")
	ctx.stringDict["qMod"] = v

	v = query.Get("modV")
	ctx.stringDict["modV"] = v
	vInt := version2Int(v)
	ctx.intDict["modV"] = vInt

	v = query.Get("csVer")
	ctx.stringDict["csVer"] = v
	vInt = version2Int(v)
	ctx.intDict["csVer"] = vInt

	v = query.Get("lobbyVer")
	ctx.stringDict["lobbyVer"] = v
	vInt = version2Int(v)
	ctx.intDict["lobbyVer"] = vInt

	v = query.Get("operatingSystem")
	ctx.stringDict["operatingSystem"] = v

	v = query.Get("operatingSystemFamily")
	ctx.stringDict["operatingSystemFamily"] = v

	v = query.Get("deviceUniqueIdentifier")
	ctx.stringDict["deviceUniqueIdentifier"] = v

	v = query.Get("deviceName")
	ctx.stringDict["deviceName"] = v

	v = query.Get("deviceModel")
	ctx.stringDict["deviceModel"] = v

	v = query.Get("network")
	ctx.stringDict["network"] = v

	return ctx
}
