package cmd_handler

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
	"hero_story.go_server/biz_server/base"
	"hero_story.go_server/biz_server/mod/dao/user/user_data"
	"hero_story.go_server/biz_server/mod/dao/user/user_lso"
	"hero_story.go_server/biz_server/msg"
	"hero_story.go_server/biz_server/network/broadcaster"
	"hero_story.go_server/comm/lazy_save"
	"hero_story.go_server/comm/log"
)

func init() {
	cmdHandlerMap[uint16(msg.MsgCode_USER_ATTK_CMD.Number())] = handleUserAttkCmd
}

func handleUserAttkCmd(ctx base.ICmdContext, pbMsgObj *dynamicpb.Message) {
	if nil == ctx ||
		nil == pbMsgObj {
		return
	}

	log.Info("收到用户攻击消息,userId = %d", ctx.GetUserId())

	userAttkCmd := &msg.UserAttkCmd{}

	pbMsgObj.Range(func(f protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		userAttkCmd.ProtoReflect().Set(f, v)
		return true
	})

	userAttkResult := &msg.UserAttkResult{
		AttkUserId:   uint32(ctx.GetUserId()),
		TargetUserId: userAttkCmd.TargetUserId,
	}

	broadcaster.Broadcast(userAttkResult)

	user := user_data.GetUserGroup().GetByUserId(int64(userAttkCmd.TargetUserId))

	if nil == user {
		return
	}

	var subtractHp int32 = 10
	user.CurrHp -= subtractHp

	userSubtractHpResult := &msg.UserSubtractHpResult{
		SubtractHp:   uint32(subtractHp),
		TargetUserId: userAttkCmd.TargetUserId,
	}

	broadcaster.Broadcast(userSubtractHpResult)

	lso := user_lso.GetUserLso(user)

	// 执行延时保存
	lazy_save.SaveOrUpdate(lso)
}
