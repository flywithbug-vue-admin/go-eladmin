package model_app_data_model

import (
	"encoding/json"
	"go-eladmin/core/mongo"
	"go-eladmin/model/a_mongo_index"
	"go-eladmin/model/shareDB"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//应用和数据模型关联

const (
	AppDataModelCollection = mongo_index.CollectionAppDataModel
)

type AppDataModel struct {
	Id         int64 `json:"id,omitempty" bson:"_id,omitempty"`
	ModelId    int64 `json:"model_id,omitempty" bson:"model_id,omitempty"`
	AppId      int64 `json:"app_id,omitempty" bson:"app_id,omitempty"`
	CreateTime int64 `json:"create_time,omitempty" bson:"create_time,omitempty"`
}

func (a AppDataModel) ToJson() string {
	js, _ := json.Marshal(a)
	return string(js)
}

func (a AppDataModel) isExist(query interface{}) bool {
	return mongo.IsExist(shareDB.DocManagerDBName(), AppDataModelCollection, query)
}

func (a AppDataModel) insert(docs ...interface{}) error {
	return mongo.Insert(shareDB.DocManagerDBName(), AppDataModelCollection, docs...)
}

func (a AppDataModel) update(selector, update interface{}) error {
	return mongo.Update(shareDB.DocManagerDBName(), AppDataModelCollection, selector, update, true)
}

func (a AppDataModel) findOne(query, selector interface{}) (AppDataModel, error) {
	ap := AppDataModel{}
	err := mongo.FindOne(shareDB.DocManagerDBName(), AppDataModelCollection, query, selector, &ap)
	return ap, err
}
func (a AppDataModel) findAll(query, selector interface{}) (results []AppDataModel, err error) {
	results = []AppDataModel{}
	err = mongo.FindAll(shareDB.DocManagerDBName(), AppDataModelCollection, query, selector, &results)
	return results, err
}

func (a AppDataModel) remove(selector interface{}) error {
	return mongo.Remove(shareDB.DocManagerDBName(), AppDataModelCollection, selector)
}

func (a AppDataModel) removeAll(selector interface{}) error {
	return mongo.RemoveAll(shareDB.DocManagerDBName(), AppDataModelCollection, selector)
}

func (a AppDataModel) totalCount(query, selector interface{}) (int, error) {
	return mongo.TotalCount(shareDB.DocManagerDBName(), AppDataModelCollection, query, selector)
}

func (a AppDataModel) findPage(page, limit int, query, selector interface{}, fields ...string) (results []AppDataModel, err error) {
	results = []AppDataModel{}
	err = mongo.FindPage(shareDB.DocManagerDBName(), AppDataModelCollection, page, limit, query, selector, &results, fields...)
	return
}

func (a AppDataModel) FindOne(query, selector interface{}) (role AppDataModel, err error) {
	role, err = a.findOne(query, selector)
	return
}
func (a AppDataModel) FindAll(query, selector interface{}) (results []AppDataModel, err error) {
	results = []AppDataModel{}
	return a.findAll(query, selector)
}

func (a AppDataModel) Exist(query interface{}) bool {
	return a.isExist(query)
}

func (a AppDataModel) Insert() error {
	a.Id, _ = mongo.GetIncrementId(shareDB.DocManagerDBName(), AppDataModelCollection)
	a.CreateTime = time.Now().Unix()
	return a.insert(a)
}

func (a AppDataModel) Update() error {
	return a.update(bson.M{"_id": a.Id}, a)
}

func (a AppDataModel) Remove() error {
	return a.remove(bson.M{"_id": a.Id})
}

func (a AppDataModel) RemoveModelId(modelId int64) error {
	return a.removeAll(bson.M{"model_id": modelId})
}

func (a AppDataModel) RemoveAppId(appId int64) error {
	return a.removeAll(bson.M{"app_id": appId})
}
