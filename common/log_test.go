package common

import (
	"fmt"
	"testing"

	g_file "github.com/flywithbug/file"
)

type FileMap struct {
	DirName string
	files   []string
}

func TestGetAllFile(t *testing.T) {

	//s, err := GetAllFile("/Users/ori/go/src/go-eladmin/log", ".log")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("slice :%v", s)
	g_file.ListDirectory("/Users/ori/go/src/go-eladmin/log", true, func(file g_file.FileInfo, err error) {
		if !file.IsDir() {
			fmt.Println(file.Name())
		}
	})

}
