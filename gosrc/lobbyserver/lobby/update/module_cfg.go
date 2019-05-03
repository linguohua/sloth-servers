package update

import (
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// AssetsBundleCfg 模块内的assets bundle配置
type AssetsBundleCfg struct {
	Name string   `json:"name"` // bundle 名字
	MD5  string   `json:"md5"`  // bundle md5，作为确定是否修改的标志
	Size int64    `json:"size"` // bundle的大小，单位是字节
	Deps []string `json:"deps"` // bundle的依赖，主要是依赖其他bundle
}

// ModuleCfg 模块配置
type ModuleCfg struct {
	Name           string `json:"name"`    // 模块名字
	Version        string `json:"version"` // 模块版本号
	versionInteger int    // 模块的版本号转换为整数方便比较

	CSVersion        string `json:"csVer"` // 模块依赖的CSHARP库版本号
	csVersionInteger int    // 模块的版本号转换为整数方便比较

	LobbyVersion        string `json:"lobbyVer"` // 模块依赖的大厅库版本号
	lobbyVersionInteger int    // 模块的版本号转换为整数方便比较

	IsDefault bool              `json:"default"` // 是否是默认模块配置，每一个模块只能有一个默认模块配置
	AbList    []AssetsBundleCfg `json:"abList"`  // 模块的assets bundle列表

	AndConditions []*ConditionCfg `json:"and"` // AND 条件集合
	OrConditions  []*ConditionCfg `json:"or"`  // OR 条件集合
}

// ConditionCfg 条件
type ConditionCfg struct {
	CondVariable string `json:"variable"` // 条件变量
	Reverse      bool   `json:"reverse"`  // 是否取反
	Operator     string `json:"operator"` // 操作符，大于，小于，等于，含有
	Value        string `json:"value"`    // 值
	valueInt     int    // Value 转换为整数
}

// value2Int 把条件中的Value字符串转换为int，方便后续比较
func (cfg *ConditionCfg) value2Int() {
	vc, ok := conditionVarCfg.VariableCfgMap[cfg.CondVariable]
	if !ok {
		log.Panicln("conditionVariableCfgIsInt, unknown variable:", cfg.CondVariable)
	}

	if vc.VType == "version" {
		cfg.valueInt = version2Int(cfg.Value)
	} else if vc.VType == "int" {
		var valueInt, err = strconv.Atoi(cfg.Value)
		if err != nil {
			log.Error("ConditionCfg.value2Int Atoi error:", err)
		} else {
			cfg.valueInt = valueInt
		}
	}
}

// verify 检查单个条件
func (cfg *ConditionCfg) verify(ctx *findContext) bool {
	var result bool
	if conditionVariableCfgIsInt(cfg.CondVariable) {
		ctxInt := ctx.getInt(cfg.CondVariable)
		// 整数类型
		switch cfg.Operator {
		case "gt":
			result = cfg.valueInt > ctxInt
		case "lt":
			result = cfg.valueInt < ctxInt
		case "eq":
			result = cfg.valueInt == ctxInt
		default:
			log.Panicln("ConditionCfg.verify unknown integer operator type:", cfg.Operator)
		}
	} else {
		// 字符串类型
		ctxStr := ctx.getString(cfg.CondVariable)
		switch cfg.Operator {
		case "ct":
			result = strings.Contains(ctxStr, cfg.Value) || strings.Contains(cfg.Value, ctxStr)
		case "eq":
			result = ctxStr == cfg.Value
		default:
			log.Panicln("ConditionCfg.verify unknown string operator type:", cfg.Operator)
		}
	}

	if cfg.Reverse {
		result = !result
	}

	return result
}

// verifyConditions 条件测试
// 相当于： return ((and1 && and2 && and...) && (or1 || or2 || or...))
func (cfg *ModuleCfg) verifyConditions(ctx *findContext) bool {
	// 先检查AND关系的conditions列表
	for _, cond := range cfg.AndConditions {
		// 如果有一个AND条件测试失败，则失败
		if !cond.verify(ctx) {
			return false
		}
	}

	// 没有OR条件，因此默认成功
	if len(cfg.OrConditions) < 1 {
		return true
	}

	// 再检查OR关系的conditions列表
	for _, cond := range cfg.OrConditions {
		// 如果有一个OR条件测试成功，则成功
		if cond.verify(ctx) {
			return true
		}
	}

	// 所有的OR条件测试都失败
	return false
}

// strings2Int 把一些字符串转换为int，方便后续的比较
func (cfg *ModuleCfg) strings2Int() {
	// 版本号
	cfg.versionInteger = version2Int(cfg.Version)
	// csharp库版本号
	cfg.csVersionInteger = version2Int(cfg.CSVersion)

	// 如果不是大厅模块，则还需要转换其依赖的大厅模块版本号
	if cfg.Name != "lobby" {
		cfg.lobbyVersionInteger = version2Int(cfg.LobbyVersion)
	}

	// 把条件中的value转换为int，如果可以转换的话
	for _, c := range cfg.AndConditions {
		c.value2Int()
	}

	// 把条件中的value转换为int，如果可以转换的话
	for _, c := range cfg.OrConditions {
		c.value2Int()
	}
}
