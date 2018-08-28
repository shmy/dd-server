package router

import (
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/handler/app"
	"github.com/shmy/dd-server/router/server"
	"github.com/shmy/dd-server/router/client"
	"github.com/shmy/dd-server/router/svcd"
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
	e.GET("/download", app.Download)
	// The health check handlers
	apiSvcd := e.Group("/sd")
	svcd.GetRoutes(apiSvcd)
	// 客户端
	apiClient := e.Group("/api/client/")
	client.GetRoutes(apiClient)
	//{
	//	var secret = viper.GetString("jsonwebtoken.client.secret")
	//	// app检查更新
	//	apiClient.GET("check_for_update", app.Update) // ok
	//	// 首页推荐
	//	apiClient.GET("recommended", video.Recommended) // ok
	//	// 视频分类大全
	//	apiClient.GET("classification", classification.Classification) // ok
	//	// 视频分类列表
	//	apiClient.GET("classification/:id", video.List) // ok
	//	// 热门搜索关键字
	//	apiClient.GET("video/hot", video.Hot) // ok
	//	// 大家都在看
	//	apiClient.GET("video/activity", activity.List) // ok
	//	// 个人播放记录
	//	apiClient.GET("video/record", activity.Record, jwt.JWT(secret, false)) // ok
	//	// 视频搜索
	//	apiClient.GET("video/search", video.Search) // ok
	//	// 视频详情
	//	apiClient.GET("video/:id", video.Detail, jwt.JWT(secret, true)) // ok
	//	//apiClient.GET("video/:id", video.Detail, jwt.JWT(secret)) // ok
	//	//// 用户注册
	//	apiClient.POST("profile/sign_up", user.SignUp) // ok
	//	// 用户登录
	//	apiClient.POST("profile/sign_in", user.SignIn) // ok
	//	// 个人详情
	//	apiClient.GET("profile/detail", user.Detail, jwt.JWT(secret, false)) // ok
	//	// 用户登出
	//	apiClient.GET("profile/sign_out", user.SignOut, jwt.JWT(secret, false)) // ok
	//
	//	// 获取所有收藏夹
	//	apiClient.GET("favorite", favorite.All, jwt.JWT(secret, false)) // ok
	//	// 更新一个收藏夹
	//	apiClient.PUT("favorite/:id", favorite.Update, jwt.JWT(secret, false)) // ok
	//	// 向收藏夹添加一个视频
	//	apiClient.POST("favorite/add_video", favorite.AddToFavorite, jwt.JWT(secret, false)) // ok
	//	// 移除一个收藏的视频
	//	apiClient.POST("favorite/remove_video", favorite.RemoveFromFavorite, jwt.JWT(secret, false)) // ok
	//	// 新建一个收藏夹
	//	apiClient.POST("favorite", favorite.Create, jwt.JWT(secret, false)) // ok
	//	// 移除一个收藏夹
	//	apiClient.DELETE("favorite/:id", favorite.Remove, jwt.JWT(secret, false)) // ok
	//	// 根据收藏夹id获取分页列表
	//	apiClient.GET("collection/:id", collection.List, jwt.JWT(secret, false)) // ok
	//
	//	// 根据id获取播单详情
	//	apiClient.GET("series/:id", series.SeriesDetail) // ok
	//
	//	// v2版本首页数据
	//	apiClient.GET("v2/video/index", v2.Index) // ok
	//	// v2版本秘密花园搜索
	//	apiClient.GET("v2/video/search_secret", v2.SearchSecret, jwt.JWT(secret, false)) // ok
	//
	//
	//	//// 测试获取vip分类
	//	//apiClient.GET("vip/list", vip.GetList) // ok
	//	//apiClient.GET("vip/classify", vip.GetClassifyList) // ok
	//	//// 测试获取vip视频分集地址
	//	//apiClient.POST("vip/detail", vip.GetPlayUrls) // ok
	//
	//}
	// 服务端端
	apiServer := e.Group("/api/server/")
	server.GetRoutes(apiServer)
}
