package video

import (
	"errors"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/model/classification"
	"github.com/shmy/dd-server/model/hot"
	"github.com/shmy/dd-server/model/video"
	"github.com/shmy/dd-server/service"
	"github.com/shmy/dd-server/util"
	"math"
	"time"
)

// 视频列表
func List(c echo.Context) error {
	cc := &util.ApiContext{c}
	id := cc.Param("id")
	qs := util.ParseQueryString(cc)
	// 判断id
	if !bson.IsObjectIdHex(id) {
		return cc.Fail(errors.New("ID格式不正确"))
	}
	objectId := bson.ObjectIdHex(id)
	// 查看分类是否存在
	classify, err := classification.M.FindById(objectId, nil)
	if err != nil {
		return cc.Fail(err)
	}
	if err == mgo.ErrNotFound {
		return cc.Fail(errors.New("该分类不存在"))
	}
	ids, err := service.FindClassifyIds(classify["_id"])
	if err != nil {
		return cc.Fail(err)
	}
	var conditions = bson.M{
		"pid": ids[0],
	}
	if len(ids) > 1 {
		conditions["pid"] = &bson.M{"$in": ids}
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
	// 关键字搜索
	keyword := cc.DefaultQueryString("keyword", "", 1)
	if keyword != "" {
		if qs["query"] == "1" {
			keyword = "^" + keyword
		}
		conditions["keyword"] = &bson.RegEx{keyword, "ig"}
	}
	// 获取总数
	total, err := video.M.Count(conditions)
	if err != nil {
		return cc.Fail(err)
	}
	paging := util.ParsePaging(cc) // 解析分页参数
	v, err := video.M.Query(conditions,
		"name, thumbnail, latest, _id, generated_at, source, language, number, region, released_at",
		qs["sort"],
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
		"title":     classify["name"],
		"result":    v,
		"total":     total,
		"page":      paging.Page,
		"per_page":  paging.Limit,
		"last_page": math.Ceil(float64(total) / float64(paging.Limit)),
	})
}

// 热门搜索
func Hot(c echo.Context) error {
	cc := &util.ApiContext{c}
	now := time.Now()
	dur, _ := time.ParseDuration("-360h") // 查询15天热门
	old := now.Add(dur)
	u, err := hot.M.Query(
		bson.M{
			"updated_at": bson.M{
				"$gte": old,
				"$lte": now,
			},
		},
		nil,
		"-index",
		0,
		12,
	)
	if err != nil {
		return cc.Fail(err)
	}
	service.GetHotsThumbnail(u)
	return cc.Success(u)
}

// 搜索
func Search(c echo.Context) error {
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
		// 查看分类是否存在
		classify, err := classification.M.FindById(qs["pid"], nil)
		if err != nil {
			return cc.Fail(err)
		}
		if err == mgo.ErrNotFound {
			return cc.Fail(errors.New("该分类不存在"))
		}
		ids, err := service.FindClassifyIds(classify["_id"])
		if err != nil {
			return cc.Fail(err)
		}
		conditions["pid"] = ids[0]
		if len(ids) > 1 {
			conditions["pid"] = &bson.M{"$in": ids}
		}
	} else {
		conditions["pid"] = &bson.M{"$nin": service.RuleOut} // 默认排除不显示的
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

// 视频详情
func Detail(c echo.Context) error {
	cc := &util.ApiContext{c}
	id := cc.Param("id")
	// 判断id
	if !bson.IsObjectIdHex(id) {
		return cc.Fail(errors.New("id格式不正确"))
	}
	objectId := bson.ObjectIdHex(id)
	ret, err := video.M.FindById(objectId, nil)
	if err != nil {
		return cc.Fail(err)
	}
	if ret == nil {
		return cc.Fail(errors.New("视频不存在"))
	}
	// 增加浏览次数
	if ret["number"] == nil {
		ret["number"] = 0
	}
	number := ret["number"].(int) + 1
	ret["number"] = number
	video.M.UpdateById(ret["_id"], bson.M{
		"number": ret["number"],
	})

	// 给出收藏次数
	ret["favorited_count"] = service.CountVideoFavorited(objectId)
	// 获取分类
	ret["classify"], err = classification.M.FindById(ret["pid"], "name")
	if err != nil {
		return cc.Fail(err)
	}

	return cc.Success(ret)
}
