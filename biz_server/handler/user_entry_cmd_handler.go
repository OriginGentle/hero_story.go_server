package handler

import (
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/types/dynamicpb"
	"hero_story.go_server/biz_server/msg"
)

func init() {
	cmdHandlerMap[uint16(msg.MsgCode_USER_ENTRY_CMD.Number())] = userEntryCmdHandler
}

func userEntryCmdHandler(conn *websocket.Conn, pbMsgObj *dynamicpb.Message) {

}
