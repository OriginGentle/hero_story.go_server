package loginsrv

import (
	"hero_story.go_server/biz_server/mod/dao/user/userdao"
	"hero_story.go_server/biz_server/mod/dao/user/userdata"
	"time"
)

func LoginByPasswordAsync(userName string, password string) *userdata.User {
	if len(userName) <= 0 ||
		len(password) <= 0 {
		return nil
	}

	user := userdao.GetUserByName(userName)

	nowTime := time.Now().UnixMilli()

	if nil == user {
		user = &userdata.User{
			UserName:   userName,
			Password:   password,
			CreateTime: nowTime,
			HeroAvatar: "Hero_Hammer",
		}
	}

	user.LastLoginTime = nowTime
	userdao.SaveOrUpdate(user)

	return user
}
