package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-sso/db/inter"
	"go-sso/db/model"
	"go-sso/pkg/api_error"
	"go-sso/pkg/email_tool"
	"go-sso/pkg/log"
	"go-sso/pkg/storage"
	"go-sso/service/api/viewset"
	"go-sso/service/middlewares"
	"go-sso/util"
)

type AuthViewset struct {
	itemInter inter.IUser
	viewset.ViewSet
}

func (this *AuthViewset) ErrorHandler(f func(c *gin.Context) error) func(c *gin.Context) {
	return func(c *gin.Context) {
		this.ViewSet.ErrorHandler(f, c)
	}
}

// @Summary user login
// @Description 1.账号密码登录 2.手机号，邮箱登录
// @Accept  json
// @Produce  json
// @Param  user body model.UserParams true "username && password"
// @Success 200 {object} viewset.Response
// @Router /api/public/v1/auth/login/ [post]
func (this *AuthViewset) Login(c *gin.Context) (err error) {
	var up model.UserParams
	err = c.ShouldBind(&up)
	if err != nil {
		log.Error(err.Error())
		this.FailResponse(c, api_error.ErrInvalid)
		return api_error.ErrInvalid
	} else {
		if u, r := this.itemInter.CheckUser(up.Account, up.Password); r {
			// 登录方式 token
			driver := middlewares.GenerateAuthDriver(middlewares.TokenAuth)
			res := driver.Login(c, u)
			return this.SuccessResponse(c, res)
		} else {
			return api_error.ErrAuth
		}
	}
}

// @Summary telephone login
// @Description telephone login steps 1.IsTelephoneExist 2. SendSmsCode 3. TelephoneLogin
// @Accept  json
// @Produce  json
// @Param  user body model.UserParams true "username && password"
// @Success 200 {object} viewset.Response
// @Router /api/public/v1/auth/telephone/login/ [post]
func (this *AuthViewset) TelephoneLogin(c *gin.Context) (err error) {
	var tl model.TelephoneLoginParams
	err = c.ShouldBind(&tl)
	err = this.VerifySmsCode(tl.Telephone, tl.Code)
	if err != nil {
		return
	}
	// 登录方式 token
	u, err := this.itemInter.GetUserByAccount(tl.Telephone)
	if err != nil {
		return
	}
	driver := middlewares.GenerateAuthDriver(middlewares.TokenAuth)
	res := driver.Login(c, u)
	return this.SuccessResponse(c, res)
}

// @Summary user register
// @Description register by username, telephone and password
// @Accept  json
// @Produce  json
// @Param  user body model.UserParams true "username && password"
// @Success 200 {object} viewset.Response
// @Router /api/public/v1/auth/register/ [post]
func (this *AuthViewset) Register(c *gin.Context) (err error) {
	var rp model.RegisterParams
	var user model.User
	err = c.ShouldBind(&rp)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	err = rp.Validate()
	if err != nil {
		return err
	}
	errs := this.CheckRegisterParams(&rp)
	if len(errs) > 0 {
		this.FailResponse(c, api_error.ErrInvalid, errs)
		return
	}
	newPassword, err := util.GeneratePassword(rp.Password)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	user.Username = rp.Username
	user.Telephone = rp.Telephone
	user.Email = rp.Email
	user.Password = newPassword
	if _, err = this.itemInter.Create(&user); err != nil {
		log.Error(err.Error())
		return err
	} else {
		return this.SuccessResponse(c, "注册成功")
	}
}

// 检测注册用户参数
func (this *AuthViewset) CheckRegisterParams(rp *model.RegisterParams) map[string]string {
	errs := make(map[string]string)
	// 检查参数是否合法
	if !this.itemInter.IsValid(rp.Username, "username") {
		errs["username"] = "用户名至少3位以上字母开头"
	}
	if !this.itemInter.IsValid(rp.Telephone, "telephone") {
		errs["telephone"] = "手机号格式错误"
	}
	if rp.Email != "" && !this.itemInter.IsValid(rp.Email, "email") {
		errs["email"] = "email格式错误"
	}
	if len(errs) > 0 {
		return errs
	}
	// 检查是否重复注册
	if this.itemInter.Exists(rp.Username, "username") {
		errs["username"] = "用户已经存在"
	}
	if this.itemInter.Exists(rp.Telephone, "telephone") {
		errs["telephone"] = "手机号已经存在"
	}
	if rp.Email != "" && this.itemInter.Exists(rp.Email, "email") {
		errs["email"] = "email已经存在"
	}
	return errs
}

// @Summary 账号注册手机号校验是否合法
// @Description send telephone verify code
// @Accept  json
// @Produce  json
// @Param  user body model.UserParams true "username && password"
// @Success 200 {object} viewset.Response
// @Router /api/public/v1/auth/check-telephone/ [post]
func (this *AuthViewset) CheckTelephoneValid(c *gin.Context) (err error) {
	errs := make(map[string]string)
	// 检查参数是否合法
	telephone := c.Query("telephone")
	if !this.itemInter.IsValid(telephone, "telephone") {
		errs["telephone"] = "手机号格式错误"
	}
	if this.itemInter.Exists(telephone, "telephone") {
		errs["telephone"] = "手机号已经存在"
	}
	if len(errs) != 0 {
		this.FailResponse(c, api_error.ErrInvalid, errs)
		return
	}
	return this.SuccessBlankResponse(c)
}

// @Summary check telephone
// @Description check telephone whether exist
// @Accept  json
// @Produce  json
// @Param  user body  true "telephone"
// @Success 200 {object} viewset.Response
// @Router /api/public/v1/auth/check-telephone-exist/ [post]
func (this *AuthViewset) CheckTelephoneExist(c *gin.Context) (err error) {
	errs := make(map[string]string)
	// 检查参数是否合法
	telephone := c.Query("telephone")
	if !this.itemInter.IsValid(telephone, "telephone") {
		errs["telephone"] = "手机号格式错误"
	}
	if !this.itemInter.Exists(telephone, "telephone") {
		errs["telephone"] = "手机号不存在"
	}
	if len(errs) != 0 {
		this.FailResponse(c, api_error.ErrInvalid, errs)
		return
	}
	return this.SuccessBlankResponse(c)
}

// @Summary 发送手机验证码
// @Description send telephone verify code
// @Accept  json
// @Produce  json
// @Param  user body model.UserParams true "username && password"
// @Success 200 {object} viewset.Response
// @Router /api/public/v1/auth/register/ [post]
func (this *AuthViewset) SendSmsCode(c *gin.Context) (err error) {
	// TODO 发送短信
	code := "123456"
	telephone := c.Query("telephone")
	if ok := this.itemInter.IsValid(telephone, "telephone"); !ok {
		return api_error.ErrInvalid
	}
	cacheStore := storage.GetStore()
	cacheStore.SetCache(telephone, code)
	return this.SuccessBlankResponse(c)
}

// @Summary telephone check
// @Description 手机验证码确认
func (this *AuthViewset) VerifySmsCode(telephone, code string) (err error) {
	cacheStore := storage.GetStore()
	if rightCode, ok := cacheStore.GetCache(telephone); ok && rightCode.(string) == code {
		return
	}
	return nil
}

// @Summary 发送邮箱验证码
// @Description register by username and password
// @Accept  json
// @Produce  json
// @Param  user body model.UserParams true "username && password"
// @Success 200 {object} viewset.Response
// @Router /api/public/v1/auth/send-email-code/ [post]
func (this *AuthViewset) SendEmailCode(c *gin.Context) (err error) {
	email := c.Query("email")
	if ok := this.itemInter.IsValid(email, "email"); !ok {
		return api_error.ErrInvalid
	}
	cacheStore := storage.GetStore()
	code := util.RandomCode()
	err = email_tool.SendEmailCode(code, email)
	if err != nil {
		return err
	}
	cacheStore.SetCache(email, code)
	return this.SuccessResponse(c, gin.H{"url": ""})
}

func (this *AuthViewset) VerifyEmailCode(email, code string) (err error) {
	cacheStore := storage.GetStore()
	if rightCode, ok := cacheStore.GetCache(email); ok && rightCode.(string) == code {
		return
	}
	return nil
}

// @Summary 重置密码
// @Description reset password by telephone or email
// @Accept  json
// @Produce  json
// @Param  user body model.UserParams true "username && password"
// @Success 200 {object} viewset.Response
// @Router /api/public/v1/auth/reset-password/ [post]
func (this *AuthViewset) ResetPassword(c *gin.Context) (err error) {
	var rp model.ResetPasswordParams
	err = c.ShouldBind(&rp)
	if err != nil {
		return
	}
	user, err := this.itemInter.GetUserByAccount(rp.Account)
	if err != nil {
		return
	}
	switch rp.VerifyType {
	case "email":
		err = this.VerifyEmailCode(user.Email, rp.Code)
	case "telephone":
		err = this.VerifySmsCode(user.Telephone, rp.Code)
	default:
		err = errors.New("验证类型错误")
	}
	if err != nil {
		return
	}
	err = this.itemInter.ChangePassword(user, rp.NewPassword)
	if err != nil {
		return
	}
	return this.SuccessBlankResponse(c)
}

// @Summary 修改密码
// @Description change password by username and password
// @Accept  json
// @Produce  json
// @Param  user body model.ChangePasswordParams true "raw_password && new_password"
// @Success 200 {object} viewset.Response
// @Router /api/v1/auth/change-password/ [post]
func (this *AuthViewset) ChangePassword(c *gin.Context) (err error) {
	var cp model.ChangePasswordParams
	err = c.ShouldBind(&cp)
	if err != nil {
		log.Error(err)
		return err
	}
	err = cp.Validate()
	if err != nil {
		log.Error(err)
		return err
	}
	username := middlewares.GetCurrentUser(c).Username
	if u, ok := this.itemInter.CheckUser(username, cp.RawPassword); ok {
		err = this.itemInter.ChangePassword(u, cp.NewPassword)
		if err != nil {
			log.Error(err)
			return
		}
		return this.SuccessBlankResponse(c)
	}
	return errors.New("原密码错误")
}
