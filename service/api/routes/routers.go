package routes

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-sso/conf"
	"go-sso/pkg/log"
	v1 "go-sso/service/api/v1"
	"go-sso/service/middlewares"
	"net/http"
	"time"
)


var (
	router = newRouter()
	swagHandler gin.HandlerFunc
)

const (
	path   = "/"
	pkPath = "/:id/"
)

func init()  {
	initRouter()
}

func GetRouter() *gin.Engine {
	return router
}

func StartServer() *http.Server {
	c := conf.GetConfig().Common
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", c.HttpPort),
		Handler:        router,
		ReadTimeout:    c.ReadTimeout * time.Second,
		WriteTimeout:   c.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Info(fmt.Sprintf("server start at 0.0.0.0:%d", c.HttpPort))
	return server
}

// 跨域
func crosOrigin() {
	c := conf.GetConfig().Cors
	config := cors.DefaultConfig()
	if len(c) > 0 {
		config.AllowOrigins = c
		router.Use(cors.New(config))
	}
}

func newRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func initRouter() *gin.Engine {
	c := conf.GetConfig().Common
	if c.Debug {
		crosOrigin()
	} else {
		gin.SetMode(gin.ReleaseMode) // 关闭gin debug
	}
	router.Use(
		middlewares.AuthMiddleware(
			[]string{"jwt"},
			middlewares.CreatePathSkipper(),
			"/api/public/",
		),
	)
	router.Use(
		middlewares.PermissionMiddleware(
			middlewares.CreatePathSkipper(),
			"/api/public/",
		),
	)
	AuthRouterInit()
	WxRouterInit()
	router.Use(middlewares.ErrEmailWriter())
	return router
}

func AuthRouterInit() {
	r := v1.GetAuthVS()
	public := router.Group("/api/public/v1/auth/")
	{
		public.POST("/login/", r.ErrorHandler(r.Login))
		public.POST("/register/", r.ErrorHandler(r.Register))
	}
}

func WxRouterInit() {
	r := v1.GetWxVS()
	wx := router.Group("/api/public/v1/wx/")
	{
		wx.POST("/login/", r.ErrorHandler(r.Login))
		wx.POST("/get-userinfo/", r.ErrorHandler(r.GetUserInfo))
	}
}