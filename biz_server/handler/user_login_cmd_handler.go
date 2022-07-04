package handler

import (
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
	"hero_story.go_server/biz_server/msg"
	"hero_story.go_server/comm/log"
)

func init() {
	cmdHandlerMap[uint16(msg.MsgCode_USER_LOGIN_CMD.Number())] = userLoginCmdHandler
}

func userLoginCmdHandler(conn *websocket.Conn, pbMsgObj *dynamicpb.Message) {
	if nil == conn ||
		nil == pbMsgObj {
		return
	}

	userLoginCmd := &msg.UserLoginCmd{}

	pbMsgObj.Range(func(f protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		userLoginCmd.ProtoReflect().Set(f, v)
		return true
	})

	log.Info("收到用户登录消息！userName = %s,password = %s",
		userLoginCmd.GetUserName(),
		userLoginCmd.GetPassword(),
	)

	userLoginResult := &msg.UserLoginResult{
		UserId:     1,
		UserName:   userLoginCmd.UserName,
		HeroAvatar: "Hero_Shaman",
	}

	byteArray, err := msg.Encode(userLoginResult)

	if nil != err {
		return
	}

	if err := conn.WriteMessage(websocket.BinaryMessage, byteArray); nil != err {
		log.Error("%+v", err)
	}
}
