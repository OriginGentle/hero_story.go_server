package broadcaster

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"hero_story.go_server/biz_server/base"
	"sync"
)

var innerMap = &sync.Map{}

// AddCmdCtx 添加指令上下文分组
func AddCmdCtx(sessionId string, cmdCtx base.ICmdContext) {
	if len(sessionId) <= 0 ||
		nil == cmdCtx {
		return
	}

	innerMap.Store(sessionId, cmdCtx)
}

// RemoveCmdCtxBySessionId 移除指令上下文分组
func RemoveCmdCtxBySessionId(sessionId string) {
	if len(sessionId) <= 0 {
		return
	}

	innerMap.Delete(sessionId)
}

// Broadcast 广播消息
func Broadcast(msgObj protoreflect.ProtoMessage) {
	if nil == msgObj {
		return
	}

	innerMap.Range(func(key interface{}, val interface{}) bool {
		if nil == key ||
			nil == val {
			return true
		}

		val.(base.ICmdContext).Write(msgObj)
		return true
	})
}
