package mongo

import (
	"github.com/flywithbug/log4go"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func getMgoSession(db string) *mgo.Session {
	tm, ok := sessionMap[db]
	if !ok {
		log4go.Info("db error:%s", db)
		return nil
	}
	s := tm.session.Copy()
	return s
}

func connect(db, collection string) (*mgo.Session, *mgo.Collection) {
	s := getMgoSession(db)
	c := s.DB(db).C(collection)
	return s, c
}

func Collection(db, collection string) (*mgo.Session, *mgo.Collection) {
	s := getMgoSession(db)
	c := s.DB(db).C(collection)
	return s, c
}

func getDb(db string) (*mgo.Session, *mgo.Database) {
	tm, ok := sessionMap[db]
	if !ok {
		return nil, nil
	}
	ms := tm.session.Copy()
	return ms, ms.DB(db)
}

func Insert(db, collection string, docs ...interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	err := c.Insert(docs...)
	return err
}

func IsExist(db, collection string, query interface{}) bool {
	ms, c := connect(db, collection)
	defer ms.Close()
	count, _ := c.Find(query).Count()
	return count > 0
}

func FindOne(db, collection string, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Select(selector).One(result)
}

func FindAll(db, collection string, query, selector, results interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Select(selector).All(results)
}

func TotalCount(db, collection string, query, selector interface{}) (int, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Select(selector).Count()
}

//filterNull: 过滤空值对象。空值对象不更新
func Update(db, collection string, selector, update interface{}, filterNull bool) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	if filterNull {
		update = bson.M{"$set": update}
	}
	return c.Update(selector, update)
}

/*
selector := bson.M{"name": "Tom"}
data := bson.M{"$set": bson.M{"age": 22}}
filterNull: 过滤空值对象。空值对象不更新
*/
func UpdateAll(db, collection string, selector, data interface{}, filterNull bool) (changInfo *mgo.ChangeInfo, err error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	if filterNull {
		data = bson.M{"$set": data}
	}
	changInfo, err = c.UpdateAll(selector, data)
	return
}

func Remove(db, collection string, selector interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Remove(selector)
}

func RemoveAll(db, collection string, selector interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	_, err := c.RemoveAll(selector)
	return err
}

func FindPage(db, collection string, page, limit int, query, selector, results interface{}, fields ...string) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Select(selector).Sort(fields...).Skip(page * limit).Limit(limit).All(results)
}

func FindIter(db, collection string, query interface{}) *mgo.Iter {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Iter()
}

//insert one or multi documents
func BulkInsert(db, collection string, docs ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.Insert(docs...)
	return bulk.Run()
}

func BulkRemove(db, collection string, selector ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(db, collection)
	defer ms.Close()

	bulk := c.Bulk()
	bulk.Remove(selector...)
	return bulk.Run()
}

func BulkRemoveAll(db, collection string, selector ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.RemoveAll(selector...)
	return bulk.Run()
}

func BulkUpdate(db, collection string, pairs ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.Update(pairs...)
	return bulk.Run()
}

func BulkUpdateAll(db, collection string, pairs ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.UpdateAll(pairs...)
	return bulk.Run()
}

func BulkUpsert(db, collection string, pairs ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.Upsert(pairs...)
	return bulk.Run()
}

func PipeAll(db, collection string, pipeline, result interface{}, allowDiskUse bool) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	var pipe *mgo.Pipe
	if allowDiskUse {
		pipe = c.Pipe(pipeline).AllowDiskUse()
	} else {
		pipe = c.Pipe(pipeline)
	}
	return pipe.All(result)
}

func PipeOne(db, collection string, pipeline, result interface{}, allowDiskUse bool) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	var pipe *mgo.Pipe
	if allowDiskUse {
		pipe = c.Pipe(pipeline).AllowDiskUse()
	} else {
		pipe = c.Pipe(pipeline)
	}
	return pipe.One(result)
}

func PipeIter(db, collection string, pipeline interface{}, allowDiskUse bool) *mgo.Iter {
	ms, c := connect(db, collection)
	defer ms.Close()
	var pipe *mgo.Pipe
	if allowDiskUse {
		pipe = c.Pipe(pipeline).AllowDiskUse()
	} else {
		pipe = c.Pipe(pipeline)
	}

	return pipe.Iter()

}

func Explain(db, collection string, pipeline, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	pipe := c.Pipe(pipeline)
	return pipe.Explain(result)
}
func GridFSCreate(db, prefix, name string) (*mgo.GridFile, error) {
	ms, d := getDb(db)
	defer ms.Close()
	gridFs := d.GridFS(prefix)
	return gridFs.Create(name)
}

func GridFSFindOne(db, prefix string, query, result interface{}) error {
	ms, d := getDb(db)
	defer ms.Close()
	gridFs := d.GridFS(prefix)
	return gridFs.Find(query).One(result)
}

func GridFSFindAll(db, prefix string, query, result interface{}) error {
	ms, d := getDb(db)
	defer ms.Close()
	gridFs := d.GridFS(prefix)
	return gridFs.Find(query).All(result)
}

func GridFSOpen(db, prefix, name string) (*mgo.GridFile, error) {
	ms, d := getDb(db)
	defer ms.Close()
	gridFs := d.GridFS(prefix)
	return gridFs.Open(name)
}

func GridFSRemove(db, prefix, name string) error {
	ms, d := getDb(db)
	defer ms.Close()
	gridFs := d.GridFS(prefix)
	return gridFs.Remove(name)
}
