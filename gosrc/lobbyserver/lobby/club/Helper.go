package club

import (
	"log"
	"net/http"
	"time"
	"lobbyserver/lobby"
	"github.com/golang/protobuf/proto"
)
// sendGenericError 发送一个错误码到客户端
func sendGenericError(w http.ResponseWriter, errCode ClubOperError) {
	gr := MsgCubOperGenericReply{}
	var err32 = int32(errCode)
	gr.ErrorCode = &err32

	b, err := proto.Marshal(&gr)
	if err != nil {
		log.Println("sendGenericReply, marshal error:", err)
		return
	}

	cr := &MsgClubReply{}
	var replyCode = int32(ClubReplyCode_RCError)
	cr.ReplyCode = &replyCode
	cr.Content = b

	b2, err := proto.Marshal(cr)
	if err != nil {
		log.Println("sendGenericReply, marshal error:", err)
		return
	}

	w.Write(b2)
}

func sendGenericErrorWithExtraString(w http.ResponseWriter, errCode ClubOperError, extra string) {
	gr := MsgCubOperGenericReply{}
	var err32 = int32(errCode)
	gr.ErrorCode = &err32
	gr.Extra = &extra

	b, err := proto.Marshal(&gr)
	if err != nil {
		log.Println("sendGenericErrorWithExtraString, marshal error:", err)
		return
	}

	cr := &MsgClubReply{}
	var replyCode = int32(ClubReplyCode_RCError)
	cr.ReplyCode = &replyCode
	cr.Content = b

	b2, err := proto.Marshal(cr)
	if err != nil {
		log.Println("sendGenericErrorWithExtraString, marshal error:", err)
		return
	}

	w.Write(b2)
}

func sendMsgClubReply(w http.ResponseWriter, rc ClubReplyCode, b []byte) {
	cr := &MsgClubReply{}
	var replyCode = int32(rc)
	cr.ReplyCode = &replyCode
	cr.Content = b

	b2, err := proto.Marshal(cr)
	if err != nil {
		log.Println("sendGenericReply, marshal error:", err)
		return
	}

	w.Write(b2)
}

// strArray2Comma 字符串数据转为逗号分隔字符串
func strArray2Comma(ss []string) string {
	result := ""
	if len(ss) < 1 {
		return result
	}

	for i := 0; i < len(ss)-1; i++ {
		result = result + ss[i] + ","
	}

	result = result + ss[len(ss)-1]

	return result
}

// stringArrayRemove 从字符串数组中删除一个元素
func stringArrayRemove(ss []string, s string) []string {
	for i, v := range ss {
		if v == s {
			ss = append(ss[:i], ss[i+1:]...)
			return ss
		}
	}

	return ss
}


func unixTimeInSeconsSince2010() uint32 {
	return uint32(time.Now().Unix() / 60)
}

//userIDs 需要通知的用户的ID   clubName 操作的俱乐部名字   userName 操作的管理员名字  ev 事件类型
func sendClubEventMails(userIDs []string, text string) {
	// var text = "未知消息类型"
	// if ev == ClubEventType_CEVT_Deny {
	// 	//申请被拒绝
	// 	text = "您申请加入 " + clubName + " 俱乐部被 " + userName + " 拒绝!"
	// } else if ev == ClubEventType_CEVT_ClubDisband {
	// 	//解散
	// 	text = clubName + " 俱乐部被 " + userName + " 解散!"
	// } else if ev == ClubEventType_CEVT_Kickout {
	// 	//被踢
	// 	text = "您被 " + userName + " 踢出 " + clubName + " 俱乐部!"
	// } else if ev == ClubEventType_CEVT_Approval {
	// 	//同意加入
	// 	text = "您成功加入了 " + clubName + " 俱乐部，赶紧加入俱乐部的牌局吧!"
	// }
	// var myMail = &webdata.Mail{}
	// myMail.GameID = 10088
	// //用户ID
	// myMail.Type = 1
	// myMail.Subject = "俱乐部邮件"
	// myMail.Title = "俱乐部"
	// myMail.Text = text
	// myMail.ExpirationTime = "2018-04-04 14:25:22"
	title := "牌友圈邮件"

	mailUtil := lobby.MailUtil()
	for _, userID := range userIDs {
		mailUtil.SendMail(userID, text, title)
	}
}
