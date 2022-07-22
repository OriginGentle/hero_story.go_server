package handler

import (
	"fmt"
	"hero_story.go_server/biz_server/base"
	"hero_story.go_server/biz_server/mod/dao/user/user_data"
	"hero_story.go_server/biz_server/mod/dao/user/user_lock"
	"hero_story.go_server/biz_server/mod/dao/user/user_lso"
	"hero_story.go_server/biz_server/msg"
	"hero_story.go_server/biz_server/network/broadcaster"
	"hero_story.go_server/comm/lazy_save"
	"hero_story.go_server/comm/log"
)

// OnUserQuit 用户退游戏时执行
func OnUserQuit(ctx base.ICmdContext) {
	if nil == ctx {
		return
	}

	log.Info("用户离线,userId = %d", ctx.GetUserId())

	// 加登出锁
	key := fmt.Sprintf("UserQuit_%d", ctx.GetUserId())
	user_lock.TryLock(key)

	// 广播用户退出消息
	broadcaster.Broadcast(&msg.UserQuitResult{
		QuitUserId: uint32(ctx.GetUserId()),
	})

	user := user_data.GetUserGroup().GetByUserId(ctx.GetUserId())

	if nil == user {
		return
	}

	userLso := user_lso.GetUserLso(user)
	lazy_save.Discard(userLso)

	log.Info("用户离线，立即保存数据！userId = %d", ctx.GetUserId())

	userLso.SaveOrUpdate(func() {
		user_lock.UnLock(key)
	})
}
