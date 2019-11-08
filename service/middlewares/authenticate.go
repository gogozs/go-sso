package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-qiuplus/db/model"
	"go-qiuplus/pkg/api_error"
	"go-qiuplus/pkg/auth"
	"go-qiuplus/service/api/viewset"
	"net/http"
)

var driverList = map[string]func() Auth{
	"cookie": func() Auth {
		return auth.NewCookieAuthDriver()
	},
	"jwt": func() Auth {
		return auth.NewJwtAuthDriver()
	},
}

type Auth interface {
	Check(c *gin.Context) error
	User(c *gin.Context) interface{}
	Login(http *http.Request, w http.ResponseWriter, user *model.User) interface{}
	Logout(http *http.Request, w http.ResponseWriter) bool
}

func RegisterGlobalAuthDriver(authKey string, key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		driver := GenerateAuthDriver(authKey)
		c.Set(key, driver)
		c.Next()
	}
}

// 支持多种认证
func AuthMiddleware(authList []string, skipper Skipper, prefixes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if skipper(c, prefixes...) {
			c.Next()
			return
		}
		var err error
		for _, authKey := range authList {
			driver := GenerateAuthDriver(authKey)
			if err = (*driver).Check(c); err == nil {
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

func GenerateAuthDriver(s string) *Auth {
	var authDriver Auth
	fmt.Println(driverList[s])
	authDriver = driverList[s]()
	return &authDriver
}

func GetCurrentUser(c *gin.Context, key string) map[string]interface{} {
	authDriver, _ := c.MustGet(key).(*Auth)
	return (*authDriver).User(c).(map[string]interface{})
}

func User(c *gin.Context) map[string]interface{} {
	return GetCurrentUser(c, "jwt_auth")
}
