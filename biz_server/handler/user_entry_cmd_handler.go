package handler

import (
	"google.golang.org/protobuf/types/dynamicpb"
	"hero_story.go_server/biz_server/base"
	"hero_story.go_server/biz_server/mod/dao/user/user_data"
	"hero_story.go_server/biz_server/msg"
	"hero_story.go_server/biz_server/network/broadcaster"
	"hero_story.go_server/comm/log"
)

func init() {
	cmdHandlerMap[uint16(msg.MsgCode_USER_ENTRY_CMD.Number())] = handleUserEntryCmd
}

func handleUserEntryCmd(ctx base.ICmdContext, pbMsgObj *dynamicpb.Message) {
	if nil == ctx ||
		nil == pbMsgObj {
		return
	}

	log.Info("收到用户入场消息,userId = %d",
		ctx.GetUserId(),
	)

	user := user_data.GetUserGroup().GetByUserId(ctx.GetUserId())

	if nil == user {
		log.Error("为找到用户数据，userId = %d",
			ctx.GetUserId(),
		)
		return
	}

	userEntryResult := &msg.UserEntryResult{
		UserId:     uint32(ctx.GetUserId()),
		UserName:   user.UserName,
		HeroAvatar: user.HeroAvatar,
	}

	broadcaster.Broadcast(userEntryResult)
}
