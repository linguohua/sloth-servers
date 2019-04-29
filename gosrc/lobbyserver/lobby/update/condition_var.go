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
			"version":["eq","gt","lt"],
			"string":["ct","eq"],
			"int":["eq","gt","lt"]
		},
		"variables":{
			"qMode": {
				"type":"string"
			},
			"modV": {
				"type":"version"
			},
			"csVer": {
				"type":"version"
			},
			"lobbyVer": {
				"type":"version"
			},
			"operatingSystem": {
				"type":"version"
			},
			"operatingSystemFamily": {
				"type":"version"
			},
			"deviceUniqueIdentifier": {
				"type":"version"
			},
			"deviceName": {
				"type":"version"
			},
			"deviceModel": {
				"type":"version"
			},
			"network": {
				"type":"version"
			}
		}
	}
	`
)

var (
	// 条件变量配置，从conditonVarJSON marshal
	conditionVarCfg *ConditionVariableCfg
)

// VariableCfg 条件变量配置
type VariableCfg struct {
	VType string `json:"type"`
}

// ConditionVariableCfg 条件变量配置
type ConditionVariableCfg struct {
	OperatorCfgMap map[string][]string    `json:"operators"`
	VariableCfgMap map[string]VariableCfg `json:"variables"`
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
