package user_handler

import (
	"encoding/json"
	"go-eladmin/common"
	"go-eladmin/model"
	"go-eladmin/model/model_user"
	"go-eladmin/model/model_verify"
	"go-eladmin/server/handler/handler_common"
	"net/http"
	"strings"

	"github.com/flywithbug/log4go"
	"github.com/gin-gonic/gin"
)

type ParaUser struct {
	Password    string `json:"password"`
	OldPassword string `json:"old_password"`
	Code        string `json:"code"`
	Mail        string `json:"mail"`
}

func (u ParaUser) ToJson() string {
	js, _ := json.Marshal(u)
	return string(js)
}

func validPasswordHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	password := c.Query("password")
	if len(password) == 0 {
		log4go.Info(handler_common.RequestId(c) + "password need")
		aRes.SetErrorInfo(http.StatusBadRequest, "password need")
		return
	}
	user := new(model_user.User)
	user.Password = password
	user.Id = common.UserId(c)
	if !user.CheckPassword() {
		log4go.Info(handler_common.RequestId(c) + "password not right")
		aRes.SetErrorInfo(http.StatusBadRequest, "password not right")
		return
	}
	aRes.SetSuccess()
}

func updatePasswordHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	para := new(ParaUser)
	err := c.BindJSON(para)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid"+err.Error())
		return
	}
	c.Set(common.KeyContextPara, para.ToJson())

	if strings.EqualFold(para.Password, para.OldPassword) {
		aRes.SetErrorInfo(http.StatusBadRequest, "password not changed")
		return
	}
	if len(para.Password) == 0 || len(para.OldPassword) == 0 {
		log4go.Info(handler_common.RequestId(c) + "para not right")
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid")
		return
	}
	user := new(model_user.User)
	user.Password = para.OldPassword
	user.Id = common.UserId(c)
	if !user.CheckPassword() {
		log4go.Info(handler_common.RequestId(c) + "password not right")
		aRes.SetErrorInfo(http.StatusBadRequest, "password not right")
		return
	}
	user.Password = para.Password
	err = user.UpdatePassword()
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusInternalServerError, "system error: "+err.Error())
		return
	}
	aRes.SetSuccess()
}

func updateMailHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	para := new(ParaUser)
	err := c.BindJSON(para)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid"+err.Error())
		return
	}
	c.Set(common.KeyContextPara, para.ToJson())

	if len(para.Mail) == 0 || len(para.Password) == 0 || len(para.Code) == 0 {
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid")
		return
	}

	if !model_verify.CheckVerify(para.Mail, para.Code) {
		aRes.SetErrorInfo(http.StatusBadRequest, "code not right")
		return
	}
	var user model_user.User
	user.Id = common.UserId(c)
	user.Password = para.Password
	if !user.CheckPassword() {
		aRes.SetErrorInfo(http.StatusBadRequest, "password not right")
		return
	}
	user.Email = para.Mail
	user.UpdateMail()
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusInternalServerError, "system error: "+err.Error())
		return
	}
	aRes.SetSuccess()
}
