package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-sso/internal/apierror"
	"go-sso/internal/middlewares/auth"
	"go-sso/internal/model"
	"go-sso/internal/repository/cache"
	"go-sso/internal/repository/storage"
	"go-sso/internal/repository/storage/mysql"
	"go-sso/internal/service/viewset"
	"go-sso/pkg/email_tool"
	"go-sso/pkg/json"
	"go-sso/pkg/log"
	"go-sso/pkg/sms"
	"go-sso/registry"
	"go-sso/util"
	"net/http"
	"time"
)

const (
	oauthExpired  = time.Second * 300
	codeLength    = 12
	smsCodeLength = 6
)

type AuthViewset struct {
	viewset.ViewSet
	storage storage.Storage
	cache   cache.CacheClient
	sms     sms.ISms
}

func NewAuthViewset(storage storage.Storage) AuthViewset {
	return AuthViewset{
		storage: storage,
		cache:   registry.GetCacheStore(),
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
	var up mysql.UserParams
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
	res, err := driver.Login(c, u)
	if err != nil {
		return err
	}

	redirectUrl := c.Query("redirect_url")
	if redirectUrl == "" {
		return v.SuccessResponse(c, res)
	}

	url, err := GenerateOauthUrl(v.cache, redirectUrl, u)
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
	var tl mysql.TelephoneLoginParams
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
	res, err := driver.Login(c, u)
	if err != nil {
		return err
	}
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
	var rp mysql.RegisterParams
	var user mysql.User
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
func (v AuthViewset) CheckRegisterParams(rp *mysql.RegisterParams) map[string]string {
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
	telephone := c.Query("telephone")
	if err := v.checkTelephoneValid(telephone); err != nil {
		v.FailResponse(c, apierror.ErrInvalid, err)
		return err
	}
	code := util.RandomCode(4)
	if err := sms.SendLoginSms(v.sms, telephone, code); err != nil {
		return err
	}
	cacheStore := registry.GetCacheStore()
	if err := cacheStore.SetCache(telephone, code); err != nil {
		return err
	}
	return v.SuccessBlankResponse(c)
}

// @Summary telephone check
// @Description 手机验证码确认
func (v AuthViewset) VerifySmsCode(telephone, code string) (err error) {
	cacheStore := registry.GetCacheStore()
	rightCode, err := cacheStore.GetCache(telephone)
	if err != nil {
		return err
	}
	if rightCode != code {
		return apierror.NewParamsError("短信验证码错误")
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
	code := util.RandomCode(smsCodeLength)
	err = email_tool.SendEmailCode(code, email)
	if err != nil {
		return err
	}
	if err := cacheStore.SetCache(email, code); err != nil {
		return err
	}
	return v.SuccessBlankResponse(c)
}

func (v AuthViewset) VerifyEmailCode(email, code string) (err error) {
	cacheStore := registry.GetCacheStore()
	rightCode, err := cacheStore.GetCache(email)
	if err != nil {
		return err
	}
	if rightCode != code {
		return apierror.NewParamsError("邮件验证码错误")
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
	var rp mysql.ResetPasswordParams
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
	var cp mysql.ChangePasswordParams
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

// @Summary oauth
// @Description oauth check
// @Accept  json
// @Produce  json
// @Param  user body model.OauthRequest true "code"
// @Success 200 {object} model.OauthResponse
// @Router /api/v1/auth/oauth-check/ [post]
func (v AuthViewset) CheckOauthCode(c *gin.Context) (err error) {
	var r model.OauthRequest
	err = c.ShouldBind(&r)
	if err != nil {
		log.Error(err)
		return err
	}
	if r.Code == "" {
		return apierror.NewParamsError("无效的校验码")
	}
	user, err := GetUserByOauthCode(v.cache, r.Code)
	if err != nil {
		return err
	}

	return v.SuccessResponse(c, model.OauthResponse{
		UserId:   user.ID,
		Username: user.Username,
	})
}

func GenerateOauthUrl(c cache.CacheClient, url string, u *mysql.User) (string, error) {
	code, err := GenerateAuthCode(c, u)
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

func GenerateAuthCode(c cache.CacheClient, u *mysql.User) (string, error) {
	code := util.RandomCode(codeLength)
	b, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	if err := c.SetCacheExpired(code, string(b), oauthExpired); err != nil {
		return "", err
	}

	return code, nil
}

func GetUserByOauthCode(c cache.CacheClient, code string) (u *mysql.User, err error) {
	v, err := c.GetCache(code)
	if err != nil {
		return nil, err
	}
	if v == "" {
		return nil, apierror.NewParamsError("无效的校验码")
	}
	u = &mysql.User{}
	if err := json.Unmarshal([]byte(v), &u); err != nil {
		return nil, err
	}

	return u, nil
}

func (v AuthViewset) getUserByPassword(up mysql.UserParams) (*mysql.User, error) {
	u, r := v.storage.GetUser(up.Account, up.Password)
	if !r {
		return nil, apierror.ErrAuth
	}

	return u, nil
}

func (v AuthViewset) getUserByTelephone(telephone string) (*mysql.User, error) {
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
