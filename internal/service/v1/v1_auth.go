package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-sso/internal/apierror"
	"go-sso/internal/middlewares/auth"
	"go-sso/internal/repository"
	"go-sso/internal/repository/mysql/model"
	"go-sso/internal/service/viewset"
	"go-sso/pkg/email_tool"
	"go-sso/pkg/json"
	"go-sso/pkg/log"
	"go-sso/registry"
	"go-sso/util"
	"net/http"
)

type AuthViewset struct {
	viewset.ViewSet
	storage repository.Storage
}

func NewAuthViewset(storage repository.Storage) AuthViewset {
	return AuthViewset{
		storage: storage,
	}
}

func (v AuthViewset) ErrorHandler(f func(c *gin.Context) error) func(c *gin.Context) {
	return func(c *gin.Context) {
		v.ViewSet.ErrorHandler(f, c)
	}
}

// @Summary user login
// @Description 1.账号密码登录 2.手机号，邮箱登录
// @Accept  json
// @Produce  json
// @Param  user body model.UserParams true "username && password"
// @Success 200 {object} viewset.Response
// @Router /api/public/v1/auth/login/ [post]
func (v AuthViewset) Login(c *gin.Context) (err error) {
	var up model.UserParams
	err = c.ShouldBind(&up)
	if err != nil {
		log.Error(err.Error())
		v.FailResponse(c, apierror.ErrInvalid)
		return apierror.ErrInvalid
	}
	u, err := v.getUserByPassword(up)
	if err != nil {
		return err
	}
	// 登录方式 token
	driver := auth.GenerateAuthDriver(auth.TokenAuth)
	res := driver.Login(c, u)

	redirectUrl := c.Query("redirect_url")
	if redirectUrl == "" {
		return v.SuccessResponse(c, res)
	}

	url, err := GenerateOauthUrl(redirectUrl, u)
	if err != nil {
		return err
	}
	c.Redirect(http.StatusMovedPermanently, url)
	return nil
}

// @Summary telephone login
// @Description telephone login steps 1.IsTelephoneExist 2. SendSmsCode 3. TelephoneLogin
// @Accept  json
// @Produce  json
// @Param  user body model.UserParams true "username && password"
// @Success 200 {object} viewset.Response
// @Router /api/public/v1/auth/telephone/login/ [post]
func (v AuthViewset) TelephoneLogin(c *gin.Context) (err error) {
	var tl model.TelephoneLoginParams
	if err = c.ShouldBind(&tl); err != nil {
		return err
	}
	err = v.VerifySmsCode(tl.Telephone, tl.Code)
	if err != nil {
		return
	}
	// 登录方式 token
	u, err := v.getUserByTelephone(tl.Telephone)
	if err != nil {
		return
	}
	driver := auth.GenerateAuthDriver(auth.TokenAuth)
	res := driver.Login(c, u)
	return v.SuccessResponse(c, res)
}

// @Summary user register
// @Description register by username, telephone and password
// @Accept  json
// @Produce  json
// @Param  user body model.UserParams true "username && password"
// @Success 200 {object} viewset.Response
// @Router /api/public/v1/auth/register/ [post]
func (v AuthViewset) Register(c *gin.Context) (err error) {
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
	errs := v.CheckRegisterParams(&rp)
	if len(errs) > 0 {
		v.FailResponse(c, apierror.ErrInvalid, errs)
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
	if _, err = registry.GetStorage().Create(&user); err != nil {
		log.Error(err.Error())
		return err
	} else {
		return v.SuccessResponse(c, "注册成功")
	}
}

// 检测注册用户参数
func (v AuthViewset) CheckRegisterParams(rp *model.RegisterParams) map[string]string {
	errs := make(map[string]string)
	// 检查参数是否合法
	if !registry.GetStorage().IsValid(rp.Username, "username") {
		errs["username"] = "用户名至少3位以上字母开头"
	}
	if !registry.GetStorage().IsValid(rp.Telephone, "telephone") {
		errs["telephone"] = "手机号格式错误"
	}
	if rp.Email != "" && !registry.GetStorage().IsValid(rp.Email, "email") {
		errs["email"] = "email格式错误"
	}
	if len(errs) > 0 {
		return errs
	}
	// 检查是否重复注册
	if registry.GetStorage().Exists(rp.Username, "username") {
		errs["username"] = "用户已经存在"
	}
	if registry.GetStorage().Exists(rp.Telephone, "telephone") {
		errs["telephone"] = "手机号已经存在"
	}
	if rp.Email != "" && registry.GetStorage().Exists(rp.Email, "email") {
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
func (v AuthViewset) CheckTelephoneValid(c *gin.Context) (err error) {
	// 检查参数是否合法
	telephone := c.Query("telephone")
	if err := v.checkTelephoneValid(telephone); err != nil {
		v.FailResponse(c, apierror.ErrInvalid, err)
		return err
	}

	return v.SuccessBlankResponse(c)
}

// @Summary check telephone
// @Description check telephone whether exist
// @Accept  json
// @Produce  json
// @Param  user body  true "telephone"
// @Success 200 {object} viewset.Response
// @Router /api/public/v1/auth/check-telephone-exist/ [post]
func (v AuthViewset) CheckTelephoneExist(c *gin.Context) (err error) {
	// 检查参数是否合法
	telephone := c.Query("telephone")
	if err := v.checkTelephoneExists(telephone); err != nil {
		v.FailResponse(c, apierror.ErrInvalid, err)
		return err
	}

	return v.SuccessBlankResponse(c)
}

// @Summary 发送手机验证码
// @Description send telephone verify code
// @Accept  json
// @Produce  json
// @Param  user body model.UserParams true "username && password"
// @Success 200 {object} viewset.Response
// @Router /api/public/v1/auth/register/ [post]
func (v AuthViewset) SendSmsCode(c *gin.Context) (err error) {
	// TODO 发送短信
	code := "123456"
	telephone := c.Query("telephone")
	if err := v.checkTelephoneValid(telephone); err != nil {
		v.FailResponse(c, apierror.ErrInvalid, err)
		return err
	}
	cacheStore := registry.GetCacheStore()
	cacheStore.SetCache(telephone, code)
	return v.SuccessBlankResponse(c)
}

// @Summary telephone check
// @Description 手机验证码确认
func (v AuthViewset) VerifySmsCode(telephone, code string) (err error) {
	cacheStore := registry.GetCacheStore()
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
func (v AuthViewset) SendEmailCode(c *gin.Context) (err error) {
	email := c.Query("email")
	if ok := registry.GetStorage().IsValid(email, "email"); !ok {
		return apierror.ErrInvalid
	}
	cacheStore := registry.GetCacheStore()
	code := util.RandomCode()
	err = email_tool.SendEmailCode(code, email)
	if err != nil {
		return err
	}
	cacheStore.SetCache(email, code)
	return v.SuccessResponse(c, gin.H{"url": ""})
}

func (v AuthViewset) VerifyEmailCode(email, code string) (err error) {
	cacheStore := registry.GetCacheStore()
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
func (v AuthViewset) ResetPassword(c *gin.Context) (err error) {
	var rp model.ResetPasswordParams
	err = c.ShouldBind(&rp)
	if err != nil {
		return
	}
	user, err := registry.GetStorage().GetUserByAccount(rp.Account)
	if err != nil {
		return
	}
	switch rp.VerifyType {
	case "email":
		err = v.VerifyEmailCode(user.Email, rp.Code)
	case "telephone":
		err = v.VerifySmsCode(user.Telephone, rp.Code)
	default:
		err = errors.New("验证类型错误")
	}
	if err != nil {
		return
	}
	err = registry.GetStorage().ChangePassword(user, rp.NewPassword)
	if err != nil {
		return
	}
	return v.SuccessBlankResponse(c)
}

// @Summary 修改密码
// @Description change password by username and password
// @Accept  json
// @Produce  json
// @Param  user body model.ChangePasswordParams true "raw_password && new_password"
// @Success 200 {object} viewset.Response
// @Router /api/v1/auth/change-password/ [post]
func (v AuthViewset) ChangePassword(c *gin.Context) (err error) {
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
	username := auth.GetCurrentUser(c).Username
	if u, ok := registry.GetStorage().GetUser(username, cp.RawPassword); ok {
		err = registry.GetStorage().ChangePassword(u, cp.NewPassword)
		if err != nil {
			log.Error(err)
			return
		}
		return v.SuccessBlankResponse(c)
	}
	return errors.New("原密码错误")
}

func GenerateOauthUrl(url string, u *model.User) (string, error) {
	code, err := GenerateAuthCode(u)
	if err != nil {
		return "", apierror.ErrInternal
	}
	m := make(map[string]interface{})
	m["code"] = code
	url, err = util.BuildUrlQuery(url, m)
	if err != nil {
		return "", err
	}

	return url, err
}

func GenerateAuthCode(u *model.User) (string, error) {
	code := util.RandomAuthCode()
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	registry.GetCacheStore().SetCache(code, string(b))

	return code, nil
}

func CheckAuthCode(code string) (*model.User, error) {
	cache, exist := registry.GetCacheStore().GetCache(code)
	if !exist {
		return nil, apierror.ErrAuth
	}
	v, ok := cache.(string)
	if !ok {
		registry.GetCacheStore().RemoveCache(code)
		return nil, apierror.ErrAuth
	}
	var u model.User
	if err := json.Unmarshal([]byte(v), &u); err != nil {
		return nil, apierror.ErrAuth
	}

	return &u, nil
}

func (v AuthViewset) getUserByPassword(up model.UserParams) (*model.User, error) {
	u, r := v.storage.GetUser(up.Account, up.Password)
	if !r {
		return nil, apierror.ErrAuth
	}

	return u, nil
}

func (v AuthViewset) getUserByTelephone(telephone string) (*model.User, error) {
	u, err := v.storage.GetUserByAccount(telephone)
	if err != nil {
		return nil, apierror.ErrAuth
	}

	return u, nil
}

func (v AuthViewset) checkTelephoneValid(telephone string) error {
	if !registry.GetStorage().IsValid(telephone, "telephone") {
		return apierror.ErrInvalid
	}
	if registry.GetStorage().Exists(telephone, "telephone") {
		return apierror.ErrInvalid
	}

	return nil
}

func (v AuthViewset) checkTelephoneExists(telephone string) error {
	if !v.storage.IsValid(telephone, "telephone") {
		return apierror.ErrInvalid
	}
	if !v.storage.Exists(telephone, "telephone") {
		return apierror.ErrInvalid
	}

	return nil
}
