package sms

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"go-sso/pkg/log"
	"strings"
	"sync"
)

var aliyunSms *AliyunSms
var ali sync.Once

type AliyunSms struct {
	accessKeyId     string
	accessKeySecret string
}

func NewAliyunSms(key, secret string) Sms {
	ali.Do(func() {
		aliyunSms = &AliyunSms{
			accessKeyId:     key,
			accessKeySecret: secret,
		}
	})
	return aliyunSms
}

// 发送短信
// telephone 手机号
// templateName 模板名称
// params  短信内容参数
func (this *AliyunSms) Send(telephone string, templateName string, params ...string) (err error) {
	client, err := dysmsapi.NewClientWithAccessKey(
		"cn-shanghai",
		this.accessKeyId,
		this.accessKeySecret,
	)

	if err != nil {
		return
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	tmpl := GetAliSmsTmp(templateName)

	request.PhoneNumbers = telephone
	request.SignName = tmpl.SignName                                 // 签名 【逗比】
	request.TemplateCode = tmpl.TemplateId                           // 模板id ${code}
	request.TemplateParam = fmt.Sprintf(tmpl.TemplateParams, params) // 替换模板参数

	response, err := client.SendSms(request)
	if err != nil {
		log.Error(err)
	}
	log.Info("send sms succeed", response)
	return nil
}

func (this *AliyunSms) SendBatch(telephoneArr []string, templateName string, params ...string) error {
	return this.Send(strings.Join(telephoneArr, ","), templateName, params...)
}
