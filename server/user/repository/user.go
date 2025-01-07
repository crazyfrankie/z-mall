package repository

import (
	"context"
	"database/sql"
	"zmall/server/user/domain"
	"zmall/server/user/repository/dao"
)

type UserRepo struct {
	dao *dao.UserDao
}

func NewUserRepo(dao *dao.UserDao) *UserRepo {
	return &UserRepo{dao: dao}
}

func (repo *UserRepo) CreateUser(ctx context.Context, user domain.User) error {
	return repo.dao.Create(ctx, repo.domainToDao(user))
}

func (repo *UserRepo) FindByWechat(ctx context.Context, openId string) (domain.User, error) {
	user, err := repo.dao.FindUser(ctx, openId)
	if err != nil {
		return domain.User{}, err
	}

	return repo.daoToDomain(user), nil
}

func (repo *UserRepo) domainToDao(user domain.User) dao.User {
	return dao.User{
		ID:       user.ID,
		OpenId:   user.OpenId,
		NickName: user.NickName,
		Avatar:   user.Avatar,
		UserName: &sql.NullString{
			String: user.UserName,
			Valid:  user.UserName != "",
		},
		Password: &sql.NullString{
			String: user.Password,
			Valid:  user.Password != "",
		},
		Status: user.Status,
		Role:   user.Role,
	}
}

func (repo *UserRepo) daoToDomain(user dao.User) domain.User {
	return domain.User{
		ID:       user.ID,
		OpenId:   user.OpenId,
		NickName: user.NickName,
		Avatar:   user.Avatar,
		UserName: user.UserName.String,
		Password: user.Password.String,
		Status:   user.Status,
		Role:     user.Role,
		Ctime:    user.Ctime,
		Utime:    user.Utime,
	}
}
