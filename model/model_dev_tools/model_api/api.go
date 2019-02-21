package Api_api

import (
	"encoding/json"
	"go-eladmin/core/mongo"
	"go-eladmin/model/a_mongo_index"
	"go-eladmin/model/model_app"
	"go-eladmin/model/shareDB"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionApi = mongo_index.CollectionApi
)

const (
	ApiPermissionALL    = "Api_ALL"
	ApiPermissionCreate = "Api_CREATE"
	ApiPermissionSelect = "Api_SELECT"
	ApiPermissionEdit   = "Api_EDIT"
	ApiPermissionDelete = "Api_DELETE"
)

//业务线 模块
type Api struct {
	Id          int64                   `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string                  `json:"name,omitempty" bson:"name,omitempty"`
	Path        string                  `json:"path,omitempty" bson:"path,omitempty"` //业务请求Path
	Desc        string                  `json:"desc,omitempty" bson:"desc,omitempty"`
	CreateTime  int64                   `json:"create_time,omitempty" bson:"create_time,omitempty"`
	UpdateTime  int64                   `json:"update_time,omitempty" bson:"update_time,omitempty"`
	Apps        []model_app.Application `json:"apps,omitempty" bson:"apps,omitempty"`
	ResponseId  int64                   `json:"response_id,omitempty" bson:"response_id,omitempty"`   //返回模型Id
	ParameterId int64                   `json:"parameter_id,omitempty" bson:"parameter_id,omitempty"` //请求模型Id

	//数据生成时获取
	Response  interface{} `json:"response,omitempty" bson:"response,omitempty"`   //返回模型
	Parameter interface{} `json:"parameter,omitempty" bson:"parameter,omitempty"` //请求模型
}

func (m Api) ToJson() string {
	js, _ := json.Marshal(m)
	return string(js)
}

func (m Api) isExist(query interface{}) bool {
	return mongo.IsExist(shareDB.DocManagerDBName(), CollectionApi, query)
}

func (m Api) insert(docs ...interface{}) error {
	return mongo.Insert(shareDB.DocManagerDBName(), CollectionApi, docs...)
}

func (m Api) update(selector, update interface{}) error {
	return mongo.Update(shareDB.DocManagerDBName(), CollectionApi, selector, update, true)
}

func (m Api) findOne(query, selector interface{}) (Api, error) {
	ap := Api{}
	err := mongo.FindOne(shareDB.DocManagerDBName(), CollectionApi, query, selector, &ap)
	return ap, err
}
func (m Api) findAll(query, selector interface{}) (results []Api, err error) {
	results = []Api{}
	err = mongo.FindAll(shareDB.DocManagerDBName(), CollectionApi, query, selector, &results)
	return results, err
}

func (m Api) remove(selector interface{}) error {
	return mongo.Remove(shareDB.DocManagerDBName(), CollectionApi, selector)
}

func (m Api) removeAll(selector interface{}) error {
	return mongo.RemoveAll(shareDB.DocManagerDBName(), CollectionApi, selector)
}

func (m Api) totalCount(query, selector interface{}) (int, error) {
	return mongo.TotalCount(shareDB.DocManagerDBName(), CollectionApi, query, selector)
}

func (m Api) findPage(page, limit int, query, selector interface{}, fields ...string) (results []Api, err error) {
	results = []Api{}
	err = mongo.FindPage(shareDB.DocManagerDBName(), CollectionApi, page, limit, query, selector, &results, fields...)
	return
}

func (m Api) FindOne() (Api, error) {
	return m.findOne(bson.M{"_id": m.Id}, nil)
}

func (m Api) FindAll(query, selector interface{}) (results []Api, err error) {
	results = []Api{}
	return m.findAll(query, selector)
}

func (m Api) Exist(query interface{}) bool {
	return m.isExist(query)
}

func (m Api) Insert() (int64, error) {
	m.Id, _ = mongo.GetIncrementId(shareDB.DocManagerDBName(), CollectionApi)
	m.CreateTime = time.Now().Unix()
	m.UpdateTime = m.CreateTime
	err := m.insert(m)
	return m.Id, err
}

func (m Api) Update() error {
	return m.update(bson.M{"_id": m.Id}, m)
}

func (m Api) Remove() error {
	return m.remove(bson.M{"_id": m.Id})
}

func (m Api) TotalCount(query, selector interface{}) (int, error) {
	return m.totalCount(query, selector)
}

func (m Api) FindPageFilter(page, limit int, query, selector interface{}, fields ...string) ([]Api, error) {
	result, err := m.findPage(page, limit, query, selector, fields...)
	for index := range result {
		result[index].Apps, _ = result[index].fetchApps()
	}
	return result, err
}
