package handler

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
	"hero_story.go_server/biz_server/msg"
	"hero_story.go_server/comm/log"
)

func init() {
	cmdHandlerMap[uint16(msg.MsgCode_USER_LOGIN_CMD.Number())] = userLoginCmdHandler
}

func userLoginCmdHandler(ctx ICmdContext, pbMsgObj *dynamicpb.Message) {
	if nil == ctx ||
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

	ctx.Write(userLoginResult)
}
