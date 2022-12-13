package websocket

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"github.com/gorilla/websocket"
	"hero_story.go_server/biz_server/base"
	"hero_story.go_server/biz_server/msg"
	"hero_story.go_server/comm/log"
	"hero_story.go_server/gateway_server/cluster"
	"hero_story.go_server/gateway_server/cluster/biz_server_finder"
	"time"
)

const oneSecond = 1000
const readMsgCountPreSecond = 16
const msgSize = 64 * 1024

// CmdContextImpl 就是ICmdContext的Websocket实现
type CmdContextImpl struct {
	userId          int64
	clientIpAddr    string
	Conn            *websocket.Conn
	sendMsgQ        chan []byte // BlockQueue
	SessionId       int32
	GatewayServerId int32
	GameServerId    int32
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
	if nil == byteArray || nil == ctx.Conn || nil == ctx.sendMsgQ {
		return
	}

	ctx.sendMsgQ <- byteArray
}

func (ctx *CmdContextImpl) SendError(errorCode int, errorInfo string) {
	if len(errorInfo) <= 0 && errorCode < 0 {
		return
	}
}

func (ctx *CmdContextImpl) Disconnect() {
	if nil != ctx.Conn {
		_ = ctx.Conn.Close()
	}
}

// LoopSendMsg 循环向客户端发送消息
// 内容通过协程实现
func (ctx *CmdContextImpl) LoopSendMsg() {
	ctx.sendMsgQ = make(chan []byte, 64)

	go func() { // 专门启动一个协程，负责发送消息
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

// LoopReadMsg 循环读取客户端发送的消息
// 游戏客户端 --> 网关服务器
func (ctx *CmdContextImpl) LoopReadMsg() {
	if nil == ctx.Conn {
		return
	}
	ctx.Conn.SetReadLimit(msgSize)

	t0 := int64(0)
	counter := 0

	// 循环读取从游戏客户端发送的消息
	// 转发给游戏服务器
	for {
		// 接收从游戏客户端发来的消息
		msgType, msgData, err := ctx.Conn.ReadMessage()
		if nil != err {
			log.Error("游戏客户端消息读取异常,err = %+v", err)
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
		serverJobType := biz_server_finder.GetServerJobTypeByMsgCode(msgCode)

		// 创建到游戏服务的连接
		bizServerConn, err := biz_server_finder.GetBizServerConn(
			serverJobType,
			getBizServerFindStrategy(ctx, serverJobType),
		)

		if nil != err {
			log.Error("%+v", err)
			continue
		}

		func() {
			defer func() {
				if err := recover(); nil != err {
					log.Error("发生异常，%+v", err)
				}
			}()
			log.Info("收到游戏客户端消息并转发")

			// 包装消息
			innerMsg := &msg.InternalServerMsg{
				GatewayServerId: ctx.GatewayServerId,
				SessionId:       ctx.SessionId,
				UserId:          ctx.GetUserId(),
				MsgData:         msgData,
			}
			innerMsgByteArray := innerMsg.ToByteArray()

			// 将游戏客户端消息转发给游戏服
			if err := bizServerConn.WriteMessage(msgType, innerMsgByteArray); nil != err {
				log.Error("消息转发失败，%+v", err)
			}
		}()
	}

	// 主动向游戏服务器发送断开连接消息
	userOfflineEvent := &base.UserOfflineEvent{
		GatewayServerId: ctx.GatewayServerId,
		SessionId:       ctx.SessionId,
		UserId:          ctx.userId,
	}

	byteArray, _ := json.Marshal(userOfflineEvent)
	// 利用广播通知其他服务器用户已经下线
	cluster.GetEtcdCli().Put(context.TODO(),
		"hero_story.go_server/publish/user_offline", string(byteArray))
}

func getBizServerFindStrategy(ctx *CmdContextImpl, ServerJobType base.ServerJobType) biz_server_finder.FindStrategy {
	if nil == ctx {
		return nil
	}

	if base.ServerJobTypeLogin == ServerJobType {
		return &biz_server_finder.RandomFindStrategy{}
	}

	if base.ServerJobTypeGame == ServerJobType {
		return &biz_server_finder.CompositeFindStrategy{

			Finder1: &biz_server_finder.IdFindStrategy{
				ServerId: ctx.GameServerId,
			},

			Finder2: &biz_server_finder.LeastLoadFindStrategy{
				WriteServerId: &ctx.GameServerId,
			},
		}
	}

	return nil
}
