package auth

import (
	"github.com/gin-gonic/gin"
	"go-sso/internal/apierror"
	"go-sso/internal/middlewares/skipper"
	"go-sso/internal/repository/mysql/model"
	"go-sso/internal/service/viewset"
)

type AuthType string

const (
	CookieAuth AuthType = "cookie"
	JwtAuth    AuthType = "jwt"
	TokenAuth  AuthType = "token"
)

var driverList = map[AuthType]func() Auth{
	CookieAuth: func() Auth {
		return NewCookieAuthDriver()
	},
	JwtAuth: func() Auth {
		return NewJwtAuthDriver()
	},
	TokenAuth: func() Auth {
		return NewTokenAuthManager()
	},
}

type Auth interface {
	Check(c *gin.Context) error                         // 校验
	User(c *gin.Context) interface{}                    // 获取用户
	Login(c *gin.Context, user *model.User) interface{} // 登录
	Logout(c *gin.Context) bool                         // 登出
}

// 支持多种认证
// authList: 认证方式
// skipper: 跳过路由
func AuthMiddleware(authList []AuthType, skipper skipper.Skipper, m map[string]struct{}, prefixes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if skipper(c, m, prefixes...) {
			c.Next()
			return
		}
		if len(authList) == 0 {
			authList = []AuthType{TokenAuth} // default
		}
		var err error
		for _, authKey := range authList {
			driver := GenerateAuthDriver(authKey)
			if err = driver.Check(c); err == nil {
				c.Set("authKey", authKey)
				c.Next()
				return
			}
		}
		appG := viewset.ViewSet{}
		c.Header("WWW-Authenticate", "Token realm=\"Authorization Required\"")
		if e, ok := err.(apierror.ApiError); ok {
			appG.FailResponse(c, e)
		} else {
			appG.FailResponse(c, apierror.ErrAuth)
		}
		c.Abort()
	}
}

func GenerateAuthDriver(s AuthType) Auth {
	authDriver := driverList[s]()
	return authDriver
}

// 获取当前用户
func GetCurrentUser(c *gin.Context) *model.User {
	if authKey, ok := c.Get("authKey"); ok {
		driver := GenerateAuthDriver(authKey.(AuthType))
		return driver.User(c).(*model.User)
	}
	return &model.AnonymousUser
}

// 登出
func Logout(c *gin.Context) {
	if authKey, ok := c.Get("authKey"); ok {
		driver := GenerateAuthDriver(authKey.(AuthType))
		driver.Logout(c)
	}
}
