package server

import (
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/handler/user"
	"github.com/shmy/dd-server/handler/server/video"
)

func GetRoutes (apiServer *echo.Group) {
	// 用户分页列表
	apiServer.GET("users", user.List)
	// 视频分类列表
	apiServer.GET("classification/:id", video.List) // ok
	// 视频详情
	apiServer.GET("video/:id", video.Detail) // ok
}
