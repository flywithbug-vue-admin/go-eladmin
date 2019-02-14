package menu_handler

import "go-eladmin/server/handler/handler_common"

var Routers = []handler_common.GinHandleFunc{
	{
		Handler:    addMenuHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "POST",
		Route:      "/menu",
	},
	{
		Handler:    getMenuHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/menu",
	},
	{
		Handler:    updateMenuHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "PUT",
		Route:      "/menu",
	},
	{
		Handler:    removeMenuHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "DELETE",
		Route:      "/menu",
	},
	{
		Handler:    getMenuListHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/menu/list",
	},
	{
		Handler:    getMenuTreeHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/menu/tree",
	},
	{
		Handler:    getMenuBuildHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/menu/build",
	},
}
