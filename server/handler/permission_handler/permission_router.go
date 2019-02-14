package permission_handler

import "go-eladmin/server/handler/handler_common"

var Routers = []handler_common.GinHandleFunc{
	{
		Handler:    addPermissionHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "POST",
		Route:      "/permission",
	},
	{
		Handler:    getPermissionHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/permission",
	},
	{
		Handler:    updatePermissionHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "PUT",
		Route:      "/permission",
	},
	{
		Handler:    removePermissionHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "DELETE",
		Route:      "/permission",
	},
	{
		Handler:    getPermissionTreeHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/permission/tree",
	},
	{
		Handler:    getPermissionListHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/permission/list",
	},
}
