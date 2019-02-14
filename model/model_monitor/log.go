package model_monitor

import (
	"encoding/json"
	"go-eladmin/core/mongo"
	"go-eladmin/model/a_mongo_index"
	"go-eladmin/model/shareDB"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	logCollection = mongo_index.CollectionLog
)

type Log struct {
	Time         string        `json:"time,omitempty" bson:"time,omitempty"`
	Code         string        `json:"code,omitempty" bson:"code,omitempty"`
	Info         string        `json:"info,omitempty" bson:"info,omitempty"`
	Level        int           `json:"level,omitempty" bson:"level,omitempty"`
	Flag         string        `json:"flag,omitempty" bson:"flag,omitempty"`
	ClientIp     string        `json:"client_ip,omitempty" bson:"client_ip,omitempty"`
	Method       string        `json:"method,omitempty" bson:"method,omitempty"`
	Path         string        `json:"path,omitempty" bson:"path,omitempty"`
	RequestId    string        `json:"request_id,omitempty" bson:"request_id,omitempty"`
	Latency      time.Duration `json:"latency,omitempty" bson:"latency,omitempty"`
	StatusCode   int           `json:"status_code,omitempty" bson:"status_code,omitempty"`
	UserId       int64         `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Para         interface{}   `json:"para,omitempty" bson:"para,omitempty"`
	ResponseCode int           `json:"response_code,omitempty" bson:"response_code,omitempty"`
	Response     interface{}   `json:"response,omitempty" bson:"response,omitempty"`

	StartTime int64  `json:"start_time,omitempty" bson:"start_time,omitempty"`
	EndTime   int64  `json:"end_time,omitempty" bson:"end_time,omitempty"`
	UUID      string `json:"uuid,omitempty" bson:"uuid,omitempty"`
}

func (l *Log) ReSet() {
	l.Time = ""
	l.Code = ""
	l.Info = ""
	l.Level = 0
	l.Flag = ""
	l.ClientIp = ""
	l.Method = ""
	l.Path = ""
	l.RequestId = ""
	l.Latency = 0
	l.StatusCode = 0
	l.UserId = 0
	l.Para = nil
	l.ResponseCode = 0
	l.StartTime = 0
	l.EndTime = 0
	l.UUID = ""
}

func (l Log) ToJson() string {
	js, _ := json.Marshal(l)
	return string(js)
}

func (l Log) isExist(query interface{}) bool {
	return mongo.IsExist(shareDB.MonitorDBName(), logCollection, query)
}

func (l Log) insert(docs ...interface{}) error {
	return mongo.Insert(shareDB.MonitorDBName(), logCollection, docs...)
}

func (l Log) update(selector, update interface{}) error {
	return mongo.Update(shareDB.MonitorDBName(), logCollection, selector, update, true)
}

func (l Log) findOne(query, selector interface{}) (Log, error) {
	ap := Log{}
	err := mongo.FindOne(shareDB.MonitorDBName(), logCollection, query, selector, &ap)
	return ap, err
}
func (l Log) findAll(query, selector interface{}) (results []Log, err error) {
	results = []Log{}
	err = mongo.FindAll(shareDB.MonitorDBName(), logCollection, query, selector, &results)
	return results, err
}

func (l Log) remove(selector interface{}) error {
	return mongo.Remove(shareDB.MonitorDBName(), logCollection, selector)
}

func (l Log) removeAll(selector interface{}) error {
	return mongo.RemoveAll(shareDB.MonitorDBName(), logCollection, selector)
}

func (l Log) totalCount(query, selector interface{}) (int, error) {
	return mongo.TotalCount(shareDB.MonitorDBName(), logCollection, query, selector)
}

func (l Log) findPage(page, limit int, query, selector interface{}, fields ...string) (results []Log, err error) {
	results = []Log{}
	err = mongo.FindPage(shareDB.MonitorDBName(), logCollection, page, limit, query, selector, &results, fields...)
	return
}

func (l Log) pipeAll(pipeline, result interface{}, allowDiskUse bool) error {
	return mongo.PipeAll(shareDB.MonitorDBName(), logCollection, pipeline, result, allowDiskUse)
}

func (l Log) pipeOne(pipeline, result interface{}, allowDiskUse bool) error {
	return mongo.PipeOne(shareDB.MonitorDBName(), logCollection, pipeline, result, allowDiskUse)
}

func (l Log) explain(pipeline, result interface{}) (results []Log, err error) {
	err = mongo.Explain(shareDB.MonitorDBName(), logCollection, pipeline, result)
	return
}

func (l Log) Update() error {
	//if !l.isExist(bson.M{"request_id": l.RequestId}) {
	//	return l.Insert()
	//}
	return l.update(bson.M{"request_id": l.RequestId}, l)
}

func (l Log) Insert() error {
	if l.isExist(bson.M{"request_id": l.RequestId}) {
		return l.Update()
	}
	return l.insert(l)
}

func (l Log) TotalCount(query, selector interface{}) (int, error) {
	return l.totalCount(query, selector)
}
func (l Log) FindPageFilter(page, limit int, query, selector interface{}, fields ...string) ([]Log, error) {
	return l.findPage(page, limit, query, selector, fields...)
}
