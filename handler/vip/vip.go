package vip

import (
	"github.com/labstack/echo"
	"github.com/shmy/dd-server/util"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"strings"
	url2 "net/url"
)

func GetDetail (c echo.Context) error {
	cc := util.ApiContext{ c }
	url := cc.DefaultFormValueString("url", "", true)
	if url == "" {
		return cc.Fail(errors.New("请输入url地址"))
	}
	res := util.HttpGet("http://660e.com/?url=" + url)
	defer res.Close()
	if res == nil {
		return cc.Fail(errors.New("解析失败"))
	}
	doc, err := goquery.NewDocumentFromReader(res)
	if err != nil {
		return cc.Fail(errors.New("解析失败"))
	}
	src := doc.Find("#player").AttrOr("src", "")
	print(src)
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
