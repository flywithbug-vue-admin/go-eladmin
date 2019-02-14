package system_handler

import "go-eladmin/server/handler/handler_common"

var Routers = []handler_common.GinHandleFunc{
	{
		Handler:    systemHandle,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "get",
		Route:      "/system",
	},
}
