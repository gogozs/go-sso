package auth

import (
	"github.com/gin-gonic/gin"
	"go-sso/conf"
	"go-sso/db/model"
	"go-sso/db/query"
	"go-sso/pkg/api_error"
	"go-sso/pkg/log"
	"go-sso/util"
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
func (this *jwtAuthManager) Check(c *gin.Context) error {
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
	user, err := query.UserQ.GetUserByAccount(username)
	if err != nil {
		return err
	}
	c.Set("User", user)
	return nil
}

// 获取user
func (this *jwtAuthManager) User(c *gin.Context) interface{} {
	if user, exist := c.Get("User"); exist {
		return user
	} else {
		err := this.Check(c)
		if err != nil {
			log.Error(err)
			panic(err)
		} else {
			user, _ := c.Get("User")
			return user
		}
	}
}

func (this *jwtAuthManager) Login(c *gin.Context, user *model.User) interface{} {
	token, _ := util.GenerateToken(user.Username, user.Password)
	return token
}

func (this *jwtAuthManager) Logout(c *gin.Context) bool {
	// TODO: 逻辑补充
	return true
}
