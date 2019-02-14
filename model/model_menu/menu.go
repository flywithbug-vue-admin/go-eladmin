package model_menu

import (
	"encoding/json"
	"fmt"
	"go-eladmin/core/mongo"
	"go-eladmin/model/a_mongo_index"
	"go-eladmin/model/model_menu_role"
	"go-eladmin/model/model_role"
	"go-eladmin/model/shareDB"
	"time"

	"github.com/flywithbug/log4go"

	"gopkg.in/mgo.v2/bson"
)

const (
	menuCollection = mongo_index.CollectionMenu
	//MenuMenuAll    = "MENU_ALL"
	MenuPermissionSelect = "MENU_SELECT"
	MenuPermissionCreate = "MENU_CREATE"
	MenuPermissionEdit   = "MENU_EDIT"
	MenuPermissionDelete = "MENU_DELETE"

	MenuTypeList  = 1
	MenuTypeTree  = 2
	MenuTypeBuild = 3
)

type meta struct {
	Title      string `json:"title,omitempty" bson:"title,omitempty"`
	Icon       string `json:"icon,omitempty" bson:"icon,omitempty"`
	NoCache    bool   `json:"noCache,omitempty" bson:"noCache,omitempty"`
	Breadcrumb bool   `json:"breadcrumb,omitempty" bson:"breadcrumb,omitempty"`
}

type Menu struct {
	Id         int64                      `json:"id,omitempty" bson:"_id,omitempty"`
	PId        int64                      `json:"pid,omitempty" bson:"pid"` //父节点ID
	Sort       int                        `json:"sort" bson:"sort"`
	Icon       string                     `json:"icon,omitempty" bson:"icon,omitempty"`
	Name       string                     `json:"name,omitempty" bson:"name,omitempty"`
	Label      string                     `json:"label,omitempty" bson:"label,omitempty"`
	Path       string                     `json:"path,omitempty" bson:"path,omitempty"`
	AlwaysShow bool                       `json:"alwaysShow" bson:"alwaysShow"`
	Component  string                     `json:"component,omitempty" bson:"component,omitempty"`
	IFrame     bool                       `json:"iframe" bson:"iframe"`
	CreateTime int64                      `json:"createTime,omitempty" bson:"createTime,omitempty"`
	Children   []Menu                     `json:"children,omitempty" bson:"children,omitempty"`
	Roles      []model_role.Role          `json:"roles,omitempty" bson:"roles,omitempty"`
	Meta       meta                       `json:"meta,omitempty" bson:"meta,omitempty"`
	menuRoles  []model_menu_role.MenuRole `json:"menu_roles,omitempty"`
}

func (m Menu) ToJson() string {
	js, _ := json.Marshal(m)
	return string(js)
}

func (m Menu) isExist(query interface{}) bool {
	return mongo.IsExist(shareDB.DocManagerDBName(), menuCollection, query)
}

func (m Menu) insert(docs ...interface{}) error {
	return mongo.Insert(shareDB.DocManagerDBName(), menuCollection, docs...)
}

func (m Menu) update(selector, update interface{}) error {
	return mongo.Update(shareDB.DocManagerDBName(), menuCollection, selector, update, true)
}

func (m Menu) findOne(query, selector interface{}) (Menu, error) {
	ap := Menu{}
	err := mongo.FindOne(shareDB.DocManagerDBName(), menuCollection, query, selector, &ap)
	return ap, err
}
func (m Menu) findAll(query, selector interface{}) (results []Menu, err error) {
	results = []Menu{}
	err = mongo.FindAll(shareDB.DocManagerDBName(), menuCollection, query, selector, &results)
	return results, err
}

func (m Menu) remove(selector interface{}) error {
	return mongo.Remove(shareDB.DocManagerDBName(), menuCollection, selector)
}

func (m Menu) removeAll(selector interface{}) error {
	return mongo.RemoveAll(shareDB.DocManagerDBName(), menuCollection, selector)
}

func (m Menu) totalCount(query, selector interface{}) (int, error) {
	return mongo.TotalCount(shareDB.DocManagerDBName(), menuCollection, query, selector)
}

func (m Menu) findPage(page, limit int, query, selector interface{}, fields ...string) (results []Menu, err error) {
	results = []Menu{}
	err = mongo.FindPage(shareDB.DocManagerDBName(), menuCollection, page, limit, query, selector, &results, fields...)
	return
}

func (m Menu) pipeAll(pipeline, result interface{}, allowDiskUse bool) error {
	return mongo.PipeAll(shareDB.DocManagerDBName(), menuCollection, pipeline, result, allowDiskUse)
}

func (m Menu) pipeOne(pipeline, result interface{}, allowDiskUse bool) error {
	return mongo.PipeOne(shareDB.DocManagerDBName(), menuCollection, pipeline, result, allowDiskUse)
}

func (m Menu) explain(pipeline, result interface{}) (results []Menu, err error) {
	err = mongo.Explain(shareDB.DocManagerDBName(), menuCollection, pipeline, result)
	return
}

func (m Menu) Exist() bool {
	return m.isExist(bson.M{"_id": m.Id})
}

func (m Menu) Insert() (int64, error) {
	if m.PId != 0 && !m.isExist(bson.M{"_id": m.PId}) {
		return -1, fmt.Errorf("pid  not exist")
	}
	list := m.Roles
	m.Id, _ = mongo.GetIncrementId(shareDB.DocManagerDBName(), menuCollection)
	m.CreateTime = time.Now().Unix()
	m.Roles = nil
	m.Children = nil
	m.IFrame = true
	m.AlwaysShow = true
	err := m.insert(m)
	if err != nil {
		return -1, err
	}
	js, _ := json.Marshal(m)
	log4go.Info("Marshal:%s", string(js))
	m.Roles = list
	m.updateMenuRole()
	return m.Id, nil
}

func (m Menu) updateMenuRole() {
	mr := model_menu_role.MenuRole{}
	mr.RemoveMenuId(m.Id)
	for _, role := range m.Roles {
		if role.Exist() {
			mr.RoleId = role.Id
			mr.MenuId = m.Id
			mr.Insert()
		}
	}

}

func (m Menu) Update() error {
	if m.PId == m.Id {
		m.PId = 0
	}
	m.updateMenuRole()
	m.Roles = nil
	m.Children = nil
	selector := bson.M{"_id": m.Id}
	if len(m.Children) > 0 {
		m.AlwaysShow = true
	}
	return m.update(selector, m)
}

func (m Menu) Remove() error {
	mr := model_menu_role.MenuRole{}
	mr.RemoveMenuId(m.Id)
	m.removeAllChildren()
	return m.remove(bson.M{"_id": m.Id})
}

func (m Menu) FindAll(query, selector interface{}) (results []Menu, err error) {
	return m.findAll(query, selector)
}

func (m Menu) TotalCount(query, selector interface{}) (int, error) {
	return m.totalCount(query, selector)
}
func (m Menu) FindPageListFilter(page, limit int, query, selector interface{}, fields ...string) ([]Menu, error) {
	results, err := m.findPage(page, limit, query, selector, fields...)
	if err != nil {
		return nil, err
	}
	makeTreeList(results, selector, MenuTypeList)
	return results, err
}

func (m Menu) FindPageTreeFilter(page, limit int, query, selector interface{}, fields ...string) ([]Menu, error) {
	results, err := m.findPage(page, limit, query, selector, fields...)
	if err != nil {
		return nil, err
	}
	makeTreeList(results, selector, MenuTypeTree)
	return results, err
}

func (m Menu) FetchTreeList(selector interface{}) (results []Menu, err error) {
	results, err = m.findAll(bson.M{"pid": 0}, selector)
	if err != nil {
		return
	}
	makeTreeList(results, selector, MenuTypeList)
	return
}

func (m Menu) FindOneTree() (menu Menu, err error) {
	menu, err = m.findOne(bson.M{"_id": m.Id}, nil)
	if err != nil {
		return
	}
	list := []Menu{menu}
	makeTreeList(list, nil, MenuTypeTree)
	return list[0], nil
}

func (m *Menu) findChildren(selector interface{}) error {
	results, err := m.findPage(0, 0, bson.M{"pid": m.Id}, selector, "+sort")
	if err != nil {
		return err
	}
	m.Children = results
	return nil
}

func makeTreeList(list []Menu, selector interface{}, menuType int) {
	for index := range list {
		err := list[index].findChildren(selector)
		if err != nil {
			return
		}
		item := list[index]
		if menuType == MenuTypeList {
			mr := model_menu_role.MenuRole{}
			results, _ := mr.FindAll(bson.M{"menu_id": item.Id}, nil)
			list[index].Roles = make([]model_role.Role, 0)
			var role model_role.Role
			for _, mr = range results {
				role.Id = mr.RoleId
				role, err := role.FindOneTree(nil)
				if err != nil {
					log4go.Info(err.Error())
				} else {
					role.Label = role.Alias
					role.Alias = ""
					list[index].Roles = append(list[index].Roles, role)
				}
			}
		}

		if selector == nil {
			list[index].Meta = meta{
				Title: item.Name,
				Icon:  item.Icon,
			}
		} else {
			//list[index].Roles
			list[index].Label = item.Name
			list[index].Name = ""
		}
		makeTreeList(list[index].Children, selector, menuType)
	}
}

func (m Menu) checkMenuSelectPermission(roles []model_role.Role) bool {
	mr := model_menu_role.MenuRole{}
	results, _ := mr.FindAll(bson.M{"menu_id": m.Id}, nil)
	for _, item := range results {
		for _, role := range roles {
			if role.Id == 10000 {
				return true
			}
			if role.Id == item.RoleId {
				return true
			}
		}
	}
	return false
}

func (m Menu) FindPageBuildFilter(roles []model_role.Role, page, limit int, query, selector interface{}, fields ...string) ([]Menu, error) {
	results, err := m.findPage(page, limit, query, selector, fields...)
	if err != nil {
		return nil, err
	}
	makeRoleTreeList(results, selector, roles)
	return results, err
}

func (m *Menu) findChildrenFilter(selector interface{}, roles []model_role.Role) error {
	results, err := m.findPage(0, 0, bson.M{"pid": m.Id}, selector, "+sort")
	if err != nil {
		return err
	}
	list := make([]Menu, 0)
	for _, item := range results {
		if item.checkMenuSelectPermission(roles) {
			list = append(list, item)
		}
	}
	m.Children = list
	return nil
}

func makeRoleTreeList(list []Menu, selector interface{}, roles []model_role.Role) {
	for index := range list {
		err := list[index].findChildrenFilter(selector, roles)
		if err != nil {
			return
		}
		item := list[index]
		list[index].Meta = meta{
			Title: item.Name,
			Icon:  item.Icon,
		}
		//list[index].Label = item.Name
		makeRoleTreeList(list[index].Children, selector, roles)
	}
}

func (m Menu) removeAllChildren() {
	m.findChildren(nil)
	for index := range m.Children {
		m.Children[index].Remove()
	}
}
