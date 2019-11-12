package middlewares

import (
	"github.com/gin-gonic/gin"
	"go-sso/db/model"
	"go-sso/pkg/api_error"
	"go-sso/pkg/auth"
	"go-sso/service/api/viewset"
)

var driverList = map[string]func() Auth{
	"cookie": func() Auth {
		return auth.NewCookieAuthDriver()
	},
	"jwt": func() Auth {
		return auth.NewJwtAuthDriver()
	},
	"token": func() Auth {
		return auth.NewTokenAuthManager()
	},
}

type Auth interface {
	Check(c *gin.Context) error                         // 校验
	User(c *gin.Context) interface{}                    // 获取用户
	Login(c *gin.Context, user *model.User) interface{} // 登录
	Logout(c *gin.Context) bool                         // 登出
}

func RegisterGlobalAuthDriver(authKey string, key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		driver := GenerateAuthDriver(authKey)
		c.Set(key, driver)
		c.Next()
	}
}

type option struct {
	authList []string
	prefixes []string
}

// 支持多种认证
// authList: 认证方式
// skipper: 跳过路由
func AuthMiddleware(authList []string, skipper Skipper, prefixes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if skipper(c, prefixes...) {
			c.Next()
			return
		}
		if len(authList) == 0 {
			authList = []string{"jwt"} // default
		}
		var err error
		for _, authKey := range authList {
			driver := GenerateAuthDriver(authKey)
			if err = driver.Check(c); err == nil {
				c.Set("authDriver", authKey)
				c.Next()
				return
			}
		}
		appG := viewset.ViewSet{}
		c.Header("WWW-Authenticate", "Token realm=\"Authorization Required\"")
		if e, ok := err.(api_error.ApiError); ok {
			appG.FailResponse(c, e)
		} else {
			appG.FailResponse(c, api_error.ErrAuth)
		}
		c.Abort()
	}
}

func GenerateAuthDriver(s string) Auth {
	authDriver := driverList[s]()
	return authDriver
}

// 获取当前用户
func GetCurrentUser(c *gin.Context) model.User {
	if authKey, ok := c.Get("authKey"); ok {
		driver := GenerateAuthDriver(authKey.(string))
		return driver.User(c).(model.User)
	}
	return model.AnonymousUser
}

// 登出
func Logout(c *gin.Context) {
	if authKey, ok := c.Get("authKey"); ok {
		driver := GenerateAuthDriver(authKey.(string))
		driver.Logout(c)
	}
}
