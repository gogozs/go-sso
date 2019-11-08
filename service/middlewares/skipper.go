package middlewares

import (
	"github.com/gin-gonic/gin"
	"strings"
)

type Skipper func(c *gin.Context, prefixes ...string) bool


func CreatePathSkipper() Skipper {
	return func(c *gin.Context, prefixes ...string) bool {
		//method := c.Request.Method
		path := c.Request.URL.Path
		for _, p := range prefixes {
			if strings.HasPrefix(path, p) {
				return true
			}
		}
		return false
	}
}