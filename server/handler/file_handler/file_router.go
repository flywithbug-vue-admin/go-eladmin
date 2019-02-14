package file_handler

import "go-eladmin/server/handler/handler_common"

var Routers = []handler_common.GinHandleFunc{
	{
		Handler:    uploadImageHandler, //上传图片
		RouterType: handler_common.RouterTypeNormal,
		Method:     "POST",
		Route:      "/upload/image",
	},
	{
		Handler:    loadImageHandler, //加载图片
		RouterType: handler_common.RouterTypeNormal,
		Method:     "GET",
		Route:      "/image/:path/:filename",
	},
}
