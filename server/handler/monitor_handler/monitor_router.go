package monitor_handler

import "go-eladmin/server/handler/handler_common"

var Routers = []handler_common.GinHandleFunc{
	{
		Handler:    visitHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/visit/info",
	},
	{
		Handler:    visitCountHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "POST",
		Route:      "/visit",
	},
	{
		Handler:    chartListHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/visit/chart",
	},
}
