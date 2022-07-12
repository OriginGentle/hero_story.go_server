package async_op

import (
	"sync"
)

var workerArray = [2048]*worker{}

var initWorkerLocker = &sync.Mutex{}

func Process(bindId int, asyncOp func(), continueWith func()) {
	if nil == asyncOp {
		return
	}

	currWorker := getCurrWorker(bindId)

	if nil != currWorker {
		currWorker.process(asyncOp, continueWith)
	}
}

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
