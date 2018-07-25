package service

import (
	"github.com/shmy/dd-server/model/collection"
	"github.com/globalsign/mgo/bson"
	"fmt"
)

func CheckIsFavorited (uid interface{}, vid interface{}) bool {
	isFavorited := false
	count, err := collection.M.Count(bson.M{
		"_uid": uid,
		"_vid": vid,
	})
	fmt.Println(count, uid, vid)
	if err != nil {
		return isFavorited
	}
	if count != 0 {
		isFavorited = true
	}
	return isFavorited
}
