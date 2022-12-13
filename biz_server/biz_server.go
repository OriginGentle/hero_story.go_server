package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	clientv3 "go.etcd.io/etcd/client/v3"
	"hero_story.go_server/biz_server/base"
	mywebsocket "hero_story.go_server/biz_server/network/websocket"
	"hero_story.go_server/comm/log"
	"net/http"
	"strings"
	"time"
)

var pServerId *int
var pBindHost *string
var pBindPort *int
var pSjtArray *string
var pEtcdEndpointArray *string

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

	pServerId = flag.Int("server_id", 0, "业务服务器 Id")
	pBindHost = flag.String("bind_host", "127.0.0.1", "绑定主机地址")
	pBindPort = flag.Int("bind_port", 12345, "绑定端口号")
	pSjtArray = flag.String("server_job_type_array", "", "服务器职责类型数组")
	pEtcdEndpointArray = flag.String("etcd_endpoint_array", "127.0.0.1:2379", "Etcd 节点地址数组")
	flag.Parse() // 解析命令行参数

	sjtArray := base.StringToServerJobTypeArray(*pSjtArray)
	etcdEndpointArray := strings.Split(*pEtcdEndpointArray, ",")

	registerTheServer(etcdEndpointArray, *pServerId, *pBindHost, *pBindPort, sjtArray)

	log.Info("启动业务服务器, serverId = %d, serverAddr = %s:%d, serverJobTypeArray = %s",
		*pServerId, *pBindHost, *pBindPort, *pSjtArray)

	http.HandleFunc("/websocket", webSocketHandshake)
	_ = http.ListenAndServe(fmt.Sprintf("%s:%d", *pBindHost, *pBindPort), nil)
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

	log.Info("有新客户端连入")

	sessionId += 1

	myConn := &mywebsocket.GatewayServerConn{
		WsConn: conn,
	}

	myConn.LoopSendMsg()
	myConn.LoopReadMsg()
}

// 注册本服务器
func registerTheServer(etcdEndpointArray []string, serverId int, bindHost string,
	bindPort int, sjtArray []base.ServerJobType) {

	etcdCli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdEndpointArray,
		DialTimeout: 5 * time.Second,
	})

	if nil != err {
		log.Error("%+v", err)
		return
	}

	reportData := &base.BizServerData{
		ServerId:   int32(serverId),
		ServerAddr: fmt.Sprintf("%s:%d", bindHost, bindPort),
		SjtArray:   sjtArray,
	}

	go func() {
		grantResp, _ := etcdCli.Grant(context.TODO(), 10)

		for {
			time.Sleep(5 * time.Second)

			// 更新负载数 ( 总人数 )
			reportData.LoadCount = base.GetLoadStat().GetTotalCount()

			leaseKeepLiveResp, _ := etcdCli.KeepAliveOnce(context.TODO(), grantResp.ID)
			byteArray, _ := json.Marshal(reportData)

			_, _ = etcdCli.Put(
				context.TODO(),
				fmt.Sprintf("hero_story.go_server/biz_server_%d", serverId), // hero_story.go_server/biz_server_1001
				string(byteArray),
				clientv3.WithLease(leaseKeepLiveResp.ID),
			)
		}
	}()

	go func() {
		watchChan := etcdCli.Watch(context.TODO(), "hero_story.go_server/publish/user_offline", clientv3.WithPrefix())

		for resp := range watchChan {
			for _, event := range resp.Events {
				switch event.Type {
				case 0: // PUT
					strVal := string(event.Kv.Value)
					log.Info("收到玩家下线通知, " + strVal)

					userOfflineEvent := &base.UserOfflineEvent{}
					_ = json.Unmarshal([]byte(strVal), userOfflineEvent)

					base.GetLoadStat().DeleteUserId(
						userOfflineEvent.GatewayServerId,
						userOfflineEvent.UserId,
					)

				case 1: // DELETE
				}
			}
		}
	}()
}
