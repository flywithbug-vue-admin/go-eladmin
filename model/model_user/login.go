package model_user

import (
	"encoding/json"
	"errors"
	"go-eladmin/core/mongo"
	"go-eladmin/model/a_mongo_index"
	"go-eladmin/model/shareDB"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type state int

const (
	loginCollection = mongo_index.CollectionLogin
	//StatusNormal    = 0
	StatusLogin  = 1
	StatusLogout = 2
)

type Login struct {
	Id         int64  `json:"id,omitempty" bson:"_id,omitempty"`
	UserId     int64  `bson:"user_id"`     // 用户ID
	Token      string `bson:"token"`       // 用户TOKEN
	CreateTime int64  `bson:"create_time"` // 登录日期
	LoginIp    string `bson:"login_ip"`    // 登录IP
	Status     state  `bson:"status"`      //status 1 已登录，2表示退出登录
	Forbidden  bool   `bson:"forbidden"`   //false 表示未禁言
	UserAgent  string `bson:"user_agent"`  //用户UA
	UpdatedAt  int64  `json:"updated_at,omitempty" bson:"updated_at"`
}

func (l Login) ToJson() string {
	js, _ := json.Marshal(l)
	return string(js)
}

func UserLogin(userID int64, userAgent, token, ip string) (l *Login, err error) {
	l = new(Login)
	l.Id, _ = mongo.GetIncrementId(shareDB.DocManagerDBName(), loginCollection)
	l.UserId = userID
	l.UserAgent = userAgent
	l.Token = token
	l.CreateTime = time.Now().Unix()
	l.UpdatedAt = l.CreateTime
	l.Status = StatusLogin
	l.LoginIp = ip
	err = l.Insert()
	return
}

func (l Login) FindAll() ([]Login, error) {
	var results []Login
	err := mongo.FindAll(shareDB.DocManagerDBName(), userCollection, nil, nil, &results)
	return results, err
}

func (l *Login) Insert() error {
	if l.UserId == 0 {
		return errors.New("user_id can not be 0")
	}
	return mongo.Insert(shareDB.DocManagerDBName(), loginCollection, l)
}

//status 0 退出登录，1 登录
//	return mongo.Update(shareDB.DocManagerDBName(), todoCollection, bson.M{"_id": t.Id}, bson.M{"$set": bson.M{"title": t.Title, "completed": t.Completed, "updated_at": t.UpdatedAt}})
func UpdateLoginStatus(token string, status int) error {
	updateAt := time.Now().Unix()
	return mongo.Update(shareDB.DocManagerDBName(), loginCollection, bson.M{"token": token}, bson.M{"status": status, "updated_at": updateAt}, true)
}

func FindLoginByToken(token string) (l *Login, err error) {
	l = new(Login)
	err = mongo.FindOne(shareDB.DocManagerDBName(), loginCollection, bson.M{"token": token}, bson.M{"status": StatusLogout}, &l)
	return
}
