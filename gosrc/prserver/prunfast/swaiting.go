package prunfast

import (
	"pokerface"

	"github.com/golang/protobuf/proto"
)

// SWaiting 等待状态
type SWaiting struct {
	room *Room
}

// getStateConst 等待状态标志
func (s *SWaiting) getStateConst() pokerface.RoomState {
	return pokerface.RoomState_SRoomWaiting
}

func (s *SWaiting) getStateName() string {
	return "SWaiting"
}

// onMessage 等待状态下，处理用户消息
func (s *SWaiting) onMessage(iu IUser, gmsg *pokerface.GameMessage) {
	// 等待状态下，只处理用户ready消息，当所有人都Ready后
	// 进入SPlaying状态
	player := s.room.getPlayerByUserID(iu.userID())
	if player == nil {
		s.room.cl.Panic("got ready message but player is null")
		return
	}
	op := gmsg.GetOps()
	switch op {
	case int32(pokerface.MessageCode_OPPlayerReady):
		s.room.cl.Println("got player ready:", player.chairID)
		if player.state == pokerface.PlayerState_PSReady {
			// 重复收到ready消息
			break
		}

		player.state = pokerface.PlayerState_PSReady

		// 有用户状态发生了变化，更新给所有客户端
		s.room.updateRoomInfo2All()

		allReady := true
		for _, p := range s.room.players {
			if p.state != pokerface.PlayerState_PSReady {
				allReady = false
				break
			}
		}

		if !allReady {
			break
		}

		// 人数未够
		if len(s.room.players) < s.room.config.playerNumAcquired {
			break
		}

		if s == s.room.state {
			s.room.state2(s, pokerface.RoomState_SRoomPlaying)
		}
		break
	case int32(pokerface.MessageCode_OPKickout):
		// 踢人
		if player.userID() != s.room.ownerID {
			sendKickoutError(player, pokerface.KickoutResult_KickoutResult_FailedNeedOwner)
			return
		}

		if s.room.handRoundStarted > 0 {
			sendKickoutError(player, pokerface.KickoutResult_KickoutResult_FailedGameHasStartted)
			return
		}

		kickoutMsg := &pokerface.MsgKickout{}
		err := proto.Unmarshal(gmsg.Data, kickoutMsg)
		if err != nil {
			s.room.cl.Println("proc kickout request failed:", err)
			return
		}

		if kickoutMsg.GetVictimUserID() == player.userID() {
			s.room.cl.Println("proc kickout request failed,can't kickout self")
			return
		}

		victim := s.room.getPlayerByUserID(kickoutMsg.GetVictimUserID())
		if victim == nil {
			sendKickoutError(player, pokerface.KickoutResult_KickoutResult_FailedPlayerNotExist)
			return
		}

		s.doKickoutPlayer(player, victim)

		if s.room.rbl == nil {
			s.room.rbl = newRoomBlockList(s.room)
		}

		s.room.rbl.blockUser(victim.userID())
		break
	default:
		s.room.cl.Println("Waiting state can not process:", op)
		break
	}
}

func (s *SWaiting) doKickoutPlayer(player *PlayerHolder, victim *PlayerHolder) {
	msg := &pokerface.MsgKickoutResult{}
	var result32 = int32(pokerface.KickoutResult_KickoutResult_Success)
	msg.Result = &result32
	var byWhoNick = player.user.getInfo().nick
	if byWhoNick == "" {
		byWhoNick = player.userID()
	}

	var owerID = player.userID()
	msg.ByWhoUserID = &owerID
	msg.ByWhoNick = &byWhoNick
	var victimUserID = victim.userID()
	msg.VictimUserID = &victimUserID
	var victimNick = victim.user.getInfo().nick
	if victimNick == "" {
		victimNick = victimUserID
	}
	msg.VictimNick = &victimNick

	for _, p := range s.room.players {
		p.sendMsg(msg, int32(pokerface.MessageCode_OPKickout))
	}

	// 断开websocket连接
	victim.user.detach()

	// 离线处理
	victim.allowedLeave = true // 强制准许离开
	s.room.onUserOffline(victim.user, false)
}

func sendKickoutError(player *PlayerHolder, result pokerface.KickoutResult) {
	msg := &pokerface.MsgKickoutResult{}
	var status32 = int32(result)
	msg.Result = &status32

	player.sendMsg(msg, int32(pokerface.MessageCode_OPKickout))
}

// onStateEnter 进入等待状态
func (s *SWaiting) onStateEnter() {
	s.room.cl.Println("room enter waiting state")
}

// onStateLeave 离开等待状态
func (s *SWaiting) onStateLeave() {
	// DO nothing!
	s.room.cl.Println("room leave waiting state")
}

// onPlayerEnter 等待状态下处理用户进入
func (s *SWaiting) onPlayerEnter(player *PlayerHolder) {

	room := s.room
	player.state = pokerface.PlayerState_PSNone

	// 玩家进入后立即发送房间信息，在延后一会客户端会发送
	// ready消息上来，然后又需要发送房间信息
	room.updateRoomInfo2All()
}

// onPlayerLeave 处理用户离线
func (s *SWaiting) onPlayerLeave(player *PlayerHolder) {
	// 等待状态允许玩家自由离开
	room := s.room
	player.state = pokerface.PlayerState_PSOffline

	// 如果是不允许离开，则不删除player
	if !player.allowedLeave {
		room.updateRoomInfo2All()
		return
	}

	room.stateRemovePlayer(player)

	if len(room.players) < 1 {
		// 如果玩家都离开了，则回到空闲状态
		room.state2(s, pokerface.RoomState_SRoomIdle)
	}

	room.updateRoomInfo2All()
}

// onPlayerReEnter 等待状态下处理用户重入
func (s *SWaiting) onPlayerReEnter(player *PlayerHolder) {
	// 等待状态下不应该发生用户重新进入，
	// 因为如果用户离线，等待状态下，会立即删除了用户对象
	// 这里需要注释掉，因为虽然用户离线，但是服务器可能还未能
	// 感应到，因此player对象还存在，故客户端重连上来后，就会进入这里
	// s.room.cl.Panic("waiting state should not have re-enter event")
	player.state = pokerface.PlayerState_PSNone
	s.room.updateRoomInfo2All()
}
