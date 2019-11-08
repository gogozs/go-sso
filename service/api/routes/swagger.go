// +build doc

package routes

import (
	"fmt"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go-qiuplus/conf"
	_ "go-qiuplus/docs"
)



var swag = initSwag() // 注意应先于中间件执行, 控制初始化顺序

func initSwag() error {
	port := conf.GetConfig().Common.HttpPort
	jsonPath := fmt.Sprintf("http://localhost:%d/swagger/doc.json", port)
	url := ginSwagger.URL(jsonPath) // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return nil
}