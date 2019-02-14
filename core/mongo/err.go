package mongo

import (
	"strings"
)

const (
	ErrExistConnectionDB      = "exist connection db"
	ErrNoConnection           = "no connection"
	ErrCannotSwitchCollection = "can not switch collection"
	ErrMongoObjDestroyed      = "the mongo object has been destoryed"
)

func EqualError(err error, str string) bool {
	return str == err.Error() || strings.HasPrefix(err.Error(), str) || strings.Index(err.Error(), str) > -1
}
