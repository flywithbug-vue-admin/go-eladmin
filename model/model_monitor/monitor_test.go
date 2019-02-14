package model_monitor

import (
	"fmt"
	"go-eladmin/core/mongo"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func TestLog_TotalCount(t *testing.T) {
	mongo.RegisterMongo("127.0.0.1:27017", "monitor")
	visit := VisitApi{}
	query := bson.M{"time_date": bson.M{"$regex": "2019-01-18", "$options": "i"}}
	count, _ := visit.TotalSumCount(query)
	fmt.Println(count)
}
