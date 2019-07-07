package drivers

import (
	"github.com/gin-gonic/gin"
	"go-weixin/config"
	"go-weixin/pkg/apierror"
	"go-weixin/pkg/log"
	"go-weixin/util"
	"net/http"
	"strings"
	"time"
)

type jwtAuthManager struct {
	secret string
	exp    time.Duration
	alg    string
}

func NewJwtAuthDriver() *jwtAuthManager {
	jwt := config.GetConfig().Jwt
	return &jwtAuthManager{
		secret: jwt.SECRET,
		exp:    jwt.EXP,
		alg:    jwt.ALG,
	}
}

// Check the token of request header is valid or not.
func (jwtAuth *jwtAuthManager) Check(c *gin.Context) bool {
	token := c.Request.Header.Get("Authorization")
	token = strings.Replace(token, "Token ", "", -1)
	if token == "" {
		return false
	}
	clamis, err := util.ParseToken(token)
	if err != nil {
		log.Error(err)
		return false
	} else if time.Now().Unix() > clamis.ExpiresAt {
		log.Info(apierror.GetMsg(apierror.ERROR_AUTH))
		return false
	}
	username := clamis.Username
	user, err := users.GetUser(username)
	if err != nil {
		log.Info(apierror.GetMsg(apierror.ERROR_AUTH))
		return false
	} else {
		c.Set("User", user)
		return true
	}
}

// 获取user
func (jwtAuth *jwtAuthManager) User(c *gin.Context) interface{} {
	if user, exist := c.Get("User"); exist {
		return user
	} else {
		token := c.Request.Header.Get("Authorization")
		token = strings.Replace(token, "Token ", "", -1)
		if token == "" {
			return nil
		}
		clamis, err := util.ParseToken(token)
		if err != nil {
			log.Info(err)
			panic(err)
		} else if time.Now().Unix() > clamis.ExpiresAt {
			log.Info(apierror.GetMsg(apierror.ERROR_AUTH))
			return nil
		}
		username := clamis.Username
		user, err := users.GetUser(username)
		if err != nil {
			log.Info(apierror.GetMsg(apierror.ERROR_AUTH))
			panic(err)
		} else {
			c.Set("User", user)
			return user
		}
	}
}

func (jwtAuth *jwtAuthManager) Login(http *http.Request, w http.ResponseWriter, user *users.User) interface{} {
	token, _ := util.GenerateToken(user.Username, user.Password)
	return token
}

func (jwtAuth *jwtAuthManager) Logout(http *http.Request, w http.ResponseWriter) bool {
	// TODO: 逻辑补充
	return true
}
