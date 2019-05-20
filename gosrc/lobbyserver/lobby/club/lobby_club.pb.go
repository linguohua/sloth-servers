// Code generated by protoc-gen-go. DO NOT EDIT.
// source: lobby_club.proto

/*
Package club is a generated protocol buffer package.

It is generated from these files:
	lobby_club.proto

It has these top-level messages:
	MsgClubReply
	MsgClubDisplayInfo
	MsgClubMemberInfo
	MsgClubBaseInfo
	MsgCubOperGenericReply
	MsgClubInfo
	MsgClubLoadMyClubsReply
	MsgClubLoadUpdateReply
	MsgClubLoadMembersReply
	MsgCreateClubReply
	MsgClubEvent
	MsgClubLoadEventsReply
	MsgClubRoomInfo
	MsgClubLoadRoomsReply
	MsgClubFundEvent
	MsgClubLoadFundEventsReply
	MsgClubLoadReplayRoomsReply
*/
package club

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type MsgClubReply struct {
	ReplyCode        *int32 `protobuf:"varint,1,req,name=replyCode" json:"replyCode,omitempty"`
	Content          []byte `protobuf:"bytes,2,opt,name=content" json:"content,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *MsgClubReply) Reset()                    { *m = MsgClubReply{} }
func (m *MsgClubReply) String() string            { return proto.CompactTextString(m) }
func (*MsgClubReply) ProtoMessage()               {}
func (*MsgClubReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *MsgClubReply) GetReplyCode() int32 {
	if m != nil && m.ReplyCode != nil {
		return *m.ReplyCode
	}
	return 0
}

func (m *MsgClubReply) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

type MsgClubDisplayInfo struct {
	Nick             *string `protobuf:"bytes,1,req,name=nick" json:"nick,omitempty"`
	Sex              *uint32 `protobuf:"varint,2,opt,name=sex" json:"sex,omitempty"`
	HeadIconURL      *string `protobuf:"bytes,3,opt,name=headIconURL" json:"headIconURL,omitempty"`
	AvatarID         *int32  `protobuf:"varint,4,opt,name=avatarID" json:"avatarID,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *MsgClubDisplayInfo) Reset()                    { *m = MsgClubDisplayInfo{} }
func (m *MsgClubDisplayInfo) String() string            { return proto.CompactTextString(m) }
func (*MsgClubDisplayInfo) ProtoMessage()               {}
func (*MsgClubDisplayInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *MsgClubDisplayInfo) GetNick() string {
	if m != nil && m.Nick != nil {
		return *m.Nick
	}
	return ""
}

func (m *MsgClubDisplayInfo) GetSex() uint32 {
	if m != nil && m.Sex != nil {
		return *m.Sex
	}
	return 0
}

func (m *MsgClubDisplayInfo) GetHeadIconURL() string {
	if m != nil && m.HeadIconURL != nil {
		return *m.HeadIconURL
	}
	return ""
}

func (m *MsgClubDisplayInfo) GetAvatarID() int32 {
	if m != nil && m.AvatarID != nil {
		return *m.AvatarID
	}
	return 0
}

// 俱乐部成员信息
type MsgClubMemberInfo struct {
	UserID           *string             `protobuf:"bytes,1,req,name=userID" json:"userID,omitempty"`
	DisplayInfo      *MsgClubDisplayInfo `protobuf:"bytes,2,opt,name=displayInfo" json:"displayInfo,omitempty"`
	Online           *bool               `protobuf:"varint,3,opt,name=online" json:"online,omitempty"`
	XXX_unrecognized []byte              `json:"-"`
}

func (m *MsgClubMemberInfo) Reset()                    { *m = MsgClubMemberInfo{} }
func (m *MsgClubMemberInfo) String() string            { return proto.CompactTextString(m) }
func (*MsgClubMemberInfo) ProtoMessage()               {}
func (*MsgClubMemberInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *MsgClubMemberInfo) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *MsgClubMemberInfo) GetDisplayInfo() *MsgClubDisplayInfo {
	if m != nil {
		return m.DisplayInfo
	}
	return nil
}

func (m *MsgClubMemberInfo) GetOnline() bool {
	if m != nil && m.Online != nil {
		return *m.Online
	}
	return false
}

// 俱乐部的基本信息
type MsgClubBaseInfo struct {
	ClubNumber       *string `protobuf:"bytes,1,req,name=clubNumber" json:"clubNumber,omitempty"`
	ClubName         *string `protobuf:"bytes,2,opt,name=clubName" json:"clubName,omitempty"`
	ClubID           *string `protobuf:"bytes,3,opt,name=clubID" json:"clubID,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *MsgClubBaseInfo) Reset()                    { *m = MsgClubBaseInfo{} }
func (m *MsgClubBaseInfo) String() string            { return proto.CompactTextString(m) }
func (*MsgClubBaseInfo) ProtoMessage()               {}
func (*MsgClubBaseInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *MsgClubBaseInfo) GetClubNumber() string {
	if m != nil && m.ClubNumber != nil {
		return *m.ClubNumber
	}
	return ""
}

func (m *MsgClubBaseInfo) GetClubName() string {
	if m != nil && m.ClubName != nil {
		return *m.ClubName
	}
	return ""
}

func (m *MsgClubBaseInfo) GetClubID() string {
	if m != nil && m.ClubID != nil {
		return *m.ClubID
	}
	return ""
}

// 俱乐部操作通用回复，免得定义太多消息体
type MsgCubOperGenericReply struct {
	ErrorCode        *int32  `protobuf:"varint,1,req,name=errorCode" json:"errorCode,omitempty"`
	Extra            *string `protobuf:"bytes,2,opt,name=extra" json:"extra,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *MsgCubOperGenericReply) Reset()                    { *m = MsgCubOperGenericReply{} }
func (m *MsgCubOperGenericReply) String() string            { return proto.CompactTextString(m) }
func (*MsgCubOperGenericReply) ProtoMessage()               {}
func (*MsgCubOperGenericReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *MsgCubOperGenericReply) GetErrorCode() int32 {
	if m != nil && m.ErrorCode != nil {
		return *m.ErrorCode
	}
	return 0
}

func (m *MsgCubOperGenericReply) GetExtra() string {
	if m != nil && m.Extra != nil {
		return *m.Extra
	}
	return ""
}

// 俱乐部信息
type MsgClubInfo struct {
	BaseInfo         *MsgClubBaseInfo `protobuf:"bytes,1,opt,name=baseInfo" json:"baseInfo,omitempty"`
	CreatorUserID    *string          `protobuf:"bytes,2,opt,name=creatorUserID" json:"creatorUserID,omitempty"`
	ClubLevel        *int32           `protobuf:"varint,3,opt,name=clubLevel" json:"clubLevel,omitempty"`
	Points           *int32           `protobuf:"varint,4,opt,name=points" json:"points,omitempty"`
	Wanka            *int32           `protobuf:"varint,5,opt,name=wanka" json:"wanka,omitempty"`
	Candy            *int32           `protobuf:"varint,6,opt,name=candy" json:"candy,omitempty"`
	MaxMember        *int32           `protobuf:"varint,7,opt,name=maxMember" json:"maxMember,omitempty"`
	JoinForbit       *bool            `protobuf:"varint,8,opt,name=joinForbit" json:"joinForbit,omitempty"`
	HasUnReadEvents  *bool            `protobuf:"varint,9,opt,name=hasUnReadEvents" json:"hasUnReadEvents,omitempty"`
	CreateRoomOption *int32           `protobuf:"varint,10,opt,name=createRoomOption" json:"createRoomOption,omitempty"`
	PayRoomOption    *int32           `protobuf:"varint,11,opt,name=payRoomOption" json:"payRoomOption,omitempty"`
	CreateTime       *int32           `protobuf:"varint,12,opt,name=createTime" json:"createTime,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *MsgClubInfo) Reset()                    { *m = MsgClubInfo{} }
func (m *MsgClubInfo) String() string            { return proto.CompactTextString(m) }
func (*MsgClubInfo) ProtoMessage()               {}
func (*MsgClubInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *MsgClubInfo) GetBaseInfo() *MsgClubBaseInfo {
	if m != nil {
		return m.BaseInfo
	}
	return nil
}

func (m *MsgClubInfo) GetCreatorUserID() string {
	if m != nil && m.CreatorUserID != nil {
		return *m.CreatorUserID
	}
	return ""
}

func (m *MsgClubInfo) GetClubLevel() int32 {
	if m != nil && m.ClubLevel != nil {
		return *m.ClubLevel
	}
	return 0
}

func (m *MsgClubInfo) GetPoints() int32 {
	if m != nil && m.Points != nil {
		return *m.Points
	}
	return 0
}

func (m *MsgClubInfo) GetWanka() int32 {
	if m != nil && m.Wanka != nil {
		return *m.Wanka
	}
	return 0
}

func (m *MsgClubInfo) GetCandy() int32 {
	if m != nil && m.Candy != nil {
		return *m.Candy
	}
	return 0
}

func (m *MsgClubInfo) GetMaxMember() int32 {
	if m != nil && m.MaxMember != nil {
		return *m.MaxMember
	}
	return 0
}

func (m *MsgClubInfo) GetJoinForbit() bool {
	if m != nil && m.JoinForbit != nil {
		return *m.JoinForbit
	}
	return false
}

func (m *MsgClubInfo) GetHasUnReadEvents() bool {
	if m != nil && m.HasUnReadEvents != nil {
		return *m.HasUnReadEvents
	}
	return false
}

func (m *MsgClubInfo) GetCreateRoomOption() int32 {
	if m != nil && m.CreateRoomOption != nil {
		return *m.CreateRoomOption
	}
	return 0
}

func (m *MsgClubInfo) GetPayRoomOption() int32 {
	if m != nil && m.PayRoomOption != nil {
		return *m.PayRoomOption
	}
	return 0
}

func (m *MsgClubInfo) GetCreateTime() int32 {
	if m != nil && m.CreateTime != nil {
		return *m.CreateTime
	}
	return 0
}

// 加载自己的俱乐部
type MsgClubLoadMyClubsReply struct {
	Clubs            []*MsgClubInfo `protobuf:"bytes,1,rep,name=clubs" json:"clubs,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *MsgClubLoadMyClubsReply) Reset()                    { *m = MsgClubLoadMyClubsReply{} }
func (m *MsgClubLoadMyClubsReply) String() string            { return proto.CompactTextString(m) }
func (*MsgClubLoadMyClubsReply) ProtoMessage()               {}
func (*MsgClubLoadMyClubsReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *MsgClubLoadMyClubsReply) GetClubs() []*MsgClubInfo {
	if m != nil {
		return m.Clubs
	}
	return nil
}

// 加载俱乐部更新回复
type MsgClubLoadUpdateReply struct {
	ClubsUpdated     []*MsgClubInfo `protobuf:"bytes,1,rep,name=clubsUpdated" json:"clubsUpdated,omitempty"`
	ClubIDsRemoved   []string       `protobuf:"bytes,2,rep,name=clubIDsRemoved" json:"clubIDsRemoved,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *MsgClubLoadUpdateReply) Reset()                    { *m = MsgClubLoadUpdateReply{} }
func (m *MsgClubLoadUpdateReply) String() string            { return proto.CompactTextString(m) }
func (*MsgClubLoadUpdateReply) ProtoMessage()               {}
func (*MsgClubLoadUpdateReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *MsgClubLoadUpdateReply) GetClubsUpdated() []*MsgClubInfo {
	if m != nil {
		return m.ClubsUpdated
	}
	return nil
}

func (m *MsgClubLoadUpdateReply) GetClubIDsRemoved() []string {
	if m != nil {
		return m.ClubIDsRemoved
	}
	return nil
}

// 服务器回复请求成员列表
type MsgClubLoadMembersReply struct {
	Members          []*MsgClubMemberInfo `protobuf:"bytes,1,rep,name=members" json:"members,omitempty"`
	Cursor           *int32               `protobuf:"varint,2,opt,name=cursor" json:"cursor,omitempty"`
	XXX_unrecognized []byte               `json:"-"`
}

func (m *MsgClubLoadMembersReply) Reset()                    { *m = MsgClubLoadMembersReply{} }
func (m *MsgClubLoadMembersReply) String() string            { return proto.CompactTextString(m) }
func (*MsgClubLoadMembersReply) ProtoMessage()               {}
func (*MsgClubLoadMembersReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *MsgClubLoadMembersReply) GetMembers() []*MsgClubMemberInfo {
	if m != nil {
		return m.Members
	}
	return nil
}

func (m *MsgClubLoadMembersReply) GetCursor() int32 {
	if m != nil && m.Cursor != nil {
		return *m.Cursor
	}
	return 0
}

// 创建俱乐部服务器给客户端的回复
type MsgCreateClubReply struct {
	ClubInfo         *MsgClubInfo `protobuf:"bytes,1,opt,name=clubInfo" json:"clubInfo,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *MsgCreateClubReply) Reset()                    { *m = MsgCreateClubReply{} }
func (m *MsgCreateClubReply) String() string            { return proto.CompactTextString(m) }
func (*MsgCreateClubReply) ProtoMessage()               {}
func (*MsgCreateClubReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *MsgCreateClubReply) GetClubInfo() *MsgClubInfo {
	if m != nil {
		return m.ClubInfo
	}
	return nil
}

// 俱乐部事件，注意和MsgClubNotification不同，事件是系统生成，不可设置
type MsgClubEvent struct {
	EvtType          *int32              `protobuf:"varint,1,req,name=evtType" json:"evtType,omitempty"`
	Id               *uint32             `protobuf:"varint,2,req,name=Id" json:"Id,omitempty"`
	GeneratedTime    *uint32             `protobuf:"varint,3,req,name=generatedTime" json:"generatedTime,omitempty"`
	To               *string             `protobuf:"bytes,4,opt,name=to" json:"to,omitempty"`
	Content          []byte              `protobuf:"bytes,5,opt,name=content" json:"content,omitempty"`
	Unread           *bool               `protobuf:"varint,6,opt,name=unread" json:"unread,omitempty"`
	NeedHandle       *bool               `protobuf:"varint,7,opt,name=needHandle" json:"needHandle,omitempty"`
	UserID1          *string             `protobuf:"bytes,8,opt,name=userID1" json:"userID1,omitempty"`
	DisplayInfo1     *MsgClubDisplayInfo `protobuf:"bytes,9,opt,name=displayInfo1" json:"displayInfo1,omitempty"`
	ApprovalResult   *int32              `protobuf:"varint,10,opt,name=approvalResult" json:"approvalResult,omitempty"`
	XXX_unrecognized []byte              `json:"-"`
}

func (m *MsgClubEvent) Reset()                    { *m = MsgClubEvent{} }
func (m *MsgClubEvent) String() string            { return proto.CompactTextString(m) }
func (*MsgClubEvent) ProtoMessage()               {}
func (*MsgClubEvent) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *MsgClubEvent) GetEvtType() int32 {
	if m != nil && m.EvtType != nil {
		return *m.EvtType
	}
	return 0
}

func (m *MsgClubEvent) GetId() uint32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *MsgClubEvent) GetGeneratedTime() uint32 {
	if m != nil && m.GeneratedTime != nil {
		return *m.GeneratedTime
	}
	return 0
}

func (m *MsgClubEvent) GetTo() string {
	if m != nil && m.To != nil {
		return *m.To
	}
	return ""
}

func (m *MsgClubEvent) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

func (m *MsgClubEvent) GetUnread() bool {
	if m != nil && m.Unread != nil {
		return *m.Unread
	}
	return false
}

func (m *MsgClubEvent) GetNeedHandle() bool {
	if m != nil && m.NeedHandle != nil {
		return *m.NeedHandle
	}
	return false
}

func (m *MsgClubEvent) GetUserID1() string {
	if m != nil && m.UserID1 != nil {
		return *m.UserID1
	}
	return ""
}

func (m *MsgClubEvent) GetDisplayInfo1() *MsgClubDisplayInfo {
	if m != nil {
		return m.DisplayInfo1
	}
	return nil
}

func (m *MsgClubEvent) GetApprovalResult() int32 {
	if m != nil && m.ApprovalResult != nil {
		return *m.ApprovalResult
	}
	return 0
}

// 服务器回复请求事件列表
type MsgClubLoadEventsReply struct {
	Events           []*MsgClubEvent `protobuf:"bytes,1,rep,name=events" json:"events,omitempty"`
	Cursor           *int32          `protobuf:"varint,2,opt,name=cursor" json:"cursor,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *MsgClubLoadEventsReply) Reset()                    { *m = MsgClubLoadEventsReply{} }
func (m *MsgClubLoadEventsReply) String() string            { return proto.CompactTextString(m) }
func (*MsgClubLoadEventsReply) ProtoMessage()               {}
func (*MsgClubLoadEventsReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *MsgClubLoadEventsReply) GetEvents() []*MsgClubEvent {
	if m != nil {
		return m.Events
	}
	return nil
}

func (m *MsgClubLoadEventsReply) GetCursor() int32 {
	if m != nil && m.Cursor != nil {
		return *m.Cursor
	}
	return 0
}

// 俱乐部房间信息
type MsgClubRoomInfo struct {
	RoomType         *int32  `protobuf:"varint,1,req,name=roomType" json:"roomType,omitempty"`
	RoomRuleJSON     *string `protobuf:"bytes,2,opt,name=roomRuleJSON" json:"roomRuleJSON,omitempty"`
	PlayerNumber     *int32  `protobuf:"varint,3,opt,name=playerNumber" json:"playerNumber,omitempty"`
	RoomState        *int32  `protobuf:"varint,4,opt,name=roomState" json:"roomState,omitempty"`
	RoomNumber       *string `protobuf:"bytes,5,opt,name=roomNumber" json:"roomNumber,omitempty"`
	RoomUUID         *string `protobuf:"bytes,6,opt,name=roomUUID" json:"roomUUID,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *MsgClubRoomInfo) Reset()                    { *m = MsgClubRoomInfo{} }
func (m *MsgClubRoomInfo) String() string            { return proto.CompactTextString(m) }
func (*MsgClubRoomInfo) ProtoMessage()               {}
func (*MsgClubRoomInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *MsgClubRoomInfo) GetRoomType() int32 {
	if m != nil && m.RoomType != nil {
		return *m.RoomType
	}
	return 0
}

func (m *MsgClubRoomInfo) GetRoomRuleJSON() string {
	if m != nil && m.RoomRuleJSON != nil {
		return *m.RoomRuleJSON
	}
	return ""
}

func (m *MsgClubRoomInfo) GetPlayerNumber() int32 {
	if m != nil && m.PlayerNumber != nil {
		return *m.PlayerNumber
	}
	return 0
}

func (m *MsgClubRoomInfo) GetRoomState() int32 {
	if m != nil && m.RoomState != nil {
		return *m.RoomState
	}
	return 0
}

func (m *MsgClubRoomInfo) GetRoomNumber() string {
	if m != nil && m.RoomNumber != nil {
		return *m.RoomNumber
	}
	return ""
}

func (m *MsgClubRoomInfo) GetRoomUUID() string {
	if m != nil && m.RoomUUID != nil {
		return *m.RoomUUID
	}
	return ""
}

// 服务器回复请求俱乐部房间列表
type MsgClubLoadRoomsReply struct {
	Rooms            []*MsgClubRoomInfo `protobuf:"bytes,1,rep,name=rooms" json:"rooms,omitempty"`
	Cursor           *int32             `protobuf:"varint,2,opt,name=cursor" json:"cursor,omitempty"`
	TotalRoomCount   *int32             `protobuf:"varint,3,opt,name=totalRoomCount" json:"totalRoomCount,omitempty"`
	XXX_unrecognized []byte             `json:"-"`
}

func (m *MsgClubLoadRoomsReply) Reset()                    { *m = MsgClubLoadRoomsReply{} }
func (m *MsgClubLoadRoomsReply) String() string            { return proto.CompactTextString(m) }
func (*MsgClubLoadRoomsReply) ProtoMessage()               {}
func (*MsgClubLoadRoomsReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *MsgClubLoadRoomsReply) GetRooms() []*MsgClubRoomInfo {
	if m != nil {
		return m.Rooms
	}
	return nil
}

func (m *MsgClubLoadRoomsReply) GetCursor() int32 {
	if m != nil && m.Cursor != nil {
		return *m.Cursor
	}
	return 0
}

func (m *MsgClubLoadRoomsReply) GetTotalRoomCount() int32 {
	if m != nil && m.TotalRoomCount != nil {
		return *m.TotalRoomCount
	}
	return 0
}

// 俱乐部基金事件
type MsgClubFundEvent struct {
	EvtType          *int32  `protobuf:"varint,1,req,name=evtType" json:"evtType,omitempty"`
	GeneratedTime    *uint32 `protobuf:"varint,2,req,name=generatedTime" json:"generatedTime,omitempty"`
	UserID           *string `protobuf:"bytes,3,req,name=userID" json:"userID,omitempty"`
	Amount           *int32  `protobuf:"varint,4,req,name=amount" json:"amount,omitempty"`
	Total            *int32  `protobuf:"varint,5,req,name=total" json:"total,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *MsgClubFundEvent) Reset()                    { *m = MsgClubFundEvent{} }
func (m *MsgClubFundEvent) String() string            { return proto.CompactTextString(m) }
func (*MsgClubFundEvent) ProtoMessage()               {}
func (*MsgClubFundEvent) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *MsgClubFundEvent) GetEvtType() int32 {
	if m != nil && m.EvtType != nil {
		return *m.EvtType
	}
	return 0
}

func (m *MsgClubFundEvent) GetGeneratedTime() uint32 {
	if m != nil && m.GeneratedTime != nil {
		return *m.GeneratedTime
	}
	return 0
}

func (m *MsgClubFundEvent) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *MsgClubFundEvent) GetAmount() int32 {
	if m != nil && m.Amount != nil {
		return *m.Amount
	}
	return 0
}

func (m *MsgClubFundEvent) GetTotal() int32 {
	if m != nil && m.Total != nil {
		return *m.Total
	}
	return 0
}

// 加载俱乐部基金事件的回复
type MsgClubLoadFundEventsReply struct {
	Events           []*MsgClubFundEvent `protobuf:"bytes,1,rep,name=events" json:"events,omitempty"`
	Cursor           *int32              `protobuf:"varint,2,opt,name=cursor" json:"cursor,omitempty"`
	XXX_unrecognized []byte              `json:"-"`
}

func (m *MsgClubLoadFundEventsReply) Reset()                    { *m = MsgClubLoadFundEventsReply{} }
func (m *MsgClubLoadFundEventsReply) String() string            { return proto.CompactTextString(m) }
func (*MsgClubLoadFundEventsReply) ProtoMessage()               {}
func (*MsgClubLoadFundEventsReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

func (m *MsgClubLoadFundEventsReply) GetEvents() []*MsgClubFundEvent {
	if m != nil {
		return m.Events
	}
	return nil
}

func (m *MsgClubLoadFundEventsReply) GetCursor() int32 {
	if m != nil && m.Cursor != nil {
		return *m.Cursor
	}
	return 0
}

// 加载回播房间的回复
type MsgClubLoadReplayRoomsReply struct {
	GZipBytes        []byte `protobuf:"bytes,1,opt,name=gZipBytes" json:"gZipBytes,omitempty"`
	Cursor           *int32 `protobuf:"varint,2,opt,name=cursor" json:"cursor,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *MsgClubLoadReplayRoomsReply) Reset()                    { *m = MsgClubLoadReplayRoomsReply{} }
func (m *MsgClubLoadReplayRoomsReply) String() string            { return proto.CompactTextString(m) }
func (*MsgClubLoadReplayRoomsReply) ProtoMessage()               {}
func (*MsgClubLoadReplayRoomsReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{16} }

func (m *MsgClubLoadReplayRoomsReply) GetGZipBytes() []byte {
	if m != nil {
		return m.GZipBytes
	}
	return nil
}

func (m *MsgClubLoadReplayRoomsReply) GetCursor() int32 {
	if m != nil && m.Cursor != nil {
		return *m.Cursor
	}
	return 0
}

func init() {
	proto.RegisterType((*MsgClubReply)(nil), "club.MsgClubReply")
	proto.RegisterType((*MsgClubDisplayInfo)(nil), "club.MsgClubDisplayInfo")
	proto.RegisterType((*MsgClubMemberInfo)(nil), "club.MsgClubMemberInfo")
	proto.RegisterType((*MsgClubBaseInfo)(nil), "club.MsgClubBaseInfo")
	proto.RegisterType((*MsgCubOperGenericReply)(nil), "club.MsgCubOperGenericReply")
	proto.RegisterType((*MsgClubInfo)(nil), "club.MsgClubInfo")
	proto.RegisterType((*MsgClubLoadMyClubsReply)(nil), "club.MsgClubLoadMyClubsReply")
	proto.RegisterType((*MsgClubLoadUpdateReply)(nil), "club.MsgClubLoadUpdateReply")
	proto.RegisterType((*MsgClubLoadMembersReply)(nil), "club.MsgClubLoadMembersReply")
	proto.RegisterType((*MsgCreateClubReply)(nil), "club.MsgCreateClubReply")
	proto.RegisterType((*MsgClubEvent)(nil), "club.MsgClubEvent")
	proto.RegisterType((*MsgClubLoadEventsReply)(nil), "club.MsgClubLoadEventsReply")
	proto.RegisterType((*MsgClubRoomInfo)(nil), "club.MsgClubRoomInfo")
	proto.RegisterType((*MsgClubLoadRoomsReply)(nil), "club.MsgClubLoadRoomsReply")
	proto.RegisterType((*MsgClubFundEvent)(nil), "club.MsgClubFundEvent")
	proto.RegisterType((*MsgClubLoadFundEventsReply)(nil), "club.MsgClubLoadFundEventsReply")
	proto.RegisterType((*MsgClubLoadReplayRoomsReply)(nil), "club.MsgClubLoadReplayRoomsReply")
}

func init() { proto.RegisterFile("lobby_club.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 850 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x54, 0x4f, 0x6f, 0xdb, 0xc6,
	0x13, 0x85, 0x24, 0x53, 0x16, 0x47, 0x94, 0x2c, 0xf1, 0x17, 0xdb, 0xc4, 0xaf, 0x17, 0x81, 0x2d,
	0x1a, 0x5d, 0x6a, 0x20, 0xbe, 0xb5, 0xbd, 0x14, 0xb1, 0x9b, 0x54, 0x85, 0x1d, 0x03, 0xb2, 0x55,
	0xa0, 0xbd, 0x14, 0x4b, 0x71, 0xea, 0xb0, 0x21, 0x77, 0x89, 0xe5, 0x52, 0xb5, 0xbe, 0x41, 0xef,
	0xfd, 0x7a, 0xfd, 0x30, 0xc5, 0xcc, 0xae, 0x64, 0x52, 0x4d, 0x7a, 0x13, 0x47, 0xf3, 0xe7, 0xcd,
	0x9b, 0xf7, 0x16, 0x26, 0xb9, 0x4a, 0x92, 0xed, 0xaf, 0xeb, 0xbc, 0x4e, 0x2e, 0x4a, 0xad, 0x8c,
	0x0a, 0x8f, 0xe8, 0x77, 0x7c, 0x09, 0xc1, 0x6d, 0xf5, 0x78, 0x95, 0xd7, 0xc9, 0x12, 0xcb, 0x7c,
	0x1b, 0x4e, 0xc1, 0xd7, 0xf4, 0xe3, 0x4a, 0xa5, 0x18, 0x75, 0x66, 0xdd, 0xb9, 0x17, 0x9e, 0xc0,
	0xf1, 0x5a, 0x49, 0x83, 0xd2, 0x44, 0xdd, 0x59, 0x67, 0x1e, 0xc4, 0x3f, 0x41, 0xe8, 0x6a, 0xae,
	0xb3, 0xaa, 0xcc, 0xc5, 0x76, 0x21, 0x7f, 0x53, 0x61, 0x00, 0x47, 0x32, 0x5b, 0x7f, 0xe0, 0x22,
	0x3f, 0x1c, 0x42, 0xaf, 0xc2, 0x27, 0x2e, 0x18, 0x85, 0xff, 0x83, 0xe1, 0x7b, 0x14, 0xe9, 0x62,
	0xad, 0xe4, 0x6a, 0x79, 0x13, 0xf5, 0x66, 0x9d, 0xb9, 0x1f, 0x4e, 0x60, 0x20, 0x36, 0xc2, 0x08,
	0xbd, 0xb8, 0x8e, 0x8e, 0x66, 0x9d, 0xb9, 0x17, 0x27, 0x30, 0x75, 0x7d, 0x6f, 0xb1, 0x48, 0x50,
	0x73, 0xdb, 0x31, 0xf4, 0xeb, 0x0a, 0x29, 0xc9, 0x36, 0xfe, 0x0a, 0x86, 0xe9, 0xf3, 0x54, 0x1e,
	0x30, 0xbc, 0x8c, 0x2e, 0x78, 0xb1, 0x8f, 0xa0, 0x1a, 0x43, 0x5f, 0xc9, 0x3c, 0x93, 0xc8, 0x53,
	0x07, 0xf1, 0x5b, 0x38, 0x71, 0x59, 0xaf, 0x45, 0x85, 0x9c, 0x12, 0x02, 0x50, 0xf5, 0xbb, 0x9a,
	0x66, 0xba, 0x29, 0x13, 0x18, 0x70, 0x4c, 0x14, 0xc8, 0x23, 0x7c, 0x6a, 0x44, 0x91, 0xc5, 0xb5,
	0x85, 0x1f, 0x7f, 0x03, 0x67, 0xd4, 0xa8, 0x4e, 0xee, 0x4a, 0xd4, 0x6f, 0x51, 0xa2, 0xce, 0xd6,
	0x7b, 0x0a, 0x51, 0x6b, 0xa5, 0x1b, 0x14, 0x8e, 0xc0, 0xc3, 0x27, 0xa3, 0x85, 0xed, 0x15, 0xff,
	0xd5, 0x85, 0xa1, 0x43, 0xc1, 0x08, 0x5e, 0xc2, 0x20, 0x71, 0x68, 0xa2, 0x0e, 0x2f, 0x74, 0xda,
	0x5a, 0x68, 0x0f, 0xf5, 0x14, 0x46, 0x6b, 0x8d, 0xc2, 0x28, 0xbd, 0xb2, 0x9c, 0x58, 0x6c, 0x53,
	0xf0, 0x29, 0xfd, 0x06, 0x37, 0x98, 0x33, 0x3c, 0x8f, 0xe0, 0x96, 0x2a, 0x93, 0xa6, 0xb2, 0xdc,
	0x12, 0x82, 0x3f, 0x84, 0xfc, 0x20, 0x22, 0x6f, 0xf7, 0xb9, 0x16, 0x32, 0xdd, 0x46, 0x7d, 0xfe,
	0x9c, 0x82, 0x5f, 0x88, 0x27, 0xcb, 0x7a, 0x74, 0xcc, 0xa1, 0x10, 0xe0, 0x77, 0x95, 0xc9, 0x37,
	0x4a, 0x27, 0x99, 0x89, 0x06, 0x44, 0x5e, 0x78, 0x0e, 0x27, 0xef, 0x45, 0xb5, 0x92, 0x4b, 0x14,
	0xe9, 0xf7, 0x1b, 0xa4, 0xee, 0x3e, 0xff, 0x11, 0xc1, 0x84, 0x71, 0xe1, 0x52, 0xa9, 0xe2, 0xae,
	0x34, 0x99, 0x92, 0x11, 0x70, 0x9b, 0x53, 0x18, 0x95, 0x62, 0xdb, 0x08, 0x0f, 0x77, 0xdd, 0x6d,
	0xc1, 0x43, 0x56, 0x60, 0x14, 0xf0, 0xf9, 0xbf, 0x85, 0x73, 0xb7, 0xef, 0x8d, 0x12, 0xe9, 0xed,
	0x96, 0x7e, 0x55, 0x96, 0xd2, 0x19, 0x78, 0xb4, 0x60, 0x15, 0x75, 0x66, 0xbd, 0xf9, 0xf0, 0x72,
	0xda, 0x62, 0x87, 0x98, 0x89, 0x7f, 0xb6, 0xe7, 0x70, 0xc5, 0xab, 0x32, 0x25, 0x30, 0x5c, 0xfb,
	0x12, 0x02, 0xae, 0xb5, 0xb1, 0xf4, 0x93, 0x2d, 0xc2, 0x33, 0x18, 0xdb, 0x0b, 0x57, 0x4b, 0x2c,
	0xd4, 0x06, 0xd3, 0xa8, 0x3b, 0xeb, 0xcd, 0xfd, 0xf8, 0xbe, 0x8d, 0x8b, 0x49, 0x72, 0xb8, 0xe6,
	0x70, 0x5c, 0xd8, 0x6f, 0xd7, 0xf6, 0xbc, 0xd5, 0xb6, 0x2d, 0xe3, 0x75, 0xad, 0x2b, 0xa5, 0xf9,
	0x64, 0x5e, 0xfc, 0xb5, 0xf5, 0x10, 0x73, 0xf0, 0xec, 0xbe, 0xcf, 0xad, 0xec, 0x1a, 0x42, 0xf8,
	0xc8, 0xaa, 0x7f, 0x77, 0xf6, 0x9e, 0xe5, 0x23, 0x90, 0x41, 0x71, 0x63, 0x1e, 0xb6, 0xe5, 0x4e,
	0x6e, 0x00, 0xdd, 0x05, 0xa1, 0xef, 0xce, 0x47, 0x74, 0x80, 0x47, 0x52, 0x27, 0xed, 0xce, 0x64,
	0xf7, 0x38, 0x0c, 0xd0, 0x35, 0x8a, 0xb5, 0xe1, 0x37, 0x0d, 0x4e, 0xea, 0x08, 0xd8, 0x73, 0x52,
	0xa3, 0x48, 0x59, 0x1e, 0x03, 0xba, 0x96, 0x44, 0x4c, 0x7f, 0x10, 0x32, 0xcd, 0x91, 0xf5, 0x31,
	0xa0, 0x22, 0xeb, 0xcb, 0x57, 0x2c, 0x0e, 0x3f, 0xbc, 0x80, 0xa0, 0x61, 0xcc, 0x57, 0xac, 0x8c,
	0xff, 0x72, 0xe6, 0x19, 0x8c, 0x45, 0x59, 0x6a, 0xb5, 0x11, 0xf9, 0x12, 0xab, 0x3a, 0x37, 0x56,
	0x31, 0xf1, 0x4d, 0xeb, 0x92, 0x56, 0x66, 0x96, 0x9d, 0x18, 0xfa, 0x68, 0x55, 0x67, 0xc9, 0x0e,
	0x5b, 0xbd, 0x2d, 0x17, 0x87, 0x3c, 0xff, 0xd9, 0xd9, 0x1b, 0x9e, 0x44, 0xc8, 0x93, 0x27, 0x30,
	0xd0, 0x4a, 0x15, 0x0d, 0xc2, 0x5e, 0x40, 0x40, 0x91, 0x65, 0x9d, 0xe3, 0x8f, 0xf7, 0x77, 0xef,
	0x9c, 0xad, 0x5e, 0x40, 0x40, 0x68, 0x51, 0xbb, 0xa7, 0xa1, 0xb7, 0xf3, 0x0a, 0xe5, 0xde, 0x1b,
	0x61, 0xd0, 0x99, 0x2b, 0x04, 0xa0, 0x90, 0x4b, 0xf3, 0x76, 0xcf, 0x1b, 0xc5, 0x56, 0xab, 0xc5,
	0x35, 0xb3, 0xe8, 0xc7, 0x08, 0xa7, 0x8d, 0xc5, 0x08, 0x8d, 0xdb, 0xeb, 0x0b, 0xf0, 0x28, 0x75,
	0xb7, 0x56, 0xdb, 0xfb, 0x7b, 0xd4, 0x07, 0x9b, 0x11, 0x7f, 0x46, 0x19, 0x91, 0x53, 0xc2, 0x95,
	0xaa, 0xa5, 0xb1, 0xf8, 0xe2, 0x0c, 0x26, 0xae, 0xf4, 0x4d, 0x2d, 0xd3, 0x4f, 0x28, 0xe4, 0x5f,
	0xaa, 0xb0, 0x62, 0x79, 0x7e, 0x6c, 0x7b, 0xfc, 0x0c, 0x8e, 0xa1, 0x2f, 0x0a, 0xee, 0x7d, 0xb4,
	0x7b, 0xc7, 0x78, 0x66, 0xe4, 0xd1, 0x67, 0xfc, 0x00, 0xff, 0x6f, 0x6c, 0xb4, 0x1f, 0xe7, 0xd6,
	0xfa, 0xf2, 0xe0, 0x5c, 0x67, 0xad, 0xbd, 0x9e, 0xc1, 0x1d, 0x9e, 0xec, 0x3b, 0xf8, 0xac, 0xc9,
	0x13, 0xd2, 0x0d, 0x1a, 0x6c, 0x4d, 0xc1, 0x7f, 0xfc, 0x25, 0x2b, 0x5f, 0x6f, 0x0d, 0x56, 0x6c,
	0x92, 0xe0, 0xb0, 0xc3, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x7d, 0x5c, 0x2e, 0x40, 0xed, 0x06,
	0x00, 0x00,
}
