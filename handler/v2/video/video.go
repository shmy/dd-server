package video

import (
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/util"
	"github.com/globalsign/mgo/bson"
	"errors"
	"github.com/shmy/dd-server/service"
	"github.com/shmy/dd-server/model/video"
)

func Index (c echo.Context) error {
	cc := util.ApiContext{ c }
	_id := cc.DefaultQueryString("id", "", true)
	if !bson.IsObjectIdHex(_id) {
		return cc.Fail(errors.New("ID格式不正确"))
	}
	id := bson.ObjectIdHex(_id)
	var result = make(bson.M)
	result["banner"] = []bson.M{
		{
			"name": "支付宝又发红包了，天天可领，点击立即去领",
			"image":
			"http://n.sinaimg.cn/finance/transform/20170330/z8Fu-fycwyns3693714.jpg",
			"action": bson.M {
				"type": "alipay_readpack",
				"data": "Nl7FJ976sg",
			},
		},
	}
	// 获取最新资源
	ids, _ := service.FindClassifyIds(id)
	var conditions = bson.M{"pid": ids[0]}
	if len(ids) > 1 {
		conditions["pid"] = &bson.M{"$in": ids}
	}
	var err error
	// 获取最新资源结果集
	result["latests"], err = video.M.Query(conditions,
		"_id, name, thumbnail, latest, generated_at",
		"-generated_at",
		0,
		8,
	)
	if err != nil {
		return cc.Fail(err)
	}
	// 获取最热资源结果集
	result["hots"], err = video.M.Query(conditions,
		"_id, name, thumbnail, latest, generated_at",
		"-number",
		0,
		8,
	)
	if err != nil {
		return cc.Fail(err)
	}
	// 添加两条广告
	// TODO 自动读取广告
	result["ads"] = []bson.M{
		{
			"image":
			"https://img.zcool.cn/community/0145735928d586a801216a3e141620.png@1280w_1l_2o_100sh.webp",
			"height": 0.4,
			"action": bson.M{
				"type": "webview",
				"data": "https://www.jd.com",
			},
		},
		{
			"image":
			"https://img.zcool.cn/community/0145735928d586a801216a3e141620.png@1280w_1l_2o_100sh.webp",
			"height": 0.4,
			"action": bson.M{
				"type": "alert",
				"data": "该广告位招租，联系QQ: 2635970493",
			},
		},
	}

	return cc.Success(result)
}