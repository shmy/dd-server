package service

import "github.com/globalsign/mgo/bson"

// 排除的分类
var RuleOut = []bson.ObjectId {
	bson.ObjectIdHex("5b6bd55a50456c5fb99610f5"), // 伦理片
	bson.ObjectIdHex("5b6c1f84adcfce70593225a9"), // 福利片
}

func IsInRuleOut (target interface{}) bool {
	for _, v := range RuleOut {
		if v == target {
			return true
		}
	}
	return false
}