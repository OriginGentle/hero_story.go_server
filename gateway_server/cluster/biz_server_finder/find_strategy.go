package biz_server_finder

import (
	"fmt"
	"github.com/pkg/errors"
	"hero_story.go_server/comm/log"
	"math/rand"
	"sync"
)

// FindStrategy 查找业务服务器策略接口
type FindStrategy interface {
	// 给定业务服务器实例字典，执行查找过程
	doFind(bizServerInstanceMap *sync.Map) (*BizServerInstance, error)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// RandomFindStrategy 随机选择策略
type RandomFindStrategy struct {
}

func (finder *RandomFindStrategy) doFind(bizServerInstanceMap *sync.Map) (*BizServerInstance, error) {
	if nil == bizServerInstanceMap {
		return nil, errors.New("业务服务器字典为空!")
	}

	// 用于收集服务器实例的数组
	varArray := make([]interface{}, 1)
	bizServerInstanceMap.Range(func(_, val interface{}) bool {
		varArray = append(varArray, val)
		return true
	})

	count := len(varArray)
	if count <= 0 {
		return nil, errors.New("业务服务器实例数量为 0!")
	}

	randIndex := rand.Int31n(int32(count))
	findVal := varArray[randIndex]
	return findVal.(*BizServerInstance), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// IdFindStrategy 根据服务器Id查找
type IdFindStrategy struct {
	ServerId int32
}

func (finder *IdFindStrategy) doFind(bizServerInstanceMap *sync.Map) (*BizServerInstance, error) {
	if nil == bizServerInstanceMap {
		return nil, errors.New("业务服务器字典为空!")
	}

	val, ok := bizServerInstanceMap.Load(finder.ServerId)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Id = %d 的服务器不存在", finder.ServerId))
	}

	log.Info("IdFindStrategy 找到了 Id = %d 的业务服务器", finder.ServerId)
	return val.(*BizServerInstance), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// LeastLoadFindStrategy 根据服务器负载数查找，寻找负载最小业务服务器
type LeastLoadFindStrategy struct {
	WriteServerId *int32
}

func (finder *LeastLoadFindStrategy) doFind(bizServerInstanceMap *sync.Map) (*BizServerInstance, error) {
	if nil == bizServerInstanceMap {
		return nil, errors.New("业务服务器字典为空!")
	}

	var findBizServerInstance *BizServerInstance = nil
	var minLoadCount int32 = 999999

	bizServerInstanceMap.Range(func(_, val any) bool {
		if nil == val {
			return true
		}

		currBizServerInstance := val.(*BizServerInstance)
		if currBizServerInstance.LoadCount < minLoadCount {
			findBizServerInstance = currBizServerInstance
			minLoadCount = currBizServerInstance.LoadCount
		}
		return true
	})

	if nil == findBizServerInstance {
		return nil, errors.New("未找到负载最小的业务服务器!")
	}

	log.Info("LeastLoadFindStrategy 找到了Id = %d 的业务服务器", findBizServerInstance.ServerId)
	*finder.WriteServerId = findBizServerInstance.ServerId
	return findBizServerInstance, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type CompositeFindStrategy struct {
	Finder1 FindStrategy
	Finder2 FindStrategy
}

func (finder *CompositeFindStrategy) doFind(bizServerInstanceMap *sync.Map) (*BizServerInstance, error) {
	if nil == bizServerInstanceMap {
		return nil, errors.New("业务服务器字典为空!")
	}

	var bizServerInstance *BizServerInstance
	// 先用策略1进行查找
	if nil != finder.Finder1 {
		bizServerInstance, _ = finder.Finder1.doFind(bizServerInstanceMap)
	}

	if nil != bizServerInstance {
		return bizServerInstance, nil
	}

	// 如果策略1找不到，就使用策略2进行查找
	if nil != finder.Finder2 {
		return finder.Finder2.doFind(bizServerInstanceMap)
	}

	// 如果策略2也找不到，直接返回错误
	return nil, errors.New("CompositeFindStrategy 未找到业务服务器实例!")
}
