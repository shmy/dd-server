package video

import (
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/util"
	"github.com/globalsign/mgo/bson"
	"errors"
	"github.com/shmy/dd-server/service"
	"github.com/shmy/dd-server/model/video"
	"math"
	"github.com/shmy/dd-server/model/classification"
	"github.com/globalsign/mgo"
)

func SearchSecret (c echo.Context) error {
	cc := &util.ApiContext{c}
	keyword := cc.DefaultQueryString("keyword", "", 1)
	query := cc.DefaultQueryString("query", "2", 1)
	sort := cc.DefaultQueryString("sort", "1", 1)
	pid := cc.DefaultQueryString("pid", "", 1)
	source := cc.DefaultQueryString("source", "", 1)
	if keyword == "" {
		return cc.Fail(errors.New("请输入搜索关键字"))
	}
	paging := util.ParsePaging(cc)
	if query == "1" {
		keyword = "^" + keyword
	}
	// 排序
	if sort == "1" {
		sort = "-generated_at"
	} else {
		sort = "+generated_at"
	}
	conditions := bson.M{
		"keyword": &bson.RegEx{keyword, "ig"},
	}
	// 有分类
	if pid != "" {
		// 判断id
		if !bson.IsObjectIdHex(pid) {
			return cc.Fail(errors.New("ID格式不正确"))
		}
		objectId := bson.ObjectIdHex(pid)
		// 查看分类是否存在
		_, err := classification.M.FindById(objectId, nil)
		if err != nil {
			return cc.Fail(err)
		}
		if err == mgo.ErrNotFound {
			return cc.Fail(errors.New("该分类不存在"))
		}
		conditions["pid"] = objectId
	} else {
		conditions["pid"] = &bson.M{"$in": service.RuleOut} // 默认排除不显示的
	}
	if source != "" { // 来源搜索
		conditions["source"] = source
	}
	// 获取总数
	total, err := video.M.Count(conditions)
	if err != nil {
		return cc.Fail(err)
	}
	v, err := video.M.Query(
		conditions,
		"name, thumbnail, latest, generated_at, _id, source",
		sort,
		paging.Offset,
		paging.Limit,
	)
	if err != nil {
		return cc.Fail(err)
	}
	//for _, el := range v { // 兼容旧版本
	//	el["quality"] = el["latest"]
	//}
	return cc.Success(&echo.Map{
		"result":    v,
		"total":     total,
		"page":      paging.Page,
		"per_page":  paging.Limit,
		"last_page": math.Ceil(float64(total) / float64(paging.Limit)),
	})
}
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
			"name": "支付宝又双叒叕发红包了，天天可领，戳我去领",
			"image":
			"https://dd.shmy.tech/static/ads/alipay/alipay_redpack.png",
			"action": bson.M {
				"type": "alipay_readpack",
				"data": "Nl7FJ976sg",
			},
		},
		{
			"name": "漫威宇宙十年电影合集",
			"image":
			"https://dd.shmy.tech/static/web_client/static/img/1.df0423e.jpg",
			"action": bson.M {
				"type": "series",
				"data": bson.M{
					"_id": "5b716e8fb8dacd1f59f942bb",
					"name": "漫威宇宙十年电影合集",
				},
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
		"_id, name, thumbnail, latest, generated_at, source",
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
				"type": "alert",
				"data": "该广告位招租，联系QQ: 2635970493",
			},
		},
		{
			"image":
			"https://dd.shmy.tech/static/ads/jd/jd.webp",
			"height": 0.24,
			"action": bson.M{
				"type": "webview",
				"data": "https://www.jd.com",
			},
		},
	}

	return cc.Success(result)
}