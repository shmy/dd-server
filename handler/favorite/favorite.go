package favorite

import (
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/util"
	"github.com/shmy/dd-server/model/favorite"
	"github.com/shmy/dd-server/handler/middleware/jwt"
	"github.com/globalsign/mgo/bson"
	"time"
	"errors"
	"github.com/shmy/dd-server/model/video"
	"github.com/shmy/dd-server/model/collection"
	"github.com/shmy/dd-server/service"
)

// 获取所有收藏夹
func All (c echo.Context) error {
	cc := util.ApiContext{ c }
	user := cc.Get("user")
	userClaims := user.(*jwt.ClienJwtClaims)
	r, err := favorite.M.Query(bson.M{
		"_uid": bson.ObjectIdHex(userClaims.Id),
	},
	nil,
	"-created_at",
	0,
	100) // TODO 默认倒叙一百条
	if err != nil {
		return cc.Fail(err)
	}
	service.ListCountCollection(r)
	return cc.Success(r)
}
// 创建一个收藏夹
func Create (c echo.Context) error {
	cc := util.ApiContext{ c }
	user := cc.Get("user")
	userClaims := user.(*jwt.ClienJwtClaims)
	name := cc.DefaultFormValueString("name", "", true)
	if name == "" {
		return cc.Fail(errors.New("请输入收藏夹名称"))
	}
	data := bson.M{
		"_id": bson.NewObjectId(),
		"name": name,
		"_uid": bson.ObjectIdHex(userClaims.Id),
		"created_at": time.Now(),
	}
	_, err := favorite.M.Insert(data)
	if err != nil {
		return cc.Fail(err)
	}
	return cc.Success(data)
}
// 更新一个收藏夹
func Update (c echo.Context) error {
	cc := util.ApiContext{ c }
	user := cc.Get("user")
	userClaims := user.(*jwt.ClienJwtClaims)
	name := cc.DefaultFormValueString("name", "", true)
	if name == "" {
		return cc.Fail(errors.New("请输入收藏夹名称"))
	}
	id := bson.ObjectIdHex(cc.Param("id"))
	// 判断收藏夹是不存在
	r, err := favorite.M.FindById(id, "_uid")
	if err != nil {
		return cc.Fail(err)
	}
	if r == nil {
		return cc.Fail(errors.New("收藏夹不存在"))
	}
	if r["_uid"] != bson.ObjectIdHex(userClaims.Id) {
		return cc.Fail(errors.New("访问失败，这不是你的收藏夹"))
	}
	data := bson.M{
		"name": name,
		"_uid": bson.ObjectIdHex(userClaims.Id),
	}
	r, err = favorite.M.UpdateById(id, data)
	if err != nil {
		return cc.Fail(err)
	}
	return cc.Success(r)
}
// 删除一个收藏夹
func Remove (c echo.Context) error {
	cc := util.ApiContext{ c }
	user := cc.Get("user")
	userClaims := user.(*jwt.ClienJwtClaims)
	id := bson.ObjectIdHex(cc.Param("id"))
	// 判断收藏夹是不存在
	r, err := favorite.M.FindById(id, "_uid")
	if err != nil {
		return cc.Fail(err)
	}
	if r == nil {
		return cc.Fail(errors.New("收藏夹不存在"))
	}
	if r["_uid"] != bson.ObjectIdHex(userClaims.Id) {
		return cc.Fail(errors.New("访问失败，这不是你的收藏夹"))
	}
	data := bson.M{
		"_fid": id,
	}
	// 删除相关收藏
	collection.M.RemoveAll(data)
	// 删除收藏夹
	favorite.M.RemoveById(id)
	return cc.Success(r)
}
// 添加一个资源到收藏夹
func AddToFavorite (c echo.Context) error {
	cc := util.ApiContext{ c }
	user := cc.Get("user")
	userClaims := user.(*jwt.ClienJwtClaims)
	uid := bson.ObjectIdHex(userClaims.Id)
	_vid := cc.DefaultFormValueString("vid", "", true)
	_fid := cc.DefaultFormValueString("fid", "", true)
	if _vid == "" {
		return cc.Fail(errors.New("请输入视频id"))
	}
	if !bson.IsObjectIdHex(_vid) {
		return cc.Fail(errors.New("视频id格式不正确"))
	}
	if _fid == "" {
		return cc.Fail(errors.New("请输入收藏夹id"))
	}
	if !bson.IsObjectIdHex(_fid) {
		return cc.Fail(errors.New("收藏夹id格式不正确"))
	}
	vid := bson.ObjectIdHex(_vid)
	fid := bson.ObjectIdHex(_fid)
	// 判断视频是否存在
	v, err := video.M.FindById(vid, "_id")
	if err != nil {
		return cc.Fail(err)
	}
	if v == nil {
		return cc.Fail(errors.New("视频不存在"))
	}
	// 判断收藏夹是否存在
	f, err := favorite.M.FindById(fid, "_uid")
	if err != nil {
		return cc.Fail(err)
	}
	if f == nil {
		return cc.Fail(errors.New("收藏夹不存在"))
	}
	// 判断是否是该用户的收藏夹
	if f["_uid"] != uid {
		return cc.Fail(errors.New("收藏夹不属于你"))
	}
	data := bson.M{
		//"_id": bson.NewObjectId(),
		"_vid": vid,
		"_uid": uid,
		//"created_at": time.Now(),
	}
	// 判断该人是否已经收藏过该视频了
	count, err := collection.M.Count(data)
	if err != nil {
		return cc.Fail(err)
	}
	if count != 0 {
		return cc.Fail(errors.New("已收藏过该视频了"))
	}
	data["_id"] = bson.NewObjectId()
	data["_fid"] = fid
	data["created_at"] = time.Now()
	_, err = collection.M.Insert(data)
	if err != nil {
		return cc.Fail(err)
	}
	return cc.Success(bson.M{
		"favorited": true,
		"favorited_count": service.CountVideoFavorited(vid),
	})
}
// 从收藏夹删除一个资源
func RemoveFromFavorite (c echo.Context) error {
	cc := util.ApiContext{ c }
	user := cc.Get("user")
	userClaims := user.(*jwt.ClienJwtClaims)
	_vid := cc.DefaultFormValueString("vid", "", true)
	if _vid == "" {
		return cc.Fail(errors.New("请输入视频id"))
	}
	if !bson.IsObjectIdHex(_vid) {
		return cc.Fail(errors.New("视频id格式不正确"))
	}
	vid := bson.ObjectIdHex(_vid)
 	collection.M.RemoveAll(bson.M{
		"_vid": vid,
		"_uid": bson.ObjectIdHex(userClaims.Id),
	})
	return cc.Success(bson.M{
		"favorited": false,
		"favorited_count": service.CountVideoFavorited(vid),
	})
}