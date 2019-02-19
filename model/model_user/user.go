package model_user

import (
	"encoding/json"
	"fmt"
	"go-eladmin/core/mongo"
	"go-eladmin/model/a_mongo_index"
	"go-eladmin/model/model_role"
	"go-eladmin/model/model_user_role"
	"go-eladmin/model/shareDB"
	"math/rand"
	"time"

	"github.com/flywithbug/log4go"

	"gopkg.in/mgo.v2/bson"
)

const (
	userCollection = mongo_index.CollectionUser
	//UserPermissionAll    = "USER_ALL"
	UserPermissionSelect = "USER_SELECT"
	UserPermissionCreate = "USER_CREATE"
	UserPermissionEdit   = "USER_EDIT"
	UserPermissionDelete = "USER_DELETE"
)

type User struct {
	Id          int64             `json:"id,omitempty" bson:"_id,omitempty"`
	Username    string            `json:"username,omitempty" bson:"username,omitempty"` //用户名
	Password    string            `json:"password,omitempty" bson:"password,omitempty"`
	Avatar      string            `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Email       string            `json:"email,omitempty" bson:"email,omitempty"`
	Phone       string            `json:"phone,omitempty" bson:"phone,omitempty"`
	Gender      int               `json:"gender,omitempty" bson:"gender,omitempty"` // 1男 2女
	Name        string            `json:"name,omitempty" bson:"name,omitempty"`     // 名字！
	Nick        string            `json:"nick,omitempty" bson:"nick,omitempty"`     // 昵称
	Title       string            `json:"title,omitempty" bson:"title,omitempty"`
	Enabled     bool              `json:"enabled,omitempty" bson:"enabled"`      //1 激活，
	Note        string            `json:"note,omitempty"  bson:"note,omitempty"` //备注,
	CreateTime  int64             `json:"createTime,omitempty"  bson:"create_time,omitempty"`
	Roles       []model_role.Role `json:"roles,omitempty" bson:"roles,omitempty"`
	RolesString []string          `json:"roles_string,omitempty" bson:"roles_string,omitempty"`
	Permissions []string          `json:"permissions,omitempty" bson:"permissions,omitempty"`
}

func (u User) ToJson() string {
	js, _ := json.Marshal(u)
	return string(js)
}

/*
	OperationModel
*/
func (u User) insert(docs ...interface{}) error {
	return mongo.Insert(shareDB.DocManagerDBName(), userCollection, docs...)
}

func (u User) isExist(query interface{}) bool {
	return mongo.IsExist(shareDB.DocManagerDBName(), userCollection, query)
}

func (u User) findOne(query, selector interface{}) (User, error) {
	us := User{}
	err := mongo.FindOne(shareDB.DocManagerDBName(), userCollection, query, selector, &us)
	return us, err
}

func (u User) findAll(query, selector interface{}) (results []User, err error) {
	results = []User{}
	err = mongo.FindAll(shareDB.DocManagerDBName(), userCollection, query, selector, &results)
	return results, err
}

func (u User) findPage(page, limit int, query, selector interface{}, fields ...string) (results []User, err error) {
	results = []User{}
	err = mongo.FindPage(shareDB.DocManagerDBName(), userCollection, page, limit, query, selector, &results, fields...)
	return
}

//data := bson.M{"$set": bson.M{"age": 22}}
func (u User) update(selector, update interface{}) error {
	return mongo.Update(shareDB.DocManagerDBName(), userCollection, selector, update, true)
}

func (u User) remove(selector interface{}) error {
	return mongo.Remove(shareDB.DocManagerDBName(), userCollection, selector)
}

func (u User) removeAll(selector interface{}) error {
	return mongo.RemoveAll(shareDB.DocManagerDBName(), userCollection, selector)
}

/*
	userModify
*/
func (u *User) Insert() error {
	u.Id, _ = mongo.GetIncrementId(shareDB.DocManagerDBName(), userCollection)
	u.Enabled = true
	u.CreateTime = time.Now().Unix()
	list := u.Roles
	u.Roles = nil
	if len(u.Avatar) == 0 {
		u.Avatar = "https://s2.ax1x.com/2019/01/12/FjDbjg.png"
	}
	if u.Password == "" {
		u.Password = createCaptcha()
	}
	err := u.insert(u)
	if err != nil {
		return err
	}

	u.Roles = list
	u.updateUserRoles()
	return nil
}

func (u User) updateUserRoles() {
	if len(u.Roles) == 0 {
		return
	}
	ur := model_user_role.UserRole{}
	ur.RemoveUserId(u.Id)
	for _, role := range u.Roles {
		if role.Exist() {
			ur.RoleId = role.Id
			ur.UserId = u.Id
			ur.Insert()
		}
	}
}

func (u User) Update() error {
	selector := bson.M{"_id": u.Id}
	u.updateUserRoles()
	u.Roles = nil
	err := u.update(selector, u)
	return err
}

func (u User) UpdateAvatar() error {
	selector := bson.M{"_id": u.Id}
	if len(u.Avatar) == 0 {
		return fmt.Errorf("avata is null")
	}
	return u.update(selector, bson.M{"avatar": u.Avatar})
}

func (u User) Remove() error {
	if u.Id == 10000 {
		return fmt.Errorf("超级管理员无法被删除")
	}
	selector := bson.M{"_id": u.Id}
	err := u.remove(selector)
	if err != nil {
		return err
	}
	ur := model_user_role.UserRole{}
	ur.RemoveUserId(u.Id)
	return nil
}

func (u User) FindAll() ([]User, error) {
	return u.findAll(nil, bson.M{"password": 0})
}

func (u User) FindOne() (User, error) {
	return u.findOne(bson.M{"_id": u.Id}, bson.M{"password": 0})
}

func (u User) FindTreeOne() (User, error) {
	query := bson.M{"_id": u.Id}
	selector := bson.M{"password": 0}
	u, err := u.findOne(query, selector)
	if err != nil {
		return u, err
	}
	list := []User{u}
	makeTreeList(list, selector)
	return list[0], nil
}

func LoginUser(username, pass string) (user User, err error) {
	user = User{}
	user, err = user.findOne(bson.M{"username": username, "password": pass}, bson.M{"password": 0})
	return
}

func (u User) totalCount(query, selector interface{}) (int, error) {
	return mongo.TotalCount(shareDB.DocManagerDBName(), userCollection, query, selector)
}

func (u User) FindPageFilter(page, limit int, query, selector interface{}, fields ...string) ([]User, error) {
	results, err := u.findPage(page, limit, query, selector, fields...)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (u User) FindPageTreeFilter(page, limit int, query, selector interface{}, fields ...string) ([]User, error) {
	results, err := u.findPage(page, limit, query, selector, fields...)
	if err != nil {
		return nil, err
	}
	makeTreeList(results, selector)
	return results, nil
}

func (u User) TotalCount(query, selector interface{}) (int, error) {
	return u.totalCount(query, selector)
}

func makeTreeList(list []User, selector interface{}) error {
	for index := range list {
		ur := model_user_role.UserRole{}
		results, _ := ur.FindAll(bson.M{"user_id": list[index].Id}, selector)
		list[index].Roles = make([]model_role.Role, 0)
		rolesString := make([]string, 0)
		var role model_role.Role
		for _, item := range results {
			role.Id = item.RoleId
			role, err := role.FindOneTree(nil)
			role.Label = role.Alias
			role.Alias = ""
			if err != nil {
				log4go.Info(err.Error())
			} else {
				list[index].Roles = append(list[index].Roles, role)
				rolesString = append(rolesString, role.PerStrings...)
			}
		}
		list[index].RolesString = rolesString
	}

	return nil
}

func createCaptcha() string {
	return fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
}

func (u User) CheckPassword() bool {
	if u.isExist(bson.M{"password": u.Password, "_id": u.Id}) {
		return true
	}
	return false
}

func (u User) UpdatePassword() error {
	return u.update(bson.M{"_id": u.Id}, bson.M{"password": u.Password})
}

func (u User) UpdateMail() error {
	return u.update(bson.M{"_id": u.Id}, bson.M{"email": u.Email})
}

func (u User) FindRoles() (User, error) {
	u, err := u.findOne(bson.M{"_id": u.Id}, bson.M{"password": 0})
	if err != nil {
		return u, err
	}
	ur := model_user_role.UserRole{}
	ur.UserId = u.Id
	results, _ := ur.FindAll(bson.M{"user_id": u.Id}, nil)
	u.Roles = make([]model_role.Role, 0)
	var role model_role.Role
	for _, item := range results {
		role.Id = item.RoleId
		role, err := role.FindSimple(bson.M{"_id": 1, "name": 1})
		role.Label = role.Alias
		role.Alias = ""
		if err != nil {
			log4go.Info(err.Error())
		} else {
			u.Roles = append(u.Roles, role)
		}
	}
	return u, nil
}
