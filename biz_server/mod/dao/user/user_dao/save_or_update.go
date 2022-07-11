package user_dao

import (
	"hero_story.go_server/biz_server/base"
	"hero_story.go_server/biz_server/mod/dao/user/user_data"
	"hero_story.go_server/comm/log"
)

const sqlSaveOrUpdate = `
insert into t_user(
	user_name,password,hero_avatar,create_time,last_login_time
) value (
	?,?,?,?,?
) on duplicate key update last_login_time = ?
`

func SaveOrUpdate(user *user_data.User) {
	if nil == user {
		return
	}

	stmt, err := base.MysqlDB.Prepare(sqlSaveOrUpdate)

	if nil != err {
		log.Error("%+v", err)
		return
	}

	result, err := stmt.Exec(
		user.UserName,
		user.Password,
		user.HeroAvatar,
		user.CreateTime,
		user.LastLoginTime,
		user.LastLoginTime,
	)

	if nil != err {
		log.Error("%+v", err)
		return
	}

	rowId, err := result.LastInsertId()

	if nil != err {
		log.Error("%+v", err)
	}

	user.UserId = rowId
}
