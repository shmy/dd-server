package model

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
	"github.com/shmy/dd-server/config"
	"github.com/spf13/viper"
	"strings"
	"time"
)

var Db *mgo.Database

func init() {
	config.Init("")
	dialInfo := &mgo.DialInfo{
		Addrs:     viper.GetStringSlice("mongodb.address"),
		Direct:    false,
		Timeout:   time.Second * 20, // 连接超时
		Source:    "admin",
		Username:  viper.GetString("mongodb.username"),
		Password:  viper.GetString("mongodb.password"),
		PoolLimit: 4096, // Session.SetPoolLimit
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatal("the database connection failed:", err)
		panic(err)
	}
	//defer session.Close()
	Db = session.DB(viper.GetString("mongodb.database"))
}

type Model struct {
	Collection *mgo.Collection
}

// 不支持泛型！！！！！ 伤不起！！！
/**
  解析选择器
  可以传入字符串以逗号分隔
  可以传入字符串数组
*/
func parserSelector(selector interface{}) bson.M {
	switch data := selector.(type) {
	case string:
		{
			strSlice := strings.Split(data, ",")
			ret := make(bson.M)
			for _, v := range strSlice {
				ret[strings.TrimSpace(v)] = 1
			}
			return ret
		}
	case []string:
		{
			ret := make(bson.M)
			for _, v := range data {
				ret[strings.TrimSpace(v)] = 1
			}
			return ret
		}
	default:

		return nil
	}
}

/**
  解析排序语法
  可以传入字符串以逗号分隔
  可以传入字符串数组
*/
func parserSort(sort interface{}) []string {
	switch data := sort.(type) {
	case []string:
		{
			return data
		}
	//case bson.M: {
	//	strSlice := make([]string, 0)
	//	for key, val := range data {
	//		if val != nil {
	//			key = "+" + key
	//		} else {
	//			key = "-" + key
	//		}
	//		strSlice = append(strSlice, key)
	//	}
	//	fmt.Println(strSlice)
	//	return strSlice
	//}
	case string:
		{
			return strings.Split(data, ",")
		}
	default:
		return nil
	}
}

func parserObjectId(id interface{}) bson.ObjectId {
	switch data := id.(type) {
	case bson.ObjectId:
		{
			return data
		}
	case string:
		{
			return bson.ObjectIdHex(data)
		}
	default:
		return ""
	}
}

// 条件 排序 分页查询
func (m *Model) Query(
	conditions bson.M,
	selector interface{},
	sort interface{},
	skip int,
	limit int,
) ([]bson.M, error) {
	ret := make([]bson.M, 0)
	selector = parserSelector(selector)
	sort = parserSort(sort)

	err := m.Collection.
		Find(conditions).
		Select(selector).
		Sort(sort.([]string)...).
		Skip(skip).
		Limit(limit).
		All(&ret)
	if err != nil {
		log.Error("DB Query:", err, lager.Data{
			"collection": m.Collection.FullName,
			"conditions": conditions,
			"selector":   selector,
			"sort":       sort,
			"skip":       skip,
			"limit":      limit,
		})
	}
	return ret, err
}

// 获取所有数据
func (m *Model) FindAll(selector interface{}) ([]bson.M, error) {
	return m.Find(nil, selector)
}

// 按条件查询多数据
func (m *Model) Find(conditions bson.M, selector interface{}) ([]bson.M, error) {
	ret := make([]bson.M, 0)
	query := m.Collection.Find(conditions)
	if selector != nil {
		selector = parserSelector(selector)
		query = query.Select(selector)
	}
	err := query.All(&ret)
	if err != nil {
		log.Error("DB Find:", err, lager.Data{
			"collection": m.Collection.FullName,
			"conditions": conditions,
		})
	}
	return ret, err
}

// 根据id查找一条数据
func (m *Model) FindById(objectId interface{}, selector interface{}) (bson.M, error) {
	var ret bson.M
	objectId = parserObjectId(objectId)
	query := m.Collection.FindId(objectId)
	if selector != nil {
		selector = parserSelector(selector)
		query = query.Select(selector)
	}
	err := query.One(&ret)
	if err != nil && err != mgo.ErrNotFound {
		log.Error("DB FindById:", err, lager.Data{
			"collection": m.Collection.FullName,
			"objectId":   objectId,
		})
		return nil, err
	}
	return ret, nil
}

// 根据条件查找查找一条数据
func (m *Model) FindOne(conditions bson.M, selector interface{}) (bson.M, error) {
	var ret bson.M
	query := m.Collection.Find(conditions)
	if selector != nil {
		selector = parserSelector(selector)
		query = query.Select(selector)
	}
	err := query.One(&ret)
	if err != nil && err != mgo.ErrNotFound {
		log.Error("DB FindOne:", err, lager.Data{
			"collection": m.Collection.FullName,
			"conditions": conditions,
		})
		return nil, err
	}
	return ret, nil
}

// 按条件计数
func (m *Model) Count(conditions bson.M) (int, error) {
	count, err := m.Collection.Find(&conditions).Count()
	if err != nil {
		log.Error("DB Count:", err, lager.Data{
			"collection": m.Collection.FullName,
			"conditions": conditions,
		})
	}
	return count, err
}

// 插入一条数据
func (m *Model) Insert(doc bson.M) (bson.M, error) {
	err := m.Collection.Insert(&doc)
	if err != nil {
		log.Error("DB Insert:", err, lager.Data{
			"collection": m.Collection.FullName,
			"document":   doc,
		})
	}
	return doc, err
}

// 插入多条数据
func (m *Model) BulkInsert(docs []bson.M) ([]bson.M, error) {
	bulk := m.Collection.Bulk()
	bulk.Unordered()
	bulk.Insert(&docs)
	_, err := bulk.Run()
	if err != nil {
		log.Error("DB BulkInsert:", err, lager.Data{
			"collection": m.Collection.FullName,
			"documents":  docs,
		})
	}
	return docs, err
}

// 更新或创建
func (m *Model) Upsert(conditions bson.M, doc bson.M) (bson.M, error) {
	_, err := m.Collection.Upsert(conditions, &doc)
	if err != nil {
		log.Error("DB Upsert:", err, lager.Data{
			"collection": m.Collection.FullName,
			"conditions": conditions,
			"document":   doc,
		})
	}
	return doc, err
}

// 更新一条数据
func (m *Model) Update(conditions bson.M, doc bson.M) (bson.M, error) {
	err := m.Collection.Update(conditions, &bson.M{
		"$set": doc,
	})
	if err != nil {
		log.Error("DB Update:", err, lager.Data{
			"collection": m.Collection.FullName,
			"conditions": conditions,
			"document":   doc,
		})
		return nil, err
	}
	return m.FindOne(conditions, nil)
}

// 按ID更新一条数据
func (m *Model) UpdateById(objectId interface{}, doc bson.M) (bson.M, error) {
	objectId = parserObjectId(objectId)
	err := m.Collection.UpdateId(objectId, &bson.M{
		"$set": doc,
	})
	if err != nil {
		log.Error("DB UpdateId:", err, lager.Data{
			"collection": m.Collection.FullName,
			"objectId":   objectId,
			"document":   doc,
		})
		return nil, err
	}
	return m.FindById(objectId, nil)
}

// 更新多条数据
func (m *Model) BulkUpdate(conditions bson.M, doc []bson.M) ([]bson.M, error) {
	_, err := m.Collection.UpdateAll(conditions, &bson.M{
		"$set": doc,
	})
	if err != nil {
		log.Error("DB UpdateAll:", err, lager.Data{
			"collection": m.Collection.FullName,
			"conditions": conditions,
			"document":   doc,
		})
		return nil, err
	}
	return m.Find(conditions, nil)
}

// 按id删除一个数据
func (m *Model) RemoveById(objectId bson.ObjectId) bool {
	err := m.Collection.RemoveId(objectId)
	if err != nil {
		log.Error("DB RemoveOne:", err, lager.Data{
			"collection": m.Collection.FullName,
			"objectId": objectId,
		})
		return false
	}
	return true
}

// 按条件删除一个数据
func (m *Model) RemoveOne(conditions bson.M) bool {
	err := m.Collection.Remove(conditions)
	if err != nil {
		log.Error("DB RemoveOne:", err, lager.Data{
			"collection": m.Collection.FullName,
			"conditions": conditions,
		})
		return false
	}
	return true
}

// 按条件删除多个数据
func (m *Model) RemoveAll(conditions bson.M) bool {
	_, err := m.Collection.RemoveAll(conditions)
	if err != nil {
		log.Error("DB RemoveOne:", err, lager.Data{
			"collection": m.Collection.FullName,
			"conditions": conditions,
		})
		return false
	}
	return true
}