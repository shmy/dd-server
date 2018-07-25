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
)

// 获取所有收藏夹
func All (c echo.Context) error {
	cc := util.ApiContext{ c }
	user := cc.Get("user")
	userClaims := user.(*jwt.ClienJwtClaims)
	r, err := favorite.M.Find(bson.M{
		"_uid": bson.ObjectIdHex(userClaims.Id),
	}, nil)
	if err != nil {
		return cc.Fail(err)
	}
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

// 添加一个资源到收藏夹

func AddToFavorite (c echo.Context) error {
	cc := util.ApiContext{ c }
	user := cc.Get("user")
	userClaims := user.(*jwt.ClienJwtClaims)
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
	_uid := userClaims.Id
	uid := bson.ObjectIdHex(_uid)
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
	return cc.Success(data)
}