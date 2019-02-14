package common

import (
	"errors"
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
	if len(version2) == 0 {
		version2 = "0"
	}
	if len(version1) == 0 {
		version1 = "0"
	}
	if !checkVersionOK(version1) || !checkVersionOK(version2) {
		return CompareVersionStateFailed, errors.New("version string not ok")
	}
	length1 := len(version1)
	length2 := len(version2)
	min := length1
	if min > length2 {
		min = length2
	}
	versionN1, _ := strconv.Atoi(strings.Replace(version1[0:min], ".", "", -1))
	versionN2, _ := strconv.Atoi(strings.Replace(version2[0:min], ".", "", -1))

	if versionN1 > versionN2 {
		return CompareVersionStateGreater, nil
	} else if versionN1 < versionN2 {
		return CompareVersionStateLess, nil
	}
	if length1 == length2 {
		return CompareVersionStateEqual, nil
	} else if length1 > length2 {
		return CompareVersionStateGreater, nil
	} else {
		return CompareVersionStateLess, nil
	}
	return CompareVersionStateFailed, nil
}

func checkVersionOK(version string) bool {
	if strings.HasPrefix(version, ".") || strings.HasSuffix(version, ".") {
		return false
	}
	if strings.Contains(version, "..") {
		return false
	}
	str := strings.Replace(version, ".", "", -1)
	_, err := strconv.Atoi(str)
	if err != nil {
		return false
	}
	return true
}
