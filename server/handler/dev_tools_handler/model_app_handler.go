package dev_tools_handler

import (
	"fmt"
	"go-eladmin/common"
	"go-eladmin/model"
	"go-eladmin/model/model_app"
	"go-eladmin/model/model_app_data_model"
	"go-eladmin/model/model_dev_tools/model_data_model"
	"go-eladmin/server/handler/check_permission"
	"go-eladmin/server/handler/handler_common"
	"net/http"
	"strconv"

	"github.com/flywithbug/log4go"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

type paraModel struct {
	Model     model_app_data_model.AppDataModel `json:"model"`
	App       model_app.Application             `json:"app"`
	Option    string                            `json:"option"`
	PopStatus bool                              `json:"pop_status"`
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

func modifyAppModelVersionHandler(c *gin.Context) {
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
	para := new(model_app_data_model.AppDataModel)
	err := c.BindJSON(para)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid"+err.Error())
		return
	}
	c.Set(common.KeyContextPara, para.ToJson())
	if para.Id == 0 {
		msg := fmt.Sprintf("id is 0")
		log4go.Info(handler_common.RequestId(c) + msg)
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid:"+msg)
		return
	}

	dm := model_data_model.DataModel{}
	if para.ModelId > 0 {
		if !dm.Exist(bson.M{"_id": para.ModelId}) {
			msg := fmt.Sprintf("para invalid: modelId:%d not exist", para.ModelId)
			log4go.Info(handler_common.RequestId(c) + msg)
			aRes.SetErrorInfo(http.StatusBadRequest, msg)
			return
		}
	}
	if para.AppId > 0 {
		app := model_app.Application{}
		if !app.Exist(bson.M{"_id": para.AppId}) {
			msg := fmt.Sprintf("para invalid: appId:%d not exist", para.AppId)
			log4go.Info(handler_common.RequestId(c) + msg)
			aRes.SetErrorInfo(http.StatusBadRequest, msg)
			return
		}
	}
	oldDM, err := para.FindOne(bson.M{"_id": para.Id}, nil)
	if err != nil {
		if err != nil {
			log4go.Error(handler_common.RequestId(c) + err.Error())
			aRes.SetErrorInfo(http.StatusBadRequest, "model find error:"+err.Error())
			return
		}
	}
	if len(para.StartVersion) == 0 {
		para.StartVersion = oldDM.StartVersion
	}
	if len(para.EndVersion) == 0 {
		para.EndVersion = oldDM.EndVersion
	}
	err = para.Update()
	if err != nil {
		log4go.Error(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "model update error"+err.Error())
		return
	}
	aRes.SetSuccess()
}

/**
模型关联的App数据获取
*/
func getModelRelationAppListHandler(c *gin.Context) {
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
	list := make([]paraModel, 0)
	modelId, _ := strconv.ParseInt(c.Query("modelId"), 10, 64)
	if modelId == 0 {
		aRes.AddResponseInfo("list", list)
		return
	}
	query := bson.M{"model_id": modelId}
	appModel := model_app_data_model.AppDataModel{}
	results, err := appModel.FindAll(query, nil)
	if err != nil {
		log4go.Error(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "model find error"+err.Error())
		return
	}
	for _, item := range results {
		app := model_app.Application{}
		app.Id = item.AppId
		app, err := app.FindSimpleOne(bson.M{"_id": item.AppId}, nil)
		if err != nil {
			item.Remove()
			continue
		}
		parM := paraModel{}
		parM.Model = item
		parM.App = app
		list = append(list, parM)
	}
	aRes.AddResponseInfo("list", list)
}

/*
TODO 应用管理的模型数据 用于生产模型代码
*/
func appRelationModelListHandler(c *gin.Context) {
	//aRes := model.NewResponse()
	//defer func() {
	//	c.Set(common.KeyContextResponse, aRes)
	//	c.JSON(http.StatusOK, aRes)
	//}()
	//if check_permission.CheckNoPermission(c, model_data_model.DataModelPermissionSelect) {
	//	log4go.Info(handler_common.RequestId(c) + "has no permission")
	//	aRes.SetErrorInfo(http.StatusBadRequest, "has no permission")
	//	return
	//}
	//list := make([]paraModel, 0)
	//version := c.Query("version")
	//
	//appId, _ := strconv.ParseInt(c.Query("appId"), 10, 64)
	//if appId == 0 {
	//	aRes.AddResponseInfo("list", list)
	//	return
	//}
	//query := bson.M{"app_id": appId}
	////		if app.isExist(bson.M{"version": app.Version, "app_id": app.AppId, "_id": bson.M{"$ne": app.Id}}) {
	//if len(version) > 0 {
	//	vNum := common.TransformVersionToInt(version)
	//	query["start_v_num"] = bson.M{"$lte": vNum}
	//	query["end_v_num"] = bson.M{"$gte": vNum}
	//}
	//
	//appModel := model_app_data_model.AppDataModel{}
	//results, err := appModel.FindAll(query, nil)
	//if err != nil {
	//	log4go.Error(handler_common.RequestId(c) + err.Error())
	//	aRes.SetErrorInfo(http.StatusUnauthorized, "apps find error"+err.Error())
	//	return
	//}
	//
	//aRes.AddResponseInfo("list", results)
}

func updateAppModelRelationHandler(c *gin.Context) {

}

func removeAppModelRelationHandler(c *gin.Context) {
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
	para := new(model_app_data_model.AppDataModel)
	err := c.BindJSON(para)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid"+err.Error())
		return
	}
	c.Set(common.KeyContextPara, para.ToJson())
	if para.Id == 0 {
		msg := fmt.Sprintf("id is 0")
		log4go.Info(handler_common.RequestId(c) + msg)
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid:"+msg)
		return
	}
	err = para.Remove()
	if err != nil {
		log4go.Error(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "model delete error"+err.Error())
		return
	}
	aRes.SetSuccess()
}

func addAppModelRelationHandler(c *gin.Context) {
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
	para := new(model_app_data_model.AppDataModel)
	err := c.BindJSON(para)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid"+err.Error())
		return
	}
	c.Set(common.KeyContextPara, para.ToJson())
	if para.ModelId == 0 {
		msg := fmt.Sprintf("model_id is 0")
		log4go.Info(handler_common.RequestId(c) + msg)
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid:"+msg)
		return
	}
	if para.AppId == 0 {
		msg := fmt.Sprintf("app_id is 0")
		log4go.Info(handler_common.RequestId(c) + msg)
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid:"+msg)
		return
	}
	err = para.Insert()
	if err != nil {
		log4go.Error(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "model delete error"+err.Error())
		return
	}
	aRes.SetSuccess()

}
