package auth

import (
	"go-sso/internal/registry"
	"go-sso/internal/repository/mysql/model"
	"go-sso/internal/service/apierror"
	"go-sso/pkg/json"
	"go-sso/util"
)

func GetUserByPassword(up model.UserParams) (*model.User, error) {
	u, r := registry.GetStorage().CheckUser(up.Account, up.Password)
	if !r {
		return nil, apierror.ErrAuth
	}

	return u, nil
}

func GetUserByTelephone(telephone string) (*model.User, error) {
	u, err := registry.GetStorage().GetUserByAccount(telephone)
	if err != nil {
		return nil, apierror.ErrAuth
	}

	return u, nil
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

func CheckTelephoneValid(telephone string) error {
	if !registry.GetStorage().IsValid(telephone, "telephone") {
		return apierror.ErrInvalid
	}
	if registry.GetStorage().Exists(telephone, "telephone") {
		return apierror.ErrInvalid
	}

	return nil
}

func CheckTelephoneExists(telephone string) error {
	if !registry.GetStorage().IsValid(telephone, "telephone") {
		return apierror.ErrInvalid
	}
	if !registry.GetStorage().Exists(telephone, "telephone") {
		return apierror.ErrInvalid
	}

	return nil
}
