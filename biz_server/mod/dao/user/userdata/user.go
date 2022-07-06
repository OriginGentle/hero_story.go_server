package userdata

type User struct {
	UserId        int64  `db:"user_id"`
	UserName      string `db:"user_name"`
	Password      string `db:"password"`
	HeroAvatar    string `db:"hero_avatar"`
	CreateTime    int64  `db:"create_time"`
	LastLoginTime int64  `db:"last_login_time"`
}
