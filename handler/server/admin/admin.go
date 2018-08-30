package admin

import (
	"errors"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/model/admin"
	"github.com/shmy/dd-server/util"
	"regexp"
	"unicode/utf8"
)

// 登录
func SignIn(c echo.Context) error {
	cc := &util.ApiContext{c}
	_username := cc.GetJSONValue("username")
	_password := cc.GetJSONValue("password")
	if _username == nil {
		return cc.Fail(errors.New("请输入用户名"))
	}
	username := _username.(string)
	if !regexp.MustCompile("^[a-zA-Z0-9_-]{4,16}$").MatchString(username) {
		return cc.Fail(errors.New("用户名只能包含字母，数字和下划线，至少4个字符，最多16个字符"))
	}
	if _password == nil {
		return cc.Fail(errors.New("请输入密码"))
	}
	password := _password.(string)
	if utf8.RuneCountInString(password) < 6 {
		return cc.Fail(errors.New("密码至少6个字符"))
	}
	// 判断用户是否存在
	u, err := admin.M.FindOne(bson.M{
		"username": username,
		"password": *util.GenerateThePasswordWithMD5(password),
	}, nil)
	if err != nil {
		return cc.Fail(err)
	}
	if u == nil {
		return cc.Fail(errors.New("用户名或密码错误"))
	}
	// 生成token
	token, err := util.GenerateTheToken(u["_id"], "admin")
	if err != nil {
		return cc.Fail(err)
	}
	// 更新用户token
	_, err = admin.M.UpdateById(u["_id"], bson.M{
		"token": token,
	})
	if err != nil {
		return cc.Fail(err)
	}
	//_username := u["username"].(string)
	return cc.Success(&echo.Map{
		"_id":      u["_id"],
		"token":    token,
		"username": u["username"],
	})
}
