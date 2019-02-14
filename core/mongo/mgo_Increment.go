package mongo

import (
	"fmt"

	"github.com/flywithbug/log4go"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	collectionsIncrementIds = "increment_ids"
	IncrementIdsStarValue   = 10000
)

var (
	IncrementDBStarId int64 = IncrementIdsStarValue
)

func SetIncrementDBStarId(startId int64) {
	IncrementDBStarId = startId
}

func makeIncrementName(incrementName string) string {
	return fmt.Sprintf("increment_%s", incrementName)
}

type incrementIds struct {
	Id            string `json:"_id" bson:"_id"`                       //自增标记Id
	SequenceValue int64  `json:"sequence_value" bson:"sequence_value"` //自增 id值
}

func createIncrementIds(dbName, incrementName string) error {
	cou := incrementIds{
		Id:            incrementName,
		SequenceValue: IncrementDBStarId,
	}
	return Insert(dbName, collectionsIncrementIds, cou)
}

func GetIncrementId(dbName, incrementName string) (int64, error) {
	incrementName = makeIncrementName(incrementName)
	increment := new(incrementIds)
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"sequence_value": 1}},
		ReturnNew: true,
	}
	ms, c := connect(dbName, collectionsIncrementIds)
	defer ms.Close()
	_, err := c.Find(bson.M{"_id": incrementName}).Apply(change, increment)
	if err != nil {
		err = createIncrementIds(dbName, incrementName)
		if err != nil {
			log4go.Error(err.Error())
			return -1, err
		}
		increment.SequenceValue = IncrementDBStarId
	}
	return increment.SequenceValue, nil
}
