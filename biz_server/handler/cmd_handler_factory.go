package handler

import (
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/types/dynamicpb"
)

// CmdHandlerFunc 自定义的消息处理函数
type CmdHandlerFunc func(conn *websocket.Conn, pbMsgObj *dynamicpb.Message)

// 消息处理器字典	key = msgCode val = CmdHandlerFunc
var cmdHandlerMap = make(map[uint16]CmdHandlerFunc)

func CreateCmdHandler(msgCode uint16) CmdHandlerFunc {
	return cmdHandlerMap[msgCode]
}
