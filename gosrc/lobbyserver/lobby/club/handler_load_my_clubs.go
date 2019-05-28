package club

import (
	"gconst"
	"net/http"
	"lobbyserver/lobby"
	"github.com/julienschmidt/httprouter"
	"github.com/garyburd/redigo/redis"
	proto "github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

// onLoadMyClubs 加载自己的俱乐部
func onLoadMyClubs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("userID")

	log.Println("onLoadMyClubs, userID:", userID)

	// TODO: 用户的牌友群是否可以放在redis中
	mySQLUtil := lobby.MySQLUtil()
	clubIDs := mySQLUtil.LoadUserClubIDs(userID)

	reply := &MsgClubLoadMyClubsReply{}

	if len(clubIDs) > 0 {
		log.Printf("user %s load clubs, count:%d\n", userID, len(clubIDs))

		msgClubInfos := make([]*MsgClubInfo, 0, len(clubIDs))
		// removedClubIds := make([]string, 0, len(clubIDs))
		// check and load club from redis if need
		for _, cid := range clubIDs {
			if cid == "" {
				continue
			}

			// 检查牌友群是否在列表中，否则从数据库加载
			club, ok := clubMgr.clubs[cid]
			if !ok {
				clubInfo := mySQLUtil.LoadClubInfo(cid)
				if clubInfo == nil {
					log.Panic("onLoadMyClubs, no club for clubID:", cid)
				}

				club = newBaseClub(clubInfo.(*MsgClubInfo), cid)
				clubMgr.clubs[cid] = club
			}

			// clubInfo := club.constructMsgClubInfo()
			msgClubInfos = append(msgClubInfos, club.clubInfo)
		}

		// if len(removedClubIds) > 0 {
		// 	conn.Send("MULTI")
		// 	// 表明已经有俱乐部不存在，因此需要删除掉
		// 	for _, v := range removedClubIds {
		// 		conn.Send("SREM", stateless.PlayerClubSetPrefix+userID, v)
		// 	}
		// 	conn.Do("EXEC")
		// }

		// 检查是否未读事件
		if len(msgClubInfos) > 0 {
			conn := lobby.Pool().Get()
			defer conn.Close()

			conn.Send("MULTI")
			// 表明已经有俱乐部不存在，因此需要删除掉
			for _, cinfo := range msgClubInfos {
				clubID := cinfo.BaseInfo.GetClubID()
				conn.Send("SCARD", gconst.LobbyClubUnReadEventUserSetPrefix+clubID+":"+userID)
			}

			ints, err := redis.Ints(conn.Do("EXEC"))

			if err == nil {
				for i, cinfo := range msgClubInfos {
					hasUnReadEvents := false
					if ints[i] > 0 {
						hasUnReadEvents = true
					}
					cinfo.HasUnReadEvents = &hasUnReadEvents
				}
			}
		}

		reply.Clubs = msgClubInfos
	}

	b, err := proto.Marshal(reply)
	if err != nil {
		log.Println("onCreateClub, marshal error:", err)
		sendGenericError(w, ClubOperError_CERR_Encode_Decode)
		return
	}

	sendMsgClubReply(w, ClubReplyCode_RCOperation, b)
}
