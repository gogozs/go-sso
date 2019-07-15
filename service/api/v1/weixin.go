package v1

import (
	"github.com/gin-gonic/gin"
	"go-weixin/pkg/log"
	"go-weixin/sdk"
)

// token 认证
func ViewWx(c *gin.Context) {
	s := sdk.TokenSignature{}
	err := c.BindQuery(&s)
	if err != nil {
		log.Error(err)
		c.JSON(200, false)
	} else if s.Confirm() {
		c.JSON(200, true)
	} else {
		c.JSON(200, false)
	}
}
