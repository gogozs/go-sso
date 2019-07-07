package util

import (
	"github.com/gin-gonic/gin"
	"github.com/Unknwon/com"
	"go-weixin/config"
)


func GetPage(c *gin.Context) int {
	cf := config.GetConfig().Common
	result := 0
	page, _ := com.StrTo(c.Query("page")).Int()
	if page >0 {
		result = (page -1 ) * cf.PageSize
	}
	return result
}
