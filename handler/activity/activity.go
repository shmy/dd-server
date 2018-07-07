package activity

import (
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/model/activity"
	"github.com/shmy/dd-server/service"
	"github.com/shmy/dd-server/util"
	"math"
)

func List(c echo.Context) error {
	cc := util.ApiContext{c}
	paging := util.ParsePaging(&cc)
	total, err := activity.M.Count(nil)
	if err != nil {
		return cc.Fail(err)
	}
	// 获取分页数据
	ret, err := activity.M.Query(nil,
		"_id, uid, vid, updated_at",
		"-updated_at,-created_at",
		paging.Offset,
		paging.Limit,
	)
	if err != nil {
		return cc.Fail(err)
	}
	// 设置关联到用户 和 视频
	service.ListActivity(ret)
	return cc.Success(&echo.Map{
		"result":    ret,
		"total":     total,
		"page":      paging.Page,
		"per_page":  paging.Limit,
		"last_page": math.Ceil(float64(total) / float64(paging.Limit)),
	})
}
