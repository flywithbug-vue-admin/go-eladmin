package mongo_index

import (
	"go-eladmin/model/shareDB"

	"gopkg.in/mgo.v2"
)

const (
	CollectionUser           = "user"
	CollectionLogin          = "login"
	CollectionApp            = "application"
	CollectionAppManager     = "app_manager"
	CollectionAppVersion     = "app_version"
	CollectionPermission     = "permission"
	CollectionRole           = "role"
	CollectionRolePermission = "role_permission"
	CollectionUserRole       = "user_role"
	CollectionVerify         = "verify"
	CollectionMenu           = "menu"
	CollectionMenuRole       = "menu_role"
	CollectionFile           = "file"
	CollectionPicture        = "picture"
)

func docManagerIndex() []Index {
	//唯一约束
	var Indexes = []Index{
		{
			DBName:     shareDB.DocManagerDBName(),
			Collection: CollectionFile,
			Index: mgo.Index{
				Key:        []string{"md5"},
				Unique:     true,
				DropDups:   true,
				Background: false,
				Sparse:     true,
				Name:       "c_file_f_md5_index",
			},
		},
		{
			DBName:     shareDB.DocManagerDBName(),
			Collection: CollectionPicture,
			Index: mgo.Index{
				Key:        []string{"md5"},
				Unique:     true,
				DropDups:   true,
				Background: false,
				Sparse:     true,
				Name:       "c_picture_f_md5_index",
			},
		},
	}
	Indexes = append(Indexes, userPermissionIndex()...)
	Indexes = append(Indexes, appIndex()...)
	return Indexes
}

func appIndex() []Index {
	//唯一约束
	var Indexes = []Index{
		{
			DBName:     shareDB.DocManagerDBName(),
			Collection: CollectionApp,
			Index: mgo.Index{
				Key:        []string{"bundle_id"},
				Unique:     true,
				DropDups:   true,
				Background: false,
				Sparse:     true,
				Name:       "c_app_f_bundle_id_index",
			},
		},
		{
			DBName:     shareDB.DocManagerDBName(),
			Collection: CollectionApp,
			Index: mgo.Index{
				Key:        []string{"name"},
				Unique:     true,
				DropDups:   true,
				Background: false,
				Sparse:     true,
				Name:       "c_app_f_name_index",
			},
		},
		{
			DBName:     shareDB.DocManagerDBName(),
			Collection: CollectionAppManager,
			Index: mgo.Index{
				Key:        []string{"user_id", "app_id"},
				Unique:     true,
				DropDups:   true,
				Background: false,
				Sparse:     true,
				Name:       "c_userId_f_appId_index",
			},
		},
		{
			DBName:     shareDB.DocManagerDBName(),
			Collection: CollectionAppVersion,
			Index: mgo.Index{
				Key:        []string{"version", "app_id"},
				Unique:     true,
				DropDups:   true,
				Background: false,
				Sparse:     true,
				Name:       "c_appVersion_f_version_f_appId_index",
			},
		},
	}
	return Indexes
}

func userPermissionIndex() []Index {
	//唯一约束
	var Indexes = []Index{
		{
			DBName:     shareDB.DocManagerDBName(),
			Collection: CollectionPermission,
			Index: mgo.Index{
				Key:        []string{"alias"},
				Unique:     true,
				DropDups:   true,
				Background: false,
				Sparse:     true,
				Name:       "c_permission_f_alias_index",
			},
		},
		{
			DBName:     shareDB.DocManagerDBName(),
			Collection: CollectionPermission,
			Index: mgo.Index{
				Key:        []string{"name"},
				Unique:     true,
				DropDups:   true,
				Background: false,
				Sparse:     true,
				Name:       "c_permission_f_name_index",
			},
		},
		{
			DBName:     shareDB.DocManagerDBName(),
			Collection: CollectionUser,
			Index: mgo.Index{
				Key:        []string{"username"},
				Unique:     true,
				DropDups:   true,
				Background: false,
				Sparse:     true,
				Name:       "c_user_f_username_index",
			},
		},
		{
			DBName:     shareDB.DocManagerDBName(),
			Collection: CollectionUser,
			Index: mgo.Index{
				Key:        []string{"email"},
				Unique:     true,
				DropDups:   true,
				Background: false,
				Sparse:     true,
				Name:       "c_user_f_email_index",
			},
		},
		{
			DBName:     shareDB.DocManagerDBName(),
			Collection: CollectionRole,
			Index: mgo.Index{
				Key:        []string{"name"},
				Unique:     true,
				DropDups:   true,
				Background: false,
				Sparse:     true,
				Name:       "c_role_f_name_index",
			},
		},
		{
			DBName:     shareDB.DocManagerDBName(),
			Collection: CollectionRole,
			Index: mgo.Index{
				Key:        []string{"alias"},
				Unique:     true,
				DropDups:   true,
				Background: false,
				Sparse:     true,
				Name:       "c_role_f_alias_index",
			},
		},
		{
			DBName:     shareDB.DocManagerDBName(),
			Collection: CollectionRolePermission,
			Index: mgo.Index{
				Key:        []string{"role_id", "permission_id"},
				Unique:     true,
				DropDups:   true,
				Background: false,
				Sparse:     true,
				Name:       "c_role_f_permission_f_role_index",
			},
		},
		{
			DBName:     shareDB.DocManagerDBName(),
			Collection: CollectionUserRole,
			Index: mgo.Index{
				Key:        []string{"role_id", "user_id"},
				Unique:     true,
				DropDups:   true,
				Background: false,
				Sparse:     true,
				Name:       "c_user_f_role_index",
			},
		},
		{
			DBName:     shareDB.DocManagerDBName(),
			Collection: CollectionMenuRole,
			Index: mgo.Index{
				Key:        []string{"role_id", "menu_id"},
				Unique:     true,
				DropDups:   true,
				Background: false,
				Sparse:     true,
				Name:       "c_menu_f_role_index",
			},
		},
	}
	return Indexes
}
