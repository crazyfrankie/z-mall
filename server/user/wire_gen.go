// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package user

import (
	"gorm.io/gorm"
	"zmall/server/user/repository"
	"zmall/server/user/repository/dao"
	"zmall/server/user/service"
	"zmall/server/user/web"
)

// Injectors from wire.go:

func NewModule(db *gorm.DB) *Module {
	userDao := dao.NewUserDao(db)
	userRepo := repository.NewUserRepo(userDao)
	userService := service.NewUserService(userRepo)
	userHandler := web.NewUserHandler(userService)
	module := &Module{
		Hdl: userHandler,
	}
	return module
}
