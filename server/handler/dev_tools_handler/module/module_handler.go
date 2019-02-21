package module

import (
	"fmt"
	"go-eladmin/common"
	"go-eladmin/model"
	"go-eladmin/model/model_dev_tools/model_module"
	"go-eladmin/server/handler/check_permission"
	"go-eladmin/server/handler/handler_common"
	"net/http"

	"github.com/flywithbug/log4go"
	"github.com/gin-gonic/gin"
)

//添加业务线
func addModuleHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	if check_permission.CheckNoPermission(c, model_module.DataModulePermissionCreate) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusForbidden, "has no permission")
		return
	}
	para := new(model_module.Module)
	err := c.BindJSON(para)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid"+err.Error())
		return
	}
	c.Set(common.KeyContextPara, para.ToJson())
	if len(para.Name) == 0 {
		msg := fmt.Sprintf("name:%s not right", para.Name)
		log4go.Info(handler_common.RequestId(c) + msg)
		aRes.SetErrorInfo(http.StatusBadRequest, msg)
		return
	}
	id, err := para.Insert()
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusInternalServerError, "server invalid: "+err.Error())
		return
	}
	aRes.AddResponseInfo("id", id)
}
