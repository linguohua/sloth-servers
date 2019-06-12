package pddz

import (
	"pokerface"
	"runtime/debug"
	"time"
)

const (
	timeWaitClientReply = 180
)

// RoomDisband 用于解散房间的数据结构
// 等待玩家回复，记录每个玩家的回复
// 并处理解散逻辑
type RoomDisband struct {
	waitItems []*RoomDisbandWait // 等待其他玩家回复
	applicant *PlayerHolder      // 申请发起者
	room      *Room

	state pokerface.DisbandState
	chanx chan bool

	startTime time.Time

	// SPlaying 等待
	chany chan bool
}

// RoomDisbandWait 所有其他玩家对解散房间请求的响应
type RoomDisbandWait struct {
	player  *PlayerHolder
	isReply bool
	isAgree bool
}

// sendNotify2All 发送解散状态通知给所有玩家
func (rdb *RoomDisband) sendNotify2All() {
	msg := &pokerface.MsgDisbandNotify{}
	var state32 = int32(rdb.state)
	msg.DisbandState = &state32
	var applicantChairID32 = int32(rdb.applicant.chairID)
	msg.Applicant = &applicantChairID32

	var timeRemain = timeWaitClientReply - time.Now().Sub(rdb.startTime).Seconds()
	if timeRemain < 0 {
		timeRemain = 0
	}
	var timeRemain32 = int32(timeRemain)
	msg.Countdown = &timeRemain32

	agrees := make([]int32, 0, len(rdb.waitItems))
	rejects := make([]int32, 0, len(rdb.waitItems))
	waits := make([]int32, 0, len(rdb.waitItems))

	// 同意者列表
	for _, wi := range rdb.waitItems {
		if wi.isReply && wi.isAgree {
			agrees = append(agrees, int32(wi.player.chairID))
		}
	}

	// 拒绝者列表
	for _, wi := range rdb.waitItems {
		if wi.isReply && !wi.isAgree {
			rejects = append(rejects, int32(wi.player.chairID))
		}
	}

	// 等待者列表
	for _, wi := range rdb.waitItems {
		if !wi.isReply {
			waits = append(waits, int32(wi.player.chairID))
		}
	}

	msg.Waits = waits
	msg.Agrees = agrees
	msg.Rejects = rejects

	for _, p := range rdb.room.players {
		p.sendMsg(msg, int32(pokerface.MessageCode_OPDisbandNotify))
	}
}

// startDisband 开始解散流程
func (rdb *RoomDisband) startDisband() {
	defer func() {
		if r := recover(); r != nil {
			roomExceptionCount++
			debug.PrintStack()
			rdb.room.cl.Printf("-----This Disban GR will die, roomID:%s, Recovered in startDisband:%v\n", rdb.room.ID, r)
		}
	}()

	rdb.startTime = time.Now()
	rdb.state = pokerface.DisbandState_Waiting
	rdb.chany = make(chan bool, 1)

	rdb.sendNotify2All()

	if len(rdb.waitItems) > 0 {
		rdb.chanx = make(chan bool, 1)
		timeout := false

		// 等待其他玩家回复
		select {
		case <-rdb.chanx:
			break
		case <-time.After((timeWaitClientReply + 3) * time.Second):
			timeout = true
			break
		}

		if timeout {
			rdb.room.cl.Println("wait disband answer timeout, assume all agree")
			// rdb.cancelDisband(DisbandState_DoneWithWaitReplyTimeout)
			for _, wi := range rdb.waitItems {
				if !wi.isReply {
					wi.isAgree = true
					wi.isReply = true
				}
			}

			rdb.sendNotify2All()
		}
	}

	allAgree := true
	for _, wi := range rdb.waitItems {
		if wi.isReply && !wi.isAgree {
			allAgree = false
			break
		}
	}

	if !allAgree {
		rdb.room.cl.Println("room disband cancel, someone disagree")
		rdb.cancelDisband(pokerface.DisbandState_DoneWithOtherReject)
		return
	}

	rdb.doDisband()
}

// onDisbandAnswer 处理玩家的解散咨询回复
func (rdb *RoomDisband) onDisbandAnswer(player *PlayerHolder, replyMsg *pokerface.MsgDisbandAnswer) {
	var wii *RoomDisbandWait
	for _, wi := range rdb.waitItems {
		if player == wi.player {
			wii = wi
			break
		}
	}

	if wii == nil {
		rdb.room.cl.Println("not expected player when wait disband answer")
		return
	}

	if wii.isReply {
		rdb.room.cl.Println("user already reply for disband, discarded answer")
		return
	}

	wii.isReply = true
	if replyMsg != nil {
		wii.isAgree = replyMsg.GetAgree()
	} else {
		// 离线玩家，默认为同意
		wii.isAgree = true
	}

	rdb.room.cl.Printf("onDisbandAnswer, room %s, user %s answer disband request:%+v\n",
		rdb.room.roomNumber, player.userID(), wii.isAgree)

	if !wii.isAgree {
		// 有人拒绝了
		rdb.chanx <- true
		return
	}

	allReply := true
	for _, wi := range rdb.waitItems {
		if !wi.isReply {
			allReply = false
			break
		}
	}

	rdb.sendNotify2All()

	if allReply {
		rdb.chanx <- true
	}
}

// doDisband 执行解散流程
func (rdb *RoomDisband) doDisband() {
	// 请求房间管理服务器删除房间
	// if rdb.room.requireRoomServer2Delete() == false {
	// 	rdb.room.cl.Println("requireRoomServer2Delete failed, will not disband:", rdb.room.ID)
	// 	rdb.cancelDisband(DisbandState_DoneWithRoomServerNotResponse)
	// 	return
	// }

	// 通知所有人，解散已经完成
	rdb.state = pokerface.DisbandState_Done
	rdb.sendNotify2All()

	// rdb.room.disband = nil

	// 通知SPlaying等待结束，房间解散，游戏结束
	rdb.chany <- false

	// 发送结算到所有客户端
	if rdb.room.handRoundFinished > 0 {
		rdb.room.onGameOver(nil)
	}

	roomMgr.forceDisbandRoom(rdb.room, pokerface.RoomDeleteReason_DisbandByApplication)
	rdb.room.disband = nil
}

// cancelDisband 取消解散
func (rdb *RoomDisband) cancelDisband(state pokerface.DisbandState) {
	rdb.room.cl.Println("disband cancel, room ID:", rdb.room.ID)
	rdb.room.disband = nil

	rdb.state = state
	rdb.sendNotify2All()

	// 通知SPlaying等待结束，继续游戏
	rdb.chany <- true
}
