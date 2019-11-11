package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"go-sso/conf"
	model2 "go-sso/db/model"
	"net/http"
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

func (cookie *cookieAuthManager) Check(c *gin.Context) error {
	// read cookie
	session, err := store.Get(c.Request, cookie.name)
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

func (cookie *cookieAuthManager) User(c *gin.Context) interface{} {
	// get model user
	session, err := store.Get(c.Request, cookie.name)
	if err != nil {
		return session.Values
	}
	return session.Values
}

func (cookie *cookieAuthManager) Login(http *http.Request, w http.ResponseWriter, user *model2.User) interface{} {
	// write cookie
	session, err := store.Get(http, cookie.name)
	if err != nil {
		return false
	}
	session.Values["id"] = user.ID
	session.Save(http, w)
	return true
}

func (cookie *cookieAuthManager) Logout(http *http.Request, w http.ResponseWriter) bool {
	// del cookie
	session, err := store.Get(http, cookie.name)
	if err != nil {
		return false
	}
	session.Values["id"] = nil
	session.Save(http, w)
	return true
}
