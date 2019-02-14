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

const (
	visitUIdCollection = mongo_index.CollectionVisitUId
)

type VisitUId struct {
	ClientIp string `json:"client_ip,omitempty" bson:"client_ip,omitempty"`
	UUID     string `json:"uuid,omitempty" bson:"uuid,omitempty"`
	UserId   int64  `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Count    int    `json:"count,omitempty" bson:"count,omitempty"`         //访问次数
	TimeDate string `json:"time_date,omitempty" bson:"time_date,omitempty"` //2018-06-10
	Total    int    `json:"total,omitempty" bson:"total,omitempty"`
}

func (v *VisitUId) ReSet() {
	v.ClientIp = ""
	v.UUID = ""
	v.UserId = 0
	v.Count = 0
	v.TimeDate = ""
	v.Total = 0
}

func (v VisitUId) ToJson() string {
	js, _ := json.Marshal(v)
	return string(js)
}

func (v VisitUId) isExist(query interface{}) bool {
	return mongo.IsExist(shareDB.MonitorDBName(), visitUIdCollection, query)
}

func (v VisitUId) insert(docs ...interface{}) error {
	return mongo.Insert(shareDB.MonitorDBName(), visitUIdCollection, docs...)
}

func (v VisitUId) update(selector, update interface{}) error {
	return mongo.Update(shareDB.MonitorDBName(), visitUIdCollection, selector, update, true)
}

func (v VisitUId) findOne(query, selector interface{}) (VisitUId, error) {
	ap := VisitUId{}
	err := mongo.FindOne(shareDB.MonitorDBName(), visitUIdCollection, query, selector, &ap)
	return ap, err
}
func (v VisitUId) findAll(query, selector interface{}) (results []VisitUId, err error) {
	results = []VisitUId{}
	err = mongo.FindAll(shareDB.MonitorDBName(), visitUIdCollection, query, selector, &results)
	return results, err
}

func (v VisitUId) remove(selector interface{}) error {
	return mongo.Remove(shareDB.MonitorDBName(), visitUIdCollection, selector)
}

func (v VisitUId) removeAll(selector interface{}) error {
	return mongo.RemoveAll(shareDB.MonitorDBName(), visitUIdCollection, selector)
}

func (v VisitUId) totalCount(query, selector interface{}) (int, error) {
	return mongo.TotalCount(shareDB.MonitorDBName(), visitUIdCollection, query, selector)
}

func (v VisitUId) findPage(page, limit int, query, selector interface{}, fields ...string) (results []VisitUId, err error) {
	results = []VisitUId{}
	err = mongo.FindPage(shareDB.MonitorDBName(), visitUIdCollection, page, limit, query, selector, &results, fields...)
	return
}

func (v VisitUId) pipeAll(pipeline, result interface{}, allowDiskUse bool) error {
	return mongo.PipeAll(shareDB.MonitorDBName(), visitUIdCollection, pipeline, result, allowDiskUse)
}

func (v VisitUId) pipeOne(pipeline, result interface{}, allowDiskUse bool) error {
	return mongo.PipeOne(shareDB.MonitorDBName(), visitUIdCollection, pipeline, result, allowDiskUse)
}

func (v VisitUId) explain(pipeline, result interface{}) (results []VisitUId, err error) {
	err = mongo.Explain(shareDB.MonitorDBName(), visitUIdCollection, pipeline, result)
	return
}

func (v VisitUId) Insert() error {
	if len(v.ClientIp) == 0 || len(v.UUID) == 0 {
		return fmt.Errorf("client_ip or uuid is null")
	}
	return v.insert(v)
}

func (v VisitUId) IncrementVisitUId() (int, error) {
	if len(v.ClientIp) == 0 || len(v.UUID) == 0 {
		return -1, fmt.Errorf("client_ip or uuid is null")
	}
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"count": 1}, "$set": bson.M{"user_id": v.UserId}},
		ReturnNew: true,
	}
	_, c := mongo.Collection(shareDB.MonitorDBName(), visitUIdCollection)
	_, err := c.Find(bson.M{"client_ip": v.ClientIp, "uuid": v.UUID, "time_date": v.TimeDate}).Apply(change, v)
	if err != nil {
		v.Count = 1
		err = v.Insert()
		if err != nil {
			return -1, err
		}
	}
	return v.Count, nil
}

func (v VisitUId) TotalCount(query, selector interface{}) (int, error) {
	return v.totalCount(query, selector)
}

func (v VisitUId) FindPageFilter(page, limit int, query, selector interface{}, fields ...string) ([]VisitUId, error) {
	return v.findPage(page, limit, query, selector, fields...)
}

func (v VisitUId) TotalSumCount(query interface{}) (int, error) {
	match := bson.M{"$match": query}
	group := bson.M{"$group": bson.M{"_id": "time_date", "total": bson.M{"$sum": "$count"}}}
	pipeline := []bson.M{
		match,
		group,
	}
	v.pipeOne(pipeline, &v, true)
	return v.Total, nil
}
