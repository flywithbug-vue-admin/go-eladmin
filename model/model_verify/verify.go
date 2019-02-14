package model_verify

import (
	"encoding/json"
	"fmt"
	"go-eladmin/core/mongo"
	"go-eladmin/model/a_mongo_index"
	"go-eladmin/model/shareDB"
	"time"

	"math/rand"

	"gopkg.in/mgo.v2/bson"
)

const (
	verifyCollection = mongo_index.CollectionVerify
)

type VerificationCode struct {
	Id         int64  `json:"id,omitempty" bson:"_id,omitempty"`
	Value      string `json:"value,omitempty" bson:"value,omitempty"` //请求code的value
	Code       string `json:"code,omitempty" bson:"code,omitempty"`   //返回Code
	Vld        int64  `json:"vld,omitempty" bson:"vld,omitempty"`     //有效期
	CreateTime int64  `json:"create_time,omitempty" bson:"create_time,omitempty"`
	Scenes     string `json:"scenes,omitempty" bson:"scenes,omitempty"` //场景
	Status     bool   `json:"status"`
}

func (v VerificationCode) ToJson() string {
	js, _ := json.Marshal(v)
	return string(js)
}

func (v VerificationCode) insert(docs ...interface{}) error {
	return mongo.Insert(shareDB.DocManagerDBName(), verifyCollection, docs...)
}

func (v VerificationCode) isExist(query interface{}) bool {
	return mongo.IsExist(shareDB.DocManagerDBName(), verifyCollection, query)
}

func (v VerificationCode) findOne(query, selector interface{}) (VerificationCode, error) {
	ap := VerificationCode{}
	err := mongo.FindOne(shareDB.DocManagerDBName(), verifyCollection, query, selector, &ap)
	return ap, err
}

func (v VerificationCode) findAll(query, selector interface{}) (results []VerificationCode, err error) {
	results = []VerificationCode{}
	err = mongo.FindAll(shareDB.DocManagerDBName(), verifyCollection, query, selector, &results)
	return results, err
}

func (v VerificationCode) update(selector, update interface{}) error {
	return mongo.Update(shareDB.DocManagerDBName(), verifyCollection, selector, update, true)
}

func (v VerificationCode) remove(selector interface{}) error {
	return mongo.Remove(shareDB.DocManagerDBName(), verifyCollection, selector)
}

func (v VerificationCode) removeAll(selector interface{}) error {
	return mongo.RemoveAll(shareDB.DocManagerDBName(), verifyCollection, selector)
}

func (v VerificationCode) totalCount(query, selector interface{}) (int, error) {
	return mongo.TotalCount(shareDB.DocManagerDBName(), verifyCollection, query, selector)
}

func (v VerificationCode) findPage(page, limit int, query, selector interface{}, fields ...string) (results []VerificationCode, err error) {
	results = []VerificationCode{}
	err = mongo.FindPage(shareDB.DocManagerDBName(), verifyCollection, page, limit, query, selector, &results, fields...)
	return
}

func (v VerificationCode) TotalCount(query, selector interface{}) (int, error) {
	return v.totalCount(query, selector)
}
func (v VerificationCode) FindPageFilter(page, limit int, query, selector interface{}, fields ...string) ([]VerificationCode, error) {
	return v.findPage(page, limit, query, selector, fields...)
}

func (v VerificationCode) Insert() error {
	v.CreateTime = time.Now().Unix()
	v.Status = true
	return v.insert(v)
}

func GeneralVerifyData(value string) (string, error) {
	var verify VerificationCode
	verify.Code = verify.generalVCode(value)
	verify.Value = value
	verify.Vld = time.Now().Unix() + 300
	err := verify.Insert()
	return verify.Code, err
}

func (v VerificationCode) generalVCode(value string) string {
	rand.Int()
	vCode := createCaptcha()
	if v.isExist(bson.M{"status": true, "code": vCode, "value": value, "vld": bson.M{"$gte": time.Now().Unix()}}) {
		vCode = v.generalVCode(value)
	}
	return vCode
}

func CheckVerify(value, code string) bool {
	var verify VerificationCode
	if verify.isExist(bson.M{"status": true, "code": code, "value": value, "vld": bson.M{"$gte": time.Now().Unix()}}) {
		updateMarked(value, code)
		return true
	}
	return false
}

func updateMarked(value, code string) {
	var verify VerificationCode
	verify.update(bson.M{"code": code, "value": value, "status": true}, bson.M{"status": false})
}

func createCaptcha() string {
	return fmt.Sprintf("%04v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
}
