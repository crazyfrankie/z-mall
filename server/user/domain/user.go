package domain

type User struct {
	ID int64
	// 微信用户唯一标识
	OpenId   string
	UserName string
	NickName string
	Avatar   string
	Password string
	Role     string
	// 用户状态,为软删除做准备
	Status uint8
	Ctime  int64
	Utime  int64
}

type WeChatInfo struct {
	OpenID  string
	UnionID string
}
