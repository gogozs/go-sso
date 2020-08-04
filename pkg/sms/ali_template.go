package sms

import "go-sso/conf"

var aliSms = map[string]*conf.AliSmsTemplate{}

// 通过名称获取模板
func GetAliSmsTmp(name string) *conf.AliSmsTemplate {
	return aliSms[name]
}

func InitAliConfig(config *conf.Config) {
	aliConfig := config.AliSms

	for _, item := range aliConfig {
		aliSms[item.TemplateName] = &item
	}

	InitSms(config)
}
