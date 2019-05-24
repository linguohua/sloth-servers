package club
import(
	"lobbyserver/lobby"
)


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

// GetClub 获取牌友圈
func (mgr *MyClubMgr)GetClub(clubID string) interface{} {
	club, ok := clubMgr.clubs[clubID]
	if !ok {
		return nil
	}
	return club
}

// IsUserPermisionCreateRoom 判断用户是否有权限创建房间
func (mgr *MyClubMgr) IsUserPermisionCreateRoom(userID string, clubID string) bool {
	mySQLUtil := lobby.MySQLUtil()
	role := mySQLUtil.LoadUserClubRole(userID, clubID)
	if role == int32(ClubRoleType_CRoleTypeCreator) || role == int32(ClubRoleType_CRoleTypeMgr) {
			return true
	}

	return false

}

// IsUserPermisionDeleteRoom 判断用户是否有权限创建房间
func (mgr *MyClubMgr) IsUserPermisionDeleteRoom(userID string, clubID string) bool {
	mySQLUtil := lobby.MySQLUtil()
	role := mySQLUtil.LoadUserClubRole(userID, clubID)
	if role == int32(ClubRoleType_CRoleTypeCreator) || role == int32(ClubRoleType_CRoleTypeMgr) {
			return true
	}

	return false

}

// IsClubMember 判断是否是牌友圈成员
func (mgr *MyClubMgr) IsClubMember(userID string, clubID string) bool {
	mySQLUtil := lobby.MySQLUtil()
	role := mySQLUtil.LoadUserClubRole(userID, clubID)
	if role == int32(ClubRoleType_CRoleTypeNone) {
		return false
	}

	return true

}

// func (mgr *MyClubMgr)addClub(club *Club) {
// 	_, ok := mgr.clubs[club.ID]
// 	if !ok {
// 		mgr.clubs[club.ID] = club
// 	}
// }