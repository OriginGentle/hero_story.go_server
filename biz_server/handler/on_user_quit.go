package handler

import (
	"hero_story.go_server/biz_server/base"
	"hero_story.go_server/biz_server/msg"
	"hero_story.go_server/biz_server/network/broadcaster"
	"hero_story.go_server/comm/log"
)

// OnUserQuit 用户退游戏时执行
func OnUserQuit(ctx base.ICmdContext) {
	if nil == ctx {
		return
	}

	log.Info("用户离线,userId = %d", ctx.GetUserId())

	broadcaster.Broadcast(&msg.UserQuitResult{
		QuitUserId: uint32(ctx.GetUserId()),
	})
}
