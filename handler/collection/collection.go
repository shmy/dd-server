package collection

import (
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/util"
	"github.com/globalsign/mgo/bson"
	"errors"
	"github.com/shmy/dd-server/model/collection"
	"math"
	"github.com/shmy/dd-server/handler/middleware/jwt"
	"github.com/shmy/dd-server/service"
)

func List (c echo.Context) error {
	cc := &util.ApiContext{c}
	id := cc.Param("id")
	// 判断id
	if !bson.IsObjectIdHex(id) {
		return cc.Fail(errors.New("ID格式不正确"))
	}
	objectId := bson.ObjectIdHex(id)
	conditions := bson.M{
		"_fid": objectId,
	}
	// 获取总数
	total, err := collection.M.Count(conditions)
	if err != nil {
		return cc.Fail(err)
	}

	paging := util.ParsePaging(cc) // 解析分页参数
	v, err := collection.M.Query(conditions,
		nil,
		"-created_at",
		paging.Offset,
		paging.Limit,
	)
	if err != nil {
		return cc.Fail(err)
	}
	if len(v) != 0 {
		user := cc.Get("user")
		userClaims := user.(*jwt.ClienJwtClaims)
		if v[0]["_uid"] != bson.ObjectIdHex(userClaims.Id) {
			return cc.Fail(errors.New("访问失败，这不是你的收藏夹"))
		}
	}
	service.ListCollection(v)
	return cc.Success(&echo.Map{
		"result":    v,
		"total":     total,
		"page":      paging.Page,
		"per_page":  paging.Limit,
		"last_page": math.Ceil(float64(total) / float64(paging.Limit)),
	})
}