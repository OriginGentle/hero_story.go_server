package biz_server_finder

import (
	"github.com/gorilla/websocket"
	"hero_story.go_server/biz_server/msg"
	"hero_story.go_server/comm/log"
	"hero_story.go_server/gateway_server/base"
	"sync"
)

const bizServerUrl = "ws://127.0.0.1:12345/websocket"

var bizServerConn *websocket.Conn
var locker = &sync.Mutex{}

func GetBizServerConn() (*websocket.Conn, error) {
	if nil != bizServerConn {
		return bizServerConn, nil
	}

	locker.Lock()
	defer locker.Unlock()

	if nil != bizServerConn {
		return bizServerConn, nil
	}

	newConn, _, err := websocket.DefaultDialer.Dial(bizServerUrl, nil)

	if nil != err {
		log.Error("业务服务器连接建立异常，%+v", err)
		return nil, err
	}

	bizServerConn = newConn

	// 循环读取从游戏服返回的消息
	// 转发给客户端
	go func() {
		for {

			// 读取从游戏服务器返回的消息
			_, msgData, err := bizServerConn.ReadMessage()

			if nil != err {
				log.Error("消息读取异常，%+v", err)
			}

			innerMsg := &msg.InternalServerMsg{}
			innerMsg.FromByteArray(msgData)

			log.Info("从业务服务器返回结果,sessionId = %d,userId = %d", innerMsg.SessionId, innerMsg.UserId)

			// 根据SessionId,获取客户端到网关服务器的上下文
			// 通过它来向客户端发送消息
			cmdCtx := base.GetCmdContextImplGroup().GetBySessionId(innerMsg.SessionId)

			if nil == cmdCtx {
				log.Error("未找到上下文对象")
				continue
			}

			if cmdCtx.GetUserId() <= 0 &&
				innerMsg.UserId > 0 {
				cmdCtx.BindUserId(innerMsg.UserId)
			}

			cmdCtx.Write(innerMsg.MsgData)
		}
	}()

	return bizServerConn, nil
}
