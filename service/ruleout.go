package service

import "github.com/globalsign/mgo/bson"

// 排除的分类
var RuleOut = []bson.ObjectId {
	bson.ObjectIdHex("5b6bd55a50456c5fb99610f5"), // 伦理片
	//bson.ObjectIdHex("5b0fd14e7cad175a34a2ea8c"): true,
}