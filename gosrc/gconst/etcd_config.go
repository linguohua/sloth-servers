package gconst

const (
	// EtcdRedisServerHost redis etcd 目录
	EtcdRedisServerHost = "/config/cache/redis"
	// EtcdDBProxyHost database etcd 目录
	EtcdDBProxyHost = "/config/dbproxy/address"
	// EtcdU8Auth u8 登录地址
	EtcdU8Auth = "/config/u8auth"
	// EtcdScheduleHost 策略服务器地址
	EtcdScheduleHost = "/config/schedule/address"
	// EtcdFileServerHost 文件服务器地址
	EtcdFileServerHost = "/config/fileserver/address"
	// EtcdUpdateDownURL etcd 更新下载地址
	EtcdUpdateDownURL = "/config/download/update"
	// EtcdShareDownURL 分享图片下载地址
	EtcdShareDownURL = "/config/download/share"
	// EtcdRankShareDownURL 排行榜分享下载
	EtcdRankShareDownURL = "/config/download/rankshare"
	// EtcdDBAgentOperSysHost AgentOperSys 代理库数据库地址
	EtcdDBAgentOperSysHost = "/config/sql/agentopersys/address"
	// EtcdDBAgentOperSysName 代理库名称
	EtcdDBAgentOperSysName = "/config/sql/agentopersys/name"
	// EtcdDBAgentOperSysUser AgentOperSys 代理库用户名
	EtcdDBAgentOperSysUser = "/config/sql/agentopersys/user"
	// EtcdDBAgentOperSysPassword AgentOperSys  代理库密码
	EtcdDBAgentOperSysPassword = "/config/sql/agentopersys/password"
	// EtcdDBGameHost 游戏数据库地址
	EtcdDBGameHost = "/config/sql/game/address"
	// EtcdDBGameName 游戏库名称
	EtcdDBGameName = "/config/sql/game/name"
	// EtcdDBGameUser 游戏库用户名
	EtcdDBGameUser = "/config/sql/game/user"
	// EtcdDBGamePassword 游戏库密码
	EtcdDBGamePassword = "/config/sql/game/password"
	// EtcdDBPlatFormHost 账号平台库地址
	EtcdDBPlatFormHost = "/config/sql/platform/address"
	// EtcdDBPlatFormName 账号平台名称
	EtcdDBPlatFormName = "/config/sql/platform/name"
	// EtcdDBPlatFormUser 账号平台用户名
	EtcdDBPlatFormUser = "/config/sql/platform/user"
	// EtcdDBPlatFormPassword 账号平台密码
	EtcdDBPlatFormPassword = "/config/sql/platform/password"
	// EtcdNginxConf  "Nginx配置文件整个文件"
	EtcdNginxConf = "/config/nginx/default/conf"
	// EtcdGameInstancesFormat 游戏服务器目录
	EtcdGameInstancesFormat = "/game-servers/instances/%s"
	// EtcdProductURL 产品的基础URL，例如ppl.qianz.com
	EtcdProductURL = "/config/baseurl"
	// EtcdAccInstanceDir ACC服务器配置目录
	EtcdAccInstanceDir = "/acc/instances"
)
