package login_srv

import (
	"fmt"
	"hero_story.go_server/biz_server/base"
	user_dao2 "hero_story.go_server/biz_server/mod/user/user_dao"
	user_data2 "hero_story.go_server/biz_server/mod/user/user_data"
	"hero_story.go_server/biz_server/mod/user/user_lock"
	"hero_story.go_server/comm/async_op"
	"time"
)

// LoginByPasswordAsync 根据用户名和密码进行登录
// 返回一个异步的业务结果
func LoginByPasswordAsync(userName string, password string) *base.AsyncBizResult {
	if len(userName) <= 0 || len(password) <= 0 {
		return nil
	}

	bizResult := &base.AsyncBizResult{}

	async_op.Process(
		async_op.StrToBindId(userName),
		func() {
			user := user_dao2.GetUserByName(userName)

			nowTime := time.Now().UnixMilli()

			if nil == user {
				user = &user_data2.User{
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
			user_dao2.SaveOrUpdate(user)

			// 将用户添加到字典
			user_data2.GetUserGroup().Add(user)
			bizResult.SetReturnedObj(user)
		},
		nil,
	)
	return bizResult
}
