package mysql

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"time"
)

// GRCRecord 牌局记录
type GRCRecord struct {
	RecordID     string
	RoomID       string
	RoomConfigID string
	ShareID      string
	WriteTime    string
	RecordData   []byte
}

// SaveGRCRecord2SqlServer 保存牌局记录到数据库
func SaveGRCRecord2SqlServer(grcRecordGUID string, roomID string, configID string,
	sharedID string, data []byte, conn *sql.DB) error {
	/* Description:	牌局录像写入数据库  目前对应关系写入在REDIS中 只需要将牌局录像写入方便 回放查看
	 [dbo].[PrPsWeb_User_Replay_Add]
		@IV_GameId               INT,	                    --游戏编号
		@IV_RePlayNo             NVARCHAR(40),	            --回放主键码唯一
		@IV_RoomNo               NVARCHAR(40),				--房间号
		@IV_RoomConfig           NVARCHAR(128),				--规则
		@IV_ReplayData           NVARCHAR(MAX),				--回放数据默认10k
		@IV_ShareID              NVARCHAR(40),              --回放码
		@IV_WriteTime            DATETIME,	                --记录时间
		@OV_ReturnMsg            VARCHAR(128) OUTPUT        --返回消息
	*/
	query := `
		declare @IV_GameID int;
		declare @IV_RePlayNo nvarchar(40);
		declare @IV_RoomNo nvarchar(40);
		declare @IV_RoomConfig nvarchar(128);
		declare @IV_ReplayData nvarchar(MAX);
		declare @IV_ShareID nvarchar(40);
		declare @IV_WriteTime datetime;
		declare @OV_ReturnMsg  varchar(128);
		declare @OV_Result int;

		set @IV_GameID = %d;
		set @IV_RePlayNo = '%s';
		set @IV_RoomNo = '%s';
		set @IV_RoomConfig = '%s';
		set @IV_ReplayData = '%x'
		set @IV_ShareID = '%s';
		set @IV_WriteTime = '%s';

		exec @OV_Result = PrPsWeb_User_Replay_Add @IV_GameID, @IV_RePlayNo, @IV_RoomNo, @IV_RoomConfig, @IV_ReplayData, @IV_ShareID, @IV_WriteTime, @OV_ReturnMsg output;
		select @OV_Result result, @OV_ReturnMsg msg;
		`

	// userID, _ := strconv.ParseInt(playerID, 10, 64)
	var t = time.Now()
	query = fmt.Sprintf(query, 10088, grcRecordGUID, roomID, configID, data, sharedID, t.Format("2006-01-02 15:04:05"))
	// log.Println("query:", query)
	stmt, err := conn.Prepare(query)
	if err != nil {
		log.Println("SaveGRCRecord2SqlServer Prepare Err:", err)
		return err
	}
	defer stmt.Close()

	var result int
	var msg sql.NullString
	row := stmt.QueryRow()
	err = row.Scan(&result, &msg)
	if err != nil {
		log.Println("SaveGRCRecord2SqlServer Scan Err:", err)
		return err
	}

	if result != 1 {
		log.Println("SaveGRCRecord2SqlServer Err:", msg.String)
		return fmt.Errorf(msg.String)
	}

	return nil

}

// DeleteGRCRecordFromSQLServer 删除牌局记录
func DeleteGRCRecordFromSQLServer(key string) {

}

// LoadGRCRcordFromSQLServer 加载牌局记录
func LoadGRCRcordFromSQLServer(grcRecordGUID string, conn *sql.DB) *GRCRecord {
	/*
		[dbo].[PrPsWeb_User_Replay_Get]
		@IV_GameId               INT,	                    --游戏编号
		@IV_RePlayNo             NVARCHAR(40),	            --回放主键码唯一
		@OV_ReturnMsg            VARCHAR(128) OUTPUT       --返回消息
	*/

	query := `
		declare @IV_GameID int;
		declare @IV_RePlayNo nvarchar(40);
		declare @OV_ReturnMsg  varchar(128);
		declare @OV_Result int;

		set @IV_GameID = %d;
		set @IV_RePlayNo = '%s';

		exec @OV_Result = PrPsWeb_User_Replay_Get @IV_GameID, @IV_RePlayNo, @OV_ReturnMsg output;
		select @OV_Result result, @OV_ReturnMsg returnMsg;
		`

	query = fmt.Sprintf(query, 10088, grcRecordGUID)
	// log.Println(query)
	stmt, err := conn.Prepare(query)
	if err != nil {
		log.Println("LoadSharedIDFromSQLServer Prepare Err:", err)
		return nil
	}
	defer stmt.Close()

	var replayNo sql.NullString
	var roomNo sql.NullString
	// roomConfig是roomConfigID
	var roomConfig sql.NullString
	var replayData sql.NullString
	var shareID sql.NullString
	var writeTime sql.NullString

	row := stmt.QueryRow()
	err = row.Scan(&replayNo, &roomNo, &roomConfig, &replayData, &shareID, &writeTime)
	if err != nil {
		log.Println("err:", err)
		return nil
	}

	var grcRecord = &GRCRecord{}
	grcRecord.RecordID = replayNo.String
	grcRecord.RoomID = roomNo.String
	grcRecord.RoomConfigID = roomConfig.String
	grcRecord.ShareID = shareID.String
	grcRecord.WriteTime = writeTime.String

	dst := make([]byte, hex.DecodedLen(len(replayData.String)))
	_, err = hex.Decode(dst, []byte(replayData.String))
	if err != nil {
		log.Fatal(err)
	}
	grcRecord.RecordData = dst

	return grcRecord
}
