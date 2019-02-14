package model_menu_role

import (
	"encoding/json"
	"go-eladmin/core/mongo"
	"go-eladmin/model/a_mongo_index"
	"go-eladmin/model/shareDB"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	menuRoleCollection = mongo_index.CollectionMenuRole
)

type MenuRole struct {
	Id         int64 `json:"id,omitempty" bson:"_id,omitempty"`
	MenuId     int64 `json:"menu_id" bson:"menu_id"`
	RoleId     int64 `json:"role_id" bson:"role_id"`
	CreateTime int64 `json:"create_time" bson:"create_time"`
}

func (m MenuRole) ToJson() string {
	js, _ := json.Marshal(m)
	return string(js)
}

func (m MenuRole) isExist(query interface{}) bool {
	return mongo.IsExist(shareDB.DocManagerDBName(), menuRoleCollection, query)
}

func (m MenuRole) insert(docs ...interface{}) error {
	return mongo.Insert(shareDB.DocManagerDBName(), menuRoleCollection, docs...)
}

func (m MenuRole) update(selector, update interface{}) error {
	return mongo.Update(shareDB.DocManagerDBName(), menuRoleCollection, selector, update, true)
}

func (m MenuRole) findOne(query, selector interface{}) (MenuRole, error) {
	ap := MenuRole{}
	err := mongo.FindOne(shareDB.DocManagerDBName(), menuRoleCollection, query, selector, &ap)
	return ap, err
}
func (m MenuRole) findAll(query, selector interface{}) (results []MenuRole, err error) {
	results = []MenuRole{}
	err = mongo.FindAll(shareDB.DocManagerDBName(), menuRoleCollection, query, selector, &results)
	return results, err
}

func (m MenuRole) remove(selector interface{}) error {
	return mongo.Remove(shareDB.DocManagerDBName(), menuRoleCollection, selector)
}

func (m MenuRole) removeAll(selector interface{}) error {
	return mongo.RemoveAll(shareDB.DocManagerDBName(), menuRoleCollection, selector)
}

func (m MenuRole) totalCount(query, selector interface{}) (int, error) {
	return mongo.TotalCount(shareDB.DocManagerDBName(), menuRoleCollection, query, selector)
}

func (m MenuRole) findPage(page, limit int, query, selector interface{}, fields ...string) (results []MenuRole, err error) {
	results = []MenuRole{}
	err = mongo.FindPage(shareDB.DocManagerDBName(), menuRoleCollection, page, limit, query, selector, &results, fields...)
	return
}

func (m MenuRole) FindOne() (role MenuRole, err error) {
	role, err = m.findOne(bson.M{"_id": m.Id}, nil)
	return
}

func (m MenuRole) FindAll(query, selector interface{}) (results []MenuRole, err error) {
	results = []MenuRole{}
	return m.findAll(query, selector)
}

func (m MenuRole) Exist(query interface{}) bool {
	return m.isExist(query)
}

func (m MenuRole) Insert() error {
	m.Id, _ = mongo.GetIncrementId(shareDB.DocManagerDBName(), menuRoleCollection)
	m.CreateTime = time.Now().Unix()
	return m.insert(m)
}

func (m MenuRole) Update() error {
	return m.update(bson.M{"_id": m.Id}, m)
}

func (m MenuRole) Remove() error {
	return m.remove(bson.M{"_id": m.Id})
}

func (m MenuRole) RemoveMenuId(menuId int64) error {
	return m.removeAll(bson.M{"menu_id": menuId})
}
func (m MenuRole) RemoveRoleId(roleId int64) error {
	return m.removeAll(bson.M{"role_id": roleId})
}

func (m MenuRole) TotalCount(query, selector interface{}) (int, error) {
	return m.totalCount(query, selector)
}
func (m MenuRole) FindPageFilter(page, limit int, query, selector interface{}, fields ...string) ([]MenuRole, error) {
	return m.findPage(page, limit, query, selector, fields...)
}
