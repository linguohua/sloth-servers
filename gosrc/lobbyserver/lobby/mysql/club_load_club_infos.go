package mysql

import (
	"database/sql"
	"lobbyserver/lobby/club"
	"time"

	log "github.com/sirupsen/logrus"
)

// 注意：这个只在服务器启动的时候才去拉取，其它时候不应该调用此方法
func loadClubInfos(cursor int, count int) (clubInfo []*club.MsgClubInfo) {
	stmt, err := dbConn.Prepare("select club_id, club_num, club_name, creator, is_league, points, wanka, candy, max_member, create_time from club limit ?, ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query(cursor, count)
	if err != nil {
		panic(err.Error())
	}

	clubInfos := make([]*club.MsgClubInfo, 0)
	for rows.Next() {
		var myClubID sql.NullString
		var clubNum sql.NullString
		var clubName sql.NullString
		var creator sql.NullString
		var isLeague sql.NullInt64
		var points sql.NullInt64
		var wanka sql.NullInt64
		var candy sql.NullInt64
		var maxMember sql.NullInt64
		var rawTime []byte

		err = rows.Scan(&myClubID, &clubNum, &clubName, &creator, &isLeague, &points, &wanka, &candy, &maxMember, &rawTime)
		if err != nil {
			panic(err.Error())
		}

		myClubInfo := &club.MsgClubInfo{}
		clubBaseInfo := &club.MsgClubBaseInfo{}
		myClubInfo.BaseInfo = clubBaseInfo

		if myClubID.Valid {
			clubBaseInfo.ClubID = &myClubID.String
		}

		if clubNum.Valid {
			clubBaseInfo.ClubNumber = &clubNum.String
		}

		if clubName.Valid {
			clubBaseInfo.ClubName = &clubName.String
		}

		if creator.Valid {
			myClubInfo.CreatorUserID = &creator.String
		}

		if isLeague.Valid {
			clubLevel := int32(isLeague.Int64)
			myClubInfo.ClubLevel = &clubLevel
		}

		if points.Valid {
			pointsInt32 := int32(points.Int64)
			myClubInfo.Points = &pointsInt32
		}

		if wanka.Valid {
			wankaInt32 := int32(wanka.Int64)
			myClubInfo.Wanka = &wankaInt32
		}

		if candy.Valid {
			candyInt32 := int32(candy.Int64)
			myClubInfo.Candy = &candyInt32
		}

		if maxMember.Valid {
			maxMemberInt32 := int32(maxMember.Int64)
			myClubInfo.MaxMember = &maxMemberInt32
		}

		if len(rawTime) > 0 {
			createTime, err := time.Parse("2006-01-02 15:04:05", string(rawTime))
			if err == nil {
				creatime := int32(createTime.Unix())
				myClubInfo.CreateTime = &creatime
			} else {
				log.Error("loadClubInfos, err:", err)
			}

		}

		clubInfos = append(clubInfos, myClubInfo)

	}

	return clubInfos
}
