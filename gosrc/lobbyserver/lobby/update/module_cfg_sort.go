package update

// ByVersion 根据version排序
type ByVersion []*ModuleCfg

func (a ByVersion) Len() int      { return len(a) }
func (a ByVersion) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// 逆序，自高到底
func (a ByVersion) Less(i, j int) bool { return a[i].versionInteger > a[j].versionInteger }
