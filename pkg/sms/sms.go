package sms

import "go-sso/conf"

type Sms interface {
	Send(telephone string, templateName string, params ...string) error
	SendBatch(telephoneArr []string, templateName string, params ...string) error
}

const (
	Login = "login"
	Sign  = "sign"
)

var sms Sms

func init() {
	aliConfig := conf.GetConfig().AliConfig
	sms = NewAliyunSms(aliConfig.AccessKey, aliConfig.AccessSecret)
}

// 发送登录短信
func SendLoginSms(telephone, code string) error {
	return sms.Send(telephone, "login", code)
}

// 发送注册短信f
func SendSignSms(telephone, code string) error {
	return sms.Send(telephone, "sign", code)
}
