package video

import (
	"errors"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/model/classification"
	"github.com/shmy/dd-server/model/video"
	"github.com/shmy/dd-server/service"
	"github.com/shmy/dd-server/util"
	"math"
	"time"
	"fmt"
)
//type Video struct {
//	Name  string `json:"name" form:"name" query:"name"`
//	Email string `json:"email" form:"email" query:"email"`
//
//}


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
	//number := ret["number"].(int) + 1
	//ret["number"] = number
	//video.M.UpdateById(ret["_id"], bson.M{
	//	"number": ret["number"],
	//})

	// 给出收藏次数
	//ret["favorited_count"] = service.CountVideoFavorited(objectId)
	// 获取分类
	ret["classify"], err = classification.M.FindById(ret["pid"], "name")
	if err != nil {
		return cc.Fail(err)
	}

	return cc.Success(ret)
}

// 修改视频
func Update(c echo.Context) error {
	cc := &util.ApiContext{c}
	data := cc.GetJSONBody()
	if data == nil {
		return cc.Fail(errors.New("请提交要修改的数据"))
	}
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
	data["pid"] = bson.ObjectIdHex(data["pid"].(string))
	data["running_time"] = int(data["running_time"].(float64))
	data["number"] = int(data["number"].(float64))
	data["generated_at"], err = time.Parse("2006-01-02T15:04:05Z", data["generated_at"].(string))
	if err != nil {
		fmt.Println(err)
		delete(data, "generated_at")
	}
	ret, err = video.M.UpdateById(objectId, data)
	if err != nil {
		return cc.Fail(err)
	}
	return cc.Success(ret)
}
