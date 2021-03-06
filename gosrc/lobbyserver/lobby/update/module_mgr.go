// 更新模块的设计思路：
// 配置的核心是ModuleCfg，这个对象是由两部分内容生成的，第一是用户通过web，提交unity3d打包生成
// 模块配置文件cfg.json，第二是用户提交cfg.json时，在web页面上填写的其他额外参数，因此，ModuleCfg
// 的主要内容是来自unity3d打包生成的cfg.json，有用户通过web提交到服务器并写入redis。

// 一个ModuleCfg可以被客户端使用的前提是，ModuleCfg的版本比客户端的新，而且其依赖的CSHARP库以及lobby模块
// 必须不能高于客户端目前的，否则客户端更新完成后，由于CSHART和lobby模块不匹配导致跑不起来。
// ModuleCfg的硬性要求匹配之后，就再检查客户端是否满足ModuleCfg所配置的条件，这些条件目前主要是用于灰度发布。

// 一个模块可以配置一个称为“默认”的ModuleCfg，所谓默认的意思是，如果客户端的版本低于该ModuleCfg版本，则强制客户端
// 做更新。如果模块有多个配置为“默认”的ModuleCfg，则取版本最高者作为唯一“默认”。

package update

import (
	"sort"

	log "github.com/sirupsen/logrus"
)

var (
	mmgr *ModulesMgr
)

// ModulesMgr 模块管理
type ModulesMgr struct {
	moduels map[string]*Module
}

// Module 模块
type Module struct {
	cfgs       []*ModuleCfg // 所有更新配置，安装版本号排序
	defaultCfg *ModuleCfg   // 默认模块配置
}

// findModuleCfg 寻找合适的更新配置
func (mm *ModulesMgr) findModuleCfg(ctx *findContext) *ModuleCfg {
	qMod := ctx.getString("qMod")
	modVInt := ctx.getInt("modV")
	csVInt := ctx.getInt("csVer")

	m, ok := mm.moduels[qMod]
	if !ok {
		log.Debug("findModuleCfg, no module cfg found for:", qMod)
		return nil
	}

	for _, mcfg := range m.cfgs {
		// 没有找到版本号更加新的配置
		if mcfg.versionInteger <= modVInt {
			// 由于版本配置是按顺序排列的，因此后面的不需要再比较
			return nil
		}

		// 比较csharp版本
		if mcfg.csVersionInteger > csVInt {
			// 更新配置要求的CSHARP模块版本较高，因此不适用
			continue
		}

		// 比较大厅版本
		if qMod != "lobby" {
			// 只要非大厅版本，才检查依赖的大厅版本
			lobbyInt := ctx.getInt("lobbyVer")
			if mcfg.lobbyVersionInteger > lobbyInt {
				// 更新配置要求的大厅模块版本较高，因此不适用
				continue
			}
		}

		// 找到了版本号合适的配置，检查条件是否满足
		if mcfg.verifyConditions(ctx) {
			return mcfg
		}
	}

	return nil
}

// getDefaultCfg 获取默认更新，如果配置了默认更新，且适用则返回默认更新
func (mm *ModulesMgr) getDefaultCfg(ctx *findContext) (*ModuleCfg, int) {
	qMod := ctx.getString("qMod")
	modVInt := ctx.getInt("modV")
	csVInt := ctx.getInt("csVer")

	m, ok := mm.moduels[qMod]
	if !ok {
		log.Debug("getDefaultCfg, no module cfg found for:", qMod)
		return nil, 0
	}

	defaultCfg := m.defaultCfg
	if defaultCfg == nil {
		// 没有配置默认的更新
		return nil, 0
	}

	if defaultCfg.versionInteger <= modVInt {
		// 客户端的版本更加新，不需要更新
		return nil, 0
	}

	// 比较csharp版本
	if defaultCfg.csVersionInteger > csVInt {
		// 更新配置要求的CSHARP模块版本较高，因此不适用
		return nil, queryErrorNeedUpgradeCS
	}

	// 比较大厅版本
	if qMod != "lobby" {
		// 只要非大厅版本，才检查依赖的大厅版本
		lobbyInt := ctx.getInt("lobbyVer")
		if defaultCfg.lobbyVersionInteger > lobbyInt {
			// 更新配置要求的大厅模块版本较高，因此不适用
			return nil, queryErrorNeedUpgradeLobby
		}
	}

	return defaultCfg, 0
}

// addNewModuleCfg 增加新的更新配置
func (mm *ModulesMgr) addNewModuleCfg(mc *ModuleCfg) {
	name := mc.Name

	m, ok := mmgr.moduels[name]
	if !ok {
		m = &Module{}
		m.cfgs = make([]*ModuleCfg, 0, 16)
		mmgr.moduels[name] = m
	}

	// 转换字符串为int
	mc.strings2Int()

	var isUpdate = false
	for i, mc2 := range m.cfgs {
		if mc2.versionInteger == mc.versionInteger {
			m.cfgs[i] = mc
			isUpdate = true
			break
		}
	}

	if !isUpdate {
		m.cfgs = append(m.cfgs, mc)
		m.refresh()
	}

	saveModuleCfg2Redis(mc)
}

// deleteModuleCfgs 删除一组更新配置
func (mm *ModulesMgr) deleteModuleCfgs(mcs []*ModuleCfg) {
	moduleMap := make(map[string]*Module)
	for _, mc := range mcs {
		name := mc.Name
		// 转换字符串为int
		mc.strings2Int()

		m, ok := mmgr.moduels[name]
		if !ok {
			continue
		}

		m.removeModuleCfg(mc)
		moduleMap[name] = m
	}

	// 刷新数据
	for _, m := range moduleMap {
		m.resetDefault()
	}

	deleteModuleCfgFromRedis(mcs)
}

// refresh 刷新数据，主要是排序以及设置defalut
func (m *Module) refresh() {
	// 排序，按照版本号由高到底排序
	sort.Sort(ByVersion(m.cfgs))
	m.resetDefault()
}

func (m *Module) resetDefault() {
	m.defaultCfg = nil
	// 检查有没有默认配置
	for _, mc := range m.cfgs {
		if mc.IsDefault {
			m.defaultCfg = mc
			break
		}
	}
}

// removeModuleCfg 移除模块
func (m *Module) removeModuleCfg(mc *ModuleCfg) {
	for i, mc2 := range m.cfgs {
		if mc2.versionInteger == mc.versionInteger {
			m.cfgs = append(m.cfgs[:i], m.cfgs[i+1:]...)
			break
		}
	}
}
