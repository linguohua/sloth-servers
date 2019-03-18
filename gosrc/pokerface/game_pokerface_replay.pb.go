// Code generated by protoc-gen-go. DO NOT EDIT.
// source: game_pokerface_replay.proto

package pokerface

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// 回放房间中的玩家信息
type MsgReplayPlayerInfo struct {
	UserID           *string `protobuf:"bytes,1,req,name=userID" json:"userID,omitempty"`
	Nick             *string `protobuf:"bytes,2,opt,name=nick" json:"nick,omitempty"`
	ChairID          *int32  `protobuf:"varint,3,req,name=chairID" json:"chairID,omitempty"`
	TotalScore       *int32  `protobuf:"varint,4,opt,name=totalScore" json:"totalScore,omitempty"`
	Sex              *uint32 `protobuf:"varint,5,opt,name=sex" json:"sex,omitempty"`
	HeadIconURI      *string `protobuf:"bytes,6,opt,name=headIconURI" json:"headIconURI,omitempty"`
	AvatarID         *int32  `protobuf:"varint,7,opt,name=avatarID" json:"avatarID,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *MsgReplayPlayerInfo) Reset()                    { *m = MsgReplayPlayerInfo{} }
func (m *MsgReplayPlayerInfo) String() string            { return proto.CompactTextString(m) }
func (*MsgReplayPlayerInfo) ProtoMessage()               {}
func (*MsgReplayPlayerInfo) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *MsgReplayPlayerInfo) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *MsgReplayPlayerInfo) GetNick() string {
	if m != nil && m.Nick != nil {
		return *m.Nick
	}
	return ""
}

func (m *MsgReplayPlayerInfo) GetChairID() int32 {
	if m != nil && m.ChairID != nil {
		return *m.ChairID
	}
	return 0
}

func (m *MsgReplayPlayerInfo) GetTotalScore() int32 {
	if m != nil && m.TotalScore != nil {
		return *m.TotalScore
	}
	return 0
}

func (m *MsgReplayPlayerInfo) GetSex() uint32 {
	if m != nil && m.Sex != nil {
		return *m.Sex
	}
	return 0
}

func (m *MsgReplayPlayerInfo) GetHeadIconURI() string {
	if m != nil && m.HeadIconURI != nil {
		return *m.HeadIconURI
	}
	return ""
}

func (m *MsgReplayPlayerInfo) GetAvatarID() int32 {
	if m != nil && m.AvatarID != nil {
		return *m.AvatarID
	}
	return 0
}

// 回放记录中玩家的得分信息
type MsgReplayPlayerScoreSummary struct {
	ChairID          *int32 `protobuf:"varint,1,req,name=chairID" json:"chairID,omitempty"`
	Score            *int32 `protobuf:"varint,2,req,name=score" json:"score,omitempty"`
	WinType          *int32 `protobuf:"varint,3,req,name=winType" json:"winType,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *MsgReplayPlayerScoreSummary) Reset()                    { *m = MsgReplayPlayerScoreSummary{} }
func (m *MsgReplayPlayerScoreSummary) String() string            { return proto.CompactTextString(m) }
func (*MsgReplayPlayerScoreSummary) ProtoMessage()               {}
func (*MsgReplayPlayerScoreSummary) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

func (m *MsgReplayPlayerScoreSummary) GetChairID() int32 {
	if m != nil && m.ChairID != nil {
		return *m.ChairID
	}
	return 0
}

func (m *MsgReplayPlayerScoreSummary) GetScore() int32 {
	if m != nil && m.Score != nil {
		return *m.Score
	}
	return 0
}

func (m *MsgReplayPlayerScoreSummary) GetWinType() int32 {
	if m != nil && m.WinType != nil {
		return *m.WinType
	}
	return 0
}

// 手牌回放记录概要
type MsgReplayRecordSummary struct {
	RecordUUID       *string                        `protobuf:"bytes,1,req,name=recordUUID" json:"recordUUID,omitempty"`
	PlayerScores     []*MsgReplayPlayerScoreSummary `protobuf:"bytes,2,rep,name=playerScores" json:"playerScores,omitempty"`
	EndTime          *uint32                        `protobuf:"varint,3,req,name=endTime" json:"endTime,omitempty"`
	ShareAbleID      *string                        `protobuf:"bytes,4,opt,name=shareAbleID" json:"shareAbleID,omitempty"`
	StartTime        *uint32                        `protobuf:"varint,5,opt,name=startTime" json:"startTime,omitempty"`
	XXX_unrecognized []byte                         `json:"-"`
}

func (m *MsgReplayRecordSummary) Reset()                    { *m = MsgReplayRecordSummary{} }
func (m *MsgReplayRecordSummary) String() string            { return proto.CompactTextString(m) }
func (*MsgReplayRecordSummary) ProtoMessage()               {}
func (*MsgReplayRecordSummary) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{2} }

func (m *MsgReplayRecordSummary) GetRecordUUID() string {
	if m != nil && m.RecordUUID != nil {
		return *m.RecordUUID
	}
	return ""
}

func (m *MsgReplayRecordSummary) GetPlayerScores() []*MsgReplayPlayerScoreSummary {
	if m != nil {
		return m.PlayerScores
	}
	return nil
}

func (m *MsgReplayRecordSummary) GetEndTime() uint32 {
	if m != nil && m.EndTime != nil {
		return *m.EndTime
	}
	return 0
}

func (m *MsgReplayRecordSummary) GetShareAbleID() string {
	if m != nil && m.ShareAbleID != nil {
		return *m.ShareAbleID
	}
	return ""
}

func (m *MsgReplayRecordSummary) GetStartTime() uint32 {
	if m != nil && m.StartTime != nil {
		return *m.StartTime
	}
	return 0
}

// 回播房间记录
type MsgReplayRoom struct {
	RecordRoomType   *int32                    `protobuf:"varint,1,req,name=recordRoomType" json:"recordRoomType,omitempty"`
	StartTime        *uint32                   `protobuf:"varint,2,req,name=startTime" json:"startTime,omitempty"`
	EndTime          *uint32                   `protobuf:"varint,3,req,name=endTime" json:"endTime,omitempty"`
	RoomNumber       *string                   `protobuf:"bytes,4,req,name=roomNumber" json:"roomNumber,omitempty"`
	Players          []*MsgReplayPlayerInfo    `protobuf:"bytes,5,rep,name=players" json:"players,omitempty"`
	Records          []*MsgReplayRecordSummary `protobuf:"bytes,6,rep,name=records" json:"records,omitempty"`
	OwnerUserID      *string                   `protobuf:"bytes,7,opt,name=ownerUserID" json:"ownerUserID,omitempty"`
	XXX_unrecognized []byte                    `json:"-"`
}

func (m *MsgReplayRoom) Reset()                    { *m = MsgReplayRoom{} }
func (m *MsgReplayRoom) String() string            { return proto.CompactTextString(m) }
func (*MsgReplayRoom) ProtoMessage()               {}
func (*MsgReplayRoom) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{3} }

func (m *MsgReplayRoom) GetRecordRoomType() int32 {
	if m != nil && m.RecordRoomType != nil {
		return *m.RecordRoomType
	}
	return 0
}

func (m *MsgReplayRoom) GetStartTime() uint32 {
	if m != nil && m.StartTime != nil {
		return *m.StartTime
	}
	return 0
}

func (m *MsgReplayRoom) GetEndTime() uint32 {
	if m != nil && m.EndTime != nil {
		return *m.EndTime
	}
	return 0
}

func (m *MsgReplayRoom) GetRoomNumber() string {
	if m != nil && m.RoomNumber != nil {
		return *m.RoomNumber
	}
	return ""
}

func (m *MsgReplayRoom) GetPlayers() []*MsgReplayPlayerInfo {
	if m != nil {
		return m.Players
	}
	return nil
}

func (m *MsgReplayRoom) GetRecords() []*MsgReplayRecordSummary {
	if m != nil {
		return m.Records
	}
	return nil
}

func (m *MsgReplayRoom) GetOwnerUserID() string {
	if m != nil && m.OwnerUserID != nil {
		return *m.OwnerUserID
	}
	return ""
}

func init() {
	proto.RegisterType((*MsgReplayPlayerInfo)(nil), "pokerface.MsgReplayPlayerInfo")
	proto.RegisterType((*MsgReplayPlayerScoreSummary)(nil), "pokerface.MsgReplayPlayerScoreSummary")
	proto.RegisterType((*MsgReplayRecordSummary)(nil), "pokerface.MsgReplayRecordSummary")
	proto.RegisterType((*MsgReplayRoom)(nil), "pokerface.MsgReplayRoom")
}

func init() { proto.RegisterFile("game_pokerface_replay.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 378 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xcd, 0x0e, 0x9a, 0x40,
	0x14, 0x85, 0x33, 0x28, 0x52, 0x2e, 0xfe, 0xb4, 0x63, 0x62, 0x26, 0x31, 0x69, 0x28, 0x8b, 0x86,
	0x95, 0x4d, 0x5c, 0x77, 0xd3, 0xc4, 0x0d, 0x8b, 0x9a, 0x46, 0x65, 0x6d, 0x46, 0xbc, 0x2a, 0x11,
	0x18, 0x32, 0x83, 0xb5, 0xbc, 0x43, 0xdf, 0xa2, 0x4f, 0xd5, 0xb7, 0x69, 0xb8, 0xa4, 0x58, 0xfb,
	0x93, 0x2e, 0x67, 0x86, 0x7b, 0xce, 0xf9, 0xce, 0x05, 0xe6, 0x67, 0x99, 0xe3, 0xbe, 0x54, 0x57,
	0xd4, 0x27, 0x99, 0xe0, 0x5e, 0x63, 0x99, 0xc9, 0x7a, 0x51, 0x6a, 0x55, 0x29, 0xee, 0x76, 0xf7,
	0xc1, 0x57, 0x06, 0xd3, 0x8f, 0xe6, 0xbc, 0xa1, 0xe7, 0x4f, 0x99, 0xac, 0x51, 0x47, 0xc5, 0x49,
	0xf1, 0x31, 0x0c, 0x6e, 0x06, 0x75, 0xb4, 0x12, 0xcc, 0xb7, 0x42, 0x97, 0x0f, 0xa1, 0x5f, 0xa4,
	0xc9, 0x55, 0x58, 0x3e, 0x0b, 0x5d, 0x3e, 0x01, 0x27, 0xb9, 0xc8, 0xb4, 0x79, 0xee, 0xf9, 0x56,
	0x68, 0x73, 0x0e, 0x50, 0xa9, 0x4a, 0x66, 0xdb, 0x44, 0x69, 0x14, 0x7d, 0x9f, 0x85, 0x36, 0xf7,
	0xa0, 0x67, 0xf0, 0x8b, 0xb0, 0x7d, 0x16, 0x8e, 0xf8, 0x14, 0xbc, 0x0b, 0xca, 0x63, 0x94, 0xa8,
	0x22, 0xde, 0x44, 0x62, 0x40, 0x32, 0x2f, 0xe1, 0x85, 0xfc, 0x2c, 0x2b, 0xd9, 0xe8, 0x38, 0xcd,
	0x4c, 0xb0, 0x86, 0xf9, 0x6f, 0x69, 0x48, 0x71, 0x7b, 0xcb, 0x73, 0xa9, 0xeb, 0x5f, 0x7d, 0x19,
	0xf9, 0x8e, 0xc0, 0x36, 0x64, 0x69, 0xd1, 0x71, 0x02, 0xce, 0x3d, 0x2d, 0x76, 0x75, 0x89, 0x6d,
	0xae, 0xe0, 0x1b, 0x83, 0x59, 0x27, 0xb8, 0xc1, 0x44, 0xe9, 0xe3, 0x4f, 0x2d, 0x0e, 0xa0, 0xe9,
	0x22, 0x8e, 0x3b, 0xca, 0xf7, 0x30, 0x2c, 0x1f, 0xae, 0x46, 0x58, 0x7e, 0x2f, 0xf4, 0x96, 0x6f,
	0x17, 0x5d, 0x5f, 0x8b, 0xff, 0xa4, 0xc3, 0xe2, 0xb8, 0x4b, 0xf3, 0xd6, 0x9d, 0xa0, 0xcd, 0x45,
	0x6a, 0xfc, 0x70, 0xc8, 0x30, 0x5a, 0x51, 0x2d, 0x2e, 0x7f, 0x05, 0xae, 0xa9, 0xa4, 0xae, 0xe8,
	0x3b, 0x2a, 0x27, 0xf8, 0xce, 0x60, 0xf4, 0x48, 0xa9, 0x54, 0xce, 0x67, 0x30, 0x6e, 0xc3, 0x35,
	0x27, 0xe2, 0x69, 0x79, 0x9f, 0x86, 0x2d, 0x32, 0xf9, 0xc3, 0xb5, 0x01, 0x53, 0x2a, 0x5f, 0xdf,
	0xf2, 0x03, 0x6a, 0xd1, 0x27, 0xb0, 0x77, 0xe0, 0xb4, 0x60, 0x46, 0xd8, 0xc4, 0xf4, 0xfa, 0xdf,
	0x4c, 0xb4, 0xff, 0x25, 0x38, 0x6d, 0x00, 0x23, 0x06, 0x34, 0xf0, 0xe6, 0x6f, 0x03, 0xcf, 0x8d,
	0x4e, 0xc1, 0x53, 0xf7, 0x02, 0x75, 0xdc, 0xfe, 0x38, 0xcd, 0x46, 0xdd, 0x1f, 0x01, 0x00, 0x00,
	0xff, 0xff, 0x2f, 0xfb, 0x6a, 0x9a, 0x89, 0x02, 0x00, 0x00,
}
