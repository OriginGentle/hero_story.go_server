package base

import (
	"hero_story.go_server/comm/main_thread"
	"sync/atomic"
)

type AsyncBizResult struct {
	// 已返回对象
	returnedObj interface{}
	// 完成回调函数
	completeFunc func()
	// 是否已有返回值
	hasReturnedObj int32
	// 是否已有回调函数
	hasCompleteFunc int32
	// 是否已经调用过完成函数
	completeFuncHasAlreadyBeenCalled int32
}

// GetReturnedObj 获取返回值
func (bizResult *AsyncBizResult) GetReturnedObj() interface{} {
	return bizResult.returnedObj
}

// SetReturnedObj 设置返回值
func (bizResult *AsyncBizResult) SetReturnedObj(val interface{}) {
	if atomic.CompareAndSwapInt32(&bizResult.hasReturnedObj, 0, 1) {
		bizResult.returnedObj = val
		bizResult.doComplete()
	}
}

// OnComplete 完成回调函数
func (bizResult *AsyncBizResult) OnComplete(val func()) {
	if atomic.CompareAndSwapInt32(&bizResult.hasCompleteFunc, 0, 1) {
		bizResult.completeFunc = val

		if 1 == bizResult.hasReturnedObj {
			bizResult.doComplete()
		}
	}
}

// 执行完成回调函数
func (bizResult *AsyncBizResult) doComplete() {
	if nil == bizResult.completeFunc {
		return
	}

	// 通过CAS原语比较标记值
	// 设置成功执行完成调用函数
	if atomic.CompareAndSwapInt32(&bizResult.completeFuncHasAlreadyBeenCalled, 0, 1) {
		// 回到主线程去执行
		main_thread.Process(bizResult.completeFunc)
	}
}
