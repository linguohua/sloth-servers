package club

import (
	"net/http"
	"gconst"
	"lobbyserver/lobby"
	"github.com/garyburd/redigo/redis"
	"github.com/julienschmidt/httprouter"
	proto "github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

// 审核其他玩家加入俱乐部申请
func onJoinApprove(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("userID")

	log.Println("onJoinApprove, userID:", userID)

	var query = r.URL.Query()
	clubID := query.Get("clubID")
	applicantID := query.Get("applicantID") // 申请者的ID
	agree := query.Get("agree")             // yes表示同意，no表示不同意
	eventID := query.Get("eID")             // 填写对应事件的ID

	if clubID == "" {
		log.Println("onJoinApprove, need club id")
		sendGenericError(w, ClubOperError_CERR_Invalid_Input_Parameter)
		return
	}

	if applicantID == "" {
		log.Println("onJoinApprove, need applicantID")
		sendGenericError(w, ClubOperError_CERR_Invalid_Input_Parameter)
		return
	}

	if agree == "" {
		log.Println("onJoinApprove, need agree")
		sendGenericError(w, ClubOperError_CERR_Invalid_Input_Parameter)
		return
	}

	if eventID == "" {
		log.Println("onJoinApprove, need eventID")
		sendGenericError(w, ClubOperError_CERR_Invalid_Input_Parameter)
		return
	}

	club, ok := clubMgr.clubs[clubID]
	if !ok {
		log.Printf("onJoinApprove, club %s not found", clubID)
		sendGenericError(w, ClubOperError_CERR_Invalid_Input_Parameter)
		return
	}

	mySQLUtil := lobby.MySQLUtil()
	role := mySQLUtil.LoadUserClubRole(applicantID, clubID)
	if role != int32(ClubRoleType_CRoleTypeNone) {
		sendGenericError(w, ClubOperError_CERR_Invitee_Already_In_Club)
		return
	}

	isApplicant := isApplicant(clubID, applicantID)

	// 没申请过
	if !isApplicant {
		log.Printf("onJoinApprove, user %s not in club %s applicant list\n", userID, clubID)
		sendGenericError(w, ClubOperError_CERR_No_Applicant)
		return
	}

	clubInfo := club.clubInfo
	// 只有部长才可以批准或者拒绝
	if clubInfo.GetCreatorUserID() != userID {
		log.Printf("onJoinApprove, userID %s not creator %s\n", userID, clubInfo.GetCreatorUserID())
		sendGenericError(w, ClubOperError_CERR_Only_Creator_Can_Approve)
		return
	}

	conn := lobby.Pool().Get()
	defer conn.Close()

	conn.Send("MULTI")
	// 清理事件
	conn.Send("LREM", gconst.LobbyClubUnReadEventUserListPrefix+clubID+":"+userID, 1, eventID)
	conn.Send("SREM", gconst.LobbyClubUnReadEventUserSetPrefix+clubID+":"+userID, eventID)
	// 从申请者列表中移除
	conn.Send("SREM", gconst.LobbyClubApplicantPrefix+clubID, applicantID)

	_, err := conn.Do("EXEC")
	if err != nil {
		log.Println("onJoinApprove, SREM applicant redis error:", err)
		sendGenericError(w, ClubOperError_CERR_Database_IO)
		return
	}

	// 修改事件的批准相关参数
	appendApprovalResult2Event(eventID, clubID, agree, conn)

	userIDs := []string {applicantID}

	if "yes" != agree {
		// 发邮件通知申请者告知被拒绝
		nick, _ := redis.String(conn.Do("HGET", gconst.LobbyUserTablePrefix+userID, "nick"))
		if nick == "" {
			nick = userID
		}

		fReason := ""
		var text = nick + " 拒绝了您进入俱乐部 " + clubInfo.GetBaseInfo().GetClubName() + " 的申请 ! 拒绝理由: " + fReason
		sendClubEventMails(userIDs, text)
	} else {
		// 检查玩家是否可以加入俱乐部
		// 检查玩家已经加入过的俱乐部个数
		maxJoin := mySQLUtil.CountUserClubNumber(userID)
		if int32(maxJoin) < clubInfo.GetMaxMember() {
			// 清理事件
			_, err := conn.Do("SADD", gconst.LobbyClubMemberSetPrefix+clubID, applicantID)
			if err != nil {
				log.Panicln("save applicant's club info redis failed:", err)
			}

			mySQLUtil.AddUserToClub(applicantID, clubID, int32(ClubRoleType_CRoleTypeMember))

			// 生成玩家加入俱乐部事件
			// newJoinEvent(clubID, applicantID, conn)

			// 给刚加入俱乐部成员 发送一个邮件
			var text = "您成功加入了 " +  clubInfo.GetBaseInfo().GetClubName() + " 俱乐部，赶紧加入俱乐部的牌局吧!"
			sendClubEventMails(userIDs, text)
			// TODO: 给用户发通知

		} else {
			log.Printf("club %s owner %s agree applicant %s to join, but applicant has exceed max join limit\n", clubID,
				userID, applicantID)
		}
	}

	// 操作成功
	sendGenericError(w, ClubOperError_CERR_OK)
}

func appendApprovalResult2Event(eventID string, clubID string, agree string, conn redis.Conn) {
	b, err := redis.Bytes(conn.Do("HGET", gconst.LobbyClubEventTablePrefix+clubID, eventID))
	if err != nil {
		log.Panicln("appendApprovalResult2Event, convert value to bytes failed:", err)
	}

	e := &MsgClubEvent{}
	err = proto.Unmarshal(b, e)
	if err != nil {
		log.Panicln("appendApprovalResult2Event, unmarshal bytes to event failed:", err)
	}

	var approvalResult32 int32
	if agree == "yes" {
		approvalResult32 = 1
	} else {
		approvalResult32 = 2
	}

	e.ApprovalResult = &approvalResult32
	b, err = proto.Marshal(e)
	if err != nil {
		log.Panicln("appendApprovalResult2Event, Marshal event failed:", err)
	}

	conn.Do("HSET", gconst.LobbyClubEventTablePrefix+clubID, eventID, b)
}
