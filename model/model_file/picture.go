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
	pictureCollection = mongo_index.CollectionPicture
)

type Picture struct {
	Id         int64  `json:"id,omitempty" bson:"_id,omitempty"`
	CreateTime int64  `json:"create_time,omitempty" bson:"create_time,omitempty"`
	Height     int    `json:"height,omitempty" bson:"height,omitempty"`
	Width      int    `json:"width,omitempty" bson:"width,omitempty"`
	Size       int64  `json:"size,omitempty" bson:"size,omitempty"`  //文件大小
	Md5        string `json:"md5,omitempty" bson:"md5,omitempty"`    //文件md5值
	Path       string `json:"path,omitempty"  bson:"path,omitempty"` //本地路径
	Ext        string `json:"ext,omitempty" bson:"ext,omitempty"`    //文件后缀

}

func (p Picture) ToJson() string {
	js, _ := json.Marshal(p)
	return string(js)
}

func (p Picture) isExist(query interface{}) bool {
	return mongo.IsExist(shareDB.DocManagerDBName(), pictureCollection, query)
}

func (p Picture) insert(docs ...interface{}) error {
	return mongo.Insert(shareDB.DocManagerDBName(), pictureCollection, docs...)
}

func (p Picture) update(selector, update interface{}) error {
	return mongo.Update(shareDB.DocManagerDBName(), pictureCollection, selector, update, true)
}

func (p Picture) findOne(query, selector interface{}) (Picture, error) {
	ap := Picture{}
	err := mongo.FindOne(shareDB.DocManagerDBName(), pictureCollection, query, selector, &ap)
	return ap, err
}
func (p Picture) findAll(query, selector interface{}) (results []Picture, err error) {
	results = []Picture{}
	err = mongo.FindAll(shareDB.DocManagerDBName(), pictureCollection, query, selector, &results)
	return results, err
}

func (p Picture) remove(selector interface{}) error {
	return mongo.Remove(shareDB.DocManagerDBName(), pictureCollection, selector)
}

func (p Picture) removeAll(selector interface{}) error {
	return mongo.RemoveAll(shareDB.DocManagerDBName(), pictureCollection, selector)
}

func (p Picture) totalCount(query, selector interface{}) (int, error) {
	return mongo.TotalCount(shareDB.DocManagerDBName(), pictureCollection, query, selector)
}

func (p Picture) findPage(page, limit int, query, selector interface{}, fields ...string) (results []Picture, err error) {
	results = []Picture{}
	err = mongo.FindPage(shareDB.DocManagerDBName(), pictureCollection, page, limit, query, selector, &results, fields...)
	return
}

func (p Picture) pipeAll(pipeline, result interface{}, allowDiskUse bool) error {
	return mongo.PipeAll(shareDB.DocManagerDBName(), pictureCollection, pipeline, result, allowDiskUse)
}

func (p Picture) pipeOne(pipeline, result interface{}, allowDiskUse bool) error {
	return mongo.PipeOne(shareDB.DocManagerDBName(), pictureCollection, pipeline, result, allowDiskUse)
}

func (p Picture) explain(pipeline, result interface{}) (results []Picture, err error) {
	err = mongo.Explain(shareDB.DocManagerDBName(), pictureCollection, pipeline, result)
	return
}

func (p Picture) Exist() bool {
	return p.isExist(bson.M{"_id": p.Id})
}

func (p Picture) Insert() (id int64, err error) {
	if len(p.Path) == 0 {
		return -1, fmt.Errorf("path can not be nill")
	}
	id, _ = mongo.GetIncrementId(shareDB.DocManagerDBName(), pictureCollection)
	p.Id = id
	err = p.insert(p)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (p Picture) FindOne() (Picture, error) {
	return p.findOne(bson.M{"md5": p.Md5}, nil)
}

func (p Picture) Remove() error {
	return p.remove(bson.M{"md5": p.Md5})
}
