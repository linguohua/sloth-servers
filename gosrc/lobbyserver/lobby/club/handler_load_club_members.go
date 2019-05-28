package club

import (
	"net/http"
	"strconv"
	"gconst"
	"lobbyserver/lobby"
	"github.com/garyburd/redigo/redis"
	proto "github.com/golang/protobuf/proto"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// onLoadClubMembers 加载俱乐部成员列表
func onLoadClubMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("userID")

	log.Println("onLoadClubMembers, userID:", userID)

	var query = r.URL.Query()
	// 俱乐部ID
	clubID := query.Get("clubID")
	// cursorStr := query.Get("cursor")

	if clubID == "" {
		log.Println("onLoadClubMembers, clubID is empty, userID:", userID)
		sendGenericError(w, ClubOperError_CERR_Invalid_Input_Parameter)
		return
	}

	// cursor := 0
	// if cursorStr != "" {
	// 	cursor, _ = strconv.Atoi(cursorStr)
	// }

	club, ok := clubMgr.clubs[clubID]
	if !ok {
		log.Println("onLoadClubMembers, clubID is invalid:", clubID)
		sendGenericError(w, ClubOperError_CERR_Invalid_Input_Parameter)
		return
	}

	mySQLUtil := lobby.MySQLUtil()
	role := mySQLUtil.LoadUserClubRole(userID, clubID)
	if role == int32(ClubRoleType_CRoleTypeNone) {
		sendGenericError(w, ClubOperError_CERR_User_Not_In_Club)
		return
	}

	memberIDs := mySQLUtil.LoadClubUserIDs(clubID)

	// 获得redis连接
	conn := lobby.Pool().Get()
	defer conn.Close()

	// // 检查是否可以加载
	// if !isUserClubMember(clubID, userID, w) {
	// 	return
	// }

	// 得到的member IDs
	// memberIDs, newCusor := club.membersSscan(cursor)

	loadMemberReply := &MsgClubLoadMembersReply{}
	newCusor32 := int32(0)
	loadMemberReply.Cursor = &newCusor32

	// 填充member 列表
	loadMemberReply.Members = constructClubMemberList(memberIDs, conn, club)

	b, err := proto.Marshal(loadMemberReply)
	if err != nil {
		log.Println("onLoadClubMembers, marshal error:", err)
		sendGenericError(w, ClubOperError_CERR_Encode_Decode)
		return
	}

	sendMsgClubReply(w, ClubReplyCode_RCOperation, b)
}

func constructClubMemberList(memberIDs []string, conn redis.Conn, club *Club) []*MsgClubMemberInfo {

	clubMembers := make([]*MsgClubMemberInfo, 0, len(memberIDs))

	conn.Send("MULTI")
	for _, mID := range memberIDs {
		conn.Send("HMGET", gconst.LobbyUserTablePrefix+mID, "Nick", "Sex", "Protrait", "AvatarID")
	}

	values, err := redis.Values(conn.Do("EXEC"))

	if err != nil {
		log.Println("constructClubMemberList, redis err:", err)
		return clubMembers
	}

	for i, val := range values {
		strValues, err := redis.Strings(val, nil)
		if err != nil {
			log.Println("constructClubMemberList, redis err:", err)
			continue
		}

		displayInfo := &MsgClubDisplayInfo{}
		nick := strValues[0]
		displayInfo.Nick = &nick
		sex, _ := strconv.Atoi(strValues[1])
		sex32 := uint32(sex)
		displayInfo.Sex = &sex32
		headIconURL := strValues[2]
		displayInfo.HeadIconURL = &headIconURL
		avatarID, _ := strconv.Atoi(strValues[3])
		avatarID32 := int32(avatarID)
		displayInfo.AvatarID = &avatarID32
		userID := memberIDs[i]

		memberInfo := &MsgClubMemberInfo{}
		memberInfo.UserID = &userID
		memberInfo.DisplayInfo = displayInfo

		// online := club.isMemberOnline(userID)
		// memberInfo.Online = &online

		clubMembers = append(clubMembers, memberInfo)
	}

	return clubMembers
}
