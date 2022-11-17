package cluster

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

var etcdCli *clientv3.Client

func InitEtcdCli(etcdEndpointArray []string) {
	if nil == etcdEndpointArray || len(etcdEndpointArray) <= 0 {
		panic("初始化 Etcd 失败 <==> 参数数组为空")
	}

	var err error
	etcdCli, err = clientv3.New(clientv3.Config{
		Endpoints:   etcdEndpointArray,
		DialTimeout: 5 * time.Second,
	})

	if nil != err {
		panic(err.Error())
	}
}

func GetEtcdCli() *clientv3.Client {
	return etcdCli
}
