package user_handler

import "go-eladmin/server/handler/handler_common"

var Routers = []handler_common.GinHandleFunc{
	{
		Handler:    registerHandler,
		RouterType: handler_common.RouterTypeNormal,
		Method:     "POST",
		Route:      "/register",
	},
	{
		Handler:    loginHandler,
		RouterType: handler_common.RouterTypeNormal,
		Method:     "POST",
		Route:      "/login",
	},
	{
		Handler:    logoutHandler,
		RouterType: handler_common.RouterTypeNeedAuth,
		Route:      "/logout",
		Method:     "POST",
	},
	{
		Handler:    getUserInfoHandler, //获取userId对应的用户信息
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/user/info",
	},
	{
		Handler:    updateUserHandler, //更新当前用户信息
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "PUT",
		Route:      "/user",
	},
	{
		Handler:    addUserHandler, //添加用户当前用户信息
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "POST",
		Route:      "/user",
	},
	{
		Handler:    deleteUserHandler, //删除当前用户信息
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "DELETE",
		Route:      "/user",
	},
	{
		Handler:    getUserTreeListInfoHandler, //获取所有用户
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/user/tree",
	},
	{
		Handler:    queryListHandler, //获取所有用户
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/user/list",
	},
	{
		Handler:    validPasswordHandler, //验证密码
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "GET",
		Route:      "/user/password/valid",
	},
	{
		Handler:    updatePasswordHandler, //修改密码
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "PUT",
		Route:      "/user/password",
	},
	{
		Handler:    updateMailHandler, //修改邮箱
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "PUT",
		Route:      "/user/mail",
	},
	{
		Handler:    updateAvatar, //修改头像
		RouterType: handler_common.RouterTypeNeedAuth,
		Method:     "PUT",
		Route:      "/user/avatar",
	},
}
