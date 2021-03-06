// Code generated by protoc-gen-go. DO NOT EDIT.
// source: game_pokerface_replay.proto

package pokerface

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// 回放房间中的玩家信息
type MsgReplayPlayerInfo struct {
	UserID               *string  `protobuf:"bytes,1,req,name=userID" json:"userID,omitempty"`
	Nick                 *string  `protobuf:"bytes,2,opt,name=nick" json:"nick,omitempty"`
	ChairID              *int32   `protobuf:"varint,3,req,name=chairID" json:"chairID,omitempty"`
	TotalScore           *int32   `protobuf:"varint,4,opt,name=totalScore" json:"totalScore,omitempty"`
	Gender               *uint32  `protobuf:"varint,5,opt,name=gender" json:"gender,omitempty"`
	HeadIconURI          *string  `protobuf:"bytes,6,opt,name=headIconURI" json:"headIconURI,omitempty"`
	AvatarID             *int32   `protobuf:"varint,7,opt,name=avatarID" json:"avatarID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MsgReplayPlayerInfo) Reset()         { *m = MsgReplayPlayerInfo{} }
func (m *MsgReplayPlayerInfo) String() string { return proto.CompactTextString(m) }
func (*MsgReplayPlayerInfo) ProtoMessage()    {}
func (*MsgReplayPlayerInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_9975814582a53f74, []int{0}
}
func (m *MsgReplayPlayerInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MsgReplayPlayerInfo.Unmarshal(m, b)
}
func (m *MsgReplayPlayerInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MsgReplayPlayerInfo.Marshal(b, m, deterministic)
}
func (m *MsgReplayPlayerInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgReplayPlayerInfo.Merge(m, src)
}
func (m *MsgReplayPlayerInfo) XXX_Size() int {
	return xxx_messageInfo_MsgReplayPlayerInfo.Size(m)
}
func (m *MsgReplayPlayerInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgReplayPlayerInfo.DiscardUnknown(m)
}

var xxx_messageInfo_MsgReplayPlayerInfo proto.InternalMessageInfo

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

func (m *MsgReplayPlayerInfo) GetGender() uint32 {
	if m != nil && m.Gender != nil {
		return *m.Gender
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
	ChairID              *int32   `protobuf:"varint,1,req,name=chairID" json:"chairID,omitempty"`
	Score                *int32   `protobuf:"varint,2,req,name=score" json:"score,omitempty"`
	WinType              *int32   `protobuf:"varint,3,req,name=winType" json:"winType,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MsgReplayPlayerScoreSummary) Reset()         { *m = MsgReplayPlayerScoreSummary{} }
func (m *MsgReplayPlayerScoreSummary) String() string { return proto.CompactTextString(m) }
func (*MsgReplayPlayerScoreSummary) ProtoMessage()    {}
func (*MsgReplayPlayerScoreSummary) Descriptor() ([]byte, []int) {
	return fileDescriptor_9975814582a53f74, []int{1}
}
func (m *MsgReplayPlayerScoreSummary) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MsgReplayPlayerScoreSummary.Unmarshal(m, b)
}
func (m *MsgReplayPlayerScoreSummary) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MsgReplayPlayerScoreSummary.Marshal(b, m, deterministic)
}
func (m *MsgReplayPlayerScoreSummary) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgReplayPlayerScoreSummary.Merge(m, src)
}
func (m *MsgReplayPlayerScoreSummary) XXX_Size() int {
	return xxx_messageInfo_MsgReplayPlayerScoreSummary.Size(m)
}
func (m *MsgReplayPlayerScoreSummary) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgReplayPlayerScoreSummary.DiscardUnknown(m)
}

var xxx_messageInfo_MsgReplayPlayerScoreSummary proto.InternalMessageInfo

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
	RecordUUID           *string                        `protobuf:"bytes,1,req,name=recordUUID" json:"recordUUID,omitempty"`
	PlayerScores         []*MsgReplayPlayerScoreSummary `protobuf:"bytes,2,rep,name=playerScores" json:"playerScores,omitempty"`
	EndTime              *uint32                        `protobuf:"varint,3,req,name=endTime" json:"endTime,omitempty"`
	ShareAbleID          *string                        `protobuf:"bytes,4,opt,name=shareAbleID" json:"shareAbleID,omitempty"`
	StartTime            *uint32                        `protobuf:"varint,5,opt,name=startTime" json:"startTime,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                       `json:"-"`
	XXX_unrecognized     []byte                         `json:"-"`
	XXX_sizecache        int32                          `json:"-"`
}

func (m *MsgReplayRecordSummary) Reset()         { *m = MsgReplayRecordSummary{} }
func (m *MsgReplayRecordSummary) String() string { return proto.CompactTextString(m) }
func (*MsgReplayRecordSummary) ProtoMessage()    {}
func (*MsgReplayRecordSummary) Descriptor() ([]byte, []int) {
	return fileDescriptor_9975814582a53f74, []int{2}
}
func (m *MsgReplayRecordSummary) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MsgReplayRecordSummary.Unmarshal(m, b)
}
func (m *MsgReplayRecordSummary) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MsgReplayRecordSummary.Marshal(b, m, deterministic)
}
func (m *MsgReplayRecordSummary) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgReplayRecordSummary.Merge(m, src)
}
func (m *MsgReplayRecordSummary) XXX_Size() int {
	return xxx_messageInfo_MsgReplayRecordSummary.Size(m)
}
func (m *MsgReplayRecordSummary) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgReplayRecordSummary.DiscardUnknown(m)
}

var xxx_messageInfo_MsgReplayRecordSummary proto.InternalMessageInfo

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
	RecordRoomType       *int32                    `protobuf:"varint,1,req,name=recordRoomType" json:"recordRoomType,omitempty"`
	StartTime            *uint32                   `protobuf:"varint,2,req,name=startTime" json:"startTime,omitempty"`
	EndTime              *uint32                   `protobuf:"varint,3,req,name=endTime" json:"endTime,omitempty"`
	RoomNumber           *string                   `protobuf:"bytes,4,req,name=roomNumber" json:"roomNumber,omitempty"`
	Players              []*MsgReplayPlayerInfo    `protobuf:"bytes,5,rep,name=players" json:"players,omitempty"`
	Records              []*MsgReplayRecordSummary `protobuf:"bytes,6,rep,name=records" json:"records,omitempty"`
	OwnerUserID          *string                   `protobuf:"bytes,7,opt,name=ownerUserID" json:"ownerUserID,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *MsgReplayRoom) Reset()         { *m = MsgReplayRoom{} }
func (m *MsgReplayRoom) String() string { return proto.CompactTextString(m) }
func (*MsgReplayRoom) ProtoMessage()    {}
func (*MsgReplayRoom) Descriptor() ([]byte, []int) {
	return fileDescriptor_9975814582a53f74, []int{3}
}
func (m *MsgReplayRoom) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MsgReplayRoom.Unmarshal(m, b)
}
func (m *MsgReplayRoom) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MsgReplayRoom.Marshal(b, m, deterministic)
}
func (m *MsgReplayRoom) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgReplayRoom.Merge(m, src)
}
func (m *MsgReplayRoom) XXX_Size() int {
	return xxx_messageInfo_MsgReplayRoom.Size(m)
}
func (m *MsgReplayRoom) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgReplayRoom.DiscardUnknown(m)
}

var xxx_messageInfo_MsgReplayRoom proto.InternalMessageInfo

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

func init() { proto.RegisterFile("game_pokerface_replay.proto", fileDescriptor_9975814582a53f74) }

var fileDescriptor_9975814582a53f74 = []byte{
	// 379 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xcd, 0xee, 0x9a, 0x40,
	0x14, 0xc5, 0x33, 0x28, 0x52, 0xae, 0x5f, 0xed, 0x98, 0x98, 0x49, 0x4c, 0x1a, 0xca, 0xa2, 0x61,
	0x65, 0x13, 0xd7, 0xdd, 0x34, 0x71, 0xc3, 0xa2, 0xa6, 0x51, 0x59, 0x9b, 0x11, 0xae, 0x4a, 0x04,
	0x86, 0xcc, 0x60, 0x0d, 0x4f, 0xd1, 0x87, 0xe8, 0x53, 0xf5, 0x6d, 0x1a, 0x2e, 0x29, 0xd6, 0x7e,
	0xe4, 0xbf, 0x9c, 0x19, 0xee, 0x39, 0xe7, 0x77, 0xb8, 0xb0, 0x38, 0xcb, 0x1c, 0x0f, 0xa5, 0xba,
	0xa2, 0x3e, 0xc9, 0x18, 0x0f, 0x1a, 0xcb, 0x4c, 0xd6, 0xcb, 0x52, 0xab, 0x4a, 0x71, 0xb7, 0xbb,
	0xf7, 0xbf, 0x31, 0x98, 0x7d, 0x36, 0xe7, 0x2d, 0x3d, 0x7f, 0xc9, 0x64, 0x8d, 0x3a, 0x2c, 0x4e,
	0x8a, 0x4f, 0x60, 0x70, 0x33, 0xa8, 0xc3, 0xb5, 0x60, 0x9e, 0x15, 0xb8, 0x7c, 0x04, 0xfd, 0x22,
	0x8d, 0xaf, 0xc2, 0xf2, 0x58, 0xe0, 0xf2, 0x29, 0x38, 0xf1, 0x45, 0xa6, 0xcd, 0x73, 0xcf, 0xb3,
	0x02, 0x9b, 0x73, 0x80, 0x4a, 0x55, 0x32, 0xdb, 0xc5, 0x4a, 0xa3, 0xe8, 0x7b, 0x2c, 0xb0, 0x1b,
	0x89, 0x33, 0x16, 0x09, 0x6a, 0x61, 0x7b, 0x2c, 0x18, 0xf3, 0x19, 0x0c, 0x2f, 0x28, 0x93, 0x30,
	0x56, 0x45, 0xb4, 0x0d, 0xc5, 0x80, 0x94, 0x5e, 0xc3, 0x2b, 0xf9, 0x55, 0x56, 0xb2, 0x91, 0x72,
	0x9a, 0x31, 0x7f, 0x03, 0x8b, 0x3f, 0x02, 0x91, 0xe8, 0xee, 0x96, 0xe7, 0x52, 0xd7, 0xbf, 0x5b,
	0x33, 0xb2, 0x1e, 0x83, 0x6d, 0xc8, 0xd5, 0xa2, 0xe3, 0x14, 0x9c, 0x7b, 0x5a, 0xec, 0xeb, 0x12,
	0xdb, 0x68, 0xfe, 0x77, 0x06, 0xf3, 0x4e, 0x70, 0x8b, 0xb1, 0xd2, 0xc9, 0x2f, 0x2d, 0x0e, 0xa0,
	0xe9, 0x22, 0x8a, 0x3a, 0xd0, 0x8f, 0x30, 0x2a, 0x1f, 0xae, 0x46, 0x58, 0x5e, 0x2f, 0x18, 0xae,
	0xde, 0x2f, 0xbb, 0xca, 0x96, 0x2f, 0xa4, 0xc3, 0x22, 0xd9, 0xa7, 0x79, 0xeb, 0x4e, 0xd0, 0xe6,
	0x22, 0x35, 0x7e, 0x3a, 0x66, 0x18, 0xae, 0xa9, 0x19, 0x97, 0xbf, 0x01, 0xd7, 0x54, 0x52, 0x57,
	0xf4, 0x1d, 0x95, 0xe3, 0xff, 0x60, 0x30, 0x7e, 0xa4, 0x54, 0x2a, 0xe7, 0x73, 0x98, 0xb4, 0xe1,
	0x9a, 0x13, 0xf1, 0xb4, 0xbc, 0x4f, 0xc3, 0x16, 0x99, 0xfc, 0xe5, 0xda, 0x80, 0x29, 0x95, 0x6f,
	0x6e, 0xf9, 0x11, 0xb5, 0xe8, 0x13, 0xd8, 0x07, 0x70, 0x5a, 0x30, 0x23, 0x6c, 0x62, 0x7a, 0xfb,
	0x7f, 0x26, 0x5a, 0x81, 0x15, 0x38, 0x6d, 0x00, 0x23, 0x06, 0x34, 0xf0, 0xee, 0x5f, 0x03, 0xcf,
	0x8d, 0xce, 0x60, 0xa8, 0xee, 0x05, 0xea, 0xa8, 0xdd, 0x9d, 0xe6, 0x8f, 0xba, 0x3f, 0x03, 0x00,
	0x00, 0xff, 0xff, 0xa9, 0x25, 0xfd, 0xb8, 0x8c, 0x02, 0x00, 0x00,
}
