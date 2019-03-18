package lobby

import (
	"github.com/golang/protobuf/proto"
)

// ISessionMgr websocket mgr
type ISessionMgr interface {
	Broacast(msg []byte)
	SendTo(id string, msg []byte)
	SendProtoMsgTo(userID string, protoMsg proto.Message, opcode int32)
	UserCount() int
}
