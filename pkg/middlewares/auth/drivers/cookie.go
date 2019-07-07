package drivers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"go-weixin/config"
	"go-weixin/models/users"
	"net/http"
)

var store = sessions.NewCookieStore([]byte(config.GetConfig().Common.APP_SECRET))

type cookieAuthManager struct {
	name string
}

func NewCookieAuthDriver() *cookieAuthManager {
	return &cookieAuthManager{
		name: config.GetConfig().Cookie.NAME,
	}
}

func (cookie *cookieAuthManager) Check(c *gin.Context) bool {
	// read cookie
	session, err := store.Get(c.Request, cookie.name)
	if err != nil {
		return false
	}
	if session == nil {
		return false
	}
	if session.Values == nil {
		return false
	}
	if session.Values["id"] == nil {
		return false
	}
	return true
}

func (cookie *cookieAuthManager) User(c *gin.Context) interface{} {
	// get model user
	session, err := store.Get(c.Request, cookie.name)
	if err != nil {
		return session.Values
	}
	return session.Values
}

func (cookie *cookieAuthManager) Login(http *http.Request, w http.ResponseWriter, user *users.User) interface{} {
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
