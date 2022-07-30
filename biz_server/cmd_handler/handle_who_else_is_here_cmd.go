package cmd_handler

import (
	"google.golang.org/protobuf/types/dynamicpb"
	"hero_story.go_server/biz_server/base"
	"hero_story.go_server/biz_server/mod/dao/user/user_data"
	"hero_story.go_server/biz_server/msg"
	"hero_story.go_server/comm/log"
)

func init() {
	cmdHandlerMap[uint16(msg.MsgCode_WHO_ELSE_IS_HERE_CMD.Number())] = handleWhoElseIsHereCmd
}

func handleWhoElseIsHereCmd(ctx base.ICmdContext, _ *dynamicpb.Message) {
	if nil == ctx ||
		ctx.GetUserId() <= 0 {
		return
	}

	log.Info("收到'还有谁'的消息,userId = %d",
		ctx.GetUserId(),
	)

	whoElseIsHereResult := &msg.WhoElseIsHereResult{}

	userALL := user_data.GetUserGroup().GetUserAll()

	for _, user := range userALL {
		if nil == user {
			continue
		}

		userInfo := &msg.WhoElseIsHereResult_UserInfo{
			UserId:     uint32(user.UserId),
			UserName:   user.UserName,
			HeroAvatar: user.HeroAvatar,
		}

		if nil != user.MoveState {
			// 将数据中的移动状体 同步到 消息上的移动状态
			userInfo.MoveState = &msg.WhoElseIsHereResult_UserInfo_MoveState{
				FromPosX:  user.MoveState.FromPosX,
				FromPosY:  user.MoveState.FromPosY,
				ToPosX:    user.MoveState.ToPosX,
				ToPosY:    user.MoveState.ToPosY,
				StartTime: uint64(user.MoveState.StartTime),
			}
		}

		whoElseIsHereResult.UserInfo = append(
			whoElseIsHereResult.UserInfo,
			userInfo,
		)
	}

	ctx.Write(whoElseIsHereResult)
}
