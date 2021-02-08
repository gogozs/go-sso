package permissions

import (
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"go-sso/internal/apierror"
	"go-sso/internal/middlewares/skipper"
	"go-sso/internal/repository/storage/mysql"
	"go-sso/internal/service/viewset"
	"go-sso/pkg/permission"
	"net/http"
)

// NewAuthorizer returns the authorizer, uses a Casbin enforcer as input
func PermissionMiddleware(skipper skipper.Skipper, m map[string]struct{}, prefixes ...string) gin.HandlerFunc {
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

func (a *PermissionAuthorizer) GetUser(c *gin.Context) *mysql.User {
	user := c.MustGet("User").(*mysql.User)
	return user
}

func (a *PermissionAuthorizer) CheckPermission(c *gin.Context) bool {
	user := a.GetUser(c)
	if user.ID == 0 {
		return false // AnonymousUser
	}
	method := c.Request.Method
	path := c.Request.URL.Path
	return a.enforcer.Enforce(user.Role, path, method)
}

// RequirePermission returns the 403 Forbidden to the client
func (a *PermissionAuthorizer) RequirePermission(c *gin.Context) {
	c.JSON(http.StatusForbidden, viewset.GetFailResponse(apierror.ErrPermission, nil))
	c.Abort()
}
