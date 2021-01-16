package skipper

import (
	"github.com/gin-gonic/gin"
	"strings"
)

type Skipper func(c *gin.Context, m map[string]struct{}, prefixes ...string) bool

func CreatePathSkipper() Skipper {
	return func(c *gin.Context, m map[string]struct{}, prefixes ...string) bool {
		//method := c.Request.Method
		path := c.Request.URL.Path
		if _, ok := m[path]; ok {
			return true
		}
		for _, p := range prefixes {
			if strings.HasPrefix(path, p) {
				return true
			}
		}
		return false
	}
}
