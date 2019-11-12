package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"go-sso/conf"
	"go-sso/db/model"
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

func (this *cookieAuthManager) Check(c *gin.Context) error {
	// read this
	session, err := store.Get(c.Request, this.name)
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

func (this *cookieAuthManager) User(c *gin.Context) interface{} {
	// get model user
	session, err := store.Get(c.Request, this.name)
	if err != nil {
		log.Error(err)
	}
	return session.Values
}

func (this *cookieAuthManager) Login(c *gin.Context, user *model.User) interface{} {
	// write this
	session, err := store.Get(c.Request, this.name)
	if err != nil {
		return false
	}
	session.Values["id"] = user.ID
	session.Save(c.Request, c.Writer)
	return true
}

func (this *cookieAuthManager) Logout(c *gin.Context) bool {
	// del this
	session, err := store.Get(c.Request, this.name)
	if err != nil {
		return false
	}
	session.Values["id"] = nil
	session.Save(c.Request, c.Writer)
	return true
}
