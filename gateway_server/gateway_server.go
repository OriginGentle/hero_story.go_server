package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"hero_story.go_server/comm/log"
	"hero_story.go_server/gateway_server/base"
	mywebsocket "hero_story.go_server/gateway_server/network/websocket"
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

// 启动网关服务器
func main() {

	fmt.Println("启动网关服务器")

	log.Config("./log/gateway_server.log")

	http.HandleFunc("/websocket", webSocketHandshake)
	_ = http.ListenAndServe("127.0.0.1:54321", nil)
}

func webSocketHandshake(w http.ResponseWriter, r *http.Request) {
	if nil == w ||
		nil == r {
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

	log.Info("有新客户端连入")

	sessionId += 1

	cmdCtx := &mywebsocket.CmdContextImpl{
		Conn:      conn,
		SessionId: sessionId,
	}

	base.GetCmdContextImplGroup().Add(cmdCtx.SessionId, cmdCtx)
	defer base.GetCmdContextImplGroup().RemoveBySessionId(cmdCtx.SessionId)

	cmdCtx.LoopSendMsg()
	cmdCtx.LoopReadMsg()
}
