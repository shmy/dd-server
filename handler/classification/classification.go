package classification

import (
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/model/classification"
	"github.com/shmy/dd-server/util"
)

// 获取所有分类
func Classification(c echo.Context) error {
	cc := &util.ApiContext{c}
	r, err := classification.M.FindAll("_id, name, pid")
	if err != nil {
		return cc.Fail(err)
	}
	return cc.Success(&r)
}
