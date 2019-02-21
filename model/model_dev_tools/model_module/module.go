package model_module

import (
	"encoding/json"
	"fmt"
	"go-eladmin/core/mongo"
	"go-eladmin/model/a_mongo_index"
	"go-eladmin/model/model_app"
	"go-eladmin/model/shareDB"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionModule = mongo_index.CollectionModule
)

const (
	DataModulePermissionALL    = "Model_ALL"
	DataModulePermissionCreate = "Model_CREATE"
	DataModulePermissionSelect = "Model_SELECT"
	DataModulePermissionEdit   = "Model_EDIT"
	DataModulePermissionDelete = "Model_DELETE"
)

//业务线 模块
type Module struct {
	Id         int64                   `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string                  `json:"name,omitempty" bson:"name,omitempty"`
	ApiDomain  string                  `json:"api_domain,omitempty" bson:"api_domain,omitempty"` //无域名地址 www.flywithme.top
	Path       string                  `json:"path,omitempty" bson:"path,omitempty"`             //业务线Path
	Desc       string                  `json:"desc,omitempty" bson:"desc,omitempty"`
	CreateTime int64                   `json:"create_time,omitempty" bson:"create_time,omitempty"`
	UpdateTime int64                   `json:"update_time,omitempty" bson:"update_time,omitempty"`
	Apps       []model_app.Application `json:"apps,omitempty" bson:"apps,omitempty"`
}

func (m Module) ToJson() string {
	js, _ := json.Marshal(m)
	return string(js)
}

func (m Module) isExist(query interface{}) bool {
	return mongo.IsExist(shareDB.DocManagerDBName(), CollectionModule, query)
}

func (m Module) insert(docs ...interface{}) error {
	return mongo.Insert(shareDB.DocManagerDBName(), CollectionModule, docs...)
}

func (m Module) update(selector, update interface{}) error {
	return mongo.Update(shareDB.DocManagerDBName(), CollectionModule, selector, update, true)
}

func (m Module) findOne(query, selector interface{}) (Module, error) {
	ap := Module{}
	err := mongo.FindOne(shareDB.DocManagerDBName(), CollectionModule, query, selector, &ap)
	return ap, err
}
func (m Module) findAll(query, selector interface{}) (results []Module, err error) {
	results = []Module{}
	err = mongo.FindAll(shareDB.DocManagerDBName(), CollectionModule, query, selector, &results)
	return results, err
}

func (m Module) remove(selector interface{}) error {
	return mongo.Remove(shareDB.DocManagerDBName(), CollectionModule, selector)
}

func (m Module) removeAll(selector interface{}) error {
	return mongo.RemoveAll(shareDB.DocManagerDBName(), CollectionModule, selector)
}

func (m Module) totalCount(query, selector interface{}) (int, error) {
	return mongo.TotalCount(shareDB.DocManagerDBName(), CollectionModule, query, selector)
}

func (m Module) findPage(page, limit int, query, selector interface{}, fields ...string) (results []Module, err error) {
	results = []Module{}
	err = mongo.FindPage(shareDB.DocManagerDBName(), CollectionModule, page, limit, query, selector, &results, fields...)
	return
}

func (m Module) FindOne() (Module, error) {
	return m.findOne(bson.M{"_id": m.Id}, nil)
}

func (m Module) FindAll(query, selector interface{}) (results []Module, err error) {
	results = []Module{}
	return m.findAll(query, selector)
}

func (m Module) Exist(query interface{}) bool {
	return m.isExist(query)
}

func (m Module) Insert() (int64, error) {
	m.Id, _ = mongo.GetIncrementId(shareDB.DocManagerDBName(), CollectionModule)
	m.CreateTime = time.Now().Unix()
	m.UpdateTime = m.CreateTime
	err := m.insert(m)
	return m.Id, err
}

func (m Module) fetchApps() (results []model_app.Application, err error) {
	results = make([]model_app.Application, 0)
	for _, item := range m.Apps {
		item, err = item.FindSimpleOne(nil, bson.M{"_id": 1, "name": 1})
		if err == nil {
			results = append(results, item)
		}
	}
	return
}

//func (d DataModel) AddAttribute(a Attribute) error {
//	if d.isExistAttribute(a) {
//		return fmt.Errorf("duplicate attribute name:%s", a.Name)
//	}
//	if err := checkNameReg(a.Name); err != nil {
//		return err
//	}
//	update := bson.M{"$addToSet": bson.M{"attributes": a}}
//	change := mgo.Change{
//		Update: update,
//	}
//	a.UpdateTime = time.Now().Unix()
//	query := bson.M{"_id": d.Id}
//	ms, c := mongo.Collection(shareDB.DocManagerDBName(), dataModelCollection)
//	defer ms.Close()
//	_, err := c.Find(query).Apply(change, nil)
//	return err
//}

func (m Module) isExistApp(a model_app.Application) bool {
	selector := bson.M{"_id": m.Id, "apps.id": a.Id}
	return m.isExist(selector)
}

func (m Module) AddAppRelation(app model_app.Application) error {
	if m.isExistApp(app) {
		return fmt.Errorf("app exist")
	}

	return nil
}
func (m Module) AddAppRelations(list []model_app.Application) error {

	return nil
}
func (m Module) RemoveRelation(app model_app.Application) error {

	return nil
}

func (m Module) RemoveRelations(list []model_app.Application) error {

	return nil
}

func (m Module) Update() error {
	return m.update(bson.M{"_id": m.Id}, m)
}

func (m Module) Remove() error {
	return m.remove(bson.M{"_id": m.Id})
}

func (m Module) TotalCount(query, selector interface{}) (int, error) {
	return m.totalCount(query, selector)
}

func (m Module) FindPageFilter(page, limit int, query, selector interface{}, fields ...string) ([]Module, error) {
	result, err := m.findPage(page, limit, query, selector, fields...)
	for index := range result {
		result[index].Apps, _ = result[index].fetchApps()
	}
	return result, err
}
