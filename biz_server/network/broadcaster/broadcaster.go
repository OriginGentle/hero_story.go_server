package broadcaster

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"hero_story.go_server/biz_server/base"
	"sync"
)

var innerMap = &sync.Map{}

func AddCmdCtx(sessionUId string, cmdCtx base.ICmdContext) {
	if len(sessionUId) <= 0 ||
		nil == cmdCtx {
		return
	}

	innerMap.Store(sessionUId, cmdCtx)
}

func RemoveBySessionId(sessionUId string) {
	if len(sessionUId) <= 0 {
		return
	}

	innerMap.Delete(sessionUId)
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
