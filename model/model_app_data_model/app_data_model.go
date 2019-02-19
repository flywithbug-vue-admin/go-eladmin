package model_app_data_model

import (
	"encoding/json"
	"fmt"
	"go-eladmin/common"
	"go-eladmin/core/mongo"
	"go-eladmin/model/a_mongo_index"
	"go-eladmin/model/shareDB"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//应用和数据模型关联

const (
	AppDataModelCollection = mongo_index.CollectionAppDataModel
	MaxVersion             = 4294967294
)

type AppDataModel struct {
	Id           int64  `json:"id,omitempty" bson:"_id,omitempty"`
	ModelId      int64  `json:"model_id,omitempty" bson:"model_id,omitempty"`
	AppId        int64  `json:"app_id,omitempty" bson:"app_id,omitempty"`
	StartVersion string `json:"start_version" bson:"start_version,omitempty"`
	EndVersion   string `json:"end_version" bson:"end_version,omitempty"`
	StartVNum    int    `json:"start_v_num,omitempty" bson:"start_v_num,omitempty"`
	EndVNum      int    `json:"end_v_num,omitempty" bson:"end_v_num,omitempty"`
	CreateTime   int64  `json:"create_time,omitempty" bson:"create_time,omitempty"`
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

func (a AppDataModel) TotalCount(query, selector interface{}) (int, error) {
	return a.totalCount(query, selector)
}

func (a AppDataModel) FindPageFilter(page, limit int, query, selector interface{}, fields ...string) (apps []AppDataModel, err error) {
	apps, err = a.findPage(page, limit, query, selector, fields...)
	return
}

func (a AppDataModel) FindOne(query, selector interface{}) (adm AppDataModel, err error) {
	adm, err = a.findOne(query, selector)
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
	if len(a.StartVersion) == 0 {
		a.StartVNum = 0
	} else {
		a.StartVNum = common.TransformVersionToInt(a.StartVersion)
		if a.StartVNum == -1 {
			return fmt.Errorf("startVersion:%s not right", a.StartVersion)
		}
	}
	if len(a.EndVersion) == 0 {
		a.EndVNum = MaxVersion
	} else {
		a.EndVNum = common.TransformVersionToInt(a.EndVersion)
		if a.EndVNum == -1 {
			return fmt.Errorf("endVersion:%s not right", a.EndVersion)
		}
		if a.EndVNum < a.StartVNum {
			return fmt.Errorf("endVersion is bigger than startVersion")
		}
	}
	return a.insert(a)
}

func (a AppDataModel) Update() error {
	vNum := common.TransformVersionToInt(a.StartVersion)
	if vNum > 0 {
		a.StartVNum = vNum
	}
	vNum = common.TransformVersionToInt(a.EndVersion)
	if vNum > 0 {
		a.EndVNum = vNum
	}
	if a.EndVNum > 0 && a.EndVNum < a.StartVNum {
		return fmt.Errorf("startVersion:%s, bigger than endVersion:%s", a.StartVersion, a.EndVersion)
	}
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
