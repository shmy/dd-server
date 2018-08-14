package vip

import (
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/util"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"strings"
	url2 "net/url"
	"regexp"
)
//
//http://vip.youku.com/ajax/filter/filter_data?tag=10005&pl=30&pt=0&ar=0&mg=0&y=0&cl=1&o=0&pn=1
//
//{
//pn: 1 页码
//pl: 30 多少条
//pt: 1 资费
//mg: 0 类型
//ar: 0 区域索引
//y: 0 年份索引
//cl: 1 分类索引
//o: 0 貌似必填
//
//}
/**
10005 电影
10006 电视剧
10007 动漫
 */
const (
	classifyUrl = "http://vip.youku.com/ajax/filter/show_filter?tag=10005"
	listUrl = "http://vip.youku.com/ajax/filter/filter_data?tag=10006"
	jxSite = "https://jx.618g.com/"
	successCode = "20000"
)
// 获取优酷的分类数据
func GetClassifyList (c echo.Context) error {
	cc := util.ApiContext{ c }
	res := util.HttpGet(classifyUrl)
	defer res.Close()
	if res == nil {
		return cc.Fail(errors.New("请求失败"))
	}
	ret := util.ReadCloserToMap(res)
	if ret["code"] != successCode {
		return cc.Fail(errors.New("请求失败"))
	}
	return cc.Success(ret["result"])
}

// 获取分页列表
func GetList (c echo.Context) error {
	cc := util.ApiContext{ c }
	paramMap := map[string]string{
		"pn": "1",
		"pl": "20",
		"pt": "0",
		"mg": "0",
		"ar": "0",
		"y": "0",
		"cl": "1",
		"o": "0",
	}
	param := util.MapToQueryString(paramMap)
	param = listUrl + param

	res := util.HttpGet(param)
	defer res.Close()
	if res == nil {
		return cc.Fail(errors.New("请求失败"))
	}
	ret := util.ReadCloserToMap(res)
	if ret["code"] != successCode {
		return cc.Fail(errors.New("请求失败"))
	}
	result := ret["result"].(map[string]interface {})
	return cc.Success(result["result"])
}

// 根据播放网页获取播放地址
func GetDetail (c echo.Context) error {
	cc := util.ApiContext{ c }
	url := cc.DefaultFormValueString("url", "", true)
	if url == "" {
		return cc.Fail(errors.New("请输入url地址"))
	}
	res := util.HttpGet(jxSite + "?url=" + url)
	defer res.Close()
	if res == nil {
		return cc.Fail(errors.New("解析失败"))
	}
	doc, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		return cc.Fail(errors.New("解析失败"))
	}
	src := doc.Find("#player").AttrOr("src", "")
	//print(src)
	if !strings.HasSuffix(src, ".m3u8") {
		return cc.Fail(errors.New("解析失败"))
	}
	p, err := url2.Parse(src)
	if err != nil {
		return cc.Fail(errors.New("解析失败"))
	}
	m, err := url2.ParseQuery(p.RawQuery)
	if err != nil {
		return cc.Fail(errors.New("解析失败"))
	}
	return cc.Success(m["url"])
}

// 无法获取电视剧全部的播放列表
func GetPlayUrls (c echo.Context) error {
	cc := util.ApiContext{ c }
	url := cc.DefaultFormValueString("url", "", true)
	if url == "" {
		return cc.Fail(errors.New("请输入url地址"))
	}
	res := util.HttpGet(url)
	defer res.Close()
	if res == nil {
		return cc.Fail(errors.New("解析失败"))
	}
	doc, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		return cc.Fail(errors.New("解析失败"))
	}
	s, err := doc.Find("#bpmodule-playpage-playerdata-code").Html()
	if err != nil {
		return cc.Fail(errors.New("解析失败"))
	}
	reg := regexp.MustCompile(`window.playerAnthology= ([\s\S]*)</script>`)
	jsonStr := reg.FindStringSubmatch(s)[1]
	jsonMap := util.JsonStringToMap(jsonStr)
	//return cc.Success(util.JsonStringToMap(reg.FindStringSubmatch(s)[0]))
	return cc.Success(jsonMap["list"])
}