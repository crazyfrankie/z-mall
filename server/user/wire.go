//go:build wireinject

package user

import (
	"github.com/google/wire"
	"gorm.io/gorm"
	"zmall/server/user/repository"
	"zmall/server/user/repository/dao"
	"zmall/server/user/service"
	"zmall/server/user/web"
)

func NewModule(db *gorm.DB) *Module {
	wire.Build(
		dao.NewUserDao,
		repository.NewUserRepo,
		service.NewUserService,
		web.NewUserHandler,

		wire.Struct(new(Module), "*"),
	)
	return new(Module)
}
