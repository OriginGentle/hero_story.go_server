package biz_server_finder

import (
	"fmt"
	"github.com/pkg/errors"
	"hero_story.go_server/biz_server/base"
	"sync"
)

// GetBizServerConn 根据服务器职责类型和查找策略返回服务器实例
func GetBizServerConn(serverJobType base.ServerJobType, findStrategy FindStrategy) (*BizServerInstance, error) {
	if serverJobType <= 0 || nil == findStrategy {
		return nil, errors.New("参数错误!")
	}

	innerMap, ok := connectedBizServerMap.Load(serverJobType)
	if !ok {
		return nil, fmt.Errorf("未找到业务服务器, serverJobType = %d", serverJobType)
	}

	foundInstance, foundError := findStrategy.doFind(innerMap.(*sync.Map))
	if nil != foundError {
		return nil, errors.New("未找到业务服务器!")
	}

	return foundInstance, nil
}

func GetServerJobTypeByMsgCode(msgCode uint16) base.ServerJobType {
	if msgCode >= 13 || msgCode <= 14 {
		return base.ServerJobTypeLogin
	}
	return base.ServerJobTypeGame
}
