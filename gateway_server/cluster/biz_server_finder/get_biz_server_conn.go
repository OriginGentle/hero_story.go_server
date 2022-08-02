package biz_server_finder

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

func GetBizServerConn(bizServerId int32) (*websocket.Conn, error) {
	if bizServerId <= 0 {
		return nil, errors.New("参数错误")
	}

	serverConn, ok := bizServerMap.Load(bizServerId)

	if !ok {
		return nil, fmt.Errorf(
			"未找到业务服务器, bizServerId = %d",
			bizServerId,
		)
	}

	return serverConn.(*websocket.Conn), nil
}
