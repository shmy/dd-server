package series

import (
	"github.com/shmy/dd-server/service"
	"github.com/globalsign/mgo/bson"
	"errors"
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/util"
	"github.com/shmy/dd-server/model/series"
)

func SeriesDetail(c echo.Context) error  {
	cc := util.ApiContext{ c }
	_id := cc.Param("id")
	if !bson.IsObjectIdHex(_id) {
		return cc.Fail(errors.New("ID格式不正确"))
	}
	id := bson.ObjectIdHex(_id)
	ret, err := _series.M.FindById(id, nil)
	if err != nil {
		return cc.Fail(err)
	}
	if ret == nil {
		return cc.Fail(errors.New("播单不存在"))
	}
	service.GetSeriesDetail(ret["series"].([]interface{}))
	return cc.Success(ret)

}