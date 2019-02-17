package model_app

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-eladmin/core/mongo"
	"go-eladmin/model/a_mongo_index"
	"go-eladmin/model/model_app_data_model"
	"go-eladmin/model/model_app_manager"
	"go-eladmin/model/model_user"
	"go-eladmin/model/shareDB"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	appCollection = mongo_index.CollectionApp

	//ApplicationPermissionAll    = "APP_ALL"
	ApplicationPermissionSelect = "APP_SELECT"
	ApplicationPermissionCreate = "APP_CREATE"
	ApplicationPermissionEdit   = "APP_EDIT"
	ApplicationPermissionDelete = "APP_DELETE"
)

//修改规则，等级
// role 等级为1 的用户可以编辑
type Application struct {
	Id         int64             `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string            `json:"name,omitempty" bson:"name,omitempty"`               //应用（组件）名称
	Desc       string            `json:"desc,omitempty" bson:"desc,omitempty"`               //项目描述
	CreateTime int64             `json:"create_time,omitempty" bson:"create_time,omitempty"` //创建时间
	Icon       string            `json:"icon,omitempty" bson:"icon,omitempty"`               //icon 地址
	Owner      *model_user.User  `json:"owner,omitempty" bson:"owner,omitempty"`             //应用所有者
	OwnerId    int64             `json:"owner_id,omitempty" bson:"owner_id,omitempty"`
	BundleId   string            `json:"bundle_id,omitempty" bson:"bundle_id,omitempty"`
	Managers   []model_user.User `json:"managers,omitempty" bson:"managers,omitempty"` //管理员
}

func (a Application) ToJson() string {
	js, _ := json.Marshal(a)
	return string(js)
}

func (a Application) insert(docs ...interface{}) error {
	return mongo.Insert(shareDB.DocManagerDBName(), appCollection, docs...)
}

func (a Application) isExist(query interface{}) bool {
	return mongo.IsExist(shareDB.DocManagerDBName(), appCollection, query)
}

func (a Application) findOne(query, selector interface{}) (Application, error) {
	ap := Application{}
	err := mongo.FindOne(shareDB.DocManagerDBName(), appCollection, query, selector, &ap)
	return ap, err
}

func (a Application) findAll(query, selector interface{}) (results []Application, err error) {
	results = []Application{}
	err = mongo.FindAll(shareDB.DocManagerDBName(), appCollection, query, selector, &results)
	return results, err
}

func (a Application) totalCount(query, selector interface{}) (int, error) {
	return mongo.TotalCount(shareDB.DocManagerDBName(), appCollection, query, selector)
}

func (a Application) findPage(page, limit int, query, selector interface{}, fields ...string) (results []Application, err error) {
	results = []Application{}
	err = mongo.FindPage(shareDB.DocManagerDBName(), appCollection, page, limit, query, selector, &results, fields...)
	return
}

func (a Application) update(selector, update interface{}) error {
	return mongo.Update(shareDB.DocManagerDBName(), appCollection, selector, update, true)
}

func (a Application) remove(selector interface{}) error {
	version := AppVersion{}
	if version.isExist(bson.M{"app_id": a.Id}) {
		return fmt.Errorf("app in use")
	}
	return mongo.Remove(shareDB.DocManagerDBName(), appCollection, selector)
}

func (a Application) removeAll(selector interface{}) error {
	return mongo.RemoveAll(shareDB.DocManagerDBName(), appCollection, selector)
}

func (a Application) Exist(query interface{}) bool {
	if query == nil {
		query = bson.M{"_id": a.Id}
	}
	return a.isExist(query)
}

func (a *Application) Insert() error {
	if a.BundleId == "" {
		return errors.New("bundleId must fill")
	}
	if a.Icon == "" {
		return errors.New("icon must fill")
	}
	if a.OwnerId == 0 {
		return errors.New("ownerId must fill")
	}
	if a.Name == "" {
		return errors.New("name must fill")
	}
	if len(a.Desc) < 10 {
		return errors.New("desc length must > 10")
	}

	if a.isExist(bson.M{"bundle_id": a.BundleId}) {
		return errors.New("bundle_id already exist")
	}
	if a.isExist(bson.M{"name": a.Name}) {
		return errors.New("name already exist")
	}
	a.Id, _ = mongo.GetIncrementId(shareDB.DocManagerDBName(), appCollection)
	a.CreateTime = time.Now().Unix()
	list := a.Managers
	a.Managers = nil
	a.Owner = nil
	err := a.insert(a)
	if err != nil {
		return err
	}
	a.Managers = list
	a.updateAppManagers()
	return nil
}

func (a Application) updateAppManagers() error {
	aM := model_app_manager.AppManager{}
	aM.RemoveAppId(a.Id)
	for _, user := range a.Managers {
		if user.Id == a.OwnerId {
			continue
		}
		aM.UserId = user.Id
		aM.AppId = a.Id
		aM.Insert()
	}
	return nil
}

func (a Application) Update() error {
	a.BundleId = ""
	a.CreateTime = 0
	a.updateAppManagers()
	a.Managers = nil
	a.Owner = nil
	return a.update(bson.M{"_id": a.Id}, a)
}

func (a Application) Remove() error {
	if a.Id == 0 {
		return errors.New("id is 0")
	}
	//删除app的用户关联数据
	appM := model_app_manager.AppManager{}
	appM.RemoveAppId(a.Id)

	//删除app的模型关联数据
	adm := model_app_data_model.AppDataModel{}
	adm.RemoveAppId(a.Id)

	return a.remove(bson.M{"_id": a.Id})
}

func (a Application) fetchOwnerManagers() (model_user.User, []model_user.User, []string) {
	user := model_user.User{}
	user.Id = a.OwnerId
	user, _ = user.FindOne()

	aM := model_app_manager.AppManager{}
	aMs, _ := aM.FindAll(bson.M{"app_id": a.Id}, nil)
	users := make([]model_user.User, 0)
	managerInfo := make([]string, 0)
	for _, item := range aMs {
		u := model_user.User{}
		u.Id = item.UserId
		u, err := u.FindOne()
		if err == nil {
			users = append(users, u)
			managerInfo = append(managerInfo, u.Username)
		}
	}
	return user, users, managerInfo
}

func (a Application) FindOne(query, selector interface{}) (Application, error) {
	if query == nil {
		query = bson.M{"_id": a.Id}
	}
	a, err := a.findOne(query, selector)
	if err != nil {
		return a, err
	}
	user, managers, _ := a.fetchOwnerManagers()
	a.Owner = &user
	a.Managers = managers
	return a, nil
}

func (a Application) FindSimpleOne(query, selector interface{}) (Application, error) {
	if query == nil {
		query = bson.M{"_id": a.Id}
	}
	a, err := a.findOne(query, selector)
	if err != nil {
		return a, err
	}
	return a, nil
}

func (a Application) FindAll(query, selector interface{}) (apps []Application, err error) {
	apps, err = a.findAll(query, selector)
	makeTreeApplication(apps)
	return
}

//func FindPageApplications(page, limit int, fields ...string) (apps *[]Application, err error) {
//	return appC.findPage(page, limit, nil, nil, fields...)
//}

func (a Application) TotalCount(query, selector interface{}) (int, error) {
	return a.totalCount(query, selector)
}
func (a Application) FindPageFilter(page, limit int, query, selector interface{}, fields ...string) (apps []Application, err error) {
	apps, err = a.findPage(page, limit, query, selector, fields...)
	makeTreeApplication(apps)
	return
}

func makeTreeApplication(apps []Application) {
	for index := range apps {
		user, users, _ := apps[index].fetchOwnerManagers()
		apps[index].Owner = &user
		apps[index].Managers = users
	}
}
