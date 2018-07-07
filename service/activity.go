package service

import (
	"github.com/globalsign/mgo/bson"
	"github.com/shmy/dd-server/model/activity"
	"github.com/shmy/dd-server/model/user"
	"github.com/shmy/dd-server/model/video"
	"sync"
	"time"
)

// 插入播放记录到动态表
func AddToActivity(v bson.M, uid bson.ObjectId) error {
	ret, err := activity.M.FindOne(bson.M{"vid": v["_id"]}, nil)
	if err != nil {
		return err
	}

	// 不存在才插入
	if ret == nil {
		_, err = activity.M.Insert(bson.M{
			"_id":        bson.NewObjectId(),
			"vid":        v["_id"],
			"uid":        uid,
			"created_at": time.Now(),
			"updated_at": time.Now(),
		})
		if err != nil {
			return err
		}
	} else {
		_, err = activity.M.UpdateById(ret["_id"], bson.M{
			"updated_at": time.Now(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// 查询关联
func ListActivity(list []bson.M) []bson.M {
	wg := sync.WaitGroup{}
	for _, v := range list {
		wg.Add(1)
		// 并发的关联
		go func(v bson.M) {
			defer wg.Done()
			v["user"], _ = user.M.FindById(v["uid"],
				"username, nickname, avatar")
			v["video"], _ = video.M.FindById(v["vid"],
				"name, thumbnail")
		}(v)
	}
	wg.Wait()
	return list
}
