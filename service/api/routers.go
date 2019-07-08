package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-weixin/config"
	"net/http"
)

var router *gin.Engine

func GetRouter() *gin.Engine {
	if router == nil {
		return InitRouter()
	}
	return router
}

func InitRouter() *gin.Engine {
	router = gin.Default()
	//router.Use(erroremail.ErrEmailWriter())
	AuthRouterInit()
	return router
}

func StartServer() *http.Server {
	router := InitRouter()
	c := config.GetConfig().Common
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", c.HttpPort),
		Handler:        router,
		ReadTimeout:    c.ReadTimeout,
		WriteTimeout:   c.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	return server
}