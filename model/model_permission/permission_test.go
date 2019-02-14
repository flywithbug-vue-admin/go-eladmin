package model_permission

import (
	"encoding/json"
	"fmt"
	"go-eladmin/model/a_mongo_index"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"go-eladmin/core/mongo"
)

func TestPipe(t *testing.T) {
	mongo.RegisterMongo("127.0.0.1:27017", "doc_manager")
	permission := Permission{}

	list, err := permission.FetchTreeList(nil)
	if err != nil {
		panic(err)
	}
	js, _ := json.Marshal(list)
	fmt.Println(string(js))

}

func TestPipelineFetch(t *testing.T) {
	mongo.RegisterMongo("127.0.0.1:27017", "doc_manager")
	permission := Permission{}

	name := "user"
	sort := bson.M{"$sort": bson.M{"_id": 1}}
	match := bson.M{"$match": bson.M{"pid": 0}}
	if len(name) > 0 {
		match = bson.M{"$match": bson.M{"pid": 0, "name": bson.M{"$regex": name, "$options": "i"}}}
	}
	lookup := bson.M{"$lookup": bson.M{"from": mongo_index.CollectionPermission, "localField": "_id", "foreignField": "pid", "as": "children"}}
	pipeline := []bson.M{
		match,
		sort,
		lookup,
	}
	results, _ := permission.FindPipeline(pipeline)

	js, _ := json.Marshal(results)
	fmt.Println(string(js))

}
