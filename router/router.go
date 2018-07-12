package router

import (
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/handler/activity"
	"github.com/shmy/dd-server/handler/classification"
	"github.com/shmy/dd-server/handler/middleware/jwt"
	"github.com/shmy/dd-server/handler/sd"
	"github.com/shmy/dd-server/handler/user"
	"github.com/shmy/dd-server/handler/video"
	"github.com/spf13/viper"
	"github.com/shmy/dd-server/handler/app"
)

func Load(e *echo.Echo) {

	e.GET("/sw.js", func(c echo.Context) error {
		return c.File("public/web_client/build.sw.js")
	})
	// 根路径重定向
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(302, "/client")
	})
	// 客户端路由
	e.GET("/client*", func(c echo.Context) error {
		return c.File("public/web_client/index.html")
	})
	// 静态服务
	e.Static("/static", "public")

	// The health check handlers
	svcd := e.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}
	// 客户端
	apiClient := e.Group("/api/client/")
	{
		var secret = viper.GetString("jsonwebtoken.client.secret")
		// app检查更新
		apiClient.GET("check_for_update", app.Update) // ok
		// 首页推荐
		apiClient.GET("recommended", video.Recommended) // ok
		// 视频分类大全
		apiClient.GET("classification", classification.Classification) // ok
		// 视频分类列表
		apiClient.GET("classification/:id", video.List) // ok
		// 热门搜索关键字
		apiClient.GET("video/hot", video.Hot) // ok
		// 大家都在看
		apiClient.GET("video/activity", activity.List) // ok
		// 个人播放记录
		apiClient.GET("video/record", activity.Record, jwt.JWT(secret, false)) // ok
		// 视频搜索
		apiClient.GET("video/search", video.Search) // ok
		// 视频详情
		apiClient.GET("video/:id", video.Detail, jwt.JWT(secret, true)) // ok
		//apiClient.GET("video/:id", video.Detail, jwt.JWT(secret)) // ok
		//// 用户注册
		apiClient.POST("profile/sign_up", user.SignUp) // ok
		// 用户登录
		apiClient.POST("profile/sign_in", user.SignIn) // ok
		// 个人详情
		apiClient.GET("profile/detail", user.Detail, jwt.JWT(secret, false)) // ok
		// 用户登出
		apiClient.GET("profile/sign_out", user.SignOut, jwt.JWT(secret, false)) // ok

	}
	// 服务端端
	apiServer := e.Group("/api/server/")
	{
		// 用户列表
		apiServer.GET("users", user.List)
	}
}
