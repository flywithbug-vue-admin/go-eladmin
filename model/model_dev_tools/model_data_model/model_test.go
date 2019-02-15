package model_data_model

import (
	"go-eladmin/core/mongo"
	"testing"
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

}
