module lobbyserver

go 1.12

require (
	gconst v0.0.0
	github.com/DisposaBoy/JsonConfigReader v0.0.0-20171218180944-5ea4d0ddac55
	github.com/coreos/etcd v3.3.12+incompatible
	github.com/garyburd/redigo v1.6.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/protobuf v1.3.1
	github.com/gorilla/mux v1.7.0
	github.com/gorilla/websocket v1.4.0
	github.com/huayuego/wordfilter v0.0.0-20171103075834-036271d0abf0
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
	github.com/sirupsen/logrus v1.4.0
	gpubsub v0.0.0

	gscfg v0.0.0
)

replace gscfg => ../gscfg

replace gconst => ../gconst

replace gpubsub => ../gpubsub
