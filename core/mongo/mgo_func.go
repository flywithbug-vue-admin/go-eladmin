package mongo

import (
	"errors"

	"gopkg.in/mgo.v2"
)

// RegisterMongo register a mongodb connection.
// alias is connection's alias
// url is connection url
// db is connection's default database name
func RegisterMongo(url, db string) error {
	if _, ok := sessionMap[db]; ok {
		return errors.New(ErrExistConnectionDB)
	}
	if db == "" {
		db = "default"
	}
	aSession, err := mgo.Dial(url)
	if err != nil {

		return err
	}
	aMongo := new(tMongo)
	aMongo.session = aSession
	aMongo.db = db
	sessionMap[db] = aMongo
	return nil
}

// NewMongo return a new mongodb operator
func NewMongo(db string) Monger {
	aM := new(tMongo)
	aM.Use(db)
	return aM
}
