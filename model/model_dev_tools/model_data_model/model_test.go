package model_data_model

import (
	"encoding/json"
	"fmt"
	"go-eladmin/core/mongo"
	"go-eladmin/model/a_mongo_index"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func TestDataModel(t *testing.T) {
	mongo.RegisterMongo("127.0.0.1:27017", "doc_manager")

	dataModel := DataModel{}
	dataModel.Name = "User1"
	dataModel.Desc = "test"
	//dataModel.Attributes = make([]Attributes, 0)

	a := Attribute{}
	a.Name = "Name"
	a.Type = modelAttributeTypeString
	a1 := Attribute{}
	a1.Name = "Name1"
	a1.Type = modelAttributeTypeString
	dataModel.Attributes = append(dataModel.Attributes, a, a1)
	_, err := dataModel.Insert()
	if err != nil {
		panic(err)
	}

}

func TestDataModel_Update(t *testing.T) {
	mongo.RegisterMongo("127.0.0.1:27017", "doc_manager")
	dataModel := DataModel{}
	dataModel.Name = "User"
	dataModel.Desc = "test"
	dataModel.Id = 10021
	a := Attribute{}
	a.Name = "Name2"
	a.Type = modelAttributeTypeInt
	err := dataModel.AddAttributes([]Attribute{a})
	if err != nil {
		panic(err)
	}
}

func TestDataModel_RemoveAttribute(t *testing.T) {
	mongo.RegisterMongo("127.0.0.1:27017", "doc_manager")
	dataModel := DataModel{}
	dataModel.Id = 10021
	a := Attribute{}
	a.Name = "Name2"
	err := dataModel.RemoveAttribute(a)
	if err != nil {
		panic(err)
	}
}

func TestPipe(t *testing.T) {
	mongo.RegisterMongo("127.0.0.1:27017", "doc_manager")
	dm := DataModel{}

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
	results, _ := dm.FindPipeline(pipeline)

	js, _ := json.Marshal(results)
	fmt.Println(string(js))
}
