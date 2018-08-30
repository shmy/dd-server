package server

import (
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/handler/user"
	"github.com/shmy/dd-server/handler/server/video"
	"github.com/shmy/dd-server/handler/server/admin"
	"github.com/shmy/dd-server/handler/middleware/jwt"
	"github.com/spf13/viper"
)

func GetRoutes (apiServer *echo.Group) {
	var secret = viper.GetString("jsonwebtoken.admin.secret")
	// 用户分页列表
	apiServer.GET("users", user.List, jwt.JWT(secret, false ,1))
	// 视频分类列表
	apiServer.GET("classification/:id", video.List, jwt.JWT(secret, false ,1)) // ok
	// 视频详情
	apiServer.GET("video/:id", video.Detail, jwt.JWT(secret, false ,1)) // ok
	// 修改视频
	apiServer.PUT("video/:id", video.Update, jwt.JWT(secret, false ,1)) // ok
	// 登录
	apiServer.POST("admin/sign_in", admin.SignIn) // ok
}
