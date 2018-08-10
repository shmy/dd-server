package service

import (
	"sync"
	"github.com/globalsign/mgo/bson"
	"github.com/shmy/dd-server/model/video"
)

func GetHotsThumbnail (list []bson.M) {
	wg := sync.WaitGroup{}
	for _, v := range list {
		wg.Add(1)
		// 并发的关联
		go func(v bson.M) {
			defer wg.Done()
			v["video"], _ = video.M.FindById(v["vid"], "thumbnail, source")
		}(v)
	}
	wg.Wait()
}