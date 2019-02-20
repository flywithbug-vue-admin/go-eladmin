package user_handler

import (
	"fmt"
	"go-eladmin/common"
	"go-eladmin/email"
	"go-eladmin/model"
	"go-eladmin/model/model_user"
	"go-eladmin/server/handler/check_permission"
	"go-eladmin/server/handler/handler_common"
	"net/http"
	"strconv"
	"strings"

	"github.com/flywithbug/log4go"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

type UserRole struct {
	model_user.User
	Roles []string `json:"roles"`
}

type ParaUserEdit struct {
	model_user.User
	Enabled string `json:"enabled,omitempty" bson:"enabled,omitempty"`
}

func addUserHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()

	if check_permission.CheckNoPermission(c, model_user.UserPermissionCreate) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusForbidden, "has no permission")
		return
	}
	para := new(model_user.User)
	err := c.BindJSON(para)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid"+err.Error())
		return
	}
	c.Set(common.KeyContextPara, para.ToJson())
	err = para.Insert()
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid"+err.Error())
		return
	}
	aRes.SetSuccess()
	aRes.Msg = fmt.Sprintf("用户名：%s 密码：%s", para.Name, para.Password)
	err = email.SendMail("后台管理", "用户密码", aRes.Msg, para.Email)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + "邮件发送失败" + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "邮件发送失败"+err.Error())
		return
	}
}

func getUserInfoHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	var id int64
	ids := c.Query("id")
	id, _ = strconv.ParseInt(ids, 10, 64)
	if id == 0 {
		id = common.UserId(c)
	} else {
		if check_permission.CheckNoPermission(c, model_user.UserPermissionSelect) {
			log4go.Info(handler_common.RequestId(c) + "has no permission")
			aRes.SetErrorInfo(http.StatusForbidden, "has no permission")
			return
		}
	}
	user := model_user.User{}
	user.Id = id
	user, err := user.FindTreeOne()
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusUnauthorized, "user not found:"+err.Error())
		return
	}
	if !user.Enabled {
		log4go.Info(handler_common.RequestId(c) + "账号已停用")
		aRes.SetErrorInfo(http.StatusUnauthorized, "账号已停用，请联系管理员")
		return
	}
	roleUser := UserRole{}
	roleUser.User = user
	roleUser.Roles = user.RolesString
	if roleUser.Id == 10000 && len(roleUser.Roles) > 0 {
		roleUser.Roles = []string{"ADMIN"}
	}
	roleUser.RolesString = nil
	aRes.AddResponseInfo("user", roleUser)
}

func updateUserHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	if check_permission.CheckNoPermission(c, model_user.UserPermissionEdit) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusForbidden, "has no permission")
		return
	}
	para := new(model_user.User)
	err := c.BindJSON(para)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid"+err.Error())
		return
	}
	c.Set(common.KeyContextPara, para.ToJson())

	err = para.Update()
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "db update failed: "+err.Error())
		return
	}
	aRes.SetSuccessInfo(http.StatusOK, "success")
}

func deleteUserHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	if check_permission.CheckNoPermission(c, model_user.UserPermissionDelete) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusForbidden, "has no permission")
		return
	}
	para := new(model_user.User)
	err := c.BindJSON(para)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid"+err.Error())
		return
	}
	c.Set(common.KeyContextPara, para.ToJson())

	if common.UserId(c) == para.Id {
		log4go.Info(handler_common.RequestId(c) + "can not delete your self")
		aRes.SetErrorInfo(http.StatusForbidden, "can not delete your self")
		return
	}
	err = para.Remove()
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "db delete failed: "+err.Error())
		return
	}
	aRes.SetSuccessInfo(http.StatusOK, "success")
}

func getUserTreeListInfoHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	if check_permission.CheckNoPermission(c, model_user.UserPermissionSelect) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusForbidden, "has no permission")
		return
	}
	limit, _ := strconv.Atoi(c.Query("size"))
	page, _ := strconv.Atoi(c.Query("page"))
	sort := c.Query("sort")
	name := c.Query("username")
	email := c.Query("email")
	enabled := c.Query("enabled")

	if strings.EqualFold(sort, "-id") {
		sort = "-_id"
	} else if strings.EqualFold(sort, "+id") {
		sort = "+_id"
	} else if len(sort) == 0 {
		sort = "+_id"
	}
	if limit == 0 {
		limit = 10
	}
	if page != 0 {
		page--
	}
	userId := common.UserId(c)
	if userId <= 0 {
		aRes.SetErrorInfo(http.StatusUnauthorized, "user not found")
		return
	}
	query := bson.M{}
	if len(name) > 0 {
		query["username"] = bson.M{"$regex": name, "$options": "i"}
	}
	if len(email) > 0 {
		query["email"] = bson.M{"$regex": email, "$options": "i"}
	}

	if strings.EqualFold(enabled, "true") {
		query["enabled"] = true
	}
	if strings.EqualFold(enabled, "false") {
		query["enabled"] = false
	}

	var user = model_user.User{}
	totalCount, _ := user.TotalCount(query, nil)
	appList, err := user.FindPageTreeFilter(page, limit, query, bson.M{"password": 0}, sort)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusInternalServerError, err.Error())
		return
	}
	aRes.AddResponseInfo("list", appList)
	aRes.AddResponseInfo("total", totalCount)
}

func queryListHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	if check_permission.CheckNoPermission(c, model_user.UserPermissionSelect) {
		log4go.Info(handler_common.RequestId(c) + "has no permission")
		aRes.SetErrorInfo(http.StatusForbidden, "has no permission")
		return
	}
	limit, _ := strconv.Atoi(c.Query("size"))
	page, _ := strconv.Atoi(c.Query("page"))
	sort := c.Query("sort")
	name := c.Query("username")
	email := c.Query("email")
	enabled := c.Query("enabled")
	exc := c.Query("exc")

	if strings.EqualFold(sort, "-id") {
		sort = "-_id"
	} else if strings.EqualFold(sort, "+id") {
		sort = "+_id"
	} else if len(sort) == 0 {
		sort = "+_id"
	}
	if limit == 0 {
		limit = 10
	}
	if page != 0 {
		page--
	}
	userId := common.UserId(c)
	if userId <= 0 {
		aRes.SetErrorInfo(http.StatusUnauthorized, "user not found")
		return
	}
	query := bson.M{}
	if len(name) > 0 {
		query["username"] = bson.M{"$regex": name, "$options": "i"}

	}
	if len(email) > 0 {
		query["email"] = bson.M{"$regex": email, "$options": "i"}
	}

	if strings.EqualFold(enabled, "true") {
		query["enabled"] = true
	}
	if strings.EqualFold(enabled, "false") {
		query["enabled"] = false
	}
	if len(exc) > 0 {
		excepts := strings.Split(exc, ",")
		ids := make([]int64, len(excepts))
		for index, item := range excepts {
			ids[index], _ = strconv.ParseInt(item, 10, 64)
		}
		query["_id"] = bson.M{"$nin": ids}
	}
	var user = model_user.User{}
	totalCount, _ := user.TotalCount(query, nil)
	appList, err := user.FindPageFilter(page, limit, query, bson.M{"password": 0}, sort)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusInternalServerError, err.Error())
		return
	}
	aRes.AddResponseInfo("list", appList)
	aRes.AddResponseInfo("total", totalCount)
}

func updateAvatar(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()
	para := new(model_user.User)
	err := c.BindJSON(para)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "para invalid"+err.Error())
		return
	}
	para.Id = common.UserId(c)
	c.Set(common.KeyContextPara, para.ToJson())
	err = para.UpdateAvatar()
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, "db update error:"+err.Error())
		return
	}
	aRes.SetSuccess()
}
