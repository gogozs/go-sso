package routers

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-sso/conf"
	"go-sso/internal/middlewares/auth"
	"go-sso/internal/middlewares/error_email"
	"go-sso/internal/middlewares/permissions"
	"go-sso/internal/middlewares/skipper"
	"go-sso/internal/repository/storage"
	v1 "go-sso/internal/service/v1"
	"go-sso/pkg/log"
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

func InitRouter(config *conf.Config, storage storage.Storage) *gin.Engine {
	c := config.Common
	if c.Debug {
		crosOrigin()
	} else {
		gin.SetMode(gin.ReleaseMode) // 关闭gin debug
	}
	router.Use(
		auth.AuthMiddleware(
			[]auth.AuthType{auth.TokenAuth},
			skipper.CreatePathSkipper(),
			pathSkipMap,
			"/api/public/",
		),
	)
	router.Use(
		permissions.PermissionMiddleware(
			skipper.CreatePathSkipper(),
			pathSkipMap,
			"/api/public/",
		),
	)
	AuthRouterInit(storage)
	router.Use(error_email.ErrEmailWriter())
	return router
}

func AuthRouterInit(storage storage.Storage) {
	r := v1.NewAuthViewset(storage)
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
