package handler

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
	"hero_story.go_server/biz_server/base"
	"hero_story.go_server/biz_server/mod/dao/user/user_data"
	"hero_story.go_server/biz_server/mod/service/login/login_srv"
	"hero_story.go_server/biz_server/msg"
	"hero_story.go_server/comm/log"
)

func init() {
	cmdHandlerMap[uint16(msg.MsgCode_USER_LOGIN_CMD.Number())] = handleUserLoginCmd
}

func handleUserLoginCmd(ctx base.ICmdContext, pbMsgObj *dynamicpb.Message) {
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

	bizResult := login_srv.LoginByPasswordAsync(userLoginCmd.GetUserName(), userLoginCmd.GetPassword())

	if nil == bizResult {
		log.Error("业务结果返回空值,userName = %s",
			userLoginCmd.GetUserName(),
		)
		return
	}

	bizResult.OnComplete(func() {
		returnedObj := bizResult.GetReturnedObj()

		if nil == returnedObj {
			log.Error("用户不存在,userName = %s",
				userLoginCmd.GetUserName(),
			)
			return
		}

		user := returnedObj.(*user_data.User)

		userLoginResult := &msg.UserLoginResult{
			UserId:     uint32(user.UserId),
			UserName:   user.UserName,
			HeroAvatar: user.HeroAvatar,
		}

		ctx.BindUserId(user.UserId)
		ctx.Write(userLoginResult)
	})

}
