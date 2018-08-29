package main

import (
	"errors"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/lexkong/log"
	"github.com/shmy/dd-server/router"
	"html/template"
	"github.com/spf13/viper"
	"net/http"
	"time"
	"github.com/spf13/pflag"
	"encoding/json"
	"fmt"
	"os"
	version2 "github.com/shmy/dd-server/pkg/version"
	"github.com/shmy/dd-server/util"
)

var (
	version = pflag.BoolP("version", "v", false, "show version info.")
)
func main() {
	// 命令获取版本号信息
	pflag.Parse()
	if *version {
		v := version2.Get()
		marshalled, err := json.MarshalIndent(&v, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(marshalled))
		return
	}
	runmode := viper.GetString("runmode")
	e := echo.New()
	e.HideBanner = runmode == "release"
	e.Debug = runmode == "debug"
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		logback := false
		var message interface{} = "Internal Server Error"
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = he.Message
		}
		if code == http.StatusInternalServerError {
			log.Error("StatusInternalServerError ", err)
		}
		// 不要出错 转换为200
		if code == http.StatusUnauthorized {
			code = http.StatusOK
			logback = true
		}
		c.JSON(code, map[string]interface{}{
			"success": false,
			"logback": logback,
			"payload": nil,
			"message": message,
		})
	}
	if runmode == "release" {
		e.Use(middleware.Recover())
	}
	// 模板引擎
	t := &util.Template{
		Templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Renderer = t
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
