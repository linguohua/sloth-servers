package gpubsub

import (
	"gconst"
	"time"

	"github.com/golang/protobuf/proto"

	log "github.com/sirupsen/logrus"
)

// waitSubcriberRsp 等待游戏服务器的返回
type waitSubcriberRsp struct {
	waitChan chan bool
	rspMsg   *gconst.SSMsgBag
}

// SendAndWait 给dst发送消息（通过redis推送），并等待回复，timeout 指定超时时间
func SendAndWait(dst string, msg *gconst.SSMsgBag, timeout time.Duration) (bool, *gconst.SSMsgBag) {
	if dst == "" {
		log.Panicln("SendAndWait, need dst")
		return false, nil
	}

	if msg == nil {
		log.Panicln("SendAndWait, msg == nil")
		return false, nil
	}

	// 填上源url，以便对方可以发回回复
	msg.SourceURL = &myServerID

	var wait = &waitSubcriberRsp{}
	wait.waitChan = make(chan bool, 1)
	waitingMap[int(msg.GetSeqNO())] = wait

	PublishMsg(dst, msg)

	var rspGot = false
	select {
	case <-wait.waitChan:
		rspGot = true
		break
	case <-time.After(timeout):
		break
	}

	// 任何情况都删除这个seqNo
	delete(waitingMap, int(msg.GetSeqNO()))
	return rspGot, wait.rspMsg
}

// PublishMsg 往redis publish消息
func PublishMsg(dst string, msg *gconst.SSMsgBag) {
	bytes, err := proto.Marshal(msg)
	if err != nil {
		log.Println(err)
		return
	}

	// 获取redis链接，并退出函数时释放
	conn := pool.Get()
	defer conn.Close()

	// 写入到目标服务器的消息队列
	var targetMsgListID = gconst.MsgListPrefix + dst
	_, err = conn.Do("RPUSH", targetMsgListID, bytes)
	if err != nil {
		log.Panic("PublishMsg failed, RPUSH error:", err)
	}

	// 发送一个通知给目标服务器
	conn.Do("PUBLISH", dst, []byte("h"))
}
