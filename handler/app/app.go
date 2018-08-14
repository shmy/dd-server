package app

import (
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/util"
	"os"
	"io/ioutil"
	"encoding/json"
	"net/http"
)
type Info struct {
	Version string `json:"version"`
	ApkUrl string `json:"url"`
	Website string `json:"website"`
	Date string `json:"date"`
	Content []string `json:"content"`
}
func Update (c echo.Context) error {
	cc := util.ApiContext{ c }
	f, err := os.Open("./public/app.json")
	if err != nil {
		return cc.Fail(err)
	}
	jsonByte, err := ioutil.ReadAll(f)
	if err != nil {
		return cc.Fail(err)
	}
	var info Info
	err = json.Unmarshal(jsonByte, &info)
	if err != nil {
		return cc.Fail(err)
	}
	return cc.Success(info)
}

func Download (c echo.Context) error {
	//return c.Redirect(302, "https://www.pgyer.com/F3AF")
	f, err := os.Open("./public/app.json")
	if err != nil {
		return c.String(http.StatusInternalServerError, "StatusInternalServerError")
	}
	jsonByte, err := ioutil.ReadAll(f)
	if err != nil {
		return c.String(http.StatusInternalServerError, "StatusInternalServerError")
	}
	var info Info
	err = json.Unmarshal(jsonByte, &info)
	if err != nil {
		return c.String(http.StatusInternalServerError, "StatusInternalServerError")
	}
	return c.Render(http.StatusOK, "download.html", info)
}
