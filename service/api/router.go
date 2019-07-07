package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-weixin/config"
	"go-weixin/pkg/middlewares/erroremail"
	"net/http"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(erroremail.ErrEmailWriter())
	return r
}

func StartServer() {
	router := InitRouter()
	c := config.GetConfig().Common
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", c.HTTP_PORT),
		Handler:        router,
		ReadTimeout:    c.READ_TIMEOUT,
		WriteTimeout:   c.WRITE_TIMEOUT,
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}
