package service

import (
	"github.com/globalsign/mgo/bson"
	"github.com/shmy/dd-server/model/classification"
)

// 获取子分类
func FindClassifyIds(id interface{}) ([]interface{}, error) {
	r, err := classification.M.
		Find(bson.M{"pid": id}, nil)
	if err != nil {
		return nil, err
	}
	var ids []interface{}
	if len(r) > 0 {
		for _, val := range r {
			ids = append(ids, val["_id"])
		}
	} else {
		ids = append(ids, id)
	}
	return ids, nil
}
