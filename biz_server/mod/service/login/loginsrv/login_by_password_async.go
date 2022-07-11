package loginsrv

import (
	"fmt"
	"hero_story.go_server/biz_server/base"
	"hero_story.go_server/biz_server/mod/dao/user/user_dao"
	"hero_story.go_server/biz_server/mod/dao/user/user_data"
	"hero_story.go_server/biz_server/mod/dao/user/user_lock"
	"hero_story.go_server/comm/async_op"
	"time"
)

func LoginByPasswordAsync(userName string, password string) *base.AsyncBizResult {
	if len(userName) <= 0 ||
		len(password) <= 0 {
		return nil
	}

	bizResult := &base.AsyncBizResult{}

	async_op.Process(
		async_op.StrToBindId(userName),
		func() {
			user := user_dao.GetUserByName(userName)

			nowTime := time.Now().UnixMilli()

			if nil == user {
				user = &user_data.User{
					UserName:   userName,
					Password:   password,
					CreateTime: nowTime,
					HeroAvatar: "Hero_Hammer",
				}
			}

			// 是否有登出锁
			key := fmt.Sprintf("UserQuit_%d", user.UserId)

			// 如果存在登出锁，直接退出
			if user_lock.HasLock(key) {
				bizResult.SetReturnedObj(nil)
				return
			}

			// 更新最后登录时间
			user.LastLoginTime = nowTime
			user_dao.SaveOrUpdate(user)

			// 将用户添加到字典
			user_data.GetUserGroup().Add(user)
			bizResult.SetReturnedObj(user)
		},
		nil,
	)
	return bizResult
}
