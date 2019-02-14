package model_permission

import (
	"encoding/json"
	"fmt"
	"go-eladmin/core/mongo"
	"go-eladmin/model/a_mongo_index"
	"go-eladmin/model/model_role_permission"
	"go-eladmin/model/shareDB"

	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	permissionCollection = mongo_index.CollectionPermission
	//PPermissionAll    = "PERMISSION_ALL"
	PPermissionSelect = "PERMISSION_SELECT"
	PPermissionCreate = "PERMISSION_CREATE"
	PPermissionEdit   = "PERMISSION_EDIT"
	PPermissionDelete = "PERMISSION_DELETE"
)

type Permission struct {
	Id         int64        `json:"id,omitempty" bson:"_id,omitempty"`
	PId        int64        `json:"pid,omitempty" bson:"pid"`
	Name       string       `json:"name,omitempty" bson:"name,omitempty"`
	Alias      string       `json:"alias,omitempty" bson:"alias,omitempty"`
	Label      string       `json:"label,omitempty" bson:"_,omitempty"`
	Note       string       `json:"note,omitempty" bson:"note,omitempty"`
	CreateTime int64        `json:"create_time,omitempty" bson:"create_time,omitempty"`
	Children   []Permission `json:"children,omitempty" bson:"children,omitempty"`
}

func (p Permission) ToJson() string {
	js, _ := json.Marshal(p)
	return string(js)
}

func (p Permission) isExist(query interface{}) bool {
	return mongo.IsExist(shareDB.DocManagerDBName(), permissionCollection, query)
}

func (p Permission) insert(docs ...interface{}) error {
	return mongo.Insert(shareDB.DocManagerDBName(), permissionCollection, docs...)
}

func (p Permission) update(selector, update interface{}) error {
	return mongo.Update(shareDB.DocManagerDBName(), permissionCollection, selector, update, true)
}

func (p Permission) findOne(query, selector interface{}) (Permission, error) {
	ap := Permission{}
	err := mongo.FindOne(shareDB.DocManagerDBName(), permissionCollection, query, selector, &ap)
	return ap, err
}
func (p Permission) findAll(query, selector interface{}) (results []Permission, err error) {
	results = []Permission{}
	err = mongo.FindAll(shareDB.DocManagerDBName(), permissionCollection, query, selector, &results)
	return results, err
}

func (p Permission) remove(selector interface{}) error {
	return mongo.Remove(shareDB.DocManagerDBName(), permissionCollection, selector)
}

func (p Permission) removeAll(selector interface{}) error {
	return mongo.RemoveAll(shareDB.DocManagerDBName(), permissionCollection, selector)
}

func (p Permission) totalCount(query, selector interface{}) (int, error) {
	return mongo.TotalCount(shareDB.DocManagerDBName(), permissionCollection, query, selector)
}

func (p Permission) findPage(page, limit int, query, selector interface{}, fields ...string) (results []Permission, err error) {
	results = []Permission{}
	err = mongo.FindPage(shareDB.DocManagerDBName(), permissionCollection, page, limit, query, selector, &results, fields...)
	return
}

func (p Permission) pipeAll(pipeline, result interface{}, allowDiskUse bool) error {
	return mongo.PipeAll(shareDB.DocManagerDBName(), permissionCollection, pipeline, result, allowDiskUse)
}

func (p Permission) pipeOne(pipeline, result interface{}, allowDiskUse bool) error {
	return mongo.PipeOne(shareDB.DocManagerDBName(), permissionCollection, pipeline, result, allowDiskUse)
}

func (p Permission) explain(pipeline, result interface{}) (results []Permission, err error) {
	err = mongo.Explain(shareDB.DocManagerDBName(), permissionCollection, pipeline, result)
	return
}

func (p Permission) Exist() bool {
	return p.isExist(bson.M{"_id": p.Id})
}

func (p Permission) Insert() (int64, error) {
	if p.PId != 0 && !p.isExist(bson.M{"_id": p.PId}) {
		return -1, fmt.Errorf("pid  not exist")
	}
	p.Id, _ = mongo.GetIncrementId(shareDB.DocManagerDBName(), permissionCollection)
	p.Children = nil
	p.CreateTime = time.Now().Unix()
	return p.Id, p.insert(p)
}

func (p Permission) FindOne(selector interface{}) (Permission, error) {
	return p.findOne(bson.M{"_id": p.Id}, selector)
}

func (p Permission) FindPipeOne() (Permission, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"_id": p.Id}},
		{"$lookup": bson.M{"from": permissionCollection, "localField": "_id", "foreignField": "pid", "as": "children"}},
	}
	err := p.pipeOne(pipeline, &p, true)
	return p, err
}

func (p Permission) Update() error {
	if p.PId == p.Id {
		p.PId = 0
	}
	if p.Id == 10000 {
		return fmt.Errorf("此权限不能编辑")
	}
	if p.PId != 0 && !p.isExist(bson.M{"_id": p.PId}) {
		return fmt.Errorf("pid not exist")
	}
	p.Children = nil
	return p.update(bson.M{"_id": p.Id}, p)
}

func (p Permission) Remove() error {
	if p.Id == 10000 {
		return fmt.Errorf("此权限不能编辑")
	}
	if p.checkInUse() {
		return fmt.Errorf("permission in use")
	}
	p.removeAllChildren()
	return p.remove(bson.M{"_id": p.Id})
}

func (p Permission) checkInUse() bool {
	p.findChildren(bson.M{"_id": 1})
	rp := model_role_permission.RolePermission{}
	if rp.Exist(bson.M{"permission_id": p.Id}) {
		return true
	}
	for _, item := range p.Children {
		if rp.Exist(bson.M{"permission_id": item.Id}) {
			return true
		}
	}
	return false
}

func (p Permission) TotalCount(query, selector interface{}) (int, error) {
	return p.totalCount(query, selector)
}
func (p Permission) FindPageFilter(page, limit int, query, selector interface{}, fields ...string) ([]Permission, error) {
	results, err := p.findPage(page, limit, query, selector, fields...)
	if err != nil {
		return nil, err
	}
	err = makeTreeList(results, selector)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (p Permission) FindAll(query, selector interface{}) (apps []Permission, err error) {
	return p.findAll(query, selector)
}

func (p Permission) FindPipeline(pipeline []bson.M) (results []Permission, err error) {
	results = make([]Permission, 0)
	err = p.pipeAll(pipeline, &results, true)
	return
}

func (p Permission) FetchTreeList(selector interface{}) (results []Permission, err error) {
	results, err = p.findAll(bson.M{"pid": 0}, selector)
	if err != nil {
		return
	}
	err = makeTreeList(results, selector)
	return
}

func (p *Permission) findChildren(selector interface{}) error {
	results, err := p.findAll(bson.M{"pid": p.Id}, selector)
	if err != nil {
		return err
	}
	p.Children = results
	return nil
}

func makeTreeList(list []Permission, selector interface{}) error {
	for index := range list {
		if selector != nil {
			list[index].Label = list[index].Alias
			list[index].Alias = ""
		}
		err := list[index].findChildren(selector)
		if err != nil {
			return err
		}
		err = makeTreeList(list[index].Children, selector)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p Permission) removeAllChildren() {
	p.findChildren(nil)
	for index := range p.Children {
		p.Children[index].Remove()
	}
}
