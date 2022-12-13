package user_lso

import (
	"hero_story.go_server/biz_server/mod/user/user_data"
)

func GetUserLso(user *user_data.User) *UserLso {
	if nil == user {
		return nil
	}
	exitComp, _ := user.GetComponentMap().Load("UserLso")

	if nil != exitComp {
		return exitComp.(*UserLso)
	}

	exitComp = &UserLso{
		User: user,
	}

	exitComp, _ = user.GetComponentMap().LoadOrStore("UserLso", exitComp)

	return exitComp.(*UserLso)
}
