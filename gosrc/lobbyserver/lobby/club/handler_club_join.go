package club

import (
	"net/http"
	"gconst"
	"time"
	"fmt"
	"lobbyserver/lobby"
	"github.com/golang/protobuf/proto"
	"github.com/julienschmidt/httprouter"
	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
)

// onJoinClub 申请加入某个俱乐部
func onJoinClub(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("userID")

	log.Println("onJoinClub, userID:", userID)

	var query = r.URL.Query()
	clubNumber := query.Get("clubNumber")

	if clubNumber == "" {
		log.Println("onJoinClub, clubNumber is empty")
		sendGenericError(w, ClubOperError_CERR_Invalid_Input_Parameter)
		return
	}

	mySQLUtil := lobby.MySQLUtil()
	clubID := mySQLUtil.LoadClubIDByNumber(clubNumber)
	if clubID == "" {
		log.Println("onJoinClub, club not exist for clubNumber:", clubNumber)
		sendGenericError(w, ClubOperError_CERR_Club_Not_Exist)
		return
	}

	club, ok := clubMgr.clubs[clubID]
	if !ok {
		log.Println("onJoinClub, club not exist for clubID:", clubID)
		sendGenericError(w, ClubOperError_CERR_Club_Not_Exist)
		return
	}

	role := mySQLUtil.LoadUserClubRole(userID, clubID)
	if role != int32(ClubRoleType_CRoleTypeNone) {
		sendGenericError(w, ClubOperError_CERR_Invitee_Already_In_Club)
		return
	}


	isApplicant := isApplicant(clubID, userID)

	// 之前已经申请过
	if isApplicant {
		log.Printf("onJoinClub, user %s in club %s block list\n", userID, clubID)
		sendGenericError(w, ClubOperError_CERR_You_Already_Applicate)
		return
	}

	clubInfo := club.clubInfo
	// 检查玩家已经加入过的俱乐部个数
	maxJoin := mySQLUtil.CountUserClubNumber(userID)
	if int32(maxJoin) >= clubInfo.GetMaxMember() {
		log.Println("onJoinClub, user has join exceed max limit:", maxJoin)
		sendGenericError(w, ClubOperError_CERR_Exceed_Max_Club_Count_Limit)
		return
	}

	// 俱乐部禁止申请加入
	if clubInfo.GetJoinForbit() {
		log.Println("onJoinClub, club forbit to join")
		sendGenericError(w, ClubOperError_CERR_Club_Forbit_Join)
		// sendGenericErrorWithExtraString(w, ClubOperError_CERR_Club_Forbit_Join, jReason)
		return
	}

	// club, ok := clubsMap[clubID]
	// if ok {
	// 	club.updateUpdateTime()
	// }

	// 生成事件
	newApplicateEvent(clubID, userID, clubInfo.GetCreatorUserID())

	// 操作成功
	sendGenericError(w, ClubOperError_CERR_OK)
}

func isApplicant(clubID string, applicantUserID string) bool {
	conn := lobby.Pool().Get()
	defer conn.Close()

	isExist, err := redis.Int(conn.Do("SISMEMBER", gconst.LobbyClubApplicantPrefix+clubID, applicantUserID))
	if err != nil {
		log.Println("check isApplicant error:", err)
	}

	if isExist > 0 {
		return true
	}

	return false
}

func newApplicateEvent(clubID string, applicantUserID string, owner string) {
	// 申请事件ID
	conn := lobby.Pool().Get()
	defer conn.Close()

	cn, err := redis.Int64(conn.Do("HINCRBY", gconst.LobbyClubSysTable, "clubEventID", 1))
	if err != nil {
		log.Panicln("newApplicateEvent alloc eventID failed, redis err:", err)
	}

	clubEvent := &MsgClubEvent{}
	evtType32 := int32(ClubEventType_CEVT_NewApplicant)
	clubEvent.EvtType = &evtType32
	clubEvent.To = &owner
	generatedTime32 := uint32(time.Since(time2010).Seconds())
	clubEvent.GeneratedTime = &generatedTime32
	needHandle := true // 申请事件是需要处理的
	clubEvent.NeedHandle = &needHandle

	eventID32 := uint32(cn % int64(0x0ffffffff))
	clubEvent.Id = &eventID32
	clubEvent.UserID1 = &applicantUserID

	clubEventBytes, err := proto.Marshal(clubEvent)
	if err != nil {
		log.Panic(err)
	}

	// TODO: 后面需要增加裁剪如下各个列表的定时器

	conn.Send("MULTI")
	// 加入到俱乐部的信息列表
	conn.Send("HSET", gconst.LobbyClubEventTablePrefix+clubID, eventID32, clubEventBytes)
	conn.Send("HSET", gconst.LobbyClubNeedHandledTablePrefix+clubID, eventID32, owner)
	conn.Send("LPUSH", gconst.LobbyClubEventListPrefix+clubID, eventID32)
	conn.Send("SADD", gconst.LobbyClubApplicantPrefix+clubID, applicantUserID)
	conn.Send("LPUSH", gconst.LobbyClubUserApplicantEventPrefix + applicantUserID, fmt.Sprintf("%s,%d", clubID, eventID32))

	_, err = conn.Do("EXEC")
	if err != nil {
		log.Panic(err)
	}

	// 加入到成员的未读信息列表
	// KEYS[1] prefix
	// KEYS[2] member-set key
	// KEYS[3] eventID
	_, err = luaScriptInsertNewEvent.Do(conn, gconst.LobbyClubUnReadEventUserListPrefix+clubID+":",
		gconst.LobbyClubUnReadEventUserSetPrefix+clubID+":", gconst.LobbyClubMemberSetPrefix+clubID, eventID32)

	if err != nil {
		log.Panic(err)
	}
}
