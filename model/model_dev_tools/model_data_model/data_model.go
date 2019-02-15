package model_data_model

import (
	"encoding/json"
	"fmt"
	"go-eladmin/core/mongo"
	"go-eladmin/model/a_mongo_index"
	"go-eladmin/model/model_app"
	"go-eladmin/model/model_app_data_model"
	"go-eladmin/model/shareDB"
	"regexp"
	"time"

	"gopkg.in/mgo.v2"

	"gopkg.in/mgo.v2/bson"
)

type typeStatus int

const (
	modelAttributeTypeUndefine typeStatus = iota //待定
	//基础类型
	modelAttributeTypeString //String类型
	modelAttributeTypeInt    //Int类型
	modelAttributeTypeBool   //布尔类型
	modelAttributeTypeArray  //数组 （基础类型或者模型）

	modelAttributeTypeObject //模型

)

var (
	ModelTypeStatus = [...]string{"Undefined", "String", "Int", "Bool", "Object", "Array"}
)

const (
	DataModelPermissionALL    = "Data_Model_ALL"
	DataModelPermissionCreate = "Data_Model_CREATE"
	DataModelPermissionSelect = "Data_Model_SELECT"
	DataModelPermissionEdit   = "Data_Model_EDIT"
	DataModelPermissionDelete = "Data_Model_DELETE"
)

var (
	nameReg = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
)

const (
	dataModelCollection = mongo_index.CollectionDataModel
)

type Attribute struct {
	Type       typeStatus `json:"type,omitempty" bson:"type,omitempty"` //int string list bool
	TypeStatus string     `json:"type_status,omitempty" bson:"type_status,omitempty"`

	Name string `json:"name,omitempty" bson:"name,omitempty"`
	//attribute是数组时，数组内元素对象
	ModelName string `json:"model_name,omitempty" bson:"model_name,omitempty"`
	ModelId   int64  `json:"model_id,omitempty" bson:"model_id,omitempty"`
}

type DataModel struct {
	Id         int64                   `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string                  `json:"name,omitempty" bson:"name,omitempty"` //过滤中文名
	Alias      string                  `json:"alias,omitempty"  bson:"alias,omitempty"`
	CreateTime int64                   `json:"create_time,omitempty" bson:"create_time,omitempty"`
	Attributes []Attribute             `json:"attributes,omitempty" bson:"attributes,omitempty"` //模型的属性表
	Apps       []model_app.Application `json:"apps,omitempty" bson:"apps,omitempty"`             //不存入数据库
}

func (d DataModel) ToJson() string {
	js, _ := json.Marshal(d)
	return string(js)
}

func (d DataModel) isExist(query interface{}) bool {
	return mongo.IsExist(shareDB.DocManagerDBName(), dataModelCollection, query)
}

func (d DataModel) insert(docs ...interface{}) error {
	return mongo.Insert(shareDB.DocManagerDBName(), dataModelCollection, docs...)
}

func (d DataModel) update(selector, update interface{}) error {
	return mongo.Update(shareDB.DocManagerDBName(), dataModelCollection, selector, update, true)
}

func (d DataModel) findOne(query, selector interface{}) (DataModel, error) {
	ap := DataModel{}
	err := mongo.FindOne(shareDB.DocManagerDBName(), dataModelCollection, query, selector, &ap)
	return ap, err
}
func (d DataModel) findAll(query, selector interface{}) (results []DataModel, err error) {
	results = []DataModel{}
	err = mongo.FindAll(shareDB.DocManagerDBName(), dataModelCollection, query, selector, &results)
	return results, err
}

func (d DataModel) remove(selector interface{}) error {
	return mongo.Remove(shareDB.DocManagerDBName(), dataModelCollection, selector)
}

func (d DataModel) removeAll(selector interface{}) error {
	return mongo.RemoveAll(shareDB.DocManagerDBName(), dataModelCollection, selector)
}

func (d DataModel) totalCount(query, selector interface{}) (int, error) {
	return mongo.TotalCount(shareDB.DocManagerDBName(), dataModelCollection, query, selector)
}

func (d DataModel) findPage(page, limit int, query, selector interface{}, fields ...string) (results []DataModel, err error) {
	results = []DataModel{}
	err = mongo.FindPage(shareDB.DocManagerDBName(), dataModelCollection, page, limit, query, selector, &results, fields...)
	return
}

func (d DataModel) pipeAll(pipeline, result interface{}, allowDiskUse bool) error {
	return mongo.PipeAll(shareDB.DocManagerDBName(), dataModelCollection, pipeline, result, allowDiskUse)
}

func (d DataModel) pipeOne(pipeline, result interface{}, allowDiskUse bool) error {
	return mongo.PipeOne(shareDB.DocManagerDBName(), dataModelCollection, pipeline, result, allowDiskUse)
}

func (d DataModel) explain(pipeline, result interface{}) (results []DataModel, err error) {
	err = mongo.Explain(shareDB.DocManagerDBName(), dataModelCollection, pipeline, result)
	return
}

func (d DataModel) AddAttribute(a Attribute) error {
	if !checkNameReg(a.Name) {
		return fmt.Errorf("attribute name:%s not right", a.Name)
	}
	if d.isExistAttribute(a) {
		return fmt.Errorf("duplicate attribute name:%s", a.Name)
	}

	update := bson.M{"$addToSet": bson.M{"attributes": a}}
	change := mgo.Change{
		Update: update,
	}
	query := bson.M{"_id": d.Id}
	ms, c := mongo.Collection(shareDB.DocManagerDBName(), dataModelCollection)
	defer ms.Close()
	_, err := c.Find(query).Apply(change, nil)
	return err
}

func (d DataModel) AddAttributes(list []Attribute) error {
	for _, item := range list {
		if int(item.Type) >= len(ModelTypeStatus) || int(item.Type) < 0 {
			return fmt.Errorf("type Status:%d not found", item.Type)
		}
		item.TypeStatus = ModelTypeStatus[item.Type]
		fmt.Println("=====:", item.TypeStatus)
		if item.Type == modelAttributeTypeObject {
			m, err := d.FindOne(bson.M{"_id": item.ModelId}, nil)
			if err != nil {
				return fmt.Errorf("model attribute name:%s Type:%d id:%d not found",
					item.Name, item.Type, d.Id)
			}
			item.ModelName = m.Name
		}
		err := d.AddAttribute(item)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkNameReg(name string) bool {
	match := nameReg.FindAllString(name, -1)
	if len(match) == 0 {
		return false
	}
	return true
}

func (d DataModel) RemoveAttribute(a Attribute) error {
	if len(a.Name) == 0 {
		return fmt.Errorf("attribute name can not be nil")
	}
	selector := bson.M{"_id": d.Id}
	option := bson.M{"$pull": bson.M{"attributes": bson.M{"name": a.Name}}}
	ms, c := mongo.Collection(shareDB.DocManagerDBName(), dataModelCollection)
	defer ms.Close()
	return c.Update(selector, option)
}

func (d DataModel) isExistAttribute(a Attribute) bool {
	selector := bson.M{"_id": d.Id, "attributes.name": a.Name}
	return d.isExist(selector)
}

func (d DataModel) FindOne(query, selector interface{}) (dm DataModel, err error) {
	if query == nil {
		query = bson.M{"_id": d.Id}
	}
	dm, err = d.findOne(query, selector)
	if err != nil {
		return
	}
	dm.Apps, _ = d.fetchApplications(nil)
	return
}

func (d DataModel) fetchApplications(selector interface{}) (results []model_app.Application, err error) {
	if selector == nil {
		selector = bson.M{"_id": 1, "name": 1}
	}
	aM := model_app_data_model.AppDataModel{}
	aM.ModelId = d.Id

	list, err := aM.FindAll(bson.M{"model_id": d.Id}, nil)
	if err != nil {
		return
	}
	results = make([]model_app.Application, 0)
	for _, item := range list {
		app := model_app.Application{}
		app.Id = item.AppId
		app, err = app.FindOne(nil, nil)
		if err == nil {
			results = append(results, app)
		} else {
			item.Remove()
		}
	}
	return
}

//插入基本信息
func (d DataModel) Insert() (id int64, err error) {
	id, err = mongo.GetIncrementId(shareDB.DocManagerDBName(), dataModelCollection)
	if err != nil {
		return -1, err
	}
	if !checkNameReg(d.Name) {
		return -1, fmt.Errorf("data_model name:%s not right", d.Name)
	}
	d.CreateTime = time.Now().Unix()
	d.Id = id
	d.Apps = nil
	d.Attributes = nil
	err = d.insert(d)
	if err != nil {
		return -1, err
	}
	return id, err
}

func (d DataModel) updateApplication(list []model_app.Application) error {
	if !d.isExist(bson.M{"_id": d.Id}) {
		return fmt.Errorf("datamodel not exist,Id:%d", d.Id)
	}
	aM := model_app_data_model.AppDataModel{}
	aM.RemoveModelId(d.Id)
	for _, app := range list {
		if app.Exist() {
			aM.AppId = app.Id
			aM.ModelId = d.Id
			err := aM.Insert()
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("application not exist, Id:%d", app.Id)
		}
	}
	return nil
}

func (d DataModel) UpdateAppRelation() error {
	return d.updateApplication(d.Apps)
}

func (d DataModel) Update() error {
	if len(d.Name) > 0 && !checkNameReg(d.Name) {
		return fmt.Errorf("data_model name:%s not right", d.Name)
	}
	//d.updateApplication(d.Apps)
	d.Apps = nil
	//err := d.AddAttributes(d.Attributes)
	//if err != nil {
	//	return err
	//}
	d.Attributes = nil
	return d.update(bson.M{"_id": d.Id}, d)
}

func (d DataModel) Remove() error {
	aM := model_app_data_model.AppDataModel{}
	if aM.Exist(bson.M{"model_id": d.Id}) {
		return fmt.Errorf("model in use")
	}
	return d.remove(bson.M{"_id": d.Id})
}

func (d DataModel) TotalCount(query, selector interface{}) (int, error) {
	return d.totalCount(query, selector)
}
func (d DataModel) FindPageFilter(page, limit int, query, selector interface{}, fields ...string) ([]DataModel, error) {
	return d.findPage(page, limit, query, selector, fields...)
}
