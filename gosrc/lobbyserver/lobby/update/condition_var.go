package update

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

const (
	// 这部分作为JSON是由于网页上也要用到，用于上传配置时，给用户选择
	// 配置的有效条件，例如用于灰度发布时指定满足条件才可以更新
	conditonVarJSON = `
	{
		"operators":{
			"version":[{"n":"等于","k":"eq"},{"n":"大于","k":"gt"},{"n":"小于","k":"lt"}],
			"string":[{"n":"等于","k":"eq"},{"n":"包含","k":"ct"}],
			"int":[{"n":"等于","k":"eq"},{"n":"大于","k":"gt"},{"n":"小于","k":"lt"}]
		},
		"variables":{
			"qMode": {
                "name":"模块名字",
				"type":"string"
			},
			"modV": {
                "name":"模块版本",
				"type":"version"
			},
			"csVer": {
                "name":"csharp版本",
				"type":"version"
			},
			"lobbyVer": {
                "name":"lobby版本",
				"type":"version"
			},
			"operatingSystem": {
                "name":"操作系统",
				"type":"string"
			},
			"operatingSystemFamily": {
                "name":"操作系统集",
				"type":"string"
			},
			"deviceUniqueIdentifier": {
                "name":"设备唯一ID",
				"type":"string"
			},
			"deviceName": {
                "name":"设备名字",
				"type":"string"
			},
			"deviceModel": {
                "name":"设备类",
				"type":"string"
			},
			"network": {
                "name":"网络类型",
				"type":"string"
			}
		}
	}`
)

var (
	// 条件变量配置，从conditonVarJSON marshal
	conditionVarCfg *ConditionVariableCfg
)

// OperatorCfg 条件变量配置
type OperatorCfg struct {
	Name     string `json:"n"`
	Operator string `json:"k"`
}

// VariableCfg 条件变量配置
type VariableCfg struct {
	VType string `json:"type"`
	Name  string `json:"name"`
}

// ConditionVariableCfg 条件变量配置
type ConditionVariableCfg struct {
	OperatorCfgMap map[string][]OperatorCfg `json:"operators"`
	VariableCfgMap map[string]VariableCfg   `json:"variables"`
}

// conditionVariableCfgIsInt 判断一个变量是否是整数类型的，也即是可以比较大小的
func conditionVariableCfgIsInt(name string) bool {
	vc, ok := conditionVarCfg.VariableCfgMap[name]
	if !ok {
		log.Panicln("conditionVariableCfgIsInt, unknown variable:", name)
	}

	if vc.VType == "version" || vc.VType == "int" {
		return true
	}

	return false
}

// initConditionVariableCfg 初始化条件变量配置
func initConditionVariableCfg() {
	log.Trace("update initConditionVariableCfg")
	c := &ConditionVariableCfg{}
	err := json.Unmarshal([]byte(conditonVarJSON), c)
	if err != nil {
		log.Panicln("initConditionVariableCfg Error: ", err)
	}

	conditionVarCfg = c
}
