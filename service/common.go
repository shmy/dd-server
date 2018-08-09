package service

import (
	"github.com/globalsign/mgo/bson"
	"github.com/shmy/dd-server/model/classification"
	"reflect"
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
			var id = val["_id"].(bson.ObjectId)
			if !isContain(id, RuleOut) { // 过滤到不该显示的东西
				ids = append(ids, val["_id"])
			}
		}
	} else {
		ids = append(ids, id)
	}
	return ids, nil
}

// 判断obj是否在target中，target支持的类型arrary,slice,map
func isContain(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}