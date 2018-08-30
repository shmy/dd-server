package util

import (
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"strings"
	"io/ioutil"
	"encoding/json"
	"github.com/globalsign/mgo/bson"
)

type ApiContext struct {
	echo.Context
}

func BodyToJson(c echo.Context) (map[string]interface{}, error) {
	s, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return nil, err
	}

	var body map[string]interface{}
	if err := json.Unmarshal(s, &body); err != nil {
		return nil, err
	}

	return body, nil
}

// 解析json
func (c *ApiContext) getJsonBody() map[string]interface{} {
	var j = c.Get("jsonData")
	if j == nil {
		var err error
		j, err = BodyToJson(c)
		if err != nil {
			return nil
		}
	}
	ret := j.(map[string]interface{})
	c.Set("jsonData", ret)
	return ret
}
//---------json-----------
// 获取所有json
func (c *ApiContext) GetJSONBody () map[string]interface{} {
	return c.getJsonBody()
}
// 按key获取json
func (c *ApiContext) GetJSONValue (key string) interface{} {
	r := c.getJsonBody()
	if r == nil {
		return nil
	}
	return r[key]
}
// string
func (c *ApiContext) GetJSONString (key string) *string {
	r := c.GetJSONValue(key)
	if v, ok := r.(string); ok {
		return &v
	}
	return nil
}
// []string
func (c *ApiContext) GetJSONStringList (key string) *[]string {
	r := c.GetJSONValue(key)
	if v, ok := r.([]string); ok {
		return &v
	}
	return nil
}
// int
func (c *ApiContext) GetJSONInt (key string) *int {
	r := c.GetJSONValue(key)
	if v, ok := r.(int); ok {
		return &v
	}
	return nil
}
// []int
func (c *ApiContext) GetJSONIntList (key string) *[]int {
	r := c.GetJSONValue(key)
	if v, ok := r.([]int); ok {
		return &v
	}
	return nil
}
// map
func (c *ApiContext) GetJSONMap (key string) *bson.M {
	r := c.GetJSONValue(key)
	if v, ok := r.(bson.M); ok {
		return &v
	}
	return nil
}
// []map
func (c *ApiContext) GetJSONMapList (key string) *[]bson.M {
	r := c.GetJSONValue(key)
	if v, ok := r.([]bson.M); ok {
		return &v
	}
	return nil
}
// bool
func (c *ApiContext) GetJSONBool (key string) *bool {
	r := c.GetJSONValue(key)
	if v, ok := r.(bool); ok {
		return &v
	}
	return nil
}
// []bool
func (c *ApiContext) GetJSONBoolList (key string) *[]bool {
	r := c.GetJSONValue(key)
	if v, ok := r.([]bool); ok {
		return &v
	}
	return nil
}

//---------json-----------
// 成功
func (c *ApiContext) Success(payload interface{}) error {
	return c.JSON(http.StatusOK, &echo.Map{
		"success": true,
		"message": "ok",
		"payload": payload,
	})
}

// 失败
func (c *ApiContext) Fail(err error) error {
	return c.JSON(http.StatusOK, &echo.Map{
		"success": false,
		"message": err.Error(),
		"payload": nil,
	})
}

// 从GET中获取参数转成Int 并提供默认值
func (c *ApiContext) DefaultQueryInt(key string, defaultValue int) int {
	p := c.QueryParam(key)
	r, err := strconv.Atoi(p)
	if err != nil {
		r = defaultValue
	}
	return r
}

// 从GET中获取参数转成String去除空格 并提供默认值
func (c *ApiContext) DefaultQueryString(key string, defaultValue string, trimSpace interface{}) string {
	r := c.QueryParam(key)
	if r == "" {
		r = defaultValue
	}
	if trimSpace != nil {
		return strings.TrimSpace(r)
	}
	return r
}

// 从form-data中获取参数转成String去除空格 并提供默认值
func (c *ApiContext) DefaultFormValueString(key string, defaultValue string, trimSpace interface{}) string {
	r := c.FormValue(key)
	if r == "" {
		r = defaultValue
	}
	if trimSpace != nil {
		return strings.TrimSpace(r)
	}
	return r
}
