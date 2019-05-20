package club


// Club 牌友群
type Club struct {
	ID                    string
	clubInfo              *MsgClubInfo      // club info
	mm                    map[string]bool   // 成员列表
}

// MyClubMgr 用户管理
type MyClubMgr struct {
	clubs map[string]*Club
}

func newBaseClub(clubInfo *MsgClubInfo, clubID string) *Club {
	club := &Club{}
	club.clubInfo = clubInfo
	club.ID = clubID
	club.mm = make(map[string]bool)

	return club
}

func newClubMgr() *MyClubMgr {
	clubMgr := &MyClubMgr{}
	clubMgr.clubs = make(map[string]*Club)
	return clubMgr
}

// func (mgr *MyClubMgr)addClub(club *Club) {
// 	_, ok := mgr.clubs[club.ID]
// 	if !ok {
// 		mgr.clubs[club.ID] = club
// 	}
// }