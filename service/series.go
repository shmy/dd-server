package service

import (
	"github.com/globalsign/mgo/bson"
	"sync"
	"github.com/shmy/dd-server/model/video"
)

func GetSeriesDetail (list []interface{}) {
	wg := sync.WaitGroup{}
	for _, v := range list {
		wg.Add(1)
		// 并发的关联
		go func(v bson.M) {
			defer wg.Done()
			v["video"], _ = video.M.FindById(v["vid"], "name, language, released_at, thumbnail, source")
		}(v.(bson.M))
	}
	wg.Wait()
}