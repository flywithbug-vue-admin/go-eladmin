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
		Handler:    addAttributeHandler, //添加模型属性
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "PUT",
		Route:      "/tools/model/attribute",
	},
	{
		Handler:    removeAttribuiteHandler, //删除模型属性
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "DELETE",
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
}
