package v1

import (
	"github.com/gin-gonic/gin"
	"go-sso/pkg/wx/wx_client"
	"go-sso/service/api/viewset"
)

type WxViewset struct {
	viewset.ViewSet
}

func (a *WxViewset) ErrorHandler(f func(c *gin.Context) error) func(c *gin.Context) {
	return func(c *gin.Context) {
		a.ViewSet.ErrorHandler(f, c)
	}
}

func (a *WxViewset) Login(c *gin.Context) (err error) {
	lp := &wx_client.LoginParams{}
	err = c.ShouldBind(lp)
	if err != nil {
		return
	}
	wx := wx_client.GetWxClient()
	res, err := wx.Login(lp)
	if err != nil {
		return
	}
	return a.SuccessResponse(c, res)
}

func (a *WxViewset) GetUserInfo(c *gin.Context) (err error) {
	return
}
