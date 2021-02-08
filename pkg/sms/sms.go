package sms

import "go-sso/conf"

type ISms interface {
	Send(telephone string, templateName string, params ...string) error
	SendBatch(telephoneArr []string, templateName string, params ...string) error
}

const (
	Login = "login"
	Sign  = "sign"
)

var sms ISms

func InitSms(config *conf.Config) {
	aliConfig := config.AliConfig
	sms = NewAliyunSms(aliConfig.AccessKey, aliConfig.AccessSecret)
}

func GetSms() ISms {
	return sms
}

// 发送登录短信
func SendLoginSms(s ISms, telephone, code string) error {
	return s.Send(telephone, "login", code)
}

// 发送注册短信
func SendSignSms(s ISms, telephone, code string) error {
	return s.Send(telephone, "sign", code)
}
