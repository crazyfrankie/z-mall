package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (dao *UserDao) Create(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.Ctime = now
	user.Utime = now

	return dao.db.WithContext(ctx).Model(&User{}).Create(&user).Error
}

func (dao *UserDao) FindUser(ctx context.Context, openId string) (User, error) {
	var user User

	err := dao.db.WithContext(ctx).Model(&User{}).Where("open_id = ?", openId).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, ErrUserNotFound
		}

		return User{}, err
	}

	return user, nil
}
