package router

import (
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/handler/app"
	"github.com/shmy/dd-server/router/server"
	"github.com/shmy/dd-server/router/client"
	"github.com/shmy/dd-server/router/svcd"
)

func Load(e *echo.Echo) {

	//e.GET("/sw.js", func(c echo.Context) error {
	//	return c.File("public/web_client/build.sw.js")
	//})
	// 根路径重定向 2018-10-20 重定项到新域名
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(301, "https://v.shmy.tech")
	})
	// 客户端路由 2018-10-20 重定项到新域名
	e.GET("/client*", func(c echo.Context) error {
		return c.Redirect(301, "https://v.shmy.tech")
	})
	// 管理端路由
	e.GET("/admin*", func(c echo.Context) error {
		return c.File("public/web_admin/index.html")
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

	// 服务端端
	apiServer := e.Group("/api/server/")
	server.GetRoutes(apiServer)
}
