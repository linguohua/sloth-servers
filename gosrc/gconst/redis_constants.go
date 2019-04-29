package gconst

/*
const (
	// RoomConfigTable m 管理型服务器前缀，r表示房间，c表示配置
	RoomConfigTable = "m:r:c"
	// GsRoomTablePrefix 游戏服务器所维护的房间表格前缀
	GsRoomTablePrefix = "g:r:"
	// GameServerKeyPrefix g 表示game server，i表示instance实例
	GameServerKeyPrefix = "g:i:"
	// MgrServerKeyPrefix 管理型服务器实例表格
	MgrServerKeyPrefix = "m:i:"
	// PlayerTablePrefix 用于表格，c表示common，公用的意思，p表示用户
	PlayerTablePrefix = "c:p:"
	// AsUserTablePrefix m管理型服务器前缀，u表示用户,保存房间ID列表、用户昵称、用户性别、用户名字、用户头像url
	AsUserTablePrefix = "m:u:"
	// RoomNumberTable m管理型服务器前缀，rn表示房间号,
	RoomNumberTable = "m:rn:"
	// RoomTablePrefix c表示common，公用的意思，r表示房间，此表维护房间配置ID、创建者、房间号、房间所在的游戏服务器ID等
	RoomTablePrefix = "c:r:"
	// RoomTableACCSet ACC维护的房间set
	RoomTableACCSet = "c:r:accrs"

	// MsgListPrefix 消息队列前缀
	MsgListPrefix = "msg:"

	// ClubSysTable 俱乐部系统表格，主要是存储一些全局信息
	ClubSysTable = "m:club:sys"
	// PlayerClubSetPrefix 玩家俱乐部set
	PlayerClubSetPrefix = "m:uc:"
	// ClubEventTablePrefix 俱乐部事件哈希表，所有的消息，以ID为哈希表的field，以proto marshal后的buffer作为value
	ClubEventTablePrefix = "m:club:ev:"
	// ClubNeedHandledTablePrefix 需要部长处理的消息，用于快速确定某个消息是否需要处理，不需要转成proto才能确定
	ClubNeedHandledTablePrefix = "m:club:evh:"
	// ClubEventListPrefix 俱乐部事件列表，每一个俱乐部有一个自己的列表，俱乐部所有的事件都保存到该list中
	ClubEventListPrefix = "m:club:el:"
	// ClubFundEventListPrefix 俱乐部基金事件列表
	ClubFundEventListPrefix = "m:club:fel:"
	// ClubUnReadEventUserListPrefix 俱乐部用户事件列表，为俱乐部的每一个用户建立一个未读消息列表
	ClubUnReadEventUserListPrefix = "m:club:uel:"
	// ClubUnReadEventUserSetPrefix 俱乐部用户事件set,为俱乐部的每一个用户建立一个未读消息set，用于快速确定一个事件对于某个用户来说是否未读
	ClubUnReadEventUserSetPrefix = "m:club:ues:"
	// ClubTablePrefix 俱乐部表格，俱乐部的核心信息存储于这个前缀的哈希表中
	ClubTablePrefix = "m:club:"
	// ClubMemberSetPrefix 俱乐部成员set，每一个俱乐部都有自己的成员set
	ClubMemberSetPrefix = "m:club:m:"
	// ClubApplicantPrefix 俱乐部申请者列表
	ClubApplicantPrefix = "m:club:a:"
	// ClubReplayRoomsListPrefix 俱乐部回播房间ID列表
	ClubReplayRoomsListPrefix = "m:club:rl:"
	// ClubReplayRoomsSetPrefix 俱乐部回播房间ID set
	ClubReplayRoomsSetPrefix = "m:club:rs:"
	// MaxClubReplayRoomsNum 最多保存的俱乐部回播记录数量
	MaxClubReplayRoomsNum = 150

	// MJReplayRoomTablePrefix 麻将房间回播记录
	MJReplayRoomTablePrefix = "g:rr:"
	// ReplayRoomsReferenceSetPrefix 回播记录引用set
	ReplayRoomsReferenceSetPrefix = "g:rrr:"
	// MJRecorderShareIDTable 回播记录分享码哈希表
	MJRecorderShareIDTable = "g:rrsharedIDs"
	// MJRecorderID2ShareIDTable 回播记录索引到shareID
	MJRecorderID2ShareIDTable = "g:rr2sharedIDs"

	// MJRecorderDeletedSet 已经删除了的回播记录set
	MJRecorderDeletedSet = "g:rcds"
	// MJRecorderTablePrefix 麻将打牌记录表格
	MJRecorderTablePrefix = "g:rc:"

	// UserOrderRecord 用户的订单记录，m房间管理服务器，o表示支付，此表维护用户的支付记录
	UserOrderRecord = "m:o:"

	// RoomUnRefund 用户支付过但系统没返还给用户，m房间管理服务器，ur表示没返还的订单，此表维护未返还订单记录
	RoomUnRefund = "m:ur:"

	// RoomPayUsers 房间已经扣钻石用户列表，用于解散房间时返还钻石给用户, ru表示房间里面的用户
	RoomPayUsers = "m:ru:"

	// RoomRefundFailed 系统返回失败的订单，m房间管理服务器，rf表示返回失败的订单
	RoomRefundFailed = "m:rf:"

	// DonateTablePrefix 打赏表格， m房间管理服务器，d表示donate
	DonateTablePrefix = "m:d:"

	// UserDonatePrefix 用户打赏列表，可以是他打赏给别人或者别人打赏给他 m房间管理服务器，du表示donate user即打赏者或者被打赏者
	UserDonatePrefix = "m:du:"

	// UserBlacklist 用户黑名单，m房间管理服务器，b黑名单，u用户
	UserBlacklist = "m:b:u"

	// MonkeyAccountTablePrefix monkey账号表格，根据游戏类型区分，例如g:mk:1表示大丰麻将的monkey账号表
	MonkeyAccountTablePrefix = "g:mk:"

	// AccMonkeyAccountTable monkey账号表格，acc的monkey账号表格
	AccMonkeyAccountTable = "m:mkacc"

	// AccMonkeyTileCfgTable ACC管理的monkey打牌配置，主要是用于给代理测试人员通过web附加配置到房间
	AccMonkeyTileCfgTable = "m:mktc"

	// AccMonkeyTileCfgCategoryPrefix ACC管理的monkey打牌配置分类SET，形式如m:mktc:1:4
	AccMonkeyTileCfgCategoryPrefix = "m:mktc:"

	// GameServerDailyStatisTablePrefix 游戏每天针对玩家统计，哈希表key为g:yyyymmdd:userID:dsu(roomType
	// 其中完成局数为fh，赢牌局数为wh，创建并完成房间次数为cf，这个记录每晚上都会清理掉
	GameServerDailyStatisTablePrefix = "g:%s:%s:dsu%d"

	// RoomTypeSet 房间类型
	RoomTypeSet = "g:roomtype"

	// MJMaxReplayRoomNumber 麻将最大保存的回放记录
	MJMaxReplayRoomNumber = 50

	// WebAccountTalbePrefix 帐号
	WebAccountTalbePrefix = "m:ac:"

	//OnlineGameUserNum g:o:游戏在线人数
	OnlineGameUserNum = "g:o:"

	//PriceConfig 价格配置表
	PriceConfig = "m:p:"

	//PriceConfigDisable 待启用的价格配置表
	PriceConfigDisable = "m:pdis:"

	// PushServerList 推送服务器列表
	PushServerList = "pushserverlist"

	// AccMonkeyAccountTalbe monkey账号表格，acc的monkey账号表格
	AccMonkeyAccountTalbe = "m:mkacc"

	// RoomGameNo m 管理型服务器前缀，r表示房间，gn表示数据库生成的房间ID，即gameNo
	RoomGameNo = "m:r:gn"

	// EventPlayerConnect 用户连接事件
	EventPlayerConnect = "player_connect"

	// EventPlayerDisconnect 用户断开事件
	EventPlayerDisconnect = "player_disconnect"

	// EventPlayerLoginSuccess 用户登录成功事件
	EventPlayerLoginSuccess = "player_login_success"

	// ABTestUser 灰度用户列表
	ABTestUser = "m:abtu:"
	// ABTestUserSet 灰度用户集合
	ABTestUserSet = "m:abtus"

	// ABTestController 灰度控制表
	ABTestController = "m:abtc:"

	//SurrServerAvatarBox 头像框 gconst.SurrServerAvatarBox + userID + ":" + AvatarBoxID
	//如果该key存在,说明该玩家拥有AvatarBoxID的头像框
	SurrServerAvatarBox = "surr:ab:"

	// LobbyOnlinePlayerList 大厅在线用户列表
	LobbyOnlinePlayerList = "lobbyonlineplayerlist"

	// ClubShopPrefix 俱乐部基金模块redis前缀
	ClubShopPrefix = "m:club:sh"

	// AgentInfo 代理信息
	AgentInfo = "m:agent:"

	// PropsCfgsEnable 启用的道具配置表
	PropsCfgsEnable = "m:props:enable"

	// PropsCfgsDisable 未启用的道具配置表
	PropsCfgsDisable = "m:props:disable"

	// EventPlayerDisconnectAcc 通知acc用户断开事件
	EventPlayerDisconnectAcc = "player_disconnect_acc"

	// UserCreatRoomLock 用redis来实现创建房间锁，
	// 若这个用户已经在创建房间，那么这个用户就不能继续创建房间,只有等前一个创建完才可以继续创建
	UserCreatRoomLock = "m:lockcr:"

	// CurrentRoomGameNo 当前GameNo，逐渐加1
	CurrentRoomGameNo = "m:r:currentGN"

	// GroupRoomsSetPrefix 保存牌友群的房间
	GroupRoomsSetPrefix = "m:club:rooms:"

	// RedisKeyEventRoomInfoChange 俱乐部房间改变，发消息给俱乐部
	RedisKeyEventRoomInfoChange = "clubroominfochange"

	// GameRoomStatistics 统计牌局内数据， g表示game，st表示统计statistics
	GameRoomStatistics = "g:st:"

	// EmojiDailyStatisTablePrefix 统计表情的日常使用，哈希表key为emoji:stats:yyyymmdd:roomType
	EmojiDailyStatisTablePrefix = "emoji:stats:%s:%d"

	// UserClubTablePrefix 保存用户的所有牌友群
	UserClubTablePrefix = "cb:H:clubnum:"

	// GamePropsCfgTable 保存子游戏道具配置
	GamePropsCfgTable = "g:props:cfg"

	// ClubListKey 牌友群列表
	ClubListKey = "cb:H:clublist"
	// PropsTable 道具表
	PropsTable = "pp:tb"

	// RoomTypeKey 房间类型
	RoomTypeKey = "m:roomType:"

	// ClubMembersOther 牌友群相关
	ClubMembersOther = "cb:H:clubmembersother:"

	// Clubserverconfig 牌友群相关
	Clubserverconfig = "cb:H:clubserverconfig"
	// cb:H:clubmembersother:

	// GroupRoomSortSetPrefix 保存牌友群的房间,根据创建时间排序
	GroupRoomSortSetPrefix = "m:club:rooms:sort:"

	// GroupBigWinnerStats 统计茶馆大赢家次数, 格式:bw:groupID:date
	GroupBigWinnerStats = "bw:%s:%s"
	// GroupBigHandStats 统计茶馆对局次数
	GroupBigHandStats = "bh:%s:%s"

	// GroupSpecificRoomBigWinnerStats 统计茶馆大对应房间类型大赢家次数, 格式:bw:groupID:roomType:date
	GroupSpecificRoomBigWinnerStats = "bw:%s:%s:%s"

	// GroupSpecificRoomBigHandStats 统计茶馆大对应房间类型大赢家次数, 格式:bh:groupID:roomType:date
	GroupSpecificRoomBigHandStats = "bh:%s:%s:%s"

	// GroupStatsUpdateTime 茶馆统计的最后更新时间, 格式:update:groupID:date
	GroupStatsUpdateTime = "update:%s:%s"

	// GroupStatsConfirm 茶馆统计是否已经确认过
	GroupStatsConfirm = "comfirm:%s:%s"
	// GroupStatsSpecificeRoomConfirm 茶馆统计是否已经确认过
	GroupStatsSpecificeRoomConfirm = "comfirm:%s:%s:%s"

	// RaceRoomStateChange 竞技场房状态改变
	RaceRoomStateChange = "raceRoomStateChange"

	// GroupMemberRoomsSet 茶馆成功创建的房间
	GroupMemberRoomsSet = "m:Group:%s:%s"

	// MaxGroupReplayRoomsNum 最多保存的牌友群回播记录数量
	MaxGroupReplayRoomsNum = 300
)
*/
