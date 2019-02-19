package dev_tools_handler

import (
	"encoding/json"
	"fmt"
	"go-eladmin/common"
	"go-eladmin/model"
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
	Id             int64                        `json:"id,omitempty"`
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
	if para.ParentId > 0 && !para.Exist(bson.M{"_id": para.ParentId}) {
		msg := fmt.Sprintf("parent classId:%d not exist", para.ParentId)
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
	if para.Id == 0 {
		log4go.Info(handler_common.RequestId(c) + "id is 0")
		aRes.SetErrorInfo(http.StatusBadRequest, "id is 0")
		return
	}
	dm := model_data_model.DataModel{}
	dm.Id = para.Id
	if len(para.DropAttributes) > 0 {
		dm.RemoveAttributes(para.DropAttributes)
	}
	if len(para.Attributes) > 0 {
		err = dm.AddAttributes(para.Attributes)
		if err != nil {
			log4go.Info(handler_common.RequestId(c) + err.Error())
			aRes.SetErrorInfo(http.StatusBadRequest, "invalid: "+err.Error())
			return
		}
	}
	aRes.SetSuccess()
}

//更新model info alias或者name
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
	if para.ParentId > 0 && !para.Exist(bson.M{"_id": para.ParentId}) {
		log4go.Info(handler_common.RequestId(c) + "parent model not exist")
		aRes.SetErrorInfo(http.StatusBadRequest, "parent model not exist")
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
	if check_permission.CheckNoPermission(c, model_data_model.DataModelPermissionCreate) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusBadRequest, "has no permission")
		return
	}
	limit, _ := strconv.Atoi(c.Query("size"))
	page, _ := strconv.Atoi(c.Query("page"))
	sort := c.Query("sort")
	name := c.Query("name")
	appId, _ := strconv.ParseInt(c.Query("appId"), 10, 64)
	exc := c.Query("exc")

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

	selector := bson.M{"_id": 1, "name": 1, "alias": 1, "desc": 1, "create_time": 1}
	if len(exc) > 0 {
		excepts := strings.Split(exc, ",")
		ids := make([]int64, len(excepts))
		for index, item := range excepts {
			ids[index], _ = strconv.ParseInt(item, 10, 64)
		}
		query["_id"] = bson.M{"$nin": ids}
	}

	if appId > 0 {
		am := model_app_data_model.AppDataModel{}
		totalCount, _ := am.TotalCount(bson.M{"app_id": appId}, selector)
		result, err := am.FindPageFilter(page, limit, bson.M{"app_id": appId}, selector, sort)
		if err != nil {
			log4go.Error(handler_common.RequestId(c) + err.Error())
			aRes.SetErrorInfo(http.StatusUnauthorized, "apps find error"+err.Error())
			return
		}
		listA := make([]model_data_model.DataModel, 0)
		for _, item := range result {
			dm := model_data_model.DataModel{}
			query["_id"] = item.ModelId
			dm, err := dm.FindSimpleOne(query, selector)
			if err != nil {
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
	list, err := dm.FindPageFilter(page, limit, query, selector, sort)
	if err != nil {
		log4go.Error(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusUnauthorized, "apps find error"+err.Error())
		return
	}
	aRes.AddResponseInfo("list", list)
	aRes.AddResponseInfo("total", totalCount)
}
