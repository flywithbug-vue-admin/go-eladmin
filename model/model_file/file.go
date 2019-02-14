package model_file

import (
	"encoding/json"
	"fmt"
	"go-eladmin/core/mongo"
	"go-eladmin/model/a_mongo_index"
	"go-eladmin/model/shareDB"

	"gopkg.in/mgo.v2/bson"
)

const (
	fileCollection = mongo_index.CollectionFile
)

type File struct {
	Id   int64  `json:"id,omitempty" bson:"_id,omitempty"`
	Md5  string `json:"md5" bson:"md5"`                       //文件md5值
	Path string `json:"path"  bson:"path"`                    //本地路径
	Ext  string `json:"ext" bson:"ext"`                       //文件后缀
	Size int64  `json:"size,omitempty" bson:"size,omitempty"` //文件大小
}

func (f File) ToJson() string {
	js, _ := json.Marshal(f)
	return string(js)
}

func (f File) isExist(query interface{}) bool {
	return mongo.IsExist(shareDB.DocManagerDBName(), fileCollection, query)
}

func (f File) insert(docs ...interface{}) error {
	return mongo.Insert(shareDB.DocManagerDBName(), fileCollection, docs...)
}

func (f File) update(selector, update interface{}) error {
	return mongo.Update(shareDB.DocManagerDBName(), fileCollection, selector, update, true)
}

func (f File) findOne(query, selector interface{}) (File, error) {
	ap := File{}
	err := mongo.FindOne(shareDB.DocManagerDBName(), fileCollection, query, selector, &ap)
	return ap, err
}
func (f File) findAll(query, selector interface{}) (results []File, err error) {
	results = []File{}
	err = mongo.FindAll(shareDB.DocManagerDBName(), fileCollection, query, selector, &results)
	return results, err
}

func (f File) remove(selector interface{}) error {
	return mongo.Remove(shareDB.DocManagerDBName(), fileCollection, selector)
}

func (f File) removeAll(selector interface{}) error {
	return mongo.RemoveAll(shareDB.DocManagerDBName(), fileCollection, selector)
}

func (f File) totalCount(query, selector interface{}) (int, error) {
	return mongo.TotalCount(shareDB.DocManagerDBName(), fileCollection, query, selector)
}

func (f File) findPage(page, limit int, query, selector interface{}, fields ...string) (results []File, err error) {
	results = []File{}
	err = mongo.FindPage(shareDB.DocManagerDBName(), fileCollection, page, limit, query, selector, &results, fields...)
	return
}

func (f File) pipeAll(pipeline, result interface{}, allowDiskUse bool) error {
	return mongo.PipeAll(shareDB.DocManagerDBName(), fileCollection, pipeline, result, allowDiskUse)
}

func (f File) pipeOne(pipeline, result interface{}, allowDiskUse bool) error {
	return mongo.PipeOne(shareDB.DocManagerDBName(), fileCollection, pipeline, result, allowDiskUse)
}

func (f File) explain(pipeline, result interface{}) (results []File, err error) {
	err = mongo.Explain(shareDB.DocManagerDBName(), fileCollection, pipeline, result)
	return
}

func (f File) Exist() bool {
	return f.isExist(bson.M{"_id": f.Id})
}

func (f File) Insert() (id int64, err error) {
	if len(f.Path) == 0 {
		return -1, fmt.Errorf("path can not be nill")
	}
	id, _ = mongo.GetIncrementId(shareDB.DocManagerDBName(), fileCollection)
	f.Id = id
	err = f.insert(f)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (f File) FindOne() (File, error) {
	return f.findOne(bson.M{"md5": f.Md5}, nil)
}

func (f File) Remove() error {
	return f.remove(bson.M{"md5": f.Md5})
}
