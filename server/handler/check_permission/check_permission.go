package check_permission

import (
	"go-eladmin/common"
	"go-eladmin/model/model_app"
	"go-eladmin/model/model_user"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	SUPERADMIN = "ADMIN"
)

func CheckNoPermission(c *gin.Context, permission string) bool {
	id := common.UserId(c)
	if id == 10000 {
		return false
	}
	user := model_user.User{}
	user.Id = id
	user, err := user.FindTreeOne()
	if err != nil {
		return true
	}
	for index := range user.RolesString {
		item := user.RolesString[index]
		if strings.EqualFold(item, SUPERADMIN) {
			return false
		}
		if strings.EqualFold(item, permission) {
			return false
		}
		if strings.HasSuffix(item, "ALL") {
			splits := strings.Split(item, "_")
			if strings.HasPrefix(permission, splits[0]) {
				return false
			}
		}
	}
	return true
}

func CheckNoAppManagerPermission(c *gin.Context, app model_app.Application) bool {
	app, err := app.FindOne(nil, nil)
	if err != nil {
		return true
	}

	userId := common.UserId(c)
	if app.Owner.Id == userId {
		return false
	}
	for _, item := range app.Managers {
		if item.Id == userId {
			return false
		}
	}
	return true
}

func CheckNoAppVersionManagerPermission(c *gin.Context, appV model_app.AppVersion) bool {
	app := model_app.Application{}
	app.Id = appV.AppId
	app, err := app.FindOne(nil, nil)
	if err != nil {
		return true
	}
	userId := common.UserId(c)
	if app.Owner.Id == userId {
		return false
	}
	for _, item := range app.Managers {
		if item.Id == userId {
			return false
		}
	}
	return true
}
