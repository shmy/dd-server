package user

import (
	"errors"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/handler/middleware/jwt"
	"github.com/shmy/dd-server/model/user"
	"github.com/shmy/dd-server/util"
	"net/http"
	"regexp"
	"time"
	"unicode/utf8"
)
const KEY = "SCU29489T427d03a73594b376a9471a70fc9c23555b4b2b37d718d"
// 注册
func SignUp(c echo.Context) error {
	cc := &util.ApiContext{c}
	username := cc.DefaultFormValueString("username", "", 1)
	password := cc.DefaultFormValueString("password", "", 1)
	rePassword := cc.DefaultFormValueString("re_password", "", 1)
	if username == "" {
		return cc.Fail(errors.New("请输入用户名"))
	}
	if !regexp.MustCompile("^[a-zA-Z0-9_-]{4,16}$").MatchString(username) {
		return cc.Fail(errors.New("用户名只能包含字母，数字和下划线，至少4个字符，最多16个字符"))
	}
	if password == "" {
		return cc.Fail(errors.New("请输入密码"))
	}
	if utf8.RuneCountInString(password) < 6 {
		return cc.Fail(errors.New("密码至少6个字符"))
	}
	if rePassword == "" {
		return cc.Fail(errors.New("请确认密码"))
	}
	if password != rePassword {
		return cc.Fail(errors.New("两次密码输入不一致"))
	}
	// 判断用户是否存在
	u, err := user.M.FindOne(bson.M{"username": username}, nil)
	if err != nil {
		return cc.Fail(err)
	}
	if u != nil {
		return cc.Fail(errors.New("用户名已存在"))
	}
	var Avatar = ""
	// 获取用户头像
	file, err := c.FormFile("avatar")
	// 如果上传了文件 但是出错
	if err != nil && err != http.ErrMissingFile {
		return cc.Fail(errors.New("上传头像发生错误"))
	}
	// 如果拿到上传的头像
	if file != nil {
		src, err := file.Open()
		defer src.Close()
		if err != nil {
			return cc.Fail(errors.New("上传头像发生错误"))
		}
		fp, err := util.SaveUploadFile(src, file.Filename)
		if err != nil {
			return cc.Fail(errors.New("保存头像发生错误"))
		}
		Avatar = fp
	}
	objectId := bson.NewObjectId()
	// 生成token
	token, err := util.GenerateTheToken(objectId, "client")
	if err != nil {
		return cc.Fail(err)
	}
	uu := bson.M{
		"_id":        objectId,
		"nickname":   "",
		"username":   username,
		"password":   *util.GenerateThePassword(password),
		"email":      "",
		"avatar":     Avatar,
		"integral":   0,
		"level":      1,
		"token":      token,
		"created_at": time.Now(),
	}
	_, err = user.M.Insert(uu)
	if err != nil {
		return cc.Fail(err)
	}
	_username := uu["username"].(string)
	//t := time.Now().Format("2006-01-02 15:04:05")
	http.Get("https://sc.ftqq.com/" + KEY + ".send?text=有人注册了&desp=" + _username + "刚刚注册了。")

	return cc.Success(&echo.Map{
		"token":    token,
		"username": _username,
		"avatar":   uu["avatar"],
	})
}

// 登录
func SignIn(c echo.Context) error {
	cc := &util.ApiContext{c}
	username := cc.DefaultFormValueString("username", "", 1)
	password := cc.DefaultFormValueString("password", "", 1)
	if username == "" {
		return cc.Fail(errors.New("请输入用户名"))
	}
	if !regexp.MustCompile("^[a-zA-Z0-9_-]{4,16}$").MatchString(username) {
		return cc.Fail(errors.New("用户名只能包含字母，数字和下划线，至少4个字符，最多16个字符"))
	}
	if password == "" {
		return cc.Fail(errors.New("请输入密码"))
	}
	if utf8.RuneCountInString(password) < 6 {
		return cc.Fail(errors.New("密码至少6个字符"))
	}
	// 判断用户是否存在
	u, err := user.M.FindOne(bson.M{
		"username": username,
		"password": *util.GenerateThePassword(password),
	}, nil)
	if err != nil {
		return cc.Fail(err)
	}
	if u == nil {
		return cc.Fail(errors.New("用户名或密码错误"))
	}
	// 生成token
	token, err := util.GenerateTheToken(u["_id"], "client")
	if err != nil {
		return cc.Fail(err)
	}
	// 更新用户token
	_, err = user.M.UpdateById(u["_id"], bson.M{
		"token": token,
	})
	if err != nil {
		return cc.Fail(err)
	}
	_username := u["username"].(string)
	//t := time.Now().Format("2006-01-02 15:04:05")
	http.Get("https://sc.ftqq.com/" + KEY + ".send?text=有人登陆了&desp=" + _username + "刚刚登陆了。")
	return cc.Success(&echo.Map{
		"token":    token,
		"username": _username,
	})
}

// 个人详情
func Detail(c echo.Context) error {
	cc := util.ApiContext{c}
	id := cc.Get("user").(*jwt.ClienJwtClaims).Id
	u, err := user.M.FindById(bson.ObjectIdHex(id), nil)
	if err != nil {
		return cc.Fail(err)
	}
	// 移除敏感字段
	delete(u, "password")
	delete(u, "token")
	return cc.Success(u)
}

// 注销登录
func SignOut(c echo.Context) error {
	cc := util.ApiContext{c}
	id := cc.Get("user").(*jwt.ClienJwtClaims).Id
	// 更新用户token
	_, err := user.M.UpdateById(id, bson.M{
		"token": "",
	})
	if err != nil {
		return cc.Fail(err)
	}
	return cc.Success(echo.Map{
		"token":    "",
		"username": "",
	})
}
