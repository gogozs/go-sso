package util

import (
	"github.com/gin-gonic/gin"
	"go-ops/pkg/settings"
	"github.com/Unknwon/com"
)


func GetPage(c *gin.Context) int {
	result := 0
	page, _ := com.StrTo(c.Query("page")).Int()
	if page >0 {
		result = (page -1 ) * settings.PageSize
	}
	return result
}
