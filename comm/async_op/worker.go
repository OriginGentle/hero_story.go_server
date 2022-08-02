package async_op

import (
	"hero_story.go_server/comm/log"
	"hero_story.go_server/comm/main_thread"
)

type worker struct {
	taskQ chan func()
}

// process 处理异步过程
// 将异步操作放到队列里，并不立即执行
func (w *worker) process(asyncOp func(), continueWith func()) {
	if nil == asyncOp {
		log.Error("异步操作为空")
		return
	}
	if nil == w.taskQ {
		log.Error("任务队列尚未初始化")
		return
	}

	w.taskQ <- func() {
		asyncOp()

		if nil != continueWith {
			main_thread.Process(continueWith)
		}
	}
}

func (w *worker) loopExecTask() {
	if nil == w.taskQ {
		log.Error("任务队列尚未初始化")
		return
	}

	for {
		task := <-w.taskQ

		if nil == task {
			continue
		}

		func() {
			defer func() {
				if err := recover(); nil != err {
					log.Error("发生异常,%+v", err)
				}
			}()

			task()
		}()
	}
}
