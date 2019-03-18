package lobby

import (
	"encoding/xml"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	errUserIDIsEmpty = 1
)

func handleUpdateDiamond(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength < 1 {
		log.Println("parseAccessoryMessage failed, content length is zero")
		return
	}

	message := make([]byte, r.ContentLength)
	n, _ := r.Body.Read(message)
	if n != int(r.ContentLength) {
		log.Println("parseAccessoryMessage failed, can't read request body")
		return
	}

	log.Println("message:", string(message))
	//var buf = []byte(`<?xml version='1.0' encoding='utf-8'?><req code='sys_sy'><msg code="4" id="1"><prop type="2" userid="10021619" TypeValue="10" fee="500" gameID="10888" propsID="1000" CYMoney="0" TypeID="0" charge="0" /></msg></req>`)

	type MyProp struct {
		Type      string `xml:"type,attr"`
		UserID    string `xml:"userid,attr"`
		TypeValue string `xml:"TypeValue,attr"`
		Fee       string `xml:"fee,attr"`
		GameID    string `xml:"gameID,attr"`
		PropsID   string `xml:"propsID,attr"`
		CYMoney   string `xml:"CYMoney,attr"`
		TypeID    string `xml:"TypeID,attr"`
		Charge    string `xml:"charge,attr"`
		EmailID   string `xml:"emailID,attr"`
	}

	type Message struct {
		Code string `xml:"code,attr"`
		ID   string `xml:"id,attr"`
		Prop MyProp `xml:"prop"`
	}

	type Req struct {
		Code string  `xml:"code,attr"`
		Msg  Message `xml:"msg"`
	}

	req := Req{}
	err := xml.Unmarshal(message, &req)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	log.Println("req:", req)

	msg := req.Msg
	prop := msg.Prop
	userID := prop.UserID
	if userID == "" {
		errMsg := fmt.Sprintf(`{"errorCode":%d}`, errUserIDIsEmpty)
		w.Write([]byte(errMsg))
		return
	}

	log.Println("userID:", prop.UserID)

	// 回复购买钻石服务器
	result := []byte(`{"errorCode":0}`)
	w.Write(result)

	// 更新用户钻石
	user := userMgr.getUserByID(userID)
	if user == nil {
		log.Println("user offline")
		return
	}

	// TODO: llwant mysql
	// diamond, err := webdata.QueryDiamond(userID)
	// if err != nil {
	// 	log.Println("can't get user diamond")
	// 	return
	// }

	user.updateMoney(uint32(0))
	// var updateUserMoney = &MsgUpdateUserMoney{}
	// var userDiamond = uint32(diamond)
	// updateUserMoney.Diamond = &userDiamond

	// if prop.EmailID != "" {
	// 	var activityType = uint32(ActivityType_Email)
	// 	updateUserMoney.ActivityType = &activityType
	// }

	// user.sendMsg(updateUserMoney, int32(MessageCode_OPUpdateUserMoney))
}
