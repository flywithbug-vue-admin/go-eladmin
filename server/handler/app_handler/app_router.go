package app_handler

import "go-eladmin/server/handler/handler_common"

var Routers = []handler_common.GinHandleFunc{
	{
		Handler:    addHandler, //添加应用
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "POST",
		Route:      "/app",
	},
	{
		Handler:    editHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "PUT",
		Route:      "/app",
	},
	{
		Handler:    delHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "DELETE",
		Route:      "/app",
	},
	{
		Handler:    listHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/app/list",
	},
	{
		Handler:    simpleListHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/app/list/simple",
	},

	//app Version Handler

	{
		Handler:    addAppVersionHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "POST",
		Route:      "/app/version",
	},
	{
		Handler:    removeAppVersionHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "DELETE",
		Route:      "/app/version",
	},
	{
		Handler:    getAppVersionListHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/app/version/list",
	},
	{
		Handler:    updateAppVersionHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "PUT",
		Route:      "/app/version",
	},
	{
		Handler:    queryAppVersion,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/app/version",
	},
}
