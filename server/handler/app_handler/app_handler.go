package app_handler

import (
	"go-eladmin/common"
	"go-eladmin/model"
	"go-eladmin/model/model_app"
	"go-eladmin/server/handler/check_permission"
	"go-eladmin/server/handler/handler_common"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/flywithbug/log4go"
	"github.com/gin-gonic/gin"
)

type appPara struct {
	model_app.Application
}

func addHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	if check_permission.CheckNoPermission(c, model_app.ApplicationPermissionCreate) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusForbidden, "has no permission")
		return
	}
	app := new(appPara)
	err := c.BindJSON(app)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid: "+err.Error())
		return
	}
	c.Set(common.KeyContextPara, app.ToJson())
	if app.BundleId == "" {
		aRes.SetErrorInfo(http.StatusBadRequest, "BundleId must fill")
		return
	}
	if app.Icon == "" {
		aRes.SetErrorInfo(http.StatusBadRequest, "Icon must fill")
		return
	}
	if app.Name == "" {
		aRes.SetErrorInfo(http.StatusBadRequest, "Name must fill")
		return
	}
	if len(app.Desc) < 10 {
		aRes.SetErrorInfo(http.StatusBadRequest, "Desc must fill")
		return
	}
	app.OwnerId = common.UserId(c)
	err = app.Insert()
	if err != nil {
		log4go.Error(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, err.Error())
		return
	}
	aRes.AddResponseInfo("app", app)
}

func editHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	//if check_permission.CheckNoPermission(c, model_app.ApplicationPermissionEdit) {
	//	log4go.Info(handler_common.RequestId(c) + "has no permission")
	//	aRes.SetErrorInfo(http.StatusForbidden, "has no permission")
	//	return
	//}
	para := new(model_app.Application)
	err := c.BindJSON(para)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid: "+err.Error())
		return
	}
	c.Set(common.KeyContextPara, para.ToJson())
	if check_permission.CheckNoAppManagerPermission(c, *para) &&
		check_permission.CheckNoPermission(c, model_app.ApplicationPermissionEdit) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusForbidden, "has no permission")
		return
	}
	err = para.Update()
	if err != nil {
		log4go.Error(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "update failed: "+err.Error())
		return
	}
	aRes.SetSuccessInfo(http.StatusOK, "success")
}

func listHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	if check_permission.CheckNoPermission(c, model_app.ApplicationPermissionSelect) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusForbidden, "has no permission")
		return
	}
	limit, _ := strconv.Atoi(c.Query("size"))
	page, _ := strconv.Atoi(c.Query("page"))
	sort := c.Query("sort")
	name := c.Query("name")
	owner := c.Query("owner")

	if strings.EqualFold(sort, "-id") {
		sort = "-_id"
	} else if strings.EqualFold(sort, "+id") {
		sort = "+_id"
	} else if len(sort) == 0 {
		sort = "+_id"
	}
	if limit == 0 {
		limit = 10
	}
	if page != 0 {
		page--
	}
	userId := common.UserId(c)
	if userId <= 0 {
		aRes.SetErrorInfo(http.StatusUnauthorized, "user not found")
		return
	}
	query := bson.M{}
	if len(name) > 0 {
		query["name"] = bson.M{"$regex": name, "$options": "i"}
	}
	if len(owner) > 0 {
		query["owner"] = bson.M{"$regex": owner, "$options": "i"}
	}
	selector := bson.M{"_id": 1, "name": 1, "bundle_id": 1, "icon": 1, "create_time": 1, "owner_id": 1, "desc": 1}
	var app = model_app.Application{}
	totalCount, _ := app.TotalCount(query, nil)
	appList, err := app.FindPageFilter(page, limit, query, selector, sort)
	if err != nil {
		log4go.Error(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusUnauthorized, "apps find error"+err.Error())
		return
	}
	aRes.AddResponseInfo("list", appList)
	aRes.AddResponseInfo("total", totalCount)
}

func simpleListHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	if check_permission.CheckNoPermission(c, model_app.ApplicationPermissionSelect) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusForbidden, "has no permission")
		return
	}
	var app = model_app.Application{}
	selector := bson.M{"_id": 1, "name": 1, "icon": 1, "create_time": 1, "owner_id": 1}
	arrList, err := app.FindAll(nil, selector)
	if err != nil {
		log4go.Info(err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, err.Error())
		return
	}
	aRes.AddResponseInfo("list", arrList)
}

func delHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	if check_permission.CheckNoPermission(c, model_app.ApplicationPermissionDelete) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusForbidden, "has no permission")
		return
	}
	para := new(model_app.Application)
	err := c.BindJSON(para)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid: "+err.Error())
		return
	}
	if check_permission.CheckNoAppManagerPermission(c, *para) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusForbidden, "has no permission")
		return
	}
	c.Set(common.KeyContextPara, para.ToJson())
	err = para.Remove()
	if err != nil {
		log4go.Error(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusInternalServerError, "para invalid: "+err.Error())
		return
	}
	aRes.SetSuccess()
}
