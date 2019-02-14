package model_user_role

import (
	"encoding/json"
	"go-eladmin/core/mongo"
	"go-eladmin/model/a_mongo_index"
	"go-eladmin/model/shareDB"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	UserRoleCollection = mongo_index.CollectionUserRole
)

type UserRole struct {
	Id         int64 `json:"id,omitempty" bson:"_id,omitempty"`
	UserId     int64 `json:"user_id,omitempty" bson:"user_id,omitempty"`
	RoleId     int64 `json:"role_id,omitempty" bson:"role_id,omitempty"`
	CreateTime int64 `json:"create_time,omitempty" bson:"create_time,omitempty"`
}

func (r UserRole) ToJson() string {
	js, _ := json.Marshal(r)
	return string(js)
}

func (r UserRole) isExist(query interface{}) bool {
	return mongo.IsExist(shareDB.DocManagerDBName(), UserRoleCollection, query)
}

func (r UserRole) insert(docs ...interface{}) error {
	return mongo.Insert(shareDB.DocManagerDBName(), UserRoleCollection, docs...)
}

func (r UserRole) update(selector, update interface{}) error {
	return mongo.Update(shareDB.DocManagerDBName(), UserRoleCollection, selector, update, true)
}

func (r UserRole) findOne(query, selector interface{}) (UserRole, error) {
	ap := UserRole{}
	err := mongo.FindOne(shareDB.DocManagerDBName(), UserRoleCollection, query, selector, &ap)
	return ap, err
}
func (r UserRole) findAll(query, selector interface{}) (results []UserRole, err error) {
	results = []UserRole{}
	err = mongo.FindAll(shareDB.DocManagerDBName(), UserRoleCollection, query, selector, &results)
	return results, err
}

func (r UserRole) remove(selector interface{}) error {
	return mongo.Remove(shareDB.DocManagerDBName(), UserRoleCollection, selector)
}

func (r UserRole) removeAll(selector interface{}) error {
	return mongo.RemoveAll(shareDB.DocManagerDBName(), UserRoleCollection, selector)
}

func (r UserRole) totalCount(query, selector interface{}) (int, error) {
	return mongo.TotalCount(shareDB.DocManagerDBName(), UserRoleCollection, query, selector)
}

func (r UserRole) findPage(page, limit int, query, selector interface{}, fields ...string) (results []UserRole, err error) {
	results = []UserRole{}
	err = mongo.FindPage(shareDB.DocManagerDBName(), UserRoleCollection, page, limit, query, selector, &results, fields...)
	return
}

func (r UserRole) FindOne(query, selector interface{}) (role UserRole, err error) {
	role, err = r.findOne(query, selector)
	return
}
func (r UserRole) FindAll(query, selector interface{}) (results []UserRole, err error) {
	results = []UserRole{}
	return r.findAll(query, selector)
}

func (r UserRole) Exist(query interface{}) bool {
	return r.isExist(query)
}

func (r UserRole) Insert() error {
	r.Id, _ = mongo.GetIncrementId(shareDB.DocManagerDBName(), UserRoleCollection)
	r.CreateTime = time.Now().Unix()
	return r.insert(r)
}

func (r UserRole) Update() error {
	return r.update(bson.M{"_id": r.Id}, r)
}

func (r UserRole) Remove() error {
	return r.remove(bson.M{"_id": r.Id})
}

func (r UserRole) RemoveUserId(userId int64) error {
	return r.removeAll(bson.M{"user_id": userId})
}
