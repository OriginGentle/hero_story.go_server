package websocket

import (
	"encoding/binary"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/reflect/protoreflect"
	"hero_story.go_server/biz_server/handler"
	"hero_story.go_server/biz_server/msg"
	"hero_story.go_server/comm/log"
	"hero_story.go_server/comm/main_thread"
	"time"
)

const oneSecond = 1000
const readMsgCountPreSecond = 16
const msgSize = 64 * 1024

// CmdContextImpl 就是ICmdContext的Websocket实现
type CmdContextImpl struct {
	userId       int64
	clientIpAddr string
	Conn         *websocket.Conn
	sendMsgQ     chan protoreflect.ProtoMessage // BlockQueue
	SessionId    int32
}

func (ctx *CmdContextImpl) BindUserId(val int64) {
	ctx.userId = val
}

func (ctx *CmdContextImpl) GetUserId() int64 {
	return ctx.userId
}

func (ctx *CmdContextImpl) GetClientIpAddr() string {
	return ctx.clientIpAddr
}

func (ctx *CmdContextImpl) Write(msgObj protoreflect.ProtoMessage) {
	if nil == msgObj ||
		nil == ctx.Conn ||
		nil == ctx.sendMsgQ {
		return
	}
	ctx.sendMsgQ <- msgObj
}

func (ctx *CmdContextImpl) SendError(errorCode int, errorInfo string) {
	if len(errorInfo) <= 0 &&
		errorCode < 0 {
		return
	}
}

func (ctx *CmdContextImpl) Disconnect() {
	if nil != ctx.Conn {
		_ = ctx.Conn.Close()
	}
}

func (ctx *CmdContextImpl) LoopSendMsg() {
	ctx.sendMsgQ = make(chan protoreflect.ProtoMessage, 64)

	go func() { // 启动一个协程，负责发送消息
		for {
			msgObj := <-ctx.sendMsgQ

			if nil == msgObj {
				continue
			}

			byteArray, err := msg.Encode(msgObj)

			if nil != err {
				log.Error("消息编码异常，err = %+v", err)
				return
			}

			if err := ctx.Conn.WriteMessage(websocket.BinaryMessage, byteArray); nil != err {
				log.Error("消息发送异常，err = %+v", err)
			}
		}
	}()
}

func (ctx *CmdContextImpl) LoopReadMsg() {
	if nil == ctx.Conn {
		return
	}

	ctx.Conn.SetReadLimit(msgSize)

	t0 := int64(0)
	counter := 0

	for {
		_, msgData, err := ctx.Conn.ReadMessage()

		if nil != err {
			log.Error("消息读取异常，err = %+v", err)
			break
		}

		t1 := time.Now().UnixMilli()

		if t1-t0 > oneSecond {
			t0 = t1
			counter = 0
		}

		if counter >= readMsgCountPreSecond {
			log.Error("消息发送过于频繁")
			continue
		}
		counter++

		msgCode := binary.BigEndian.Uint16(msgData[2:4])
		newMsgX, err := msg.Decode(msgData[4:], int16(msgCode))

		if nil != err {
			log.Error("消息解码错误,msgCode = %d ,err = %+v ",
				msgCode, err,
			)
			continue
		}

		log.Info("收到客户端消息,msgCode = %d,msgName = %s",
			msgCode, newMsgX.Descriptor().Name(),
		)

		cmdHandler := handler.CreateCmdHandler(msgCode)

		if nil == cmdHandler {
			log.Error("未找到指令处理器,msgCode = %d",
				msgCode,
			)
			continue
		}

		main_thread.Process(func() {
			cmdHandler(ctx, newMsgX)
		})
	}

	handler.OnUserQuit(ctx)
}
