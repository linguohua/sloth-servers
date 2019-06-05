package gconst

const (
	// LobbyMsgListPrefix 消息队列前缀
	LobbyMsgListPrefix = "l:msg:"
	// LobbyRoomNumberTablePrefix m管理型服务器前缀，rn表示房间号,
	LobbyRoomNumberTablePrefix = "l:rn:"
	// LobbyPlayerTablePrefix 用于表格，p表示用户
	LobbyPlayerTablePrefix = "l:p:"
	// LobbyRoomConfigTable r表示房间，c表示配置
	LobbyRoomConfigTable = "l:rcfg"
	// LobbyUserTablePrefix m管理型服务器前缀，u表示用户,保存房间ID列表、用户昵称、用户性别、用户名字、用户头像url
	LobbyUserTablePrefix = "l:u:"
	// LobbyUserSet 用户列表
	LobbyUserSet = "l:uset"
	// LobbyRoomTablePrefix c表示common，公用的意思，r表示房间，此表维护房间配置ID、创建者、房间号、房间所在的游戏服务器ID等
	LobbyRoomTablePrefix = "l:r:"
	// LobbyRoomTableSet ACC维护的房间set
	LobbyRoomTableSet = "l:rset"
	//LobbyPriceConfigPrefix 价格配置表
	LobbyPriceConfigPrefix = "l:price:"
	//LobbyPriceConfigDisablePrefix 待启用的价格配置表
	LobbyPriceConfigDisablePrefix = "l:pricedis:"
	// LobbyUserBlacklistSet 用户黑名单，m房间管理服务器，b黑名单，u用户
	LobbyUserBlacklistSet = "l:b:u"
	// LobbyMonkeyAccountTalbe monkey账号表格，acc的monkey账号表格
	LobbyMonkeyAccountTalbe = "l:mkacc"
	// LobbyPropsCfgTable 保存子游戏道具配置
	LobbyPropsCfgTable = "l:props:cfg"
	// LobbyPropsTable 道具表
	LobbyPropsTable = "l:tb"

	// LobbyPayOrderPrefix 房间支付记录
	LobbyPayOrderPrefix = "l:pay:"
	// // LobbyRoomPayUsersPrefix 房间已经扣钻石用户列表，用于解散房间时返还钻石给用户, ru表示房间里面的用户
	// LobbyRoomPayUsersPrefix = "l:pay:ru:"
	// // LobbyPayUserOrderPrefix 用户的订单记录，m房间管理服务器，o表示支付，此表维护用户的支付记录
	// LobbyPayUserOrderPrefix = "l:pay:o:"
	// // LobbyPayRoomRefundFailedPrefix 系统返回失败的订单，m房间管理服务器，rf表示返回失败的订单
	// LobbyPayRoomRefundFailedPrefix = "l:pay:rf:"

	// LobbyUserCreatRoomLockPrefix 用redis来实现创建房间锁，
	// 若这个用户已经在创建房间，那么这个用户就不能继续创建房间,只有等前一个创建完才可以继续创建
	LobbyUserCreatRoomLockPrefix = "l:lockcr:"

	// LobbyChatMessagePrefix 用户聊天消息
	LobbyChatMessagePrefix = "l:chat:"

	// LobbyMailPrefix 用户邮件
	LobbyMailPrefix = "l:mail:"

	// LobbyMailSortSetPrefix 用户有序的邮件列表
	LobbyMailSortSetPrefix = "l:mailsort:"

	// LobbyDatabaseConfig 数据库配置
	LobbyDatabaseConfig = "l:db:cfg"

	// LobbyWeChatConfig 微信配置
	LobbyWeChatConfig = "l:wc:cfg"

	// LobbyUpgradeModuleSet 所有配置了更新的模块的名字set
	LobbyUpgradeModuleSet = "l:upset"

	// LobbyUpgradeModuleHashPrefix 所有配置了更新的模块的名字set
	LobbyUpgradeModuleHashPrefix = "l:up:"

	// LobbyUserDonatePrefix 用户道具表
	LobbyUserDonatePrefix = "l:d:"

	// LobbyClubApplicantPrefix 俱乐部申请者列表
	LobbyClubApplicantPrefix = "l:club:a:"

	// LobbyClubEventTablePrefix 俱乐部事件哈希表，所有的消息，以ID为哈希表的field，以proto marshal后的buffer作为value
	LobbyClubEventTablePrefix = "l:club:ev:"
	// LobbyClubNeedHandledTablePrefix 需要部长处理的消息，用于快速确定某个消息是否需要处理，不需要转成proto才能确定
	LobbyClubNeedHandledTablePrefix = "l:club:evh:"
	// LobbyClubEventListPrefix 俱乐部事件列表，每一个俱乐部有一个自己的列表，俱乐部所有的事件都保存到该list中
	LobbyClubEventListPrefix = "l:club:el:"
	// LobbyClubUnReadEventUserListPrefix 俱乐部用户事件列表，为俱乐部的每一个用户建立一个未读消息列表
	LobbyClubUnReadEventUserListPrefix = "l:club:ur:uel:"
	// LobbyClubUnReadEventUserSetPrefix 俱乐部用户事件set,为俱乐部的每一个用户建立一个未读消息set，用于快速确定一个事件对于某个用户来说是否未读
	LobbyClubUnReadEventUserSetPrefix = "l:club:ur:ues:"
	// LobbyClubSysTable 俱乐部系统表格，主要是存储一些全局信息
	LobbyClubSysTable = "m:club:sys"
	// LobbyClubMemberSetPrefix 牌友群中的成员set
	LobbyClubMemberSetPrefix = "l:club:m:"
	// LobbyClubRoomSetPrefix 牌友群房间Set
	LobbyClubRoomSetPrefix = "l:club:rs:"
	// LobbyClubUserApplicantEventPrefix 用户的申请记录列表
	LobbyClubUserApplicantEventPrefix = "l:club:uel:"
)
