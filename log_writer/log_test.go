package log_writer

import (
	"fmt"
	"go-eladmin/core/mongo"
	"testing"
)

func TestLog_Insert(t *testing.T) {
	mongo.RegisterMongo("127.0.0.1:27017", "log")

	log1 := new(Log)
	log1.Code = "23232"

	log := new(Log)
	log.Code = "23232"
	fmt.Println(log)
	log.Insert()
}
