package model_app_manager

import (
	"encoding/json"
	"go-eladmin/core/mongo"
	"go-eladmin/model/a_mongo_index"
	"go-eladmin/model/shareDB"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	AppManagerCollection = mongo_index.CollectionAppManager
)

type AppManager struct {
	Id         int64 `json:"id,omitempty" bson:"_id,omitempty"`
	UserId     int64 `json:"user_id,omitempty" bson:"user_id,omitempty"`
	AppId      int64 `json:"app_id,omitempty" bson:"app_id,omitempty"`
	CreateTime int64 `json:"create_time,omitempty" bson:"create_time,omitempty"`
}

func (r AppManager) ToJson() string {
	js, _ := json.Marshal(r)
	return string(js)
}

func (r AppManager) isExist(query interface{}) bool {
	return mongo.IsExist(shareDB.DocManagerDBName(), AppManagerCollection, query)
}

func (r AppManager) insert(docs ...interface{}) error {
	return mongo.Insert(shareDB.DocManagerDBName(), AppManagerCollection, docs...)
}

func (r AppManager) update(selector, update interface{}) error {
	return mongo.Update(shareDB.DocManagerDBName(), AppManagerCollection, selector, update, true)
}

func (r AppManager) findOne(query, selector interface{}) (AppManager, error) {
	ap := AppManager{}
	err := mongo.FindOne(shareDB.DocManagerDBName(), AppManagerCollection, query, selector, &ap)
	return ap, err
}
func (r AppManager) findAll(query, selector interface{}) (results []AppManager, err error) {
	results = []AppManager{}
	err = mongo.FindAll(shareDB.DocManagerDBName(), AppManagerCollection, query, selector, &results)
	return results, err
}

func (r AppManager) remove(selector interface{}) error {
	return mongo.Remove(shareDB.DocManagerDBName(), AppManagerCollection, selector)
}

func (r AppManager) removeAll(selector interface{}) error {
	return mongo.RemoveAll(shareDB.DocManagerDBName(), AppManagerCollection, selector)
}

func (r AppManager) totalCount(query, selector interface{}) (int, error) {
	return mongo.TotalCount(shareDB.DocManagerDBName(), AppManagerCollection, query, selector)
}

func (r AppManager) findPage(page, limit int, query, selector interface{}, fields ...string) (results []AppManager, err error) {
	results = []AppManager{}
	err = mongo.FindPage(shareDB.DocManagerDBName(), AppManagerCollection, page, limit, query, selector, &results, fields...)
	return
}

func (r AppManager) FindOne(query, selector interface{}) (role AppManager, err error) {
	role, err = r.findOne(query, selector)
	return
}
func (r AppManager) FindAll(query, selector interface{}) (results []AppManager, err error) {
	results = []AppManager{}
	return r.findAll(query, selector)
}

func (r AppManager) Exist(query interface{}) bool {
	return r.isExist(query)
}

func (r AppManager) Insert() error {
	r.Id, _ = mongo.GetIncrementId(shareDB.DocManagerDBName(), AppManagerCollection)
	r.CreateTime = time.Now().Unix()
	return r.insert(r)
}

func (r AppManager) Update() error {
	return r.update(bson.M{"_id": r.Id}, r)
}

func (r AppManager) Remove() error {
	return r.remove(bson.M{"_id": r.Id})
}

func (r AppManager) RemoveUserId(userId int64) error {
	return r.removeAll(bson.M{"user_id": userId})
}

func (r AppManager) RemoveAppId(appId int64) error {
	return r.removeAll(bson.M{"app_id": appId})
}
