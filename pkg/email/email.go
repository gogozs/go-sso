package email

import (
	"go-qiuplus/conf"
	"net/smtp"
	"strings"
)

var (
	user     string
	password string
	host     string
	admin    []string
)

func init() {
	e := conf.GetConfig().Email
	user = e.User
	password = e.Password
	host = e.Host
	admin = []string{e.Admin}
}

func SendEmail(to []string, subject, body, mailtype string) error {
	if to == nil {
		to = admin
	}
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}
	tostr := strings.Join(to, ",")
	msg := []byte("To: " + tostr + "\r\nFrom: " + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	err := smtp.SendMail(host, auth, user, to, msg)
	return err
}
