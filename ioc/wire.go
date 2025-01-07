//go:build wireinject

package ioc

import (
	"github.com/google/wire"
	
	"zmall/server/user"
)

var BaseSet = wire.NewSet(InitDB)

func NewApp() *App {
	wire.Build(
		BaseSet,

		user.NewModule,

		InitWebServer,

		wire.FieldsOf(new(*user.Module), "Hdl"),
		wire.Struct(new(App), "*"),
	)

	return new(App)
}
