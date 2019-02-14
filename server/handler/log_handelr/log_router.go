package log_handelr

import "go-eladmin/server/handler/handler_common"

var Routers = []handler_common.GinHandleFunc{
	{
		Handler:    getLogListHandler, //添加应用
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/log/list",
	},
}
