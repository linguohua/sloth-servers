package gconst

const (
	// GameServerRoomTypeSet 房间类型
	GameServerRoomTypeSet = "g:roomtype"
	// GameServerRoomTablePrefix 游戏服务器所维护的房间表格前缀
	GameServerRoomTablePrefix = "g:r:"
	// GameServerInstancePrefix g 表示game server，i表示instance实例
	GameServerInstancePrefix = "g:i:"
	// GameServerMJRecorderTablePrefix 麻将打牌记录表格
	GameServerMJRecorderTablePrefix = "g:rc:"
	// GameServerMJReplayRoomTablePrefix 麻将房间回播记录
	GameServerMJReplayRoomTablePrefix = "g:rr:"
	// GameServerReplayRoomsReferenceSetPrefix 回播记录引用set
	GameServerReplayRoomsReferenceSetPrefix = "g:rrr:"
	// GameServerMJRecorderDeletedSet 已经删除了的回播记录set
	GameServerMJRecorderDeletedSet = "g:rcds"
	// GameServerMJRecorderShareIDTable 回播记录分享码哈希表
	GameServerMJRecorderShareIDTable = "g:rrsharedIDs"
	// GameServerMonkeyAccountTablePrefix monkey账号表格，根据游戏类型区分，例如g:mk:1表示大丰麻将的monkey账号表
	GameServerMonkeyAccountTablePrefix = "g:mk:"
	// GameServerRoomStatisticsPrefix 统计牌局内数据， g表示game，st表示统计statistics
	GameServerRoomStatisticsPrefix = "g:st:"
	// GameServerDailyStatisTablePrefix 游戏每天针对玩家统计，哈希表key为g:yyyymmdd:userID:dsu(roomType
	// 其中完成局数为fh，赢牌局数为wh，创建并完成房间次数为cf，这个记录每晚上都会清理掉
	GameServerDailyStatisTablePrefix = "g:%s:%s:dsu%d"
	//GameServerOnlineUserNumPrefix g:o:游戏在线人数
	GameServerOnlineUserNumPrefix = "g:o:"
)
