package auth

import (
	"github.com/gin-gonic/gin"
	"go-qiuplus/conf"
	model2 "go-qiuplus/db/model"
	"go-qiuplus/db/query"
	"go-qiuplus/pkg/api_error"
	"go-qiuplus/pkg/log"
	"go-qiuplus/util"
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
	jwt := conf.GetConfig().Jwt
	return &jwtAuthManager{
		secret: jwt.SECRET,
		exp:    jwt.EXP,
		alg:    jwt.ALG,
	}
}

// Check the token of request header is valid or not.
func (jwtAuth *jwtAuthManager) Check(c *gin.Context) error {
	token := c.Request.Header.Get("Authorization")
	token = strings.Replace(token, "Token ", "", -1)
	if token == "" {
		return api_error.ErrTokenInvalid
	}
	clamis, err := util.ParseToken(token)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	username := clamis.Username
	user, err := query.UserQ.GetUserByName(username)
	if err != nil {
		return err
	} else {
		c.Set("User", user)
		return nil
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
			c.Set("User", model2.AnonymousUser)
			return nil
		}
		clamis, err := util.ParseToken(token)
		if err != nil {
			log.Error(err.Error())
			panic(err)
		} else if time.Now().Unix() > clamis.ExpiresAt {
			return nil
		}
		username := clamis.Username
		user, err := query.UserQ.GetUserByName(username)
		if err != nil {
			panic(err)
		} else {
			c.Set("User", user)
			return user
		}
	}
}

func (jwtAuth *jwtAuthManager) Login(http *http.Request, w http.ResponseWriter, user *model2.User) interface{} {
	token, _ := util.GenerateToken(user.Username, user.Password)
	return token
}

func (jwtAuth *jwtAuthManager) Logout(http *http.Request, w http.ResponseWriter) bool {
	// TODO: 逻辑补充
	return true
}
