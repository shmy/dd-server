package video

import (
	"errors"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
	"github.com/lexkong/log"
	"github.com/shmy/dd-server/model/classification"
	"github.com/shmy/dd-server/model/hot"
	"github.com/shmy/dd-server/model/video"
	"github.com/shmy/dd-server/service"
	"github.com/shmy/dd-server/util"
	"math"
	"time"
	"github.com/shmy/dd-server/handler/middleware/jwt"
)

// 推荐列表
func Recommended(c echo.Context) error {
	cc := &util.ApiContext{c}
	r, err := classification.M.Find(bson.M{"pid": nil}, nil)
	if err != nil {
		return cc.Fail(err)
	}
	var latest = []*echo.Map{}
	for _, val := range r {
		ids, _ := service.FindClassifyIds(val["_id"])
		var conditions = bson.M{"pid": ids[0]}
		if len(ids) > 1 {
			conditions["pid"] = &bson.M{"$in": ids}
		}
		// 获取结果集
		v, err := video.M.Query(conditions,
			"_id, name, thumbnail, latest, generated_at",
			"-generated_at",
			0,
			16,
		)
		if err != nil {
			return cc.Fail(err)
		}
		latest = append(latest, &echo.Map{
			"name":     "最近更新的" + val["name"].(string),
			"_id":      val["_id"],
			"children": &v,
		})
	}
	return cc.Success(&echo.Map{
		"latest": latest,
	})
}

// 视频列表
func List(c echo.Context) error {
	cc := &util.ApiContext{c}
	id := cc.Param("id")
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
	// 获取总数
	total, err := video.M.Count(conditions)
	if err != nil {
		return cc.Fail(err)
	}
	paging := util.ParsePaging(cc) // 解析分页参数
	v, err := video.M.Query(conditions,
		"name, thumbnail, latest, _id, generated_at, source",
		"-generated_at",
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
	//query := cc.DefaultQueryString("query", "2", 1)
	//sort := cc.DefaultQueryString("sort", "1", 1)
	//pid := cc.DefaultQueryString("pid", "", 1)
	//source := cc.DefaultQueryString("source", "", 1)
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

// 视频详情
func Detail(c echo.Context) error {
	cc := &util.ApiContext{c}
	id := cc.Param("id")
	from := cc.DefaultQueryString("from", "", 1)
	//if from == "" {
	//	from = "normal"
	//}
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
	// 更新或者添加记录
	//user := cc.Get("user")
	//if user != nil {
	//	userClaims := user.(*jwt2.ClienJwtClaims)
	//	uid := bson.ObjectIdHex(userClaims.Id)
	//	service.AddToActivity(ret, uid)
	//}
	// 查询该人是否收藏
	user := cc.Get("user")
	if user != nil {
		userClaims := user.(*jwt.ClienJwtClaims)
		uid := bson.ObjectIdHex(userClaims.Id)
		ret["favorited"] = service.CheckIsFavorited(uid, ret["_id"])
	} else {
		// 没收藏
		ret["favorited"] = false
	}
	// 给出收藏次数
	ret["favorited_count"] = service.CountVideoFavorited(objectId)
	// 设置热搜
	if from == "search" && ret != nil {
		red, _ := hot.M.FindOne(bson.M{
			"vid": ret["_id"],
		}, nil)
		if red == nil {
			u := bson.M{
				"_id":        bson.NewObjectId(),
				"name":       ret["name"],
				"index":      1,
				"vid":        ret["_id"],
				"created_at": time.Now(),
				"updated_at": time.Now(),
			}
			_, err := hot.M.Insert(u)
			if err != nil {
				log.Warn("ADD HOT:" + err.Error())
			}
		} else {
			var index int
			if red["index"] == nil {
				index = 0
			}
			index = red["index"].(int)
			_, err := hot.M.UpdateById(red["_id"], bson.M{
				"index": index + 1,
				"updated_at": time.Now(),
			})
			if err != nil {
				log.Warn("UPDATE HOT:" + err.Error())
			}
		}
	}
	// 获取分类
	ret["classify"], err = classification.M.FindById(ret["pid"], "name")
	if err != nil {
		return cc.Fail(err)
	}
	// TODO 自动读取广告
	ret["ads"] = []bson.M {
		{
			"image":
			"https://dd.shmy.tech/static/ads/eye/eye.webp",
			"height": 0.666,
			"action": bson.M{
				"type": "webview",
				"data": "https://dd.shmy.tech/static/ads/eye",
			},
		},
	}
	return cc.Success(ret)
}
