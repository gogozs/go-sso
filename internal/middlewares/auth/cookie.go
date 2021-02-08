package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"go-sso/conf"
	"go-sso/internal/repository/storage/mysql"
	"go-sso/pkg/log"
)

var store = sessions.NewCookieStore([]byte(conf.GetConfig().Common.AppSecret))

type cookieAuthManager struct {
	name string
}

func NewCookieAuthDriver() *cookieAuthManager {
	return &cookieAuthManager{
		name: conf.GetConfig().Cookie.NAME,
	}
}

func (a *cookieAuthManager) Check(c *gin.Context) error {
	// read a
	session, err := store.Get(c.Request, a.name)
	if err != nil {
		return errors.New("session invalid")
	}
	if session == nil {
		return errors.New("session invalid")
	}
	if session.Values == nil {
		return errors.New("session invalid")
	}
	if session.Values["id"] == nil {
		return errors.New("session invalid")
	}
	return nil
}

func (a *cookieAuthManager) User(c *gin.Context) interface{} {
	// get model user
	session, err := store.Get(c.Request, a.name)
	if err != nil {
		log.Error(err)
	}
	return session.Values
}

func (a *cookieAuthManager) Login(c *gin.Context, user *mysql.User) (interface{}, error) {
	// write a
	session, err := store.Get(c.Request, a.name)
	if err != nil {
		return false, err
	}
	session.Values["id"] = user.ID
	if err := session.Save(c.Request, c.Writer); err != nil {
		return false, err
	}
	return true, nil
}

func (a *cookieAuthManager) Logout(c *gin.Context) bool {
	// del a
	session, err := store.Get(c.Request, a.name)
	if err != nil {
		return false
	}
	session.Values["id"] = nil
	if err := session.Save(c.Request, c.Writer); err != nil {
		return false
	}
	return true
}
