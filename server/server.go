package server

import (
	"fmt"
	"go-eladmin/server/handler"
	"go-eladmin/server/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/contrib/static"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartServer(port, staticPath, rPrefix, authPrefix string) {
	r := gin.New()
	r.Use(middleware.Logger(), gin.Recovery())
	r.Use(static.Serve("/", static.LocalFile(staticPath, true)))
	r.LoadHTMLGlob(staticPath + "/index.html")
	r.NoRoute(NoRoute)
	cors.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "UUID", "Content-Type", "Content-Length"},
		ExposeHeaders:    []string{},
		MaxAge:           12 * time.Hour,
		AllowCredentials: false,
	}))
	r.Use(middleware.CookieMiddleware())
	handler.RegisterRouters(r, rPrefix, authPrefix)
	//err := r.Run(port)
	s := &http.Server{
		Addr:           port,
		Handler:        r,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   50 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	if err != nil {
		panic(fmt.Errorf("server启动失败 %s", err.Error()))
	}
}
