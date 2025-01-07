package user

import "zmall/server/user/web"

type Handler = web.UserHandler

type Module struct {
	Hdl *Handler
}
