package v1

import (
	"github.com/gin-gonic/gin"
	"go-weixin/pkg/log"
	"go-weixin/sdk"
)

// token 认证
func ViewWx(c *gin.Context) {
	appG := Gin{C: c}
	s := sdk.TokenSignature{}
	err := c.BindQuery(&s)
	if err != nil {
		log.Error(err)
		appG.SuccessResponse(false)
	} else if s.Confirm() {
		appG.SuccessResponse(true)
	} else {
		appG.SuccessResponse(false)
	}
}
