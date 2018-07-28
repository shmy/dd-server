package service

import (
	"github.com/shmy/dd-server/model/collection"
	"github.com/globalsign/mgo/bson"
	"fmt"
	"sync"
	"github.com/shmy/dd-server/model/video"
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

// 查询关联
func ListCollection(list []bson.M) {
	wg := sync.WaitGroup{}
	for _, v := range list {
		wg.Add(1)
		// 并发的关联
		go func(v bson.M) {
			defer wg.Done()
			v["video"], _ = video.M.FindById(v["_vid"],
				"name, latest, generated_at, thumbnail")
		}(v)
	}
	wg.Wait()
	//return list
}

// 查询关联
func ListCountCollection(list []bson.M) {
	wg := sync.WaitGroup{}
	for _, v := range list {
		wg.Add(1)
		// 并发的关联
		go func(v bson.M) {
			defer wg.Done()
			v["count"], _ = collection.M.Count(bson.M{
				"_fid": v["_id"],
			})
		}(v)
	}
	wg.Wait()
	//return list
}