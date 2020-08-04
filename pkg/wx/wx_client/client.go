package wx_client

import (
	"fmt"
	"github.com/pkg/errors"
	"go-sso/conf"
	"go-sso/pkg/json"
	"go-sso/pkg/log"
	"go-sso/pkg/request"
	"go-sso/util/encryption"
)

var (
	wx *wxClient
)

type wxClient struct {
	AppId  string
	Secret string
}

type LoginResponse struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    string `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

type UserInfo struct {
	AvatarURL string `json:"avatarUrl"`
	City      string `json:"city"`
	Country   string `json:"country"`
	Gender    int    `json:"gender"`
	NickName  string `json:"nickName"`
	OpenID    string `json:"openId"`
	Province  string `json:"province"`
	UnionID   string `json:"unionId"`
	Watermark struct {
		Appid     string `json:"appid"`
		Timestamp int    `json:"timestamp"`
	} `json:"watermark"`
}

func InitWeixin(config *conf.Config) {
	c := config.Weixin
	wx = &wxClient{
		AppId:  c.AppID,
		Secret: c.AppSecret,
	}
}

func GetWxClient() *wxClient {
	return wx
}

type LoginParams struct {
	Code string
}

func (this wxClient) Login(lp *LoginParams) (lr LoginResponse, err error) {
	url := "https://api.weixin.qq.com/sns/jscode2session"
	params := map[string]interface{}{
		"appid":      this.AppId,
		"secret":     this.Secret,
		"js_code":    lp.Code,
		"grant_type": "authorization_code",
	}
	fmt.Println(lp.Code)
	res, statusCode, _ := request.Get(url, nil, params)
	if statusCode != 200 {
		err = errors.New(fmt.Sprintf("login request error, status_code: %d", statusCode))
		log.Error(err)
		return
	}
	err = json.Unmarshal(res, &lr)
	if err != nil {
		return
	}
	if lr.ErrCode != "" {
		return lr, errors.New(lr.ErrMsg)
	}
	return
}

func (this *wxClient) getUserInfo(encryptedData, sessionKey, iv string) (*UserInfo, error) {
	s, err := decode(encryptedData, sessionKey, iv)
	if err != nil {
		return nil, err
	}
	u := &UserInfo{}
	err = json.Unmarshal([]byte(s), u)
	return u, err
}

func decode(encryptedData, sessionKey, iv string) (string, error) {
	iv, _ = encryption.Base64Decode(iv)
	aesKey, _ := encryption.Base64Decode(sessionKey)
	return encryption.Dncrypt(encryptedData, []byte(aesKey), []byte(iv))
}
