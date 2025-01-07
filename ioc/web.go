package ioc

import (
	"github.com/gin-gonic/gin"

	"zmall/server/user"
)

func InitWebServer(user *user.Handler) *gin.Engine {
	server := gin.Default()

	user.RegisterRoute(server)

	return server
}
