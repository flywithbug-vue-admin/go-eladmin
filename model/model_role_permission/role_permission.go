package model_role_permission

import (
	"encoding/json"
	"go-eladmin/core/mongo"
	"go-eladmin/model/a_mongo_index"
	"go-eladmin/model/shareDB"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	rolePermissionCollection = mongo_index.CollectionRolePermission
)

type RolePermission struct {
	Id           int64 `json:"id,omitempty" bson:"_id,omitempty"`
	RoleId       int64 `json:"role_id" bson:"role_id"`
	PermissionId int64 `json:"permission_id" bson:"permission_id"`
	CreateTime   int64 `json:"create_time" bson:"create_time"`
}

func (r RolePermission) ToJson() string {
	js, _ := json.Marshal(r)
	return string(js)
}

func (r RolePermission) isExist(query interface{}) bool {
	return mongo.IsExist(shareDB.DocManagerDBName(), rolePermissionCollection, query)
}

func (r RolePermission) insert(docs ...interface{}) error {
	return mongo.Insert(shareDB.DocManagerDBName(), rolePermissionCollection, docs...)
}

func (r RolePermission) update(selector, update interface{}) error {
	return mongo.Update(shareDB.DocManagerDBName(), rolePermissionCollection, selector, update, true)
}

func (r RolePermission) findOne(query, selector interface{}) (RolePermission, error) {
	ap := RolePermission{}
	err := mongo.FindOne(shareDB.DocManagerDBName(), rolePermissionCollection, query, selector, &ap)
	return ap, err
}
func (r RolePermission) findAll(query, selector interface{}) (results []RolePermission, err error) {
	results = []RolePermission{}
	err = mongo.FindAll(shareDB.DocManagerDBName(), rolePermissionCollection, query, selector, &results)
	return results, err
}

func (r RolePermission) remove(selector interface{}) error {
	return mongo.Remove(shareDB.DocManagerDBName(), rolePermissionCollection, selector)
}

func (r RolePermission) removeAll(selector interface{}) error {
	return mongo.RemoveAll(shareDB.DocManagerDBName(), rolePermissionCollection, selector)
}

func (r RolePermission) totalCount(query, selector interface{}) (int, error) {
	return mongo.TotalCount(shareDB.DocManagerDBName(), rolePermissionCollection, query, selector)
}

func (r RolePermission) findPage(page, limit int, query, selector interface{}, fields ...string) (results []RolePermission, err error) {
	results = []RolePermission{}
	err = mongo.FindPage(shareDB.DocManagerDBName(), rolePermissionCollection, page, limit, query, selector, &results, fields...)
	return
}

func (r RolePermission) FindOne() (role RolePermission, err error) {
	role, err = r.findOne(bson.M{"_id": r.Id}, nil)
	return
}

func (r RolePermission) FindAll(query, selector interface{}) (results []RolePermission, err error) {
	results = []RolePermission{}
	return r.findAll(query, selector)
}

func (r RolePermission) Exist(query interface{}) bool {
	return r.isExist(query)
}

func (r RolePermission) Insert() error {
	r.Id, _ = mongo.GetIncrementId(shareDB.DocManagerDBName(), rolePermissionCollection)
	r.CreateTime = time.Now().Unix()
	return r.insert(r)
}

func (r RolePermission) Update() error {
	return r.update(bson.M{"_id": r.Id}, r)
}

func (r RolePermission) Remove() error {
	return r.remove(bson.M{"_id": r.Id})
}

func (r RolePermission) RemoveRoleId(roleId int64) error {
	return r.removeAll(bson.M{"role_id": roleId})
}

func (r RolePermission) TotalCount(query, selector interface{}) (int, error) {
	return r.totalCount(query, selector)
}
func (r RolePermission) FindPageFilter(page, limit int, query, selector interface{}, fields ...string) ([]RolePermission, error) {
	return r.findPage(page, limit, query, selector, fields...)
}
