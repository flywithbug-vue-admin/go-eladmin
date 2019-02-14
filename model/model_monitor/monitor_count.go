package model_monitor

import (
	"encoding/json"
	"fmt"
	"go-eladmin/core/mongo"
	"go-eladmin/model/a_mongo_index"
	"go-eladmin/model/shareDB"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	monitorCountCollection = mongo_index.CollectionMonitorCount
)

type MonitorCount struct {
	Count    int    `json:"count,omitempty" bson:"count,omitempty"`
	Monitor  string `json:"monitor,omitempty" bson:"monitor,omitempty"`
	TimeDate string `json:"time_date,omitempty" bson:"time_date,omitempty"` //2006-01-02:15 小时计算计算
	Total    int    `json:"total,omitempty" bson:"total,omitempty"`
}

func (v MonitorCount) ToJson() string {
	js, _ := json.Marshal(v)
	return string(js)
}

func (v MonitorCount) isExist(query interface{}) bool {
	return mongo.IsExist(shareDB.MonitorDBName(), monitorCountCollection, query)
}

func (v MonitorCount) insert(docs ...interface{}) error {
	return mongo.Insert(shareDB.MonitorDBName(), monitorCountCollection, docs...)
}

func (v MonitorCount) update(selector, update interface{}) error {
	return mongo.Update(shareDB.MonitorDBName(), monitorCountCollection, selector, update, true)
}

func (v MonitorCount) findOne(query, selector interface{}) (MonitorCount, error) {
	ap := MonitorCount{}
	err := mongo.FindOne(shareDB.MonitorDBName(), monitorCountCollection, query, selector, &ap)
	return ap, err
}
func (v MonitorCount) findAll(query, selector interface{}) (results []MonitorCount, err error) {
	results = []MonitorCount{}
	err = mongo.FindAll(shareDB.MonitorDBName(), monitorCountCollection, query, selector, &results)
	return results, err
}

func (v MonitorCount) remove(selector interface{}) error {
	return mongo.Remove(shareDB.MonitorDBName(), monitorCountCollection, selector)
}

func (v MonitorCount) removeAll(selector interface{}) error {
	return mongo.RemoveAll(shareDB.MonitorDBName(), monitorCountCollection, selector)
}

func (v MonitorCount) totalCount(query, selector interface{}) (int, error) {
	return mongo.TotalCount(shareDB.MonitorDBName(), monitorCountCollection, query, selector)
}

func (v MonitorCount) findPage(page, limit int, query, selector interface{}, fields ...string) (results []MonitorCount, err error) {
	results = []MonitorCount{}
	err = mongo.FindPage(shareDB.MonitorDBName(), monitorCountCollection, page, limit, query, selector, &results, fields...)
	return
}

func (v MonitorCount) pipeAll(pipeline, result interface{}, allowDiskUse bool) error {
	return mongo.PipeAll(shareDB.MonitorDBName(), monitorCountCollection, pipeline, result, allowDiskUse)
}

func (v MonitorCount) pipeOne(pipeline, result interface{}, allowDiskUse bool) error {
	return mongo.PipeOne(shareDB.MonitorDBName(), monitorCountCollection, pipeline, result, allowDiskUse)
}

func (v MonitorCount) explain(pipeline, result interface{}) (results []MonitorCount, err error) {
	err = mongo.Explain(shareDB.MonitorDBName(), monitorCountCollection, pipeline, result)
	return
}

func (v MonitorCount) Insert() error {
	timeF := time.Now().Format(TimeLayout)
	v.TimeDate = timeF[:10]
	return v.insert(v)
}

func (v MonitorCount) FindOne(query interface{}) (MonitorCount, error) {
	return v.findOne(query, nil)
}

func (v MonitorCount) IncrementMonitorCount() (int, error) {
	if len(v.Monitor) == 0 {
		return -1, fmt.Errorf("monitor is null")
	}
	timeF := time.Now().Format(TimeLayout)
	v.TimeDate = timeF[:10]
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"count": 1}},
		ReturnNew: true,
	}
	_, c := mongo.Collection(shareDB.MonitorDBName(), monitorCountCollection)
	_, err := c.Find(bson.M{"monitor": v.Monitor, "time_date": v.TimeDate}).Apply(change, v)
	if err != nil {
		v.Count = 1
		err = v.Insert()
		if err != nil {
			return -1, err
		}
	}
	return v.Count, nil
}

func (v MonitorCount) TotalCount(query, selector interface{}) (int, error) {
	return v.totalCount(query, selector)
}

func (v MonitorCount) FindPageFilter(page, limit int, query, selector interface{}, fields ...string) ([]MonitorCount, error) {
	return v.findPage(page, limit, query, selector, fields...)
}

func (v MonitorCount) TotalSumCount(query interface{}) (int, error) {
	match := bson.M{"$match": query}
	group := bson.M{"$group": bson.M{"_id": "time_date", "total": bson.M{"$sum": "$count"}}}
	pipeline := []bson.M{
		match,
		group,
	}
	v.pipeOne(pipeline, &v, true)
	return v.Total, nil
}
