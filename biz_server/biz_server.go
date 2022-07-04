package main

import (
	"encoding/binary"
	"github.com/gorilla/websocket"
	"hero_story.go_server/biz_server/handler"
	"hero_story.go_server/biz_server/msg"
	"hero_story.go_server/comm/log"
	"hero_story.go_server/comm/main_thread"
	"net/http"
)

var upGrader = &websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	log.Config("./log/biz_server.log")

	http.HandleFunc("/websocket", webSocketHandshake)
	_ = http.ListenAndServe("127.0.0.1:12345", nil)
}

func webSocketHandshake(w http.ResponseWriter, r *http.Request) {
	if nil == w || nil == r {
		return
	}

	conn, err := upGrader.Upgrade(w, r, nil)

	if nil != err {
		log.Error("Websocket upgrade error, %+v", err)
		return
	}

	defer func() {
		_ = conn.Close()
	}()

	log.Info("有新客户端联入,客户端地址,%+v", conn.LocalAddr())

	for {
		_, msgData, err := conn.ReadMessage()

		if nil != err {
			log.Error("%+v", err)
			break
		}

		log.Info("%+v", msgData)

		msgCode := binary.BigEndian.Uint16(msgData[2:4])
		newMsgX, err := msg.Decode(msgData[4:], int16(msgCode))

		if nil != err {
			log.Error("消息解码错误,msgCode = %d ,err = %+v ",
				msgCode, err,
			)
			continue
		}

		log.Info("收到客户端消息,msgCode = %d,msgName = %s",
			msgCode, newMsgX.Descriptor().Name(),
		)

		cmdHandler := handler.CreateCmdHandler(msgCode)

		if nil == cmdHandler {
			log.Error("未找到指令处理器,msgCode = %d",
				msgCode,
			)
			continue
		}

		main_thread.Process(func() {
			cmdHandler(conn, newMsgX)
		})
	}
}
