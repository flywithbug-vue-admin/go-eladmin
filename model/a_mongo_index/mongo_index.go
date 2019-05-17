package mongo_index

import (
	"go-eladmin/config"
	"go-eladmin/core/mongo"

	"github.com/flywithbug/log4go"
	"gopkg.in/mgo.v2"
)

type Index struct {
	Collection string
	DBName     string
	Index      mgo.Index
	DropIndex  []string
}

func CreateMgoIndex() {
	Indexes := formatIndex()
	aMCfg := config.Conf().DBConfig
	for _, aMongoIndex := range Indexes {
		_, c := mongo.Collection(aMongoIndex.DBName, aMongoIndex.Collection)
		if len(aMongoIndex.DropIndex) > 0 {
			for _, idxName := range aMongoIndex.DropIndex {
				if err := c.DropIndexName(idxName); err != nil {
					log4go.Warn(err.Error())
				}
			}
		}
		if aMCfg.ForceSync {
			if err := c.DropIndexName(aMongoIndex.Index.Name); err != nil {
				log4go.Warn(err.Error())
			}
		}
		err := c.EnsureIndex(aMongoIndex.Index)
		if err != nil {
			log4go.Warn(err.Error())
		}
	}
}

func formatIndex() []Index {
	Indexes := append(docManagerIndex(), monitorIndex()...)
	Indexes = append(Indexes, devToolsIndex()...)

	return Indexes
}
