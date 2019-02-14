package menu_handler

import (
	"go-eladmin/common"
	"go-eladmin/model"
	"go-eladmin/model/model_menu"
	"go-eladmin/model/model_role"
	"go-eladmin/model/model_user"
	"go-eladmin/server/handler/check_permission"
	"go-eladmin/server/handler/handler_common"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/flywithbug/log4go"
	"github.com/gin-gonic/gin"
)

func addMenuHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	if check_permission.CheckNoPermission(c, model_menu.MenuPermissionSelect) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusBadRequest, "has no permission")
		return
	}

	para := new(model_menu.Menu)
	err := c.BindJSON(para)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid"+err.Error())
		return
	}
	c.Set(common.KeyContextPara, para.ToJson())

	if !para.IFrame && strings.HasPrefix(para.Path, "http") {
		aRes.SetErrorInfo(http.StatusBadRequest, "外链必须以http或者https开头")
		return
	}

	_, err = para.Insert()
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusInternalServerError, "server invalid"+err.Error())
		return
	}
	aRes.SetSuccess()
}

func getMenuHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	if check_permission.CheckNoPermission(c, model_menu.MenuPermissionSelect) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusBadRequest, "has no permission")
		return
	}

	ids := c.Query("id")
	var menu = model_menu.Menu{}
	id, _ := strconv.Atoi(ids)
	menu.Id = int64(id)
	menu, err := menu.FindOneTree()
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "invalid: "+err.Error())
		return
	}
	aRes.AddResponseInfo("menu", menu)
}

func updateMenuHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	if check_permission.CheckNoPermission(c, model_menu.MenuPermissionEdit) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusBadRequest, "has no permission")
		return
	}
	para := new(model_menu.Menu)
	err := c.BindJSON(para)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid"+err.Error())
		return
	}
	c.Set(common.KeyContextPara, para.ToJson())

	if !para.IFrame && !strings.HasPrefix(para.Path, "http") {
		aRes.SetErrorInfo(http.StatusBadRequest, "外链必须以http或者https开头")
		return
	}
	err = para.Update()
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "invalid: "+err.Error())
		return
	}
}

func removeMenuHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	if check_permission.CheckNoPermission(c, model_menu.MenuPermissionDelete) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusBadRequest, "has no permission")
		return
	}
	//need id
	para := new(model_menu.Menu)
	err := c.BindJSON(para)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid"+err.Error())
		return
	}
	c.Set(common.KeyContextPara, para.ToJson())

	err = para.Remove()
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "invalid: "+err.Error())
		return
	}
}

func getMenuListHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()

	//if check_permission.CheckNoPermission(c, model_menu.MenuPermissionSelect) {
	//	log4go.Info(handler_common.RequestId(c) + "has no permission")
	//	aRes.SetErrorInfo(http.StatusBadRequest, "has no permission")
	//	return
	//}
	var role = model_menu.Menu{}
	limit, _ := strconv.Atoi(c.Query("size"))
	page, _ := strconv.Atoi(c.Query("page"))
	sort := "+sort"
	name := c.Query("name")

	if page != 0 {
		page--
	}
	query := bson.M{"pid": 0}
	if len(name) > 0 {
		query["name"] = bson.M{"$regex": name, "$options": "i"}
	}
	totalCount, _ := role.TotalCount(query, nil)
	list, err := role.FindPageListFilter(page, limit, query, nil, sort)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusUnauthorized, "apps find error"+err.Error())
		return
	}
	aRes.AddResponseInfo("list", list)
	aRes.AddResponseInfo("total", totalCount)
}

func getMenuTreeHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	//if check_permission.CheckNoPermission(c, model_menu.MenuPermissionSelect) {
	//	log4go.Info(handler_common.RequestId(c) + "has no permission")
	//	aRes.SetErrorInfo(http.StatusBadRequest, "has no permission")
	//	return
	//}
	var role = model_menu.Menu{}
	query := bson.M{"pid": 0}
	selector := bson.M{"_id": 1, "name": 1}
	list, err := role.FindPageTreeFilter(0, 0, query, selector)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusUnauthorized, "app version list find error"+err.Error())
		return
	}
	aRes.AddResponseInfo("list", list)
}

func getMenuBuildHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	//if check_permission.CheckNoPermission(c, model_menu.MenuPermissionSelect) {
	//	log4go.Info(handler_common.RequestId(c) + "has no permission")
	//	aRes.SetErrorInfo(http.StatusBadRequest, "has no permission")
	//	return
	//}
	sort := "+sort"
	var role = model_menu.Menu{}
	query := bson.M{"pid": 0}
	roles := getUserRoles(c)
	//js, _ := json.Marshal(roles)
	//log4go.Info(handler_common.RequestId(c) + string(js))

	list, err := role.FindPageBuildFilter(roles, 0, 0, query, nil, sort)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusUnauthorized, "app version list find error"+err.Error())
		return
	}
	aRes.AddResponseInfo("list", list)
}

func getUserRoles(c *gin.Context) []model_role.Role {
	id := common.UserId(c)
	user := model_user.User{}
	user.Id = id
	user, _ = user.FindTreeOne()
	return user.Roles
}
