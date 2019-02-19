package model_data_model

import (
	"encoding/json"
	"fmt"
	"go-eladmin/core/mongo"
	"go-eladmin/model/a_mongo_index"
	"go-eladmin/model/model_app"
	"go-eladmin/model/model_app_data_model"
	"go-eladmin/model/model_user"
	"go-eladmin/model/shareDB"
	"regexp"
	"time"

	"gopkg.in/mgo.v2"

	"gopkg.in/mgo.v2/bson"
)

type typeStatus int

const (
	modelAttributeTypeString = "String" //String类型
	modelAttributeTypeInt    = "Int"    //Int类型
	modelAttributeTypeFloat  = "Float"  //浮点数
	modelAttributeTypeBool   = "Bool"   //布尔类型
	modelAttributeTypeArray  = "Array"  //数组 （只能是其他的类型模型 必须要有ModelId）
	modelAttributeTypeObject = "Object" //模型
)

const (
	DataModelPermissionALL    = "Model_ALL"
	DataModelPermissionCreate = "Model_CREATE"
	DataModelPermissionSelect = "Model_SELECT"
	DataModelPermissionEdit   = "Model_EDIT"
	DataModelPermissionDelete = "Model_DELETE"
)

var nameReg = regexp.MustCompile(`^[A-Z][A-Za-z0-9_]+$`)

const (
	dataModelCollection = mongo_index.CollectionDataModel
)

type Attribute struct {
	Type string `json:"type,omitempty" bson:"type,omitempty"` //数据类型
	Name string `json:"name,omitempty" bson:"name,omitempty"`
	//attribute是数组时，数组内元素对象
	ModelName string `json:"model_name,omitempty" bson:"model_name,omitempty"`
	ModelId   int64  `json:"model_id,omitempty" bson:"model_id,omitempty"`
	Default   string `json:"default,omitempty" bson:"default,omitempty"`   //默认值
	Required  bool   `json:"required" bson:"required,omitempty"`           //是否必填 RequestPara 使用
	Comments  string `json:"comments,omitempty" bson:"comments,omitempty"` //属性说明
}

func (a Attribute) ToJson() string {
	js, _ := json.Marshal(a)
	return string(js)
}

/*
TODO 历史版本记录
*/
type DataModel struct {
	Id         int64                   `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string                  `json:"name,omitempty" bson:"name,omitempty"` //过滤中文名
	Alias      string                  `json:"alias,omitempty"  bson:"alias,omitempty"`
	Desc       string                  `json:"desc,omitempty" bson:"desc,omitempty"` //模型说明
	CreateTime int64                   `json:"create_time,omitempty" bson:"create_time,omitempty"`
	Attributes []Attribute             `json:"attributes,omitempty" bson:"attributes,omitempty"` //模型的属性表
	Apps       []model_app.Application `json:"apps,omitempty" bson:"apps,omitempty"`             //不存入数据库
	Owner      model_user.User         `json:"owner,omitempty" bson:"owner,omitempty"`           //模型负责人（初始为创建人）
	ParentId   int64                   `json:"parent_id,omitempty" bson:"parent_id,omitempty"`
	Parent     interface{}             `json:"parent,omitempty" bson:"parent,omitempty"`
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

func (d DataModel) Exist(query interface{}) bool {
	return d.isExist(query)
}

func (d DataModel) AddAttribute(a Attribute) error {
	if d.isExistAttribute(a) {
		return fmt.Errorf("duplicate attribute name:%s", a.Name)
	}
	if err := checkNameReg(a.Name); err != nil {
		return err
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
	//删除之前同名属性
	d.RemoveAttributes(list)
	for _, item := range list {
		switch item.Type {
		case modelAttributeTypeString,
			modelAttributeTypeInt,
			modelAttributeTypeFloat,
			modelAttributeTypeBool:
			//基础数据类型不处理
		case modelAttributeTypeObject:
			m, err := d.FindSimpleOne(bson.M{"_id": item.ModelId}, nil)
			if err != nil {
				return fmt.Errorf("model attribute name:%s Type:%s id:%d not found",
					item.Name, item.Type, d.Id)
			}
			item.ModelName = m.Name
		case modelAttributeTypeArray:
			if item.ModelId > 0 {
				m, err := d.FindSimpleOne(bson.M{"_id": item.ModelId}, nil)
				if err != nil {
					return fmt.Errorf("model attribute name:%s Type:%d id:%d not found",
						item.Name, item.Type, d.Id)
				}
				item.ModelName = m.Name
			} else {
				switch item.ModelName {
				case modelAttributeTypeString,
					modelAttributeTypeInt,
					modelAttributeTypeFloat,
					modelAttributeTypeBool:
					//基础数据类型不处理直接使用
				default:
					return fmt.Errorf("属性类型未定义")
				}
				return fmt.Errorf("数组元素属性类型未指定")
			}

		default:
			return fmt.Errorf("属性类型未定义")
		}
		err := d.AddAttribute(item)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkNameReg(name string) error {
	match := nameReg.FindString(name)
	if len(match) == 0 {
		return fmt.Errorf("attribute :%s not right (note:^[A-Z][A-Za-z0-9_]+$)", name)
	}
	return nil
}

func (d DataModel) RemoveAttributes(list []Attribute) error {
	for _, item := range list {
		d.RemoveAttribute(item)
	}
	return nil
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

/**
查询模型详细信息
*/
func (d DataModel) FindOne(query, selector interface{}) (dm DataModel, err error) {
	if query == nil {
		query = bson.M{"_id": d.Id}
	}
	dm, err = d.findOne(query, selector)
	if err != nil {
		return dm, err
	}
	if dm.ParentId > 0 {
		query = bson.M{"_id": dm.ParentId}
		parent, err := d.findOne(query, selector)
		if err != nil {
			d.remove(query)
		} else {
			dm.Parent = parent
		}
	}
	dm.Apps, _ = d.fetchApplications(nil)
	result := make([]DataModel, 1)
	result[0] = dm
	err = fetchOwnerAndAttributes(result)
	return result[0], err
}

func (d DataModel) FindSimpleOne(query, selector interface{}) (dm DataModel, err error) {
	if query == nil {
		query = bson.M{"_id": d.Id}
	}
	dm, err = d.findOne(query, selector)
	if err != nil {
		return
	}
	return dm, err
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
	if err = checkNameReg(d.Name); err != nil {
		return -1, err
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
		if app.Exist(bson.M{"_id": app.Id}) {
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
	if err := checkNameReg(d.Name); err != nil {
		return err
	}
	d.Apps = nil
	d.Attributes = nil
	return d.update(bson.M{"_id": d.Id}, d)
}

func (d DataModel) Remove() error {
	aM := model_app_data_model.AppDataModel{}
	if aM.Exist(bson.M{"model_id": d.Id}) {
		return fmt.Errorf("model in use")
	}
	if d.isExist(bson.M{"parent_id": d.Id}) {
		return fmt.Errorf("has son model")
	}
	return d.remove(bson.M{"_id": d.Id})
}

func (d DataModel) TotalCount(query, selector interface{}) (int, error) {
	return d.totalCount(query, selector)
}

func (d DataModel) FindPageFilter(page, limit int, query, selector interface{}, fields ...string) ([]DataModel, error) {
	result, err := d.findPage(page, limit, query, selector, fields...)
	if err != nil {
		return nil, err
	}
	//err = fetchOwnerAndAttributes(result)
	return result, err
}

func fetchOwnerAndAttributes(result []DataModel) error {
	for index := 0; index < len(result); index++ {
		if result[index].Owner.Id > 0 {
			user := model_user.User{}
			user.Id = result[index].Owner.Id
			user, err := user.FindOne()
			if err != nil {
				return err
			}
			result[index].Owner = user
		}
		for index1 := range result[index].Attributes {
			if result[index].Attributes[index1].ModelId > 0 {
				dm := DataModel{}
				dm.Id = result[index].Attributes[index1].ModelId
				dm, err := dm.FindSimpleOne(bson.M{"_id": dm.Id}, nil)
				if err != nil {
					return err
				}
				result[index].Attributes[index1].ModelName = dm.Name
			}
		}
	}
	return nil
}
