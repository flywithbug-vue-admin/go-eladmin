package verify_handler

import (
	"encoding/json"
	"go-eladmin/common"
	"go-eladmin/email"
	"go-eladmin/model"
	"go-eladmin/model/model_verify"
	"go-eladmin/server/handler/handler_common"
	"net/http"
	"strconv"

	"gopkg.in/mgo.v2/bson"

	"github.com/flywithbug/log4go"
	"github.com/gin-gonic/gin"
)

type mailVerifyPara struct {
	Mail string `json:"mail"`
}

func (v mailVerifyPara) ToJson() string {
	js, _ := json.Marshal(v)
	return string(js)
}

func sendVerifyMailHanlder(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	para := new(mailVerifyPara)
	err := c.BindJSON(para)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "Param invalid"+err.Error())
		return
	}
	c.Set(common.KeyContextPara, para.ToJson())
	if !email.MailVerify(para.Mail) {
		aRes.SetErrorInfo(http.StatusBadRequest, "mail invalid")
		return
	}
	vCode, err := model_verify.GeneralVerifyData(para.Mail)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusInternalServerError, "invalid"+err.Error())
		return
	}
	err = email.SendVerifyCode("", vCode, para.Mail)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusInternalServerError, err.Error())
		return
	}
	aRes.SetSuccess()
}

func getVerifyListHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	size, _ := strconv.Atoi(c.Query("size"))
	page, _ := strconv.Atoi(c.Query("page"))
	value := c.Query("value")
	code := c.Query("code")
	scenes := c.Query("scenes")
	if size == 0 {
		size = 10
	}
	if page != 0 {
		page--
	}
	query := bson.M{}
	if len(value) > 0 {
		query["value"] = value
	}

	if len(code) > 0 {
		query["code"] = code
	}

	if len(scenes) > 0 {
		query["scenes"] = scenes
	}
	var v = model_verify.VerificationCode{}
	totalCount, _ := v.TotalCount(query, nil)
	results, err := v.FindPageFilter(page, size, query, nil, "-_id")
	if err != nil {
		log4go.Error(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusInternalServerError, "list find error"+err.Error())
		return
	}
	aRes.AddResponseInfo("list", results)
	aRes.AddResponseInfo("total", totalCount)
}
