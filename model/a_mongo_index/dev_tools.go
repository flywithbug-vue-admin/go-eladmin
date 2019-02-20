package mongo_index

import (
	"go-eladmin/model/shareDB"

	"gopkg.in/mgo.v2"
)

const (
	CollectionDataModel    = "data_model"
	CollectionAppDataModel = "app_data_model"
	CollectionModule       = "module"
	CollectionModuleApp    = "module_app"
)

func devToolsIndex() []Index {
	var Indexes = []Index{
		{
			DBName:     shareDB.DocManagerDBName(),
			Collection: CollectionDataModel,
			Index: mgo.Index{
				Key:        []string{"name"},
				Unique:     true,
				DropDups:   true,
				Background: false,
				Sparse:     true,
				Name:       "c_data_model_f_name_index",
			},
		},
		{
			DBName:     shareDB.DocManagerDBName(),
			Collection: CollectionAppDataModel,
			Index: mgo.Index{
				Key:        []string{"model_id", "app_id"},
				Unique:     true,
				DropDups:   true,
				Background: false,
				Sparse:     true,
				Name:       "c_modelId_f_appId_index",
			},
		},
		{
			DBName:     shareDB.DocManagerDBName(),
			Collection: CollectionModuleApp,
			Index: mgo.Index{
				Key:        []string{"module_id", "app_id"},
				Unique:     true,
				DropDups:   true,
				Background: false,
				Sparse:     true,
				Name:       "c_moduleId_f_appId_index",
			},
		},
	}
	return Indexes
}
