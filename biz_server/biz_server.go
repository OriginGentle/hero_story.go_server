package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	clientv3 "go.etcd.io/etcd/client/v3"
	mywebsocket "hero_story.go_server/biz_server/network/websocket"
	"hero_story.go_server/comm/log"
	"net/http"
	"time"
)

const SERVER_ADDR = "127.0.0.1:12345"

type bizServerData struct {
	ServerId   int32  `json:"serverId"`
	ServerAddr string `json:"serverAddr"`
}

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
	fmt.Printf("start bizServer")
	log.Config("./log/biz_server.log")

	registerTheServer()

	http.HandleFunc("/websocket", webSocketHandshake)
	_ = http.ListenAndServe(SERVER_ADDR, nil)
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

	//ctx := &myWebsocket.CmdContextImpl{
	//	Conn:      conn,
	//	SessionId: sessionId,
	//}
	//
	//// 将指令上下文添加到分组
	//// 当断开连接时移除指令上下文
	//broadcaster.AddCmdCtx(sessionId, ctx)
	//defer broadcaster.RemoveBySessionId(sessionId)
	//
	//// 循环发送消息
	//ctx.LoopSendMsg()
	//// 循环读取消息
	//ctx.LoopReadMsg()

	myConn := &mywebsocket.GatewayServerConn{
		WsConn: conn,
	}

	myConn.LoopSendMsg()
	myConn.LoopReadMsg()
}

// 注册本服务器信息
func registerTheServer() {
	etcdCli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})

	if nil != err {
		log.Error("etcd客户端创建异常，%+v", err)
		return
	}

	reportData := &bizServerData{
		ServerId:   1001,
		ServerAddr: SERVER_ADDR,
	}

	go func() {
		for {
			time.Sleep(5 * time.Second)

			byteArray, _ := json.Marshal(reportData)
			// 申请租约，10s过期
			grantResp, _ := etcdCli.Grant(context.TODO(), 10)
			_, _ = etcdCli.Put(
				context.TODO(),
				"hero_story.go_server/biz_server_1001",
				string(byteArray),
				clientv3.WithLease(grantResp.ID),
			)
		}
	}()
}
