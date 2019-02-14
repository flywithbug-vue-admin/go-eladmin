package dev_tools_handler

import (
	"encoding/json"
	"fmt"
	"go-eladmin/common"
	"go-eladmin/model"
	"go-eladmin/model/model_dev_tools/model_data_model"
	"go-eladmin/server/handler/check_permission"
	"go-eladmin/server/handler/handler_common"
	"net/http"
	"regexp"
	"strconv"

	"gopkg.in/mgo.v2/bson"

	"github.com/flywithbug/log4go"

	"github.com/gin-gonic/gin"
)

var (
	nameReg = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
)

type paraAttribute struct {
	ModelId   int64                      `json:"model_id"`
	Attribute model_data_model.Attribute `json:"attribute"`
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
	id, err := para.Insert()
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusInternalServerError, "server invalid: "+err.Error())
		return
	}
	aRes.AddResponseInfo("id", id)
}

func addAttributeHandler(c *gin.Context) {
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
	err = dm.AddAttribute(para.Attribute)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "invalid: "+err.Error())
		return
	}
	aRes.SetSuccess()
}

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

func removeAttribuiteHandler(c *gin.Context) {
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
	para := new(paraAttribute)
	err := c.BindJSON(para)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid"+err.Error())
		return
	}
	c.Set(common.KeyContextPara, para.ToJson())

	dm := model_data_model.DataModel{}
	dm.Id = para.ModelId
	err = dm.RemoveAttribute(para.Attribute)
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
