package base

import "sync"

type ClientCmdContextGroup struct {
	innerMap *sync.Map
}

var cmdContextImplGroupInstance = &ClientCmdContextGroup{
	innerMap: &sync.Map{},
}

// GetCmdContextImplGroup 获取单例对象
func GetCmdContextImplGroup() *ClientCmdContextGroup {
	return cmdContextImplGroupInstance
}

// Add 添加客户端指令上下文
func (group *ClientCmdContextGroup) Add(sessionId int32, cmdCtx ClientCmdContext) {
	if nil == cmdCtx {
		return
	}

	group.innerMap.Store(sessionId, cmdCtx)
}

// RemoveBySessionId 根据会话 Id 移除客户端指令上下文
func (group *ClientCmdContextGroup) RemoveBySessionId(sessionId int32) {
	if sessionId <= 0 {
		return
	}

	group.innerMap.Delete(sessionId)
}

// GetBySessionId 根据会话 Id 获取客户端指令上下文
func (group *ClientCmdContextGroup) GetBySessionId(sessionId int32) ClientCmdContext {
	if sessionId <= 0 {
		return nil
	}

	cmdCtx, ok := group.innerMap.Load(sessionId)

	if !ok {
		return nil
	}

	return cmdCtx.(ClientCmdContext)
}
