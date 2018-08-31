package ad

import (
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/util"
	"github.com/shmy/dd-server/model/ad"
	"math"
	"errors"
	"time"
	"github.com/globalsign/mgo/bson"
)

// 广告分页
func List (c echo.Context) error {
	cc := util.ApiContext{ c }
	// 解析分页
	paging := util.ParsePaging(&cc)
	// 获取总数
	total, err := ad.M.Count(nil)
	if err != nil {
		return cc.Fail(err)
	}
	// 获取结果集
	v, err := ad.M.Query(
		nil,
		"name, type, data, created_at",
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

// 新增广告
func Create (c echo.Context) error {
	cc := util.ApiContext{ c }
	name := cc.GetJSONString("name")
	_type := cc.GetJSONString("type")
	data := cc.GetJSONString("data")
	if name == nil || *name == "" {
		return cc.Fail(errors.New("广告名称不能为空"))
	}
	if _type == nil || *_type == "" {
		return cc.Fail(errors.New("广告类型不能为空"))
	}
	if data == nil || *data == "" {
		return cc.Fail(errors.New("广告参数不能为空"))
	}
	ret, err := ad.M.Insert(map[string]interface{} {
		"name": name,
		"type": _type,
		"data": data,
		"created_at": time.Now(),
	})
	if err != nil {
		return cc.Fail(err)
	}
	return cc.Success(ret)
}

func Delete (c echo.Context) error {
	cc := util.ApiContext{ c }
	_id := cc.Param("id")
	if !bson.IsObjectIdHex(_id) {
		return cc.Fail(errors.New("ID格式不正确"))
	}
	id := bson.ObjectIdHex(_id)
	if !ad.M.RemoveById(id) {
		return cc.Fail(errors.New("删除失败"))
	}
	return cc.Success(nil)
}