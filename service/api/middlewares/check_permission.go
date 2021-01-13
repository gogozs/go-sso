package middlewares

import (
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"go-sso/pkg/permission"
	"go-sso/service/api/api_error"
	"go-sso/service/api/viewset"
	"go-sso/storage/mysql/model"
	"net/http"
)

// NewAuthorizer returns the authorizer, uses a Casbin enforcer as input
func PermissionMiddleware(skipper Skipper, m map[string]struct{}, prefixes ...string) gin.HandlerFunc {
	enforcer := permission.GetEnforcer()
	pa := &PermissionAuthorizer{enforcer: enforcer}

	return func(c *gin.Context) {
		if !skipper(c, m, prefixes...) {
			if !pa.CheckPermission(c) {
				pa.RequirePermission(c)
			}
		}
	}
}

type PermissionAuthorizer struct {
	enforcer *casbin.Enforcer
}

func (this *PermissionAuthorizer) GetUser(c *gin.Context) *model.User {
	user := c.MustGet("User").(*model.User)
	return user
}

func (this *PermissionAuthorizer) CheckPermission(c *gin.Context) bool {
	user := this.GetUser(c)
	if user.ID == 0 {
		return false // AnonymousUser
	}
	method := c.Request.Method
	path := c.Request.URL.Path
	return this.enforcer.Enforce(user.Role, path, method)
}

// RequirePermission returns the 403 Forbidden to the client
func (this *PermissionAuthorizer) RequirePermission(c *gin.Context) {
	c.JSON(http.StatusForbidden, viewset.GetFailResponse(api_error.ErrPermission, nil))
	c.Abort()
}
