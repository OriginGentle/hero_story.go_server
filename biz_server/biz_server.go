package main

import (
	"github.com/gorilla/websocket"
	"hero_story.go_server/biz_server/network/broadcaster"
	myWebsocket "hero_story.go_server/biz_server/network/websocket"
	"hero_story.go_server/comm/log"
	"net/http"
)

var upGrader = &websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var sessionId int32 = 0

// 启动业务服务器
func main() {
	log.Config("./log/biz_server.log")

	http.HandleFunc("/websocket", webSocketHandshake)
	_ = http.ListenAndServe("127.0.0.1:12345", nil)
}

func webSocketHandshake(w http.ResponseWriter, r *http.Request) {
	if nil == w || nil == r {
		return
	}

	conn, err := upGrader.Upgrade(w, r, nil)

	if nil != err {
		log.Error("Websocket upgrade error, %+v", err)
		return
	}

	defer func() {
		_ = conn.Close()
	}()

	log.Info("有新客户端连入,客户端地址,%+v", conn.LocalAddr())

	sessionId++

	ctx := &myWebsocket.CmdContextImpl{
		Conn:      conn,
		SessionId: sessionId,
	}

	// 将指令上下文添加到分组
	// 当断开连接时移除指令上下文
	broadcaster.AddCmdCtx(sessionId, ctx)
	defer broadcaster.RemoveBySessionId(sessionId)

	// 循环发送消息
	ctx.LoopSendMsg()
	// 循环读取消息
	ctx.LoopReadMsg()
}
