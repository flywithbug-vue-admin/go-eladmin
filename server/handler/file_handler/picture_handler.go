package file_handler

import (
	"fmt"
	"go-eladmin/common"
	"go-eladmin/model"
	"go-eladmin/server/handler/handler_common"
	"net/http"

	"github.com/flywithbug/log4go"
	"github.com/gin-gonic/gin"
)

func uploadImageHandler(c *gin.Context) {
	aRes := model.NewResponse()
	defer func() {
		c.Set(common.KeyContextResponse, aRes)
		c.JSON(http.StatusOK, aRes)
	}()

	//gin将het/http包的FormFile函数封装到c.Request
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, fmt.Sprintf("get file err : %s", err.Error()))
		return
	}

	imgPath, err := saveImageFile(file, header)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		aRes.SetErrorInfo(http.StatusBadRequest, fmt.Sprintf("write file err : %s", err.Error()))
		return
	}
	aRes.SetResponseDataInfo("imagePath", imgPath)
}

func loadImageHandler(c *gin.Context) {
	path := c.Param("path")
	filename := c.Param("filename")
	//log4go.Info(handler_common.RequestId(c) + "loadImageHandler: %s %s", path, filename)
	if path == "" || filename == "" {
		return
	}
	size := c.Query("size")
	imgPath, err := loadImageFile(path, filename, size)
	if err != nil {
		log4go.Info(handler_common.RequestId(c) + err.Error())
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	http.ServeFile(c.Writer, c.Request, imgPath)

}
