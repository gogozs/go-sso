package v1

import (
	"github.com/gin-gonic/gin"
	"go-weixin/pkg/apierror"
	"go-weixin/pkg/log"
	"go-weixin/service/models"
	"go-weixin/util"
)

// get token login
func ViewLogin(c *gin.Context) {
	appG := Gin{C: c}
	var user models.User
	err := c.BindJSON(user)
	if err != nil {
		log.Error(err)
		appG.FailResponse(apierror.INVALID_PARAMS)
	} else {
		if r := util.CheckUser(user.Username, user.Password); r {
			token, err := util.GenerateToken(user.Username, user.Password)
			if err != nil {
				appG.FailResponse(apierror.INVALID_PARAMS)
			} else {
				appG.SuccessResponse(gin.H{"token": token})
			}
		} else {
			appG.FailResponse(apierror.INVALID_PARAMS)
		}
	}
}


func ViewRegister(c *gin.Context) {
	appG := Gin{C: c}
	var user models.User
	err := c.BindJSON(user)
	if err != nil {
		log.Error(err)
		appG.FailResponse(apierror.INVALID_PARAMS)
	} else {
		if err := util.CreateUser(user); err != nil {
			appG.FailResponse(apierror.INVALID_PARAMS)
		} else {
			appG.SuccessResponse("注册成功")
		}
	}
}
