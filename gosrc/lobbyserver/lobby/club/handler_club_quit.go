package club

import (
	"net/http"
	"gconst"
	"time"
	"lobbyserver/lobby"
	"github.com/julienschmidt/httprouter"
	"github.com/garyburd/redigo/redis"
	proto "github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

// onQuit 主动退出俱乐部
func onQuit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("userID")
	log.Println("onQuit, userID:", userID)

	var query = r.URL.Query()
	// 俱乐部ID
	clubID := query.Get("clubID")

	if clubID == "" {
		log.Println("onQuit, need clubID")
		sendGenericError(w, ClubOperError_CERR_Invalid_Input_Parameter)
		return
	}

	// 获得redis连接
	conn := lobby.Pool().Get()
	defer conn.Close()

	// 检查合法性
	ok := isQuitAble(clubID, userID, w)
	if !ok {
		return
	}

	mySQLUtil := lobby.MySQLUtil()
	mySQLUtil.RemoveUserFromClub(userID, clubID)

	redisClearClubUserData(clubID, userID, conn)

	//newUserLeaveEvent(clubID, userID, conn)  玩家离开俱乐部 现在不发消息

	// 重新加载所有俱乐部到客户端
	onLoadMyClubs(w, r, ps)
}

func isQuitAble(clubID string, userID string, w http.ResponseWriter) bool {
	// 判斷牌友群是否存在
	club, ok := clubMgr.clubs[clubID]
	if !ok {
		log.Printf("isQuitAble, club %s not exist", clubID)
		sendGenericError(w, ClubOperError_CERR_Club_Not_Exist)
		return false
	}

	// 判断用户是否在牌友圈中
	mySQLUtil := lobby.MySQLUtil()
	role := mySQLUtil.LoadUserClubRole(userID, clubID)
	if role == int32(ClubRoleType_CRoleTypeNone) {
		log.Printf("isQuitAble, user %s not in club %s", userID, clubID)
		sendGenericError(w, ClubOperError_CERR_User_Not_In_Club)
		return false
	}

	// 用户不能是群主
	clubInfo := club.clubInfo
	creatorUserID := clubInfo.GetCreatorUserID()
	if userID == creatorUserID {
		log.Printf("isQuitAble, user %s is owner %s, can't quit\n", creatorUserID, userID)
		sendGenericError(w, ClubOperError_CERR_Owner_Can_not_quit)
		return false
	}

	return true
}

func redisClearClubUserData(clubID string, userID string, conn redis.Conn) {
	// 移除用户的club数据
	conn.Send("MULTI")
	// 移除club的用户数据
	conn.Send("SREM", gconst.LobbyClubMemberSetPrefix+clubID, userID)

	// 清理事件相关数据
	conn.Send("DEL", gconst.LobbyClubUnReadEventUserListPrefix+clubID+":"+userID)
	conn.Send("DEL", gconst.LobbyClubUnReadEventUserSetPrefix+clubID+":"+userID)

	conn.Do("EXEC")
}

func newUserLeaveEvent(clubID string, userID string, conn redis.Conn) {
	// 事件ID
	cn, err := redis.Int64(conn.Do("HINCRBY", gconst.LobbyClubSysTable, "clubEventID", 1))
	if err != nil {
		log.Panicln("newJoinEvent alloc eventID failed, redis err:", err)
	}

	clubEvent := &MsgClubEvent{}
	evtType32 := int32(ClubEventType_CEVT_Quit)
	clubEvent.EvtType = &evtType32
	to := ""
	clubEvent.To = &to
	generatedTime32 := uint32(time.Since(time2010).Seconds())
	clubEvent.GeneratedTime = &generatedTime32
	needHandle := false // 通知事件是需要处理的
	clubEvent.NeedHandle = &needHandle

	eventID32 := uint32(cn % int64(0x0ffffffff))
	clubEvent.Id = &eventID32

	clubEvent.UserID1 = &userID

	clubEventBytes, err := proto.Marshal(clubEvent)
	if err != nil {
		log.Panic(err)
	}

	// TODO: 后面需要增加裁剪如下各个列表的定时器

	conn.Send("MULTI")
	// 加入到俱乐部的信息列表
	conn.Send("HSET", gconst.LobbyClubEventTablePrefix+clubID, eventID32, clubEventBytes)
	conn.Send("LPUSH", gconst.LobbyClubEventListPrefix+clubID, eventID32)
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
