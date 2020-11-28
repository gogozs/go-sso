package auth

import (
	"github.com/gin-gonic/gin"
	"go-sso/db/model"
	"go-sso/pkg/log"
	"go-sso/registry"
	"go-sso/service/api/api_error"
	"math/rand"
	"strings"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	total   = len(letters)
)

type TokenAuthManager struct {
}

func NewTokenAuthManager() *TokenAuthManager {
	return &TokenAuthManager{}
}

func (a *TokenAuthManager) RandomToken() string {
	b := make([]byte, 16)
	for i := 0; i < 16; i++ {
		b[i] = letters[rand.Intn(total)]
	}
	return string(b)
}

func (a *TokenAuthManager) Check(c *gin.Context) error {
	token := c.Request.Header.Get("Authorization")
	token = strings.Replace(token, "Token ", "", -1)
	if token == "" {
		return api_error.ErrTokenInvalid
	}
	cacheStore := registry.GetCacheStore()
	if v, ok := cacheStore.GetCache(token); ok {
		user := v.(model.User)
		c.Set("User", &user)
		return nil
	}
	return api_error.ErrTokenInvalid
}

func (a *TokenAuthManager) User(c *gin.Context) interface{} {
	if user, exist := c.Get("User"); exist {
		return user
	} else {
		err := a.Check(c)
		if err != nil {
			log.Error(err)
			panic(err)
		} else {
			user, _ := c.Get("User")
			return user
		}
	}
}

func (a *TokenAuthManager) Login(c *gin.Context, u *model.User) interface{} {
	token := a.RandomToken()
	cacheStore := registry.GetCacheStore()
	// single sign on
	if oldToken, ok := cacheStore.GetCache(u.Username); ok {
		cacheStore.RemoveCache(oldToken.(string))
	}
	cacheStore.SetCache(token, *u)
	cacheStore.SetCache(u.Username, token)
	return gin.H{"token": token, "username": u.Username, "user_id": u.ID}
}

func (a *TokenAuthManager) Logout(c *gin.Context) bool {
	token := c.Request.Header.Get("Authorization")
	token = strings.Replace(token, "Token ", "", -1)
	cacheStore := registry.GetCacheStore()
	cacheStore.RemoveCache(token)
	return true
}
