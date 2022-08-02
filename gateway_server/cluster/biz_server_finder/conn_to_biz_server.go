package biz_server_finder

import (
	"fmt"
	"github.com/gorilla/websocket"
	"hero_story.go_server/biz_server/msg"
	"hero_story.go_server/comm/log"
	"hero_story.go_server/gateway_server/base"
	"sync"
)

var bizServerMap = &sync.Map{}

// 连接到游戏服务器
func connToBizServer(bizServerId int32, bizServerAddr string) {
	if bizServerId <= 0 || len(bizServerAddr) <= 0 {
		return
	}

	_, ok := bizServerMap.Load(bizServerId)

	if ok {
		return
	}

	// 创建到游戏服务器的连接
	newConn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/websocket", bizServerAddr), nil)

	if nil != err {
		log.Error("%+v", err)
		return
	}

	log.Info("已经连接到业务服务器,%s", bizServerAddr)

	//
	// 循环读取游戏服发来的消息,
	// 转发给客户端
	//
	go func() {
		_, ok := bizServerMap.Load(bizServerId)

		if ok {
			return
		}

		bizServerMap.Store(bizServerId, newConn)
		defer bizServerMap.Delete(bizServerId)

		for {
			// 读取从游戏服返回来的消息
			_, msgData, err := newConn.ReadMessage()

			if nil != err {
				log.Error("%+v", err)
			}

			innerMsg := &msg.InternalServerMsg{}
			innerMsg.FromByteArray(msgData)

			log.Info("从业务服务器返回结果, sessionId = %d, userId = %d", innerMsg.SessionId, innerMsg.UserId)

			// 这个是客户端到网关服务器的上下文对象,
			// 通过它来发送消息给客户端
			cmdCtx := base.GetCmdContextImplGroup().GetBySessionId(innerMsg.SessionId)

			if nil == cmdCtx {
				log.Error("未找到指令上下文")
				continue
			}

			if 0 != innerMsg.Disconnect {
				log.Info(
					"服务器强制断开玩家连接, sessionId = %d, userId = %d",
					innerMsg.SessionId,
					innerMsg.UserId,
				)
				cmdCtx.Disconnect()
				return
			}

			if cmdCtx.GetUserId() <= 0 &&
				innerMsg.UserId > 0 {
				cmdCtx.BindUserId(innerMsg.UserId)
			}

			cmdCtx.Write(innerMsg.MsgData)
		}
	}()
}
