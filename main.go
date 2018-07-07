package main

import (
	"errors"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/lexkong/log"
	"github.com/shmy/dd-server/router"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func main() {
	runmode := viper.GetString("runmode")
	e := echo.New()
	e.HideBanner = runmode == "release"
	e.Debug = runmode == "debug"
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		var message interface{} = "Internal Server Error"
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = he.Message
		}
		if code == http.StatusInternalServerError {
			log.Error("StatusInternalServerError ", err)
		}
		c.JSON(code, map[string]interface{}{
			"success": false,
			"payload": nil,
			"message": message,
		})
	}
	if runmode == "release" {
		e.Use(middleware.Recover())
	}
	e.Use(middleware.CORS())

	router.Load(e)

	// 启动时自检
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Info("The router has been deployed successfully.")
	}()
	host := viper.GetString("server.host")
	port := viper.GetString("server.port")

	log.Infof("Start to listening the incoming requests on http address: %s", host+":"+port)
	s := &http.Server{
		Addr:         host + ":" + port,
		ReadTimeout:  20 * time.Minute,
		WriteTimeout: 20 * time.Minute,
	}
	e.Logger.Fatal(e.StartServer(s))

}

// 自己Ping自己
func pingServer() error {
	host := viper.GetString("server.host")
	port := viper.GetString("server.port")
	url := "http://" + host + ":" + port + "/sd/health"
	for i := 0; i < 20; i++ {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			return nil
		}
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	//noinspection GoErrorStringFormat
	return errors.New("Cannot connect to the router.")
}
