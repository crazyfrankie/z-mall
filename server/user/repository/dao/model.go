package dao

import "database/sql"

type User struct {
	ID       int64           `gorm:"primaryKey;autoIncrement;comment:用户ID"`
	OpenId   string          `gorm:"type:varchar(255);uniqueIndex;not null;comment:微信用户唯一标识"`
	NickName string          `gorm:"type:varchar(50);comment:用户昵称"`
	Avatar   string          `gorm:"type:varchar(255);comment:头像URL"`
	UserName *sql.NullString `gorm:"type:varchar(50);uniqueIndex;not null;comment:用户名"`
	Password *sql.NullString `gorm:"type:varchar(255);not null;comment:用户密码"`
	Role     string          `gorm:"type:varchar(20);default:'user';not null;comment:用户角色"`
	Status   uint8           `gorm:"type:tinyint;default:1;not null;comment:用户状态"`
	Ctime    int64           `gorm:"autoCreateTime:milli;comment:创建时间"`
	Utime    int64           `gorm:"autoUpdateTime:milli;comment:更新时间"`
}
