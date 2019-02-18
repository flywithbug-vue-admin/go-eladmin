package dev_tools_handler

import "go-eladmin/server/handler/handler_common"

var Routers = []handler_common.GinHandleFunc{
	{
		Handler:    addDataModelHandler, //添加模型
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "POST",
		Route:      "/tools/model",
	},
	{
		Handler:    updateDataModelHandler, //模型修改
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "PUT",
		Route:      "/tools/model",
	},
	{
		Handler:    updateApplicationRelationHandler, //添加应用和模型关联
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "PUT",
		Route:      "/tools/model/apps",
	},
	{
		Handler:    modifyAttributeHandler, //添加模型属性
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "PUT",
		Route:      "/tools/model/attribute",
	},
	{
		Handler:    removeDataModelHandler, //模型删除
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "DELETE",
		Route:      "/tools/model",
	},
	{
		Handler:    getDataModelHandler, //获取模型数据
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/tools/model",
	},
	{
		Handler:    listHandler, //获取模型数据
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/tools/model/list",
	},
	{
		Handler:    getModelRelationAppListHandler, //获取模型数据
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/tools/model_apps",
	},
	{
		Handler:    modeifyAppModelHandler, //获取模型数据
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "PUT",
		Route:      "/tools/model_apps",
	},
}
