package email_tool

import (
	"bytes"
	"html/template"
)

var (
	validCodeTmpl = `<div> 您在注册账号，验证码为{{ . }} </div>`
	activateTmpl   = `<div> 请点击 <a href="{{ . }}">激活链接</a> 激活账号</div`
)

func RegisterTmpl(code string) string {
	tmpl, _ := template.New("tmp").Parse(validCodeTmpl)
	var b bytes.Buffer
	_ = tmpl.Execute(&b, code)
	return b.String()
}

func ActivateTmpl(url string) string {
	tmpl, _ := template.New("tmp").Parse(validCodeTmpl)
	var b bytes.Buffer
	_ = tmpl.Execute(&b, url)
	return b.String()
}
