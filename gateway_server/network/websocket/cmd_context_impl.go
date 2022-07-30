package websocket

import (
	"github.com/gorilla/websocket"
	"hero_story.go_server/biz_server/msg"
	"hero_story.go_server/comm/log"
	"hero_story.go_server/gateway_server/cluster/biz_server_finder"
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
	sendMsgQ     chan []byte // BlockQueue
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

func (ctx *CmdContextImpl) Write(byteArray []byte) {
	if nil == byteArray ||
		nil == ctx.Conn ||
		nil == ctx.sendMsgQ {
		return
	}
	ctx.sendMsgQ <- byteArray
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

// LoopSendMsg 循环向客户端发送消息
func (ctx *CmdContextImpl) LoopSendMsg() {
	ctx.sendMsgQ = make(chan []byte, 64)

	go func() { // 启动一个协程，负责发送消息
		for {
			byteArray := <-ctx.sendMsgQ

			if nil == byteArray {
				continue
			}

			func() {
				defer func() {
					if err := recover(); nil != err {
						log.Error("发生异常，%+v", err)
					}
				}()

				if err := ctx.Conn.WriteMessage(websocket.BinaryMessage, byteArray); nil != err {
					log.Error("消息发送异常，%+v", err)
				}
			}()
		}
	}()
}

// LoopReadMsg 循环读取从客户端发送的消息
// 游戏客户端 --> 网关服务器
func (ctx *CmdContextImpl) LoopReadMsg() {
	if nil == ctx.Conn {
		return
	}

	ctx.Conn.SetReadLimit(msgSize)

	t0 := int64(0)
	counter := 0

	// 创建到游戏服的连接
	bizServerConn, err := biz_server_finder.GetBizServerConn()

	if nil != err {
		log.Error("%+v", err)
		return
	}

	// 循环读取从游戏客户端发送的消息
	// 转发给游戏服务器
	for {
		msgType, msgData, err := ctx.Conn.ReadMessage()

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

		func() {
			defer func() {
				if err := recover(); nil != err {
					log.Error("发生异常，%+v", err)
				}
			}()

			log.Info("收到客户端消息并转发")

			// 包装消息
			innerMsg := &msg.InternalServerMsg{
				GatewayServerId: 0,
				SessionId:       ctx.SessionId,
				UserId:          ctx.GetUserId(),
				MsgData:         msgData,
			}

			innerMsgByteArray := innerMsg.ToByteArray()

			// 将客户端消息转发给游戏服
			if err := bizServerConn.WriteMessage(msgType, innerMsgByteArray); nil != err {
				log.Error("消息转发失败，%+v", err)
			}
		}()
	}
}
