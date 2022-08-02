package async_op

import (
	"sync"
)

// 工人数组
var workerArray = [2048]*worker{}

// 初始化工人用的锁
var initWorkerLocker = &sync.Mutex{}

// Process 处理异步过程
// asyncOp 异步函数, 将被放到一个新协程里去执行...
// continueWith 则是回到主线程继续执行的函数
func Process(bindId int, asyncOp func(), continueWith func()) {
	if nil == asyncOp {
		return
	}

	currWorker := getCurrWorker(bindId)

	if nil != currWorker {
		currWorker.process(asyncOp, continueWith)
	}
}

// 根据 bindId 获取一个工人
func getCurrWorker(bindId int) *worker {
	if bindId < 0 {
		bindId = -bindId
	}

	workerIndex := bindId % len(workerArray)
	currWorker := workerArray[workerIndex]

	if nil != currWorker {
		return currWorker
	}

	initWorkerLocker.Lock()
	defer initWorkerLocker.Unlock()

	if nil != currWorker {
		return currWorker
	}

	currWorker = &worker{
		taskQ: make(chan func(), 2048),
	}

	workerArray[workerIndex] = currWorker
	go currWorker.loopExecTask()

	return currWorker
}
