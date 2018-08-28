package server

import (
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/handler/user"
)

func GetRoutes (apiServer *echo.Group) {
	apiServer.GET("users", user.List)
}
