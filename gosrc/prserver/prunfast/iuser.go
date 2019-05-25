package prunfast

import "github.com/gorilla/websocket"

// IUser 用户接口
type IUser interface {
	// userID 获得64位唯一用户ID
	userID() string
	// send 发送消息
	send(bytes []byte)
	// onWebsocketClosed 通知websocket断开
	onWebsocketClosed(ws *websocket.Conn)
	// onWebsocketMessage 处理websocket消息
	onWebsocketMessage(ws *websocket.Conn, message []byte)
	// 重新绑定一个websocket
	rebind(ws *websocket.Conn)
	// 关闭连接并断开和room关联
	detach()
	getRoom() *Room
	// 获取用户信息，包括昵称、性别、头像URI
	getInfo() *UserInfo
	// 更新用户信息
	updateInfo()
	// 关闭websocket连接
	closeWebsocket()

	// 发送ping给玩家
	sendPing()

	// 发送pong给玩家
	sendPong(msg string)

	setFromWeb(isFromWeb bool)

	isFromWeb() bool
}
