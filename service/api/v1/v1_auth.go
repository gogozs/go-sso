package v1

import (
	"github.com/gin-gonic/gin"
	"go-qiuplus/db/inter"
	"go-qiuplus/db/model"
	"go-qiuplus/pkg/api_error"
	"go-qiuplus/pkg/log"
	"go-qiuplus/service/api/viewset"
	"go-qiuplus/util"
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

}

// @Summary user register
// @Description register by username and password
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
		this.FailResponse(c, api_error.ErrInvalid, errs)
		return
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

// @Summary 检测用户是否合法 手机，邮箱是否被占用
// @Description register by username and password
// @Accept  json
// @Produce  json
// @Param  user body model.UserParams true "username && password"
// @Success 200 {object} viewset.Response
// @Router /api/public/v1/auth/register/ [post]
func (this *AuthViewset) UserCheck(c *gin.Context) (err error) {

}

// @Summary 发送手机验证码
// @Description register by username and password
// @Accept  json
// @Produce  json
// @Param  user body model.UserParams true "username && password"
// @Success 200 {object} viewset.Response
// @Router /api/public/v1/auth/register/ [post]
func (this *AuthViewset) SendSmsCode(c *gin.Context) (err error) {

}

// @Summary 发送邮箱验证码
// @Description register by username and password
// @Accept  json
// @Produce  json
// @Param  user body model.UserParams true "username && password"
// @Success 200 {object} viewset.Response
// @Router /api/public/v1/auth/register/ [post]
func (this *AuthViewset) SendEmailCode(c *gin.Context) (err error) {

}

// @Summary 重置密码
// @Description register by username and password
// @Accept  json
// @Produce  json
// @Param  user body model.UserParams true "username && password"
// @Success 200 {object} viewset.Response
// @Router /api/public/v1/auth/register/ [post]
func (this *AuthViewset) ResetPassword(c *gin.Context) (err error) {

}
