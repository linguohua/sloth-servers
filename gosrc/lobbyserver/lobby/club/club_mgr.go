package club

// MyClubMgr 用户管理
type MyClubMgr struct {
	clubs map[string]*MsgClubInfo
}

func newClubMgr() *MyClubMgr {
	clubMgr := &MyClubMgr{}
	clubMgr.clubs = make(map[string]*MsgClubInfo)
	return clubMgr
}