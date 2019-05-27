package club

import (
	"net/http"
	"gconst"
	"lobbyserver/lobby"
	"github.com/julienschmidt/httprouter"
	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
)

func forceDisabandClubRooms(clubID string) {
	conn := lobby.Pool().Get()
	defer conn.Close()

	roomIDs, err := redis.Strings(conn.Do("SMEMBERS", gconst.LobbyClubRoomSetPrefix+clubID))
	if err != nil && err != redis.ErrNil {
		log.Println("forceDisabandClubRooms, SMEMBERS redis error:", err)
		return
	}

	roomUtil := lobby.RoomUtil()

	for _, roomID := range roomIDs {
		errCode := roomUtil.ForceDeleteRoom(roomID)
		if errCode != 0 {
			log.Errorf("ForceDeleteRoom %s failed, code:%d", roomID, errCode)
		}
	}
}

// onDisbandClub 解散俱乐部
func onDisbandClub(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("userID")
	clubID := r.URL.Query().Get("clubID")

	log.Printf("onDisbandClub, userID:%s, clubID:%s", userID, clubID)

	if clubID == "" {
		log.Println("onDisbandClub, need clubID")
		sendGenericError(w, ClubOperError_CERR_Invalid_Input_Parameter)
		return
	}

	club, ok := clubMgr.clubs[clubID]
	if !ok {
		log.Println("onDisbandClub, no club found for clubID:", clubID)
		sendGenericError(w, ClubOperError_CERR_Invalid_Input_Parameter)
		return
	}

	clname := club.clubInfo.GetBaseInfo().GetClubName()

	mySQLUtil := lobby.MySQLUtil()
	role := mySQLUtil.LoadUserClubRole(userID, clubID)
	if role == int32(ClubRoleType_CRoleTypeNone) {
		log.Printf("onDisbandClub, user %s not in club %s", userID, clubID)
		sendGenericError(w, ClubOperError_CERR_User_Not_In_Club)
		return
	}

	if role != int32(ClubRoleType_CRoleTypeCreator) {
		log.Printf("onDisbandClub, user %s no owner, can not disband club", userID)
		sendGenericError(w, ClubOperError_CERR_Club_Only_Owner_Can_Disband)
		return
	}

	// TODO: 拉取所有成员，给他们发通知
	memberIDs := mySQLUtil.LoadClubUserIDs(clubID)

	redisClearClubData(clubID)

	// 如果club拥有房间，强制解散所有房间，不管是否空闲
	forceDisabandClubRooms(clubID)

	mySQLUtil.DeleteClub(clubID)

	delete(clubMgr.clubs, clubID)


	// if clubBusyRoomCount(clubID) > 0 {
	// 	log.Printf("onDisbandClub, club has playing rooms %s\n", clubID)
	// 	sendGenericError(w, ClubOperError_CERR_Club_Has_Room_In_PlayingState)
	// 	return
	// }

	// // 获取所有俱乐部成员
	// memberIDs, err := redis.Strings(conn.Do("SMEMBERS", stateless.ClubMemberSetPrefix+clubID))
	// if err != nil && err != redis.ErrNil {
	// 	log.Println("onDisbandClub, SMEMBERS redis error:", err)
	// 	sendGenericError(w, ClubOperError_CERR_Database_IO)
	// 	return
	// }

	// applicantIDs, err := redis.Strings(conn.Do("SMEMBERS", stateless.ClubApplicantPrefix+clubID))
	// if err != nil && err != redis.ErrNil {
	// 	log.Println("onDisbandClub, SMEMBERS redis error:", err)
	// 	sendGenericError(w, ClubOperError_CERR_Database_IO)
	// 	return
	// }

	// clname, _ := redis.String(conn.Do("HGET", stateless.ClubTablePrefix+clubID, "clname"))

	// // 清理俱乐部数据
	// err = redisClearClubData(conn, clubID, userID)
	// if err != nil {
	// 	log.Println("onDisbandClub, redisClearClubData, redis err:", err)
	// 	sendGenericError(w, ClubOperError_CERR_Database_IO)
	// 	return
	// }

	// log.Printf("club %s has be disbanded, by userID %s\n", clubID, userID)

	// // 从clubMap中删除
	// deleteClubFromMap(clubID)

	// // 如果club拥有房间，强制解散所有房间，不管是否空闲
	// forceDeleteAllClubRooms(clubID)
	conn := lobby.Pool().Get()
	defer conn.Close()

	nick, _ := redis.String(conn.Do("HGET", gconst.LobbyUserTablePrefix+userID, "Nick"))
	if nick == "" {
		nick = userID
	}
	//  发送邮件给所有俱乐部成员，通知他们俱乐部已经解散
	var text = clname + " 俱乐部已被 " + nick + " 解散!"
	sendClubEventMails(memberIDs, text)

	// //  发送邮件给所有申请者，通知他们其申请的俱乐部已经解散
	applicantIDs, err := redis.Strings(conn.Do("SMEMBERS", gconst.LobbyClubApplicantPrefix+clubID))
	if err != nil && err != redis.ErrNil {
		log.Println("onDisbandClub, SMEMBERS redis error:", err)
		sendGenericError(w, ClubOperError_CERR_Database_IO)
		return
	}
	sendClubEventMails(applicantIDs, text)

	// // 把当前剩余的俱乐部给返回去
	onLoadMyClubs(w, r, ps)
}

// redisClearClubData 删除俱乐部所有相关表格
func redisClearClubData( clubID string) {
	conn := lobby.Pool().Get()
	defer conn.Close()

	luaScriptRemoveMemberEventList.Do(conn, gconst.LobbyClubUnReadEventUserListPrefix+clubID+":",
		gconst.LobbyClubUnReadEventUserSetPrefix+clubID+":", gconst.LobbyClubMemberSetPrefix+clubID)

	conn.Send("MULTI")
	conn.Send("DEL", gconst.LobbyClubMemberSetPrefix+clubID)
	conn.Send("DEL", gconst.LobbyClubApplicantPrefix+clubID)
	conn.Send("DEL", gconst.LobbyClubEventTablePrefix+clubID)
	conn.Send("DEL", gconst.LobbyClubEventListPrefix+clubID)
	conn.Send("DEL", gconst.LobbyClubNeedHandledTablePrefix+clubID)
	conn.Do("EXEC")
}
