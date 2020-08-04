package email_tool

import (
	"go-sso/conf"
	"go-sso/pkg/log"
	"gopkg.in/gomail.v2"
)

var (
	user     string
	password string
	host     string
	port     int
	admin    []string
)

func InitEmail(config *conf.Config) {
	e := config.Email
	user = e.User
	password = e.Password
	host = e.Host
	port = e.Port
	admin = e.Admin
}

func SendEmail(to []string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", user)
	if to == nil {
		to = admin
	}
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(host, port, user, password)
	err := d.DialAndSend(m)
	if err != nil {
		log.Error(err)
	}
	return err
}
