package common

import (
	"fmt"
	"testing"
)

func TestVersionCompare(t *testing.T) {

	result, err := VersionCompare("1.1.1", "")
	if err != nil {
		panic(err)
	}
	fmt.Println("compare: ", result)
	result, err = VersionCompare("1.1.1", "1.1.1.1.1")
	if err != nil {
		panic(err)
	}
	fmt.Println("compare: ", result)

	result, err = VersionCompare("2.1.1", "1.1.2")
	if err != nil {
		panic(err)
	}
	fmt.Println("compare: ", result)

	result, err = VersionCompare("2.1.1.1", "1.1.2")
	if err != nil {
		panic(err)
	}
	fmt.Println("compare: ", result)

	result, err = VersionCompare("2.1.1.1", "2.1.1.1")
	if err != nil {
		panic(err)
	}
	fmt.Println("compare: ", result)

	//fmt.Println(checkVersionOK("1.1.1.1"))
	//fmt.Println(checkVersionOK("1.1..1.1"))
	//fmt.Println(checkVersionOK("a1.1.1.1"))
	//fmt.Println(checkVersionOK("1.1.1."))
	//fmt.Println(checkVersionOK(".1.1.1"))

}

func TestVersionCompare2(t *testing.T) {
	num := TransformVersionToInt("255.255.255.255")
	fmt.Println(num)
}
