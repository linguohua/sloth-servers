package club

import (
	"net/http"
	"lobbyserver/lobby"
	"github.com/golang/protobuf/proto"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// onCreateClub 创建一个新俱乐部
func onCreateClub(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("userID")
	log.Println("onCreateClub:", userID)

	var query = r.URL.Query()
	clname := query.Get("clname")
	if clname == "" {
		log.Println("onCreateClub, need clname")
		sendGenericError(w, ClubOperError_CERR_Invalid_Input_Parameter)
		return
	}

	if len(clname) > 20 {
		log.Println("onCreateClub, clname too long:", len(clname))
		sendGenericError(w, ClubOperError_CERR_Club_Name_Too_Long)
		return
	}


	mySQLUtil := lobby.MySQLUtil()
	clubCount := mySQLUtil.CountUserClubNumber(userID)

	// 检查是否数量超过上限
	if clubCount >= maxClubCreatePerUser {
		log.Printf("onCreateClub, user %s already has %d clubs, can't create new\n", userID, clubCount)
		sendGenericError(w, ClubOperError_CERR_Exceed_Max_Club_Count_Limit)
		return
	}

	clubID, clubNumber, errCode := mySQLUtil.CreateClub(clname, userID, defaultLevel, defaultWanka, defaultCandy, maxMemberPerClub)
	if errCode != 0 {
		log.Panic("Create club error, errCode:", errCode)
	}

	clubBaseInfo := &MsgClubBaseInfo{}
	clubBaseInfo.ClubID = &clubID
	clubBaseInfo.ClubNumber = &clubNumber
	clubBaseInfo.ClubName = &clname

	clubInfo := &MsgClubInfo{}
	clubInfo.BaseInfo = clubBaseInfo
	clubLevel := int32(defaultLevel)
	clubInfo.ClubLevel = &clubLevel
	clubInfo.CreatorUserID = &userID
	wanka := int32(defaultWanka)
	clubInfo.Wanka = &wanka
	candy := int32(defaultCandy)
	clubInfo.Candy = &candy
	maxMember := int32(maxMemberPerClub)
	clubInfo.MaxMember = &maxMember

	club := newBaseClub(clubInfo, clubID)
	log.Println("club:", club)
	_, ok := clubMgr.clubs[club.ID]
	if !ok {
		clubMgr.clubs[club.ID] = club
	}

	cr := &MsgCreateClubReply{}
	cr.ClubInfo = clubInfo

	b, err := proto.Marshal(cr)
	if err != nil {
		log.Println("onCreateClub, marshal error:", err)
		sendGenericError(w, ClubOperError_CERR_Encode_Decode)
		return
	}

	sendMsgClubReply(w, ClubReplyCode_RCOperation, b)
}
