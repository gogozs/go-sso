package v1

import (
	"github.com/gin-gonic/gin"
	"go-sso/db/inter"
	"go-sso/db/model"
	"go-sso/pkg/api_error"
	"go-sso/pkg/log"
	"go-sso/pkg/storage"
	"go-sso/service/api/viewset"
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
		if r := this.itemInter.CheckUser(up.Account, up.Password); r {
			token, err := util.GenerateToken(up.Account, up.Password)
			if err != nil {
				return api_error.ErrInternal
			} else {
				return this.SuccessResponse(c, gin.H{"token": token})
			}
		} else {
			return api_error.ErrAuth
		}
	}
}

// @Summary telephone login
// @Description 手机验证码登录
// @Accept  json
// @Produce  json
// @Param  user body model.UserParams true "username && password"
// @Success 200 {object} viewset.Response
// @Router /api/public/v1/auth/login/ [post]
func (this *AuthViewset) TelephoneLogin(c *gin.Context) (err error) {
	return
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
	if this.itemInter.Exists(rp.Username, "username") {
		errs["username"] = "用户已经存在"
	}
	if this.itemInter.Exists(rp.Telephone, "telephone") {
		errs["telephone"] = "手机号已经存在"
	}
	if this.itemInter.Exists(rp.Email, "email") {
		errs["email"] = "email已经存在"
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
	if this.itemInter.Exists(rp.Email, "email") {
		errs["email"] = "email已经存在"
	}
	return errs
}

// @Summary 发送手机验证码
// @Description register by username and password
// @Accept  json
// @Produce  json
// @Param  user body model.UserParams true "username && password"
// @Success 200 {object} viewset.Response
// @Router /api/public/v1/auth/register/ [post]
func (this *AuthViewset) SendSmsCode(c *gin.Context) (err error) {
	return this.SuccessBlackResponse(c)
}


// @Summary telephone check
// @Description 手机验证码确认
// @Accept  json
// @Produce  json
// @Param  user body model.UserParams true "username && password"
// @Success 200 {object} viewset.Response
// @Router /api/public/v1/auth/register/ [post]
func (this *AuthViewset) SmsCodeValid(c *gin.Context) (err error) {
	// TODO	接口检测
	return this.SuccessBlackResponse(c)
}

// @Summary 发送邮箱验证码
// @Description register by username and password
// @Accept  json
// @Produce  json
// @Param  user body model.UserParams true "username && password"
// @Success 200 {object} viewset.Response
// @Router /api/public/v1/auth/register/ [post]
func (this *AuthViewset) SendEmailCode(c *gin.Context) (err error) {
	email := c.Query("email")
	if ok := this.itemInter.IsValid(email, "email"); !ok {
		return api_error.ErrInvalid
	}
	cacheStore := storage.GetStore()
	code := util.RandomCode()
	cacheStore.SetCache(email, code)
	// TODO 发送邮件

	return this.SuccessBlackResponse(c)
}

// @Summary email check
// @Description email验证码确认
// @Accept  json
// @Produce  json
// @Param  user body model.UserParams true "username && password"
// @Success 200 {object} viewset.Response
// @Router /api/public/v1/auth/register/ [post]
func (this *AuthViewset) EmailCodeValid(c *gin.Context) (err error) {
	// TODO	接口检测

	return this.SuccessBlackResponse(c)
}

// @Summary 重置密码
// @Description register by username and password
// @Accept  json
// @Produce  json
// @Param  user body model.UserParams true "username && password"
// @Success 200 {object} viewset.Response
// @Router /api/public/v1/auth/register/ [post]
func (this *AuthViewset) ResetPassword(c *gin.Context) (err error) {

	return this.SuccessBlackResponse(c)
}
