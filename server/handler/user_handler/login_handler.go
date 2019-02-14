package user_handler

import (
	"go-eladmin/common"
	"go-eladmin/core/jwt"
	"go-eladmin/model"
	"go-eladmin/model/model_user"
	"go-eladmin/server/handler/handler_common"
	"go-eladmin/server/sync_map"
	"net/http"

	"github.com/flywithbug/log4go"
	"github.com/gin-gonic/gin"
)

func loginHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	user := new(model_user.User)
	err := c.BindJSON(user)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid"+err.Error())
		return
	}
	c.Set(common.KeyContextPara, user.ToJson())
	*user, err = model_user.LoginUser(user.Username, user.Password)
	if err != nil {
		log4go.Error(err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "username or password not right")
		return
	}
	if !user.Enabled {
		aRes.SetErrorInfo(http.StatusBadRequest, "账号已停用，请联系管理员")
		return
	}
	claims := jwt.NewCustomClaims(user.Id)
	token, err := jwt.GenerateToken(claims)
	if err != nil {
		log4go.Error(err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "token generate error"+err.Error())
		return
	}
	userAg := c.GetHeader(common.KeyUserAgent)
	_, err = model_user.UserLogin(user.Id, userAg, token, c.ClientIP())
	if err != nil {
		log4go.Error(err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "token generate error"+err.Error())
		return
	}
	sync_map.SetKeyValue(token)
	aRes.SetResponseDataInfo("token", token)
}

func registerHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	user := new(model_user.User)
	err := c.BindJSON(user)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid"+err.Error())
		return
	}
	c.Set(common.KeyContextPara, user.ToJson())

	if user.Username == "" {
		aRes.SetErrorInfo(http.StatusBadRequest, "username can not be nil")
		return
	}
	if user.Password == "" {
		aRes.SetErrorInfo(http.StatusBadRequest, "Password can not be nil")
		return
	}
	if user.Email == "" {
		aRes.SetErrorInfo(http.StatusBadRequest, "email can not be nil")
		return
	}
	err = user.Insert()
	if err != nil {
		aRes.SetErrorInfo(http.StatusBadRequest, err.Error())
		return
	}
	aRes.AddResponseInfo("user", user)
}

func logoutHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	token := common.UserToken(c)
	sync_map.RemoveKey(token)
	model_user.UpdateLoginStatus(token, model_user.StatusLogout)
	aRes.SetSuccessInfo(http.StatusOK, "success")
}
