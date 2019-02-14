package index_handler

import "go-eladmin/server/handler/handler_common"

var Routers = []handler_common.GinHandleFunc{
	{
		Handler:    indexHandler, //添加应用
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/index",
	},
}
