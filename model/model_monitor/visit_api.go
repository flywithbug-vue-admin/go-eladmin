package model_monitor

import (
	"encoding/json"
	"fmt"
	"go-eladmin/core/mongo"
	"go-eladmin/model/a_mongo_index"
	"go-eladmin/model/shareDB"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type VisitApi struct {
	Path     string      `json:"path,omitempty" bson:"path,omitempty"`
	TimeDate string      `json:"time_date,omitempty" bson:"time_date,omitempty"` //2006-01-02:15 小时计算计算
	Count    int         `json:"count,omitempty" bson:"count,omitempty"`
	Method   string      `json:"method,omitempty" bson:"method,omitempty"`
	Total    int         `json:"total,omitempty" bson:"total,omitempty"`
	Para     interface{} `json:"para,omitempty" bson:"para,omitempty"`
}

const (
	visitApiApiCollection = mongo_index.CollectionVisitApi
)

func (v VisitApi) ToJson() string {
	js, _ := json.Marshal(v)
	return string(js)
}

func (v VisitApi) isExist(query interface{}) bool {
	return mongo.IsExist(shareDB.MonitorDBName(), visitApiApiCollection, query)
}

func (v VisitApi) insert(docs ...interface{}) error {
	return mongo.Insert(shareDB.MonitorDBName(), visitApiApiCollection, docs...)
}

func (v VisitApi) update(selector, update interface{}) error {
	return mongo.Update(shareDB.MonitorDBName(), visitApiApiCollection, selector, update, true)
}

func (v VisitApi) findOne(query, selector interface{}) (Log, error) {
	ap := Log{}
	err := mongo.FindOne(shareDB.MonitorDBName(), visitApiApiCollection, query, selector, &ap)
	return ap, err
}
func (v VisitApi) findAll(query, selector interface{}) (results []Log, err error) {
	results = []Log{}
	err = mongo.FindAll(shareDB.MonitorDBName(), visitApiApiCollection, query, selector, &results)
	return results, err
}

func (v VisitApi) remove(selector interface{}) error {
	return mongo.Remove(shareDB.MonitorDBName(), visitApiApiCollection, selector)
}

func (v VisitApi) removeAll(selector interface{}) error {
	return mongo.RemoveAll(shareDB.MonitorDBName(), visitApiApiCollection, selector)
}

func (v VisitApi) totalCount(query, selector interface{}) (int, error) {
	return mongo.TotalCount(shareDB.MonitorDBName(), visitApiApiCollection, query, selector)
}

func (v VisitApi) findPage(page, limit int, query, selector interface{}, fields ...string) (results []VisitApi, err error) {
	results = []VisitApi{}
	err = mongo.FindPage(shareDB.MonitorDBName(), visitApiApiCollection, page, limit, query, selector, &results, fields...)
	return
}
func (v VisitApi) FindPipeline(pipeline []bson.M) (results []VisitApi, err error) {
	results = make([]VisitApi, 0)
	err = v.pipeAll(pipeline, &results, true)
	return
}

func (v VisitApi) pipeAll(pipeline, result interface{}, allowDiskUse bool) error {
	return mongo.PipeAll(shareDB.MonitorDBName(), visitApiApiCollection, pipeline, result, allowDiskUse)
}

func (v VisitApi) pipeOne(pipeline, result interface{}, allowDiskUse bool) error {
	return mongo.PipeOne(shareDB.MonitorDBName(), visitApiApiCollection, pipeline, result, allowDiskUse)
}

func (v VisitApi) explain(pipeline, result interface{}) (results []VisitApi, err error) {
	err = mongo.Explain(shareDB.MonitorDBName(), visitApiApiCollection, pipeline, result)
	return
}

func (v VisitApi) Insert() error {
	if len(v.Path) == 0 || len(v.Method) == 0 {
		return fmt.Errorf("path or method is null")
	}
	return v.insert(v)
}

func (v VisitApi) IncrementVisitApi() (int, error) {
	if len(v.Path) == 0 || len(v.Method) == 0 {
		return -1, fmt.Errorf("path or method is null")
	}
	update := bson.M{"$inc": bson.M{"count": 1}}
	if v.Para != nil {
		update["para"] = v.Para
	}
	change := mgo.Change{
		Update:    update,
		ReturnNew: true,
	}
	_, c := mongo.Collection(shareDB.MonitorDBName(), visitApiApiCollection)
	_, err := c.Find(bson.M{"path": v.Path, "method": v.Method, "time_date": v.TimeDate}).Apply(change, v)
	if err != nil {
		v.Count = 1
		err = v.Insert()
		if err != nil {
			return -1, err
		}
	}
	return v.Count, nil
}

func (v VisitApi) TotalCount(query, selector interface{}) (int, error) {
	return v.totalCount(query, selector)
}

func (v VisitApi) FindPageFilter(page, limit int, query, selector interface{}, fields ...string) ([]VisitApi, error) {
	return v.findPage(page, limit, query, selector, fields...)
}

func (v VisitApi) TotalSumCount(query interface{}) (int, error) {
	match := bson.M{"$match": query}
	group := bson.M{"$group": bson.M{"_id": "time_date", "total": bson.M{"$sum": "$count"}}}
	pipeline := []bson.M{
		match,
		group,
	}
	v.pipeOne(pipeline, &v, true)
	return v.Total, nil
}
