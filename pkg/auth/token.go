package auth

import (
	"github.com/gin-gonic/gin"
	"go-sso/db/model"
	"go-sso/pkg/api_error"
	"go-sso/pkg/log"
	"go-sso/pkg/storage"
	"math/rand"
	"strings"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	total = len(letters)
)

type TokenAuthManager struct {
}

func NewTokenAuthManager() *TokenAuthManager {
	return &TokenAuthManager{}
}

func (this *TokenAuthManager) RandomToken() string {
	b := make([]byte, 16)
	for i:=0;i<16;i++ {
		b[i] = letters[rand.Intn(total)]
	}
	return string(b)
}


func (this *TokenAuthManager) Check(c *gin.Context) error {
	token := c.Request.Header.Get("Authorization")
	token = strings.Replace(token, "Token ", "", -1)
	if token == "" {
		return api_error.ErrTokenInvalid
	}
	cacheStore := storage.GetStore()
	if v, ok := cacheStore.GetCache(token); ok {
		var user model.User
		//err := json.Unmarshal(v.([]byte), &user)
		user = v.(model.User)
		c.Set("User", &user)
		return nil
	}
	return api_error.ErrTokenInvalid
}

func (this *TokenAuthManager) User(c *gin.Context) interface{} {
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

func (this *TokenAuthManager) Login(c *gin.Context, u *model.User) interface{} {
	token := this.RandomToken()
	cacheStore := storage.GetStore()
	cacheStore.SetCache(token, *u)
	cacheStore.SetCache(u.Username, token)
	return gin.H{"token": token}
}

func (this *TokenAuthManager) Logout(c *gin.Context) bool {
	token := c.Request.Header.Get("Authorization")
	token = strings.Replace(token, "Token ", "", -1)
	cacheStore := storage.GetStore()
	cacheStore.RemoveCache(token)
	return true
}
