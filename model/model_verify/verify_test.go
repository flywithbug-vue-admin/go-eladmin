package model_verify

import (
	"fmt"
	"go-eladmin/core/mongo"
	"testing"
)

func TestGenerVerifyData(t *testing.T) {
	mongo.RegisterMongo("127.0.0.1:27017", "doc_manager")

	source := "addadad232323a"
	vCode, err := GeneralVerifyData(source)
	fmt.Println(vCode)

	if err != nil {
		panic(err)
	}

	fmt.Println("result:", CheckVerify(source, vCode))
	//time.Sleep(time.Second * 2)
	fmt.Println("result:", CheckVerify(source, vCode))

}
