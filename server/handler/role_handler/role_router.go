package role_handler

import "go-eladmin/server/handler/handler_common"

var Routers = []handler_common.GinHandleFunc{
	{
		Handler:    addRoleHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "POST",
		Route:      "/role",
	},
	{
		Handler:    getRoleHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/role",
	},
	{
		Handler:    updateRoleHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "PUT",
		Route:      "/role",
	},
	{
		Handler:    removeRoleHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "DELETE",
		Route:      "/role",
	},
	{
		Handler:    getRoleListHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/role/list",
	},
	{
		Handler:    getRoleTreeHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/role/tree",
	},
}
