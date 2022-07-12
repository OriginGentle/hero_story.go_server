package log

import (
	"fmt"
	"log"
	"sync"
)

var writer *dailyFileWriter
var infoLogger, errorLogger *log.Logger

func Config(outputFileName string) {
	if len(outputFileName) <= 0 {
		panic("输出文件名为空")
	}

	writer = &dailyFileWriter{
		fileName:       outputFileName,
		lastYearDay:    -1,
		fileSwitchLock: &sync.Mutex{},
	}

	infoLogger = log.New(
		writer, "[ INFO ]",
		log.Ltime|log.Lmicroseconds|log.Lshortfile,
	)

	errorLogger = log.New(
		writer, "[ ERROR ]",
		log.Ltime|log.Lmicroseconds|log.Lshortfile,
	)
}

func Info(format string, valArray ...interface{}) {
	_ = infoLogger.Output(
		2,
		fmt.Sprintf(format, valArray...),
	)
}

func Error(format string, valArray ...interface{}) {
	_ = errorLogger.Output(
		2,
		fmt.Sprintf(format, valArray...),
	)
}
