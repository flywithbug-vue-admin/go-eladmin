package module

import "go-eladmin/server/handler/handler_common"

var Routers = []handler_common.GinHandleFunc{
	{
		Handler:    addModuleHandler, //添加模块
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "POST",
		Route:      "/tools/module",
	},
}
