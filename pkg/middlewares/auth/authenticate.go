package auth

import (
	"github.com/gin-gonic/gin"
	"go-weixin/pkg/apierror"
	"go-weixin/pkg/middlewares/auth/drivers"
	"net/http"
)

var driverList = map[string]func() Auth{
	"cookie": func() Auth {
		return drivers.NewCookieAuthDriver()
	},
	"jwt": func() Auth {
		return drivers.NewJwtAuthDriver()
	},
}

type Auth interface {
	Check(c *gin.Context) bool
	User(c *gin.Context) interface{}
	Login(http *http.Request, w http.ResponseWriter, user *users.User) interface{}
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
func AuthMiddleware(authList []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		loginSuccess := false // 认证结果
		for _, authKey := range authList {
			driver := GenerateAuthDriver(authKey)
			if (*driver).Check(c) {
				loginSuccess = true
				if loginSuccess {
					break
				}
			}
		}
		if !loginSuccess {
			appG := app.Gin{C: c}
			c.Header("WWW-Authenticate", "Token realm=\"Authorization Required\"")
			appG.Response(http.StatusUnauthorized, apierror.UNAUTHORIZED, gin.H{"msg": apierror.GetMsg(apierror.UNAUTHORIZED)})
			c.Abort()
		}
		c.Next()
	}
}

func GenerateAuthDriver(string string) *Auth {
	var authDriver Auth
	authDriver = driverList[string]()
	return &authDriver
}

func GetCurrentUser(c *gin.Context, key string) map[string]interface{} {
	authDriver, _ := c.MustGet(key).(*Auth)
	return (*authDriver).User(c).(map[string]interface{})
}

func User(c *gin.Context) map[string]interface{} {
	return GetCurrentUser(c, "jwt_auth")
}
