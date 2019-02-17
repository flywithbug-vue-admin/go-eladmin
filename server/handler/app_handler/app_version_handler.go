package app_handler

import (
	"encoding/json"
	"fmt"
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

func addAppVersionHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	if check_permission.CheckNoPermission(c, model_app.ApplicationPermissionCreate) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusBadRequest, "has no permission")
		return
	}
	para := new(model_app.AppVersion)
	err := c.BindJSON(para)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid: "+err.Error())
		return
	}
	c.Set(common.KeyContextPara, para.ToJson())
	err = para.Insert()

	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusInternalServerError, "para invalid: "+err.Error())
		return
	}

	aRes.AddResponseInfo("app", para)
}

func updateAppVersionHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()

	para := new(model_app.AppVersion)
	err := c.BindJSON(para)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid: "+err.Error())
		return
	}

	c.Set(common.KeyContextPara, para.ToJson())
	js, _ := json.Marshal(para)
	fmt.Println(string(js))
	appVersion, _ := para.FindOne()
	if check_permission.CheckNoAppVersionManagerPermission(c, appVersion) &&
		check_permission.CheckNoPermission(c, model_app.ApplicationPermissionEdit) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusBadRequest, "has no permission")
		return
	}
	err = para.Update()

	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusInternalServerError, "para invalid: "+err.Error())
		return
	}
	aRes.SetSuccess()
}

func getAppVersionListHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	if check_permission.CheckNoPermission(c, model_app.ApplicationPermissionSelect) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusBadRequest, "has no permission")
		return
	}

	appId, _ := strconv.Atoi(c.Query("app_id"))
	limit, _ := strconv.Atoi(c.Query("size"))
	page, _ := strconv.Atoi(c.Query("page"))
	sort := c.Query("sort")
	version := c.Query("version")
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
		log4go.Info(handler_common.RequestId(c) + "user not found")
		aRes.SetErrorInfo(http.StatusUnauthorized, "user not found")
		return
	}
	query := bson.M{}
	if appId > 0 {
		query = bson.M{"app_id": appId}
	}

	if len(version) > 0 {
		query["version"] = bson.M{"$gte": version}
	}

	var appV = model_app.AppVersion{}
	totalCount, _ := appV.TotalCount(query, nil)
	appList, err := appV.FindPageFilter(page, limit, query, nil, sort)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusUnauthorized, "app version list find error"+err.Error())
		return
	}
	for _, appV := range appList {
		appV.VersionNum = common.TransformVersionToInt(appV.Version)
		appV.Update()
	}

	aRes.AddResponseInfo("list", appList)
	aRes.AddResponseInfo("total", totalCount)
}

func removeAppVersionHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	if check_permission.CheckNoPermission(c, model_app.ApplicationPermissionDelete) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusBadRequest, "has no permission")
		return
	}
	para := new(model_app.AppVersion)
	err := c.BindJSON(para)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid: "+err.Error())
		return
	}
	c.Set(common.KeyContextPara, para.ToJson())
	if check_permission.CheckNoAppVersionManagerPermission(c, *para) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusBadRequest, "has no permission")
		return
	}
	err = para.Remove()
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusInternalServerError, "para invalid: "+err.Error())
		return
	}
	aRes.SetSuccess()
}
