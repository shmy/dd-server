package user

import (
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/util"
	"github.com/shmy/dd-server/model/user"
	"math"
)

// 用户列表
func List (c echo.Context) error {
	cc := util.ApiContext{ c }
	// 解析分页
	paging := util.ParsePaging(&cc)
	// 获取总数
	total, err := user.M.Count(nil)
	if err != nil {
		return cc.Fail(err)
	}
	// 获取结果集
	v, err := user.M.Query(
		nil,
		"username, created_at, avatar, email, integral, level, nickname",
		"-created_at",
		paging.Offset,
		paging.Limit,
	)
	if err != nil {
		return cc.Fail(err)
	}

	return cc.Success(&echo.Map{
		"result":    v,
		"total":     total,
		"page":      paging.Page,
		"per_page":  paging.Limit,
		"last_page": math.Ceil(float64(total) / float64(paging.Limit)),
	})
}
