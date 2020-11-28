package auth

import (
	"go-sso/db/inter"
	"go-sso/db/model"
	"go-sso/pkg/json"
	"go-sso/registry"
	"go-sso/service/api/api_error"
	"go-sso/util"
)

func GetUserByPassword(up model.UserParams) (*model.User, error) {
	u, r := inter.GetQuery().CheckUser(up.Account, up.Password)
	if !r {
		return nil, api_error.ErrAuth
	}

	return u, nil
}

func GetUserByTelephone(telephone string) (*model.User, error) {
	u, err := inter.GetQuery().GetUserByAccount(telephone)
	if err != nil {
		return nil, api_error.ErrAuth
	}

	return u, nil
}

func GenerateOauthUrl(url string, u *model.User) (string, error) {
	code, err := GenerateAuthCode(u)
	if err != nil {
		return "", api_error.ErrInternal
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
		return nil, api_error.ErrAuth
	}
	v, ok := cache.(string)
	if !ok {
		registry.GetCacheStore().RemoveCache(code)
		return nil, api_error.ErrAuth
	}
	var u model.User
	if err := json.Unmarshal([]byte(v), &u); err != nil {
		return nil, api_error.ErrAuth
	}

	return &u, nil
}

func CheckTelephoneValid(telephone string) error {
	if !inter.GetQuery().IsValid(telephone, "telephone") {
		return api_error.ErrInvalid
	}
	if inter.GetQuery().Exists(telephone, "telephone") {
		return api_error.ErrInvalid
	}

	return nil
}

func CheckTelephoneExists(telephone string) error {
	if !inter.GetQuery().IsValid(telephone, "telephone") {
		return api_error.ErrInvalid
	}
	if !inter.GetQuery().Exists(telephone, "telephone") {
		return api_error.ErrInvalid
	}

	return nil
}
