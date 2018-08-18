package util

import (
	"github.com/globalsign/mgo/bson"
)

/**
	query
	{"title": "精确查询", "key": "query", "value": "1"},
  	{"title": "模糊查询", "key": "query", "value": "2"},

	sort
	{"title": "最近收录", "key": "sort", "value": "1"},
    {"title": "最近上映", "key": "sort", "value": "2"},
    {"title": "最多播放", "key": "sort", "value": "3"},

	pid

	source
	{"title": "不限来源", "key": "source", "value": ""},
    {"title": "最大资源网", "key": "source", "value": "zuidazy"},
    {"title": "酷云资源网", "key": "source", "value": "kuyunzy"},

	year

	{"title": "不限年代", "key": "year", "value": ""},
	  {"title": "2018", "key": "year", "value": "2018"},
	  {"title": "2017", "key": "year", "value": "2017"},
	  {"title": "2016", "key": "year", "value": "2016"},
	  {"title": "2015", "key": "year", "value": "2015"},
	  {"title": "2014", "key": "year", "value": "2014"},
	  {"title": "2013", "key": "year", "value": "2013"},
	  {"title": "2012", "key": "year", "value": "2012"},
	  {"title": "2011", "key": "year", "value": "2011"},
	  {"title": "2010", "key": "year", "value": "2010"},
	  {"title": "00年代", "key": "year", "value": "00"},
	  {"title": "90年代", "key": "year", "value": "90"},
	  {"title": "80年代", "key": "year", "value": "80"},
	  {"title": "70年代", "key": "year", "value": "70"},
	  {"title": "更早", "key": "year", "value": "更早"},

	area
	{"title": "不限地区", "key": "area", "value": ""},
  {"title": "大陆", "key": "area", "value": "大陆"},
  {"title": "香港", "key": "area", "value": "香港"},
  {"title": "台湾", "key": "area", "value": "台湾"},
  {"title": "日本", "key": "area", "value": "日本"},
  {"title": "韩国", "key": "area", "value": "韩国"},
  {"title": "美国", "key": "area", "value": "美国"},
  {"title": "法国", "key": "area", "value": "法国"},
  {"title": "德国", "key": "area", "value": "德国"},
  {"title": "英国", "key": "area", "value": "英国"},
  {"title": "其他", "key": "area", "value": "其他"},

 */
var years = bson.M{
	"2018": "2018",
	"2017": "2017",
	"2016": "2016",
	"2015": "2015",
	"2014": "2014",
	"2013": "2013",
	"2012": "2012",
	"2011": "2011",
	"2010": "2010",
	"00": bson.M{
		"$gte": "2000", // >= 2000
		"$lte": "2010", // <= 2010
	},
	"90": bson.M{
		"$gte": "1900", // >= 1900
		"$lte": "1999", // <= 1999
	},
	"80": bson.M{
		"$gte": "1800", // >= 1800
		"$lte": "1899", // <= 1899
	},
	"70": bson.M{
		"$gte": "1700", // >= 1700
		"$lte": "1799", // <= 1799
	},
	"更早": bson.M{
		"$lte": "1699", // <= 1699
	},
}
var areas = bson.M{
	"大陆": "大陆",
	"香港": "香港",
	"台湾": "台湾",
	"日本": "日本",
	"韩国": "韩国",
	"美国": "美国",
	"法国": "法国",
	"德国": "德国",
	"英国": "英国",
	"其他": bson.M{
		"$nin": []string{"大陆","香港","台湾","日本","韩国","美国","法国","德国","英国",},
	},
}
func ParseQueryString (c *ApiContext) bson.M {
	ret := bson.M{}
	query := c.DefaultQueryString("query", "2", 1)
	sort := c.DefaultQueryString("sort", "1", 1)
	pid := c.DefaultQueryString("pid", "", 1)
	source := c.DefaultQueryString("source", "", 1)
	year := c.DefaultQueryString("year", "", 1)
	area := c.DefaultQueryString("area", "", 1)
	ret["query"] = query
	if sort == "1" {
		ret["sort"] = "-generated_at" // 收录时间降序
	} else if sort == "2" {
		ret["sort"] = "-released_at"  // 上映时间降序
	} else if sort == "3" {
		ret["sort"] = "-number"		  // 浏览次数降序
	}

	if pid != "none" && pid != "" {	  // 检查pid 历史原因保留none
		if bson.IsObjectIdHex(pid) {
			ret["pid"] = bson.ObjectIdHex(pid)
		}
	}
	if source == "zuidazy" {
		ret["source"] = "zuidazy"    // 最大资源网
	} else if source == "kuyunzy" {
		ret["source"] = "kuyunzy"	 // 酷云资源网
	}
	y := years[year]	// 解析年份
	if y != nil {
		ret["released_at"] = y
	}

	a := areas[area]	// 解析年份
	if a != nil {
		ret["region"] = a
	}
	return ret
}
