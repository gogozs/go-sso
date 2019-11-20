package email_tool

func SendEmailCode(code, email string)  error {
	body := RegisterTmpl(code)
	err := SendEmail([]string{email}, "邮件验证", body)
	return err
}
