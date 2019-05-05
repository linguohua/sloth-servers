package dfmahjong

import (
	"log"
	"mahjong"

	"github.com/golang/protobuf/proto"
)

// PlayerHolder 表示一个玩家
type PlayerHolder struct {
	room           *Room
	hStatis        *HStatis
	gStatis        *GStatis
	user           IUser
	tiles          *PlayerTiles
	expectedAction int
	sctx           *ScoreContext
	chairID        int
	state          mahjong.PlayerState

	allowedLeave bool // 服务器批准离开
}

// newPlayerHolder 新建一个玩家
func newPlayerHolder(room *Room, chairID int, iu IUser) *PlayerHolder {
	p := &PlayerHolder{}
	p.hStatis = newHStatis()
	p.gStatis = newGStatis()
	p.user = iu
	p.room = room
	p.tiles = newPlayerTiles(p)
	p.chairID = chairID
	p.state = mahjong.PlayerState_PSNone
	return p
}

// userID 获得玩家的userID
func (p *PlayerHolder) userID() string {
	return p.user.userID()
}

// sendDealMsg 发送发牌消息
func (p *PlayerHolder) sendDealMsg(msg *mahjong.MsgDeal) {
	p.sendMsg(msg, int32(mahjong.MessageCode_OPDeal))
}

// sendActoinAllowedMsg 发送允许动作给玩家
func (p *PlayerHolder) sendActoinAllowedMsg(msgAllowedAction *mahjong.MsgAllowPlayerAction) {
	p.sendMsg(msgAllowedAction, int32(mahjong.MessageCode_OPActionAllowed))
}

// sendReActoinAllowedMsg 发送允许反应给玩家
func (p *PlayerHolder) sendReActoinAllowedMsg(msgAllowedReAction *mahjong.MsgAllowPlayerReAction) {
	p.sendMsg(msgAllowedReAction, int32(mahjong.MessageCode_OPReActionAllowed))
}

// sendActionResultNotify 发送动作通知给玩家
func (p *PlayerHolder) sendActionResultNotify(msgActionResultNotify *mahjong.MsgActionResultNotify) {
	p.sendMsg(msgActionResultNotify, int32(mahjong.MessageCode_OPActionResultNotify))
}

// sendHandOver 发送手牌结束结果给玩家
func (p *PlayerHolder) sendHandOver(msgHandOver *mahjong.MsgHandOver) {
	p.sendMsg(msgHandOver, int32(mahjong.MessageCode_OPHandOver))
}

func (p *PlayerHolder) sendTipsString(tips string) {
	msgTips := &mahjong.MsgRoomShowTips{}
	msgTips.Tips = &tips
	var tipCode32 = int32(mahjong.TipCode_TCNone)
	msgTips.TipCode = &tipCode32

	p.sendMsg(msgTips, int32(mahjong.MessageCode_OPRoomShowTips))
}

func (p *PlayerHolder) sendTipsCode(tipCode mahjong.TipCode) {
	msgTips := &mahjong.MsgRoomShowTips{}
	var tipCode32 = int32(tipCode)
	msgTips.TipCode = &tipCode32

	p.sendMsg(msgTips, int32(mahjong.MessageCode_OPRoomShowTips))
}

// sendMsg 给玩家发送 GameMessage
func (p *PlayerHolder) sendMsg(pb proto.Message, ops int32) {
	gmsg := &mahjong.GameMessage{}
	gmsg.Ops = &ops

	if pb != nil {
		bytes, err := proto.Marshal(pb)

		if err != nil {
			log.Panic("marshal msg failed:", err)
			return
		}
		gmsg.Data = bytes
	}

	bytes, err := proto.Marshal(gmsg)
	if err != nil {
		log.Panic("marshal game msg failed:", err)
		return
	}

	p.user.send(bytes)
}

// resetForNextHand 新一手牌开始时，重设player对象中的相关变量
func (p *PlayerHolder) resetForNextHand() {
	p.hStatis.reset()
	p.tiles.clear()
}
