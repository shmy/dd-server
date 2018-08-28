package svcd

import (
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/handler/sd"
)

func GetRoutes (apiSvcd *echo.Group) {
	apiSvcd.GET("/health", sd.HealthCheck)
	apiSvcd.GET("/disk", sd.DiskCheck)
	apiSvcd.GET("/cpu", sd.CPUCheck)
	apiSvcd.GET("/ram", sd.RAMCheck)
}
