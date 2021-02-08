package auth

import (
	"github.com/gin-gonic/gin"
	"go-sso/internal/apierror"
	"go-sso/internal/repository/storage/mysql"
	"go-sso/pkg/json"
	"go-sso/pkg/log"
	"go-sso/registry"
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
		return apierror.ErrTokenInvalid
	}
	cacheStore := registry.GetCacheStore()
	v, err := cacheStore.GetCache(token)
	if err != nil {
		return err
	}
	if v == "" {
		return apierror.ErrTokenInvalid
	}
	user, err := mysql.UnmarshalUser(v)
	if err != nil {
		_, _ = cacheStore.RemoveCache(token)
		return apierror.ErrTokenInvalid
	}
	c.Set("User", &user)

	return apierror.ErrTokenInvalid
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

func (a *TokenAuthManager) Login(c *gin.Context, u *mysql.User) (interface{}, error) {
	token := a.RandomToken()
	cacheStore := registry.GetCacheStore()
	// single sign on
	oldToken, err := cacheStore.GetCache(u.Username)
	if err != nil {
		return nil, err
	}
	if oldToken != "" {
		_, _ = cacheStore.RemoveCache(oldToken)
	}

	b, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	if err := cacheStore.SetCache(token, string(b)); err != nil {
		return nil, err
	}
	if err := cacheStore.SetCache(u.Username, token); err != nil {
		return nil, err
	}
	return gin.H{"token": token, "username": u.Username, "user_id": u.ID}, nil
}

func (a *TokenAuthManager) Logout(c *gin.Context) bool {
	token := c.Request.Header.Get("Authorization")
	token = strings.Replace(token, "Token ", "", -1)
	cacheStore := registry.GetCacheStore()
	_, _ = cacheStore.RemoveCache(token)
	return true
}
