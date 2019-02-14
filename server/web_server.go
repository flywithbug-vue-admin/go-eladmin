package server

import (
	"fmt"
	"go-eladmin/server/middleware"
	"net/http"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func StartWeb(port, staticPath string) {
	r := gin.New()
	r.Use(middleware.Logger(), gin.Recovery())
	r.NoRoute(NoRoute)
	r.Use(static.Serve("/", static.LocalFile(staticPath, true)))
	r.LoadHTMLGlob(staticPath + "/index.html")
	err := r.Run(port)
	if err != nil {
		panic(fmt.Errorf("server启动失败 %s", err.Error()))
	}
}

func NoRoute(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", "")
}
