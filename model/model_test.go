package model

import (
	"github.com/flywithbug/mongo"
	"testing"
)

func TestUserFunctions(t *testing.T) {
	mongo.RegisterMongo("doc:doc11121014a@118.89.108.25:27017/docmanager", "docmanager")
	//mongo.DialMgo("doc:doc11121014a@118.89.108.25:27017/docmanager")

	//user := new(User)
	//user.Password = "admin"
	//user.username = "admin"
	//user.Title = "CEO"
	//user.Phone = "129"
	//user.Email = "admin@admin.com"
	//user.Sex = 1
	//
	//err := user.Insert()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(user)

}
