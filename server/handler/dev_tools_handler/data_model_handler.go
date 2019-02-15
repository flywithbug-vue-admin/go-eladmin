package dev_tools_handler

import (
	"encoding/json"
	"fmt"
	"go-eladmin/common"
	"go-eladmin/model"
	"go-eladmin/model/model_app"
	"go-eladmin/model/model_app_data_model"
	"go-eladmin/model/model_dev_tools/model_data_model"
	"go-eladmin/server/handler/check_permission"
	"go-eladmin/server/handler/handler_common"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/flywithbug/log4go"

	"github.com/gin-gonic/gin"
)

var (
	nameReg = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
)

type paraAttribute struct {
	ModelId        int64                        `json:"model_id,omitempty"`
	Attributes     []model_data_model.Attribute `json:"attributes,omitempty"`      //批量修改或增加的属性
	DropAttributes []model_data_model.Attribute `json:"drop_attributes,omitempty"` //批量删除的属性
}

func (u paraAttribute) ToJson() string {
	js, _ := json.Marshal(u)
	return string(js)
}

func addDataModelHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	if check_permission.CheckNoPermission(c, model_data_model.DataModelPermissionCreate) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusBadRequest, "has no permission")
		return
	}
	para := new(model_data_model.DataModel)
	err := c.BindJSON(para)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid"+err.Error())
		return
	}
	c.Set(common.KeyContextPara, para.ToJson())
	match := nameReg.FindAllString(para.Name, -1)
	if len(match) == 0 {
		msg := fmt.Sprintf("name:%s not right", para.Name)
		log4go.Info(handler_common.RequestId(c) + msg)
		aRes.SetErrorInfo(http.StatusBadRequest, msg)
		return
	}
	userId := common.UserId(c)
	para.Owner.Id = userId
	id, err := para.Insert()
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusInternalServerError, "server invalid: "+err.Error())
		return
	}
	aRes.AddResponseInfo("id", id)
}

func modifyAttributeHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	if check_permission.CheckNoPermission(c, model_data_model.DataModelPermissionCreate) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusBadRequest, "has no permission")
		return
	}
	para := new(paraAttribute)
	err := c.BindJSON(para)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid"+err.Error())
		return
	}
	c.Set(common.KeyContextPara, para.ToJson())
	if para.ModelId == 0 {
		log4go.Info(handler_common.RequestId(c) + "id is 0")
		aRes.SetErrorInfo(http.StatusBadRequest, "id is 0")
		return
	}
	dm := model_data_model.DataModel{}
	dm.Id = para.ModelId

	if len(para.DropAttributes) > 0 {
		dm.RemoveAttributes(para.DropAttributes)
	}
	if len(para.Attributes) > 0 {
		dm.RemoveAttributes(para.Attributes)
		err = dm.AddAttributes(para.Attributes)
		if err != nil {
			log4go.Info(handler_common.RequestId(c) + err.Error())
			aRes.SetErrorInfo(http.StatusBadRequest, "invalid: "+err.Error())
			return
		}
	}
	aRes.SetSuccess()
}

func updateApplicationRelationHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	if check_permission.CheckNoPermission(c, model_data_model.DataModelPermissionEdit) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusBadRequest, "has no permission")
		return
	}
	para := new(model_data_model.DataModel)
	err := c.BindJSON(para)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid"+err.Error())
		return
	}
	c.Set(common.KeyContextPara, para.ToJson())
	if para.Id == 0 {
		log4go.Info(handler_common.RequestId(c) + "id is 0")
		aRes.SetErrorInfo(http.StatusBadRequest, "id is 0")
		return
	}
	err = para.UpdateAppRelation()
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "invalid: "+err.Error())
		return
	}
	aRes.SetSuccess()
}

//更新alias或者name
func updateDataModelHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	if check_permission.CheckNoPermission(c, model_data_model.DataModelPermissionEdit) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusBadRequest, "has no permission")
		return
	}
	para := new(model_data_model.DataModel)
	err := c.BindJSON(para)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid"+err.Error())
		return
	}
	c.Set(common.KeyContextPara, para.ToJson())
	if para.Id == 0 {
		log4go.Info(handler_common.RequestId(c) + "id is 0")
		aRes.SetErrorInfo(http.StatusBadRequest, "id is 0")
		return
	}
	err = para.Update()
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "invalid: "+err.Error())
		return
	}
	aRes.SetSuccess()
}

func removeDataModelHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	if check_permission.CheckNoPermission(c, model_data_model.DataModelPermissionDelete) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusBadRequest, "has no permission")
		return
	}
	para := new(model_data_model.DataModel)
	err := c.BindJSON(para)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid"+err.Error())
		return
	}
	c.Set(common.KeyContextPara, para.ToJson())
	if para.Id == 0 {
		log4go.Info(handler_common.RequestId(c) + "id is 0")
		aRes.SetErrorInfo(http.StatusBadRequest, "id is 0")
		return
	}
	err = para.Remove()
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "invalid: "+err.Error())
		return
	}
	aRes.SetSuccess()
}

func getDataModelHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	if check_permission.CheckNoPermission(c, model_data_model.DataModelPermissionSelect) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusBadRequest, "has no permission")
		return
	}
	ids := c.Query("id")
	if len(ids) == 0 {
		log4go.Info(handler_common.RequestId(c) + "para invalid")
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid")
		return
	}
	id, _ := strconv.ParseInt(ids, 10, 64)
	para := model_data_model.DataModel{}
	para.Id = id
	para, err := para.FindOne(bson.M{"_id": id}, nil)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "invalid: "+err.Error())
		return
	}
	aRes.AddResponseInfo("model", para)
}

func listHandler(c *gin.Context) {
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
	limit, _ := strconv.Atoi(c.Query("size"))
	page, _ := strconv.Atoi(c.Query("page"))
	sort := c.Query("sort")
	name := c.Query("name")
	appId, _ := strconv.ParseInt(c.Query("appId"), 10, 64)
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
	if appId > 0 {
		am := model_app_data_model.AppDataModel{}
		totalCount, _ := am.TotalCount(query, nil)

		result, err := am.FindAll(bson.M{"app_id": appId}, nil)
		if err != nil {
			log4go.Error(handler_common.RequestId(c) + err.Error())
			aRes.SetErrorInfo(http.StatusUnauthorized, "apps find error"+err.Error())
			return
		}
		listA := make([]model_data_model.DataModel, 0)
		for _, item := range result {
			dm := model_data_model.DataModel{}
			dm, err := dm.FindOne(bson.M{"_id": item.ModelId}, nil)
			if err != nil {
				am.RemoveModelId(item.ModelId)
				continue
			}
			listA = append(listA, dm)
		}
		aRes.AddResponseInfo("list", listA)
		aRes.AddResponseInfo("total", totalCount)
		return
	}
	var dm = model_data_model.DataModel{}
	totalCount, _ := dm.TotalCount(query, nil)
	list, err := dm.FindPageFilter(page, limit, query, nil, sort)
	if err != nil {
		log4go.Error(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusUnauthorized, "apps find error"+err.Error())
		return
	}
	aRes.AddResponseInfo("list", list)
	aRes.AddResponseInfo("total", totalCount)
}
