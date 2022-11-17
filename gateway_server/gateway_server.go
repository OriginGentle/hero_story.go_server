package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"hero_story.go_server/comm/log"
	"hero_story.go_server/gateway_server/base"
	"hero_story.go_server/gateway_server/cluster"
	"hero_story.go_server/gateway_server/cluster/biz_server_finder"
	myWebsocket "hero_story.go_server/gateway_server/network/websocket"
	"net/http"
	"strings"
)

var upGrader = &websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var pServerId *int
var pBindHost *string
var pBindPort *int
var pEtcdEndpointArray *string
var sessionId int32 = 0

// 启动网关服务器
func main() {
	fmt.Println("启动网关服务器")
	log.Config("./log/gateway_server.log")

	pServerId = flag.Int("server_id", 0, "业务服务器 id")
	pBindHost = flag.String("bind_host", "127.0.0.1", "绑定主机地址")
	pBindPort = flag.Int("bind_port", 54321, "绑定端口号")
	pEtcdEndpointArray = flag.String("etcd_endpoint_array", "127.0.0.1:2379", "Etcd 节点地址数组")
	flag.Parse()

	etcdEndpointArray := strings.Split(*pEtcdEndpointArray, ",")
	cluster.InitEtcdCli(etcdEndpointArray)
	biz_server_finder.FindNewBizServer()

	log.Info("启动网关服务器,serverId = %d, serverAddr = %s:%d", *pServerId, *pBindHost, *pBindPort)
	http.HandleFunc("/websocket", webSocketHandshake)
	_ = http.ListenAndServe(fmt.Sprintf("%s:%d", *pBindHost, *pBindPort), nil)
}

func webSocketHandshake(writer http.ResponseWriter, request *http.Request) {
	if nil == writer || nil == request {
		return
	}

	conn, err := upGrader.Upgrade(writer, request, nil)
	if nil != err {
		log.Error("Websocket upgrade error, %v", err)
		return
	}

	defer func() {
		_ = conn.Close()
	}()

	log.Info("有新客户端连入")
	sessionId += 1
	cmdCtx := &myWebsocket.CmdContextImpl{
		Conn:         conn,
		SessionId:    sessionId,
		GameServerId: int32(*pServerId),
	}

	base.GetCmdContextImplGroup().Add(cmdCtx.SessionId, cmdCtx)
	defer base.GetCmdContextImplGroup().RemoveBySessionId(cmdCtx.SessionId)

	cmdCtx.LoopSendMsg()
	cmdCtx.LoopReadMsg()
}
