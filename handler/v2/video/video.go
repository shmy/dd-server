package video

import (
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/util"
	"github.com/globalsign/mgo/bson"
	"errors"
	"github.com/shmy/dd-server/service"
	"github.com/shmy/dd-server/model/video"
	"math"
	"time"
)

func SearchSecret (c echo.Context) error {
	cc := &util.ApiContext{c}
	keyword := cc.DefaultQueryString("keyword", "", 1)
	if keyword == "" {
		return cc.Fail(errors.New("请输入搜索关键字"))
	}
	paging := util.ParsePaging(cc)
	qs := util.ParseQueryString(cc)
	if qs["query"] == "1" {
		keyword = "^" + keyword
	}

	conditions := bson.M{
		"keyword": &bson.RegEx{keyword, "ig"},
	}
	// 有分类
	if qs["pid"] != nil {
		conditions["pid"] = qs["pid"]
	} else {
		conditions["pid"] = &bson.M{"$in": service.RuleOut} // 默认排除不显示的
	}
	if qs["source"] != nil { // 来源搜索
		conditions["source"] = qs["source"]
	}
	if qs["released_at"] != nil { // 年代搜索
		conditions["released_at"] = qs["released_at"]
	}
	if qs["region"] != nil { // 区域搜索
		conditions["region"] = qs["region"]
	}
	// 获取总数
	total, err := video.M.Count(conditions)
	if err != nil {
		return cc.Fail(err)
	}
	v, err := video.M.Query(
		conditions,
		"name, thumbnail, latest, generated_at, _id, source",
		qs["sort"],
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
			"name": "🙏🙏救救这个女孩🙏🙏",
			"image":
			"http://cf.alioss.shuidichou.com/img/ck/20181012/d0e0972a-db85-44e1-b24c-29b091018ea8!cf_mtr_200_nw",
			"action": bson.M {
				"type": "browser",
				"data": "https://www.shuidichou.com/cf/contribute/caff17ed-905e-460b-a65a-8f0e943d47ae?channel=wx_charity_hy",
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
		12,
	)
	if err != nil {
		return cc.Fail(err)
	}
	// 获取最热资源结果集
	//old, _ := time.ParseDuration("-15d")
	conditions["generated_at"] = bson.M{
		// 十五天热门
		"$gte": time.Now().Add(-time.Hour * 24 * 15),
	}
	result["hots"], err = video.M.Query(conditions,
		"_id, name, thumbnail, latest, generated_at",
		"-number",
		0,
		12,
	)
	if err != nil {
		return cc.Fail(err)
	}
	// 添加两条广告
	// TODO 自动读取广告
	result["ads"] = []bson.M{
		{
			"image":
			"http://cf.alioss.shuidichou.com/img/ck/20181012/d0e0972a-db85-44e1-b24c-29b091018ea8!cf_mtr_200_nw",
			"height": 0.4,
			"action": bson.M{
				"type": "browser",
				"data": "https://www.shuidichou.com/cf/contribute/caff17ed-905e-460b-a65a-8f0e943d47ae?channel=wx_charity_hy",
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