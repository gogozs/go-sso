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

	pathSkipMap = make(map[string]struct{})
)

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

func InitRouter(config *conf.Config) *gin.Engine {
	c := config.Common
	if c.Debug {
		crosOrigin()
	} else {
		gin.SetMode(gin.ReleaseMode) // 关闭gin debug
	}
	router.Use(
		middlewares.AuthMiddleware(
			[]middlewares.AuthType{middlewares.TokenAuth},
			middlewares.CreatePathSkipper(),
			pathSkipMap,
			"/api/public/",
		),
	)
	router.Use(
		middlewares.PermissionMiddleware(
			middlewares.CreatePathSkipper(),
			pathSkipMap,
			"/api/public/",
		),
	)
	AuthRouterInit()
	router.Use(middlewares.ErrEmailWriter())
	return router
}

func AuthRouterInit() {
	r := v1.NewAuthViewset()
	public := router.Group("/api/public/v1/auth/")
	{
		public.POST("/login/", r.ErrorHandler(r.Login))
		public.POST("/telephone-login/", r.ErrorHandler(r.TelephoneLogin))
		public.POST("/register/", r.ErrorHandler(r.Register))
		public.POST("/reset-password/", r.ErrorHandler(r.ResetPassword))
		public.POST("/send-sms-code/", r.ErrorHandler(r.SendSmsCode))
		public.POST("/send-email-code/", r.ErrorHandler(r.SendEmailCode))
		public.POST("/check-telephone-valid/", r.ErrorHandler(r.CheckTelephoneValid))
		public.POST("/check-telephone-exist/", r.ErrorHandler(r.CheckTelephoneExist))
	}
	authGroup := router.Group("/api/v1/auth/")
	{
		authGroup.POST("/change-password/", r.ErrorHandler(r.ChangePassword))
	}
}
