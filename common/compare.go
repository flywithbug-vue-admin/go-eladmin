package common

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type CompareState int

const (
	CompareVersionStateFailed CompareState = iota
	CompareVersionStateGreater
	CompareVersionStateEqual
	CompareVersionStateLess
)

var versionRe = regexp.MustCompile(`[\d.]`)

//版本号规则xx.xx.xx 只能有数字和点 version1
func VersionCompare(version1, version2 string) (CompareState, error) {
	vNum1 := TransformVersionToInt(version1)
	if vNum1 == -1 {
		return CompareVersionStateFailed, fmt.Errorf("%s not right (note: x.x.x or x.x.x.x)", version1)
	}
	vNum2 := TransformVersionToInt(version2)
	if vNum2 == -1 {
		return CompareVersionStateFailed, fmt.Errorf("%s not right (note: x.x.x or x.x.x.x)", version2)
	}
	if vNum1 > vNum2 {
		return CompareVersionStateGreater, nil
	} else if vNum1 < vNum2 {
		return CompareVersionStateLess, nil
	} else {
		return CompareVersionStateLess, nil
	}
}

func TransformVersionToInt(version string) int {
	list := strings.Split(version, ".")
	if len(list) > 4 {
		return -1
	}
	total := 0
	versions := []int{0, 0, 0, 0}
	for index := range list {
		num, err := strconv.Atoi(list[index])
		if err != nil {
			return -1
		}
		if num > 255 {
			return -1
		}
		versions[index] = num
	}
	total = versions[0]<<24 | versions[1]<<16 | versions[2]<<8 | versions[3]<<1
	return total
}
