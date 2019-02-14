package model_user

import (
	"encoding/json"
	"fmt"
	"go-eladmin/core/mongo"
	"testing"
)

func TestUser_FindRoles(t *testing.T) {
	mongo.DialMgo("127.0.0.1:27017")

	user := User{}
	user.Id = 10000
	results, err := user.FindRoles()
	if err != nil {
		panic(err)
	}
	js, _ := json.Marshal(results)
	fmt.Println(string(js))

}
