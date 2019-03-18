module prserver

go 1.12

require (
	gconst v0.0.0
	github.com/DisposaBoy/JsonConfigReader v0.0.0-20171218180944-5ea4d0ddac55
	github.com/garyburd/redigo v1.6.0
	github.com/golang/protobuf v1.3.1
	github.com/gorilla/mux v1.7.0
	github.com/gorilla/websocket v1.4.0
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b
	github.com/sirupsen/logrus v1.4.0
	gscfg v0.0.0
	pokerface v0.0.0
)

replace gscfg => ../gscfg

replace pokerface => ../pokerface

replace gconst => ../gconst
