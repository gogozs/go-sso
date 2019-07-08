package permission

import (
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"go-weixin/pkg/apierror"
	"go-weixin/service/api"
	"go-weixin/service/models"
)

// NewAuthorizer returns the authorizer, uses a Casbin enforcer as input
func PermissionMiddleware() gin.HandlerFunc {
	enforcer := Casbin()
	a := &PermissionAuthorizer{enforcer: enforcer}

	return func(c *gin.Context) {
		if !a.CheckPermission(c) {
			a.RequirePermission(c)
		}
	}
}

type PermissionAuthorizer struct {
	enforcer *casbin.Enforcer
}

func (a *PermissionAuthorizer) GetUserName(c *gin.Context) string {
	user:= c.MustGet("User").(models.User)
	return user.Role.String
}

func (a *PermissionAuthorizer) CheckPermission(c *gin.Context) bool {
	role := a.GetUserName(c)
	method := c.Request.Method
	path := c.Request.URL.Path
	return a.enforcer.Enforce(role, path, method)
}

// RequirePermission returns the 403 Forbidden to the client
func (a *PermissionAuthorizer) RequirePermission(c *gin.Context) {
	c.JSON(403, api.Response{
		Code: apierror.ERROR_AUTH,
		Msg:  apierror.GetMsg(apierror.ERROR_AUTH),
	})
	c.Abort()
}