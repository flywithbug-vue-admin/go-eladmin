package verify_handler

import "go-eladmin/server/handler/handler_common"

var Routers = []handler_common.GinHandleFunc{
	{
		Handler:    sendVerifyMailHanlder,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "POST",
		Route:      "/mail/verify",
	},
	{
		Handler:    getVerifyListHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/verify/list",
	},
}
