package app

import (
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/util"
	"io/ioutil"
	"os"
	"encoding/json"
)

func Update (c echo.Context) error {
	cc := util.ApiContext{ c }
	//cc.Set("Content-Type", "application/json; charset=UTF-8")
	f, err := os.Open("./public/app.json")
	if err != nil {
		return cc.Fail(err)
	}
	jsonByte, err := ioutil.ReadAll(f)
	if err != nil {
		return cc.Fail(err)
	}
	var j interface{}
	err = json.Unmarshal(jsonByte, &j)
	if err != nil {
		return cc.Fail(err)
	}
	return cc.Success(j)
}
