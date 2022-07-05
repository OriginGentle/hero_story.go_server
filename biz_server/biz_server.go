package main

import (
	"github.com/gorilla/websocket"
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

	ctx := &myWebsocket.CmdContextImpl{
		Conn: conn,
	}

	ctx.LoopSendMsg()
	ctx.LoopReadMsg()
}
